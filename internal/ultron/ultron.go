// Package ultron is the personal control plane above Cashtro OS.
//
// Ultron is a separated IDE / company OS dashboard. Clients authenticate with
// RBAC. Companies (Giant, ScanApp, Proximity, Empire, …) are first-class
// tenants. Cashtro OS on :8080 remains the subordinate agentic kernel —
// Ultron bridges to it; it does not fork a second agent runtime.
package ultron

import (
	"sync"
	"time"
)

const (
	// Name is the public control-plane brand.
	Name = "Ultron"
	// Version is the Ultron release.
	Version = "0.1.1"
	// Motto is the operator line.
	Motto = "One IDE. Every company. Authority with RBAC."
)

// Role is an RBAC role.
type Role string

const (
	RoleOwner    Role = "owner"
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleClient   Role = "client"
	RoleViewer   Role = "viewer"
)

// Perm is a named permission.
type Perm string

const (
	PermIDEAccess     Perm = "ide.access"
	PermCompanyRead   Perm = "company.read"
	PermCompanyWrite  Perm = "company.write"
	PermFleetRead     Perm = "fleet.read"
	PermFleetInvoke   Perm = "fleet.invoke"
	PermOSBridge      Perm = "os.bridge"
	PermUsersRead     Perm = "users.read"
	PermUsersWrite    Perm = "users.write"
	PermAuditRead     Perm = "audit.read"
)

// User is an Ultron principal.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         Role      `json:"role"`
	CompanyIDs   []string  `json:"companyIds"`
	PasswordHash string    `json:"passwordHash,omitempty"`
	Salt         string    `json:"salt,omitempty"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Public strips secrets.
func (u User) Public() User {
	u.PasswordHash = ""
	u.Salt = ""
	return u
}

// Company is one business Ultron runs (Giant, ScanApp, …).
type Company struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Sector      string    `json:"sector"`
	Status      string    `json:"status"` // active | incubating | paused
	Stack       []string  `json:"stack"`
	RepoURL     string    `json:"repoUrl,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	AgentIDs    []string  `json:"agentIds"`
	WorkerIDs   []string  `json:"workerIds"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Agent is a logical agentic bound to one or more companies.
type Agent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Kind        string    `json:"kind"` // system | user | worker
	Role        string    `json:"role"`
	CompanyID   string    `json:"companyId"`
	CashtroID   string    `json:"cashtroId,omitempty"` // bridge id on Cashtro OS
	Mode        string    `json:"mode"`                 // live | resident
	Capabilities []string `json:"capabilities"`
	Status      string    `json:"status"`
	Summary     string    `json:"summary"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Worker is an execution bind (browser, deploy, model, CI).
type Worker struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	CompanyID string    `json:"companyId"`
	Endpoint  string    `json:"endpoint,omitempty"`
	Bound     bool      `json:"bound"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Session is a bearer login.
type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// AuditLine is one authorization / control event.
type AuditLine struct {
	ID        int       `json:"id"`
	At        time.Time `json:"at"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Target    string    `json:"target,omitempty"`
	OK        bool      `json:"ok"`
	Detail    string    `json:"detail,omitempty"`
}

// About is the public Ultron card.
type About struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Motto      string `json:"motto"`
	Companies  int    `json:"companies"`
	Agents     int    `json:"agents"`
	Workers    int    `json:"workers"`
	Users      int    `json:"users"`
	CashtroURL string `json:"cashtroUrl"`
	CashtroOK  bool   `json:"cashtroOk"`
	Manifesto  string `json:"manifesto"`
}

// Plane is the Ultron control plane state.
type Plane struct {
	mu          sync.RWMutex
	now         func() time.Time
	persistPath string
	cashtroURL  string
	users       map[string]*User
	sessions    map[string]*Session
	companies   map[string]*Company
	agents      map[string]*Agent
	workers     map[string]*Worker
	audit       []AuditLine
	auditSeq    int
	bootSecret  string // owner password at first boot (cleared after seed persist)
}

// Option configures Ultron.
type Option func(*Plane)

// WithClock injects a clock.
func WithClock(now func() time.Time) Option {
	return func(p *Plane) { p.now = now }
}

// WithPersistPath sets the Ultron disk image.
func WithPersistPath(path string) Option {
	return func(p *Plane) { p.persistPath = path }
}

// WithCashtroURL sets the subordinate Cashtro OS base URL.
func WithCashtroURL(url string) Option {
	return func(p *Plane) { p.cashtroURL = url }
}

// WithBootOwnerPassword sets the initial owner password (env-friendly).
func WithBootOwnerPassword(pw string) Option {
	return func(p *Plane) { p.bootSecret = pw }
}

// New builds an empty plane; call Boot to seed.
func New(opts ...Option) *Plane {
	p := &Plane{
		now:        func() time.Time { return time.Now().UTC() },
		cashtroURL: "http://127.0.0.1:8080",
		users:      make(map[string]*User),
		sessions:   make(map[string]*Session),
		companies:  make(map[string]*Company),
		agents:     make(map[string]*Agent),
		workers:    make(map[string]*Worker),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// PersistPath returns the disk path.
func (p *Plane) PersistPath() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.persistPath
}

// CashtroURL returns the bridge target.
func (p *Plane) CashtroURL() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.cashtroURL
}
