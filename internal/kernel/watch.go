package kernel

import (
	"strings"
	"time"
)

const maxPulses = 64

// Pulse is one closed-hours heartbeat.
type Pulse struct {
	ID     int       `json:"id"`
	At     time.Time `json:"at"`
	Closed bool      `json:"closed"`
	Note   string    `json:"note"`
}

// Watch is the closed-hours control card.
type Watch struct {
	Closed    bool     `json:"closed"`
	Message   string   `json:"message"`
	LastPulse *Pulse   `json:"lastPulse,omitempty"`
	Pulses    []Pulse  `json:"pulses"`
	Image     string   `json:"image,omitempty"`
	Flowing   []string `json:"flowing"`
	Gated     []string `json:"gated"`
}

// ClosedHoursFlow is work the OS keeps doing while the desk is closed.
var ClosedHoursFlow = []string{
	"watch.pulse",
	"research.ingest",
	"research.list",
	"explorer.search",
	"memory.store",
	"memory.recall",
	"planner.backlog",
	"delivery.list",
	"delivery.create",
	"delivery.advance",
	"os.about",
	"chooser.list",
	"desk.plan",
	"desk.capture",
	"symbols.list",
	"symbols.fire",
}

// ClosedHoursGated still needs a human — outbound or a live operator.
var ClosedHoursGated = []string{
	"comms.send",
	"deploy.release",
	"operator.browse",
}

func closedMessage(closed bool) string {
	if closed {
		return "CLOSED HOURS · work still flowing through here"
	}
	return "OPEN HOURS · desk open · human on site"
}

// Closed reports whether the desk is in closed hours.
func (k *Kernel) Closed() bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.closed
}

// SetClosed flips closed hours and journals it.
func (k *Kernel) SetClosed(closed bool) {
	k.mu.Lock()
	k.closed = closed
	kind := "open"
	if closed {
		kind = "close"
	}
	msg := closedMessage(closed)
	k.seq++
	k.events = append(k.events, Event{
		Seq:     k.seq,
		At:      k.now(),
		Source:  "watch",
		Kind:    kind,
		Message: msg,
		Data:    map[string]any{"closed": closed},
	})
	if len(k.events) > maxEvents {
		k.events = append([]Event(nil), k.events[len(k.events)-maxEvents:]...)
	}
	k.mu.Unlock()
	k.persist()
}

// RecordPulse stores a heartbeat.
func (k *Kernel) RecordPulse(note string) Pulse {
	k.mu.Lock()
	k.pulseSeq++
	p := Pulse{
		ID:     k.pulseSeq,
		At:     k.now(),
		Closed: k.closed,
		Note:   strings.TrimSpace(note),
	}
	k.pulses = append(k.pulses, p)
	if len(k.pulses) > maxPulses {
		k.pulses = append([]Pulse(nil), k.pulses[len(k.pulses)-maxPulses:]...)
	}
	k.seq++
	k.events = append(k.events, Event{
		Seq:     k.seq,
		At:      k.now(),
		Source:  "watch",
		Kind:    "pulse",
		Message: closedMessage(k.closed),
		Data:    map[string]any{"id": p.ID, "note": p.Note},
	})
	if len(k.events) > maxEvents {
		k.events = append([]Event(nil), k.events[len(k.events)-maxEvents:]...)
	}
	out := p
	k.mu.Unlock()
	k.persist()
	return out
}

// Pulses returns heartbeats, newest first.
func (k *Kernel) Pulses() []Pulse {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Pulse, len(k.pulses))
	copy(out, k.pulses)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// WatchCard returns the closed-hours status.
func (k *Kernel) WatchCard() Watch {
	k.mu.RLock()
	defer k.mu.RUnlock()
	pulses := make([]Pulse, len(k.pulses))
	copy(pulses, k.pulses)
	for i, j := 0, len(pulses)-1; i < j; i, j = i+1, j-1 {
		pulses[i], pulses[j] = pulses[j], pulses[i]
	}
	w := Watch{
		Closed:  k.closed,
		Message: closedMessage(k.closed),
		Pulses:  pulses,
		Image:   k.persistPath,
		Flowing: append([]string(nil), ClosedHoursFlow...),
		Gated:   append([]string(nil), ClosedHoursGated...),
	}
	if n := len(k.pulses); n > 0 {
		last := k.pulses[n-1]
		w.LastPulse = &last
	}
	return w
}
