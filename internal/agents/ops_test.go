package agents

import "testing"

func TestOpsChainsHelpTartaria(t *testing.T) {
	ops := Ops()
	if ops.Repo != "Evolu-Jeunes/OPS" || ops.Brain != "ops" || ops.Source != "https://web-ops.shop/" {
		t.Fatalf("repo %+v", ops)
	}
	if len(ops.Chains) != 3 || ops.Agents != 60 || ops.Brains != 300 {
		t.Fatalf("swarm %+v", ops)
	}
	if !ops.LoadsTartaria || !ops.LoadsVoltron || !ops.Automated || ops.Weapons {
		t.Fatal("ops must load Tartaria and Voltron and make no weapon")
	}
	if ops.Sent || ops.Paid || ops.Dialed {
		t.Fatal("ops must not send, charge, or dial")
	}
	if ops.OpenedBy != "voltron" || ops.OrderedBy != "scrum" || ops.DecidedBy != "tartaria" {
		t.Fatalf("seats %+v", ops)
	}
}
