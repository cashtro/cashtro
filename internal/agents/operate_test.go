package agents

import (
	"os"
	"path/filepath"
	"testing"
)

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
	if len(board.Lines) != 9 || len(board.Agents) != 15 || len(board.Links) != 11 {
		t.Fatalf("lines %d agents %d links %d", len(board.Lines), len(board.Agents), len(board.Links))
	}
	if len(board.Gaps) != 3 {
		t.Fatalf("gaps = %v, want the 3 unclassified repos", board.Gaps)
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
	if proximity < 30 {
		t.Fatalf("proximity repos = %d", proximity)
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
