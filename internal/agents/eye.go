package agents

// TheEye is the defensive watch. It reads cashtro and Evolu-Jeunes
// and names a weakness. It runs its own brain, the same five-desk shape,
// and only for the watch. Epicenter Einstein stays the main brain and decides.
// The watch does not attack or copy a secret.
type Eye struct {
	Name            string   `json:"name"`
	Repo            string   `json:"repo"`
	Brain           string   `json:"brain"`
	Gate            string   `json:"gate"`
	Scope           []string `json:"scope"`
	Steps           []string `json:"steps"`
	Attacks         bool     `json:"attacks"`
	Flood           bool     `json:"flood"`
	CopiesSecrets   bool     `json:"copiesSecrets"`
	AppliesPatch    bool     `json:"appliesPatch"`
	ThroughInstinct bool     `json:"throughInstinct"`
	DecidedBy       string   `json:"decidedBy"`
	LoadsInstinct   bool     `json:"loadsInstinctBrain"`
	OwnBrain        bool     `json:"ownBrain"`
	MainBrain       string   `json:"mainBrain"`
	SameShape       bool     `json:"sameShape"`
	FunctionsOnly   bool     `json:"functionsOnly"`
	Niche           string   `json:"niche"`
	BrainSource     string   `json:"brainSource"`
	Corporation     []Desk   `json:"corporation"`
	Workers         int      `json:"workers"`
	SameBrain       bool     `json:"sameBrain"`
}

// SubBrain is one hemisphere inside a desk.
type SubBrain struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Bridges bool   `json:"bridges"`
}

// Desk is one of five. The main brain's roster sits in cashtro/epicenter.
type Desk struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Motto     string     `json:"motto"`
	Directors []string   `json:"directors"`
	Workers   int        `json:"workers"`
	SubBrains []SubBrain `json:"subBrains"`
}

// MainBrain is Epicenter Einstein's five desks. Programs do not run this roster as their own.
func MainBrain() []Desk {
	return []Desk{
		{
			ID: "architecte", Name: "L'Architecte", Motto: "Le Cerveau qui pense — le méta-cerveau",
			Directors: []string{"architect"}, Workers: 200,
			SubBrains: []SubBrain{
				{ID: "meta", Name: "Méta-Cerveau", Role: "deep-learning, self-improvement, cerveau dans le cerveau"},
				{ID: "relativite", Name: "Relativité", Role: "théorie de la relativité appliquée aux chaînes"},
				{ID: "gematria", Name: "Géomatria", Role: "géométrie sacrée et calcul symbolique"},
				{ID: "marche", Name: "Calcul du marché", Role: "calcul prédictif du marché"},
				{ID: "chaînes", Name: "Tisseur de chaînes", Role: "construit les chaînes agentics et les nœuds", Bridges: true},
			},
		},
		{
			ID: "cartographe", Name: "Le Cartographe", Motto: "Le Cerveau qui voit — la cartographie",
			Directors: []string{"planner", "research", "explorer", "memory"}, Workers: 150,
			SubBrains: []SubBrain{
				{ID: "plan", Name: "Planificateur", Role: "planification de tous les projets"},
				{ID: "map", Name: "Mappeur", Role: "mapping des projets sur la chaîne"},
				{ID: "graphify", Name: "Graphifieur", Role: "suit toute la chaîne via Graphify"},
				{ID: "monitor", Name: "Moniteur", Role: "monitorise tous les projets 24/7"},
				{ID: "philosophe", Name: "Philosophe", Role: "recherche concepts & philosophie pour améliorer la pensée et la structure", Bridges: true},
			},
		},
		{
			ID: "forgeron", Name: "Le Forgeron", Motto: "Le Cerveau qui construit — 24/7",
			Directors: []string{"operator", "deploy"}, Workers: 200,
			SubBrains: []SubBrain{
				{ID: "soudeur", Name: "Le Soudeur", Role: "frontend / UI"},
				{ID: "mecanicien", Name: "Le Mécanicien", Role: "backend / API"},
				{ID: "tisserand", Name: "Le Tisserand", Role: "intégration des chaînes"},
				{ID: "mineur", Name: "Le Mineur", Role: "données / DB"},
				{ID: "ambassadeur", Name: "L'Ambassadeur", Role: "connexion aux autres chaînes agentics", Bridges: true},
			},
		},
		{
			ID: "orfevre", Name: "L'Orfèvre", Motto: "Le Cerveau qui exécute — 24/7",
			Directors: []string{"reviewer", "delivery"}, Workers: 200,
			SubBrains: []SubBrain{
				{ID: "controleur", Name: "Le Contrôleur", Role: "QA / tests"},
				{ID: "livreur", Name: "Le Livreur", Role: "delivery / ship"},
				{ID: "polisseur", Name: "Le Polisseur", Role: "refactor / perf"},
				{ID: "documenteur", Name: "Le Documenteur", Role: "docs / structure"},
				{ID: "diplomate", Name: "Le Diplomate", Role: "connexion aux autres chaînes agentics", Bridges: true},
			},
		},
		{
			ID: "hustler", Name: "Le Hustler", Motto: "Le Cerveau qui conquiert — le gangster intelligent",
			Directors: []string{"security", "investigator", "comms", "router"}, Workers: 150,
			SubBrains: []SubBrain{
				{ID: "stratège", Name: "Le Stratège", Role: "suit tout, planification stratégique"},
				{ID: "traqueur", Name: "Le Traqueur", Role: "suit ce qui rapporte (green)"},
				{ID: "recruteur", Name: "Le Recruteur", Role: "team de hustlers"},
				{ID: "empire", Name: "Le Bâtisseur d'empire", Role: "crée de nouveaux business, automatise le tout"},
				{ID: "pontife", Name: "Le Pontife", Role: "utilise les 2 repos, toutes les chaînes, tous les nœuds, A→Z", Bridges: true},
			},
		},
	}
}

// Corporation is the main brain. A program's own desks come from NicheByID.
func Corporation() []Desk {
	return MainBrain()
}

// TheEye returns the watch roster. The working look lives in Evolu-Jeunes/Forge/brain/eye.ts.
func TheEye() Eye {
	own, _ := NicheByID("eye")
	return Eye{
		Name:            "The Eye",
		Repo:            "Evolu-Jeunes/Forge",
		Brain:           "eye",
		Gate:            "cashtro/epicenter",
		Scope:           []string{"cashtro", "Evolu-Jeunes"},
		Steps:           []string{"veille", "menace", "faille", "recherche", "garde", "instinct", "rustine", "preuve", "sceau"},
		Attacks:         false,
		Flood:           false,
		CopiesSecrets:   false,
		AppliesPatch:    false,
		ThroughInstinct: true,
		DecidedBy:       "instinct",
		LoadsInstinct:   false,
		OwnBrain:        true,
		MainBrain:       "instinct",
		SameShape:       true,
		FunctionsOnly:   true,
		Niche:           own.Niche,
		BrainSource:     "Evolu-Jeunes/Forge",
		Corporation:     own.Desks,
		Workers:         own.Workers,
		SameBrain:       false,
	}
}

// EyeChange opens a draft only after Epicenter Einstein. The draft is not applied.
func EyeChange(via, seat, caller string) (bool, string) {
	if caller != "voltron" {
		return false, "une fonction optionnelle passe par Voltron"
	}
	if via != "instinct" && via != "cashtro/epicenter" {
		return false, "The Eye passe par Epicenter Einstein avant tout changement."
	}
	if seat != "instinct" {
		return false, "Sans Epicenter Einstein, aucun changement."
	}
	return true, "brouillon ouvert. rien n'est appliqué"
}
