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
	foundPulse := false
	for _, ev := range k2.Events() {
		if ev.Kind == "pulse" && ev.Source == "watch" {
			foundPulse = true
			break
		}
	}
	if !foundPulse {
		t.Fatalf("journal lost pulse event: %#v", k2.Events())
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
