package agents

// LoadedBrain is one brain the conglomerate actually runs.
// Ask, Chain, Contradict, Improve, and the ledger all see it.
type LoadedBrain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Seat string `json:"seat"`
	Dept string `json:"dept"`
	Do   string `json:"do"`
}

// Brains is the set that must be analyzed, not only named.
// Epicenter Einstein's five, Voltron, Scrum, Forge, The Eye, Giant, the three OPS chains, and The Hustle.
func Brains() []LoadedBrain {
	return []LoadedBrain{
		{ID: "instinct", Name: "Epicenter Einstein", Seat: "instinct", Dept: "direction", Do: "Charger les cinq cerveaux et trancher en dernier."},
		{ID: "architecte", Name: "L'Architecte", Seat: "instinct", Dept: "ingenierie", Do: "Penser le méta-cerveau. Deux façons, puis la plus courte."},
		{ID: "cartographe", Name: "Le Cartographe", Seat: "instinct", Dept: "direction", Do: "Voir chaque dépôt sur Graphify avant le coup."},
		{ID: "forgeron", Name: "Le Forgeron", Seat: "instinct", Dept: "operations", Do: "Construire le brouillon. Pas une arme."},
		{ID: "orfevre", Name: "L'Orfèvre", Seat: "instinct", Dept: "controle", Do: "Exécuter le contrôle. Citer la règle."},
		{ID: "hustler", Name: "Le Hustler", Seat: "instinct", Dept: "marche", Do: "Ouvrir The Hustle. Une offre, une mesure. Il ne décide pas à la place d'Epicenter Einstein."},
		{ID: "voltron", Name: "Voltron", Seat: "voltron", Dept: "flux", Do: "Tenir chaque chaîne agentique allumée. Elles rapportent ici. Aucune n'est un site."},
		{ID: "scrum", Name: "Scrum", Seat: "scrum", Dept: "direction", Do: "Prendre chaque dépôt et chaque projet nouveau."},
		{ID: "forge", Name: "Forge", Seat: "forge", Dept: "ingenierie", Do: "Faire tourner le conglomérat de développement. Pas une arme."},
		{ID: "eye", Name: "The Eye", Seat: "eye", Dept: "controle", Do: "Nommer une faiblesse avec son propre cerveau de veille. Epicenter Einstein reste le cerveau principal et tranche. Ne pas attaquer."},
		{ID: "giant", Name: "Giant", Seat: "giant", Dept: "controle", Do: "Tenir la formation crypto. Le mint n'est pas une option."},
		{ID: "accueil", Name: "Accueil", Seat: "ops", Dept: "marche", Do: "Qualifier le lead. Ne pas appeler et ne pas envoyer de SMS."},
		{ID: "terrain", Name: "Terrain", Seat: "ops", Dept: "operations", Do: "Tenir soumission, chantier, tournée et facture en brouillon."},
		{ID: "propriete", Name: "Propriété", Seat: "ops", Dept: "ingenierie", Do: "Posséder OPS et aider les dépôts cashtro et Evolu-Jeunes."},
		{ID: "hustle", Name: "The Hustle", Seat: "instinct", Dept: "marche", Do: "Lire, tenir les deux métaux, et ne garder que le coup légal. Une offre, une mesure. Epicenter Einstein est l'ami."},
	}
}

// BrainByID returns one loaded brain.
func BrainByID(id string) (LoadedBrain, bool) {
	for _, b := range Brains() {
		if b.ID == id {
			return b, true
		}
	}
	return LoadedBrain{}, false
}

func brainKnown(id string) (map[string]string, bool) {
	b, ok := BrainByID(id)
	if !ok {
		return nil, false
	}
	return map[string]string{
		"situation":   b.Name + " est chargé. " + b.Do,
		"manque":      "La décision d'Epicenter Einstein avant un envoi, un paiement, un appel, ou un déploiement.",
		"intouchable": "Un site client déjà en ligne, un secret, et une arme.",
	}, true
}

func brainChain(id string) ([]ChainStep, bool) {
	b, ok := BrainByID(id)
	if !ok {
		return nil, false
	}
	return []ChainStep{
		{Line: id, Order: 1, Agent: "memory", Do: b.Do},
		{Line: id, Order: 2, Agent: "manager", Do: "Epicenter Einstein lit ce cerveau et décide. Le coup le plus court reste."},
	}, true
}

// Analysis is one brain, or one business line, after ask, chain, contradict, improve, and the ledger.
type Analysis struct {
	ID        string `json:"id"`
	Open      int    `json:"open"`
	Steps     int    `json:"steps"`
	Best      string `json:"best"`
	Dept      string `json:"dept"`
	Improved  bool   `json:"improved"`
	Intact    bool   `json:"intact"`
	Connected bool   `json:"connected"`
}

func deptFor(id string) string {
	if b, ok := BrainByID(id); ok {
		return b.Dept
	}
	for _, d := range Chart().Divisions {
		if d.ID == id {
			return d.Department
		}
	}
	return "direction"
}

// Analyze runs the kernel path on one id. Connected is false when any step is missing.
func Analyze(id string) Analysis {
	self := AskSelf(id, nil)
	steps, chained := Chain(id)
	verdict := Contradict(Proposal{
		Subject: id,
		Options: []Option{
			{Name: "rallonger", Cost: 5, Risk: 3, Steps: 4},
			{Name: "le plus court", Cost: 1, Risk: 1, Steps: 1},
		},
	})
	dept := deptFor(id)
	draft := id
	if chained && len(steps) > 0 {
		draft = steps[0].Do
	}
	improved, improvedOK := Improve(dept, draft)
	question := "Quelle est la position sur le tableau, avant ce coup?"
	if len(self.Asked) > 0 {
		question = self.Asked[0].Ask
	}
	moves := AppendMove(nil, "ask:"+id, question, self.Answered["position"])
	coup := draft
	if chained && len(steps) > 1 {
		coup = steps[1].Do
	}
	moves = AppendMove(moves, "chain:"+id, question, coup)
	moves = AppendMove(moves, "contradict:"+id, question, verdict.Best)
	moves = AppendMove(moves, "improve:"+id, question, improved.Result)
	intact := ChainIntact(moves)
	connected := len(self.Open) == 0 && chained && len(steps) >= 2 && verdict.Ready && verdict.Best == "le plus court" && improvedOK && intact
	return Analysis{
		ID:        id,
		Open:      len(self.Open),
		Steps:     len(steps),
		Best:      verdict.Best,
		Dept:      dept,
		Improved:  improvedOK && improved.Department == dept,
		Intact:    intact,
		Connected: connected,
	}
}

// Subjects is every line, every loaded brain, and the watch question.
// A line id is listed once even when a chain action shares that id.
func Subjects() []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(Lines())+len(Brains())+1)
	add := func(id string) {
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		out = append(out, id)
	}
	for _, ln := range Lines() {
		add(ln.ID)
	}
	for _, b := range Brains() {
		add(b.ID)
	}
	add("watch")
	return out
}

// Disconnected lists subjects whose ask, chain, contradict, improve, or ledger step is missing.
func Disconnected() []string {
	var out []string
	for _, a := range AnalyzeAll() {
		if !a.Connected {
			out = append(out, a.ID)
		}
	}
	return out
}

// AnalyzeAll runs Analyze on every subject.
func AnalyzeAll() []Analysis {
	subjects := Subjects()
	out := make([]Analysis, 0, len(subjects))
	for _, id := range subjects {
		out = append(out, Analyze(id))
	}
	return out
}
