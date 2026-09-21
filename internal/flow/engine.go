// Package flow is Cashtro's owned n8n.
//
// A Workflow is a node graph we store and execute. Catch-hooks live on
// this kernel. Propose/critique seats are internal/think. There is no
// N8N_WEBHOOK_URL, no MOONSHOT_API_KEY, and no GLM_API_KEY.
package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cashtro/cashtro/internal/think"
)

// Kind is a node we own.
type Kind string

const (
	KindWebhook  Kind = "webhook"
	KindManual   Kind = "manual"
	KindSet      Kind = "set"
	KindIf       Kind = "if"
	KindInvoke   Kind = "invoke"
	KindPropose  Kind = "propose"
	KindCritique Kind = "critique"
	KindDual     Kind = "dual"
	KindMemory   Kind = "memory"
	KindNote     Kind = "note"
	KindConfirm  Kind = "confirm"
	KindShip     Kind = "ship"
)

// Node is one box on the graph.
type Node struct {
	ID     string            `json:"id"`
	Kind   Kind              `json:"kind"`
	Name   string            `json:"name"`
	Params map[string]string `json:"params,omitempty"`
}

// Edge connects two nodes. Port is "" | "true" | "false".
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Port string `json:"port,omitempty"`
}

// Workflow is one owned automation.
type Workflow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Summary   string    `json:"summary"`
	Trigger   Kind      `json:"trigger"`
	Nodes     []Node    `json:"nodes"`
	Edges     []Edge    `json:"edges"`
	Enabled   bool      `json:"enabled"`
	Hook      string    `json:"hook,omitempty"`
	Fires     int       `json:"fires"`
	LastOK    bool      `json:"lastOk"`
	LastMsg   string    `json:"lastMsg,omitempty"`
	LastScore int       `json:"lastScore,omitempty"`
	LastAt    time.Time `json:"lastAt,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Create is the input for a new workflow.
type Create struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Trigger Kind   `json:"trigger"`
	Nodes   []Node `json:"nodes"`
	Edges   []Edge `json:"edges"`
}

// NodeLog is one node's result inside a run.
type NodeLog struct {
	ID      string `json:"id"`
	Kind    Kind   `json:"kind"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Skipped bool   `json:"skipped,omitempty"`
}

// Run is one firing of a workflow.
type Run struct {
	ID       int               `json:"id"`
	Workflow string            `json:"workflow"`
	Name     string            `json:"name"`
	At       time.Time         `json:"at"`
	OK       bool              `json:"ok"`
	Message  string            `json:"message"`
	Score    int               `json:"score,omitempty"`
	Owner    string            `json:"owner"`
	Bound    bool              `json:"bound"`
	Input    map[string]string `json:"input,omitempty"`
	Item     map[string]string `json:"item,omitempty"`
	Nodes    []NodeLog         `json:"nodes"`
}

// Invoker runs a kernel verb from an invoke/memory/note/confirm/ship node.
type Invoker func(ctx context.Context, cap string, fields map[string]string) (ok bool, message string, err error)

var (
	// ErrNotFound is returned when a workflow id is unknown.
	ErrNotFound = errors.New("workflow not found")
	// ErrInvalid is returned for a graph that cannot be created or fired.
	ErrInvalid = errors.New("invalid workflow")
	// ErrDisabled is returned when a workflow is switched off.
	ErrDisabled = errors.New("workflow disabled")
	// ErrBusy is returned when the same workflow is already mid-fire.
	ErrBusy = errors.New("workflow busy")
)

const (
	maxRuns  = 80
	maxSteps = 32
	maxFlows = 48
)

var slot = regexp.MustCompile(`\{\{([a-zA-Z0-9_.$]+)\}\}`)

// Engine is an in-memory workflow table.
type Engine struct {
	mu     sync.RWMutex
	seq    int
	runSeq int
	now    func() time.Time
	flows  map[string]*Workflow
	runs   []Run
	busy   map[string]bool
}

// Option configures an engine.
type Option func(*Engine)

// WithClock injects a clock for tests.
func WithClock(now func() time.Time) Option {
	return func(e *Engine) {
		e.now = now
	}
}

// New returns an engine seeded with owned workflows. Zero vendor keys.
func New(opts ...Option) *Engine {
	e := &Engine{
		now:   func() time.Time { return time.Now().UTC() },
		flows: make(map[string]*Workflow),
		busy:  make(map[string]bool),
	}
	for _, opt := range opts {
		opt(e)
	}
	e.seed()
	return e
}

// List returns workflows ordered by name.
func (e *Engine) List() []Workflow {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Workflow, 0, len(e.flows))
	for _, f := range e.flows {
		out = append(out, cloneWorkflow(f))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Get returns one workflow by id.
func (e *Engine) Get(id string) (Workflow, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	f, ok := e.flows[id]
	if !ok {
		return Workflow{}, ErrNotFound
	}
	return cloneWorkflow(f), nil
}

// Create adds a workflow. Default trigger is manual.
func (e *Engine) Create(in Create) (Workflow, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Workflow{}, fmt.Errorf("%w: name is required", ErrInvalid)
	}
	nodes, err := cleanNodes(in.Nodes)
	if err != nil {
		return Workflow{}, err
	}
	if len(nodes) == 0 {
		return Workflow{}, fmt.Errorf("%w: at least one node", ErrInvalid)
	}
	trigger := in.Trigger
	if trigger == "" {
		trigger = KindManual
	}
	switch trigger {
	case KindManual, KindWebhook:
	default:
		return Workflow{}, fmt.Errorf("%w: unknown trigger %q", ErrInvalid, trigger)
	}
	edges, err := cleanEdges(in.Edges, nodes)
	if err != nil {
		return Workflow{}, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.flows) >= maxFlows {
		return Workflow{}, fmt.Errorf("%w: table full", ErrInvalid)
	}
	id := e.uniqueID(name)
	now := e.now()
	f := &Workflow{
		ID:        id,
		Name:      name,
		Summary:   strings.TrimSpace(in.Summary),
		Trigger:   trigger,
		Nodes:     nodes,
		Edges:     edges,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if trigger == KindWebhook {
		f.Hook = "/api/n8n/" + id + "/hook"
	}
	e.flows[id] = f
	return cloneWorkflow(f), nil
}

// Toggle flips enabled.
func (e *Engine) Toggle(id string) (Workflow, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	f, ok := e.flows[id]
	if !ok {
		return Workflow{}, ErrNotFound
	}
	f.Enabled = !f.Enabled
	f.UpdatedAt = e.now()
	return cloneWorkflow(f), nil
}

// Runs returns recent firings, newest first.
func (e *Engine) Runs() []Run {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Run, len(e.runs))
	copy(out, e.runs)
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out
}

// Catalog is the public node catalog — what we own, not n8n's app store.
func Catalog() []map[string]string {
	return []map[string]string{
		{"kind": string(KindWebhook), "role": "trigger", "summary": "Catch-hook we own. POST /api/n8n/{id}/hook."},
		{"kind": string(KindManual), "role": "trigger", "summary": "Desk or POST /api/agents/n8n/workflow."},
		{"kind": string(KindSet), "role": "transform", "summary": "Write fields onto the item."},
		{"kind": string(KindIf), "role": "branch", "summary": "true/false ports on a field match."},
		{"kind": string(KindInvoke), "role": "verb", "summary": "Call a kernel capability."},
		{"kind": string(KindPropose), "role": "think", "summary": "Owned proposer. Replaces Moonshot."},
		{"kind": string(KindCritique), "role": "think", "summary": "Owned critic. Replaces GLM."},
		{"kind": string(KindDual), "role": "think", "summary": "Propose then critique in-process."},
		{"kind": string(KindMemory), "role": "store", "summary": "memory.store on the kernel."},
		{"kind": string(KindNote), "role": "store", "summary": "research.ingest on the kernel."},
		{"kind": string(KindConfirm), "role": "gate", "summary": "comms.send parks a human confirm."},
		{"kind": string(KindShip), "role": "ship", "summary": "delivery.create parks a mandate in idea."},
	}
}

// Enter marks a workflow busy.
func (e *Engine) Enter(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	f, ok := e.flows[id]
	if !ok {
		return ErrNotFound
	}
	if !f.Enabled {
		return ErrDisabled
	}
	if e.busy[id] {
		return ErrBusy
	}
	e.busy[id] = true
	return nil
}

// Leave clears the busy flag.
func (e *Engine) Leave(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.busy, id)
}

// Execute walks the graph. Dual seats are local.
func (e *Engine) Execute(ctx context.Context, id string, input map[string]string, invoker Invoker) (Run, error) {
	wf, err := e.Get(id)
	if err != nil {
		return Run{}, err
	}
	if !wf.Enabled {
		return Run{}, ErrDisabled
	}
	if input == nil {
		input = map[string]string{}
	}
	item := copyMap(input)
	item["owner"] = think.Owner
	item["bound"] = "false"

	byID := make(map[string]Node, len(wf.Nodes))
	for _, n := range wf.Nodes {
		byID[n.ID] = n
	}
	outgoing := make(map[string][]Edge)
	for _, edge := range wf.Edges {
		outgoing[edge.From] = append(outgoing[edge.From], edge)
	}

	start := startNode(wf)
	if start.ID == "" {
		return Run{}, fmt.Errorf("%w: no start node", ErrInvalid)
	}

	logs := make([]NodeLog, 0, len(wf.Nodes))
	ok := true
	msg := "fired " + wf.Name
	score := 0
	queue := []string{start.ID}
	seen := map[string]int{}

	for len(queue) > 0 {
		if len(logs) >= maxSteps {
			ok = false
			msg = "max steps"
			break
		}
		nid := queue[0]
		queue = queue[1:]
		seen[nid]++
		if seen[nid] > 2 {
			ok = false
			msg = "cycle"
			break
		}
		node, exists := byID[nid]
		if !exists {
			ok = false
			msg = "unknown node " + nid
			break
		}
		log, nextPort, err := e.runNode(ctx, node, item, invoker)
		logs = append(logs, log)
		if err != nil {
			ok = false
			msg = err.Error()
			break
		}
		if !log.OK {
			ok = false
			msg = log.Message
			break
		}
		if v, err := strconv.Atoi(item["score"]); err == nil {
			score = v
		}
		for _, edge := range outgoing[nid] {
			if edge.Port != "" && edge.Port != nextPort {
				continue
			}
			queue = append(queue, edge.To)
		}
	}

	return e.finish(wf.ID, ok, msg, score, input, item, logs)
}

func (e *Engine) runNode(ctx context.Context, node Node, item map[string]string, invoker Invoker) (NodeLog, string, error) {
	params := Render(node.Params, item)
	log := NodeLog{ID: node.ID, Kind: node.Kind, OK: true, Message: node.Name}
	port := ""

	switch node.Kind {
	case KindManual, KindWebhook:
		log.Message = string(node.Kind) + " catch"
		return log, port, nil
	case KindSet:
		for k, v := range params {
			if k == "" {
				continue
			}
			item[k] = v
		}
		log.Message = "set " + strconv.Itoa(len(params)) + " fields"
		return log, port, nil
	case KindIf:
		field := params["field"]
		if field == "" {
			field = "prompt"
		}
		want := params["equals"]
		got := item[field]
		match := got == want
		if want == "" {
			match = strings.TrimSpace(got) != ""
		}
		if contains := params["contains"]; contains != "" {
			match = strings.Contains(strings.ToLower(got), strings.ToLower(contains))
		}
		if match {
			port = "true"
			log.Message = "if true"
		} else {
			port = "false"
			log.Message = "if false"
		}
		item["branch"] = port
		return log, port, nil
	case KindPropose:
		d := think.Local(firstNonEmpty(params["prompt"], item["prompt"]))
		item["propose"] = d.Propose.Content
		item["prompt"] = firstNonEmpty(item["prompt"], d.Propose.Content)
		log.Message = "propose · owned"
		return log, port, nil
	case KindCritique:
		src := firstNonEmpty(params["prompt"], item["propose"], item["prompt"])
		d := think.Local(src)
		item["critique"] = d.Critique.Content
		item["score"] = strconv.Itoa(d.Score)
		log.Message = "critique · score " + item["score"]
		return log, port, nil
	case KindDual:
		d := think.Local(firstNonEmpty(params["prompt"], item["prompt"]))
		item["propose"] = d.Propose.Content
		item["critique"] = d.Critique.Content
		item["score"] = strconv.Itoa(d.Score)
		item["owner"] = d.Owner
		item["bound"] = "false"
		log.Message = d.Message
		return log, port, nil
	case KindInvoke, KindMemory, KindNote, KindConfirm, KindShip:
		if invoker == nil {
			log.OK = false
			log.Message = "no invoker"
			return log, port, nil
		}
		cap := params["cap"]
		if cap == "" {
			cap = defaultCap(node.Kind)
		}
		fields := copyMap(params)
		delete(fields, "cap")
		if fields["prompt"] == "" {
			fields["prompt"] = item["prompt"]
		}
		if fields["text"] == "" {
			fields["text"] = firstNonEmpty(item["propose"], item["critique"], item["prompt"])
		}
		if fields["claim"] == "" {
			fields["claim"] = firstNonEmpty(item["propose"], item["prompt"])
		}
		if fields["body"] == "" {
			fields["body"] = firstNonEmpty(item["propose"], item["prompt"])
		}
		if fields["name"] == "" {
			fields["name"] = firstNonEmpty(params["name"], item["name"], item["prompt"])
		}
		if fields["topic"] == "" {
			fields["topic"] = "n8n"
		}
		if fields["source"] == "" {
			fields["source"] = "n8n"
		}
		ok, message, err := invoker(ctx, cap, fields)
		if err != nil {
			log.OK = false
			log.Message = err.Error()
			return log, port, nil
		}
		log.OK = ok
		log.Message = message
		if !ok {
			return log, port, nil
		}
		return log, port, nil
	default:
		log.OK = false
		log.Message = "unknown kind " + string(node.Kind)
		return log, port, fmt.Errorf("%w: %s", ErrInvalid, node.Kind)
	}
}

func (e *Engine) finish(id string, ok bool, msg string, score int, input, item map[string]string, logs []NodeLog) (Run, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	f, exists := e.flows[id]
	if !exists {
		return Run{}, ErrNotFound
	}
	e.runSeq++
	now := e.now()
	run := Run{
		ID:       e.runSeq,
		Workflow: f.ID,
		Name:     f.Name,
		At:       now,
		OK:       ok,
		Message:  msg,
		Score:    score,
		Owner:    think.Owner,
		Bound:    false,
		Input:    copyMap(input),
		Item:     copyMap(item),
		Nodes:    logs,
	}
	e.runs = append(e.runs, run)
	if len(e.runs) > maxRuns {
		e.runs = append([]Run(nil), e.runs[len(e.runs)-maxRuns:]...)
	}
	f.Fires++
	f.LastOK = ok
	f.LastMsg = msg
	f.LastScore = score
	f.LastAt = now
	f.UpdatedAt = now
	return run, nil
}

func (e *Engine) seed() {
	now := e.now()
	put := func(f Workflow) {
		f.CreatedAt = now
		f.UpdatedAt = now
		if f.Trigger == KindWebhook {
			f.Hook = "/api/n8n/" + f.ID + "/hook"
		}
		e.flows[f.ID] = &f
	}

	put(Workflow{
		ID:      "wealth-dual",
		Name:    "Owned wealth dual",
		Summary: "Propose then critique in-process. Replaces Moonshot + GLM. No vendor keys.",
		Trigger: KindManual,
		Enabled: true,
		Nodes: []Node{
			{ID: "start", Kind: KindManual, Name: "Desk fire"},
			{ID: "dual", Kind: KindDual, Name: "Owned dual"},
			{ID: "remember", Kind: KindMemory, Name: "Remember the move", Params: map[string]string{"topic": "forge", "text": "{{propose}}"}},
			{ID: "note", Kind: KindNote, Name: "Park the finding", Params: map[string]string{"claim": "{{propose}}", "source": "n8n:wealth-dual"}},
		},
		Edges: []Edge{
			{From: "start", To: "dual"},
			{From: "dual", To: "remember"},
			{From: "remember", To: "note"},
		},
	})

	put(Workflow{
		ID:      "intake-hook",
		Name:    "Owned intake hook",
		Summary: "Catch-hook we own. Replaces N8N_WEBHOOK_URL.",
		Trigger: KindWebhook,
		Enabled: true,
		Nodes: []Node{
			{ID: "hook", Kind: KindWebhook, Name: "Catch"},
			{ID: "set", Kind: KindSet, Name: "Stamp owner", Params: map[string]string{"lane": "intake", "owner": "cashtro"}},
			{ID: "gate", Kind: KindIf, Name: "Has a body?", Params: map[string]string{"field": "prompt"}},
			{ID: "remember", Kind: KindMemory, Name: "Remember intake", Params: map[string]string{"topic": "intake", "text": "{{prompt}}"}},
			{ID: "confirm", Kind: KindConfirm, Name: "Park human confirm", Params: map[string]string{"body": "{{prompt}}"}},
			{ID: "skip", Kind: KindSet, Name: "Empty catch", Params: map[string]string{"skipped": "true"}},
		},
		Edges: []Edge{
			{From: "hook", To: "set"},
			{From: "set", To: "gate"},
			{From: "gate", To: "remember", Port: "true"},
			{From: "remember", To: "confirm"},
			{From: "gate", To: "skip", Port: "false"},
		},
	})

	put(Workflow{
		ID:      "research-line",
		Name:    "Owned research line",
		Summary: "Manual ingest → note → memory. Stays on the kernel.",
		Trigger: KindManual,
		Enabled: true,
		Nodes: []Node{
			{ID: "start", Kind: KindManual, Name: "Desk fire"},
			{ID: "note", Kind: KindNote, Name: "Ingest", Params: map[string]string{"claim": "{{prompt}}", "source": "n8n:research-line"}},
			{ID: "remember", Kind: KindMemory, Name: "Remember", Params: map[string]string{"topic": "research", "text": "{{prompt}}"}},
		},
		Edges: []Edge{
			{From: "start", To: "note"},
			{From: "note", To: "remember"},
		},
	})
}

func (e *Engine) uniqueID(name string) string {
	base := slugify(name)
	if base == "" {
		e.seq++
		return fmt.Sprintf("wf-%d", e.seq)
	}
	if _, exists := e.flows[base]; !exists {
		return base
	}
	for n := 2; ; n++ {
		id := fmt.Sprintf("%s-%d", base, n)
		if _, exists := e.flows[id]; !exists {
			return id
		}
	}
}

func startNode(wf Workflow) Node {
	incoming := map[string]int{}
	for _, e := range wf.Edges {
		incoming[e.To]++
	}
	for _, n := range wf.Nodes {
		if incoming[n.ID] == 0 {
			return n
		}
	}
	if len(wf.Nodes) > 0 {
		return wf.Nodes[0]
	}
	return Node{}
}

func defaultCap(k Kind) string {
	switch k {
	case KindMemory:
		return "memory.store"
	case KindNote:
		return "research.ingest"
	case KindConfirm:
		return "comms.send"
	case KindShip:
		return "delivery.create"
	default:
		return ""
	}
}

func cleanNodes(in []Node) ([]Node, error) {
	out := make([]Node, 0, len(in))
	seen := map[string]struct{}{}
	for _, n := range in {
		n.ID = slugify(n.ID)
		n.Name = strings.TrimSpace(n.Name)
		if n.ID == "" {
			return nil, fmt.Errorf("%w: node id required", ErrInvalid)
		}
		if n.Name == "" {
			n.Name = n.ID
		}
		switch n.Kind {
		case KindWebhook, KindManual, KindSet, KindIf, KindInvoke, KindPropose, KindCritique, KindDual, KindMemory, KindNote, KindConfirm, KindShip:
		default:
			return nil, fmt.Errorf("%w: unknown kind %q", ErrInvalid, n.Kind)
		}
		if _, ok := seen[n.ID]; ok {
			return nil, fmt.Errorf("%w: duplicate node %s", ErrInvalid, n.ID)
		}
		seen[n.ID] = struct{}{}
		n.Params = cleanParams(n.Params)
		out = append(out, n)
	}
	return out, nil
}

func cleanEdges(in []Edge, nodes []Node) ([]Edge, error) {
	known := map[string]struct{}{}
	for _, n := range nodes {
		known[n.ID] = struct{}{}
	}
	out := make([]Edge, 0, len(in))
	for _, e := range in {
		e.From = slugify(e.From)
		e.To = slugify(e.To)
		e.Port = strings.TrimSpace(e.Port)
		if e.From == "" || e.To == "" {
			return nil, fmt.Errorf("%w: edge needs from and to", ErrInvalid)
		}
		if _, ok := known[e.From]; !ok {
			return nil, fmt.Errorf("%w: unknown from %s", ErrInvalid, e.From)
		}
		if _, ok := known[e.To]; !ok {
			return nil, fmt.Errorf("%w: unknown to %s", ErrInvalid, e.To)
		}
		if e.Port != "" && e.Port != "true" && e.Port != "false" {
			return nil, fmt.Errorf("%w: port %q", ErrInvalid, e.Port)
		}
		out = append(out, e)
	}
	return out, nil
}

func cleanParams(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// FlattenJSON turns a payload into string fields the way n8n $json does.
func FlattenJSON(raw []byte) map[string]string {
	out := map[string]string{}
	if len(raw) == 0 {
		return out
	}
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		out["prompt"] = strings.TrimSpace(string(raw))
		out["body"] = out["prompt"]
		return out
	}
	for k, v := range obj {
		switch t := v.(type) {
		case string:
			out[k] = strings.TrimSpace(t)
		case float64:
			out[k] = strconv.FormatFloat(t, 'f', -1, 64)
		case bool:
			out[k] = strconv.FormatBool(t)
		default:
			b, _ := json.Marshal(t)
			out[k] = string(b)
		}
	}
	return out
}

// Render replaces {{field}} slots from the current item.
func Render(fields map[string]string, item map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = slot.ReplaceAllStringFunc(v, func(m string) string {
			key := slot.FindStringSubmatch(m)[1]
			key = strings.TrimPrefix(key, "$json.")
			key = strings.TrimPrefix(key, "item.")
			if got, ok := item[key]; ok {
				return got
			}
			return ""
		})
	}
	return out
}

func cloneWorkflow(f *Workflow) Workflow {
	out := *f
	out.Nodes = append([]Node(nil), f.Nodes...)
	for i := range out.Nodes {
		out.Nodes[i].Params = copyMap(f.Nodes[i].Params)
	}
	out.Edges = append([]Edge(nil), f.Edges...)
	return out
}

func copyMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		default:
			if b.Len() > 0 && !dash {
				b.WriteByte('-')
				dash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
