package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestInquisitorStepsUpCorporationAndBlocksUntilSmarter(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}

	res, err := k.Invoke(context.Background(), "inquisitor", kernel.Call{Capability: "inquisitor.status"})
	if err != nil || !res.OK {
		t.Fatalf("status: %+v %v", res, err)
	}
	st := res.Data.(map[string]any)
	depts, ok := st["departments"].([]Department)
	if !ok || len(depts) != 16 {
		t.Fatalf("departments = %#v", st["departments"])
	}
	if st["blocked"].(int) != 0 {
		t.Fatalf("boot must not freeze departments: %+v", st)
	}
	byID := map[string]Department{}
	for _, d := range depts {
		byID[d.ID] = d
		if d.Score < 1 {
			t.Fatalf("%s was not stepped up at boot: %+v", d.ID, d)
		}
	}
	for _, id := range []string{"corporation", "hustler", "proximity", "teal", "trading"} {
		if _, ok := byID[id]; !ok {
			t.Fatalf("missing department %s", id)
		}
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.audit",
		Payload:    []byte(`{"department":"hustler"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("audit: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.dissent",
		Payload:    []byte(`{"action":"let Hustler run every line"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("dissent: %+v %v", res, err)
	}
	dissent := res.Data.(map[string]any)
	if !strings.Contains(dissent["reject"].(string), "No.") {
		t.Fatalf("dissent reject = %#v", dissent)
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.block",
		Payload:    []byte(`{"department":"panda"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("block: %+v %v", res, err)
	}
	if !DepartmentBlocked(k, "panda") {
		t.Fatal("panda should be blocked")
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.call",
		Payload:    []byte(`{"to":"+15555550123","line":"panda"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("blocked call should refuse: %+v %v", res, err)
	}
	if !strings.Contains(res.Message, "blocked") {
		t.Fatalf("call msg = %q", res.Message)
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.release",
		Payload:    []byte(`{"department":"panda"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("release before smarter should fail: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.smarter",
		Payload:    []byte(`{"department":"panda"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("smarter: %+v %v", res, err)
	}
	if DepartmentBlocked(k, "panda") {
		t.Fatal("panda should self-release after smarter")
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.call",
		Payload:    []byte(`{"to":"+15555550123","line":"panda","prompt":"after smarter"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("call after smarter: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.optimize",
		Payload:    []byte(`{"department":"corporation","option":"one shared brain for everything"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("optimize: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.ask",
		Payload:    []byte(`{"action":"rewrite proximity","department":"proximity"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("ask: %+v %v", res, err)
	}
	ask := res.Data.(kernel.Ask)
	if ask.Status != "pending" || len(ask.Questions) < 2 {
		t.Fatalf("ask = %+v", ask)
	}
	payload, _ := json.Marshal(map[string]any{"id": ask.ID, "answers": map[string]string{"note": "specialize per-site QA"}})
	res, err = k.Invoke(context.Background(), "inquisitor", kernel.Call{
		Capability: "inquisitor.answer",
		Payload:    payload,
	})
	if err != nil || !res.OK {
		t.Fatalf("answer: %+v %v", res, err)
	}
	answered := res.Data.(kernel.Ask)
	if answered.Status != "answered" {
		t.Fatalf("answered = %+v", answered)
	}

	inq, err := k.Process("inquisitor")
	if err != nil {
		t.Fatal(err)
	}
	if len(inq.Spec.Rules) < 4 {
		t.Fatalf("inquisitor rules = %#v", inq.Spec.Rules)
	}
}
