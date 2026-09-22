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
			ID: "control", Name: "Cerveau principal", Brain: "Manager",
			Mandate:  "Synchronise les deux GitHub, les fiches, les lignes et les liens. C'est le tableau opérable, pas un neuvième produit.",
			Controls: "Manager lit state/operating.json au boot.",
			Repos:    []string{"cashtro/cashtro", "cashtro/epicenter"},
		},
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
			Mandate: "Boutique interne, aujourd'hui sans vente et sans client. Le but est de générer des ventes " +
				"en vendant les items plus cher, pour un profit. Scanner déjà capable : EAN-13, EAN-8, UPC-A, UPC-E, Code 128, Code 39, QR. " +
				"Photo trouvée en ligne seulement avec les droits ; sinon une image à nous, pas une copie. " +
				"Prix en dollars canadiens, taxes non incluses, campagne pour le Canada entier. Proximity n'est pas dans cette chaîne. " +
				"La vente live passe par Empire, sur toutes les plateformes live.",
			Controls: "CMP vend et promeut au Canada. Opérations tient le scan et la boutique interne. Contrôle vérifie les droits photo.",
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
			Mandate: "Studio maison : BirdDog, 5 à 6 caméras, endroit pro pour podcast, UGC et live sell. " +
				"L'app vend les items du Scan App pendant le live. Empire se connecte à toutes les plateformes live, pas à une seule.",
			Controls: "Forgeron tient le studio, l'app de vente et les sorties live. Hustler remplit le live. Le stock vient du scan, pas de Proximity.",
			Repos:    []string{"Evolu-Jeunes/Empire-", "Evolu-Jeunes/EmpireMedia"},
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
