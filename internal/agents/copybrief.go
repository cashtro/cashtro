package agents

// CopyBrief is the marketing guide the teams hold.
// The lessons are written here. They are not passages from the PDF.
type CopyBrief struct {
	Title     string       `json:"title"`
	Author    string       `json:"author"`
	Source    string       `json:"source"`
	Teams     []CopyTeam   `json:"teams"`
	Lessons   []CopyLesson `json:"lessons"`
	InHand    bool         `json:"inHand"`
	Copied    bool         `json:"copied"`
	Posted    bool         `json:"posted"`
	DecidedBy string       `json:"decidedBy"`
}

// CopyTeam is one marketing crew. It holds the guide and can say the lesson.
type CopyTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	InHand      bool   `json:"inHand"`
	Understands bool   `json:"understands"`
}

// CopyLesson is one method, in our words.
type CopyLesson struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Lesson string `json:"lesson"`
	Copied bool   `json:"copied"`
}

// MarketingCopy is what marketing, Panda, and Pandora hold.
func MarketingCopy() CopyBrief {
	teams := []CopyTeam{
		{ID: "marketing", Name: "Marketing digital corporate", InHand: true, Understands: true},
		{ID: "panda", Name: "Panda", InHand: true, Understands: true},
		{ID: "pandora", Name: "Pandora", InHand: true, Understands: true},
	}
	return CopyBrief{
		Title:     "7 Figure Marketing Copy",
		Author:    "Sean Vosler",
		Source:    "https://drive.google.com/file/d/1dGOdmWIacxPj9OlKTL6NAIeECnibeyJd/view",
		Teams:     teams,
		Lessons:   copyLessons(),
		InHand:    true,
		Copied:    false,
		Posted:    false,
		DecidedBy: "instinct",
	}
}

func copyLessons() []CopyLesson {
	return []CopyLesson{
		{ID: "imitation", Method: "The Imitation Game", Lesson: "Reprendre la forme d'un titre qui tient déjà l'attention, et y mettre notre sujet. Ne pas coller le titre d'origine.", Copied: false},
		{ID: "amazon", Method: "Amazon R&D", Lesson: "Lire ce qu'un lecteur voisin dit déjà vouloir. Écrire ce bénéfice.", Copied: false},
		{ID: "biais", Method: "Cognitive bias", Lesson: "Nommer le biais du lecteur et rester honnête. Une promesse fausse est refusée.", Copied: false},
		{ID: "communaute", Method: "Community arbitrage", Lesson: "La liste des bénéfices vient de ce que la communauté nomme elle-même.", Copied: false},
		{ID: "formules", Method: "Tried and true formulas", Lesson: "Choisir une forme connue, attention puis intérêt puis désir puis action, ou problème puis agitation puis solution, et la remplir avec notre offre.", Copied: false},
		{ID: "surprise", Method: "Counterintuitive structure", Lesson: "Ouvrir par ce qui surprend, puis montrer le bénéfice. La surprise ne ment pas.", Copied: false},
		{ID: "enseigner", Method: "Teach, transform, transact", Lesson: "Enseigner un geste utile, montrer le changement, puis proposer l'offre. L'enseignement n'est pas un piège.", Copied: false},
		{ID: "heros", Method: "The hero's journey", Lesson: "Le lecteur est le héros. La marque est le guide.", Copied: false},
		{ID: "role", Method: "Marketing archetypes", Lesson: "Nommer le genre d'expert que la marque est, et tenir ce rôle dans toute la copie.", Copied: false},
		{ID: "ame", Method: "Data with a soul", Lesson: "Un fait, un sentiment, et une raison de croire. Les trois ensemble.", Copied: false},
		{ID: "diamant", Method: "The diamond of persuasion", Lesson: "Encourager un désir réel, promettre clairement, ajuster le geste, aligner la pensée et l'acte, et se tenir du côté de la lutte du lecteur.", Copied: false},
		{ID: "revelation", Method: "Sell the revelation", Lesson: "Chaque partie de la campagne donne un geste qu'on peut faire, pas seulement une révélation.", Copied: false},
		{ID: "anneau", Method: "From stone to a ring", Lesson: "Chercher beaucoup. Garder une seule idée.", Copied: false},
		{ID: "ecrire", Method: "Craft the ring", Lesson: "Écrire cette idée. Le brouillon reste local tant qu'Epicenter Einstein n'a pas décidé.", Copied: false},
	}
}

// CopyUnderstand returns the lesson a marketing team can say.
// It does not return the guide.
func CopyUnderstand(id string) (CopyLesson, bool) {
	wanted := id
	for _, lesson := range copyLessons() {
		if lesson.ID == wanted {
			return lesson, true
		}
	}
	return CopyLesson{}, false
}

// CopySeat is one chain, agent, worker, or department after the upgrade.
type CopySeat struct {
	ID      string       `json:"id"`
	Kind    string       `json:"kind"`
	Primary bool         `json:"primary"`
	Lessons []CopyLesson `json:"lessons"`
	Copied  bool         `json:"copied"`
	Posted  bool         `json:"posted"`
}

// CopyUpgrade gives the guide to every chain, every agent, and every worker.
// A marketing department holds every lesson. The others take one lesson this round.
func CopyUpgrade(round int) []CopySeat {
	if round < 1 {
		round = 1
	}
	lessons := copyLessons()
	var seats []CopySeat
	seen := map[string]bool{}
	add := func(id, kind, dept string) {
		key := kind + ":" + id
		if id == "" || seen[key] {
			return
		}
		seen[key] = true
		primary := marketingPrimary(id, dept)
		got := []CopyLesson{lessons[(round-1)%len(lessons)]}
		if primary {
			got = lessons
		}
		seats = append(seats, CopySeat{ID: id, Kind: kind, Primary: primary, Lessons: got, Copied: false, Posted: false})
	}
	for _, layer := range Layers() {
		dept := ""
		if layer.ID == "hustle" {
			dept = "marche"
		}
		add(layer.ID, "chain", dept)
	}
	for _, brain := range Brains() {
		add(brain.ID, "agent", brain.Dept)
	}
	for _, dept := range Chart().Departments {
		for _, member := range dept.Members {
			kind := "agent"
			if dept.ID == "flux" || stringsContainsWorker(member.Title) {
				kind = "worker"
			}
			add(member.Agent, kind, dept.ID)
		}
	}
	for _, desk := range MainBrain() {
		add(desk.ID, "desk", "")
		add(desk.ID+"-workers", "worker", "")
	}
	for _, brain := range NicheBrains() {
		for _, desk := range brain.Desks {
			add(desk.ID, "desk", "")
			add(desk.ID+"-workers", "worker", "")
		}
	}
	for _, division := range Chart().Divisions {
		add(division.ID, "department", division.Department)
	}
	for _, team := range MarketingCopy().Teams {
		add(team.ID, "department", "marche")
	}
	return seats
}

func marketingPrimary(id, dept string) bool {
	switch id {
	case "marketing", "panda", "pandora", "scanapp", "marche", "hustle", "hustler", "accueil":
		return true
	}
	return dept == "marche"
}

func stringsContainsWorker(title string) bool {
	for i := 0; i+6 <= len(title); i++ {
		if title[i:i+6] == "worker" || title[i:i+6] == "Worker" {
			return true
		}
	}
	return false
}

// CopyHeld reports whether every marketing team holds the guide and understands it.
func CopyHeld() bool {
	brief := MarketingCopy()
	if !brief.InHand || brief.Copied || brief.Posted || len(brief.Lessons) < 11 || len(brief.Teams) == 0 {
		return false
	}
	for _, team := range brief.Teams {
		if !team.InHand || !team.Understands || team.ID == "" {
			return false
		}
	}
	for _, lesson := range brief.Lessons {
		if lesson.Copied || lesson.Lesson == "" {
			return false
		}
	}
	return true
}
