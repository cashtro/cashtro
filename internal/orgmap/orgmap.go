// Package orgmap is the joint map of Cashtro (control plane) and
// Evolu-Jeunes (delivery org). Catalog ships are named only. Live
// client production stays do-not-touch. Private GitHub trees stay
// dark until Option A.
package orgmap

import (
	"encoding/json"
	"html"
	"strings"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
)

const (
	// Title is the public name of the joint map.
	Title = "Cashtro × Evolu-Jeunes"
	// CashtroRepo is the reachable control-plane repository.
	CashtroRepo = "cashtro/cashtro"
	// EvoluOrg is the delivery GitHub organization.
	EvoluOrg = "Evolu-Jeunes"
)

// Node is one place on the joint map.
type Node struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Org    string `json:"org,omitempty"`
	Access string `json:"access"`
	Stage  string `json:"stage,omitempty"`
	URL    string `json:"url,omitempty"`
	Notes  string `json:"notes,omitempty"`
}

// Edge is a labeled relation between two nodes.
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rel  string `json:"rel"`
}

// Graph is the Cashtro + Evolu-Jeunes fleet mapped together.
type Graph struct {
	Title       string    `json:"title"`
	GeneratedAt time.Time `json:"generatedAt"`
	Summary     string    `json:"summary"`
	Links       []Link    `json:"links"`
	Nodes       []Node    `json:"nodes"`
	Edges       []Edge    `json:"edges"`
}

// Link is a shareable URL on the map.
type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Joint returns the canonical two-org map. It does not invent private
// repo slugs. Evolu-Jeunes ships come from the kernel catalog.
func Joint() Graph {
	ships := catalog.New().List()
	nodes := []Node{
		{ID: "castro", Name: "Castro", Kind: "person", Access: "ok", Notes: "Scrum Master & Software Engineer · Founder"},
		{ID: "cashtro-user", Name: "cashtro", Kind: "github-user", Org: "cashtro", Access: "ok", URL: "https://github.com/cashtro", Notes: "Public GitHub user. 1 public repo."},
		{ID: "cashtro-repo", Name: CashtroRepo, Kind: "control-plane", Org: "cashtro", Access: "ok", Stage: "active", URL: "https://github.com/cashtro/cashtro", Notes: "Cashtro OS / Voltron. Go kernel + TypeScript control plane."},
		{ID: "kernel", Name: "Cashtro OS kernel", Kind: "runtime", Org: "cashtro", Access: "ok", URL: "http://127.0.0.1:8080", Notes: "14 agentics on :8080. 7 live, 7 resident."},
		{ID: "control-plane", Name: "Control plane API", Kind: "runtime", Org: "cashtro", Access: "ok", URL: "http://127.0.0.1:8787/docs", Notes: "Recon, Prisma registry, Fastify. $2/run cap."},
		{ID: "evolu-org", Name: EvoluOrg, Kind: "github-org", Org: EvoluOrg, Access: "limited", URL: "https://github.com/Evolu-Jeunes", Notes: "Delivery org. 0 public repos. ~48 private client trees expected."},
		{ID: "evolu-member", Name: "evoluJeunes", Kind: "github-user", Org: EvoluOrg, Access: "ok", URL: "https://github.com/evoluJeunes", Notes: "Public member of the Evolu-Jeunes org."},
		{ID: "evolu-site", Name: "evolujeunes.ca", Kind: "production-site", Org: EvoluOrg, Access: "live-client", Stage: "production", URL: "https://evolujeunes.ca", Notes: "LIVE nonprofit site. Do not touch production."},
		{ID: "evolu-dark", Name: "~48 private repos", Kind: "org-wall", Org: EvoluOrg, Access: "limited", Notes: "Cursor GitHub token cannot list them. See docs/ACCESS_REQUIRED.md."},
	}
	for _, s := range ships {
		access := "catalog-only"
		notes := s.Notes
		if s.Stage == catalog.StageProduction {
			access = "live-client"
			notes = strings.TrimSpace(s.Notes + " LIVE CLIENT — do not touch.")
		}
		nodes = append(nodes, Node{
			ID:     "ship-" + s.ID,
			Name:   s.Name,
			Kind:   "catalog-ship",
			Org:    EvoluOrg,
			Access: access,
			Stage:  string(s.Stage),
			Notes:  notes,
		})
	}
	edges := []Edge{
		{From: "castro", To: "cashtro-user", Rel: "owns"},
		{From: "castro", To: "evolu-org", Rel: "owns"},
		{From: "cashtro-user", To: "cashtro-repo", Rel: "contains"},
		{From: "cashtro-repo", To: "kernel", Rel: "runs"},
		{From: "cashtro-repo", To: "control-plane", Rel: "runs"},
		{From: "control-plane", To: "evolu-org", Rel: "recon-scans"},
		{From: "evolu-org", To: "evolu-member", Rel: "has-member"},
		{From: "evolu-org", To: "evolu-site", Rel: "publishes"},
		{From: "evolu-org", To: "evolu-dark", Rel: "holds"},
	}
	for _, s := range ships {
		edges = append(edges, Edge{From: "kernel", To: "ship-" + s.ID, Rel: "catalog-seeds"})
		edges = append(edges, Edge{From: "evolu-org", To: "ship-" + s.ID, Rel: "expected-repo"})
	}
	return Graph{
		Title:       Title,
		GeneratedAt: time.Date(2026, 9, 22, 2, 0, 0, 0, time.UTC),
		Summary:     "Cashtro is the public control plane. Evolu-Jeunes is the private delivery org. They are one fleet under Castro. Client production stays dark.",
		Links: []Link{
			{Label: "cashtro/cashtro", URL: "https://github.com/cashtro/cashtro"},
			{Label: "Evolu-Jeunes org", URL: "https://github.com/Evolu-Jeunes"},
			{Label: "evoluJeunes member", URL: "https://github.com/evoluJeunes"},
			{Label: "evolujeunes.ca", URL: "https://evolujeunes.ca"},
			{Label: "Desk map", URL: "http://127.0.0.1:8080/map"},
			{Label: "Markdown map", URL: "https://github.com/cashtro/cashtro/blob/main/docs/FLEET_MAP.md"},
		},
		Nodes: nodes,
		Edges: edges,
	}
}

// JSON returns the joint graph as JSON.
func JSON() []byte {
	raw, err := json.MarshalIndent(Joint(), "", "  ")
	if err != nil {
		return []byte(`{"error":"marshal"}`)
	}
	return append(raw, '\n')
}

// Mermaid renders a flowchart that keeps both orgs on one canvas.
func Mermaid() string {
	g := Joint()
	var b strings.Builder
	b.WriteString("flowchart TB\n")
	b.WriteString("  castro[Castro]\n")
	b.WriteString("  subgraph cashtro_home[cashtro GitHub user]\n")
	b.WriteString("    cashtro_repo[cashtro/cashtro — control plane]\n")
	b.WriteString("    kernel[Cashtro OS kernel — 14 agentics]\n")
	b.WriteString("    cp[TypeScript recon + registry + API]\n")
	b.WriteString("    cashtro_repo --> kernel\n")
	b.WriteString("    cashtro_repo --> cp\n")
	b.WriteString("  end\n")
	b.WriteString("  subgraph evolu_home[Evolu-Jeunes GitHub org]\n")
	b.WriteString("    evolu_member[github.com/evoluJeunes]\n")
	b.WriteString("    evolu_site[evolujeunes.ca — LIVE do not touch]\n")
	b.WriteString("    evolu_dark[~48 private repos / access limited]\n")
	b.WriteString("    subgraph live[LIVE CLIENT — catalog named only]\n")
	for _, n := range g.Nodes {
		if n.Kind == "catalog-ship" && n.Access == "live-client" {
			b.WriteString("      " + mermaidID(n.ID) + "[" + mermaidLabel(n.Name+" / "+n.Stage) + "]\n")
		}
	}
	b.WriteString("    end\n")
	b.WriteString("    subgraph later[Safe to onboard later]\n")
	for _, n := range g.Nodes {
		if n.Kind == "catalog-ship" && n.Access != "live-client" {
			b.WriteString("      " + mermaidID(n.ID) + "[" + mermaidLabel(n.Name+" / "+n.Stage) + "]\n")
		}
	}
	b.WriteString("    end\n")
	b.WriteString("  end\n")
	b.WriteString("  castro -->|owns| cashtro_home\n")
	b.WriteString("  castro -->|owns| evolu_home\n")
	b.WriteString("  cp -->|recon scans| evolu_home\n")
	b.WriteString("  kernel -->|catalog seeds| live\n")
	b.WriteString("  kernel -->|catalog seeds| later\n")
	b.WriteString("  cp -.->|Option A PAT required| evolu_dark\n")
	return b.String()
}

// Markdown is the GitHub-rendered joint map.
func Markdown() string {
	g := Joint()
	var b strings.Builder
	b.WriteString("# " + g.Title + "\n\n")
	b.WriteString(g.Summary + "\n\n")
	b.WriteString("Open on the desk: [`http://127.0.0.1:8080/map`](http://127.0.0.1:8080/map) · JSON: [`/api/map`](http://127.0.0.1:8080/api/map).\n\n")
	b.WriteString("| Place | URL |\n| --- | --- |\n")
	for _, l := range g.Links {
		b.WriteString("| " + l.Label + " | " + l.URL + " |\n")
	}
	b.WriteString("\n```mermaid\n")
	b.WriteString(Mermaid())
	b.WriteString("```\n\n")
	b.WriteString("## What is reachable today\n\n")
	b.WriteString("- **cashtro/cashtro** — public control plane. Go kernel, Voltron, recon, registry, API.\n")
	b.WriteString("- **Evolu-Jeunes** — org exists, **0 public repos**. Token cannot list the private fleet.\n")
	b.WriteString("- **evoluJeunes** — public org member.\n")
	b.WriteString("- **evolujeunes.ca** — live nonprofit site. Do not touch.\n\n")
	b.WriteString("## Catalog ships (not GitHub-scanned)\n\n")
	b.WriteString("These names come from `internal/catalog`. Recon did not open their repos.\n\n")
	for _, n := range g.Nodes {
		if n.Kind != "catalog-ship" {
			continue
		}
		flag := "safe to onboard later"
		if n.Access == "live-client" {
			flag = "**LIVE CLIENT — do not touch**"
		}
		b.WriteString("- **" + n.Name + "** · " + n.Stage + " · " + flag + "\n")
	}
	b.WriteString("\n## Hard stop\n\n")
	b.WriteString("Do not guess private repo slugs. After Option A (`docs/ACCESS_REQUIRED.md`), re-run `pnpm recon` and this map gains real Evolu-Jeunes trees.\n")
	return b.String()
}

func mermaidID(id string) string {
	id = strings.ReplaceAll(id, "-", "_")
	return id
}

func mermaidLabel(s string) string {
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")
	s = strings.ReplaceAll(s, "\"", "'")
	return s
}

func esc(s string) string { return html.EscapeString(s) }
