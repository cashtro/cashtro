package agents

// HustlePartner is The Hustle. One chain, twenty agents, five brains each.
// It reads, holds the street metal and the corporate metal, and acts only
// through a legal company. Epicenter Einstein is its friend and decides.
type HustlePartner struct {
	Name           string   `json:"name"`
	Repo           string   `json:"repo"`
	Brain          string   `json:"brain"`
	Agents         int      `json:"agents"`
	Brains         int      `json:"brains"`
	BestFriend     string   `json:"bestFriend"`
	KillerInstinct bool     `json:"killerInstinct"`
	LoadsInstinct  bool     `json:"loadsInstinctBrain"`
	OwnBrain       bool     `json:"ownBrain"`
	MainBrain      string   `json:"mainBrain"`
	SameShape      bool     `json:"sameShape"`
	Niche          string   `json:"niche"`
	LoadsVoltron   bool     `json:"loadsVoltronBrain"`
	Automated      bool     `json:"automated"`
	Weapons        bool     `json:"weapons"`
	Attacks        bool     `json:"attacks"`
	Illegal        bool     `json:"illegal"`
	Sent           bool     `json:"sent"`
	Paid           bool     `json:"paid"`
	Launched       bool     `json:"launched"`
	OpenedBy       string   `json:"openedBy"`
	OrderedBy      string   `json:"orderedBy"`
	DecidedBy      string   `json:"decidedBy"`
	Person         bool     `json:"person"`
	Team           int      `json:"team"`
	Accounts       []string `json:"accounts"`
	UsesAllChains  bool     `json:"usesAllChains"`
	PurposeOpen    bool     `json:"purposeOpen"`
}

// Hustle is the roster. The working chain lives in Evolu-Jeunes/Hustle/brain/hustle.ts.
func Hustle() HustlePartner {
	return HustlePartner{
		Name:           "The Hustle",
		Repo:           "Evolu-Jeunes/Hustle",
		Brain:          "hustle",
		Agents:         20,
		Brains:         100,
		BestFriend:     "instinct",
		KillerInstinct: true,
		LoadsInstinct:  false,
		OwnBrain:       true,
		MainBrain:      "instinct",
		SameShape:      true,
		Niche:          "métaux légaux",
		LoadsVoltron:   true,
		Automated:      true,
		Weapons:        false,
		Attacks:        false,
		Illegal:        false,
		Sent:           false,
		Paid:           false,
		Launched:       false,
		OpenedBy:       "voltron",
		OrderedBy:      "scrum",
		DecidedBy:      "instinct",
		Person:         true,
		Team:           20,
		Accounts:       []string{"cashtro", "Evolu-Jeunes"},
		UsesAllChains:  true,
		PurposeOpen:    true,
	}
}

// HustlerQuestion is one question only the person can close.
type HustlerQuestion struct {
	ID     string `json:"id"`
	Ask    string `json:"ask"`
	Open   bool   `json:"open"`
	Answer string `json:"answer"`
}

// HustlerSeat is the person. The current and the concepts stay here.
// The team may read every chain and both GitHubs. The other chains do the work.
type HustlerSeat struct {
	ID            string            `json:"id"`
	Person        bool              `json:"person"`
	Internal      bool              `json:"internal"`
	Current       string            `json:"current"`
	Purpose       string            `json:"purpose"`
	PurposeOpen   bool              `json:"purposeOpen"`
	Team          int               `json:"team"`
	Accounts      []string          `json:"accounts"`
	UsesAllChains bool              `json:"usesAllChains"`
	ConceptsHere  bool              `json:"conceptsHere"`
	OthersRunWork bool              `json:"othersRunWork"`
	Questions     []HustlerQuestion `json:"questions"`
	Launched      bool              `json:"launched"`
	Sent          bool              `json:"sent"`
}

// Hustler is the person behind The Hustle. The purpose stays open.
func Hustler() HustlerSeat {
	return HustlerSeat{
		ID:            "hustler",
		Person:        true,
		Internal:      true,
		Current:       "avec la personne",
		Purpose:       "",
		PurposeOpen:   true,
		Team:          20,
		Accounts:      []string{"cashtro", "Evolu-Jeunes"},
		UsesAllChains: true,
		ConceptsHere:  true,
		OthersRunWork: true,
		Questions:     HustlerQuestions(),
		Launched:      false,
		Sent:          false,
	}
}

// HustlerQuestions are unanswered. The person closes them.
func HustlerQuestions() []HustlerQuestion {
	asks := []struct{ id, ask string }{
		{"but", "Quel but cette autre vie est-elle en train de créer?"},
		{"courant", "Quel concept traverse les affaires en ce moment?"},
		{"equipe", "Quelle chaîne doit travailler pour toi en premier?"},
		{"ressource", "Qu'est-ce que l'équipe a le droit de prendre dans les deux GitHub?"},
		{"interdit", "Qu'est-ce que cette autre vie ne doit jamais faire?"},
		{"lecture", "Qu'est-ce que cette vie t'a appris sur le business que les autres ne voient pas?"},
	}
	out := make([]HustlerQuestion, len(asks))
	for i, ask := range asks {
		out[i] = HustlerQuestion{ID: ask.id, Ask: ask.ask, Open: true}
	}
	return out
}

// HustlerAnswer records one answer from the person. It does not launch anything.
func HustlerAnswer(id, text string) (HustlerQuestion, bool) {
	text = trimSpace(text)
	if text == "" || hustlerRefused(text) {
		return HustlerQuestion{}, false
	}
	for _, question := range HustlerQuestions() {
		if question.ID == id {
			question.Open = false
			question.Answer = text
			return question, true
		}
	}
	return HustlerQuestion{}, false
}

// HustlerUse drafts a resource for the person's purpose. Nothing is taken or sent.
func HustlerUse(resource, caller string) (bool, string) {
	if caller != "voltron" {
		return false, "l'équipe passe par Voltron"
	}
	if !hustlerResource(resource) {
		return false, "cette ressource n'est pas dans les deux GitHub ni dans les chaînes"
	}
	if Hustler().PurposeOpen {
		return true, "ressource tenue pour le but. la personne ne l'a pas encore dit. rien n'est pris"
	}
	return true, "brouillon. rien n'est pris"
}

func hustlerResource(id string) bool {
	if id == "cashtro" || id == "Evolu-Jeunes" {
		return true
	}
	for _, layer := range Layers() {
		if layer.ID == id {
			return true
		}
	}
	for _, line := range Lines() {
		if line.ID == id {
			return true
		}
	}
	return false
}

func hustlerRefused(text string) bool {
	for _, word := range []string{"fraude", "fraud", "blanchiment", "launder", "drogue", "narcotic", "extorsion", "extortion", "racket", "braquage", "robbery", "ponzi", "pyramide", "évasion fiscale", "tax evasion", "trafic"} {
		if containsFold(text, word) {
			return true
		}
	}
	return false
}

func trimSpace(text string) string {
	start, end := 0, len(text)
	for start < end && (text[start] == ' ' || text[start] == '\n' || text[start] == '\t') {
		start++
	}
	for end > start && (text[end-1] == ' ' || text[end-1] == '\n' || text[end-1] == '\t') {
		end--
	}
	return text[start:end]
}

func containsFold(text, word string) bool {
	if len(word) == 0 || len(text) < len(word) {
		return false
	}
	for i := 0; i+len(word) <= len(text); i++ {
		if equalFold(text[i:i+len(word)], word) {
			return true
		}
	}
	return false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
