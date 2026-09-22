package agents

import "testing"

func TestQuestionsBlockUntilAnswered(t *testing.T) {
	open := OpenQuestions("scanapp", nil)
	if len(open) != 4 || open[0].ID != "marge" {
		t.Fatalf("open = %+v", open)
	}
	filled := OpenQuestions("scanapp", map[string]string{
		"marge": "à préciser", "marketplace": "pas déployé", "catalogue": "fiches seulement", "stripe": "site seulement",
	})
	if len(filled) != 0 {
		t.Fatalf("still open: %+v", filled)
	}
	if len(OpenQuestions("assign", map[string]string{"pourquoi": "Forgeron"})) != 1 {
		t.Fatal("live question should remain")
	}
}
