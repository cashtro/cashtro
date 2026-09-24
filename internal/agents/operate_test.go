package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWordPressLeavesTheOtherDepartment(t *testing.T) {
	board, err := LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if len(board.Lines) != 14 || len(board.Links) != 23 || len(board.Org.Departments) != 9 {
		t.Fatalf("lines %d links %d depts %d", len(board.Lines), len(board.Links), len(board.Org.Departments))
	}
	if len(board.Brains) != len(Brains()) || len(board.Brains) != 15 {
		t.Fatalf("brains = %d", len(board.Brains))
	}
	if gap := Disconnected(); len(gap) != 0 {
		t.Fatalf("board not analyzing: %v", gap)
	}
	var wordpress, proximity int
	wp := map[string]bool{}
	for _, ln := range board.Lines {
		if ln.ID == "wordpress" {
			wordpress = len(ln.Repos)
			for _, r := range ln.Repos {
				wp[r] = true
			}
		}
		if ln.ID == "proximity" {
			proximity = len(ln.Repos)
			for _, r := range ln.Repos {
				if wp[r] {
					t.Fatalf("wordpress repo still on the other team: %s", r)
				}
			}
		}
	}
	if wordpress < 20 || proximity == 0 || proximity > wordpress {
		t.Fatalf("wordpress repos = %d, other proximity repos = %d", wordpress, proximity)
	}
	for _, repo := range []string{"Evolu-Jeunes/Proximity", "Evolu-Jeunes/ProximityApp", "Evolu-Jeunes/Api-Proximity"} {
		if !wp[repo] {
			t.Fatalf("%s is not on the wordpress commit list", repo)
		}
	}
	for _, c := range board.Repos {
		if wp[c.FullName] && c.LineID != "wordpress" {
			t.Fatalf("%s filed as %s", c.FullName, c.LineID)
		}
	}
	root, err := findFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(filepath.Dir(root), "state", "operating.json")
	if err := SaveBoard(out, board); err != nil {
		t.Fatal(err)
	}
}

func TestBuildBoardSyncsFiches(t *testing.T) {
	dir := os.Getenv("EPICENTER_PROJECTS")
	if dir == "" {
		dir = "/tmp/epicenter/PROJECTS"
	}
	if _, err := os.Stat(dir); err != nil {
		t.Skip("fiches not on disk")
	}
	graph := os.Getenv("GRAPH_ACCOUNTS")
	if graph == "" {
		graph = "/tmp/epicenter/graph/accounts.json"
	}
	board, err := BuildBoard(dir, graph)
	if err != nil {
		t.Fatal(err)
	}
	if board.Fiches != 58 {
		t.Fatalf("fiches = %d, want 58", board.Fiches)
	}
	if len(board.Lines) != 14 || len(board.Agents) != 16 || len(board.Links) != 23 {
		t.Fatalf("lines %d agents %d links %d", len(board.Lines), len(board.Agents), len(board.Links))
	}
	if len(board.Gaps) != 0 {
		t.Fatalf("gaps = %v, want none", board.Gaps)
	}
	if UsesStripe("proximity", false) || !UsesStripe("proximity", true) || UsesStripe("scanapp", true) {
		t.Fatal("stripe rail: proximity only when asked, scanner never")
	}
	for _, id := range []string{"marketplace", "propres", "marketing", "empire", "panda"} {
		if !UsesStripe(id, false) {
			t.Fatalf("stripe should encash %s", id)
		}
	}
	if UsesStripe("trading", false) || UsesStripe("nft-giant", false) {
		t.Fatal("crypto rails stay off the card account")
	}
	if board.Graph.RepoNodes < 50 || board.Graph.Bridges < 1 {
		t.Fatalf("graph = %+v", board.Graph)
	}
	var proximity int
	for _, ln := range board.Lines {
		if ln.ID == "proximity" {
			proximity = len(ln.Repos)
		}
	}
	var wordpress int
	for _, ln := range board.Lines {
		if ln.ID == "wordpress" {
			wordpress = len(ln.Repos)
		}
	}
	if wordpress < 20 || proximity > wordpress {
		t.Fatalf("wordpress repos = %d, other proximity repos = %d", wordpress, proximity)
	}
	root, err := findFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(filepath.Dir(root), "state", "operating.json")
	if err := SaveBoard(out, board); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Fiches != 58 || loaded.Graph.Bridges != board.Graph.Bridges {
		t.Fatalf("loaded = fiches %d bridges %d", loaded.Fiches, loaded.Graph.Bridges)
	}
}
