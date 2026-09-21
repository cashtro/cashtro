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

	res, err = k.Invoke(context.Background(), "n8n", kernel.Call{
		Capability: "n8n.workflow",
		Payload:    []byte(`{"prompt":"Advance ScanApp. Keep confirms."}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("n8n.workflow: %+v %v", res, err)
	}
}

func TestN8NIsOwned(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "n8n", kernel.Call{Capability: "n8n.list"})
	if err != nil || !res.OK {
		t.Fatalf("list: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "n8n", kernel.Call{
		Capability: "n8n.workflow",
		Payload:    []byte(`{"id":"wealth-dual","prompt":"Park a retainer. Keep the human gate."}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("workflow: %+v %v", res, err)
	}
	raw, _ := json.Marshal(res.Data)
	if !json.Valid(raw) {
		t.Fatal("data not json")
	}
	var run map[string]any
	if err := json.Unmarshal(raw, &run); err != nil {
		t.Fatal(err)
	}
	if run["owner"] != "cashtro" {
		t.Fatalf("owner = %v", run["owner"])
	}
	if run["bound"] != false {
		t.Fatalf("bound = %v, owned loop must not call vendor keys", run["bound"])
	}
	res, err = k.Invoke(context.Background(), "n8n", kernel.Call{
		Capability: "n8n.workflow",
		Payload:    []byte(`{"id":"intake-hook","prompt":"new QC clinic lead"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("hook workflow: %+v %v", res, err)
	}
	if len(k.Confirms()) == 0 {
		t.Fatal("intake-hook should park a human confirm")
	}
}
