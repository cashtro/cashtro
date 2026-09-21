package kernel

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cashtro/cashtro/internal/catalog"
)

// Snapshot is the durable OS image for this repo.
type Snapshot struct {
	Version  string         `json:"version"`
	Closed   bool           `json:"closed"`
	Notes    []Note         `json:"notes"`
	Facts    []Fact         `json:"facts"`
	Confirms []Confirm      `json:"confirms"`
	Mail     []Mail         `json:"mail"`
	Pulses   []Pulse        `json:"pulses,omitempty"`
	Ships    []catalog.Ship `json:"ships,omitempty"`
	Choices  []Choice       `json:"choices,omitempty"`
	Scan     *Scan          `json:"scan,omitempty"`
	Requests []Request      `json:"requests,omitempty"`
}

// PersistPath returns the disk image path, if any.
func (k *Kernel) PersistPath() string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.persistPath
}

func (k *Kernel) persist() {
	path := k.PersistPath()
	if path == "" {
		return
	}
	_ = SaveFile(path, k)
}

// Snapshot copies durable state. Callers must not hold k.mu.
func (k *Kernel) Snapshot() Snapshot {
	k.mu.RLock()
	s := Snapshot{
		Version:  Version,
		Closed:   k.closed,
		Notes:    append([]Note(nil), k.notes...),
		Facts:    append([]Fact(nil), k.facts...),
		Confirms: append([]Confirm(nil), k.confirms...),
		Mail:     append([]Mail(nil), k.mail...),
		Pulses:   append([]Pulse(nil), k.pulses...),
		Choices:  append([]Choice(nil), k.choices...),
		Requests: append([]Request(nil), k.requests...),
	}
	if k.scan.Account != "" {
		scan := k.scan
		s.Scan = &scan
	}
	cat := k.cat
	k.mu.RUnlock()
	if cat != nil {
		s.Ships = cat.List()
	}
	return s
}

// Restore reloads durable state. Callers must not hold k.mu.
func (k *Kernel) Restore(s Snapshot) {
	k.mu.Lock()
	k.closed = s.Closed
	k.notes = append([]Note(nil), s.Notes...)
	k.facts = append([]Fact(nil), s.Facts...)
	k.confirms = append([]Confirm(nil), s.Confirms...)
	k.mail = append([]Mail(nil), s.Mail...)
	k.pulses = append([]Pulse(nil), s.Pulses...)
	k.choices = append([]Choice(nil), s.Choices...)
	k.requests = append([]Request(nil), s.Requests...)
	if s.Scan != nil {
		k.scan = *s.Scan
	}
	k.noteSeq = maxID(s.Notes, func(n Note) int { return n.ID })
	k.factSeq = maxID(s.Facts, func(f Fact) int { return f.ID })
	k.confSeq = maxID(s.Confirms, func(c Confirm) int { return c.ID })
	k.mailSeq = maxID(s.Mail, func(m Mail) int { return m.ID })
	k.pulseSeq = maxID(s.Pulses, func(p Pulse) int { return p.ID })
	k.choiceSeq = maxID(s.Choices, func(c Choice) int { return c.ID })
	k.reqSeq = maxID(s.Requests, func(r Request) int { return r.ID })
	cat := k.cat
	k.mu.Unlock()
	if cat != nil && len(s.Ships) > 0 {
		cat.Replace(s.Ships)
	}
}

// SaveFile writes the OS image as JSON.
func SaveFile(path string, k *Kernel) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	data, err := json.MarshalIndent(k.Snapshot(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// LoadFile restores the OS image. Missing files are ignored.
func LoadFile(path string, k *Kernel) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	k.Restore(s)
	return nil
}

func maxID[T any](items []T, id func(T) int) int {
	n := 0
	for _, item := range items {
		if v := id(item); v > n {
			n = v
		}
	}
	return n
}
