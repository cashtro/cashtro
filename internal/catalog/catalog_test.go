package catalog

import (
	"errors"
	"testing"
	"time"
)

func TestSeededCatalog(t *testing.T) {
	c := New()
	ships := c.List()
	if len(ships) != 10 {
		t.Fatalf("list len = %d, want 10", len(ships))
	}
	giant, err := c.Get("giant")
	if err != nil {
		t.Fatal(err)
	}
	if giant.Stage != StageProduction || giant.Client != "Giant" {
		t.Fatalf("giant = %+v", giant)
	}
	if ships[0].Stage != StageIdea {
		t.Fatalf("first stage = %s, want idea", ships[0].Stage)
	}

	got, err := c.Get("proximity")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Proximity" || got.Stage != StageProduction {
		t.Fatalf("proximity = %+v", got)
	}

	_, err = c.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing err = %v, want ErrNotFound", err)
	}
}

func TestCreateAndAdvance(t *testing.T) {
	now := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	c := New(WithClock(func() time.Time { return now }))

	ship, err := c.Create(CreateShip{
		Name:   "  Ledger desk  ",
		Client: "Solution Hypothèque QC",
		Sector: "finance",
		Stack:  []string{"Go", " Go ", "", "Azure"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ship.ID != "ledger-desk" {
		t.Fatalf("id = %q", ship.ID)
	}
	if ship.Stage != StageIdea {
		t.Fatalf("stage = %s, want idea", ship.Stage)
	}
	if len(ship.Stack) != 2 || ship.Stack[0] != "Go" || ship.Stack[1] != "Azure" {
		t.Fatalf("stack = %#v", ship.Stack)
	}

	dup, err := c.Create(CreateShip{Name: "Ledger desk"})
	if err != nil {
		t.Fatal(err)
	}
	if dup.ID != "ledger-desk-2" {
		t.Fatalf("dup id = %q", dup.ID)
	}

	next, err := c.Advance(ship.ID)
	if err != nil {
		t.Fatal(err)
	}
	if next.Stage != StageConcept {
		t.Fatalf("after first advance = %s", next.Stage)
	}

	prod, err := c.Advance(ship.ID)
	if err != nil {
		t.Fatal(err)
	}
	if prod.Stage != StageProduction {
		t.Fatalf("after second advance = %s", prod.Stage)
	}

	_, err = c.Advance(ship.ID)
	if !errors.Is(err, ErrDone) {
		t.Fatalf("third advance err = %v, want ErrDone", err)
	}
}

func TestCreateRequiresName(t *testing.T) {
	c := New()
	_, err := c.Create(CreateShip{Name: "   "})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("err = %v, want ErrInvalid", err)
	}
}

func TestNextStage(t *testing.T) {
	got, err := NextStage(StageIdea)
	if err != nil || got != StageConcept {
		t.Fatalf("idea → %s (%v)", got, err)
	}
	got, err = NextStage(StageConcept)
	if err != nil || got != StageProduction {
		t.Fatalf("concept → %s (%v)", got, err)
	}
	_, err = NextStage(StageProduction)
	if !errors.Is(err, ErrDone) {
		t.Fatalf("production err = %v", err)
	}
	_, err = NextStage("draft")
	if !errors.Is(err, ErrUnknownStage) {
		t.Fatalf("unknown err = %v", err)
	}
}

func TestProfile(t *testing.T) {
	p := New().Profile()
	if p.Name != "Castro" {
		t.Fatalf("name = %q", p.Name)
	}
	if p.Email != "alejandro@proximityagency.ca" {
		t.Fatalf("email = %q", p.Email)
	}
	if len(p.Stack) == 0 || p.Stack[0] != "Go" {
		t.Fatalf("stack should lead with Go: %#v", p.Stack)
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"BTK Avocats": "btk-avocats",
		"  Édu 2026 ": "edu-2026",
		"!!!":         "",
		"ScanApp":     "scanapp",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Fatalf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}
