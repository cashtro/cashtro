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
			From: "epicenter", To: "proximity",
			Via:    "le cerveau principal ouvre le cycle et tient l'usine de sites",
			Agents: []string{"init", "manager", "planner"},
		},
		{
			From: "proximity", To: "scanapp",
			Via:    "les sites Proximity publient les fiches du Scan App, photos propres comprises",
			Agents: []string{"operator", "deploy", "reviewer", "delivery"},
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
			Via:    "le live sell d'Empire vend les items du Scan App",
			Agents: []string{"operator", "comms", "delivery", "reviewer"},
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
