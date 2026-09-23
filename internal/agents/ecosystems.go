package agents

// Link is one connection inside or between business ecosystems.
// Agents are the kernel agentics that carry the link.
type Link struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Via    string   `json:"via"`
	Agents []string `json:"agents"`
}

// Links is the automation graph. Every kernel agentic appears at least once.
func Links() []Link {
	return []Link{
		{
			From: "voltron", To: "epicenter",
			Via:    "Voltron tient le kernel allumé. Epicenter tient le tableau. L'un sans l'autre ne fait rien.",
			Agents: []string{"init", "manager"},
		},
		{
			From: "epicenter", To: "graphify",
			Via:    "Aucun coup ne part sans le graphe des deux GitHub. Le graphe grandit, la chaîne garde le lien d'avant.",
			Agents: []string{"explorer", "memory"},
		},
		{
			From: "epicenter", To: "proximity",
			Via:    "le cerveau principal ouvre le cycle et tient l'usine de sites",
			Agents: []string{"init", "manager", "planner"},
		},
		{
			From: "graphify", To: "proximity",
			Via:    "Chaque site est un contexte gardé. Une demande client est le coup suivant sur ce contexte, pas un site nouveau.",
			Agents: []string{"memory", "planner", "security"},
		},
		{
			From: "lovable", To: "proximity",
			Via:    "Un site neuf naît dans Lovable. Son dépôt entre ensuite dans la chaîne Proximity. Lovable ne touche pas un WordPress déjà en ligne.",
			Agents: []string{"architect", "delivery"},
		},
		{
			From: "proximity", To: "panda",
			Via:    "Panda installe l'IA et les cours white-glove sur les sites clients",
			Agents: []string{"architect", "operator", "planner"},
		},
		{
			From: "nft-giant", To: "trading",
			Via:    "le token Giant relie l'art NFT aux bots qui tradent",
			Agents: []string{"architect", "security", "investigator", "router"},
		},
		{
			From: "empire", To: "scanapp",
			Via:    "le live Empire vend le stock Scan App sur toutes les plateformes live. Proximity n'est pas dans cette chaîne.",
			Agents: []string{"operator", "comms", "delivery", "reviewer", "deploy"},
		},
		{
			From: "propres", To: "marketing",
			Via:    "Pandora est la maison marketing. Elle ne poste pas toute seule. Une campagne mesurée, puis Facebook, après comms.allow.",
			Agents: []string{"planner", "comms", "explorer"},
		},
		{
			From: "marketing", To: "proximity",
			Via:    "le marketing corporate promeut les sites et mesure via Graphify",
			Agents: []string{"comms", "planner", "explorer"},
		},
		{
			From: "proximity", To: "control",
			Via:    "Un site live ne change qu'en preview, revue, puis accord. La sécurité lit le graphe. Elle n'attaque pas.",
			Agents: []string{"security", "reviewer", "comms", "deploy"},
		},
		{
			From: "marketing", To: "empire",
			Via:    "le média du live part dans la gestion marketing",
			Agents: []string{"comms", "explorer"},
		},
		{
			From: "marketing", To: "panda",
			Via:    "les campagnes vendent le service white-glove",
			Agents: []string{"comms", "research"},
		},
		{
			From: "ecole", To: "proximity",
			Via:    "Éduconnexion enseigne WordPress, PHP et ACF Pro",
			Agents: []string{"research", "memory", "planner"},
		},
		{
			From: "ecole", To: "trading",
			Via:    "l'école enseigne le Web3 que les bots exécutent",
			Agents: []string{"research", "memory", "architect"},
		},
		{
			From: "marketplace", To: "scanapp",
			Via:    "le marketplace fetch les items du Scan App pour les vendre. Deux projets, deux écosystèmes. Il ne scanne pas.",
			Agents: []string{"delivery", "comms"},
		},
		{
			From: "azure", To: "proximity",
			Via:    "l'équipe Azure gère le système branché, en même temps que le reste. Pas le Scan App, pas le marketplace.",
			Agents: []string{"deploy", "operator"},
		},
		{
			From: "optimisation", To: "direction",
			Via:    "Optimisation contredit chaque proposition et ne laisse passer que l'option la plus courte",
			Agents: []string{"reviewer", "explorer", "manager"},
		},
	}
}

// AgentIDs returns every kernel agentic used by the links.
func AgentIDs() []string {
	seen := map[string]bool{}
	var out []string
	for _, ln := range Links() {
		for _, id := range ln.Agents {
			if seen[id] {
				continue
			}
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}
