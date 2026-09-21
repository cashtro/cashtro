package agents

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBootLoadsAllAgentics(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	procs := k.Processes()
	if len(procs) != 17 {
		t.Fatalf("processes = %d, want 17", len(procs))
	}
	about := k.About()
	if about.Running != 17 || about.Live != 10 {
		t.Fatalf("about = %+v", about)
	}
	if len(k.Notes()) < 5 {
		t.Fatalf("research seeds = %d", len(k.Notes()))
	}
	if k.Catalog() == nil {
		t.Fatal("delivery did not attach catalog")
	}

	res, err := k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.about"})
	if err != nil || !res.OK {
		t.Fatalf("os.about: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "delivery", kernel.Call{Capability: "delivery.list"})
	if err != nil {
		t.Fatal(err)
	}
	ships, ok := res.Data.([]any)
	if !ok {
		// list returns []catalog.Ship
		if res.Data == nil {
			t.Fatal("empty list")
		}
	}
	_ = ships

	payload, _ := json.Marshal(map[string]string{"name": "North desk"})
	res, err = k.Invoke(context.Background(), "delivery", kernel.Call{
		Capability: "delivery.create",
		Payload:    payload,
	})
	if err != nil || !res.OK {
		t.Fatalf("create: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "router", kernel.Call{Capability: "model.status"})
	if err != nil || !res.OK {
		t.Fatalf("router status: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "router", kernel.Call{
		Capability: "model.chat",
		Payload:    []byte(`{"prompt":"ping"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("unbound chat should be ok=false: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"delivery"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("explorer: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "research", kernel.Call{Capability: "research.list"})
	if err != nil || !res.OK {
		t.Fatalf("research: %+v %v", res, err)
	}
}

func TestPulseKeepsAndForgeImproves(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	if len(k.Offers()) != 10 {
		t.Fatalf("offers = %d", len(k.Offers()))
	}
	if err := k.Stop("explorer"); err != nil {
		t.Fatal(err)
	}
	p, _ := k.Process("explorer")
	if p.Status != kernel.StatusStopped {
		t.Fatalf("status = %s", p.Status)
	}

	ctx := context.Background()
	res, err := k.Invoke(ctx, "pulse", kernel.Call{Capability: "pulse.tick"})
	if err != nil || !res.OK {
		t.Fatalf("pulse: %+v %v", res, err)
	}
	p, _ = k.Process("explorer")
	if p.Status != kernel.StatusRunning {
		t.Fatalf("explorer not restored: %+v", p)
	}
	first := k.Ledger()
	if first.Generation != 1 || first.Score < 1 || len(first.History) != 1 {
		t.Fatalf("ledger after pulse = %+v", first)
	}

	var prev int
	for i := 0; i < 11; i++ {
		res, err = k.Invoke(ctx, "forge", kernel.Call{Capability: "forge.tick"})
		if err != nil || !res.OK {
			t.Fatalf("forge %d: %+v %v", i, res, err)
		}
		led := k.Ledger()
		if led.Score <= prev {
			t.Fatalf("score did not rise: %d -> %d", prev, led.Score)
		}
		prev = led.Score
	}
	if k.Generation() != 12 {
		t.Fatalf("generation = %d", k.Generation())
	}

	res, err = k.Invoke(ctx, "wealth", kernel.Call{Capability: "wealth.ledger"})
	if err != nil || !res.OK {
		t.Fatalf("wealth: %+v %v", res, err)
	}
	res, err = k.Invoke(ctx, "router", kernel.Call{Capability: "model.dual", Payload: []byte(`{"prompt":"ping"}`)})
	if err != nil || res.OK {
		t.Fatalf("unbound dual should be ok=false: %+v %v", res, err)
	}
}
