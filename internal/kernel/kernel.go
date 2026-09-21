// Package kernel is the Cashtro AI OS control plane.
//
// Every agentic we build here is a process. It registers a spec, exposes
// capabilities, and talks on the bus. Delivery (idea → concept → production)
// is the first live subsystem. New agentics boot into this kernel — they do
// not become a second product.
package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
)

const (
	// Name is the public OS name.
	Name = "Cashtro OS"
	// Version is the kernel release.
	Version = "0.4.9"
)

// Status is a process lifecycle state.
type Status string

const (
	StatusStopped Status = "stopped"
	StatusBooting Status = "booting"
	StatusRunning Status = "running"
	StatusError   Status = "error"
)

// Kind distinguishes system daemons from user agentics.
type Kind string

const (
	KindSystem Kind = "system"
	KindUser   Kind = "user"
)

// Mode says whether an agentic can execute work today.
type Mode string

const (
	ModeLive     Mode = "live"
	ModeResident Mode = "resident"
)

// Spec is the contract an agentic publishes at registration.
type Spec struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Kind         Kind     `json:"kind"`
	Mode         Mode     `json:"mode"`
	Role         string   `json:"role"`
	Summary      string   `json:"summary"`
	Capabilities []string `json:"capabilities"`
	Autostart    bool     `json:"autostart"`
}

// Process is a running (or stopped) agentic.
type Process struct {
	PID       int       `json:"pid"`
	Spec      Spec      `json:"spec"`
	Status    Status    `json:"status"`
	StartedAt time.Time `json:"startedAt,omitempty"`
	LastBeat  time.Time `json:"lastBeat,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// Event is one line on the kernel journal.
type Event struct {
	Seq     int            `json:"seq"`
	At      time.Time      `json:"at"`
	Source  string         `json:"source"`
	Kind    string         `json:"kind"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

// Call is a capability invocation.
type Call struct {
	Capability string          `json:"capability"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

// Result is what an agentic returns from Invoke.
type Result struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Capability is a named verb the OS can route.
type Capability struct {
	ID      string `json:"id"`
	Agent   string `json:"agent"`
	Summary string `json:"summary"`
	Live    bool   `json:"live"`
}

// Agent is a loadable agentic.
type Agent interface {
	Spec() Spec
	Boot(ctx context.Context, k *Kernel) error
	Invoke(ctx context.Context, call Call) (Result, error)
}

// About is the public OS card.
type About struct {
	Name      string     `json:"name"`
	Version   string     `json:"version"`
	Motto     string     `json:"motto"`
	Kernel    Status     `json:"kernel"`
	Agents    int        `json:"agents"`
	Running   int        `json:"running"`
	Live      int        `json:"live"`
	Resident  int        `json:"resident"`
	Events    int        `json:"events"`
	Closed    bool       `json:"closed"`
	LastPulse *time.Time `json:"lastPulse,omitempty"`
	Manifesto string     `json:"manifesto"`
}

var (
	// ErrUnknownAgent is returned when an agent id is not registered.
	ErrUnknownAgent = errors.New("unknown agent")
	// ErrUnknownCapability is returned when a capability is not on the agent.
	ErrUnknownCapability = errors.New("unknown capability")
	// ErrNotRunning is returned when invoke hits a stopped process.
	ErrNotRunning = errors.New("agent not running")
)

const maxEvents = 200

// Kernel is the in-process OS.
type Kernel struct {
	mu          sync.RWMutex
	now         func() time.Time
	nextID      int
	seq         int
	procs       map[string]*Process
	agents      map[string]Agent
	caps        map[string]Capability
	events      []Event
	cat         *catalog.Catalog
	mail        []Mail
	mailSeq     int
	notes       []Note
	noteSeq     int
	facts       []Fact
	factSeq     int
	confirms    []Confirm
	confSeq     int
	persistPath string
	closed      bool
	pulses      []Pulse
	pulseSeq    int
}

// Option configures the kernel.
type Option func(*Kernel)

// WithClock injects a clock for tests.
func WithClock(now func() time.Time) Option {
	return func(k *Kernel) {
		k.now = now
	}
}

// WithPersistPath writes the OS image to disk after mutates.
func WithPersistPath(path string) Option {
	return func(k *Kernel) {
		k.persistPath = path
	}
}

// New returns an empty kernel. Call Register then Boot.
func New(opts ...Option) *Kernel {
	k := &Kernel{
		now:    func() time.Time { return time.Now().UTC() },
		procs:  make(map[string]*Process),
		agents: make(map[string]Agent),
		caps:   make(map[string]Capability),
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

// AttachCatalog lets the delivery agent publish the live board.
func (k *Kernel) AttachCatalog(cat *catalog.Catalog) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.cat = cat
}

// Catalog returns the delivery board, if the delivery agent has attached it.
func (k *Kernel) Catalog() *catalog.Catalog {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.cat
}

// Register loads an agentic into the process table as stopped.
func (k *Kernel) Register(agent Agent) {
	spec := agent.Spec()
	k.mu.Lock()
	defer k.mu.Unlock()
	k.nextID++
	k.agents[spec.ID] = agent
	k.procs[spec.ID] = &Process{PID: k.nextID, Spec: spec, Status: StatusStopped}
	live := spec.Mode == ModeLive
	for _, cap := range spec.Capabilities {
		k.caps[cap] = Capability{ID: cap, Agent: spec.ID, Live: live, Summary: spec.Summary}
	}
}

// Boot starts every autostart agentic and writes the init journal.
func (k *Kernel) Boot(ctx context.Context) error {
	k.Publish("init", "boot", Name+" "+Version+" coming up", nil)
	ids := k.autostartIDs()
	for _, id := range ids {
		if err := k.Spawn(ctx, id); err != nil {
			return err
		}
	}
	about := k.About()
	msg := fmt.Sprintf("kernel online · %d agentics · %d live", about.Agents, about.Live)
	if about.Closed {
		msg += " · desk closed · work still flowing"
	}
	k.Publish("init", "ready", msg, nil)
	return nil
}

// Spawn boots one agentic.
func (k *Kernel) Spawn(ctx context.Context, id string) error {
	k.mu.Lock()
	agent, ok := k.agents[id]
	proc, pok := k.procs[id]
	if !ok || !pok {
		k.mu.Unlock()
		return fmt.Errorf("%w: %s", ErrUnknownAgent, id)
	}
	if proc.Status == StatusRunning {
		k.mu.Unlock()
		return nil
	}
	proc.Status = StatusBooting
	k.mu.Unlock()

	if err := agent.Boot(ctx, k); err != nil {
		k.mu.Lock()
		proc.Status = StatusError
		proc.Note = err.Error()
		k.mu.Unlock()
		k.Publish(id, "error", err.Error(), nil)
		return err
	}

	now := k.now()
	k.mu.Lock()
	proc.Status = StatusRunning
	proc.StartedAt = now
	proc.LastBeat = now
	proc.Note = ""
	k.mu.Unlock()
	k.Publish(id, "spawn", proc.Spec.Name+" running", map[string]any{"mode": string(proc.Spec.Mode)})
	return nil
}

// Invoke routes a capability call to a running agentic.
func (k *Kernel) Invoke(ctx context.Context, id string, call Call) (Result, error) {
	k.mu.RLock()
	agent, ok := k.agents[id]
	proc, pok := k.procs[id]
	k.mu.RUnlock()
	if !ok || !pok {
		return Result{}, fmt.Errorf("%w: %s", ErrUnknownAgent, id)
	}
	if proc.Status != StatusRunning {
		return Result{}, fmt.Errorf("%w: %s", ErrNotRunning, id)
	}
	if !hasCap(proc.Spec.Capabilities, call.Capability) {
		return Result{}, fmt.Errorf("%w: %s", ErrUnknownCapability, call.Capability)
	}

	res, err := agent.Invoke(ctx, call)
	now := k.now()
	k.mu.Lock()
	proc.LastBeat = now
	k.mu.Unlock()
	if err != nil {
		k.Publish(id, "invoke.error", err.Error(), map[string]any{"capability": call.Capability})
		return Result{}, err
	}
	k.Publish(id, "invoke", res.Message, map[string]any{"capability": call.Capability, "ok": res.OK})
	k.persist()
	return res, nil
}

// Processes returns the process table, system first then name.
func (k *Kernel) Processes() []Process {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Process, 0, len(k.procs))
	for _, p := range k.procs {
		cp := *p
		cp.Spec.Capabilities = append([]string(nil), p.Spec.Capabilities...)
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Spec.Kind != out[j].Spec.Kind {
			return out[i].Spec.Kind == KindSystem
		}
		return out[i].Spec.Name < out[j].Spec.Name
	})
	return out
}

// Process returns one process by agent id.
func (k *Kernel) Process(id string) (Process, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	p, ok := k.procs[id]
	if !ok {
		return Process{}, fmt.Errorf("%w: %s", ErrUnknownAgent, id)
	}
	cp := *p
	cp.Spec.Capabilities = append([]string(nil), p.Spec.Capabilities...)
	return cp, nil
}

// Capabilities returns every registered verb.
func (k *Kernel) Capabilities() []Capability {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Capability, 0, len(k.caps))
	for _, c := range k.caps {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Events returns the journal, oldest first.
func (k *Kernel) Events() []Event {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Event, len(k.events))
	copy(out, k.events)
	return out
}

// Publish appends a journal line.
func (k *Kernel) Publish(source, kind, message string, data map[string]any) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.seq++
	ev := Event{
		Seq:     k.seq,
		At:      k.now(),
		Source:  source,
		Kind:    kind,
		Message: message,
		Data:    data,
	}
	k.events = append(k.events, ev)
	if len(k.events) > maxEvents {
		k.events = append([]Event(nil), k.events[len(k.events)-maxEvents:]...)
	}
}

// About returns the OS card.
func (k *Kernel) About() About {
	k.mu.RLock()
	defer k.mu.RUnlock()
	var running, live, resident int
	for _, p := range k.procs {
		if p.Status == StatusRunning {
			running++
		}
		if p.Spec.Mode == ModeLive {
			live++
		} else {
			resident++
		}
	}
	var lastPulse *time.Time
	if n := len(k.pulses); n > 0 {
		t := k.pulses[n-1].At
		lastPulse = &t
	}
	return About{
		Name:      Name,
		Version:   Version,
		Motto:     "I make teams ship: idea → concept → production.",
		Kernel:    StatusRunning,
		Agents:    len(k.procs),
		Running:   running,
		Live:      live,
		Resident:  resident,
		Events:    len(k.events),
		Closed:    k.closed,
		LastPulse: lastPulse,
		Manifesto: "Cashtro OS is under construction — the control plane for every agentic we build here. " +
			"Agents are processes. Capabilities are verbs. Mail, notes, memory, and confirms are first-class. " +
			"Delivery is live. Research is live. Watch keeps the desk flowing while things are closed. " +
			"Outbound comms still wait at the human gate. OpenRouter stays optional. " +
			"New agentics register into this kernel — they do not fork a second product.",
	}
}

func (k *Kernel) autostartIDs() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	ids := make([]string, 0, len(k.procs))
	for id, p := range k.procs {
		if p.Spec.Autostart {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func hasCap(caps []string, id string) bool {
	for _, c := range caps {
		if c == id {
			return true
		}
	}
	return false
}
