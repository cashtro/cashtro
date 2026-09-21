package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
)

const maxBuilds = 32

var nightClaims = []string{
	"Closed hours: Watch keeps Castro's agentics on the Giant desk.",
	"Companies without a roster inherit Castro's mine.",
	"Outbound comms stay gated while the desk is closed.",
	"Delivery still moves idea → concept → production overnight.",
	"Always-on builder parks Cashtro work so Castro can sleep.",
}

// Build is one overnight shift the agentics ran while the desk was closed.
type Build struct {
	ID    int         `json:"id"`
	At    time.Time   `json:"at"`
	Note  string      `json:"note"`
	Steps []BuildStep `json:"steps"`
}

// BuildStep is one verb the night shift invoked.
type BuildStep struct {
	Agent      string `json:"agent"`
	Capability string `json:"capability"`
	OK         bool   `json:"ok"`
	Message    string `json:"message"`
}

type tenantCounter interface {
	Count() (companies int, usingMine int)
}

// NightShift runs Castro's builder loop: park/advance Cashtro work,
// ingest a note, remember it. Never sends comms or deploys.
func (k *Kernel) NightShift(ctx context.Context, note string) (Build, error) {
	if note == "" {
		note = "overnight build"
	}
	b := Build{At: k.now(), Note: note, Steps: make([]BuildStep, 0, 8)}

	step := func(agent, cap string, payload any) {
		var raw json.RawMessage
		if payload != nil {
			raw, _ = json.Marshal(payload)
		}
		res, err := k.Invoke(ctx, agent, Call{Capability: cap, Payload: raw})
		msg := cap
		ok := err == nil && res.OK
		if err != nil {
			msg = err.Error()
		} else if res.Message != "" {
			msg = res.Message
		}
		b.Steps = append(b.Steps, BuildStep{Agent: agent, Capability: cap, OK: ok, Message: msg})
	}

	k.Remember("watch", "night shift · "+note)
	b.Steps = append(b.Steps, BuildStep{Agent: "memory", Capability: "memory.store", OK: true, Message: "remembered watch"})

	if cat := k.Catalog(); cat != nil {
		if ship, ok := nextOvernight(cat); ok {
			advanced, err := cat.Advance(ship.ID)
			if err != nil {
				b.Steps = append(b.Steps, BuildStep{Agent: "delivery", Capability: "delivery.advance", OK: false, Message: err.Error()})
			} else {
				b.Steps = append(b.Steps, BuildStep{Agent: "delivery", Capability: "delivery.advance", OK: true, Message: advanced.Name + " → " + string(advanced.Stage)})
			}
		} else {
			parkOvernight(cat, note, &b)
		}
		if overnightOpen(cat) < 2 {
			parkOvernight(cat, "keep the Giant line moving", &b)
		}
	}

	if n, mine, ok := tenantCount(k); ok {
		text := fmt.Sprintf("%d companies · %d using Castro", n, mine)
		k.Remember("fleet", text)
		b.Steps = append(b.Steps, BuildStep{Agent: "memory", Capability: "memory.store", OK: true, Message: text})
	}

	claim := nightClaims[len(k.Builds())%len(nightClaims)]
	step("research", "research.ingest", map[string]string{
		"claim":  claim + " " + note,
		"source": "watch.night",
		"url":    "",
	})
	step("explorer", "explorer.search", map[string]string{"query": "overnight cashtro giant"})
	step("memory", "memory.recall", map[string]string{"query": "watch"})
	_, _ = k.Post("watch", "planner", "pulse", note)
	_, _ = k.Post("watch", "research", "pulse", "keep gathering while the desk is closed")

	k.mu.Lock()
	k.buildSeq++
	b.ID = k.buildSeq
	k.builds = append(k.builds, b)
	if len(k.builds) > maxBuilds {
		k.builds = append([]Build(nil), k.builds[len(k.builds)-maxBuilds:]...)
	}
	k.seq++
	k.events = append(k.events, Event{
		Seq: k.seq, At: k.now(), Source: "watch", Kind: "build",
		Message: "night shift · " + note, Data: map[string]any{"id": b.ID, "steps": len(b.Steps)},
	})
	k.mu.Unlock()
	k.persist()
	return b, nil
}

// Builds returns night shifts, newest first.
func (k *Kernel) Builds() []Build {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Build, len(k.builds))
	copy(out, k.builds)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func nextOvernight(cat *catalog.Catalog) (catalog.Ship, bool) {
	for _, s := range cat.List() {
		if s.Stage == catalog.StageProduction {
			continue
		}
		if isOvernight(s) {
			return s, true
		}
	}
	return catalog.Ship{}, false
}

func isOvernight(s catalog.Ship) bool {
	if s.Client == "Cashtro" || s.Client == "Cashtro OS" {
		return true
	}
	blob := strings.ToLower(s.Name + " " + s.Notes)
	return strings.Contains(blob, "overnight") || strings.Contains(blob, "closed-hours") || strings.Contains(blob, "night shift")
}

func overnightOpen(cat *catalog.Catalog) int {
	n := 0
	for _, s := range cat.List() {
		if s.Stage != catalog.StageProduction && isOvernight(s) {
			n++
		}
	}
	return n
}

func parkOvernight(cat *catalog.Catalog, note string, b *Build) {
	name := "Overnight · " + strings.TrimSpace(note)
	if len(name) > 48 {
		name = name[:48]
	}
	ship, err := cat.Create(catalog.CreateShip{
		Name:   name,
		Client: "Cashtro",
		Sector: "os",
		Stack:  []string{"Go"},
		Notes:  "Closed-hours builder parked this. Agentics keep shipping while the desk is closed.",
	})
	if err != nil {
		b.Steps = append(b.Steps, BuildStep{Agent: "delivery", Capability: "delivery.create", OK: false, Message: err.Error()})
		return
	}
	b.Steps = append(b.Steps, BuildStep{Agent: "delivery", Capability: "delivery.create", OK: true, Message: "parked " + ship.Name})
}

func tenantCount(k *Kernel) (int, int, bool) {
	c, ok := k.Tenants().(tenantCounter)
	if !ok {
		return 0, 0, false
	}
	n, mine := c.Count()
	return n, mine, true
}
