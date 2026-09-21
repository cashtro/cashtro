package agents

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBootLoadsAllAgentics(t *testing.T) {
	t.Setenv("CASHTRO_CLOSED", "")
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	procs := k.Processes()
	if len(procs) != 15 {
		t.Fatalf("processes = %d, want 15", len(procs))
	}
	about := k.About()
	if about.Running != 15 || about.Live != 10 {
		t.Fatalf("about = %+v", about)
	}
	if len(k.Notes()) < 5 {
		t.Fatalf("research seeds = %d", len(k.Notes()))
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

	res, err = k.Invoke(context.Background(), "watch", kernel.Call{
		Capability: "watch.close",
		Payload:    []byte(`{"note":"things are closed"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("watch.close: %+v %v", res, err)
	}
	if !k.Closed() {
		t.Fatal("desk should be closed")
	}
	if _, err := k.Catalog().Get("closed-hours-flow"); err != nil {
		t.Fatalf("closed-hours ship: %v", err)
	}
	res, err = k.Invoke(context.Background(), "watch", kernel.Call{Capability: "watch.pulse"})
	if err != nil || !res.OK {
		t.Fatalf("watch.pulse: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "watch", kernel.Call{Capability: "watch.open"})
	if err != nil || !res.OK || k.Closed() {
		t.Fatalf("watch.open: %+v closed=%v err=%v", res, k.Closed(), err)
	}

	res, err = k.Invoke(context.Background(), "architect", kernel.Call{
		Capability: "architect.plan",
		Payload:    []byte(`{"goal":"agent OS control plane with outbound comms"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("architect.plan: %+v %v", res, err)
	}
	brief, ok := res.Data.(Brief)
	if !ok || brief.Shape != "os" || brief.Goal == "" {
		t.Fatalf("brief = %#v ok=%v", res.Data, ok)
	}
	if !containsStr(brief.Agents, "comms") || !containsStr(brief.Gates, "comms.send") {
		t.Fatalf("brief agents/gates = %#v", brief)
	}
	if len(k.Recall("design")) == 0 {
		t.Fatal("architect did not remember the brief")
	}

	res, err = k.Invoke(context.Background(), "investigator", kernel.Call{
		Capability: "investigator.trace",
		Payload:    []byte(`{"query":"watch"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("investigator.trace: %+v %v", res, err)
	}
	trace, ok := res.Data.(Trace)
	if !ok || trace.Query != "watch" || len(trace.Hits) == 0 {
		t.Fatalf("trace = %#v ok=%v", res.Data, ok)
	}
	if !traceHasKind(trace.Hits, "agent") {
		t.Fatalf("trace missed watch agent: %#v", trace.Hits)
	}
}

func traceHasKind(hits []TraceHit, kind string) bool {
	for _, h := range hits {
		if h.Kind == kind {
			return true
		}
	}
	return false
}

func containsStr(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
