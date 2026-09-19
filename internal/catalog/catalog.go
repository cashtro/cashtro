package catalog

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

// Stage is one step in the delivery line Castro runs: idea → concept → production.
type Stage string

const (
	StageIdea       Stage = "idea"
	StageConcept    Stage = "concept"
	StageProduction Stage = "production"
)

var stageOrder = []Stage{StageIdea, StageConcept, StageProduction}

// Profile is the public Cashtro Builder card.
type Profile struct {
	Name      string   `json:"name"`
	Role      string   `json:"role"`
	Motto     string   `json:"motto"`
	Languages []string `json:"languages"`
	Stack     []string `json:"stack"`
	Email     string   `json:"email"`
	Org       string   `json:"org"`
	OrgURL    string   `json:"orgUrl"`
	Team      string   `json:"team"`
}

// Ship is a mandate moving through the delivery line.
type Ship struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Client    string    `json:"client"`
	Sector    string    `json:"sector"`
	Stage     Stage     `json:"stage"`
	Stack     []string  `json:"stack"`
	URL       string    `json:"url,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateShip is the input for a new mandate.
type CreateShip struct {
	Name   string   `json:"name"`
	Client string   `json:"client"`
	Sector string   `json:"sector"`
	Stack  []string `json:"stack"`
	URL    string   `json:"url"`
	Notes  string   `json:"notes"`
}

var (
	// ErrNotFound is returned when a ship id is unknown.
	ErrNotFound = errors.New("ship not found")
	// ErrDone is returned when a ship is already in production.
	ErrDone = errors.New("already in production")
	// ErrInvalid is returned for a ship that cannot be created.
	ErrInvalid = errors.New("invalid ship")
	// ErrUnknownStage is returned when a stage is not on the line.
	ErrUnknownStage = errors.New("unknown stage")
)

// Catalog is an in-memory delivery board.
type Catalog struct {
	mu    sync.RWMutex
	seq   int
	now   func() time.Time
	ships map[string]*Ship
}

// Option configures a catalog.
type Option func(*Catalog)

// WithClock injects a clock for tests.
func WithClock(now func() time.Time) Option {
	return func(c *Catalog) {
		c.now = now
	}
}

// New returns a catalog seeded with live Evolu-Jeunes / Proximity work.
func New(opts ...Option) *Catalog {
	c := &Catalog{
		now:   func() time.Time { return time.Now().UTC() },
		ships: make(map[string]*Ship),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.seed()
	return c
}

// Profile returns the public builder card.
func (c *Catalog) Profile() Profile {
	return Profile{
		Name:      "Castro",
		Role:      "Scrum Master & Software Engineer · Founder",
		Motto:     "I make teams ship: idea → concept → production.",
		Languages: []string{"FR", "EN"},
		Stack:     []string{"Go", "PHP/WordPress", "React/Next.js/TypeScript", "Python", "Azure"},
		Email:     "alejandro@proximityagency.ca",
		Org:       "Evolu-Jeunes",
		OrgURL:    "https://github.com/Evolu-Jeunes",
		Team:      "75+ developers trained on real client mandates",
	}
}

// List returns ships ordered by stage then name.
func (c *Catalog) List() []Ship {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]Ship, 0, len(c.ships))
	for _, s := range c.ships {
		out = append(out, cloneShip(s))
	}
	sort.Slice(out, func(i, j int) bool {
		if stageIndex(out[i].Stage) != stageIndex(out[j].Stage) {
			return stageIndex(out[i].Stage) < stageIndex(out[j].Stage)
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// Get returns one ship by id.
func (c *Catalog) Get(id string) (Ship, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	s, ok := c.ships[id]
	if !ok {
		return Ship{}, ErrNotFound
	}
	return cloneShip(s), nil
}

// Create adds a new mandate at the idea stage.
func (c *Catalog) Create(in CreateShip) (Ship, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Ship{}, fmt.Errorf("%w: name is required", ErrInvalid)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.uniqueID(name)
	now := c.now()
	s := &Ship{
		ID:        id,
		Name:      name,
		Client:    strings.TrimSpace(in.Client),
		Sector:    strings.TrimSpace(in.Sector),
		Stage:     StageIdea,
		Stack:     cleanStack(in.Stack),
		URL:       strings.TrimSpace(in.URL),
		Notes:     strings.TrimSpace(in.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	}
	c.ships[id] = s
	return cloneShip(s), nil
}

// Advance moves a ship one stage forward.
func (c *Catalog) Advance(id string) (Ship, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	s, ok := c.ships[id]
	if !ok {
		return Ship{}, ErrNotFound
	}
	next, err := NextStage(s.Stage)
	if err != nil {
		return Ship{}, err
	}
	s.Stage = next
	s.UpdatedAt = c.now()
	return cloneShip(s), nil
}

// NextStage returns the following stage on the line.
func NextStage(current Stage) (Stage, error) {
	i := stageIndex(current)
	if i < 0 {
		return "", fmt.Errorf("%w: %s", ErrUnknownStage, current)
	}
	if i == len(stageOrder)-1 {
		return "", ErrDone
	}
	return stageOrder[i+1], nil
}

// Stages returns the delivery line in order.
func Stages() []Stage {
	out := make([]Stage, len(stageOrder))
	copy(out, stageOrder)
	return out
}

func (c *Catalog) seed() {
	fixed := time.Date(2026, 7, 22, 21, 14, 35, 0, time.UTC)
	items := []Ship{
		{
			ID:     "btk-avocats",
			Name:   "BTK Avocats",
			Client: "BTK Avocats",
			Sector: "law",
			Stage:  StageProduction,
			Stack:  []string{"WordPress"},
			Notes:  "Law firm platform",
		},
		{
			ID:     "md-clinic",
			Name:   "MD Clinic",
			Client: "MD Clinic",
			Sector: "health",
			Stage:  StageProduction,
			Stack:  []string{"WordPress", "ACF"},
			Notes:  "Medical clinic",
		},
		{
			ID:     "solution-hypotheque-qc",
			Name:   "Solution Hypothèque QC",
			Client: "Solution Hypothèque QC",
			Sector: "finance",
			Stage:  StageProduction,
			Stack:  []string{"WordPress"},
			Notes:  "Mortgage services",
		},
		{
			ID:     "educonnexion",
			Name:   "Éduconnexion",
			Client: "Éduconnexion",
			Sector: "education",
			Stage:  StageProduction,
			Stack:  []string{"WordPress"},
			Notes:  "Education platform",
		},
		{
			ID:     "proximity",
			Name:   "Proximity",
			Client: "Proximity Agency",
			Sector: "agency",
			Stage:  StageProduction,
			Stack:  []string{"Next.js", "TypeScript", "Azure"},
			Notes:  "Agency platform",
		},
		{
			ID:     "scanapp",
			Name:   "ScanApp",
			Client: "Evolu-Jeunes",
			Sector: "ops",
			Stage:  StageConcept,
			Stack:  []string{"TypeScript", "Python"},
			Notes:  "Scan, CRM, and AI bots",
		},
		{
			ID:     "cashtro-catalog",
			Name:   "Cashtro delivery catalog",
			Client: "Cashtro",
			Sector: "internal",
			Stage:  StageIdea,
			Stack:  []string{"Go"},
			Notes:  "Stdlib Go board for idea → concept → production",
		},
	}
	for i := range items {
		items[i].CreatedAt = fixed
		items[i].UpdatedAt = fixed
		copy := items[i]
		c.ships[copy.ID] = &copy
	}
}

func (c *Catalog) uniqueID(name string) string {
	base := slugify(name)
	if base == "" {
		c.seq++
		return fmt.Sprintf("ship-%d", c.seq)
	}
	if _, exists := c.ships[base]; !exists {
		return base
	}
	for n := 2; ; n++ {
		id := fmt.Sprintf("%s-%d", base, n)
		if _, exists := c.ships[id]; !exists {
			return id
		}
	}
}

func foldRune(r rune) rune {
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
	case 'ý', 'ÿ':
		return 'y'
	case 'œ':
		return 'o'
	default:
		return unicode.ToLower(r)
	}
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		r = foldRune(r)
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

func stageIndex(s Stage) int {
	for i, stage := range stageOrder {
		if stage == s {
			return i
		}
	}
	return -1
}

func cloneShip(s *Ship) Ship {
	out := *s
	if s.Stack != nil {
		out.Stack = append([]string(nil), s.Stack...)
	}
	return out
}
