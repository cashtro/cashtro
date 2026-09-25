package kernel

import (
	"context"
	"strings"
	"testing"
)

func TestCoreChainsKernelsWithTheirOwnAgents(t *testing.T) {
	ctx := context.Background()
	root := New()
	if root.About().Agents != 0 {
		t.Fatalf("root agents = %d", root.About().Agents)
	}
	core := root.Core()

	alpha, err := core.SpawnKernel(ctx, "alpha", "Trinity desk", nil)
	if err != nil {
		t.Fatal(err)
	}
	if alpha.Index != 0 || alpha.Parent != "root" || alpha.Prev != "" || alpha.Hash == "" || alpha.Hash != alpha.seal() {
		t.Fatalf("alpha = %+v", alpha)
	}
	if len(alpha.Agents) != 1 || alpha.Agents[0] != "alpha-agent" {
		t.Fatalf("alpha agents = %+v", alpha.Agents)
	}
	if root.About().Agents != 0 {
		t.Fatal("child agent landed on the parent")
	}
	kid, ok := core.Kernel("alpha")
	if !ok {
		t.Fatal("alpha kernel missing")
	}
	proc, err := kid.Process("alpha-agent")
	if err != nil || proc.Status != StatusRunning || kid.About().Agents != 1 {
		t.Fatalf("alpha process = %+v %v", proc, err)
	}
	res, err := kid.Invoke(ctx, "alpha-agent", Call{Capability: "alpha.run"})
	if err != nil || !res.OK {
		t.Fatalf("alpha run: %+v %v", res, err)
	}

	beta, err := core.SpawnKernel(ctx, "beta", "Trinity desk", []Spec{
		{ID: "beta-eye", Name: "Eye"},
		{ID: "beta-hand", Name: "Hand"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if beta.Index != 1 || beta.Parent != "alpha" || beta.Prev != alpha.Hash || beta.Hash != beta.seal() {
		t.Fatalf("beta = %+v", beta)
	}
	if len(beta.Agents) != 2 || beta.Agents[0] != "beta-eye" || beta.Agents[1] != "beta-hand" {
		t.Fatalf("beta agents = %+v", beta.Agents)
	}
	hand, ok := core.Kernel("beta")
	if !ok || hand.About().Agents != 2 {
		t.Fatal("beta kernel missing its agents")
	}
	if _, err := hand.Process("alpha-agent"); err == nil {
		t.Fatal("alpha agent is on beta")
	}

	chain := core.Chain()
	if len(chain) != 2 || chain[1].Prev != chain[0].Hash || chain[0].Hash != chain[0].seal() || chain[1].Hash != chain[1].seal() {
		t.Fatalf("chain = %+v", chain)
	}

	if _, err := core.SpawnKernel(ctx, "alpha", "again", nil); err == nil {
		t.Fatal("duplicate kernel accepted")
	}
	if _, err := core.SpawnKernel(ctx, "  ", "desk", nil); err == nil {
		t.Fatal("empty id accepted")
	}
	if _, err := core.SpawnKernel(ctx, "gamma", "  ", nil); err == nil {
		t.Fatal("empty infrastructure accepted")
	}
	if _, err := core.SpawnKernel(ctx, "gamma", "desk", []Spec{{ID: "same", Name: "A"}, {ID: "same", Name: "B"}}); err == nil {
		t.Fatal("duplicate agent id accepted")
	}
	if _, err := kid.CreateAgent(ctx, Spec{ID: "alpha-agent", Name: "Again"}); err == nil || !strings.Contains(err.Error(), "exists") {
		t.Fatalf("duplicate create = %v", err)
	}

	nested, err := kid.Core().SpawnKernel(ctx, "alpha-cell", "cell desk", []Spec{{ID: "cell", Name: "Cell"}})
	if err != nil || nested.Parent != "root" || nested.Agents[0] != "cell" {
		t.Fatalf("nested = %+v %v", nested, err)
	}
	if len(core.Chain()) != 2 || len(kid.Core().Chain()) != 1 {
		t.Fatal("nested chain wrote on the parent core")
	}
	if root.About().Agents != 0 || kid.About().Agents != 1 {
		t.Fatal("nested spawn changed the parent tables")
	}
}
