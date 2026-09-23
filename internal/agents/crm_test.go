package agents

import (
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestAgentCRMStaysApartFromFix2(t *testing.T) {
	book := CRM()
	if book.Repo != "Evolu-Jeunes/CRM" || book.Voice != "vapi" || book.Writes || book.Sent || book.Dialed {
		t.Fatalf("book %+v", book)
	}
	if len(book.Lines) != len(Lines())-1 || len(book.Next) < 5 {
		t.Fatalf("plan %+v", book)
	}
	for _, id := range book.Lines {
		if id == "fix2" {
			t.Fatal("fix2 must stay out of the agent book")
		}
	}
	k := kernel.New()
	for _, ln := range Lines() {
		if ln.ID == "fix2" {
			continue
		}
		note, err := CRMGrow(k, ln.ID, "note", ln.Name, "brouillon du conglomerat", false)
		if err != nil || note.Status != "brouillon" || note.Sent || note.DecidedBy != "epicenter" || !note.Remembered {
			t.Fatalf("%s note %+v err %v", ln.ID, note, err)
		}
	}
	if _, err := CRMGrow(k, "fix2", "contact", "Cuisine", "chantier", false); err == nil {
		t.Fatal("fix2 should stay in its own book")
	}
	lead, err := CRMGrow(k, "marketing", "lead", "Campagne Pandora", "mesurer avant d'élargir", false)
	if err != nil || lead.Kind != "lead" {
		t.Fatalf("lead %+v %v", lead, err)
	}
	if _, err := CRMGrow(k, "marketing", "lead", "x", "y", true); err == nil {
		t.Fatal("send should be refused")
	}
	wp, err := CRMGrow(k, "wordpress", "note", "Client Nord", "nom du site, pas le thème", false)
	if err != nil || wp.Line != "wordpress" {
		t.Fatalf("wordpress note %+v %v", wp, err)
	}
	giant, err := CRMGrow(k, "nft-giant", "note", "Lecture GNT", "pas un ordre", false)
	if err != nil || giant.Kind != "note" {
		t.Fatalf("giant note %+v %v", giant, err)
	}
	if _, err := CRMGrow(k, "marketing", "note", "clé", "sk_live_abc", false); err == nil {
		t.Fatal("secret should be filtered")
	}
	if _, err := CRMGrow(k, "ailleurs", "note", "x", "y", false); err == nil {
		t.Fatal("unknown line should be refused")
	}
}
