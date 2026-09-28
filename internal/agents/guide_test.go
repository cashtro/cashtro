package agents

import "testing"

func TestCategoriesAndCorrector(t *testing.T) {
	cats := Categories()
	seen := map[string]bool{}
	for _, cat := range cats {
		seen[cat.ID] = true
		if cat.What == "" || len(cat.Items) == 0 {
			t.Fatalf("empty category %+v", cat)
		}
		for _, item := range cat.Items {
			if item.ID == "" || item.Does == "" {
				t.Fatalf("empty item in %s: %+v", cat.ID, item)
			}
		}
	}
	for _, id := range []string{"brain", "hustler", "hustle", "lines", "voltron", "chain", "school"} {
		if !seen[id] {
			t.Fatalf("missing category %s", id)
		}
	}
	guide := CorrectChain()
	if guide.Applied || guide.How == "" || len(guide.Rules) < 5 || len(guide.Plan) < 5 {
		t.Fatalf("guide %+v", guide)
	}
	for _, cut := range guide.Plan {
		if cut.Applied || cut.Do == "" {
			t.Fatalf("cut applied %+v", cut)
		}
	}
	good := map[string]bool{}
	for _, mark := range guide.Marks {
		good[mark.ID] = mark.Good
	}
	for _, id := range []string{"edu-website", "edu-not-school", "school-in-panda", "voltron-starts", "vapi-panda", "repos-apart"} {
		if !good[id] {
			t.Fatalf("mark %s does not hold: %+v", id, guide.Marks)
		}
	}
	if good["panda-pages"] {
		t.Fatal("course pages were marked as read")
	}
}
