package agents

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBootLoadsAllAgentics(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	procs := k.Processes()
	if len(procs) != 18 {
		t.Fatalf("processes = %d, want 18", len(procs))
	}
	about := k.About()
	if about.Running != 18 || about.Live != 11 {
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

	if k.Symbols() == nil {
		t.Fatal("symbols did not attach bus")
	}
	desk := k.DeskCard()
	if desk.PendingN < 10 || desk.Scan.Unread != 163 {
		t.Fatalf("chooser desk = %+v", desk)
	}
	if len(k.Requests()) < 2 {
		t.Fatalf("desk seeds = %d", len(k.Requests()))
	}
	res, err = k.Invoke(context.Background(), "desk", kernel.Call{Capability: "desk.plan"})
	if err != nil || !res.OK {
		t.Fatalf("desk.plan: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "watch", kernel.Call{Capability: "watch.status"})
	if err != nil || !res.OK {
		t.Fatalf("watch.status: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "symbols", kernel.Call{Capability: "symbols.list"})
	if err != nil || !res.OK {
		t.Fatalf("symbols.list: %+v %v", res, err)
	}
	payload, _ = json.Marshal(map[string]string{"id": "desk-pulse"})
	res, err = k.Invoke(context.Background(), "symbols", kernel.Call{Capability: "symbols.fire", Payload: payload})
	if err != nil || !res.OK {
		t.Fatalf("desk-pulse: %+v %v", res, err)
	}
}

func TestSymbolsIntakeCreatesShip(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]string{
		"id":     "intake-to-idea",
		"name":   "Hook desk",
		"client": "Cashtro",
		"sector": "internal",
		"notes":  "caught without Zapier",
	})
	res, err := k.Invoke(context.Background(), "symbols", kernel.Call{Capability: "symbols.fire", Payload: payload})
	if err != nil || !res.OK {
		t.Fatalf("intake: %+v %v", res, err)
	}
	ship, err := k.Catalog().Get("hook-desk")
	if err != nil {
		t.Fatal(err)
	}
	if ship.Name != "Hook desk" || ship.Stage != "idea" {
		t.Fatalf("ship = %+v", ship)
	}
	facts := k.Recall("hook desk")
	if len(facts) == 0 {
		t.Fatal("intake should remember the name")
	}
}

func TestSymbolsMailEventIngestsNote(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	before := len(k.Notes())
	if _, err := k.Post("explorer", "research", "job", "walk the symbols line"); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		for _, n := range k.Notes() {
			if n.Source == "symbols" && n.Claim == "to research · job" {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("mail event did not ingest a note; before=%d after=%d %#v", before, len(k.Notes()), k.Notes()[:min(3, len(k.Notes()))])
}
