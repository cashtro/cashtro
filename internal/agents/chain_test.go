package agents

import (
	"strings"
	"testing"
)

func TestScanChainAndFicheGate(t *testing.T) {
	steps, ok := Chain("scanapp")
	if !ok || len(steps) != 9 {
		t.Fatalf("steps = %d ok %v", len(steps), ok)
	}
	if steps[0].Agent != "operator" || steps[8].Agent != "operator" {
		t.Fatalf("ends = %+v %+v", steps[0], steps[8])
	}

	bad := Fiche{Code: "123", Name: "Noix", Price: "4.00", Stock: 3, Photo: "database"}
	if ready, why := bad.Ready(); ready || why != "l'image de la base n'est pas la fiche" {
		t.Fatalf("database photo: ready=%v why=%s", ready, why)
	}
	empty := Fiche{Code: "123", Name: "Noix", Price: "4.00", Stock: 2}
	if ready, _ := empty.Ready(); ready {
		t.Fatal("missing photo should not publish")
	}
	person := Fiche{Code: "123", Name: "Noix", Price: "4.00", Stock: 2, Photo: "generated", Person: true}
	if ready, why := person.Ready(); ready || why != "personne identifiable sans consentement" {
		t.Fatalf("person: ready=%v why=%s", ready, why)
	}
	good := Fiche{Code: "123", Name: "Noix", Price: "4.00", Stock: 2, Photo: "licensed"}
	if ready, why := good.Ready(); !ready {
		t.Fatalf("good fiche blocked: %s", why)
	}
	if _, ok := Chain("nope"); ok {
		t.Fatal("unknown chain")
	}
}

func TestEveryLineRunsAClosedChain(t *testing.T) {
	for _, ln := range Lines() {
		steps, ok := Chain(ln.ID)
		if !ok || len(steps) < 2 {
			t.Fatalf("%s steps = %d ok %v", ln.ID, len(steps), ok)
		}
		if open := AskSelf(ln.ID, nil).Open; len(open) != 0 {
			t.Fatalf("%s still open: %+v", ln.ID, open)
		}
	}
	if open := AskSelf("watch", nil).Open; len(open) != 0 {
		t.Fatalf("watch open: %+v", open)
	}
	steps, ok := Chain("watch")
	if !ok || len(steps) < 2 {
		t.Fatalf("watch chain = %d ok %v", len(steps), ok)
	}
}

func TestChainReachesEveryBrain(t *testing.T) {
	if len(Lines()) != 14 {
		t.Fatalf("lines = %d", len(Lines()))
	}
	for _, b := range Brains() {
		steps, ok := Chain(b.ID)
		if !ok || len(steps) < 2 || steps[0].Agent != "memory" || steps[1].Agent != "manager" {
			t.Fatalf("%s chain = %+v ok %v", b.ID, steps, ok)
		}
		if steps[0].Do == "" || !strings.Contains(steps[1].Do, "Epicenter Einstein") {
			t.Fatalf("%s steps do not close on Epicenter Einstein: %+v", b.ID, steps)
		}
	}
	if gap := Disconnected(); len(gap) != 0 {
		t.Fatalf("chain not connected: %v", gap)
	}
}
