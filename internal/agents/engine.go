package agents

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// Constellation is one ecosystem with its own infrastructure.
type Constellation struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Infrastructure string   `json:"infrastructure"`
	Function       string   `json:"function"`
	Product        string   `json:"product"`
	Agent          string   `json:"agent"`
	Trinity        []string `json:"trinity"`
}

// GraphNode is one point on the blueprint map.
type GraphNode struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// GraphEdge plugs one point to another.
type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Via  string `json:"via"`
}

// Blueprint is the one map: constellations, links, flows, and a single anchor.
type Blueprint struct {
	Anchor         string          `json:"anchor"`
	Constellations []Constellation `json:"constellations"`
	Links          []Link          `json:"links"`
	Flows          []Flow          `json:"flows"`
	Nodes          []GraphNode     `json:"nodes"`
	Edges          []GraphEdge     `json:"edges"`
	Empty          int             `json:"empty"`
}

func trinityBuild() string {
	parts := make([]string, 0, 6)
	for _, chapter := range TheTrinity().Chapters {
		parts = append(parts, chapter.Title)
	}
	return "Trinity: " + strings.Join(parts, " · ")
}

func infrastructureOf(id string) string {
	base := ""
	switch id {
	case "control":
		base = "kernel, Voltron, Epicenter Einstein, Graphify"
	case "wordpress":
		base = "XAMPP, PHP, ACF Pro"
	case "proximity":
		base = "comptes hors thème"
	case "scanapp":
		base = "scanner interne"
	case "marketplace":
		base = "MCP, Stripe, fiche"
	case "panda":
		base = "offre white-glove"
	case "nft-giant":
		base = "cerveau propre de Giant, art, token, lecture de marché"
	case "ecole":
		base = "programme de cours"
	case "marketing":
		base = "agence, Panda, Proximity cloud, carnet des agents"
	case "fix2":
		base = "site Fix Tout, sa base, école de skills"
	case "empire":
		base = "studio live"
	case "propres":
		base = "Stripe des projets propres"
	case "trading":
		base = "modèle, aucun ordre live"
	case "fonds":
		base = "grand livre, chèque interne, TPS et TVQ"
	}
	if base == "" {
		return ""
	}
	return base + " · " + trinityBuild()
}

// Constellations gives every line its own infrastructure.
func Constellations() []Constellation {
	flows := map[string]Flow{}
	for _, f := range Flows() {
		flows[f.Division] = f
	}
	out := make([]Constellation, 0, len(Lines()))
	for _, ln := range Lines() {
		flow := flows[ln.ID]
		fn := flow.Function
		product := flow.Product
		agent := flow.Agent
		if fn == "" {
			fn = "memory.store"
		}
		if product == "" {
			if len(ln.Repos) > 0 {
				product = ln.Repos[0]
			} else {
				product = ln.Name
			}
		}
		if agent == "" {
			agent = "manager"
		}
		chapters := make([]string, 0, 6)
		for _, chapter := range TheTrinity().Chapters {
			chapters = append(chapters, chapter.ID)
		}
		out = append(out, Constellation{
			ID:             ln.ID,
			Name:           ln.Name,
			Infrastructure: infrastructureOf(ln.ID),
			Function:       fn,
			Product:        product,
			Agent:          agent,
			Trinity:        chapters,
		})
	}
	return out
}

// OneAnchor seals every constellation into a single anchor.
func OneAnchor(items []Constellation) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, item.ID+"|"+item.Infrastructure+"|"+item.Function+"|"+item.Product)
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:])
}

func emptyCount(items []Constellation) int {
	n := 0
	for _, item := range items {
		if item.ID == "" || item.Infrastructure == "" || item.Function == "" || item.Product == "" || item.Agent == "" || len(item.Trinity) != 6 {
			n++
		}
	}
	return n
}

// MakeBlueprint draws the map and plugs every constellation into the one anchor.
func MakeBlueprint(k *kernel.Kernel) Blueprint {
	items := Constellations()
	anchor := OneAnchor(items)
	nodes := []GraphNode{{ID: "anchor", Kind: "anchor", Name: anchor}}
	var edges []GraphEdge
	for _, item := range items {
		nodes = append(nodes, GraphNode{ID: "const:" + item.ID, Kind: "constellation", Name: item.Name})
		nodes = append(nodes, GraphNode{ID: "agent:" + item.Agent, Kind: "agent", Name: item.Agent})
		edges = append(edges, GraphEdge{From: "const:" + item.ID, To: "anchor", Via: item.Infrastructure})
		edges = append(edges, GraphEdge{From: "const:" + item.ID, To: "agent:" + item.Agent, Via: item.Function})
	}
	for _, ln := range Links() {
		edges = append(edges, GraphEdge{From: ln.From, To: ln.To, Via: ln.Via})
	}
	if k != nil {
		k.Remember("anchor", anchor)
	}
	return Blueprint{
		Anchor:         anchor,
		Constellations: items,
		Links:          Links(),
		Flows:          Flows(),
		Nodes:          nodes,
		Edges:          edges,
		Empty:          emptyCount(items),
	}
}

// Lens is the single reading of the engine: one anchor, no empty constellation.
func Lens(k *kernel.Kernel) map[string]any {
	board := MakeBlueprint(k)
	return map[string]any{
		"anchor":         board.Anchor,
		"constellations": len(board.Constellations),
		"empty":          board.Empty,
		"funds":          FundsPosition(k),
		"taxOn100":       TaxQuebec(10000),
		"nodes":          len(board.Nodes),
		"edges":          len(board.Edges),
		"brains":         len(Brains()),
	}
}

// EngineRun gives every constellation its infrastructure, then seals one anchor.
func EngineRun(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	task := payloadQuery(call, "task")
	if task == "" {
		task = "recevoir l'infrastructure"
	}
	items := Constellations()
	if n := emptyCount(items); n > 0 {
		return kernel.Result{OK: false, Message: "constellation vide", Data: map[string]any{"empty": n}}, nil
	}
	pulses := make([]map[string]any, 0, len(items))
	for _, item := range items {
		text := item.Infrastructure + " " + task
		if names := secretNames(text); len(names) > 0 {
			return kernel.Result{OK: false, Message: "interdit: " + strings.Join(names, ", "), Data: map[string]any{
				"constellation": item.ID,
			}}, nil
		}
		k.Remember("constellation:"+item.ID, item.Infrastructure)
		verdict := Contradict(Proposal{
			Subject: item.ID,
			Options: []Option{
				{Name: "infrastructure neuve", Cost: 4, Risk: 2, Steps: 3},
				{Name: "infrastructure déjà là", Cost: 1, Risk: 1, Steps: 1},
			},
		})
		if !verdict.Ready {
			return kernel.Result{OK: false, Message: verdict.Attack, Data: map[string]any{"constellation": item.ID}}, nil
		}
		pulses = append(pulses, map[string]any{
			"id":             item.ID,
			"infrastructure": item.Infrastructure,
			"function":       item.Function,
			"product":        item.Product,
			"kept":           verdict.Best,
		})
	}
	if cents := payloadInt(call, "cents"); cents > 0 {
		tax := TaxQuebec(cents)
		k.Remember("fonds:caisse", strconv.Itoa(tax.Total))
		pulses = append(pulses, map[string]any{"id": "tax", "total": tax.Total})
	}
	report := AnalyzeAll()
	var loose []string
	for _, item := range report {
		if !item.Connected {
			loose = append(loose, item.ID)
		}
	}
	if len(loose) > 0 {
		return kernel.Result{OK: false, Message: "analyse débranchée: " + strings.Join(loose, ", "), Data: map[string]any{
			"loose": loose,
		}}, nil
	}
	board := MakeBlueprint(k)
	k.Publish("manager", "engine", "constellations branchées", map[string]any{
		"anchor":   board.Anchor,
		"count":    len(items),
		"analyzed": len(report),
	})
	return kernel.Result{OK: true, Message: "moteur en marche. une ancre", Data: map[string]any{
		"anchor":   board.Anchor,
		"empty":    board.Empty,
		"pulses":   pulses,
		"brains":   len(Brains()),
		"analyzed": len(report),
		"lens":     Lens(k),
	}}, nil
}

// SaveBlueprint writes the map next to the operating board.
func SaveBlueprint(k *kernel.Kernel) (string, error) {
	root, err := findFile("go.mod")
	if err != nil {
		return "", err
	}
	path := filepath.Join(filepath.Dir(root), "state", "blueprint.json")
	raw, err := json.MarshalIndent(MakeBlueprint(k), "", "  ")
	if err != nil {
		return "", err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
