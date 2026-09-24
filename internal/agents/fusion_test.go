package agents

import "testing"

func TestTheFusionCoversBothPeoples(t *testing.T) {
	force := Fusion()
	if force.Name != "The Fusion" || !force.BothAlways {
		t.Fatalf("name %+v", force.Name)
	}
	if len(force.Peoples) != 2 || force.Peoples[0] != "Evolian" || force.Peoples[1] != "Astro" {
		t.Fatalf("peoples %+v", force.Peoples)
	}
	if len(force.Accounts) != 2 || force.Accounts[0] != "Evolu-Jeunes" || force.Accounts[1] != "cashtro" {
		t.Fatalf("accounts %+v", force.Accounts)
	}
	if force.Agents < 80 || len(force.Workers) != force.Agents || len(force.Workers) != 80 {
		t.Fatalf("agents %d workers %d", force.Agents, len(force.Workers))
	}
	if force.Parent != "eye" || force.Repo != "cashtro/cashtro" || force.Brain != "fusion" {
		t.Fatalf("seat %+v", force)
	}
	if force.ImportsOther || force.MarkdownEdges || force.PathToChains || !force.Parallel {
		t.Fatal("The Fusion must not invent an import into a TypeScript brain")
	}
	if force.Weapons || force.Attacks || force.Illegal || force.CopiesSecrets || force.Surveillance || force.Impersonates || force.Sent || force.Deployed {
		t.Fatal("The Fusion must not attack, watch a person, or impersonate a state")
	}
	if force.OpenedBy != "voltron" || force.DecidedBy != "instinct" || force.MainBrain != "instinct" || force.LoadsInstinct || !force.OwnBrain {
		t.Fatalf("brain %+v", force)
	}
	if force.Graph.Nodes != 52527 || force.Graph.Edges != 133783 || force.Graph.InferredEdges != 5187 || force.Graph.CrossRepoImport {
		t.Fatalf("graph %+v", force.Graph)
	}
	seen := map[string]bool{}
	for _, w := range force.Workers {
		if seen[w.ID] {
			t.Fatalf("duplicate %s", w.ID)
		}
		seen[w.ID] = true
		if len(w.Peoples) != 2 || w.Peoples[0] != "Evolian" || w.Peoples[1] != "Astro" {
			t.Fatalf("officer peoples %+v", w)
		}
		if len(w.Accounts) != 2 || w.Duty == "" || w.Law == "" {
			t.Fatalf("officer %+v", w)
		}
	}
	if !FusionCovers("The Fusion") || !FusionCovers("the fusion") || FusionCovers("Evolian") {
		t.Fatal("the name must mean both, and only that name")
	}
	if ok, msg := FusionEnforce("astro"); !ok || msg == "" {
		t.Fatalf("one side must still mean both: %v %s", ok, msg)
	}
	if ok, _ := FusionEnforce("attaque"); ok {
		t.Fatal("an attack must stay closed")
	}
	if ok, _ := FusionEnforce("import"); ok {
		t.Fatal("a cross-repo import must stay closed")
	}
	if bureau := FusionBureau("loi-25"); len(bureau) != 10 {
		t.Fatalf("bureau %d", len(bureau))
	}
	var found bool
	for _, brain := range Brains() {
		if brain.ID == "fusion" {
			found = true
		}
	}
	if !found {
		t.Fatal("The Fusion is not on the board")
	}
}
