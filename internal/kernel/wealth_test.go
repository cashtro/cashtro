package kernel

import "testing"

func TestCommitGenerationMonotonic(t *testing.T) {
	k := New()
	k.UpsertOffer(Offer{ID: "ship-pack", Name: "SKU Ship Pack", Lane: LaneCreate, Thesis: "packs"})
	if len(k.Offers()) != 1 {
		t.Fatalf("offers = %d", len(k.Offers()))
	}
	first := k.CommitGeneration(Generation{Score: 0, Offer: "SKU Ship Pack", Lane: LaneCreate, Message: "g1"})
	if first.N != 1 || first.Score < 1 {
		t.Fatalf("first = %+v", first)
	}
	second := k.CommitGeneration(Generation{Score: 0, Offer: "SKU Ship Pack", Lane: LaneCreate, Message: "g2"})
	if second.N != 2 || second.Score <= first.Score {
		t.Fatalf("second = %+v first = %+v", second, first)
	}
	got, ok := k.TouchOffer("ship-pack", second.N, 1)
	if !ok || got.Hits != 1 || got.Score < 1 {
		t.Fatalf("touch = %+v ok=%v", got, ok)
	}
	led := k.Ledger()
	if led.Generation != 2 || led.Score != second.Score || len(led.History) != 2 {
		t.Fatalf("ledger = %+v", led)
	}
}
