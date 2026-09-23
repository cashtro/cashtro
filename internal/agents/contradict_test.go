package agents

import "testing"

func TestContradictPicksTheLeanerOption(t *testing.T) {
	one := Contradict(Proposal{Subject: "site", Options: []Option{{Name: "refaire", Cost: 5, Steps: 8}}})
	if one.Ready {
		t.Fatal("one option must be contradicted")
	}
	tie := Contradict(Proposal{Options: []Option{{Name: "a", Cost: 2}, {Name: "b", Cost: 2}}})
	if tie.Ready {
		t.Fatal("a tie is not an optimized choice")
	}
	got := Contradict(Proposal{
		Subject: "fiche scan",
		Options: []Option{
			{Name: "régénérer tout le site", Cost: 8, Risk: 3, Steps: 6},
			{Name: "réutiliser le thème", Cost: 1, Risk: 1, Steps: 2},
		},
	})
	if !got.Ready || got.Best != "réutiliser le thème" || got.Score != 4 {
		t.Fatalf("verdict = %+v", got)
	}
}

func TestOptimisationDepartment(t *testing.T) {
	var found bool
	for _, d := range Chart().Departments {
		if d.ID != "optimisation" {
			continue
		}
		found = true
		if d.Chief != "reviewer" || d.ReportsTo != "ceo" || len(d.Members) != 2 {
			t.Fatalf("optimisation = %+v", d)
		}
	}
	if !found {
		t.Fatal("missing optimisation department")
	}
}

func TestContradictSitsOnEveryChain(t *testing.T) {
	for _, a := range AnalyzeAll() {
		if !a.Connected || a.Best != "le plus court" {
			t.Fatalf("%s contradict = %+v", a.ID, a)
		}
	}
}
