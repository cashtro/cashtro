package agents

import (
	"strings"
	"testing"
)

func TestQuestionsBlockUntilAnswered(t *testing.T) {
	open := OpenQuestions("scanapp", nil)
	if len(open) != 1 || open[0].ID != "catalogue" {
		t.Fatalf("open = %+v", open)
	}
	filled := OpenQuestions("scanapp", map[string]string{"catalogue": "tout le stock"})
	if len(filled) != 0 {
		t.Fatalf("scan still open: %+v", filled)
	}
	self := AskSelf("marketplace", nil)
	if len(self.Open) != 0 || self.Answered["url"] == "" || !strings.Contains(self.Answered["marge"], "20 %") {
		t.Fatalf("self = %+v", self)
	}
	if !strings.Contains(self.Answered["stripe"], "Pas Proximity") || !strings.Contains(self.Answered["stripe"], "PBTM") {
		t.Fatalf("stripe = %q", self.Answered["stripe"])
	}
	scan := AskSelf("scanapp", nil)
	if len(scan.Open) != 0 || scan.Answered["catalogue"] == "" {
		t.Fatalf("scan self = %+v", scan)
	}
	shop := OpenQuestions("marketplace", nil)
	if len(shop) != 3 {
		t.Fatalf("marketplace questions = %d", len(shop))
	}
	filled = OpenQuestions("marketplace", map[string]string{
		"marge": "à préciser", "url": "pas déployé", "stripe": "site seulement",
	})
	if len(filled) != 0 {
		t.Fatalf("still open: %+v", filled)
	}
	if len(OpenQuestions("assign", map[string]string{"pourquoi": "Forgeron"})) != 1 {
		t.Fatal("live question should remain")
	}
}
