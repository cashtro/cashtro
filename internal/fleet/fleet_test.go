package fleet

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestSeededGiantEcosystem(t *testing.T) {
	f := New()
	cos := f.List()
	if len(cos) != 6 {
		t.Fatalf("companies = %d, want 6", len(cos))
	}
	found := map[string]bool{}
	for _, c := range cos {
		found[c.ID] = true
	}
	for _, id := range []string{"giant", "scanapp", "proximity", "empire", "evolu-jeunes", "cashtro"} {
		if !found[id] {
			t.Fatalf("missing %s", id)
		}
	}
	giant, err := f.Get("giant")
	if err != nil || giant.Name != "Giant" {
		t.Fatalf("giant = %+v %v", giant, err)
	}
	own, err := f.OwnAgents("giant")
	if err != nil || len(own) != 2 {
		t.Fatalf("giant agents = %d %v", len(own), err)
	}
	if f.UsingMine("giant") {
		t.Fatal("giant has its own roster")
	}
	if !f.UsingMine("cashtro") {
		t.Fatal("cashtro should use Castro's mine")
	}
}

func TestCreateUsesMineWhenNoAgentics(t *testing.T) {
	now := time.Date(2026, 9, 21, 22, 0, 0, 0, time.UTC)
	f := New(WithClock(func() time.Time { return now }))

	co, err := f.Create(CreateCompany{Name: "Northwind", Sector: "ops", Notes: "new tenant"})
	if err != nil {
		t.Fatal(err)
	}
	if co.ID != "northwind" || co.Mark != "N" || !co.UseMine {
		t.Fatalf("created = %+v", co)
	}
	if !f.UsingMine(co.ID) {
		t.Fatal("empty roster must use Castro's")
	}
	own, err := f.OwnAgents(co.ID)
	if err != nil || len(own) != 0 {
		t.Fatalf("own = %d %v", len(own), err)
	}

	_, err = f.Create(CreateCompany{Name: "Northwind"})
	if !errors.Is(err, ErrExists) {
		t.Fatalf("dup err = %v", err)
	}

	_, err = f.Create(CreateCompany{Name: "   "})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("empty err = %v", err)
	}
}

func TestAddAgenticAndKeepMine(t *testing.T) {
	f := New()
	co, err := f.Create(CreateCompany{Name: "Ledger"})
	if err != nil {
		t.Fatal(err)
	}
	ag, err := f.AddAgentic(co.ID, CreateAgentic{
		Name: "Ledger Clerk", Role: "ops", CashtroID: "delivery",
		Capabilities: []string{"delivery.list"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ag.ID != "ledger-clerk" || ag.CashtroID != "delivery" || ag.Mode != "live" {
		t.Fatalf("agentic = %+v", ag)
	}
	got, err := f.Get(co.ID)
	if err != nil || len(got.AgentIDs) != 1 {
		t.Fatalf("company after add = %+v %v", got, err)
	}
	if !f.UsingMine(co.ID) {
		t.Fatal("UseMine stays on so they still have Castro's as backup")
	}

	mineOff := false
	ownOnly, err := f.Create(CreateCompany{
		Name: "Solo Co", UseMine: &mineOff,
		Agents: []CreateAgentic{{Name: "Solo Bot", Role: "bot"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ownOnly.UseMine || f.UsingMine(ownOnly.ID) {
		t.Fatalf("solo should not inherit: %+v", ownOnly)
	}
}

func TestImageRestoreKeepsAddedCompany(t *testing.T) {
	f := New()
	if _, err := f.Create(CreateCompany{Name: "West Desk", Sector: "ops"}); err != nil {
		t.Fatal(err)
	}
	img := f.Image()
	n, mine := f.Count()
	if n < 7 || mine < 1 {
		t.Fatalf("count = %d mine=%d", n, mine)
	}

	f2 := New()
	if err := f2.RestoreTenants(mustJSON(t, img)); err != nil {
		t.Fatal(err)
	}
	if _, err := f2.Get("west-desk"); err != nil {
		t.Fatal("west desk should restore")
	}
	if f2.UsingMine("west-desk") != true {
		t.Fatal("empty west desk should use Castro's")
	}
}

func mustJSON(t *testing.T, img Image) []byte {
	t.Helper()
	raw, err := json.Marshal(img)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestMarkAndSlug(t *testing.T) {
	if got := markOf("Evolu-Jeunes"); got != "EJ" {
		t.Fatalf("mark Evolu-Jeunes = %q", got)
	}
	if got := markOf("Giant"); got != "G" {
		t.Fatalf("mark Giant = %q", got)
	}
	if got := slugify("North Desk"); got != "north-desk" {
		t.Fatalf("slug = %q", got)
	}
}
