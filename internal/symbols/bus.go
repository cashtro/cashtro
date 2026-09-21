// Package symbols is Cashtro's internal Zapier.
//
// A Symbol is trigger → kernel verbs. Connectors are our process table,
// not a paid app store. The bus stores definitions and run history; the
// symbols agentic executes steps through the kernel. No subscription.
package symbols

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Kind is how a symbol starts.
type Kind string

const (
	// TriggerManual is a desk or API fire.
	TriggerManual Kind = "manual"
	// TriggerWebhook is POST /api/symbols/{id}/hook — catch-hook, on us.
	TriggerWebhook Kind = "webhook"
	// TriggerEvent matches a kernel journal kind (mail, note, confirm.pending).
	TriggerEvent Kind = "event"
)

// Trigger is the left-hand node of a symbol.
type Trigger struct {
	Kind  Kind   `json:"kind"`
	Match string `json:"match,omitempty"`
}

// Step is one verb on the line.
type Step struct {
	Cap    string            `json:"cap"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Symbol is one automation on the desk.
type Symbol struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Summary   string    `json:"summary"`
	On        Trigger   `json:"on"`
	Steps     []Step    `json:"steps"`
	Enabled   bool      `json:"enabled"`
	Fires     int       `json:"fires"`
	LastOK    bool      `json:"lastOk"`
	LastMsg   string    `json:"lastMsg,omitempty"`
	LastAt    time.Time `json:"lastAt,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Create is the input for a new symbol.
type Create struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	On      Trigger `json:"on"`
	Steps   []Step  `json:"steps"`
}

// StepLog is one verb's result inside a run.
type StepLog struct {
	Cap     string `json:"cap"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Run is one firing of a symbol.
type Run struct {
	ID      int               `json:"id"`
	Symbol  string            `json:"symbol"`
	Name    string            `json:"name"`
	At      time.Time         `json:"at"`
	OK      bool              `json:"ok"`
	Message string            `json:"message"`
	Input   map[string]string `json:"input,omitempty"`
	Steps   []StepLog         `json:"steps"`
}

var (
	// ErrNotFound is returned when a symbol id is unknown.
	ErrNotFound = errors.New("symbol not found")
	// ErrInvalid is returned for a symbol that cannot be created or fired.
	ErrInvalid = errors.New("invalid symbol")
	// ErrDisabled is returned when a symbol is switched off.
	ErrDisabled = errors.New("symbol disabled")
	// ErrBusy is returned when the same symbol is already mid-fire.
	ErrBusy = errors.New("symbol busy")
)

const maxRuns = 100

var slot = regexp.MustCompile(`\{\{([a-zA-Z0-9_.]+)\}\}`)

// Bus is an in-memory symbol table.
type Bus struct {
	mu     sync.RWMutex
	seq    int
	runSeq int
	now    func() time.Time
	syms   map[string]*Symbol
	runs   []Run
	busy   map[string]bool
}

// Option configures a bus.
type Option func(*Bus)

// WithClock injects a clock for tests.
func WithClock(now func() time.Time) Option {
	return func(b *Bus) {
		b.now = now
	}
}

// New returns a bus seeded with internal automations. Zero Zapier.
func New(opts ...Option) *Bus {
	b := &Bus{
		now:  func() time.Time { return time.Now().UTC() },
		syms: make(map[string]*Symbol),
		busy: make(map[string]bool),
	}
	for _, opt := range opts {
		opt(b)
	}
	b.seed()
	return b
}

// List returns symbols ordered by name.
func (b *Bus) List() []Symbol {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Symbol, 0, len(b.syms))
	for _, s := range b.syms {
		out = append(out, cloneSymbol(s))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Get returns one symbol by id.
func (b *Bus) Get(id string) (Symbol, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	s, ok := b.syms[id]
	if !ok {
		return Symbol{}, ErrNotFound
	}
	return cloneSymbol(s), nil
}

// Create adds a symbol. Default trigger is manual.
func (b *Bus) Create(in Create) (Symbol, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Symbol{}, fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if len(in.Steps) == 0 {
		return Symbol{}, fmt.Errorf("%w: at least one step", ErrInvalid)
	}
	steps, err := cleanSteps(in.Steps)
	if err != nil {
		return Symbol{}, err
	}
	on := in.On
	if on.Kind == "" {
		on.Kind = TriggerManual
	}
	switch on.Kind {
	case TriggerManual, TriggerWebhook, TriggerEvent:
	default:
		return Symbol{}, fmt.Errorf("%w: unknown trigger %q", ErrInvalid, on.Kind)
	}
	on.Match = strings.TrimSpace(on.Match)
	if on.Kind == TriggerEvent && on.Match == "" {
		return Symbol{}, fmt.Errorf("%w: event trigger needs match", ErrInvalid)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.uniqueID(name)
	now := b.now()
	s := &Symbol{
		ID:        id,
		Name:      name,
		Summary:   strings.TrimSpace(in.Summary),
		On:        on,
		Steps:     steps,
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	b.syms[id] = s
	return cloneSymbol(s), nil
}

// Toggle flips enabled.
func (b *Bus) Toggle(id string) (Symbol, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, ok := b.syms[id]
	if !ok {
		return Symbol{}, ErrNotFound
	}
	s.Enabled = !s.Enabled
	s.UpdatedAt = b.now()
	return cloneSymbol(s), nil
}

// Enter marks a symbol in-flight. False if missing, disabled, or busy.
func (b *Bus) Enter(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, ok := b.syms[id]
	if !ok {
		return ErrNotFound
	}
	if !s.Enabled {
		return ErrDisabled
	}
	if b.busy[id] {
		return ErrBusy
	}
	b.busy[id] = true
	return nil
}

// Leave clears the in-flight mark.
func (b *Bus) Leave(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.busy, id)
}

// Finish records a run against a symbol.
func (b *Bus) Finish(id string, ok bool, msg string, input map[string]string, logs []StepLog) (Run, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, found := b.syms[id]
	if !found {
		return Run{}, ErrNotFound
	}
	b.runSeq++
	now := b.now()
	run := Run{
		ID:      b.runSeq,
		Symbol:  id,
		Name:    s.Name,
		At:      now,
		OK:      ok,
		Message: msg,
		Input:   cloneMap(input),
		Steps:   append([]StepLog(nil), logs...),
	}
	s.Fires++
	s.LastOK = ok
	s.LastMsg = msg
	s.LastAt = now
	s.UpdatedAt = now
	b.runs = append(b.runs, run)
	if len(b.runs) > maxRuns {
		b.runs = append([]Run(nil), b.runs[len(b.runs)-maxRuns:]...)
	}
	return cloneRun(run), nil
}

// Runs returns history, newest first.
func (b *Bus) Runs() []Run {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Run, len(b.runs))
	for i := range b.runs {
		out[len(b.runs)-1-i] = cloneRun(b.runs[i])
	}
	return out
}

// MatchEvent reports whether an event trigger should fire for this journal line.
func MatchEvent(on Trigger, source, kind string) bool {
	if on.Kind != TriggerEvent {
		return false
	}
	switch kind {
	case "invoke", "invoke.error", "spawn", "boot", "ready", "bind":
		return false
	}
	m := strings.ToLower(strings.TrimSpace(on.Match))
	if m == "" {
		return false
	}
	src := strings.ToLower(source)
	k := strings.ToLower(kind)
	return m == k || m == src+"/"+k || m == src
}

// Render fills {{slots}} from input. Missing keys become empty.
func Render(fields map[string]string, in map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = RenderString(v, in)
	}
	return out
}

// RenderString replaces {{key}} tokens.
func RenderString(tmpl string, in map[string]string) string {
	if in == nil {
		in = map[string]string{}
	}
	return slot.ReplaceAllStringFunc(tmpl, func(m string) string {
		parts := slot.FindStringSubmatch(m)
		if len(parts) != 2 {
			return ""
		}
		return in[parts[1]]
	})
}

// FlattenJSON turns a request body into string slots for {{templates}}.
func FlattenJSON(raw []byte) map[string]string {
	raw = bytesTrim(raw)
	if len(raw) == 0 {
		return map[string]string{}
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return map[string]string{"body": strings.TrimSpace(string(raw))}
	}
	return Flatten(v)
}

// Flatten walks JSON-ish values into dotted keys.
func Flatten(v any) map[string]string {
	out := make(map[string]string)
	flattenInto(out, "", v)
	return out
}

func flattenInto(out map[string]string, prefix string, v any) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			flattenInto(out, key, val)
		}
	case []any:
		for i, val := range t {
			key := strconv.Itoa(i)
			if prefix != "" {
				key = prefix + "." + key
			}
			flattenInto(out, key, val)
		}
	case string:
		putFlat(out, prefix, t)
	case float64:
		if t == float64(int64(t)) {
			putFlat(out, prefix, strconv.FormatInt(int64(t), 10))
		} else {
			putFlat(out, prefix, strconv.FormatFloat(t, 'f', -1, 64))
		}
	case json.Number:
		putFlat(out, prefix, t.String())
	case bool:
		putFlat(out, prefix, strconv.FormatBool(t))
	case nil:
		return
	default:
		putFlat(out, prefix, fmt.Sprint(t))
	}
}

func putFlat(out map[string]string, key, val string) {
	val = strings.TrimSpace(val)
	if key == "" {
		out["body"] = val
		return
	}
	out[key] = val
}

func bytesTrim(raw []byte) []byte {
	return []byte(strings.TrimSpace(string(raw)))
}

func (b *Bus) seed() {
	fixed := time.Date(2026, 9, 21, 7, 0, 0, 0, time.UTC)
	items := []Symbol{
		{
			ID:      "desk-pulse",
			Name:    "Desk pulse",
			Summary: "Manual health check. Kernel verbs only — no paid bus.",
			On:      Trigger{Kind: TriggerManual},
			Enabled: true,
			Steps: []Step{
				{Cap: "os.about"},
				{Cap: "explorer.search", Fields: map[string]string{"query": "live"}},
			},
		},
		{
			ID:      "intake-to-idea",
			Name:    "Intake → idea",
			Summary: "Catch a webhook and park a mandate. Zapier catch-hook, hosted here.",
			On:      Trigger{Kind: TriggerWebhook},
			Enabled: true,
			Steps: []Step{
				{Cap: "delivery.create", Fields: map[string]string{
					"name":   "{{name}}",
					"client": "{{client}}",
					"sector": "{{sector}}",
					"notes":  "{{notes}}",
				}},
				{Cap: "memory.store", Fields: map[string]string{
					"topic": "intake",
					"text":  "{{name}}",
				}},
			},
		},
		{
			ID:      "mail-to-library",
			Name:    "Mail → library",
			Summary: "When an agentic posts mail, ingest a research line. Event trigger, $0.",
			On:      Trigger{Kind: TriggerEvent, Match: "mail"},
			Enabled: true,
			Steps: []Step{
				{Cap: "research.ingest", Fields: map[string]string{
					"source": "symbols",
					"claim":  "{{message}}",
				}},
			},
		},
		{
			ID:      "signal-human",
			Name:    "Signal the human",
			Summary: "Park outbound copy at the confirm gate. No send without allow.",
			On:      Trigger{Kind: TriggerManual},
			Enabled: true,
			Steps: []Step{
				{Cap: "comms.send", Fields: map[string]string{"body": "{{body}}"}},
			},
		},
	}
	for i := range items {
		items[i].CreatedAt = fixed
		items[i].UpdatedAt = fixed
		cp := items[i]
		cp.Steps = cloneSteps(items[i].Steps)
		b.syms[cp.ID] = &cp
	}
}

func (b *Bus) uniqueID(name string) string {
	base := slugify(name)
	if base == "" {
		b.seq++
		return fmt.Sprintf("symbol-%d", b.seq)
	}
	if _, exists := b.syms[base]; !exists {
		return base
	}
	for n := 2; ; n++ {
		id := fmt.Sprintf("%s-%d", base, n)
		if _, exists := b.syms[id]; !exists {
			return id
		}
	}
}

func cleanSteps(steps []Step) ([]Step, error) {
	out := make([]Step, 0, len(steps))
	for _, s := range steps {
		cap := strings.TrimSpace(s.Cap)
		if cap == "" {
			return nil, fmt.Errorf("%w: step cap is required", ErrInvalid)
		}
		out = append(out, Step{Cap: cap, Fields: cloneMap(s.Fields)})
	}
	return out, nil
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

func cloneSymbol(s *Symbol) Symbol {
	out := *s
	out.Steps = cloneSteps(s.Steps)
	return out
}

func cloneSteps(steps []Step) []Step {
	if steps == nil {
		return nil
	}
	out := make([]Step, len(steps))
	for i, s := range steps {
		out[i] = Step{Cap: s.Cap, Fields: cloneMap(s.Fields)}
	}
	return out
}

func cloneRun(r Run) Run {
	r.Input = cloneMap(r.Input)
	if r.Steps != nil {
		r.Steps = append([]StepLog(nil), r.Steps...)
	}
	return r
}

func cloneMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
