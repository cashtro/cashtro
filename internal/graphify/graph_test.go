package graphify

import (
	"strings"
	"testing"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
)

func sample(query string) Graph {
	return Build(Input{
		Query: query,
		Now:   time.Date(2026, 9, 22, 5, 30, 0, 0, time.UTC),
		Processes: []kernel.Process{
			{Spec: kernel.Spec{ID: "manager", Name: "Manager", Mode: kernel.ModeLive, Role: "epicenter", Kind: kernel.KindSystem, Capabilities: []string{"manager.graphify"}}, Status: kernel.StatusRunning},
			{Spec: kernel.Spec{ID: "explorer", Name: "Explorer", Mode: kernel.ModeLive, Role: "search", Capabilities: []string{"explorer.search"}}, Status: kernel.StatusRunning},
			{Spec: kernel.Spec{ID: "operator", Name: "Operator", Mode: kernel.ModeResident, Role: "computer-use", Capabilities: []string{"operator.browse"}}, Status: kernel.StatusRunning},
		},
		Events: []kernel.Event{
			{Seq: 1, Source: "manager", Kind: "epicenter", Message: "Epicenter online"},
			{Seq: 2, Source: "manager", Kind: "graphify", Message: "graphify query: giant"},
		},
		Ships: []catalog.Ship{
			{ID: "proximity", Name: "Proximity", Client: "Proximity Agency", Stage: catalog.StageProduction},
			{ID: "scanapp", Name: "ScanApp", Client: "Evolu-Jeunes", Stage: catalog.StageConcept},
		},
		Brains: []Brain{
			{ID: "architecte", Name: "L'Architecte", Verb: "pense", Color: "#f0c14b", Subs: []string{"Tisseur"}},
			{ID: "cartographe", Name: "Le Cartographe", Verb: "voit · Graphify", Color: "#6ea8ff", Subs: []string{"Graphifieur"}},
			{ID: "hustler", Name: "Le Hustler", Verb: "empire", Color: "#c084fc", Subs: []string{"Empire"}},
		},
		Lines: []Line{
			{ID: "nft-giant", Name: "NFT + token Giant", Brain: "Hustler", Mandate: "art + token"},
			{ID: "trading", Name: "Crypto et AI bot", Brain: "Architecte", Mandate: "bots"},
			{ID: "proximity", Name: "Proximity Agency", Brain: "Forgeron", Mandate: "usine"},
		},
		Links: []Link{
			{From: "nft-giant", To: "trading", Via: "le token Giant relie l'art NFT aux bots", Agents: []string{"manager", "explorer"}},
		},
	})
}

func TestBuildIncludesGiantAndCodeFlow(t *testing.T) {
	g := sample("")
	if g.Source != "cashtro-graphify" {
		t.Fatalf("source = %q", g.Source)
	}
	if g.Counts["nodes"] < 20 || g.Counts["edges"] < 15 {
		t.Fatalf("counts = %#v", g.Counts)
	}
	if !hasNode(g, "token:giant") {
		t.Fatal("missing GIANT hub")
	}
	if !hasNode(g, "mod:cmd") || !hasNode(g, "mod:graphify") || !hasNode(g, "mod:live") {
		t.Fatal("missing code-flow modules")
	}
	if !hasEdge(g, "mod:cmd", "mod:boot") || !hasEdge(g, "token:giant", "line:nft-giant") {
		t.Fatalf("missing flow/giant edges: %+v", g.Edges)
	}
	if len(g.Flow) < 8 {
		t.Fatalf("flow steps = %d", len(g.Flow))
	}
	if !g.Score.GiantOnline || g.Score.XP < 100 || g.Score.Production != 1 {
		t.Fatalf("score = %+v", g.Score)
	}
	giant := node(g, "token:giant")
	if giant.Size < 40 || giant.Label != "GIANT" {
		t.Fatalf("giant node = %+v", giant)
	}
}

func TestQueryNeighborhoodKeepsGiant(t *testing.T) {
	g := sample("giant")
	if !hasNode(g, "token:giant") {
		t.Fatal("query giant dropped the hub")
	}
	if hasNode(g, "ship:scanapp") && !strings.Contains(strings.ToLower(node(g, "ship:scanapp").Label), "giant") {
		// scanapp is not in the giant neighborhood
		t.Fatalf("query leaked unrelated ship: %#v", ids(g))
	}
	if !hasNode(g, "line:nft-giant") {
		t.Fatalf("giant neighborhood missing nft-giant: %v", ids(g))
	}
}

func hasNode(g Graph, id string) bool {
	return node(g, id).ID == id
}

func node(g Graph, id string) Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return Node{}
}

func hasEdge(g Graph, from, to string) bool {
	for _, e := range g.Edges {
		if e.From == from && e.To == to {
			return true
		}
	}
	return false
}

func ids(g Graph) []string {
	out := make([]string, 0, len(g.Nodes))
	for _, n := range g.Nodes {
		out = append(out, n.ID)
	}
	return out
}
