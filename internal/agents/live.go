package agents

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
)

func explorerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	q := strings.ToLower(payloadQuery(call, "query"))
	hits := make([]map[string]string, 0)
	for _, p := range k.Processes() {
		blob := strings.ToLower(p.Spec.Name + " " + p.Spec.Role + " " + p.Spec.Summary + " " + strings.Join(p.Spec.Capabilities, " "))
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "agent", "id": p.Spec.ID, "name": p.Spec.Name})
		}
	}
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			blob := strings.ToLower(s.Name + " " + s.Client + " " + s.Notes + " " + string(s.Stage))
			if q == "" || strings.Contains(blob, q) {
				hits = append(hits, map[string]string{"kind": "ship", "id": s.ID, "name": s.Name})
			}
		}
	}
	for _, n := range k.Notes() {
		blob := strings.ToLower(n.Claim + " " + n.Source + " " + n.Quote)
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "note", "id": n.Source, "name": n.Claim})
		}
	}
	for _, p := range k.Pulses() {
		blob := strings.ToLower(p.Note)
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "pulse", "id": strconv.Itoa(p.ID), "name": p.Note})
		}
	}
	for _, m := range k.Inbox("") {
		blob := strings.ToLower(m.From + " " + m.To + " " + m.Kind + " " + m.Body)
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "mail", "id": strconv.Itoa(m.ID), "name": m.Body})
		}
	}
	for _, c := range k.Confirms() {
		blob := strings.ToLower(c.Body + " " + c.Status)
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "confirm", "id": strconv.Itoa(c.ID), "name": c.Body})
		}
	}
	_, _ = k.Post("explorer", "research", "search", q)
	return kernel.Result{OK: true, Message: "explorer hit " + strconv.Itoa(len(hits)), Data: hits}, nil
}

func memoryInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "memory.store":
		topic := payloadQuery(call, "topic")
		text := payloadQuery(call, "text")
		if text == "" {
			text = payloadQuery(call, "prompt")
		}
		if topic == "" {
			topic = "working"
		}
		f := k.Remember(topic, text)
		return kernel.Result{OK: true, Message: "remembered " + f.Topic, Data: f}, nil
	case "memory.recall":
		q := payloadQuery(call, "query")
		if q == "" {
			q = payloadQuery(call, "topic")
		}
		got := k.Recall(q)
		return kernel.Result{OK: true, Message: "recalled " + strconv.Itoa(len(got)), Data: got}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func commsInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "comms.send":
		body := payloadQuery(call, "body")
		if body == "" {
			body = payloadQuery(call, "prompt")
		}
		if body == "" {
			body = "outbound draft"
		}
		c := k.RequestConfirm("comms", "comms.send", body)
		return kernel.Result{OK: true, Message: "parked confirm #" + strconv.Itoa(c.ID), Data: c}, nil
	case "comms.pending":
		return kernel.Result{OK: true, Message: "confirms", Data: k.Confirms()}, nil
	case "comms.allow":
		id := payloadInt(call, "id")
		c, err := k.DecideConfirm(id, true)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: c.Outcome, Data: c}, nil
	case "comms.deny":
		id := payloadInt(call, "id")
		c, err := k.DecideConfirm(id, false)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: c.Outcome, Data: c}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

// Brief is a deterministic design card. The model does not fill this.
type Brief struct {
	Goal     string   `json:"goal"`
	Shape    string   `json:"shape"`
	Line     []string `json:"line"`
	Agents   []string `json:"agents"`
	Gates    []string `json:"gates"`
	ClosedOK bool     `json:"closedOk"`
}

func architectInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	goal := payloadQuery(call, "goal")
	if goal == "" {
		goal = payloadQuery(call, "prompt")
	}
	if goal == "" {
		goal = "unnamed system"
	}
	brief := shapeBrief(goal)
	k.WriteNote(kernel.Note{
		Agent:  "architect",
		Source: "architect.plan",
		Claim:  brief.Shape + " brief: " + brief.Goal,
		Quote:  "line " + strings.Join(brief.Line, " → ") + " · agents " + strings.Join(brief.Agents, ", "),
	})
	k.Remember("design", brief.Goal)
	_, _ = k.Post("architect", "planner", "brief", brief.Goal)
	_, _ = k.Post("architect", "delivery", "brief", brief.Goal)
	return kernel.Result{OK: true, Message: "shaped " + brief.Shape + " brief", Data: brief}, nil
}

func shapeBrief(goal string) Brief {
	q := strings.ToLower(goal)
	shape := "app"
	switch {
	case strings.Contains(q, "os") || strings.Contains(q, "kernel"):
		shape = "os"
	case strings.Contains(q, "workflow") || strings.Contains(q, "pipeline"):
		shape = "workflow"
	case strings.Contains(q, "agent"):
		shape = "agent"
	}
	agents := []string{"architect", "delivery", "watch", "memory"}
	if strings.Contains(q, "research") || strings.Contains(q, "aos") || strings.Contains(q, "paper") {
		agents = append(agents, "research")
	}
	if strings.Contains(q, "mail") || strings.Contains(q, "outbound") || strings.Contains(q, "comms") {
		agents = append(agents, "comms")
	}
	if strings.Contains(q, "browser") || strings.Contains(q, "desktop") {
		agents = append(agents, "operator")
	}
	if strings.Contains(q, "cve") || strings.Contains(q, "sast") {
		agents = append(agents, "security")
	}
	gates := []string{"comms.send"}
	if strings.Contains(q, "production") || strings.Contains(q, "release") || strings.Contains(q, "deploy") {
		agents = append(agents, "deploy")
		gates = append(gates, "deploy.release")
	}
	return Brief{
		Goal:     strings.TrimSpace(goal),
		Shape:    shape,
		Line:     []string{"idea", "concept", "production"},
		Agents:   agents,
		Gates:    gates,
		ClosedOK: true,
	}
}

func plannerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	goal := payloadQuery(call, "goal")
	if goal == "" {
		goal = payloadQuery(call, "prompt")
	}
	if goal == "" {
		goal = "unnamed mandate"
	}
	cat := k.Catalog()
	if cat == nil {
		return kernel.Result{OK: false, Message: "delivery board not attached"}, nil
	}
	items := payloadList(call, "items")
	if len(items) == 0 {
		items = []string{goal}
	}
	created := make([]catalog.Ship, 0, len(items))
	for _, item := range items {
		ship, err := cat.Create(catalog.CreateShip{Name: item, Notes: goal, Sector: "research", Stack: []string{"Go"}})
		if err != nil {
			return kernel.Result{}, err
		}
		created = append(created, ship)
	}
	_, _ = k.Post("planner", "delivery", "backlog", goal)
	k.Remember("plan", goal)
	return kernel.Result{OK: true, Message: "parked " + strconv.Itoa(len(created)) + " on idea", Data: created}, nil
}

func researchInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "research.list":
		return kernel.Result{OK: true, Message: "notes " + strconv.Itoa(len(k.Notes())), Data: k.Notes()}, nil
	case "research.ingest", "note.write":
		n := kernel.Note{
			Agent:  "research",
			Source: payloadQuery(call, "source"),
			URL:    payloadQuery(call, "url"),
			Claim:  payloadQuery(call, "claim"),
			Quote:  payloadQuery(call, "quote"),
		}
		if n.Claim == "" {
			n.Claim = payloadQuery(call, "prompt")
		}
		if n.Source == "" {
			n.Source = "desk"
		}
		if n.Claim == "" {
			return kernel.Result{OK: false, Message: "claim required"}, nil
		}
		got := k.WriteNote(n)
		k.Remember("research", got.Claim)
		return kernel.Result{OK: true, Message: "ingested " + got.Source, Data: got}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func seedResearch(k *kernel.Kernel) {
	seeds := []kernel.Note{
		{Agent: "research", Source: "arxiv:2606.01508", URL: "https://arxiv.org/abs/2606.01508", Claim: "An AOS is a control plane over a classical OS, not a Linux replacement.", Quote: "AOS mediates tool invocations the way a classical OS mediates syscalls. First-class entities include agent identity, goals, capabilities, context, and execution records."},
		{Agent: "research", Source: "arxiv:2608.03214", URL: "https://arxiv.org/abs/2608.03214", Claim: "Split governance from runtime: intent/policy/audit vs lifecycle/routing/memory.", Quote: "Control & Governance plane owns authority and human oversight. Runtime & Coordination owns agent lifecycle and tool routing."},
		{Agent: "research", Source: "XKernel", URL: "https://github.com/JosephBerm/XKernel", Claim: "Treat agents as first-class processes with capability tokens and typed IPC.", Quote: "Unix treats processes; Kubernetes treats containers; an agent OS treats agents."},
		{Agent: "research", Source: "12-factor agents", URL: "https://github.com/humanlayer/12-factor-agents", Claim: "Own the loop in deterministic code. The model only fills structured next steps.", Quote: "Human confirm sits between selection and invocation. OpenRouter stays optional."},
		{Agent: "research", Source: "closed-hours", URL: "https://github.com/cashtro/cashtro", Claim: "An AOS stays up while the operator is away: persist, pulse, keep internal work flowing, and park outbound behind a human gate.", Quote: "Keep the work flowing through here while things are closed. Delivery, research, and memory stay live. Comms do not send until allow."},
		{Agent: "research", Source: "treg", URL: "https://treg.to", Claim: "Exa publication search is available as research.ingest input at $0.007/call when Treg is signed in.", Quote: "catalog_search → catalog_get → call. Token was expired this pass; notes still landed from open sources."},
	}
	for _, n := range seeds {
		k.WriteNote(n)
		k.Remember("research", n.Claim)
	}
}

func payloadQuery(call kernel.Call, key string) string {
	if len(call.Payload) == 0 {
		return ""
	}
	var obj map[string]any
	if err := json.Unmarshal(call.Payload, &obj); err != nil {
		return strings.TrimSpace(string(call.Payload))
	}
	if v, ok := obj[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func payloadList(call kernel.Call, key string) []string {
	if len(call.Payload) == 0 {
		return nil
	}
	var obj map[string]any
	if err := json.Unmarshal(call.Payload, &obj); err != nil {
		return nil
	}
	raw, ok := obj[key]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	case string:
		parts := strings.Split(v, "\n")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if s := strings.TrimSpace(p); s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func payloadInt(call kernel.Call, key string) int {
	if len(call.Payload) == 0 {
		return 0
	}
	var obj map[string]any
	if err := json.Unmarshal(call.Payload, &obj); err != nil {
		return 0
	}
	switch v := obj[key].(type) {
	case float64:
		return int(v)
	}
	return 0
}

type researchAgent struct {
	k *kernel.Kernel
}

func (a *researchAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "research", Name: "Research", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "library", Summary: "First-class notes. Gathers AOS papers and desk findings while the OS is under construction.",
		Capabilities: []string{"research.list", "research.ingest", "note.write"},
		Autostart:    true,
	}
}

func (a *researchAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	seedResearch(k)
	k.Publish("research", "library", "seeded construction notes", map[string]any{"notes": len(k.Notes())})
	return nil
}

func (a *researchAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	return researchInvoke(a.k, call)
}

type watchAgent struct {
	k *kernel.Kernel
}

func (a *watchAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "watch", Name: "Watch", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role: "night", Summary: "Keeps the desk flowing while things are closed. Pulse and persist. Outbound still needs a human.",
		Capabilities: []string{"watch.status", "watch.close", "watch.open", "watch.pulse"},
		Autostart:    true,
	}
}

func (a *watchAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	k.Publish("watch", "ready", "closed-hours watch online", nil)
	return nil
}

func (a *watchAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	return watchInvoke(a.k, call)
}

func watchInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	note := payloadQuery(call, "note")
	if note == "" {
		note = payloadQuery(call, "prompt")
	}
	switch call.Capability {
	case "watch.status":
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	case "watch.close":
		if note == "" {
			note = "things are closed"
		}
		k.SetClosed(true)
		k.RecordPulse(note)
		keepFlowing(k, note)
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	case "watch.open":
		if note == "" {
			note = "desk open"
		}
		k.SetClosed(false)
		k.RecordPulse(note)
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	case "watch.pulse":
		if note == "" {
			note = "pulse"
		}
		k.RecordPulse(note)
		if k.Closed() {
			keepFlowing(k, note)
		}
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func keepFlowing(k *kernel.Kernel, note string) {
	if cat := k.Catalog(); cat != nil {
		if _, err := cat.Get("closed-hours-flow"); err != nil {
			_, _ = cat.Create(catalog.CreateShip{
				Name:   "Closed-hours flow",
				Client: "Cashtro OS",
				Sector: "ops",
				Stack:  []string{"Go"},
				Notes:  "Keep the work flowing through here while things are closed.",
			})
		}
	}
	k.Remember("watch", "closed hours · "+note)
	_, _ = k.Post("watch", "research", "pulse", "keep gathering while the desk is closed")
	_, _ = k.Post("watch", "planner", "pulse", "keep the line moving")
}

const maxTraceHits = 32

// TraceHit is one node in an investigator blast radius.
type TraceHit struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Trace is a deterministic incident map. The model does not fill this.
type Trace struct {
	Query string     `json:"query"`
	Hits  []TraceHit `json:"hits"`
}

func investigatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	q := payloadQuery(call, "query")
	if q == "" {
		q = payloadQuery(call, "prompt")
	}
	if q == "" {
		if k.Closed() {
			q = "closed"
		} else {
			q = "error"
		}
	}
	needle := strings.ToLower(q)
	hits := make([]TraceHit, 0)
	add := func(kind, id, name string) {
		if len(hits) >= maxTraceHits {
			return
		}
		blob := strings.ToLower(kind + " " + id + " " + name)
		if strings.Contains(blob, needle) {
			hits = append(hits, TraceHit{Kind: kind, ID: id, Name: name})
		}
	}
	for _, p := range k.Processes() {
		add("agent", p.Spec.ID, p.Spec.Name+" "+p.Spec.Role+" "+p.Spec.Summary)
	}
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			add("ship", s.ID, s.Name+" "+s.Notes+" "+string(s.Stage))
		}
	}
	for _, n := range k.Notes() {
		add("note", n.Source, n.Claim)
	}
	for _, ev := range k.Events() {
		add("event", ev.Source+"/"+ev.Kind, ev.Message)
	}
	for _, m := range k.Inbox("") {
		add("mail", m.To, m.Body)
	}
	for _, c := range k.Confirms() {
		add("confirm", strconv.Itoa(c.ID), c.Body+" "+c.Status)
	}
	for _, p := range k.Pulses() {
		add("pulse", strconv.Itoa(p.ID), p.Note)
	}
	for _, f := range k.Recall("") {
		add("fact", f.Topic, f.Text)
	}
	k.Remember("incident", q)
	_, _ = k.Post("investigator", "security", "trace", q)
	return kernel.Result{OK: true, Message: "traced " + strconv.Itoa(len(hits)), Data: Trace{Query: q, Hits: hits}}, nil
}

const maxFindings = 32

// Finding is one security.triage hit.
type Finding struct {
	Kind   string `json:"kind"`
	Source string `json:"source"`
	Claim  string `json:"claim"`
}

// Triage is a deterministic scan of the desk. The model does not fill this.
type Triage struct {
	Clear    bool      `json:"clear"`
	Findings []Finding `json:"findings"`
}

func securityInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	findings := make([]Finding, 0)
	add := func(source, claim string) {
		if len(findings) >= maxFindings {
			return
		}
		kind := classifyFinding(claim)
		if kind == "" {
			return
		}
		findings = append(findings, Finding{Kind: kind, Source: source, Claim: claim})
	}
	for _, n := range k.Notes() {
		add("note:"+n.Source, n.Claim+" "+n.Quote)
	}
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			add("ship:"+s.ID, s.Name+" "+s.Notes)
		}
	}
	for _, ev := range k.Events() {
		add("event:"+ev.Source+"/"+ev.Kind, ev.Message)
	}
	for _, f := range k.Recall("") {
		add("fact:"+f.Topic, f.Text)
	}
	for _, m := range k.Inbox("") {
		add("mail:"+m.To, m.Body)
	}
	msg := "clear"
	if len(findings) > 0 {
		msg = "flagged " + strconv.Itoa(len(findings))
		k.Remember("security", msg)
		_, _ = k.Post("security", "investigator", "triage", msg)
	}
	return kernel.Result{OK: true, Message: msg, Data: Triage{Clear: len(findings) == 0, Findings: findings}}, nil
}

func classifyFinding(blob string) string {
	b := strings.ToLower(blob)
	switch {
	case strings.Contains(b, "cve-") || strings.Contains(b, "cve "):
		return "cve"
	case strings.Contains(b, "sast"):
		return "sast"
	case strings.Contains(b, "secret") || strings.Contains(b, "leak"):
		return "secret"
	case strings.Contains(b, "vulnerab") || strings.Contains(b, "xss") || strings.Contains(b, "sqli"):
		return "vuln"
	default:
		return ""
	}
}
