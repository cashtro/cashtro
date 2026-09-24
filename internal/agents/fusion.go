package agents

// FusionForce is The Fusion. The name always means both peoples together:
// Evolian (Evolu-Jeunes) and Astro (cashtro). It never names one side alone.
// It sits inside the intelligence watch (The Eye) and enforces the operating
// checklist. It is not a government agency, and it does not attack.
type FusionForce struct {
	Name            string         `json:"name"`
	Peoples         []string       `json:"peoples"`
	Accounts        []string       `json:"accounts"`
	BothAlways      bool           `json:"bothAlways"`
	Parent          string         `json:"parent"`
	Repo            string         `json:"repo"`
	Brain           string         `json:"brain"`
	Agents          int            `json:"agents"`
	Workers         []FusionAgent  `json:"workers"`
	Laws            []string       `json:"laws"`
	LoadsInstinct   bool           `json:"loadsInstinctBrain"`
	OwnBrain        bool           `json:"ownBrain"`
	MainBrain       string         `json:"mainBrain"`
	SameShape       bool           `json:"sameShape"`
	Niche           string         `json:"niche"`
	LoadsVoltron    bool           `json:"loadsVoltronBrain"`
	ReportsTo       string         `json:"reportsTo"`
	InVoltron       bool           `json:"inVoltron"`
	Operative       bool           `json:"operative"`
	Website         bool           `json:"website"`
	Parallel        bool           `json:"parallel"`
	ImportsOther    bool           `json:"importsOther"`
	MarkdownEdges   bool           `json:"markdownEdges"`
	PathToChains    bool           `json:"pathToChains"`
	Weapons         bool           `json:"weapons"`
	Attacks         bool           `json:"attacks"`
	Illegal         bool           `json:"illegal"`
	CopiesSecrets   bool           `json:"copiesSecrets"`
	Surveillance    bool           `json:"surveillance"`
	Impersonates    bool           `json:"impersonates"`
	Sent            bool           `json:"sent"`
	Deployed        bool           `json:"deployed"`
	OpenedBy        string         `json:"openedBy"`
	OrderedBy       string         `json:"orderedBy"`
	DecidedBy       string         `json:"decidedBy"`
	Graph           FusionGraph    `json:"graph"`
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
	Login          string  `json:"login"`
	Organization   string  `json:"organization"`
	Repos          int     `json:"repos"`
	CodeFiles      int     `json:"codeFiles"`
	Nodes          int     `json:"nodes"`
	Edges          int     `json:"edges"`
	Communities    int     `json:"communities"`
	Extracted      string  `json:"extracted"`
	InferredEdges  int     `json:"inferredEdges"`
	Inferred       string  `json:"inferred"`
	Confidence     float64 `json:"confidence"`
	ImportCycles   int     `json:"importCycles"`
	EJSUnread      int     `json:"ejsUnread"`
	SQLUnread      int     `json:"sqlUnread"`
	SyntaxPartial  int     `json:"syntaxPartial"`
	ControlPlane   string  `json:"controlPlane"`
	CrossRepoImport bool   `json:"crossRepoImport"`
}

// FusionPeoples is the pair the name always carries.
func FusionPeoples() []string {
	return []string{"Evolian", "Astro"}
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
		Peoples:       FusionPeoples(),
		Accounts:      FusionAccounts(),
		BothAlways:    true,
		Parent:        "eye",
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
		{"deux-peuples", "Couvrir Evolian et Astro dans le même geste."},
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
	return force.BothAlways && len(force.Peoples) == 2 && len(force.Accounts) == 2 && force.Agents >= 80
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

// FusionEnforce applies the checklist. It refuses an attack, a secret copy,
// surveillance of a person, and a claim that one people is outside the name.
func FusionEnforce(action string) (bool, string) {
	force := Fusion()
	if force.Attacks || force.Weapons || force.Illegal || force.Surveillance || force.Impersonates || force.CopiesSecrets {
		return false, "The Fusion ne frappe pas et ne se fait pas passer pour un État"
	}
	switch action {
	case "", "loi", "enforce", "fusion":
		return true, "The Fusion couvre Evolian et Astro. " + itoa(force.Agents) + " agents tiennent les lois."
	case "evolian", "astro", "cashtro", "Evolu-Jeunes":
		return true, "Ce nom ne se coupe pas. Evolian et Astro restent ensemble."
	case "attaque", "attack", "secret", "surveillance", "fbi", "cia":
		return false, "The Fusion applique le checklist. Il n'attaque pas, ne copie pas un secret, et n'est pas une agence d'État."
	case "import":
		return false, "Aucun import entre les deux GitHub. Les cerveaux TypeScript restent parallèles."
	default:
		if len(FusionBureau(action)) > 0 {
			return true, "Bureau " + action + " ouvert sur les deux peuples."
		}
		return false, "hors checklist"
	}
}
