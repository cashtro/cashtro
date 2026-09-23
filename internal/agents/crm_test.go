package agents

import (
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestOneLovableCRMAndDraftGrowth(t *testing.T) {
	book := CRM()
	if book.Repo != "Evolu-Jeunes/CRM" || book.Writes || book.Sent {
		t.Fatalf("book %+v", book)
	}
	if len(book.Lines) != 4 || len(book.Next) < 4 {
		t.Fatalf("plan %+v", book)
	}
	k := kernel.New()
	note, err := CRMGrow(k, "marketing", "lead", "Campagne Pandora", "mesurer avant d'élargir", false)
	if err != nil || note.Status != "brouillon" || note.Sent || note.DecidedBy != "epicenter" || !note.Remembered {
		t.Fatalf("note %+v err %v", note, err)
	}
	if _, err := CRMGrow(k, "marketing", "lead", "x", "y", true); err == nil {
		t.Fatal("send should be refused")
	}
	if _, err := CRMGrow(k, "wordpress", "note", "thème", "page", false); err == nil {
		t.Fatal("wordpress should stay out of lovable")
	}
	if _, err := CRMGrow(k, "nft-giant", "note", "ordre", "btc", false); err == nil {
		t.Fatal("giant should stay out of this crm")
	}
	if _, err := CRMGrow(k, "proximity", "lead", "client", "nom", false); err == nil {
		t.Fatal("proximity lead should be refused")
	}
	got, err := CRMGrow(k, "proximity", "contact", "Client Nord", "nom seulement", false)
	if err != nil || got.Kind != "contact" {
		t.Fatalf("proximity contact %+v %v", got, err)
	}
	if _, err := CRMGrow(k, "marketing", "note", "clé", "sk_live_abc", false); err == nil {
		t.Fatal("secret should be filtered")
	}
}
