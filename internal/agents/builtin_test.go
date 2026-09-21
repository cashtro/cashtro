package agents

import (
	"context"
	"encoding/json"
	"strings"
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

	res, err = k.Invoke(context.Background(), "teams", kernel.Call{Capability: "teams.status"})
	if err != nil || !res.OK {
		t.Fatalf("teams status: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "teams", kernel.Call{
		Capability: "teams.say",
		Payload:    []byte(`{"text":"open slack for me","channel":"teams.microsoft","tenant":"lroximity","kind":"conversation"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("teams say: %+v %v", res, err)
	}
	th, ok := res.Data.(kernel.TeamThread)
	if !ok || len(th.Messages) < 3 {
		t.Fatalf("thread = %#v", res.Data)
	}
	last := th.Messages[len(th.Messages)-1]
	if last.Role != "assistant" || !strings.Contains(strings.ToLower(last.Body), "microsoft teams") {
		t.Fatalf("reply = %+v", last)
	}
	res, err = k.Invoke(context.Background(), "teams", kernel.Call{
		Capability: "teams.thread",
		Payload:    []byte(`{"channel":"slack","tenant":"Proximity","kind":"conversation"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("slack should fail: %+v %v", res, err)
	}
}
