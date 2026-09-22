package agents

import "testing"

func TestQuestionsBlockUntilAnswered(t *testing.T) {
	open := OpenQuestions("scanapp", nil)
	if len(open) != 6 {
		t.Fatalf("open = %d", len(open))
	}
	filled := OpenQuestions("scanapp", map[string]string{
		"boutique": "Noix", "codes": "EAN-13", "photo": "générée par nous",
		"site": "Evolu-Jeunes/noix, pas en ligne", "kick": "non, site seulement", "prix": "CAD taxes incluses, Québec",
	})
	if len(filled) != 0 {
		t.Fatalf("still open: %+v", filled)
	}
	if len(OpenQuestions("assign", map[string]string{"pourquoi": "Forgeron"})) != 1 {
		t.Fatal("live question should remain")
	}
}
