package agents

import (
	"context"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestFusionHasEightyWorkersAcrossBothEcosystems(t *testing.T) {
	bureau := Fusion()
	if bureau.ID != "fusion" || bureau.Name != "The Fusion" {
		t.Fatalf("bureau identity = %+v", bureau)
	}
	if bureau.WorkerCount != 80 || len(bureau.Workers) != 80 {
		t.Fatalf("workers = %d/%d, want 80", bureau.WorkerCount, len(bureau.Workers))
	}
	if len(bureau.Scope) != 2 || bureau.Scope[0] != "cashtro" || bureau.Scope[1] != "Evolu-Jeunes" {
		t.Fatalf("scope = %v", bureau.Scope)
	}
	perLaw := map[string]int{}
	ids := map[string]bool{}
	for _, worker := range bureau.Workers {
		if ids[worker.ID] {
			t.Fatalf("duplicate worker id %s", worker.ID)
		}
		ids[worker.ID] = true
		perLaw[worker.LawID]++
		if !worker.ReadOnly || len(worker.Scope) != 2 || worker.ReportsTo != "fusion" {
			t.Fatalf("unsafe worker = %+v", worker)
		}
	}
	for _, law := range Loi().Laws {
		if perLaw[law.ID] != 10 {
			t.Fatalf("%s workers = %d, want 10", law.ID, perLaw[law.ID])
		}
	}
}

func TestFusionChecksAndRemediationStayHumanGated(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "fusion", kernel.Call{
		Capability: "fusion.status",
	})
	if err != nil || !res.OK {
		t.Fatalf("fusion.status: %+v %v", res, err)
	}
	if res.Data.(FusionBureau).WorkerCount != 80 {
		t.Fatalf("status = %+v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "fusion", kernel.Call{
		Capability: "fusion.check",
		Payload:    []byte(`{"lawId":"loi-25"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("empty evidence must hold: %+v %v", res, err)
	}
	check := res.Data.(FusionCheck)
	if check.Decision != "hold" || !check.HumanGate || check.WorkerID == "" {
		t.Fatalf("check = %+v", check)
	}

	res, err = k.Invoke(context.Background(), "fusion", kernel.Call{
		Capability: "fusion.remediate",
		Payload:    []byte(`{"finding":"secret exposé","action":"retirer du dépôt"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("remediation must wait: %+v %v", res, err)
	}
	confirm := res.Data.(kernel.Confirm)
	if confirm.Status != "pending" || confirm.Agent != "fusion" {
		t.Fatalf("confirm = %+v", confirm)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.fusion"})
	if err != nil || !res.OK || res.Data.(FusionBureau).WorkerCount != 80 {
		t.Fatalf("manager.fusion: %+v %v", res, err)
	}
}

func TestGraphifyStandingExtractIsCurrent(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.graphify",
		Payload:    []byte(`{"query":"fusion scope"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.graphify: %+v %v", res, err)
	}
	graph := res.Data.(map[string]any)
	if graph["owner"] != "cashtro" || graph["repositories"] != 63 || graph["nodes"] != 52527 || graph["edges"] != 133783 {
		t.Fatalf("graph extract = %+v", graph)
	}
}
