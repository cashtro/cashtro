package agents

import "testing"

func TestEveryAgentLearnsASkillAndAMultitask(t *testing.T) {
	lessons, err := LoadLessons("../../state/graph-lessons.json")
	if err != nil {
		t.Fatal(err)
	}
	repos := map[string]bool{}
	var fix2 int
	for _, lesson := range lessons {
		repos[lesson.Repo] = true
		if lesson.Repo == "evolu/Fix2" {
			fix2++
		}
		if lesson.Skill == "" || lesson.File == "" {
			t.Fatalf("empty lesson %+v", lesson)
		}
	}
	if len(repos) < 30 || fix2 == 0 {
		t.Fatalf("repos %d fix2 %d", len(repos), fix2)
	}
	first := Learn(lessons, nil, 1)
	if len(first) < 20 {
		t.Fatalf("gains %d", len(first))
	}
	held := map[string][]string{}
	kinds := map[string]map[string]bool{}
	for _, gain := range first {
		if gain.Skill == "" || (gain.Kind != "metier" && gain.Kind != "multitache") {
			t.Fatalf("gain %+v", gain)
		}
		held[gain.Agent] = append(held[gain.Agent], gain.Skill)
		if kinds[gain.Agent] == nil {
			kinds[gain.Agent] = map[string]bool{}
		}
		kinds[gain.Agent][gain.Kind] = true
	}
	for agent, got := range kinds {
		if !got["metier"] || !got["multitache"] {
			t.Fatalf("%s learned %+v", agent, got)
		}
	}
	voice, ok := kinds["vapi:fix2"]
	if !ok || !voice["multitache"] {
		t.Fatal("vapi fixtool did not learn a multitask skill")
	}
	second := Learn(lessons, held, 2)
	seen := map[string]bool{}
	for agent, skills := range held {
		for _, skill := range skills {
			seen[agent+" "+skill] = true
		}
	}
	for _, gain := range second {
		if seen[gain.Agent+" "+gain.Skill] {
			t.Fatalf("skill repeated: %+v", gain)
		}
	}
}
