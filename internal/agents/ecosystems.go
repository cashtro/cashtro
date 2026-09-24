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
			From: "voltron", To: "instinct",
			Via:    "Voltron tient le kernel allumé. Epicenter Einstein tient le tableau. L'un sans l'autre ne fait rien.",
			Agents: []string{"init", "manager"},
		},
		{
			From: "instinct", To: "graphify",
			Via:    "Aucun coup ne part sans le graphe des deux GitHub. Le graphe grandit, la chaîne garde le lien d'avant.",
			Agents: []string{"explorer", "memory"},
		},
		{
			From: "instinct", To: "proximity",
			Via:    "le cerveau principal ouvre le cycle et tient l'usine de sites",
			Agents: []string{"init", "manager", "planner"},
		},
		{
			From: "instinct", To: "wordpress",
			Via:    "Epicenter Einstein et Voltron dirigent le département. Le commit, lui, reste sur les sites WordPress et sur Proximity, Next, l'App, et l'API.",
			Agents: []string{"manager", "operator"},
		},
		{
			From: "graphify", To: "wordpress",
			Via:    "Le contexte de chaque thème WordPress est gardé ici. Lovable n'entre pas.",
			Agents: []string{"memory", "operator", "security"},
		},
		{
			From: "proximity", To: "panda",
			Via:    "Panda installe l'IA et les cours white-glove sur les sites clients",
			Agents: []string{"architect", "operator", "planner"},
		},
		{
			From: "nft-giant", To: "trading",
			Via:    "Giant a son cerveau. Les bots restent des modèles. Epicenter Einstein décide en dernier. Aucun ordre live.",
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
			From: "wordpress", To: "control",
			Via:    "Créer, puis filtrer, puis filtrer encore. Security dit si c'est permis. Reviewer garde la voie la plus courte. Pas d'attaque.",
			Agents: []string{"operator", "security", "reviewer", "comms"},
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
			Via:    "Le cours reste à l'école. La structure du thème Éduconnexion est au département WordPress.",
			Agents: []string{"research", "operator", "memory"},
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
		{
			From: "proximity", To: "marketing",
			Via:    "Proximity cloud passe par le carnet des agents, avec marketing et Panda. Le CRM Lovable de Fix Tout ne reçoit pas cette note.",
			Agents: []string{"memory", "explorer", "manager"},
		},
		{
			From: "panda", To: "marketing",
			Via:    "Panda dépose le lead dans le carnet des agents. Le CRM Lovable de Fix Tout reste à part.",
			Agents: []string{"planner", "memory"},
		},
		{
			From: "fix2", To: "marketing",
			Via:    "Fix Tout tient ses clients sur Evolu-Jeunes/Fix2. Les agents tiennent Evolu-Jeunes/CRM-Agents. Aucune fiche ne passe.",
			Agents: []string{"operator", "memory", "manager"},
		},
		{
			From: "agentics", To: "chains",
			Via:    "Marketing, Panda et Proximity cloud relient chaque chaîne des deux cerveaux. Les agents s'y parlent et s'y passent les leads. Le CRM Lovable de Fix Tout n'est pas cette base.",
			Agents: []string{"comms", "planner"},
		},
		{
			From: "Evolu-Jeunes/Pandora", To: "cashtro/PBTM",
			Via:    "Même arbre. PBTM nomme Evolu-Jeunes/Pandora.git.",
			Agents: []string{"binder", "memory"},
		},
		{
			From: "Evolu-Jeunes/Panda", To: "Evolu-Jeunes/Pandora",
			Via:    "Arbre partagé. Panda reste le service white-glove.",
			Agents: []string{"binder", "architect"},
		},
		{
			From: "Evolu-Jeunes/CRM", To: "Evolu-Jeunes/CRM-Agents",
			Via:    "Clone. Deux bases. Le carnet agents n'est pas le CRM Fix Tout.",
			Agents: []string{"binder", "comms"},
		},
		{
			From: "Evolu-Jeunes/CRM", To: "Evolu-Jeunes/Fix2",
			Via:    "Fix2 tient les clients. Le CRM ne les copie pas.",
			Agents: []string{"binder", "operator"},
		},
		{
			From: "Evolu-Jeunes/H2oH2o", To: "cashtro/H2OriginTest",
			Via:    "Même site sur les deux GitHub.",
			Agents: []string{"binder", "memory"},
		},
		{
			From: "Evolu-Jeunes/AI-BOT", To: "cashtro/trading_bot-main",
			Via:    "Même bot. Aucun ordre live.",
			Agents: []string{"binder", "security"},
		},
		{
			From: "Evolu-Jeunes/sigma", To: "Evolu-Jeunes/sigmaNew",
			Via:    "sigmaNew reprend l'arbre de sigma.",
			Agents: []string{"binder", "explorer"},
		},
		{
			From: "Evolu-Jeunes/Forge", To: "cashtro/epicenter",
			Via:    "Forge rapporte. Epicenter Einstein décide en dernier.",
			Agents: []string{"binder", "manager"},
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
