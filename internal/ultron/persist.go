package ultron

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Snapshot is the durable Ultron image.
type Snapshot struct {
	Version   string     `json:"version"`
	Users     []User     `json:"users"`
	Sessions  []Session  `json:"sessions"`
	Companies []Company  `json:"companies"`
	Agents    []Agent    `json:"agents"`
	Workers   []Worker   `json:"workers"`
	Audit     []AuditLine `json:"audit,omitempty"`
}

func (p *Plane) persist() {
	path := p.PersistPath()
	if path == "" {
		return
	}
	_ = SaveFile(path, p)
}

// Snapshot copies durable state.
func (p *Plane) Snapshot() Snapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s := Snapshot{Version: Version}
	for _, u := range p.users {
		s.Users = append(s.Users, *u)
	}
	for _, sess := range p.sessions {
		s.Sessions = append(s.Sessions, *sess)
	}
	for _, c := range p.companies {
		s.Companies = append(s.Companies, *c)
	}
	for _, a := range p.agents {
		s.Agents = append(s.Agents, *a)
	}
	for _, w := range p.workers {
		s.Workers = append(s.Workers, *w)
	}
	s.Audit = append([]AuditLine(nil), p.audit...)
	return s
}

// Restore reloads durable state.
func (p *Plane) Restore(s Snapshot) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.users = make(map[string]*User)
	p.sessions = make(map[string]*Session)
	p.companies = make(map[string]*Company)
	p.agents = make(map[string]*Agent)
	p.workers = make(map[string]*Worker)
	for i := range s.Users {
		u := s.Users[i]
		p.users[u.ID] = &u
	}
	for i := range s.Sessions {
		sess := s.Sessions[i]
		p.sessions[sess.Token] = &sess
	}
	for i := range s.Companies {
		c := s.Companies[i]
		p.companies[c.ID] = &c
	}
	for i := range s.Agents {
		a := s.Agents[i]
		p.agents[a.ID] = &a
	}
	for i := range s.Workers {
		w := s.Workers[i]
		p.workers[w.ID] = &w
	}
	p.audit = append([]AuditLine(nil), s.Audit...)
	max := 0
	for _, line := range p.audit {
		if line.ID > max {
			max = line.ID
		}
	}
	p.auditSeq = max
}

// SaveFile writes Ultron JSON.
func SaveFile(path string, p *Plane) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	data, err := json.MarshalIndent(p.Snapshot(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

// LoadFile restores Ultron JSON.
func LoadFile(path string, p *Plane) error {
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
	p.Restore(s)
	return nil
}
