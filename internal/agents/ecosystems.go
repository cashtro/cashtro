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
// The agency watch is appended, one node per line.
func Links() []Link {
	return append(baseLinks(), WatchLinks()...)
}

func baseLinks() []Link {
	return []Link{
		{
			From: "epicenter", To: "proximity",
			Via:    "le cerveau principal ouvre le cycle et tient l'usine de sites",
			Agents: []string{"init", "manager", "planner"},
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
			From: "marketing", To: "proximity",
			Via:    "le marketing corporate promeut les sites et mesure via Graphify",
			Agents: []string{"comms", "planner", "explorer"},
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

// WatchLinks ties the agency to every other line, so a project cannot sit
// off the board. The link is a node in the chain, not a new process.
func WatchLinks() []Link {
	out := make([]Link, 0, len(Lines()))
	for _, ln := range Lines() {
		if ln.ID == "agence" {
			continue
		}
		out = append(out, Link{
			From:   "agence",
			To:     ln.ID,
			Via:    "l'agence lit cette ligne à chaque boot. Rien ne sort du tableau.",
			Agents: []string{"memory", "security", "investigator"},
		})
	}
	return out
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
