package agents

import "testing"

func TestOpsChainsHelpInstinct(t *testing.T) {
	ops := Ops()
	if ops.Repo != "Evolu-Jeunes/OPS" || ops.Brain != "ops" || ops.Website || !ops.InVoltron || !ops.Operative || ops.ReportsTo != "voltron" {
		t.Fatalf("repo %+v", ops)
	}
	if len(ops.Chains) != 3 || ops.Agents != 60 || ops.Brains != 300 {
		t.Fatalf("swarm %+v", ops)
	}
	if !ops.LoadsInstinct || !ops.LoadsVoltron || !ops.Automated || ops.Weapons {
		t.Fatal("ops must load Epicenter Einstein and Voltron and make no weapon")
	}
	if ops.Sent || ops.Paid || ops.Dialed {
		t.Fatal("ops must not send, charge, or dial")
	}
	if ops.OpenedBy != "voltron" || ops.OrderedBy != "scrum" || ops.DecidedBy != "instinct" {
		t.Fatalf("seats %+v", ops)
	}
}
