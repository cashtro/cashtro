package agents

import (
	"strconv"
	"strings"
	"sync"
)

// PandaSection is one part of the Panda website.
type PandaSection struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Fact      string `json:"fact"`
	File      string `json:"file,omitempty"`
	Connected bool   `json:"connected"`
}

// PandaBreak is one place the backend line does not meet.
type PandaBreak struct {
	ID   string `json:"id"`
	Fact string `json:"fact"`
	File string `json:"file,omitempty"`
}

// RoadItem is later work. Nothing on it is built or called.
type RoadItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Do     string `json:"do"`
	Built  bool   `json:"built"`
	Called bool   `json:"called"`
}

// ClassDraft is one class on the school desk. The price is a draft. Nothing is charged or called.
type ClassDraft struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Charged bool   `json:"charged"`
	Called  bool   `json:"called"`
	Version string `json:"version"`
}

// PandaReport is the scan of the Panda site plus the school V1 desk.
type PandaReport struct {
	Repo        string         `json:"repo"`
	Version     string         `json:"version"`
	LocalSchool bool           `json:"localSchool"`
	Breaks      []PandaBreak   `json:"breaks"`
	Sections    []PandaSection `json:"sections"`
	Roadmap     []RoadItem     `json:"roadmap"`
	Classes     []ClassDraft   `json:"classes"`
}

// Roadmap is the work that waits. Dashboards are not added. Nothing is called.
func Roadmap() []RoadItem {
	return []RoadItem{
		{ID: "voltron-builder", Name: "Build an agent from Voltron", Do: "A person builds an agent from Voltron in plain language. An admin dashboard and a client dashboard come later, for the admin and for the user.", Built: false, Called: false},
		{ID: "panda-ai", Name: "Panda AI sections", Do: "The automated AI sections inside the Panda site take the new wave after the school V1 is patched.", Built: false, Called: false},
		{ID: "panda-marketing", Name: "Digital marketing", Do: "Digital marketing on Panda becomes automated, with new dashboards and new presentations.", Built: false, Called: false},
		{ID: "agent-crm", Name: "Agent CRM", Do: "One CRM is for the agents: share knowledge, book, close a client, and call. The agents run one workflow through it. A call stays off.", Built: false, Called: false},
		{ID: "fix2-crm", Name: "Fix Tout CRM", Do: "The other CRM stays in the Fix Tout repo. The admin dashboard, the client dashboard, and the project front end land in that repo later. Connections come after that push.", Built: false, Called: false},
		{ID: "einstein-return", Name: "Einstein learns by sending back", Do: "Einstein auto-learns when he can send the analysis back to the other agents, through the agent CRM or a direct read of every agent. He writes his own input after the analysis.", Built: false, Called: false},
	}
}

// ScanPanda reads the Graphify catalog of Evolu-Jeunes/Panda and the kernel lines.
// The school V1 on this desk is separate from a missing school route in that catalog.
func ScanPanda(lessons []Lesson) PandaReport {
	var panda []Lesson
	for _, lesson := range lessons {
		if lesson.Repo == "evolu/Panda" || lesson.Repo == "Evolu-Jeunes/Panda" || strings.Contains(lesson.File, "/Panda/") {
			panda = append(panda, lesson)
		}
	}
	fileOf := func(pred func(Lesson) bool) string {
		for _, lesson := range panda {
			if pred(lesson) {
				return lesson.File
			}
		}
		return ""
	}
	has := func(pred func(Lesson) bool) bool {
		return fileOf(pred) != ""
	}
	blob := func(lesson Lesson) string {
		return strings.ToLower(lesson.Skill + " " + lesson.File)
	}
	storeFile := ""
	storeSkill := ""
	for _, lesson := range panda {
		if strings.Contains(blob(lesson), "achat") {
			storeFile = lesson.File
			storeSkill = lesson.Skill
			break
		}
	}
	demoFile := fileOf(func(lesson Lesson) bool {
		return strings.Contains(blob(lesson), "dashboard-demo")
	})
	schoolFile := fileOf(func(lesson Lesson) bool {
		text := blob(lesson)
		if strings.Contains(text, "educonnexion") {
			return false
		}
		return strings.Contains(text, "school") || strings.Contains(text, "cours") || strings.Contains(text, "classe") || strings.Contains(text, "ecole")
	})
	priceFile := fileOf(func(lesson Lesson) bool {
		text := blob(lesson)
		return strings.Contains(text, "price") || strings.Contains(text, "prix") || strings.Contains(text, "drawer")
	})
	_, ecoleChain := Chain("ecole")
	_, pandaChain := Chain("panda")
	marketRepos := []string{}
	if market, ok := LineByID("marketplace"); ok {
		marketRepos = market.Repos
	}
	pandaIsMarket := false
	for _, repo := range marketRepos {
		if repo == "Evolu-Jeunes/Panda" {
			pandaIsMarket = true
		}
	}
	sections := []PandaSection{
		{ID: "school", Name: "School", Fact: "Classes are sold inside the Panda website.", File: schoolFile, Connected: schoolFile != ""},
		{ID: "store", Name: "Store", Fact: "A small store lives inside the Panda website.", File: storeFile, Connected: storeFile != ""},
		{ID: "prices", Name: "Prices", Fact: "A drawer gives the prices.", File: priceFile, Connected: priceFile != ""},
		{ID: "ai", Name: "AI sections", Fact: "Automated AI sections are in the site. The new wave is later.", Connected: has(func(lesson Lesson) bool {
			return strings.Contains(blob(lesson), "openai") || strings.Contains(blob(lesson), "/ai")
		})},
		{ID: "marketing", Name: "Digital marketing", Fact: "Digital marketing on this site is automated later.", Connected: has(func(lesson Lesson) bool {
			return strings.Contains(blob(lesson), "marketing")
		})},
	}
	breaks := []PandaBreak{}
	if schoolFile == "" {
		purchase := "purchases"
		if storeSkill != "" {
			purchase = storeSkill
		}
		breaks = append(breaks, PandaBreak{ID: "school-route", Fact: "The extracted Panda backend has no school route. The catalog connects " + purchase + " in the Panda server, plus signup, auth, and Supabase.", File: storeFile})
	}
	if priceFile == "" {
		breaks = append(breaks, PandaBreak{ID: "price-drawer", Fact: "The price drawer is not in the extracted Panda symbols."})
	}
	if !pandaIsMarket {
		where := "Scan App"
		if len(marketRepos) > 0 {
			where = marketRepos[0]
		}
		breaks = append(breaks, PandaBreak{ID: "store-line", Fact: "Panda's store is achatRoutes inside Evolu-Jeunes/Panda. The kernel marketplace line is " + where + "."})
	}
	if ecoleChain && pandaChain {
		breaks = append(breaks, PandaBreak{ID: "two-chains", Fact: "The school chain and the Panda chain are still two chains. The teaching lives in the Panda website."})
	}
	if demoFile != "" {
		breaks = append(breaks, PandaBreak{ID: "admin-demo", Fact: "Admin controls in the extract sit on a demo file. The admin dashboard and the client dashboard are later.", File: demoFile})
	}
	return PandaReport{
		Repo:        "Evolu-Jeunes/Panda",
		Version:     "v1",
		LocalSchool: true,
		Breaks:      breaks,
		Sections:    sections,
		Roadmap:     Roadmap(),
		Classes:     []ClassDraft{},
	}
}

// ReadPandaScan loads the Graphify catalog beside this module and scans Panda.
func ReadPandaScan() PandaReport {
	var lessons []Lesson
	if path, err := findFile("state", "graph-lessons.json"); err == nil {
		if loaded, err := LoadLessons(path); err == nil {
			lessons = loaded
		}
	}
	return ScanPanda(lessons)
}

// SchoolV1 parks classes for the Panda school. A price stays a draft.
type SchoolV1 struct {
	mu      sync.Mutex
	n       int
	classes []ClassDraft
}

// NewSchool returns an empty school desk.
func NewSchool() *SchoolV1 {
	return &SchoolV1{}
}

// Add parks one class. Charged and Called stay false.
func (s *SchoolV1) Add(name, price string) (ClassDraft, error) {
	name = strings.TrimSpace(name)
	price = strings.TrimSpace(price)
	if name == "" {
		return ClassDraft{}, errClassName
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
	class := ClassDraft{
		ID:      strings.TrimSpace(strings.ReplaceAll(strings.ToLower(name), " ", "-")) + "-" + strconv.Itoa(s.n),
		Name:    name,
		Price:   price,
		Charged: false,
		Called:  false,
		Version: "v1",
	}
	s.classes = append(s.classes, class)
	return class, nil
}

// List returns the parked classes.
func (s *SchoolV1) List() []ClassDraft {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ClassDraft, len(s.classes))
	copy(out, s.classes)
	return out
}

var errClassName = classError("a class needs a name")

type classError string

func (e classError) Error() string { return string(e) }
