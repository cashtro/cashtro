package agents

import (
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestVapiIsPandaCustomerService(t *testing.T) {
	desk := Vapi()
	if desk.Project != "panda" || desk.Database != "db:panda" || desk.Role != "service client" || desk.Voice != "vapi" {
		t.Fatalf("desk %+v", desk)
	}
	if !strings.Contains(desk.Prompt, "service client") || !strings.Contains(desk.Prompt, "db:panda") {
		t.Fatalf("prompt = %s", desk.Prompt)
	}

	k := kernel.New()
	if _, err := VapiFile(k, "marketing", "note", "Campagne", "brouillon"); err == nil {
		t.Fatal("marketing must not have Vapi")
	}
	if _, err := VapiFile(k, "fix2", "contact", "Cuisine", "chantier"); err == nil {
		t.Fatal("Fix Tout must not have Vapi")
	}
	if _, err := VapiFile(k, "panda", "note", "Client", "écrire dans db:marketing"); err == nil {
		t.Fatal("Panda customer service must stay on db:panda")
	}
	got, err := VapiFile(k, "panda", "note", "Client", "question sur le cours")
	if err != nil || got.Database != "db:panda" || got.Dialed || got.Project != "panda" {
		t.Fatalf("panda %+v %v", got, err)
	}
	if len(k.Recall("vapi:db:panda")) != 1 {
		t.Fatal("the note should sit on Panda's database")
	}
}
