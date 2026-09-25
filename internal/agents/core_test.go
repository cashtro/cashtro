package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestCoreCreatesKernelChainWithoutGrowingParent(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	if k.About().Agents != 15 {
		t.Fatalf("boot agents = %d", k.About().Agents)
	}

	raw, _ := json.Marshal(map[string]any{
		"id": "alpha",
		"agents": []map[string]string{
			{"id": "scout", "name": "Scout"},
		},
	})
	res, err := k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.core", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("os.core alpha: %+v %v", res, err)
	}
	data := res.Data.(map[string]any)
	link := data["link"].(kernel.KernelLink)
	if link.Parent != "root" || link.Index != 0 || !strings.Contains(link.Infrastructure, "Trinity") || len(link.Agents) != 1 || link.Agents[0] != "scout" {
		t.Fatalf("alpha link = %+v", link)
	}
	if data["parentAgents"].(int) != 15 || k.About().Agents != 15 {
		t.Fatalf("parent grew: %+v about=%d", data["parentAgents"], k.About().Agents)
	}

	raw, _ = json.Marshal(map[string]any{
		"id":             "beta",
		"infrastructure": "desk beta",
		"agents": []map[string]string{
			{"id": "eye", "name": "Eye"},
			{"id": "hand", "name": "Hand"},
		},
	})
	res, err = k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.core", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("os.core beta: %+v %v", res, err)
	}
	chain := k.Core().Chain()
	if len(chain) != 2 || chain[1].Parent != "alpha" || chain[1].Prev != chain[0].Hash || chain[1].Infrastructure != "desk beta" {
		t.Fatalf("chain = %+v", chain)
	}
	if chain[0].Hash == "" || chain[1].Hash == "" || chain[1].Hash == chain[0].Hash {
		t.Fatal("chain hash broke")
	}
	beta, ok := k.Core().Kernel("beta")
	if !ok || beta.About().Agents != 2 {
		t.Fatal("beta kernel missing its own agents")
	}
	run, err := beta.Invoke(context.Background(), "hand", kernel.Call{Capability: "hand.run"})
	if err != nil || !run.OK {
		t.Fatalf("hand run: %+v %v", run, err)
	}
	if _, err := k.Process("scout"); err == nil {
		t.Fatal("scout registered on the parent")
	}
	if k.About().Agents != 15 {
		t.Fatalf("parent agents = %d", k.About().Agents)
	}

	res, err = k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.core", Payload: []byte(`{"id":"alpha"}`)})
	if err != nil || res.OK {
		t.Fatalf("duplicate kernel: %+v %v", res, err)
	}
}
