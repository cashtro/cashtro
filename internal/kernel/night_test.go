package kernel

import (
	"context"
	"testing"

	"github.com/cashtro/cashtro/internal/catalog"
)

func TestNightShiftAdvancesCashtroAndStaysGated(t *testing.T) {
	k := New()
	cat := catalog.New()
	k.AttachCatalog(cat)
	k.Register(&stub{spec: Spec{ID: "research", Name: "Research", Autostart: true, Capabilities: []string{"research.ingest"}}})
	k.Register(&stub{spec: Spec{ID: "explorer", Name: "Explorer", Autostart: true, Capabilities: []string{"explorer.search"}}})
	k.Register(&stub{spec: Spec{ID: "memory", Name: "Memory", Autostart: true, Capabilities: []string{"memory.store", "memory.recall"}}})
	k.Register(&stub{spec: Spec{ID: "planner", Name: "Planner", Autostart: true, Capabilities: []string{"planner.backlog"}}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}

	before, err := cat.Get("cashtro-catalog")
	if err != nil || before.Stage != catalog.StageIdea {
		t.Fatalf("seed = %+v %v", before, err)
	}

	b, err := k.NightShift(context.Background(), "keep building asap")
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Steps) < 4 {
		t.Fatalf("steps = %+v", b.Steps)
	}
	for _, s := range b.Steps {
		if s.Agent == "comms" || s.Capability == "comms.send" || s.Capability == "deploy.release" {
			t.Fatalf("night shift must stay gated: %+v", s)
		}
	}

	after, err := cat.Get("cashtro-catalog")
	if err != nil || after.Stage != catalog.StageConcept {
		t.Fatalf("cashtro-catalog should advance to concept: %+v %v", after, err)
	}
	if overnightOpen(cat) < 2 {
		t.Fatal("night shift should keep a second Cashtro ship on the line")
	}
	if len(k.Confirms()) != 0 {
		t.Fatal("night shift must not park outbound comms")
	}
	if len(k.Builds()) != 1 {
		t.Fatalf("builds = %d", len(k.Builds()))
	}
}
