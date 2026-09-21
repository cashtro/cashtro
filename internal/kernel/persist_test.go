package kernel

import (
	"path/filepath"
	"testing"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/fleet"
)

func TestPersistRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "os.json")
	k := New(WithPersistPath(path))
	cat := catalog.New()
	k.AttachCatalog(cat)
	k.SetClosed(true)
	k.RecordPulse("night")
	k.WriteNote(Note{Agent: "watch", Source: "test", Claim: "persisted"})
	if err := SaveFile(path, k); err != nil {
		t.Fatal(err)
	}

	k2 := New(WithPersistPath(path))
	cat2 := catalog.New()
	k2.AttachCatalog(cat2)
	if err := LoadFile(path, k2); err != nil {
		t.Fatal(err)
	}
	if !k2.Closed() {
		t.Fatal("closed should restore")
	}
	if len(k2.Pulses()) != 1 || len(k2.Notes()) == 0 {
		t.Fatalf("pulses=%d notes=%d", len(k2.Pulses()), len(k2.Notes()))
	}
	if len(cat2.List()) == 0 {
		t.Fatal("ships should restore")
	}
}

func TestPersistRestoresCompanies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "os.json")
	k := New(WithPersistPath(path))
	f := fleet.New()
	k.AttachTenants(f)
	if _, err := f.Create(fleet.CreateCompany{Name: "West Desk", Sector: "ops"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveFile(path, k); err != nil {
		t.Fatal(err)
	}

	k2 := New(WithPersistPath(path))
	f2 := fleet.New()
	k2.AttachTenants(f2)
	if err := LoadFile(path, k2); err != nil {
		t.Fatal(err)
	}
	got, err := f2.Get("west-desk")
	if err != nil || got.Name != "West Desk" {
		t.Fatalf("company restore = %+v %v", got, err)
	}
}
