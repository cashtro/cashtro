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
