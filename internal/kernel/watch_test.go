package kernel

import (
	"context"
	"testing"
	"time"
)

func TestClosedHoursPulse(t *testing.T) {
	now := time.Date(2026, 9, 21, 4, 0, 0, 0, time.UTC)
	k := New(WithClock(func() time.Time { return now }))

	if k.Closed() {
		t.Fatal("desk should boot open")
	}
	k.SetClosed(true)
	if !k.Closed() {
		t.Fatal("expected closed hours")
	}
	p := k.RecordPulse("things are closed")
	if p.ID != 1 || !p.Closed || p.Note != "things are closed" {
		t.Fatalf("pulse = %+v", p)
	}

	card := k.WatchCard()
	if !card.Closed || card.LastPulse == nil || card.LastPulse.ID != 1 {
		t.Fatalf("card = %+v", card)
	}
	if len(card.Flowing) == 0 || len(card.Gated) == 0 {
		t.Fatal("duty roster missing")
	}
	if !k.About().Closed {
		t.Fatal("about should show closed")
	}

	k.SetClosed(false)
	if k.Closed() {
		t.Fatal("expected desk open")
	}
}

func TestRunClosedPulses(t *testing.T) {
	k := New()
	k.SetClosed(true)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		k.RunClosedPulses(ctx, 8*time.Millisecond)
		close(done)
	}()
	deadline := time.Now().Add(400 * time.Millisecond)
	for time.Now().Before(deadline) {
		if len(k.Pulses()) >= 1 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if len(k.Pulses()) < 1 {
		t.Fatal("watchdog did not pulse")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("watchdog did not stop")
	}
}
