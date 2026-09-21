package flow

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestSeededEngineOwnsTheLine(t *testing.T) {
	e := New()
	list := e.List()
	if len(list) != 3 {
		t.Fatalf("flows = %d, want 3", len(list))
	}
	wf, err := e.Get("wealth-dual")
	if err != nil {
		t.Fatal(err)
	}
	if wf.Trigger != KindManual || !wf.Enabled {
		t.Fatalf("wealth-dual = %+v", wf)
	}
	hook, err := e.Get("intake-hook")
	if err != nil {
		t.Fatal(err)
	}
	if hook.Hook != "/api/n8n/intake-hook/hook" {
		t.Fatalf("hook = %q", hook.Hook)
	}
	kinds := Catalog()
	if len(kinds) < 8 {
		t.Fatalf("catalog = %d", len(kinds))
	}
}

func TestExecuteDualIsLocal(t *testing.T) {
	now := time.Date(2026, 9, 21, 23, 0, 0, 0, time.UTC)
	e := New(WithClock(func() time.Time { return now }))
	var caps []string
	invoker := func(ctx context.Context, cap string, fields map[string]string) (bool, string, error) {
		caps = append(caps, cap)
		if fields["text"] == "" && fields["claim"] == "" {
			t.Fatalf("empty fields for %s: %#v", cap, fields)
		}
		return true, "ok " + cap, nil
	}
	if err := e.Enter("wealth-dual"); err != nil {
		t.Fatal(err)
	}
	defer e.Leave("wealth-dual")

	run, err := e.Execute(context.Background(), "wealth-dual", map[string]string{
		"prompt": "Advance ScanApp. Keep confirms.",
	}, invoker)
	if err != nil {
		t.Fatal(err)
	}
	if !run.OK || run.Bound || run.Owner != "cashtro" {
		t.Fatalf("run = %+v", run)
	}
	if run.Score < 70 {
		t.Fatalf("score = %d", run.Score)
	}
	if !strings.Contains(run.Item["propose"], "human gate") {
		t.Fatalf("propose = %q", run.Item["propose"])
	}
	if got := strings.Join(caps, ","); got != "memory.store,research.ingest" {
		t.Fatalf("caps = %s", got)
	}
	if len(e.Runs()) != 1 {
		t.Fatalf("runs = %d", len(e.Runs()))
	}
}

func TestWebhookBranchAndToggle(t *testing.T) {
	e := New()
	called := 0
	invoker := func(ctx context.Context, cap string, fields map[string]string) (bool, string, error) {
		called++
		return true, cap, nil
	}
	run, err := e.Execute(context.Background(), "intake-hook", map[string]string{"prompt": "new lead"}, invoker)
	if err != nil || !run.OK {
		t.Fatalf("hook: %+v %v", run, err)
	}
	if called != 2 {
		t.Fatalf("called = %d, want memory+confirm", called)
	}

	empty, err := e.Execute(context.Background(), "intake-hook", map[string]string{}, invoker)
	if err != nil || !empty.OK {
		t.Fatalf("empty: %+v %v", empty, err)
	}
	if empty.Item["skipped"] != "true" {
		t.Fatalf("empty item = %#v", empty.Item)
	}

	off, err := e.Toggle("intake-hook")
	if err != nil || off.Enabled {
		t.Fatalf("toggle = %+v %v", off, err)
	}
	if err := e.Enter("intake-hook"); err != ErrDisabled {
		t.Fatalf("enter disabled = %v", err)
	}
}

func TestCreateAndRender(t *testing.T) {
	e := New()
	wf, err := e.Create(Create{
		Name:    "North desk",
		Summary: "owned",
		Trigger: KindWebhook,
		Nodes: []Node{
			{ID: "hook", Kind: KindWebhook, Name: "Catch"},
			{ID: "set", Kind: KindSet, Name: "Stamp", Params: map[string]string{"hello": "hi {{prompt}}"}},
		},
		Edges: []Edge{{From: "hook", To: "set"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if wf.ID != "north-desk" || wf.Hook != "/api/n8n/north-desk/hook" {
		t.Fatalf("created = %+v", wf)
	}
	got := Render(map[string]string{"x": "n={{prompt}}"}, map[string]string{"prompt": "scan"})
	if got["x"] != "n=scan" {
		t.Fatalf("render = %#v", got)
	}
	flat := FlattenJSON([]byte(`{"prompt":"hi","n":2}`))
	if flat["prompt"] != "hi" || flat["n"] != "2" {
		t.Fatalf("flat = %#v", flat)
	}
}

func TestCreateRejectsBadGraph(t *testing.T) {
	e := New()
	_, err := e.Create(Create{Name: "x"})
	if err == nil {
		t.Fatal("expected empty nodes error")
	}
	_, err = e.Create(Create{
		Name:  "x",
		Nodes: []Node{{ID: "a", Kind: "nope"}},
	})
	if err == nil {
		t.Fatal("expected unknown kind")
	}
}
