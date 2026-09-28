package agents

import "testing"

func TestVoltronHoldsEveryChain(t *testing.T) {
	ok, loose := VoltronHolds()
	if !ok {
		t.Fatalf("loose chain %s", loose)
	}
	layers := Layers()
	if len(layers) != 6 {
		t.Fatalf("layers = %d", len(layers))
	}
	seen := map[string]bool{}
	for _, layer := range layers {
		seen[layer.ID] = true
		if layer.Website {
			t.Fatalf("%s is a website", layer.ID)
		}
	}
	for _, id := range []string{"ops", "hustle", "eye", "forge", "giant", "fusion"} {
		if !seen[id] {
			t.Fatalf("missing %s", id)
		}
	}
	if len(Lines()) != 14 || len(Links()) != 23 {
		t.Fatalf("lines %d links %d", len(Lines()), len(Links()))
	}
	os := TheOS()
	if os.Role != "OS of the agentics" || os.Brain != "instinct" || !os.Acts || os.Replies || len(os.Stages) != 3 || os.Stages[0] != "idea" || os.Stages[2] != "production" {
		t.Fatalf("os %+v", os)
	}
	got := map[string]bool{}
	for _, start := range os.Starts {
		if start.StartedBy != "voltron" || !start.Acts || start.Replies || start.Task == "" {
			t.Fatalf("start %+v", start)
		}
		got[start.ID] = true
	}
	for _, id := range []string{"ops", "hustle", "eye", "forge", "giant", "fusion", "nft-giant", "trading", "giant-trading", "maitre"} {
		if !got[id] {
			t.Fatalf("Voltron does not start %s", id)
		}
	}
	for _, id := range []string{"wordpress", "ecole", "fix2", "educonnexion"} {
		if got[id] {
			t.Fatalf("Voltron started a website %s", id)
		}
	}
}
