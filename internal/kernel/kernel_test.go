package kernel

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stub struct {
	spec    Spec
	booted  bool
	bootErr error
	last    Call
}

func (s *stub) Spec() Spec { return s.spec }
func (s *stub) Boot(ctx context.Context, k *Kernel) error {
	s.booted = true
	return s.bootErr
}
func (s *stub) Invoke(ctx context.Context, call Call) (Result, error) {
	s.last = call
	return Result{OK: true, Message: "ok", Data: call.Capability}, nil
}

func TestBootAndInvoke(t *testing.T) {
	now := time.Date(2026, 9, 19, 22, 0, 0, 0, time.UTC)
	k := New(WithClock(func() time.Time { return now }))
	live := &stub{spec: Spec{
		ID: "delivery", Name: "Delivery", Kind: KindUser, Mode: ModeLive,
		Capabilities: []string{"delivery.list"}, Autostart: true,
	}}
	resident := &stub{spec: Spec{
		ID: "explorer", Name: "Explorer", Kind: KindUser, Mode: ModeResident,
		Capabilities: []string{"explorer.search"}, Autostart: true,
	}}
	k.Register(live)
	k.Register(resident)

	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !live.booted || !resident.booted {
		t.Fatal("expected autostart boot")
	}

	about := k.About()
	if about.Name != Name || about.Agents != 2 || about.Running != 2 || about.Live != 1 {
		t.Fatalf("about = %+v", about)
	}

	res, err := k.Invoke(context.Background(), "delivery", Call{Capability: "delivery.list"})
	if err != nil || !res.OK {
		t.Fatalf("invoke: %+v %v", res, err)
	}

	_, err = k.Invoke(context.Background(), "delivery", Call{Capability: "nope"})
	if !errors.Is(err, ErrUnknownCapability) {
		t.Fatalf("cap err = %v", err)
	}
	_, err = k.Invoke(context.Background(), "missing", Call{Capability: "x"})
	if !errors.Is(err, ErrUnknownAgent) {
		t.Fatalf("agent err = %v", err)
	}

	evs := k.Events()
	if len(evs) < 3 {
		t.Fatalf("events = %d", len(evs))
	}
	if evs[0].Kind != "boot" {
		t.Fatalf("first event = %+v", evs[0])
	}
}

func TestSpawnIdempotentAndBootError(t *testing.T) {
	k := New()
	ok := &stub{spec: Spec{ID: "a", Name: "A", Autostart: false, Capabilities: []string{"a.ping"}}}
	bad := &stub{spec: Spec{ID: "b", Name: "B", Autostart: false}, bootErr: errors.New("fail")}
	k.Register(ok)
	k.Register(bad)

	if err := k.Spawn(context.Background(), "a"); err != nil {
		t.Fatal(err)
	}
	if err := k.Spawn(context.Background(), "a"); err != nil {
		t.Fatal(err)
	}
	if err := k.Spawn(context.Background(), "b"); err == nil {
		t.Fatal("expected boot error")
	}
	p, err := k.Process("b")
	if err != nil || p.Status != StatusError {
		t.Fatalf("bad process = %+v %v", p, err)
	}
}

func TestInvokeRequiresRunning(t *testing.T) {
	k := New()
	k.Register(&stub{spec: Spec{ID: "a", Name: "A", Capabilities: []string{"a.ping"}}})
	_, err := k.Invoke(context.Background(), "a", Call{Capability: "a.ping"})
	if !errors.Is(err, ErrNotRunning) {
		t.Fatalf("err = %v", err)
	}
}

func TestKeepAliveRespawns(t *testing.T) {
	k := New()
	k.Register(&stub{spec: Spec{ID: "a", Name: "A", Capabilities: []string{"a.ping"}, Autostart: true}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := k.Stop("a"); err != nil {
		t.Fatal(err)
	}
	spawned, err := k.KeepAlive(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(spawned) != 1 || spawned[0] != "a" {
		t.Fatalf("spawned = %#v", spawned)
	}
	p, _ := k.Process("a")
	if p.Status != StatusRunning {
		t.Fatalf("status = %s", p.Status)
	}
}
