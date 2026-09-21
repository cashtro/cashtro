package kernel

import (
	"errors"
	"testing"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
)

func TestChooserTakeAndSkip(t *testing.T) {
	now := time.Date(2026, 9, 21, 5, 0, 0, 0, time.UTC)
	k := New(WithClock(func() time.Time { return now }))
	k.AttachCatalog(catalog.New(catalog.WithClock(func() time.Time { return now })))

	k.RecordScan(Scan{Account: "alejandro@proximityagency.ca", Device: "Samsung Gmail", Inbox: 620, Unread: 163})
	mail := k.Offer(Choice{Key: "mail-clic", Kind: KindMail, Title: "Clic Inspection v2", Client: "Proximity", Summary: "zip landed"})
	verb := k.Offer(Choice{Key: "verb-pulse", Kind: KindVerb, Title: "Pulse the desk"})
	dup := k.Offer(Choice{Key: "mail-clic", Kind: KindMail, Title: "should not duplicate"})
	if dup.ID != mail.ID || dup.Title != "Clic Inspection v2" {
		t.Fatalf("offer should be idempotent: %+v", dup)
	}

	desk := k.DeskCard()
	if desk.Scan.Unread != 163 || desk.PendingN != 2 {
		t.Fatalf("desk = %+v", desk)
	}

	taken, err := k.Take(mail.ID)
	if err != nil || taken.Status != StatusTaken || taken.ShipID == "" {
		t.Fatalf("take = %+v %v", taken, err)
	}
	if _, err := k.Catalog().Get(taken.ShipID); err != nil {
		t.Fatalf("ship: %v", err)
	}
	if _, err := k.Take(mail.ID); !errors.Is(err, ErrAlreadyDecided) {
		t.Fatalf("second take err = %v", err)
	}

	skipped, err := k.Skip(verb.ID)
	if err != nil || skipped.Status != StatusSkipped {
		t.Fatalf("skip = %+v %v", skipped, err)
	}
	if _, err := k.Skip(99); !errors.Is(err, ErrUnknownChoice) {
		t.Fatalf("missing skip err = %v", err)
	}

	desk = k.DeskCard()
	if len(desk.Pending) != 0 || len(desk.Taken) != 1 || len(desk.Skipped) != 1 {
		t.Fatalf("after decide desk = %+v", desk)
	}
}
