package orgmap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJointMapsBothHomes(t *testing.T) {
	g := Joint()
	if g.Title != Title {
		t.Fatalf("title = %q", g.Title)
	}
	ids := map[string]Node{}
	for _, n := range g.Nodes {
		ids[n.ID] = n
	}
	for _, id := range []string{"castro", "cashtro-repo", "evolu-org", "evolu-dark", "ship-proximity", "ship-scanapp"} {
		if _, ok := ids[id]; !ok {
			t.Fatalf("missing node %s", id)
		}
	}
	if ids["cashtro-repo"].URL != "https://github.com/cashtro/cashtro" {
		t.Fatalf("cashtro url = %s", ids["cashtro-repo"].URL)
	}
	if ids["evolu-org"].Access != "limited" {
		t.Fatalf("evolu access = %s", ids["evolu-org"].Access)
	}
	if ids["ship-proximity"].Access != "live-client" {
		t.Fatalf("proximity access = %s", ids["ship-proximity"].Access)
	}

	var ownsCashtro, ownsEvolu, recon bool
	for _, e := range g.Edges {
		if e.From == "castro" && e.To == "cashtro-user" && e.Rel == "owns" {
			ownsCashtro = true
		}
		if e.From == "castro" && e.To == "evolu-org" && e.Rel == "owns" {
			ownsEvolu = true
		}
		if e.From == "control-plane" && e.To == "evolu-org" && e.Rel == "recon-scans" {
			recon = true
		}
	}
	if !ownsCashtro || !ownsEvolu || !recon {
		t.Fatalf("missing joint edges ownsCashtro=%v ownsEvolu=%v recon=%v", ownsCashtro, ownsEvolu, recon)
	}
}

func TestMermaidConnectsOrgs(t *testing.T) {
	m := Mermaid()
	for _, needle := range []string{
		"cashtro/cashtro",
		"Evolu-Jeunes",
		"castro -->|owns| cashtro_home",
		"castro -->|owns| evolu_home",
		"recon scans",
		"BTK Avocats",
		"ScanApp",
		"LIVE CLIENT",
	} {
		if !strings.Contains(m, needle) {
			t.Fatalf("mermaid missing %q\n%s", needle, m)
		}
	}
}

func TestJSONAndPage(t *testing.T) {
	var g Graph
	if err := json.Unmarshal(JSON(), &g); err != nil {
		t.Fatal(err)
	}
	if len(g.Nodes) < 10 || len(g.Edges) < 10 {
		t.Fatalf("graph too small nodes=%d edges=%d", len(g.Nodes), len(g.Edges))
	}
	page := string(PageHTML())
	for _, needle := range []string{
		"Cashtro × Evolu-Jeunes",
		"id=\"cashtro\"",
		"id=\"evolu-jeunes\"",
		"https://github.com/cashtro/cashtro",
		"https://github.com/Evolu-Jeunes",
		"evolujeunes.ca",
		"do not touch",
		"ScanApp",
		"Proximity",
	} {
		if !strings.Contains(page, needle) {
			t.Fatalf("page missing %q", needle)
		}
	}
	md := Markdown()
	if !strings.Contains(md, "```mermaid") || !strings.Contains(md, CashtroRepo) {
		t.Fatalf("markdown incomplete")
	}
	root := repoRoot(t)
	got, err := os.ReadFile(filepath.Join(root, "docs/FLEET_MAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != md {
		t.Fatalf("docs/FLEET_MAP.md drifted from orgmap.Markdown(); regenerate it")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
