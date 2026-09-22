package castro

import "testing"

func TestMapHasBothHomesAndEveryAgentic(t *testing.T) {
	cashtroHome, evolu := Map()
	if cashtroHome.Name != "cashtro/cashtro" {
		t.Fatalf("cashtro = %q", cashtroHome.Name)
	}
	if evolu.Name != "Evolu-Jeunes" {
		t.Fatalf("evolu = %q", evolu.Name)
	}
	wantAgents := []string{
		"init", "delivery", "router", "research", "explorer", "operator",
		"reviewer", "architect", "deploy", "security", "memory", "comms",
		"planner", "investigator",
	}
	got := map[string]bool{}
	for _, n := range cashtroHome.Nodes {
		got[n] = true
	}
	for _, id := range wantAgents {
		if !got[id] {
			t.Fatalf("missing agentic %s in %#v", id, cashtroHome.Nodes)
		}
	}
	need := []string{"evolujeunes.ca", "btk-avocats", "proximity", "scanapp", "private-repos-access-limited"}
	egot := map[string]bool{}
	for _, n := range evolu.Nodes {
		egot[n] = true
	}
	for _, id := range need {
		if !egot[id] {
			t.Fatalf("missing evolu node %s", id)
		}
	}
}
