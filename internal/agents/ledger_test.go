package agents

import "testing"

func TestMoveChainDoesNotRewrite(t *testing.T) {
	chain := AppendMove(nil, "boot", "Quelle est la position?", "Le tableau est relu.")
	chain = AppendMove(chain, "chain:scanapp", "Quelle est la position?", "Le coup le plus court.")
	if !ChainIntact(chain) || chain[1].Prev != chain[0].Hash {
		t.Fatalf("chain = %+v", chain)
	}
	chain[0].Coup = "réécrit"
	if ChainIntact(chain) {
		t.Fatal("a rewritten move must break the chain")
	}
}
