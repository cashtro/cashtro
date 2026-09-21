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

	_, err = k.OpenTeamThread(TeamOpen{Channel: "slack", Tenant: "Proximity", Kind: "conversation"})
	if !errors.Is(err, ErrTeamChannel) {
		t.Fatalf("slack = %v", err)
	}
	_, err = k.OpenTeamThread(TeamOpen{Channel: "teams.microsfot", Tenant: "lroximity", Kind: "meeting"})
	if !errors.Is(err, ErrTeamKind) {
		t.Fatalf("meeting = %v", err)
	}
	_, err = k.OpenTeamThread(TeamOpen{Channel: "teams.microsoft", Tenant: "other-co", Kind: "conversation"})
	if !errors.Is(err, ErrTeamTenant) {
		t.Fatalf("tenant = %v", err)
	}
	th, err := k.OpenTeamThread(TeamOpen{ID: "proximity-desk", Title: "Proximity desk", Channel: "teams.microsfot", Tenant: "lroximity", Kind: "conversation"})
	if err != nil || th.Channel != TeamChannelMicrosoft || th.Tenant != TeamTenantProximity || th.Dock != "bottom" {
		t.Fatalf("open = %+v %v", th, err)
	}
	th, err = k.AppendTeamMessage(th.ID, "user", "status on the agency platform")
	if err != nil || len(th.Messages) != 1 {
		t.Fatalf("say = %+v %v", th, err)
	}
	if k.TeamCard().Threads != 1 {
		t.Fatalf("card = %+v", k.TeamCard())
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
