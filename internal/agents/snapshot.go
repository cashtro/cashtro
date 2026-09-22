package agents

import (
	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/graphify"
	"github.com/cashtro/cashtro/internal/kernel"
)

// Map builds the Graphify snapshot without writing mail or journal.
func Map(k *kernel.Kernel, query string) graphify.Graph {
	a := &managerAgent{k: k}
	return a.graph(query)
}

// Snapshot is the Giant HUD tick: one read, no side effects.
func Snapshot(k *kernel.Kernel, query string) map[string]any {
	cat := k.Catalog()
	var ships []catalog.Ship
	var profile catalog.Profile
	if cat != nil {
		ships = cat.List()
		profile = cat.Profile()
	}
	return map[string]any{
		"os":       k.About(),
		"agents":   k.Processes(),
		"ships":    ships,
		"events":   k.Events(),
		"notes":    k.Notes(),
		"confirms": k.Confirms(),
		"mail":     k.Inbox(""),
		"memory":   k.Recall(""),
		"graph":    Map(k, query),
		"brains":   Brains(),
		"lines":    Lines(),
		"profile":  profile,
	}
}
