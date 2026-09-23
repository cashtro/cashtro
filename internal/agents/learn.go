package agents

import (
	"encoding/json"
	"os"
)

// Lesson is one skill taken from Graphify. The file is the node it came from.
type Lesson struct {
	Repo  string `json:"repo"`
	Skill string `json:"skill"`
	File  string `json:"file"`
}

// Gain is one new skill. Learning adds. It does not erase what the agent already knows.
type Gain struct {
	Agent string `json:"agent"`
	Repo  string `json:"repo"`
	Skill string `json:"skill"`
	Kind  string `json:"kind"`
	File  string `json:"file"`
	Round int    `json:"round"`
}

type pupil struct {
	id   string
	home string
}

// LoadLessons reads the Graphify catalog.
func LoadLessons(path string) ([]Lesson, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var lessons []Lesson
	if err := json.Unmarshal(raw, &lessons); err != nil {
		return nil, err
	}
	return lessons, nil
}

// Learn gives every agent a new skill and a multitask skill.
// An idle crew, including Vapi when no call is open, is in this list on purpose.
func Learn(lessons []Lesson, held map[string][]string, round int) []Gain {
	if round < 1 {
		round = 1
	}
	byRepo := map[string][]Lesson{}
	var all []Lesson
	for _, lesson := range lessons {
		if lesson.Skill == "" {
			continue
		}
		byRepo[lesson.Repo] = append(byRepo[lesson.Repo], lesson)
		all = append(all, lesson)
	}
	var gains []Gain
	for _, p := range pupils(byRepo) {
		known := map[string]bool{}
		for _, skill := range held[p.id] {
			known[skill] = true
		}
		pool := byRepo[p.home]
		if len(pool) == 0 {
			pool = all
		}
		trade := pick(pool, known, round)
		gains = append(gains, Gain{Agent: p.id, Repo: trade.Repo, Skill: trade.Skill, Kind: "metier", File: trade.File, Round: round})
		known[trade.Skill] = true
		extra := pick(all, known, round+len(all))
		gains = append(gains, Gain{Agent: p.id, Repo: extra.Repo, Skill: extra.Skill, Kind: "multitache", File: extra.File, Round: round})
	}
	return gains
}

func pupils(byRepo map[string][]Lesson) []pupil {
	var out []pupil
	seen := map[string]bool{}
	add := func(id, home string) {
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		out = append(out, pupil{id: id, home: home})
	}
	for _, member := range Chart().Members() {
		add(member.Agent, "cashtro/cashtro")
	}
	for _, desk := range VapiDesks() {
		home := "cashtro/cashtro"
		if desk.Project == "fix2" {
			home = "evolu/Fix2"
		}
		add("vapi:"+desk.Project, home)
	}
	for repo := range byRepo {
		add("repo:"+repo, repo)
	}
	return out
}

func pick(pool []Lesson, known map[string]bool, salt int) Lesson {
	if len(pool) == 0 {
		return Lesson{Repo: "graphify", Skill: "observer", File: "graphify"}
	}
	start := salt % len(pool)
	if start < 0 {
		start = 0
	}
	for i := 0; i < len(pool); i++ {
		lesson := pool[(start+i)%len(pool)]
		if !known[lesson.Skill] {
			if lesson.Repo == "" {
				lesson.Repo = "graphify"
			}
			return lesson
		}
	}
	a := pool[start]
	b := pool[(start+1)%len(pool)]
	return Lesson{Repo: a.Repo, Skill: a.Skill + " + " + b.Skill, File: a.File}
}
