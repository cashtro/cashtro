package agents

import (
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestVapiKeepsOneDatabasePerProject(t *testing.T) {
	desks := VapiDesks()
	if len(desks) != len(Lines()) {
		t.Fatalf("desks = %d lines = %d", len(desks), len(Lines()))
	}
	seen := map[string]bool{}
	for _, desk := range desks {
		if desk.Voice != "vapi" || desk.Shared || desk.Database == "" || seen[desk.Database] {
			t.Fatalf("desk %+v", desk)
		}
		if !strings.Contains(desk.Prompt, desk.Database) {
			t.Fatalf("prompt missing database: %s", desk.Prompt)
		}
		seen[desk.Database] = true
	}
	fix, ok := VapiDeskBy("fix2")
	market, mok := VapiDeskBy("marketing")
	panda, pok := VapiDeskBy("panda")
	if !ok || !mok || !pok || fix.Database == market.Database || market.Database == panda.Database {
		t.Fatalf("fix %v market %v panda %v", fix, market, panda)
	}
	if fix.Database != "Evolu-Jeunes/Fix2" {
		t.Fatalf("fix database = %s", fix.Database)
	}

	k := kernel.New()
	if _, err := VapiFile(k, "marketing", "note", "Campagne", "écrire dans db:panda"); err == nil {
		t.Fatal("a marketing call must not name panda's database")
	}
	got, err := VapiFile(k, "marketing", "note", "Campagne", "brouillon marketing")
	if err != nil || got.Database != "db:marketing" || got.Dialed {
		t.Fatalf("marketing %+v %v", got, err)
	}
	other, err := VapiFile(k, "panda", "note", "Client", "brouillon panda")
	if err != nil || other.Database != "db:panda" {
		t.Fatalf("panda %+v %v", other, err)
	}
	home, err := VapiFile(k, "fix2", "contact", "Cuisine", "chantier fix2")
	if err != nil || home.Database != "Evolu-Jeunes/Fix2" {
		t.Fatalf("fix2 %+v %v", home, err)
	}
	for _, fact := range k.Recall("vapi:db:marketing") {
		if strings.Contains(fact.Text, "Cuisine") || strings.Contains(fact.Text, "panda") {
			t.Fatalf("marketing database saw another project: %+v", fact)
		}
	}
	if len(k.Recall("vapi:Evolu-Jeunes/Fix2")) != 1 || len(k.Recall("vapi:db:panda")) != 1 {
		t.Fatal("each database should keep its own note")
	}
}
