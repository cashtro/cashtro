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

func TestLedgerSealsEveryAnalysis(t *testing.T) {
	for _, a := range AnalyzeAll() {
		if !a.Intact || !a.Connected {
			t.Fatalf("%s ledger = %+v", a.ID, a)
		}
	}
	for _, b := range Brains() {
		moves := AppendMove(nil, "ask:"+b.ID, "Quelle est la position?", b.Do)
		moves = AppendMove(moves, "chain:"+b.ID, "Quelle est la position?", "Epicenter Einstein décide.")
		if !ChainIntact(moves) {
			t.Fatalf("%s moves broke", b.ID)
		}
		moves[0].Coup = "réécrit"
		if ChainIntact(moves) {
			t.Fatalf("%s rewrite stayed intact", b.ID)
		}
	}
}
