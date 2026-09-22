package ultron

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"
)

// Boot seeds companies, fleet, and owner if the plane is empty.
func (p *Plane) Boot() error {
	p.mu.Lock()
	empty := len(p.companies) == 0
	p.mu.Unlock()
	if empty {
		p.seedCompanies()
		p.seedFleet()
	}
	p.mu.Lock()
	needOwner := len(p.users) == 0
	p.mu.Unlock()
	if needOwner {
		if err := p.seedOwner(); err != nil {
			return err
		}
	}
	if err := p.seedGiantClient(); err != nil {
		return err
	}
	p.Audit("ultron", "boot", Version, true, Name+" "+Version)
	p.persist()
	return nil
}

func (p *Plane) seedOwner() error {
	pw := p.bootSecret
	if pw == "" {
		pw = "ultron-change-me"
	}
	u := &User{
		ID:         "castro",
		Email:      "alejandro@proximityagency.ca",
		Name:       "Castro",
		Role:       RoleOwner,
		CompanyIDs: nil, // owner sees all
		Active:     true,
		CreatedAt:  p.now(),
	}
	if err := p.SetPassword(u, pw); err != nil {
		return err
	}
	p.mu.Lock()
	p.users[u.ID] = u
	p.bootSecret = ""
	p.mu.Unlock()
	return nil
}

func (p *Plane) seedGiantClient() error {
	const email = "client@giant.local"
	p.mu.RLock()
	for _, u := range p.users {
		if strings.EqualFold(u.Email, email) {
			p.mu.RUnlock()
			return nil
		}
	}
	p.mu.RUnlock()
	u := &User{
		ID:         "giant-client",
		Email:      email,
		Name:       "Giant Client",
		Role:       RoleClient,
		CompanyIDs: []string{"giant"},
		Active:     true,
		CreatedAt:  p.now(),
	}
	if err := p.SetPassword(u, "giant-client-1"); err != nil {
		return err
	}
	p.mu.Lock()
	p.users[u.ID] = u
	p.mu.Unlock()
	p.Audit("ultron", "user.seed", u.ID, true, "giant client tenancy")
	return nil
}

func (p *Plane) seedCompanies() {
	now := p.now()
	items := []Company{
		{ID: "giant", Name: "Giant", Slug: "giant", Sector: "platform", Status: "active", Stack: []string{"Go", "TypeScript"}, Notes: "Flagship multi-tenant platform company.", AgentIDs: []string{"giant-conductor", "giant-deploy"}, WorkerIDs: []string{"giant-ci"}},
		{ID: "scanapp", Name: "ScanApp", Slug: "scanapp", Sector: "ops", Status: "active", Stack: []string{"TypeScript", "Python"}, Notes: "Scan, CRM, and AI bots.", AgentIDs: []string{"scan-ingest", "scan-qa"}, WorkerIDs: []string{"scan-worker"}},
		{ID: "proximity", Name: "Proximity", Slug: "proximity", Sector: "agency", Status: "active", Stack: []string{"Next.js", "TypeScript", "Azure"}, RepoURL: "https://github.com/Evolu-Jeunes", Notes: "Agency platform + client delivery.", AgentIDs: []string{"prox-delivery", "prox-reviewer"}, WorkerIDs: []string{"prox-azure"}},
		{ID: "empire", Name: "Empire", Slug: "empire", Sector: "holding", Status: "incubating", Stack: []string{"Go", "Next.js"}, Notes: "Holding / empire orchestration layer.", AgentIDs: []string{"empire-planner"}, WorkerIDs: []string{}},
		{ID: "evolu-jeunes", Name: "Evolu-Jeunes", Slug: "evolu-jeunes", Sector: "org", Status: "active", Stack: []string{"WordPress", "TypeScript"}, RepoURL: "https://github.com/Evolu-Jeunes", Notes: "Org that ships many client mandates.", AgentIDs: []string{"ej-delivery"}, WorkerIDs: []string{}},
		{ID: "cashtro", Name: "Cashtro", Slug: "cashtro", Sector: "os", Status: "active", Stack: []string{"Go"}, RepoURL: "https://github.com/cashtro/cashtro", Notes: "Subordinate agentic OS kernel bridged by Ultron.", AgentIDs: []string{"cashtro-bridge"}, WorkerIDs: []string{"cashtro-http"}},
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range items {
		c.CreatedAt = now
		c.UpdatedAt = now
		cp := c
		p.companies[cp.ID] = &cp
	}
}

func (p *Plane) seedFleet() {
	now := p.now()
	agents := []Agent{
		{ID: "giant-conductor", Name: "Giant Conductor", Kind: "user", Role: "orchestrator", CompanyID: "giant", Mode: "live", Capabilities: []string{"giant.plan", "giant.ship"}, Status: "running", Summary: "Runs Giant product loops from Ultron."},
		{ID: "giant-deploy", Name: "Giant Deploy", Kind: "user", Role: "release", CompanyID: "giant", CashtroID: "deploy", Mode: "live", Capabilities: []string{"deploy.release"}, Status: "running", Summary: "Bridges Giant releases through Cashtro deploy."},
		{ID: "scan-ingest", Name: "Scan Ingest", Kind: "user", Role: "ingest", CompanyID: "scanapp", Mode: "live", Capabilities: []string{"scan.ingest"}, Status: "running", Summary: "Ingests ScanApp jobs."},
		{ID: "scan-qa", Name: "Scan QA", Kind: "user", Role: "qa", CompanyID: "scanapp", CashtroID: "reviewer", Mode: "live", Capabilities: []string{"reviewer.watch"}, Status: "running", Summary: "QA via Cashtro reviewer."},
		{ID: "prox-delivery", Name: "Proximity Delivery", Kind: "user", Role: "ship", CompanyID: "proximity", CashtroID: "delivery", Mode: "live", Capabilities: []string{"delivery.list", "delivery.advance"}, Status: "running", Summary: "Delivery line for Proximity mandates."},
		{ID: "prox-reviewer", Name: "Proximity Reviewer", Kind: "user", Role: "qa", CompanyID: "proximity", CashtroID: "reviewer", Mode: "live", Capabilities: []string{"reviewer.watch"}, Status: "running", Summary: "Desk evidence before ship."},
		{ID: "empire-planner", Name: "Empire Planner", Kind: "user", Role: "backlog", CompanyID: "empire", CashtroID: "planner", Mode: "live", Capabilities: []string{"planner.backlog"}, Status: "running", Summary: "Turns empire goals into idea ships."},
		{ID: "ej-delivery", Name: "EJ Delivery", Kind: "user", Role: "ship", CompanyID: "evolu-jeunes", CashtroID: "delivery", Mode: "live", Capabilities: []string{"delivery.list"}, Status: "running", Summary: "Evolu-Jeunes board view."},
		{ID: "cashtro-bridge", Name: "Cashtro Bridge", Kind: "system", Role: "bridge", CompanyID: "cashtro", Mode: "live", Capabilities: []string{"os.about", "os.heartbeat"}, Status: "running", Summary: "Health + heartbeat into Cashtro OS."},
	}
	workers := []Worker{
		{ID: "giant-ci", Name: "Giant CI", Kind: "ci", CompanyID: "giant", Bound: false, Status: "resident", Notes: "Bind CI webhook when ready."},
		{ID: "scan-worker", Name: "Scan Worker", Kind: "compute", CompanyID: "scanapp", Bound: false, Status: "resident", Notes: "Bind ScanApp worker pool."},
		{ID: "prox-azure", Name: "Proximity Azure", Kind: "cloud", CompanyID: "proximity", Bound: false, Status: "resident", Notes: "Azure APIs for agency platform."},
		{ID: "cashtro-http", Name: "Cashtro HTTP", Kind: "os", CompanyID: "cashtro", Endpoint: "http://127.0.0.1:8080", Bound: true, Status: "running", Notes: "Subordinate Cashtro OS always-on."},
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, a := range agents {
		a.UpdatedAt = now
		cp := a
		p.agents[cp.ID] = &cp
	}
	for _, w := range workers {
		w.UpdatedAt = now
		cp := w
		p.workers[cp.ID] = &cp
	}
}

// ListCompanies returns companies visible to the user.
func (p *Plane) ListCompanies(u User) []Company {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Company, 0, len(p.companies))
	for _, c := range p.companies {
		if AllowsCompany(u, c.ID) {
			out = append(out, *c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// GetCompany returns one company if allowed.
func (p *Plane) GetCompany(u User, id string) (Company, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	c, ok := p.companies[id]
	if !ok {
		return Company{}, ErrNotFound
	}
	if !AllowsCompany(u, id) {
		return Company{}, ErrForbidden
	}
	return *c, nil
}

// UpsertCompany creates or updates a company.
func (p *Plane) UpsertCompany(actor User, in Company) (Company, error) {
	if !Can(actor.Role, PermCompanyWrite) {
		return Company{}, ErrForbidden
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Company{}, ErrInvalid
	}
	now := p.now()
	p.mu.Lock()
	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = slugify(name)
	}
	if id == "" {
		return Company{}, ErrInvalid
	}
	existing, ok := p.companies[id]
	if !ok {
		c := Company{
			ID: id, Name: name, Slug: id, Sector: in.Sector, Status: orDefault(in.Status, "incubating"),
			Stack: append([]string(nil), in.Stack...), RepoURL: in.RepoURL, Notes: in.Notes,
			AgentIDs: append([]string(nil), in.AgentIDs...), WorkerIDs: append([]string(nil), in.WorkerIDs...),
			CreatedAt: now, UpdatedAt: now,
		}
		p.companies[id] = &c
		out := c
		p.mu.Unlock()
		p.Audit(actor.Email, "company.create", id, true, name)
		p.persist()
		return out, nil
	}
	existing.Name = name
	if in.Sector != "" {
		existing.Sector = in.Sector
	}
	if in.Status != "" {
		existing.Status = in.Status
	}
	if in.Stack != nil {
		existing.Stack = append([]string(nil), in.Stack...)
	}
	if in.RepoURL != "" {
		existing.RepoURL = in.RepoURL
	}
	if in.Notes != "" {
		existing.Notes = in.Notes
	}
	existing.UpdatedAt = now
	out := *existing
	p.mu.Unlock()
	p.Audit(actor.Email, "company.update", id, true, name)
	p.persist()
	return out, nil
}

// ListAgents returns fleet agents, optionally company-filtered.
func (p *Plane) ListAgents(u User, companyID string) []Agent {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Agent, 0, len(p.agents))
	for _, a := range p.agents {
		if companyID != "" && a.CompanyID != companyID {
			continue
		}
		if !AllowsCompany(u, a.CompanyID) {
			continue
		}
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ListWorkers returns workers visible to the user.
func (p *Plane) ListWorkers(u User, companyID string) []Worker {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Worker, 0, len(p.workers))
	for _, w := range p.workers {
		if companyID != "" && w.CompanyID != companyID {
			continue
		}
		if !AllowsCompany(u, w.CompanyID) {
			continue
		}
		out = append(out, *w)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// About returns the Ultron card.
func (p *Plane) About(cashtroOK bool) About {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return About{
		Name:       Name,
		Version:    Version,
		Motto:      Motto,
		Companies:  len(p.companies),
		Agents:     len(p.agents),
		Workers:    len(p.workers),
		Users:      len(p.users),
		CashtroURL: p.cashtroURL,
		CashtroOK:  cashtroOK,
		Manifesto:  "Ultron is the personal IDE OS above Cashtro. Run Giant, ScanApp, Proximity, Empire and every company from one authorized desk. Clients get RBAC. Cashtro stays the agentic kernel — Ultron commands it over HTTP. No Cursor token required on the server.",
	}
}

// Audit appends a control-plane event.
func (p *Plane) Audit(actor, action, target string, ok bool, detail string) {
	p.mu.Lock()
	p.auditSeq++
	line := AuditLine{ID: p.auditSeq, At: p.now(), Actor: actor, Action: action, Target: target, OK: ok, Detail: detail}
	p.audit = append(p.audit, line)
	if len(p.audit) > 500 {
		p.audit = append([]AuditLine(nil), p.audit[len(p.audit)-500:]...)
	}
	p.mu.Unlock()
}

// AuditLog returns newest-first lines.
func (p *Plane) AuditLog() []AuditLine {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]AuditLine, len(p.audit))
	copy(out, p.audit)
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out
}

func orDefault(v, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		r = fold(r)
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return fmt.Sprintf("c-%d", time.Now().Unix()%100000)
	}
	return out
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
	case 'ò', 'ó', 'ô', 'ö', 'õ':
		return 'o'
	case 'ù', 'ú', 'û', 'ü':
		return 'u'
	case 'ñ':
		return 'n'
	}
	return unicode.ToLower(r)
}
