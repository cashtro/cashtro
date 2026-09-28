package agents

import (
	"context"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestRepairCrewReadsTheWholeChain(t *testing.T) {
	crew := RepairCrew()
	if len(crew) != 20 || crew[0].ID != "maitre" || crew[19].ID != "einstein" {
		t.Fatalf("crew = %d lead %+v", len(crew), crew[0])
	}
	for i, agent := range crew {
		if agent.Order != i+1 || agent.Craft == "" || agent.Name == "" {
			t.Fatalf("seat %+v", agent)
		}
	}
	report := SurveyChain()
	if report.Lead != "maitre" || report.Decided != "instinct" || report.Applied || len(report.Repair) < 8 {
		t.Fatalf("report %+v", report)
	}
	found := map[string]bool{}
	for _, br := range report.Breaks {
		if br.Fact == "" || br.Where == "" {
			t.Fatalf("empty break %+v", br)
		}
		found[br.ID] = true
	}
	if !found["site-dans-la-chaine"] || !found["repo-deux-lignes"] || !found["nom-double"] {
		t.Fatalf("breaks = %+v", report.Breaks)
	}
	if found["vapi-hors-panda"] || found["agent-hors-table"] || found["ordre"] || found["voltron"] {
		t.Fatalf("false break %+v", report.Breaks)
	}
	one := Chain1()
	if one.Name != "Chain 1" || one.ID != "chain-1" || one.Agentic || one.Repo != "Evolu-Jeunes/educonnexion" || one.Owner == "" || one.With != "Proximity agency" {
		t.Fatalf("chain 1 %+v", one)
	}
	for _, ln := range Lines() {
		for _, repo := range ln.Repos {
			if repo == one.Repo {
				t.Fatalf("%s still holds Chain 1", ln.ID)
			}
		}
	}

	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	if k.About().Agents != 15 {
		t.Fatalf("boot agents = %d", k.About().Agents)
	}
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.repair"})
	if err != nil || !res.OK {
		t.Fatalf("repair: %+v %v", res, err)
	}
	got := res.Data.(ChainReport)
	if len(got.Agents) != 20 || got.Lead != "maitre" {
		t.Fatalf("invoke %+v", got.Lead)
	}
}
