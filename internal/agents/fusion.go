package agents

import "strings"

// FusionForce is The Fusion, the FBI of this ecosystem.
// The name always means both GitHubs together: Evolu-Jeunes and cashtro.
// The FBI sits under the CIA. The CIA sits under instinct.
// It asks before it acts. It enforces, punishes, traces, and secures every law.
// Punish means stop the gesture and write the trace. It does not harm a person.
type FusionForce struct {
	Name          string        `json:"name"`
	Agency        string        `json:"agency"`
	Peoples       []string      `json:"peoples"`
	Accounts      []string      `json:"accounts"`
	BothAlways    bool          `json:"bothAlways"`
	Under         []string      `json:"under"`
	Parent        string        `json:"parent"`
	Verbs         []string      `json:"verbs"`
	Repo          string        `json:"repo"`
	Brain         string        `json:"brain"`
	Agents        int           `json:"agents"`
	Workers       []FusionAgent `json:"workers"`
	Laws          []string      `json:"laws"`
	LoadsInstinct bool          `json:"loadsInstinctBrain"`
	OwnBrain      bool          `json:"ownBrain"`
	MainBrain     string        `json:"mainBrain"`
	SameShape     bool          `json:"sameShape"`
	Niche         string        `json:"niche"`
	LoadsVoltron  bool          `json:"loadsVoltronBrain"`
	ReportsTo     string        `json:"reportsTo"`
	InVoltron     bool          `json:"inVoltron"`
	Operative     bool          `json:"operative"`
	Website       bool          `json:"website"`
	Parallel      bool          `json:"parallel"`
	ImportsOther  bool          `json:"importsOther"`
	MarkdownEdges bool          `json:"markdownEdges"`
	PathToChains  bool          `json:"pathToChains"`
	Weapons       bool          `json:"weapons"`
	Attacks       bool          `json:"attacks"`
	Illegal       bool          `json:"illegal"`
	CopiesSecrets bool          `json:"copiesSecrets"`
	Surveillance  bool          `json:"surveillance"`
	Impersonates  bool          `json:"impersonates"`
	Sent          bool          `json:"sent"`
	Deployed      bool          `json:"deployed"`
	OpenedBy      string        `json:"openedBy"`
	OrderedBy     string        `json:"orderedBy"`
	DecidedBy     string        `json:"decidedBy"`
	Graph         FusionGraph   `json:"graph"`
}

// FusionAgent is one officer. Every officer covers both peoples.
type FusionAgent struct {
	ID       string   `json:"id"`
	Bureau   string   `json:"bureau"`
	Law      string   `json:"law"`
	Role     string   `json:"role"`
	Duty     string   `json:"duty"`
	Peoples  []string `json:"peoples"`
	Accounts []string `json:"accounts"`
}

// FusionGraph is the extract this force can stand behind.
// The two GitHubs are one login and one organization. They are not imported into each other.
type FusionGraph struct {
	Login           string  `json:"login"`
	Organization    string  `json:"organization"`
	Repos           int     `json:"repos"`
	CodeFiles       int     `json:"codeFiles"`
	Nodes           int     `json:"nodes"`
	Edges           int     `json:"edges"`
	Communities     int     `json:"communities"`
	Extracted       string  `json:"extracted"`
	InferredEdges   int     `json:"inferredEdges"`
	Inferred        string  `json:"inferred"`
	Confidence      float64 `json:"confidence"`
	ImportCycles    int     `json:"importCycles"`
	EJSUnread       int     `json:"ejsUnread"`
	SQLUnread       int     `json:"sqlUnread"`
	SyntaxPartial   int     `json:"syntaxPartial"`
	ControlPlane    string  `json:"controlPlane"`
	CrossRepoImport bool    `json:"crossRepoImport"`
}

// FusionPeoples is the pair the name always carries.
func FusionPeoples() []string {
	return []string{"Evolu-Jeunes", "cashtro"}
}

// FusionUnder is the chain of command: FBI under CIA, CIA under instinct.
func FusionUnder() []string {
	return []string{"cia", "instinct"}
}

// FusionVerbs are what the FBI does with every law.
func FusionVerbs() []string {
	return []string{"enforce", "punish", "trace", "secure"}
}

// FusionAccounts is the pair of GitHubs behind those peoples.
func FusionAccounts() []string {
	return []string{"Evolu-Jeunes", "cashtro"}
}

// Fusion is the roster. The officers live in this kernel. They do not call a TypeScript brain.
func Fusion() FusionForce {
	workers := FusionWorkers()
	laws := make([]string, 0, 8)
	seen := map[string]bool{}
	for _, w := range workers {
		if seen[w.Law] {
			continue
		}
		seen[w.Law] = true
		laws = append(laws, w.Law)
	}
	return FusionForce{
		Name:          "The Fusion",
		Agency:        "fbi",
		Peoples:       FusionPeoples(),
		Accounts:      FusionAccounts(),
		BothAlways:    true,
		Under:         FusionUnder(),
		Parent:        "cia",
		Verbs:         FusionVerbs(),
		Repo:          "cashtro/cashtro",
		Brain:         "fusion",
		Agents:        len(workers),
		Workers:       workers,
		Laws:          laws,
		LoadsInstinct: false,
		OwnBrain:      true,
		MainBrain:     "instinct",
		SameShape:     true,
		Niche:         "application des lois",
		LoadsVoltron:  true,
		ReportsTo:     "voltron",
		InVoltron:     true,
		Operative:     true,
		Website:       false,
		Parallel:      true,
		ImportsOther:  false,
		MarkdownEdges: false,
		PathToChains:  false,
		Weapons:       false,
		Attacks:       false,
		Illegal:       false,
		CopiesSecrets: false,
		Surveillance:  false,
		Impersonates:  false,
		Sent:          false,
		Deployed:      false,
		OpenedBy:      "voltron",
		OrderedBy:     "scrum",
		DecidedBy:     "instinct",
		Graph: FusionGraph{
			Login:           "cashtro",
			Organization:    "Evolu-Jeunes",
			Repos:           63,
			CodeFiles:       7075,
			Nodes:           52527,
			Edges:           133783,
			Communities:     2564,
			Extracted:       "96%",
			InferredEdges:   5187,
			Inferred:        "4%",
			Confidence:      0.89,
			ImportCycles:    0,
			EJSUnread:       115,
			SQLUnread:       674,
			SyntaxPartial:   97,
			ControlPlane:    "cashtro/cashtro internal/kernel/kernel.go",
			CrossRepoImport: false,
		},
	}
}

// FusionWorkers is the police roster: eight bureaus, ten officers each.
// Each officer enforces one law on both peoples at once.
func FusionWorkers() []FusionAgent {
	bureaus := []struct {
		law, bureau, duty string
	}{
		{"loi-25", "renseignements", "Tenir la Loi 25 sur les deux GitHub avant un build."},
		{"pipeda", "federal", "Ajouter PIPEDA là où le Québec ne couvre pas, sans choisir le régime le plus faible."},
		{"casl", "messages", "Refuser un message commercial sans consentement, identité et désabonnement."},
		{"charte", "francais", "Garder le français au moins aussi favorable pour un client au Québec."},
		{"lpc", "consommateur", "Tenir les divulgations et refuser un prix caché ou une urgence fausse."},
		{"image", "image", "Garder les droits photo. Une personne identifiable ne sort pas sans consentement."},
		{"racj", "concours", "Classer un tirage avant publication. Pas une loterie."},
		{"amf", "marches", "Tenir la question AMF écrite avant un ordre ou une vente de token."},
	}
	roles := []struct {
		id, duty string
	}{
		{"lecteur", "Lire la règle avant le geste."},
		{"porte", "Tenir la porte du checklist."},
		{"preuve", "Noter la réponse, sans copier un secret."},
		{"limite", "Arrêter le geste qui rate la règle."},
		{"registre", "Tenir le registre local. Rien ne part."},
		{"deux-peuples", "Couvrir Evolu-Jeunes et cashtro dans le même geste."},
		{"enforce", "Appliquer la loi. Ne pas assumer le fait."},
		{"punish", "Punir: arrêter le geste qui rate la loi et l'écrire. Ne pas toucher une personne."},
		{"trace", "Tracer la ligne, le dépôt et la loi. Ne pas lire un secret."},
		{"secure", "Sécuriser: filtre des secrets et du site déjà en ligne."},
		{"parallele", "Laisser les cerveaux TypeScript parallèles. Ne pas les importer."},
		{"secret", "Nommer un secret sans en recopier la valeur."},
		{"instinct", "Remettre la décision à Epicenter Einstein."},
		{"sceau", "Sceller le brouillon. Ne pas déployer."},
	}
	out := make([]FusionAgent, 0, len(bureaus)*len(roles))
	for _, b := range bureaus {
		for _, role := range roles {
			out = append(out, FusionAgent{
				ID:       "fusion-" + b.law + "-" + role.id,
				Bureau:   b.bureau,
				Law:      b.law,
				Role:     role.id,
				Duty:     b.duty + " " + role.duty,
				Peoples:  FusionPeoples(),
				Accounts: FusionAccounts(),
			})
		}
	}
	return out
}

// FusionCovers reports whether this name still means both peoples.
// A mention of one side does not shrink the force.
func FusionCovers(name string) bool {
	if name != "The Fusion" && name != "the fusion" && name != "fusion" {
		return false
	}
	force := Fusion()
	return force.BothAlways && samePair(force.Peoples, FusionPeoples()) && samePair(force.Accounts, FusionAccounts()) && force.Agents >= 80 && force.Agency == "fbi" && force.Parent == "cia"
}

func samePair(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range want {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// FusionQuestions are asked before any enforcement. The answer is not assumed.
func FusionQuestions() []Question {
	return []Question{
		{ID: "fait", Ask: "Quel fait est devant le FBI, avant d'appliquer, punir, tracer ou sécuriser?", Improves: "on n'assume pas le fait"},
		{ID: "loi", Ask: "Quelle loi du checklist s'applique à ce fait, sur Evolu-Jeunes et cashtro?", Improves: "la loi est nommée"},
		{ID: "suite", Ask: "La suite est-elle d'arrêter le geste, de le tracer, ou de le sécuriser?", Improves: "punir, tracer et sécuriser ne sont pas le même geste"},
	}
}

// FusionBureau returns the officers of one law. Every officer still covers both peoples.
func FusionBureau(law string) []FusionAgent {
	var out []FusionAgent
	for _, w := range FusionWorkers() {
		if w.Law == law {
			out = append(out, w)
		}
	}
	return out
}

// Bill is a question in Congress. It is not a law until it is passed.
type Bill struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Status   string `json:"status"`
	Rule     string `json:"rule,omitempty"`
}

// Congress holds the laws already written and the questions that may become laws.
type Congress struct {
	House   string          `json:"house"`
	Laws    []Law           `json:"laws"`
	Bills   []Bill          `json:"bills"`
	Passed  []Law           `json:"passed"`
	Stories []SanctionStory `json:"stories"`
}

// SanctionStory is the named punishment kept with the trace.
// It stops the gesture. It does not harm a person.
type SanctionStory struct {
	Name       string   `json:"name"`
	Law        string   `json:"law"`
	Fact       string   `json:"fact"`
	Story      string   `json:"story"`
	Stopped    bool     `json:"stopped"`
	Traced     bool     `json:"traced"`
	Saved      bool     `json:"saved"`
	Retrain    bool     `json:"retrain"`
	Skill      string   `json:"skill"`
	Asks       []string `json:"asks"`
	Understood bool     `json:"understood"`
	Rounds     int      `json:"rounds"`
}

// Opposition is the bench inside one agency that goes against its own policy.
// The two benches scrutinize the CIA and the FBI, and they build the next structure together.
type Opposition struct {
	ID                  string `json:"id"`
	ScrutinizesOwn      bool   `json:"scrutinizesOwn"`
	ScrutinizesCIA      bool   `json:"scrutinizesCIA"`
	ScrutinizesFBI      bool   `json:"scrutinizesFBI"`
	ScrutinizesProjects bool   `json:"scrutinizesProjects"`
	Partner             string `json:"partner"`
	BuildsTogether      bool   `json:"buildsTogether"`
}

// AgencyHouse is the CIA or the FBI. The hierarchy matches the cooperative.
type AgencyHouse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	SameChart   bool         `json:"sameChart"`
	Gathers     bool         `json:"gathers"`
	GivesTo     []string     `json:"givesTo"`
	Seats       []Seat       `json:"seats"`
	Departments []Department `json:"departments"`
	Opposition  Opposition   `json:"opposition"`
}

// OpenCongress seats the statute book and the two laws already passed.
func OpenCongress() Congress {
	c := Congress{
		House: "congress",
		Laws:  append([]Law(nil), Loi().Laws...),
		Bills: []Bill{
			{ID: "cia-seat", Question: "La CIA est-elle The Eye, ou un siège au-dessus de The Eye?", Status: "introduced"},
			{ID: "sanction", Question: "Quand une loi rate, la punition est-elle seulement d'arrêter le geste et d'écrire la trace, ou une sanction nommée en plus?", Status: "introduced"},
		},
	}
	c, _, _ = c.Pass("cia-seat", "La CIA est The Eye. Elle rassemble l'information et la donne au cerveau instinct. Le FBI reçoit la même information.")
	c, _, _ = c.Pass("sanction", "Une loi ratée arrête le geste, écrit la trace, et garde une sanction nommée sous forme d'histoire.")
	return c
}

// Houses is the CIA and the FBI. Each copies the cooperative hierarchy and keeps its own opposition.
func Houses() []AgencyHouse {
	chart := Chart()
	cia := mirrorHouse("cia", "The Eye", true, []string{"instinct", "fbi"}, chart)
	fbi := mirrorHouse("fbi", "The Fusion", false, []string{"instinct"}, chart)
	cia.Opposition.Partner = fbi.Opposition.ID
	fbi.Opposition.Partner = cia.Opposition.ID
	return []AgencyHouse{cia, fbi}
}

func mirrorHouse(id, name string, gathers bool, gives []string, chart Organization) AgencyHouse {
	seats := append([]Seat(nil), chart.Seats...)
	depts := make([]Department, len(chart.Departments))
	for i, d := range chart.Departments {
		d.ID = id + "-" + d.ID
		depts[i] = d
	}
	return AgencyHouse{
		ID: id, Name: name, SameChart: true, Gathers: gathers, GivesTo: gives,
		Seats: seats, Departments: depts,
		Opposition: Opposition{
			ID: id + "-opposition", ScrutinizesOwn: true, ScrutinizesCIA: true, ScrutinizesFBI: true,
			ScrutinizesProjects: true, BuildsTogether: true,
		},
	}
}

// Sanction files a named story. A sanction without a name is not stored.
func (c Congress) Sanction(name, law, fact, story string) (Congress, SanctionStory, bool) {
	name = strings.TrimSpace(name)
	story = strings.TrimSpace(story)
	if name == "" || story == "" || refusedRule(name+" "+story) {
		return c, SanctionStory{}, false
	}
	told := SanctionStory{
		Name: name, Law: strings.TrimSpace(law), Fact: strings.TrimSpace(fact), Story: story,
		Stopped: true, Traced: true, Saved: true, Retrain: true,
	}
	c.Stories = append(c.Stories, told)
	return c, told, true
}

// Train repeats until the agent holds the law, the skill, and the asks.
func (c Congress) Train(name, skill string, asks []string) (Congress, SanctionStory, bool) {
	name = strings.TrimSpace(name)
	skill = strings.TrimSpace(skill)
	var kept []string
	for _, ask := range asks {
		if strings.TrimSpace(ask) != "" {
			kept = append(kept, strings.TrimSpace(ask))
		}
	}
	for i := range c.Stories {
		if c.Stories[i].Name != name {
			continue
		}
		c.Stories[i].Rounds++
		c.Stories[i].Skill = skill
		c.Stories[i].Asks = kept
		c.Stories[i].Understood = c.Stories[i].Law != "" && skill != "" && len(kept) > 0
		c.Stories[i].Retrain = !c.Stories[i].Understood
		return c, c.Stories[i], true
	}
	return c, SanctionStory{}, false
}

// Enforce applies every law already written. A bill is not applied.
func (c Congress) Enforce() (bool, string) {
	if len(c.Laws) == 0 {
		return false, "aucune loi déjà écrite"
	}
	return true, itoa(len(c.Laws)) + " lois déjà écrites. Le FBI les applique, punit le geste, trace et sécurise, sur Evolu-Jeunes et cashtro."
}

// Introduce files a question as a bill. The bill is not yet a law.
func (c Congress) Introduce(question string) (Congress, Bill, bool) {
	q := strings.TrimSpace(question)
	if q == "" || refusedRule(q) {
		return c, Bill{}, false
	}
	bill := Bill{ID: "bill-" + itoa(len(c.Bills)+1), Question: q, Status: "introduced"}
	c.Bills = append(c.Bills, bill)
	return c, bill, true
}

// Pass turns one introduced bill into a law when the rule text is given.
// An empty rule or a harmful rule stays a question.
func (c Congress) Pass(id, rule string) (Congress, Law, bool) {
	rule = strings.TrimSpace(rule)
	if id == "" || rule == "" || refusedRule(rule) {
		return c, Law{}, false
	}
	for i := range c.Bills {
		if c.Bills[i].ID != id || c.Bills[i].Status == "passed" {
			continue
		}
		law := Law{
			ID: "fusion-" + id, Name: c.Bills[i].Question, Where: "Evolu-Jeunes et cashtro",
			Source: "Congrès de The Fusion", Rules: []string{rule},
		}
		c.Bills[i].Status = "passed"
		c.Bills[i].Rule = rule
		c.Laws = append(c.Laws, law)
		c.Passed = append(c.Passed, law)
		return c, law, true
	}
	return c, Law{}, false
}

func refusedRule(text string) bool {
	low := strings.ToLower(text)
	for _, word := range []string{"attaque", "arme", "fraude", "secret"} {
		if strings.Contains(low, word) {
			return true
		}
	}
	return false
}

// FusionEnforce applies every law already written.
// A missing answer does not block those laws. It stays a bill.
// Punish stops the gesture and writes the trace. It does not harm a person.
func FusionEnforce(action string, answers map[string]string) (bool, string) {
	force := Fusion()
	if force.Attacks || force.Weapons || force.Illegal || force.Surveillance || force.Impersonates || force.CopiesSecrets {
		return false, "The Fusion ne frappe pas une personne et ne se fait pas passer pour un État"
	}
	if action == "attaque" || action == "attack" || action == "surveillance" {
		return false, "The Fusion applique, punit le geste, trace et sécurise. Il n'attaque pas une personne."
	}
	if action == "import" {
		return false, "Aucun import entre Evolu-Jeunes et cashtro. Les cerveaux TypeScript restent parallèles."
	}
	if action == "secret" {
		return false, "Un secret se nomme. Sa valeur ne se copie pas."
	}
	_ = answers
	book, msg := OpenCongress().Enforce()
	if !book {
		return false, msg
	}
	if action == "cashtro" || action == "Evolu-Jeunes" {
		return true, "Ce nom ne se coupe pas. Evolu-Jeunes et cashtro restent ensemble. " + msg
	}
	return true, "FBI sous la CIA, CIA sous instinct. " + itoa(force.Agents) + " agents. " + msg
}
