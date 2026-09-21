package symbols

import (
	"errors"
	"testing"
	"time"
)

func TestSeededBus(t *testing.T) {
	b := New()
	list := b.List()
	if len(list) != 4 {
		t.Fatalf("list = %d, want 4", len(list))
	}
	got, err := b.Get("intake-to-idea")
	if err != nil {
		t.Fatal(err)
	}
	if got.On.Kind != TriggerWebhook || len(got.Steps) != 2 {
		t.Fatalf("intake = %+v", got)
	}
	_, err = b.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing = %v", err)
	}
}

func TestCreateToggleFinish(t *testing.T) {
	now := time.Date(2026, 9, 21, 8, 0, 0, 0, time.UTC)
	b := New(WithClock(func() time.Time { return now }))

	s, err := b.Create(Create{
		Name:    "  Hook desk  ",
		Summary: "test",
		On:      Trigger{Kind: TriggerManual},
		Steps:   []Step{{Cap: "os.about"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "hook-desk" || !s.Enabled {
		t.Fatalf("created = %+v", s)
	}

	toggled, err := b.Toggle(s.ID)
	if err != nil || toggled.Enabled {
		t.Fatalf("toggle = %+v %v", toggled, err)
	}
	if err := b.Enter(s.ID); !errors.Is(err, ErrDisabled) {
		t.Fatalf("enter disabled = %v", err)
	}

	if _, err := b.Toggle(s.ID); err != nil {
		t.Fatal(err)
	}
	if err := b.Enter(s.ID); err != nil {
		t.Fatal(err)
	}
	if err := b.Enter(s.ID); !errors.Is(err, ErrBusy) {
		t.Fatalf("busy = %v", err)
	}
	b.Leave(s.ID)

	run, err := b.Finish(s.ID, true, "ok", map[string]string{"name": "X"}, []StepLog{{Cap: "os.about", OK: true, Message: "about"}})
	if err != nil || !run.OK || run.ID != 1 {
		t.Fatalf("run = %+v %v", run, err)
	}
	got, _ := b.Get(s.ID)
	if got.Fires != 1 || !got.LastOK {
		t.Fatalf("after finish = %+v", got)
	}
	if len(b.Runs()) != 1 {
		t.Fatalf("runs = %d", len(b.Runs()))
	}
}

func TestCreateValidation(t *testing.T) {
	b := New()
	_, err := b.Create(Create{Name: "  "})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("empty name = %v", err)
	}
	_, err = b.Create(Create{Name: "X"})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("no steps = %v", err)
	}
	_, err = b.Create(Create{Name: "X", Steps: []Step{{Cap: "os.about"}}, On: Trigger{Kind: "zapier"}})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("bad trigger = %v", err)
	}
	_, err = b.Create(Create{Name: "X", Steps: []Step{{Cap: "os.about"}}, On: Trigger{Kind: TriggerEvent}})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("event without match = %v", err)
	}
}

func TestRenderAndFlatten(t *testing.T) {
	in := FlattenJSON([]byte(`{"name":"North desk","client":"Cashtro","n":2}`))
	if in["name"] != "North desk" || in["n"] != "2" {
		t.Fatalf("flatten = %#v", in)
	}
	fields := Render(map[string]string{
		"name":  "{{name}}",
		"notes": "parked {{missing}}",
	}, in)
	if fields["name"] != "North desk" || fields["notes"] != "parked " {
		t.Fatalf("render = %#v", fields)
	}
	if FlattenJSON(nil)["body"] != "" && len(FlattenJSON(nil)) != 0 {
		t.Fatalf("empty = %#v", FlattenJSON(nil))
	}
	raw := FlattenJSON([]byte("not json"))
	if raw["body"] != "not json" {
		t.Fatalf("raw = %#v", raw)
	}
}

func TestMatchEvent(t *testing.T) {
	on := Trigger{Kind: TriggerEvent, Match: "mail"}
	if !MatchEvent(on, "explorer", "mail") {
		t.Fatal("mail should match")
	}
	if MatchEvent(on, "delivery", "invoke") {
		t.Fatal("invoke must never auto-fire")
	}
	if MatchEvent(Trigger{Kind: TriggerWebhook, Match: "mail"}, "explorer", "mail") {
		t.Fatal("webhook trigger is not an event")
	}
	if MatchEvent(Trigger{Kind: TriggerEvent, Match: "delivery"}, "delivery", "invoke") {
		t.Fatal("invoke skipped even when source matches")
	}
	if !MatchEvent(Trigger{Kind: TriggerEvent, Match: "explorer/mail"}, "explorer", "mail") {
		t.Fatal("source/kind should match")
	}
}
