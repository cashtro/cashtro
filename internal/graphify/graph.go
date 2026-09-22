// Package graphify is the live Cashtro knowledge graph.
//
// Graphify MCP is the external map. This package is the kernel-side
// graph the desk can render every second: how the code unfolds, the
// five brains, the business ecosystems, and Giant as the world hub.
package graphify

import (
	"sort"
	"strings"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
)

// Node is one entity on the Graphify map.
type Node struct {
	ID     string         `json:"id"`
	Kind   string         `json:"kind"`
	Label  string         `json:"label"`
	Detail string         `json:"detail"`
	Group  string         `json:"group"`
	Color  string         `json:"color"`
	Size   int            `json:"size"`
	Heat   float64        `json:"heat"`
	Meta   map[string]any `json:"meta,omitempty"`
}

// Edge is one relationship.
type Edge struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Kind  string `json:"kind"`
	Label string `json:"label"`
}

// FlowStep is one beat in how the code actually runs.
type FlowStep struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	File  string `json:"file"`
	Verb  string `json:"verb"`
}

// Score is the game HUD meter derived from live kernel state.
type Score struct {
	XP          int    `json:"xp"`
	Level       int    `json:"level"`
	Title       string `json:"title"`
	Next        int    `json:"next"`
	Production  int    `json:"production"`
	Concept     int    `json:"concept"`
	Idea        int    `json:"idea"`
	Live        int    `json:"live"`
	Resident    int    `json:"resident"`
	Quests      int    `json:"quests"`
	Events      int    `json:"events"`
	GiantOnline bool   `json:"giantOnline"`
}

// Graph is a Graphify-shaped snapshot.
type Graph struct {
	Query     string         `json:"query"`
	Source    string         `json:"source"`
	Nodes     []Node         `json:"nodes"`
	Edges     []Edge         `json:"edges"`
	Flow      []FlowStep     `json:"flow"`
	Score     Score          `json:"score"`
	Counts    map[string]int `json:"counts"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// Brain is a five-mind card passed in from the agents roster.
type Brain struct {
	ID    string
	Name  string
	Verb  string
	Color string
	Subs  []string
}

// Line is a business line passed in from the Manager.
type Line struct {
	ID       string
	Name     string
	Brain    string
	Mandate  string
	Repos    []string
	Controls string
}

// Link is one ecosystem automation.
type Link struct {
	From   string
	To     string
	Via    string
	Agents []string
}

// Input is everything the map needs from the running OS.
type Input struct {
	Query     string
	Now       time.Time
	Processes []kernel.Process
	Events    []kernel.Event
	Ships     []catalog.Ship
	Brains    []Brain
	Lines     []Line
	Links     []Link
}

const (
	giantID     = "token:giant"
	osID        = "os:cashtro"
	graphifyID  = "mod:graphify"
	sourceName  = "cashtro-graphify"
	levelStride = 250
)

// Build assembles the live map. Query filters to a neighborhood when set.
func Build(in Input) Graph {
	if in.Now.IsZero() {
		in.Now = time.Now().UTC()
	}
	nodes := map[string]Node{}
	var edges []Edge
	heat := heatFrom(in)

	put := func(n Node) {
		if n.Size == 0 {
			n.Size = 16
		}
		if h, ok := heat[n.ID]; ok && h > n.Heat {
			n.Heat = h
		}
		nodes[n.ID] = n
	}
	link := func(from, to, kind, label string) {
		if from == "" || to == "" || from == to {
			return
		}
		edges = append(edges, Edge{From: from, To: to, Kind: kind, Label: label})
	}

	put(Node{
		ID: osID, Kind: "os", Label: "Cashtro OS", Group: "os", Color: "#ff4d3d", Size: 28,
		Detail: "Kernel Voltron · idea → concept → production",
		Meta:   map[string]any{"version": kernel.Version},
	})
	put(Node{
		ID: giantID, Kind: "token", Label: "GIANT", Group: "giant", Color: "#f0c14b", Size: 56,
		Detail: "Token de l'écosystème · pont NFT, trading, live Empire. Le boss de la map.",
		Meta:   map[string]any{"repos": []string{"Evolu-Jeunes/Giant", "Evolu-Jeunes/Nft"}, "role": "world-hub"},
	})

	for _, b := range in.Brains {
		bid := "brain:" + b.ID
		put(Node{
			ID: bid, Kind: "brain", Label: b.Name, Group: "brain", Color: b.Color, Size: 26,
			Detail: b.Verb, Meta: map[string]any{"subs": b.Subs},
		})
		link(osID, bid, "steers", "cerveau")
	}

	for _, ln := range in.Lines {
		lid := "line:" + ln.ID
		color := "#8b93a7"
		size := 20
		if ln.ID == "nft-giant" || ln.ID == "trading" {
			color = "#f0c14b"
			size = 28
		}
		put(Node{
			ID: lid, Kind: "line", Label: ln.Name, Group: "eco", Color: color, Size: size,
			Detail: ln.Mandate, Meta: map[string]any{"repos": ln.Repos, "controls": ln.Controls, "brain": ln.Brain},
		})
		if b, ok := brainID(in.Brains, ln.Brain); ok {
			link("brain:"+b, lid, "owns", ln.Brain)
		}
		link(osID, lid, "line", "business")
		if ln.ID == "nft-giant" || ln.ID == "trading" || ln.ID == "empire" {
			link(giantID, lid, "token", "Giant")
		}
	}

	for _, lk := range in.Links {
		from := "line:" + lk.From
		to := "line:" + lk.To
		link(from, to, "ecosystem", lk.Via)
		for _, agent := range lk.Agents {
			link("agent:"+agent, from, "carries", lk.Via)
			link("agent:"+agent, to, "carries", lk.Via)
		}
	}

	for _, p := range in.Processes {
		aid := "agent:" + p.Spec.ID
		color := "#8b93a7"
		if p.Spec.Mode == kernel.ModeLive {
			color = "#3dffa6"
		}
		put(Node{
			ID: aid, Kind: "agent", Label: p.Spec.Name, Group: "crew", Color: color, Size: 18,
			Detail: p.Spec.Role + " · " + string(p.Spec.Mode) + " · " + string(p.Status),
			Meta: map[string]any{
				"role": p.Spec.Role, "mode": p.Spec.Mode, "status": p.Status,
				"capabilities": p.Spec.Capabilities,
			},
		})
		link(osID, aid, "process", string(p.Spec.Kind))
	}

	for _, s := range in.Ships {
		sid := "ship:" + s.ID
		color := "#6ea8ff"
		switch s.Stage {
		case catalog.StageConcept:
			color = "#f0c14b"
		case catalog.StageProduction:
			color = "#3dffa6"
		}
		put(Node{
			ID: sid, Kind: "ship", Label: s.Name, Group: "quest", Color: color, Size: 14,
			Detail: string(s.Stage) + " · " + s.Client, Meta: map[string]any{"stage": s.Stage, "sector": s.Sector},
		})
		link(osID, sid, "quest", string(s.Stage))
		if s.ID == "proximity" {
			link("line:proximity", sid, "ships", "usine")
		}
		if s.ID == "scanapp" {
			link("line:scanapp", sid, "ships", "marketplace")
		}
		if s.ID == "educonnexion" {
			link("line:ecole", sid, "ships", "école")
		}
	}

	for _, step := range codeFlow() {
		put(Node{
			ID: step.ID, Kind: "module", Label: step.Label, Group: "code", Color: "#d7dde8", Size: 15,
			Detail: step.File + " · " + step.Verb, Meta: map[string]any{"file": step.File, "verb": step.Verb},
		})
	}
	flowEdges := [][3]string{
		{"mod:cmd", "mod:boot", "appelle"},
		{"mod:boot", "mod:register", "charge"},
		{"mod:register", "mod:kernel-boot", "spawn"},
		{"mod:kernel-boot", "mod:server", "écoute"},
		{"mod:server", "mod:desk", "GET /"},
		{"mod:desk", "mod:invoke", "POST invoke"},
		{"mod:invoke", graphifyID, "manager.graphify"},
		{graphifyID, "mod:live", "GET /api/live"},
		{"mod:live", "mod:desk", "HUD tick"},
		{"mod:invoke", "mod:catalog", "delivery"},
		{"mod:invoke", "mod:model", "router"},
	}
	for _, e := range flowEdges {
		link(e[0], e[1], "flow", e[2])
	}
	link(osID, "mod:cmd", "runs", "entry")
	link(graphifyID, giantID, "maps", "Giant hub")
	link("brain:cartographe", graphifyID, "voit", "Graphify")
	link("agent:manager", graphifyID, "queries", "manager.graphify")
	link("agent:explorer", graphifyID, "searches", "explorer")

	g := Graph{
		Query:     strings.TrimSpace(in.Query),
		Source:    sourceName,
		Flow:      codeFlow(),
		Score:     score(in),
		UpdatedAt: in.Now,
	}
	if g.Query != "" {
		nodes, edges = neighborhood(nodes, edges, g.Query)
	}
	g.Nodes = values(nodes)
	g.Edges = uniqueEdges(edges)
	g.Counts = counts(g.Nodes, g.Edges)
	return g
}

func codeFlow() []FlowStep {
	return []FlowStep{
		{ID: "mod:cmd", Label: "cmd/cashtro", File: "cmd/cashtro/main.go", Verb: "boot process"},
		{ID: "mod:boot", Label: "agents.Boot", File: "internal/agents/builtin.go", Verb: "register 15 seats"},
		{ID: "mod:register", Label: "kernel.Register", File: "internal/kernel/kernel.go", Verb: "process table"},
		{ID: "mod:kernel-boot", Label: "kernel.Boot", File: "internal/kernel/kernel.go", Verb: "spawn autostart"},
		{ID: "mod:server", Label: "server.New", File: "internal/server/server.go", Verb: "HTTP mux"},
		{ID: "mod:desk", Label: "Giant HUD", File: "internal/server/index.html", Verb: "render live desk"},
		{ID: "mod:invoke", Label: "kernel.Invoke", File: "internal/kernel/kernel.go", Verb: "route capability"},
		{ID: graphifyID, Label: "graphify.Build", File: "internal/graphify/graph.go", Verb: "live map"},
		{ID: "mod:live", Label: "/api/live", File: "internal/server/server.go", Verb: "HUD snapshot"},
		{ID: "mod:catalog", Label: "catalog.Advance", File: "internal/catalog/catalog.go", Verb: "idea → production"},
		{ID: "mod:model", Label: "model.Bus", File: "internal/model/bus.go", Verb: "both lanes"},
	}
}

func score(in Input) Score {
	var live, resident, idea, concept, prod int
	for _, p := range in.Processes {
		if p.Spec.Mode == kernel.ModeLive {
			live++
		} else {
			resident++
		}
	}
	for _, s := range in.Ships {
		switch s.Stage {
		case catalog.StageIdea:
			idea++
		case catalog.StageConcept:
			concept++
		case catalog.StageProduction:
			prod++
		}
	}
	xp := len(in.Events)*3 + prod*120 + concept*40 + idea*10 + live*15 + resident*6
	if xp < 1 {
		xp = 1
	}
	level := 1 + xp/levelStride
	title := "Cadet"
	switch {
	case level >= 8:
		title = "Giant"
	case level >= 6:
		title = "Architecte"
	case level >= 4:
		title = "Forgeron"
	case level >= 2:
		title = "Operator"
	}
	return Score{
		XP: xp, Level: level, Title: title, Next: level * levelStride,
		Production: prod, Concept: concept, Idea: idea,
		Live: live, Resident: resident, Quests: len(in.Ships),
		Events: len(in.Events), GiantOnline: true,
	}
}

func heatFrom(in Input) map[string]float64 {
	out := map[string]float64{}
	if len(in.Events) == 0 {
		return out
	}
	last := in.Events[len(in.Events)-1].Seq
	for _, ev := range in.Events {
		age := float64(last - ev.Seq)
		h := 1.0
		if age > 0 {
			h = 1.0 / (1.0 + age/8.0)
		}
		bump := func(id string) {
			if h > out[id] {
				out[id] = h
			}
		}
		bump("agent:" + ev.Source)
		if strings.Contains(strings.ToLower(ev.Message+ev.Kind), "giant") {
			bump(giantID)
		}
		if strings.Contains(strings.ToLower(ev.Kind+ev.Message), "graphify") {
			bump(graphifyID)
			bump("brain:cartographe")
		}
		if ev.Kind == "ecosystem" || ev.Kind == "mail" {
			bump(osID)
		}
	}
	for _, p := range in.Processes {
		if p.Spec.Mode == kernel.ModeLive && p.Status == kernel.StatusRunning {
			id := "agent:" + p.Spec.ID
			if out[id] < 0.35 {
				out[id] = 0.35
			}
		}
	}
	out[giantID] = max(out[giantID], 0.7)
	return out
}

func neighborhood(nodes map[string]Node, edges []Edge, q string) (map[string]Node, []Edge) {
	q = strings.ToLower(strings.TrimSpace(q))
	seed := map[string]bool{}
	for id, n := range nodes {
		blob := strings.ToLower(n.ID + " " + n.Kind + " " + n.Label + " " + n.Detail + " " + n.Group)
		if strings.Contains(blob, q) {
			seed[id] = true
		}
	}
	if len(seed) == 0 {
		return nodes, edges
	}
	keep := map[string]bool{}
	for id := range seed {
		keep[id] = true
	}
	for _, e := range edges {
		if seed[e.From] || seed[e.To] {
			keep[e.From] = true
			keep[e.To] = true
		}
	}
	filtered := map[string]Node{}
	for id, n := range nodes {
		if keep[id] {
			filtered[id] = n
		}
	}
	var next []Edge
	for _, e := range edges {
		if keep[e.From] && keep[e.To] {
			next = append(next, e)
		}
	}
	return filtered, next
}

func values(nodes map[string]Node) []Node {
	out := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func uniqueEdges(edges []Edge) []Edge {
	seen := map[string]bool{}
	out := make([]Edge, 0, len(edges))
	for _, e := range edges {
		key := e.From + "|" + e.To + "|" + e.Kind
		if seen[key] {
			continue
		}
		if _, ok := seen[e.To+"|"+e.From+"|"+e.Kind]; ok && e.Kind != "flow" {
			continue
		}
		seen[key] = true
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].From != out[j].From {
			return out[i].From < out[j].From
		}
		return out[i].To < out[j].To
	})
	return out
}

func counts(nodes []Node, edges []Edge) map[string]int {
	out := map[string]int{"nodes": len(nodes), "edges": len(edges)}
	for _, n := range nodes {
		out[n.Kind]++
	}
	return out
}

func brainID(brains []Brain, name string) (string, bool) {
	for _, b := range brains {
		if b.Name == name || b.ID == name {
			return b.ID, true
		}
		n := strings.ToLower(name)
		if strings.Contains(strings.ToLower(b.Name), n) || strings.Contains(n, b.ID) {
			return b.ID, true
		}
	}
	return "", false
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
