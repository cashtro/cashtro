package agents

import (
	"strings"
	"testing"
)

func TestMarketingTeamsHoldAndUnderstandTheGuide(t *testing.T) {
	if !CopyHeld() {
		t.Fatal("the marketing teams do not hold the guide")
	}
	brief := MarketingCopy()
	if brief.Title != "7 Figure Marketing Copy" || brief.Author != "Sean Vosler" {
		t.Fatalf("title %+v", brief)
	}
	if !strings.Contains(brief.Source, "1dGOdmWIacxPj9OlKTL6NAIeECnibeyJd") {
		t.Fatal("the file is not in hand")
	}
	want := []string{"marketing", "panda", "pandora"}
	if len(brief.Teams) != len(want) {
		t.Fatalf("teams %+v", brief.Teams)
	}
	for i, id := range want {
		if brief.Teams[i].ID != id || !brief.Teams[i].InHand || !brief.Teams[i].Understands {
			t.Fatalf("team %+v", brief.Teams[i])
		}
	}
	lesson, ok := CopyUnderstand("heros")
	if !ok || lesson.Copied || !strings.Contains(lesson.Lesson, "héros") {
		t.Fatalf("lesson %+v %v", lesson, ok)
	}
	if _, ok := CopyUnderstand("copie le pdf"); ok {
		t.Fatal("understanding must not return the file")
	}
	blob := brief.Title
	for _, lesson := range brief.Lessons {
		blob += lesson.Lesson
	}
	for _, stolen := range []string{"getwsodo", "Fan Their Flames", "David Ogilvy", "F. Scott Fitzgerald"} {
		if strings.Contains(blob, stolen) {
			t.Fatalf("the brief copied the guide: %s", stolen)
		}
	}
}

func TestEveryChainAgentAndWorkerUpgrades(t *testing.T) {
	first := CopyUpgrade(1)
	if len(first) < len(Brains())+len(Layers()) {
		t.Fatalf("upgrade %d", len(first))
	}
	seen := map[string]CopySeat{}
	for _, seat := range first {
		seen[seat.Kind+":"+seat.ID] = seat
		if seat.Copied || seat.Posted || len(seat.Lessons) == 0 {
			t.Fatalf("seat %+v", seat)
		}
	}
	for _, id := range []string{"ops", "hustle", "eye", "forge", "giant"} {
		seat, ok := seen["chain:"+id]
		if !ok {
			t.Fatalf("chain %s missing", id)
		}
		if id == "hustle" && len(seat.Lessons) != 14 {
			t.Fatalf("hustle lessons %d", len(seat.Lessons))
		}
		if id == "forge" && len(seat.Lessons) != 1 {
			t.Fatalf("forge lessons %d", len(seat.Lessons))
		}
	}
	for _, brain := range Brains() {
		if _, ok := seen["agent:"+brain.ID]; !ok {
			t.Fatalf("agent %s missing", brain.ID)
		}
	}
	market, ok := seen["department:marketing"]
	if !ok || !market.Primary || len(market.Lessons) != 14 {
		t.Fatalf("marketing %+v", market)
	}
	if _, ok := seen["worker:operator"]; !ok {
		t.Fatal("a worker did not upgrade")
	}
	if _, ok := seen["worker:architecte-workers"]; !ok {
		t.Fatal("a desk crew did not upgrade")
	}
	second := CopyUpgrade(2)
	var forge1, forge2 string
	for _, seat := range first {
		if seat.Kind == "chain" && seat.ID == "forge" {
			forge1 = seat.Lessons[0].ID
		}
	}
	for _, seat := range second {
		if seat.Kind == "chain" && seat.ID == "forge" {
			forge2 = seat.Lessons[0].ID
		}
	}
	if forge1 == "" || forge1 == forge2 {
		t.Fatalf("round did not move %s %s", forge1, forge2)
	}
}
