package kernel

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPersistRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cashtro.json")
	k := New(WithPersistPath(path))
	k.Register(&stub{spec: Spec{ID: "a", Name: "A", Capabilities: []string{"a.ping"}, Autostart: true}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}
	k.WriteNote(Note{Agent: "research", Source: "test", Claim: "persist me"})
	k.Remember("os", "disk image")
	k.SetClosed(true)
	k.RecordPulse("overnight")
	k.Offer(Choice{Key: "mail-persist", Kind: KindMail, Title: "Persist me"})
	k.RecordScan(Scan{Account: "alejandro@proximityagency.ca", Inbox: 3, Unread: 1})
	if err := SaveFile(path, k); err != nil {
		t.Fatal(err)
	}

	k2 := New(WithPersistPath(path))
	k2.Register(&stub{spec: Spec{ID: "a", Name: "A", Capabilities: []string{"a.ping"}, Autostart: true}})
	if err := k2.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := LoadFile(path, k2); err != nil {
		t.Fatal(err)
	}
	if len(k2.Notes()) != 1 || k2.Notes()[0].Claim != "persist me" {
		t.Fatalf("notes = %#v", k2.Notes())
	}
	if len(k2.Recall("disk")) != 1 {
		t.Fatalf("facts = %#v", k2.Recall("disk"))
	}
	if !k2.Closed() {
		t.Fatal("closed hours did not restore")
	}
	if got := k2.Pulses(); len(got) != 1 || got[0].Note != "overnight" {
		t.Fatalf("pulses = %#v", k2.Pulses())
	}
	if got := k2.Choices(""); len(got) != 1 || got[0].Key != "mail-persist" {
		t.Fatalf("choices = %#v", k2.Choices(""))
	}
	if k2.LastScan().Unread != 1 {
		t.Fatalf("scan = %#v", k2.LastScan())
	}
}

func TestLoadMissingFile(t *testing.T) {
	k := New()
	if err := LoadFile(filepath.Join(t.TempDir(), "nope.json"), k); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected missing")
	}
}
