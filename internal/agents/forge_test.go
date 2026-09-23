package agents

import "testing"

func TestForgePartnersWithTartaria(t *testing.T) {
	team := Forge()
	if team.Repo != "Evolu-Jeunes/Forge" || team.Brain != "forge" || !team.LoadsTartaria || !team.LoadsVoltron || team.Weapons {
		t.Fatalf("brain %+v", team)
	}
	if team.Agents != 100 || team.Assistants != 100 || len(team.Departments) != 4 {
		t.Fatalf("swarm %+v", team)
	}
	if team.OpenedBy != "voltron" || team.OrderedBy != "scrum" || team.DecidedBy != "tartaria" {
		t.Fatalf("seats %+v", team)
	}
	if team.Attacks || team.Flood || team.Deployed {
		t.Fatal("the side brain must not attack, flood, or deploy")
	}
	joined := team.Partners[0] + team.Partners[1] + team.Partners[2]
	if joined != "voltronscrumtartaria" {
		t.Fatalf("partners %s", joined)
	}
}
