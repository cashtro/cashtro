package agents

// Link is one connection inside or between business ecosystems.
// Agents are the kernel agentics that carry the link.
type Link struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Via    string   `json:"via"`
	Agents []string `json:"agents"`
}

func withVoice(agents ...string) []string {
	out := make([]string, 0, len(agents)+1)
	out = append(out, "vapi")
	out = append(out, agents...)
	return out
}

// Links is the automation graph. Every kernel agentic appears at least once.
// Vapi rides every hop so talk and outbound calls share line context.
func Links() []Link {
	return []Link{
		{
			From: "epicenter", To: "proximity",
			Via:    "le cerveau principal ouvre le cycle et tient l'usine de sites",
			Agents: withVoice("init", "manager", "planner"),
		},
		{
			From: "proximity", To: "scanapp",
			Via:    "les sites Proximity publient les fiches du Scan App, photos propres comprises",
			Agents: withVoice("operator", "deploy", "reviewer", "delivery"),
		},
		{
			From: "proximity", To: "panda",
			Via:    "Panda installe l'IA et les cours white-glove sur les sites clients",
			Agents: withVoice("architect", "operator", "planner"),
		},
		{
			From: "nft-giant", To: "trading",
			Via:    "le token Giant relie l'art NFT aux bots qui tradent",
			Agents: withVoice("architect", "security", "investigator", "router"),
		},
		{
			From: "empire", To: "scanapp",
			Via:    "le live sell d'Empire vend les items du Scan App",
			Agents: withVoice("operator", "comms", "delivery", "reviewer"),
		},
		{
			From: "marketing", To: "proximity",
			Via:    "le marketing corporate promeut les sites et mesure via Graphify",
			Agents: withVoice("comms", "planner", "explorer"),
		},
		{
			From: "marketing", To: "empire",
			Via:    "le média du live part dans la gestion marketing",
			Agents: withVoice("comms", "explorer"),
		},
		{
			From: "marketing", To: "panda",
			Via:    "les campagnes vendent le service white-glove",
			Agents: withVoice("comms", "research"),
		},
		{
			From: "ecole", To: "proximity",
			Via:    "Éduconnexion enseigne WordPress, PHP et ACF Pro",
			Agents: withVoice("research", "memory", "planner"),
		},
		{
			From: "ecole", To: "trading",
			Via:    "l'école enseigne le Web3 que les bots exécutent",
			Agents: withVoice("research", "memory", "architect"),
		},
		{
			From: "pandora", To: "ecole",
			Via:    "Teal + Vapi: le cerveau écoute, les stagiaires posent, la voix parle sur l'infra",
			Agents: withVoice("teal", "memory", "comms", "planner"),
		},
		{
			From: "pandora", To: "proximity",
			Via:    "Vapi talks factory tickets and intern pairing across Proximity sites",
			Agents: withVoice("teal", "manager", "delivery"),
		},
		{
			From: "pandora", To: "scanapp",
			Via:    "Vapi sells Scan App fiches on the call (human-gated outbound)",
			Agents: withVoice("comms", "delivery"),
		},
		{
			From: "pandora", To: "panda",
			Via:    "Vapi coaches white-glove installs on the Panda line",
			Agents: withVoice("architect", "planner"),
		},
		{
			From: "pandora", To: "nft-giant",
			Via:    "Vapi narrates NFT + Giant token drops",
			Agents: withVoice("architect", "security"),
		},
		{
			From: "pandora", To: "marketing",
			Via:    "Vapi runs corporate campaign callbacks",
			Agents: withVoice("comms", "explorer"),
		},
		{
			From: "pandora", To: "empire",
			Via:    "Vapi sits in the Empire live-sell booth",
			Agents: withVoice("operator", "delivery"),
		},
		{
			From: "pandora", To: "trading",
			Via:    "Vapi briefs Architecte bots and Giant token moves",
			Agents: withVoice("router", "investigator"),
		},
		{
			From: "epicenter", To: "teal",
			Via:    "Cashtro manager publishes the fresh Evolu-Jeunes/Teal voice repo across the fleet",
			Agents: withVoice("init", "manager", "teal"),
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
