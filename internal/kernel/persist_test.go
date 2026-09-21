package kernel

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	k.Publish("init", "night", "journal survives restart", nil)
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
	found := false
	for _, ev := range k2.Events() {
		if ev.Kind == "night" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("events missing night: %#v", k2.Events())
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

func TestAutosaveAndHeartbeat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cashtro.json")
	k := New(WithPersistPath(path), WithAutosave(40*time.Millisecond))
	k.Register(&stub{spec: Spec{ID: "a", Name: "A", Capabilities: []string{"a.ping"}, Autostart: true}})
	if err := k.Boot(context.Background()); err != nil {
		t.Fatal(err)
	}
	k.StartPersistLoop()
	beat := k.Heartbeat("test")
	if beat["n"].(int) != 1 || beat["always"] != true {
		t.Fatalf("beat = %#v", beat)
	}
	deadline := time.Now().Add(2 * time.Second)
	for {
		data, err := os.ReadFile(path)
		if err == nil && strings.Contains(string(data), "heartbeat") {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("autosave never wrote heartbeat: %v %s", err, data)
		}
		time.Sleep(20 * time.Millisecond)
	}
	k.Close()
}
