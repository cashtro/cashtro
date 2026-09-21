// Package fleet is the Giant company plane on Cashtro OS.
//
// Companies are tenants. Each one can bring its own agentics. If it
// does not, it uses Castro's kernel agentics (mine). New companies
// register here so the desk can evolve without forking a second OS.
package fleet

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

// MineIDs are Castro's kernel agentics. A company with no roster
// inherits this set instead of shipping empty.
var MineIDs = []string{
	"delivery", "planner", "reviewer", "explorer", "research",
	"memory", "comms", "architect", "deploy", "security",
	"investigator", "operator",
}

var (
	// ErrNotFound is an unknown company or agentic.
	ErrNotFound = errors.New("not found")
	// ErrInvalid is a company or agentic that cannot be created.
	ErrInvalid = errors.New("invalid")
	// ErrExists is a duplicate company or agentic id.
	ErrExists = errors.New("already exists")
)

// Company is one tenant on the Giant desk.
type Company struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Mark      string    `json:"mark"`
	Sector    string    `json:"sector"`
	Status    string    `json:"status"`
	Stack     []string  `json:"stack"`
	Notes     string    `json:"notes,omitempty"`
	RepoURL   string    `json:"repoUrl,omitempty"`
	UseMine   bool      `json:"useMine"`
	AgentIDs  []string  `json:"agentIds"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Agentic is a company-owned process. CashtroID binds it to Castro's kernel.
type Agentic struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CompanyID    string    `json:"companyId"`
	CashtroID    string    `json:"cashtroId,omitempty"`
	Mode         string    `json:"mode"`
	Capabilities []string  `json:"capabilities"`
	Summary      string    `json:"summary"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// CreateCompany is the input for a new tenant.
type CreateCompany struct {
	Name    string          `json:"name"`
	Sector  string          `json:"sector"`
	Status  string          `json:"status"`
	Stack   []string        `json:"stack"`
	Notes   string          `json:"notes"`
	RepoURL string          `json:"repoUrl"`
	UseMine *bool           `json:"useMine"`
	Agents  []CreateAgentic `json:"agents"`
}

// CreateAgentic is the input for a company-owned agentic.
type CreateAgentic struct {
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Summary      string   `json:"summary"`
	CashtroID    string   `json:"cashtroId"`
	Mode         string   `json:"mode"`
	Capabilities []string `json:"capabilities"`
}

// Fleet is the in-memory company + agentic registry.
type Fleet struct {
	mu        sync.RWMutex
	now       func() time.Time
	seq       int
	companies map[string]*Company
	agents    map[string]*Agentic
}

// Option configures a fleet.
type Option func(*Fleet)

// WithClock injects a clock for tests.
func WithClock(now func() time.Time) Option {
	return func(f *Fleet) { f.now = now }
}

// New returns a fleet seeded with the Giant ecosystem.
func New(opts ...Option) *Fleet {
	f := &Fleet{
		now:       func() time.Time { return time.Now().UTC() },
		companies: make(map[string]*Company),
		agents:    make(map[string]*Agentic),
	}
	for _, opt := range opts {
		opt(f)
	}
	f.seed()
	return f
}

// List returns companies ordered by name.
func (f *Fleet) List() []Company {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]Company, 0, len(f.companies))
	for _, c := range f.companies {
		out = append(out, cloneCompany(c))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Get returns one company.
func (f *Fleet) Get(id string) (Company, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	c, ok := f.companies[id]
	if !ok {
		return Company{}, fmt.Errorf("%w: company %s", ErrNotFound, id)
	}
	return cloneCompany(c), nil
}

// OwnAgents returns company-owned agentics, name-sorted.
func (f *Fleet) OwnAgents(companyID string) ([]Agentic, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if _, ok := f.companies[companyID]; !ok {
		return nil, fmt.Errorf("%w: company %s", ErrNotFound, companyID)
	}
	out := make([]Agentic, 0)
	for _, a := range f.agents {
		if a.CompanyID == companyID {
			out = append(out, cloneAgentic(a))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// UsingMine reports whether the company falls back to Castro's agentics.
func (f *Fleet) UsingMine(id string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	c, ok := f.companies[id]
	if !ok {
		return true
	}
	if c.UseMine {
		return true
	}
	return len(c.AgentIDs) == 0
}

// GetAgentic returns one company-owned agentic.
func (f *Fleet) GetAgentic(id string) (Agentic, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	a, ok := f.agents[id]
	if !ok {
		return Agentic{}, fmt.Errorf("%w: agentic %s", ErrNotFound, id)
	}
	return cloneAgentic(a), nil
}

// Create adds a tenant. No roster means they use Castro's agentics.
func (f *Fleet) Create(in CreateCompany) (Company, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Company{}, fmt.Errorf("%w: name is required", ErrInvalid)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	id := slugify(name)
	if id == "" {
		f.seq++
		id = fmt.Sprintf("co-%d", f.seq)
	}
	if _, exists := f.companies[id]; exists {
		return Company{}, fmt.Errorf("%w: %s", ErrExists, id)
	}

	useMine := len(in.Agents) == 0
	if in.UseMine != nil {
		useMine = *in.UseMine
	}
	now := f.now()
	c := &Company{
		ID:        id,
		Name:      name,
		Mark:      markOf(name),
		Sector:    orDefault(in.Sector, "platform"),
		Status:    orDefault(in.Status, "incubating"),
		Stack:     cleanStack(in.Stack),
		Notes:     strings.TrimSpace(in.Notes),
		RepoURL:   strings.TrimSpace(in.RepoURL),
		UseMine:   useMine,
		CreatedAt: now,
		UpdatedAt: now,
	}
	f.companies[id] = c
	for _, a := range in.Agents {
		if _, err := f.addAgenticLocked(c, a); err != nil {
			return Company{}, err
		}
	}
	return cloneCompany(c), nil
}

// AddAgentic binds a company-owned agentic onto a tenant.
func (f *Fleet) AddAgentic(companyID string, in CreateAgentic) (Agentic, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c, ok := f.companies[companyID]
	if !ok {
		return Agentic{}, fmt.Errorf("%w: company %s", ErrNotFound, companyID)
	}
	return f.addAgenticLocked(c, in)
}

func (f *Fleet) addAgenticLocked(c *Company, in CreateAgentic) (Agentic, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Agentic{}, fmt.Errorf("%w: agentic name is required", ErrInvalid)
	}
	id := slugify(name)
	if id == "" {
		f.seq++
		id = fmt.Sprintf("%s-agent-%d", c.ID, f.seq)
	}
	if _, exists := f.agents[id]; exists {
		id = slugify(c.ID + "-" + name)
	}
	if id == "" || f.agents[id] != nil {
		return Agentic{}, fmt.Errorf("%w: %s", ErrExists, id)
	}
	now := f.now()
	mode := strings.TrimSpace(in.Mode)
	if mode == "" {
		if strings.TrimSpace(in.CashtroID) != "" {
			mode = "live"
		} else {
			mode = "resident"
		}
	}
	a := &Agentic{
		ID:           id,
		Name:         name,
		Role:         orDefault(in.Role, "agentic"),
		CompanyID:    c.ID,
		CashtroID:    strings.TrimSpace(in.CashtroID),
		Mode:         mode,
		Capabilities: cleanStack(in.Capabilities),
		Summary:      strings.TrimSpace(in.Summary),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if a.Summary == "" {
		if a.CashtroID != "" {
			a.Summary = "Bound to Castro's " + a.CashtroID + " on the kernel."
		} else {
			a.Summary = "Company agentic. Bind a Cashtro id to execute through Castro's kernel."
		}
	}
	f.agents[a.ID] = a
	c.AgentIDs = append(append([]string{}, c.AgentIDs...), a.ID)
	c.UpdatedAt = now
	return cloneAgentic(a), nil
}

func (f *Fleet) seed() {
	fixed := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	type seedAgent struct {
		id, name, role, cashtro, mode, summary string
		caps                                   []string
	}
	type seedCo struct {
		id, name, mark, sector, status, notes string
		stack                                 []string
		useMine                               bool
		agents                                []seedAgent
	}
	items := []seedCo{
		{
			id: "giant", name: "Giant", mark: "G", sector: "platform", status: "active",
			notes: "Flagship multi-tenant platform", stack: []string{"Go", "TypeScript"},
			agents: []seedAgent{
				{id: "giant-conductor", name: "Giant Conductor", role: "orchestrator", mode: "live",
					caps: []string{"giant.plan", "giant.ship"}, summary: "Runs Giant product loops from the desk."},
				{id: "giant-deploy", name: "Giant Deploy", role: "release", cashtro: "deploy", mode: "live",
					caps: []string{"deploy.release"}, summary: "Giant releases through Castro's deploy."},
			},
		},
		{
			id: "scanapp", name: "ScanApp", mark: "S", sector: "ops", status: "active",
			notes: "Scan, CRM, and AI bots", stack: []string{"TypeScript", "Python"},
			agents: []seedAgent{
				{id: "scan-ingest", name: "Scan Ingest", role: "ingest", mode: "live",
					caps: []string{"scan.ingest"}, summary: "Ingests ScanApp jobs."},
				{id: "scan-qa", name: "Scan QA", role: "qa", cashtro: "reviewer", mode: "live",
					caps: []string{"reviewer.watch"}, summary: "QA via Castro's reviewer."},
			},
		},
		{
			id: "proximity", name: "Proximity", mark: "P", sector: "agency", status: "active",
			notes: "Agency platform + client delivery", stack: []string{"Next.js", "TypeScript", "Azure"},
			agents: []seedAgent{
				{id: "prox-delivery", name: "Proximity Delivery", role: "ship", cashtro: "delivery", mode: "live",
					caps: []string{"delivery.list", "delivery.advance"}, summary: "Delivery line for Proximity mandates."},
				{id: "prox-reviewer", name: "Proximity Reviewer", role: "qa", cashtro: "reviewer", mode: "live",
					caps: []string{"reviewer.watch"}, summary: "Desk evidence before ship."},
			},
		},
		{
			id: "empire", name: "Empire", mark: "E", sector: "holding", status: "incubating",
			notes: "Holding / empire orchestration", stack: []string{"Go", "Next.js"},
			agents: []seedAgent{
				{id: "empire-planner", name: "Empire Planner", role: "backlog", cashtro: "planner", mode: "live",
					caps: []string{"planner.backlog"}, summary: "Turns empire goals into idea ships."},
			},
		},
		{
			id: "evolu-jeunes", name: "Evolu-Jeunes", mark: "EJ", sector: "org", status: "active",
			notes: "Org that ships many client mandates", stack: []string{"WordPress", "TypeScript"},
			agents: []seedAgent{
				{id: "ej-delivery", name: "EJ Delivery", role: "ship", cashtro: "delivery", mode: "live",
					caps: []string{"delivery.list"}, summary: "Evolu-Jeunes board view."},
			},
		},
		{
			id: "cashtro", name: "Cashtro", mark: "C", sector: "os", status: "active",
			notes: "Castro's agentic OS kernel — mine", stack: []string{"Go"},
			useMine: true,
			agents: []seedAgent{
				{id: "cashtro-bridge", name: "Cashtro Bridge", role: "bridge", cashtro: "init", mode: "live",
					caps: []string{"os.about"}, summary: "Health into Castro's kernel."},
			},
		},
	}
	for _, item := range items {
		c := &Company{
			ID: item.id, Name: item.name, Mark: item.mark, Sector: item.sector, Status: item.status,
			Stack: append([]string(nil), item.stack...), Notes: item.notes, UseMine: item.useMine,
			CreatedAt: fixed, UpdatedAt: fixed,
		}
		ids := make([]string, 0, len(item.agents))
		for _, sa := range item.agents {
			a := &Agentic{
				ID: sa.id, Name: sa.name, Role: sa.role, CompanyID: c.ID, CashtroID: sa.cashtro,
				Mode: sa.mode, Capabilities: append([]string(nil), sa.caps...), Summary: sa.summary,
				CreatedAt: fixed, UpdatedAt: fixed,
			}
			f.agents[a.ID] = a
			ids = append(ids, a.ID)
		}
		c.AgentIDs = ids
		f.companies[c.ID] = c
	}
}

func cloneCompany(c *Company) Company {
	out := *c
	out.Stack = append([]string(nil), c.Stack...)
	out.AgentIDs = append([]string(nil), c.AgentIDs...)
	return out
}

func cloneAgentic(a *Agentic) Agentic {
	out := *a
	out.Capabilities = append([]string(nil), a.Capabilities...)
	return out
}

func orDefault(v, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return strings.TrimSpace(v)
}

func markOf(name string) string {
	cleaned := strings.NewReplacer("-", " ", "_", " ").Replace(name)
	parts := strings.Fields(cleaned)
	var b strings.Builder
	for _, p := range parts {
		r := []rune(p)
		if len(r) == 0 {
			continue
		}
		b.WriteRune(unicode.ToUpper(r[0]))
		if b.Len() >= 2 {
			break
		}
	}
	if b.Len() == 0 {
		return "?"
	}
	return b.String()
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		r = fold(r)
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

func fold(r rune) rune {
	switch r {
	case 'à', 'á', 'â', 'ä', 'ã':
		return 'a'
	case 'ç':
		return 'c'
	case 'è', 'é', 'ê', 'ë':
		return 'e'
	case 'ì', 'í', 'î', 'ï':
		return 'i'
	case 'ñ':
		return 'n'
	case 'ò', 'ó', 'ô', 'ö', 'õ':
		return 'o'
	case 'ù', 'ú', 'û', 'ü':
		return 'u'
	default:
		return unicode.ToLower(r)
	}
}

func cleanStack(stack []string) []string {
	seen := make(map[string]struct{}, len(stack))
	out := make([]string, 0, len(stack))
	for _, item := range stack {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		key := strings.ToLower(item)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
	}
	return out
}
