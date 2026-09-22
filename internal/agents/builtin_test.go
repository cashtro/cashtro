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
	if len(procs) != 15 {
		t.Fatalf("processes = %d, want 15", len(procs))
	}
	about := k.About()
	if about.Running != 15 || about.Live != 8 {
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

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.status"})
	if err != nil || !res.OK {
		t.Fatalf("manager.status: %+v %v", res, err)
	}
	if got := res.Data.(map[string]any)["repos"]; got != 57 {
		t.Fatalf("manager.repos = %v, want 57", got)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.lines"})
	if err != nil || !res.OK {
		t.Fatalf("manager.lines: %+v %v", res, err)
	}
	lines, ok := res.Data.([]Line)
	if !ok || len(lines) != 8 {
		t.Fatalf("lines = %#v", res.Data)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.line",
		Payload:    []byte(`{"id":"proximity"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.line: %+v %v", res, err)
	}
	if res.Data.(Line).Brain != "Forgeron" {
		t.Fatalf("proximity brain = %+v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.assign",
		Payload:    []byte(`{"brain":"Forgeron","repo":"Evolu-Jeunes/Btkavocat"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.assign: %+v %v", res, err)
	}
}
