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
	if len(procs) != 14 {
		t.Fatalf("processes = %d, want 14", len(procs))
	}
	about := k.About()
	if about.Running != 14 || about.Live != 13 {
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

	res, err = k.Invoke(context.Background(), "investigator", kernel.Call{
		Capability: "investigator.trace",
		Payload:    []byte(`{"query":"AOS"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("investigator: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "architect", kernel.Call{
		Capability: "architect.plan",
		Payload:    []byte(`{"goal":"Ship overnight keep-alive desk"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("architect: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "security", kernel.Call{
		Capability: "security.triage",
		Payload:    []byte(`{"query":"AOS"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("security: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "deploy", kernel.Call{
		Capability: "deploy.release",
		Payload:    []byte(`{"target":"cashtro-os"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("deploy: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "reviewer", kernel.Call{
		Capability: "reviewer.watch",
		Payload:    []byte(`{"target":"cashtro-os"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("reviewer: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "operator", kernel.Call{
		Capability: "operator.browse",
		Payload:    []byte(`{"url":"http://127.0.0.1:8080/"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("operator: %+v %v", res, err)
	}
}
