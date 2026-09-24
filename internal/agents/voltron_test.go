package agents

import "testing"

func TestVoltronHoldsEveryChain(t *testing.T) {
	ok, loose := VoltronHolds()
	if !ok {
		t.Fatalf("loose chain %s", loose)
	}
	layers := Layers()
	if len(layers) != 5 {
		t.Fatalf("layers = %d", len(layers))
	}
	seen := map[string]bool{}
	for _, layer := range layers {
		seen[layer.ID] = true
		if layer.Website {
			t.Fatalf("%s is a website", layer.ID)
		}
	}
	for _, id := range []string{"ops", "hustle", "eye", "forge", "giant"} {
		if !seen[id] {
			t.Fatalf("missing %s", id)
		}
	}
	if len(Lines()) != 14 || len(Links()) != 30 {
		t.Fatalf("lines %d links %d", len(Lines()), len(Links()))
	}
}
