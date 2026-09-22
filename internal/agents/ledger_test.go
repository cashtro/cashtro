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

func TestAgencyWatchesEveryLine(t *testing.T) {
	seen := map[string]bool{}
	for _, ln := range WatchLinks() {
		if ln.From != "agence" {
			t.Fatalf("watch from %s", ln.From)
		}
		seen[ln.To] = true
	}
	for _, ln := range Lines() {
		if ln.ID == "agence" {
			continue
		}
		if !seen[ln.ID] {
			t.Fatalf("line %s is off the board", ln.ID)
		}
	}
	if len(Links()) != len(baseLinks())+len(WatchLinks()) {
		t.Fatalf("links = %d", len(Links()))
	}
}
