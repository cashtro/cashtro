package kernel

import (
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

	k.SetClosed(false)
	if k.Closed() {
		t.Fatal("expected desk open")
	}
	if got := k.Pulses(); len(got) != 1 {
		t.Fatalf("pulses = %+v", got)
	}
}
