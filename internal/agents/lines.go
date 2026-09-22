package agents

// Line is one business the main brain controls.
type Line struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Brain    string   `json:"brain"`
	Mandate  string   `json:"mandate"`
	Repos    []string `json:"repos"`
	Controls string   `json:"controls"`
}

// Lines is the roster the Manager uses to steer both GitHub accounts.
func Lines() []Line {
	return []Line{
		{
			ID: "proximity", Name: "Proximity Agency", Brain: "Forgeron",
			Mandate: "Usine de sites WordPress, PHP et ACF Pro pour Proximity Agency. " +
				"Couche de déploiement Azure pilotée par le cerveau principal : confection optimisée, programmation, mise en ligne.",
			Controls: "Manager assigne, Forgeron confectionne, Orfèvre déploie sur Azure, Cartographe surveille via Graphify.",
			Repos: []string{
				"Evolu-Jeunes/Proximity", "Evolu-Jeunes/ProximityApp", "Evolu-Jeunes/Proximity-Agentic",
				"Evolu-Jeunes/Api-Proximity", "Evolu-Jeunes/Proxy", "Evolu-Jeunes/Plugin",
				"Evolu-Jeunes/Btkavocat", "Evolu-Jeunes/MD-clinic", "Evolu-Jeunes/MB-Health",
				"Evolu-Jeunes/Neuro-equilibre", "Evolu-Jeunes/Hypotheque.ca", "Evolu-Jeunes/hypotheque",
				"Evolu-Jeunes/immobilier", "Evolu-Jeunes/AsselinCPA", "Evolu-Jeunes/CPA",
				"Evolu-Jeunes/Expertise", "Evolu-Jeunes/GroulxGroulx", "Evolu-Jeunes/Corps-Art",
				"Evolu-Jeunes/Alaska", "Evolu-Jeunes/Al-Soudani", "Evolu-Jeunes/AstroPoulet",
				"Evolu-Jeunes/bubbles", "Evolu-Jeunes/clic", "Evolu-Jeunes/Decoland",
				"Evolu-Jeunes/h20", "Evolu-Jeunes/h20landing", "Evolu-Jeunes/H2oH2o",
				"cashtro/H2OriginTest", "Evolu-Jeunes/noix", "Evolu-Jeunes/Noix_landing",
				"Evolu-Jeunes/lanordique", "Evolu-Jeunes/Newnordique", "Evolu-Jeunes/nordiqueweb",
				"Evolu-Jeunes/MtlEvolution", "Evolu-Jeunes/multiservices", "Evolu-Jeunes/Nhl",
				"Evolu-Jeunes/Pollo", "Evolu-Jeunes/demo-repository",
			},
		},
		{
			ID: "scanapp", Name: "Scan App Marketplace", Brain: "Hustler",
			Mandate: "Prendre la donnée du Scan App et la vendre. Système automatisé de vente, promotion, marketing et média. " +
				"Générer et trouver de vraies photos de fiche, pas les photos mal traitées de la base.",
			Controls: "Hustler vend et promeut. Forgeron améliore les fiches. Mémoire du scan reste dans le repo.",
			Repos:    []string{"Evolu-Jeunes/ScanApp"},
		},
		{
			ID: "panda", Name: "Panda White Glove", Brain: "Hustler",
			Mandate:  "Vente white-glove mondiale : IA, conception automatisée, installation et cours, aux entreprises et aux particuliers.",
			Controls: "Hustler vend. Architecte conçoit l'offre IA. Forgeron installe.",
			Repos:    []string{"Evolu-Jeunes/Panda"},
		},
		{
			ID: "nft-giant", Name: "NFT + token Giant", Brain: "Hustler",
			Mandate: "NFT pour vendre l'art en ligne avec utilité, branché sur Giant, le token de l'écosystème. " +
				"Tirages : marketing digital offert, volume, prix, partenariats de visibilité.",
			Controls: "Hustler opère art, token et tirages. Architecte relie les contrats et les chaînes.",
			Repos:    []string{"Evolu-Jeunes/Nft", "Evolu-Jeunes/Giant"},
		},
		{
			ID: "ecole", Name: "École tech", Brain: "Cartographe",
			Mandate:  "Enseigner tout ce qui est tech. Éduconnexion est le côté école.",
			Controls: "Cartographe tient le programme. Forgeron tient le site WordPress.",
			Repos:    []string{"Evolu-Jeunes/educonnexion"},
		},
		{
			ID: "marketing", Name: "Marketing digital corporate", Brain: "Hustler",
			Mandate:  "Système automatisé de gestion pour le marketing digital entier des business corporate en ligne.",
			Controls: "Hustler pilote campagnes et gestion. Cartographe mesure.",
			Repos:    []string{"Evolu-Jeunes/CRM"},
		},
		{
			ID: "empire", Name: "Empire Media", Brain: "Forgeron",
			Mandate: "Live dans une maison équipée (BirdDog, 5 à 6 caméras) : podcast, UGC, live sell. " +
				"Branché sur une app pour vendre les items pendant le live. Diffusion maison, les streamers ne dépendent pas d'une plateforme externe.",
			Controls: "Forgeron tient le studio et l'app de vente. Hustler remplit le live.",
			Repos:    []string{"Evolu-Jeunes/Empire-", "Evolu-Jeunes/EmpireMedia"},
		},
		{
			ID: "pandora", Name: "Pandora / PBTM", Brain: "Hustler",
			Mandate: "Pandora Business Technology and Marketing. Brainstorm, intern stages, Vapi voice across every bridge project. " +
				"Teal tient Teams + le desk stagiaires. Vapi parle et compose (confirm humain) sur l'infra, pas seulement Teams.",
			Controls: "Teal scrape Teams. Vapi is the kernel voice lane on every line. Hustler vend. Cursor lance depuis le desk et la carte Teams.",
			Repos:    []string{"cashtro/Pandora", "Evolu-Jeunes/Pandora", "Evolu-Jeunes/Panda"},
		},
		{
			ID: "trading", Name: "Crypto et AI bot", Brain: "Architecte",
			Mandate: "Trader et automatiser plusieurs méthodes : exchanges, blockchain, coins, gems, Web3. " +
				"Branché sur l'écosystème token Giant pour voir comment automatiser l'ensemble.",
			Controls: "Architecte modélise. Hustler suit ce qui rapporte. Giant est le pont token.",
			Repos: []string{
				"Evolu-Jeunes/AI-BOT", "Evolu-Jeunes/bot", "Evolu-Jeunes/Blockchain-Trading-",
				"Evolu-Jeunes/TradingBotCodex", "cashtro/trading_bot-main",
				"Evolu-Jeunes/sigma", "Evolu-Jeunes/sigmaNew", "Evolu-Jeunes/Giant",
			},
		},
	}
}

// LineByID returns one business line.
func LineByID(id string) (Line, bool) {
	for _, ln := range Lines() {
		if ln.ID == id {
			return ln, true
		}
	}
	return Line{}, false
}
