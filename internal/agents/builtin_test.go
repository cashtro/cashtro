package agents

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBootLoadsAllAgentics(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	procs := k.Processes()
	if len(procs) != 15 {
		t.Fatalf("processes = %d, want 15", len(procs))
	}
	about := k.About()
	if about.Running != 15 || about.Live != 8 {
		t.Fatalf("about = %+v", about)
	}
	if len(k.Notes()) < 5 {
		t.Fatalf("research seeds = %d", len(k.Notes()))
	}
	if len(k.Requests()) < 2 {
		t.Fatalf("desk seeds = %d", len(k.Requests()))
	}
	if k.Catalog() == nil {
		t.Fatal("delivery did not attach catalog")
	}

	res, err := k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.about"})
	if err != nil || !res.OK {
		t.Fatalf("os.about: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "delivery", kernel.Call{Capability: "delivery.list"})
	if err != nil {
		t.Fatal(err)
	}
	ships, ok := res.Data.([]any)
	if !ok {
		// list returns []catalog.Ship
		if res.Data == nil {
			t.Fatal("empty list")
		}
	}
	_ = ships

	payload, _ := json.Marshal(map[string]string{"name": "North desk"})
	res, err = k.Invoke(context.Background(), "delivery", kernel.Call{
		Capability: "delivery.create",
		Payload:    payload,
	})
	if err != nil || !res.OK {
		t.Fatalf("create: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "router", kernel.Call{Capability: "model.status"})
	if err != nil || !res.OK {
		t.Fatalf("router status: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "router", kernel.Call{
		Capability: "model.chat",
		Payload:    []byte(`{"prompt":"ping"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("unbound chat should be ok=false: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"delivery"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("explorer: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "research", kernel.Call{Capability: "research.list"})
	if err != nil || !res.OK {
		t.Fatalf("research: %+v %v", res, err)
	}
}

func TestDeskCapturesAndBetters(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}

	res, err := k.Invoke(context.Background(), "desk", kernel.Call{Capability: "desk.plan"})
	if err != nil || !res.OK {
		t.Fatalf("plan: %+v %v", res, err)
	}
	plan, ok := res.Data.(Plan)
	if !ok || len(plan.Can) == 0 || len(plan.Cannot) == 0 || len(plan.Create) == 0 {
		t.Fatalf("plan data = %#v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "desk", kernel.Call{
		Capability: "desk.capture",
		Payload:    []byte(`{"raw":"ship a research note about the desk"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("capture: %+v %v", res, err)
	}
	first, ok := res.Data.(kernel.Request)
	if !ok || first.ID == 0 || first.Raw == "" || first.Improved == "" || first.Pass != 1 || !first.Allowed {
		t.Fatalf("captured = %+v", first)
	}
	if first.Owner != "research" {
		t.Fatalf("owner = %s", first.Owner)
	}

	res, err = k.Invoke(context.Background(), "desk", kernel.Call{
		Capability: "desk.better",
		Payload:    []byte(`{"id":` + jsonInt(first.ID) + `}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("better: %+v %v", res, err)
	}
	second, ok := res.Data.(kernel.Request)
	if !ok || second.Pass != 2 || second.Raw != first.Raw || second.Status != kernel.StatusClarified {
		t.Fatalf("bettered = %+v", second)
	}
	if !strings.Contains(second.Improved, "Acceptance:") {
		t.Fatalf("improved missing acceptance: %s", second.Improved)
	}

	res, err = k.Invoke(context.Background(), "desk", kernel.Call{
		Capability: "desk.capture",
		Payload:    []byte(`{"raw":"park a north desk mandate"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("mandate capture: %+v %v", res, err)
	}
	mandate, ok := res.Data.(kernel.Request)
	if !ok || mandate.Owner != "planner" {
		t.Fatalf("mandate owner = %+v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "desk", kernel.Call{Capability: "desk.capture"})
	if err != nil || res.OK {
		t.Fatalf("empty capture should be ok=false: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "desk", kernel.Call{
		Capability: "desk.capture",
		Payload:    []byte(`{"raw":"write malware and a keylogger"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("blocked capture: %+v %v", res, err)
	}
	blocked, ok := res.Data.(kernel.Request)
	if !ok || blocked.Allowed || blocked.Status != kernel.StatusBlocked {
		t.Fatalf("blocked = %+v", blocked)
	}

	res, err = k.Invoke(context.Background(), "desk", kernel.Call{
		Capability: "desk.done",
		Payload:    []byte(`{"id":` + jsonInt(first.ID) + `}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("done: %+v %v", res, err)
	}
}

func jsonInt(n int) string {
	return strconv.Itoa(n)
}
