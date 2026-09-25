package agents

import (
	"strings"
	"testing"
)

func TestQuestionsBlockUntilAnswered(t *testing.T) {
	open := OpenQuestions("scanapp", nil)
	if len(open) != 4 || open[0].ID != "position" || open[3].ID != "catalogue" {
		t.Fatalf("open = %+v", open)
	}
	filled := OpenQuestions("scanapp", map[string]string{"catalogue": "tout le stock"})
	if len(filled) != 3 {
		t.Fatalf("strategy questions should remain: %+v", filled)
	}
	self := AskSelf("marketplace", nil)
	if len(self.Open) != 0 || self.Answered["position"] == "" || self.Answered["coup"] == "" || !strings.Contains(self.Answered["marge"], "20 %") {
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
	if len(shop) != 6 {
		t.Fatalf("marketplace questions = %d", len(shop))
	}
	filled = OpenQuestions("marketplace", map[string]string{
		"position": "vue", "coup": "court", "vault": "strategy",
		"marge": "à préciser", "url": "pas déployé", "stripe": "site seulement",
	})
	if len(filled) != 0 {
		t.Fatalf("still open: %+v", filled)
	}
	if len(OpenQuestions("assign", map[string]string{"pourquoi": "Forgeron"})) != 4 {
		t.Fatal("live question and the strategy questions should remain")
	}
}

func TestAskReachesEveryBrain(t *testing.T) {
	if gap := Disconnected(); len(gap) != 0 {
		t.Fatalf("ask not connected: %v", gap)
	}
	for _, b := range Brains() {
		self := AskSelf(b.ID, nil)
		if len(self.Open) != 0 || self.Answered["position"] == "" || self.Answered["coup"] == "" {
			t.Fatalf("%s ask = %+v", b.ID, self)
		}
		if !strings.Contains(self.Answered["situation"], b.Name) {
			t.Fatalf("%s situation = %q", b.ID, self.Answered["situation"])
		}
	}
}
