package kernel

import (
	"context"
	"errors"
	"testing"
)

func TestMailboxNotesMemoryConfirm(t *testing.T) {
	k := New()
	k.Register(&stub{spec: Spec{ID: "a", Name: "A", Capabilities: []string{"a.ping"}, Autostart: true}})
	k.Register(&stub{spec: Spec{ID: "b", Name: "B", Capabilities: []string{"b.ping"}, Autostart: true}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}

	if _, err := k.Post("a", "missing", "ping", "x"); !errors.Is(err, ErrUnknownAgent) {
		t.Fatalf("post missing = %v", err)
	}
	m, err := k.Post("a", "b", "job", "walk the line")
	if err != nil || m.To != "b" {
		t.Fatalf("post = %+v %v", m, err)
	}
	if len(k.Inbox("b")) != 1 {
		t.Fatalf("inbox = %#v", k.Inbox("b"))
	}

	n := k.WriteNote(Note{Agent: "research", Source: "arxiv", URL: "https://arxiv.org/abs/2606.01508", Claim: "AOS is a control plane", Quote: "not a replacement for Linux"})
	if n.ID == 0 || len(k.Notes()) != 1 {
		t.Fatalf("note = %+v", n)
	}

	k.Remember("os", "OpenRouter is optional")
	got := k.Recall("openrouter")
	if len(got) != 1 {
		t.Fatalf("recall = %#v", got)
	}

	c := k.RequestConfirm("comms", "comms.send", "ping Castro")
	if c.Status != "pending" {
		t.Fatalf("confirm = %+v", c)
	}
	decided, err := k.DecideConfirm(c.ID, true)
	if err != nil || decided.Status != "allowed" {
		t.Fatalf("decide = %+v %v", decided, err)
	}

	if _, err := k.SaveRequest(Request{}); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("empty request = %v", err)
	}
	req, err := k.SaveRequest(Request{Raw: "  ship the desk  ", Title: "Ship the desk", Improved: "first pass", Pass: 1, Allowed: true})
	if err != nil || req.ID == 0 || req.Raw != "ship the desk" {
		t.Fatalf("capture = %+v %v", req, err)
	}
	req.Pass = 2
	req.Improved = "second pass"
	req.Raw = "should not stick"
	bettered, err := k.SaveRequest(req)
	if err != nil || bettered.Raw != "ship the desk" || bettered.Pass != 2 || bettered.Improved != "second pass" {
		t.Fatalf("better = %+v %v", bettered, err)
	}
	done, err := k.AdvanceRequest(bettered.ID, StatusDone)
	if err != nil || done.Status != StatusDone {
		t.Fatalf("advance = %+v %v", done, err)
	}
	if len(k.Requests()) != 1 {
		t.Fatalf("requests = %#v", k.Requests())
	}
	if _, err := k.RequestByID(99); !errors.Is(err, ErrUnknownRequest) {
		t.Fatalf("missing = %v", err)
	}
}

func TestInvokeCap(t *testing.T) {
	k := New()
	k.Register(&stub{spec: Spec{ID: "delivery", Name: "D", Capabilities: []string{"delivery.list"}, Autostart: true}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}
	res, err := k.InvokeCap(context.Background(), "delivery.list", Call{})
	if err != nil || !res.OK {
		t.Fatalf("invoke cap: %+v %v", res, err)
	}
	_, err = k.InvokeCap(context.Background(), "nope", Call{})
	if !errors.Is(err, ErrUnknownCapability) {
		t.Fatalf("err = %v", err)
	}
}
