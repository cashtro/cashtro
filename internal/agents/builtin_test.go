package agents

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cashtro/cashtro/internal/graphify"
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

	if len(AgentIDs()) != 15 {
		t.Fatalf("agentics in ecosystems = %d, want 15 (%v)", len(AgentIDs()), AgentIDs())
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.ecosystems"})
	if err != nil || !res.OK {
		t.Fatalf("manager.ecosystems: %+v %v", res, err)
	}
	if links, ok := res.Data.([]Link); !ok || len(links) < 8 {
		t.Fatalf("links = %#v", res.Data)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.automate"})
	if err != nil || !res.OK {
		t.Fatalf("manager.automate: %+v %v", res, err)
	}
	gotAgents := res.Data.(map[string]any)["agents"].([]string)
	if len(gotAgents) != 15 {
		t.Fatalf("automate agents = %v", gotAgents)
	}
	if len(k.Inbox("operator")) == 0 || len(k.Inbox("security")) == 0 || len(k.Inbox("init")) == 0 {
		t.Fatal("ecosystem mail did not reach the agentics")
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.assign",
		Payload:    []byte(`{"brain":"Forgeron","repo":"Evolu-Jeunes/Btkavocat"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.assign: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.brains"})
	if err != nil || !res.OK {
		t.Fatalf("manager.brains: %+v %v", res, err)
	}
	if brains, ok := res.Data.([]Brain); !ok || len(brains) != 5 {
		t.Fatalf("brains = %#v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.graphify",
		Payload:    []byte(`{"query":"giant"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.graphify: %+v %v", res, err)
	}

	before := len(k.Events())
	snap := Snapshot(k, "giant")
	if len(k.Events()) != before {
		t.Fatal("snapshot wrote journal")
	}
	g := snap["graph"].(graphify.Graph)
	if g.Counts["token"] < 1 || g.Score.GiantOnline == false {
		t.Fatalf("snapshot graph = %+v", g.Counts)
	}
}
