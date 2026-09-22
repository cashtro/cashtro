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
			ID: "scanapp", Name: "Scan App", Brain: "Forgeron",
			Mandate: "Projet déployé : la fonction scan. https://inventory-scanner-theta.vercel.app/ " +
				"gère les produits, scanne les codes et met le stock à jour. " +
				"Codes lus : EAN-13, EAN-8, UPC-A, UPC-E, Code 128, Code 39, QR. " +
				"Ce n'est pas le marketplace. Empire vend aussi ces produits en live.",
			Controls: "Opérations tient le scanner déployé. Empire vend ce stock en live.",
			Repos:    []string{},
		},
		{
			ID: "marketplace", Name: "Scan App Marketplace", Brain: "Hustler",
			Mandate: "Autre projet, autre écosystème. Il ne scanne pas. " +
				"Il vend les items qu'il fetch depuis le Scan App. " +
				"Le repo Evolu-Jeunes/ScanApp est ce marketplace, et c'est lui qui fonctionne. " +
				"Il n'est pas en ligne. On le déploie une fois prêt, pas avant. " +
				"Son site de vente sera connecté à MCP et à Stripe. Pas le live Empire. " +
				"Encore sans vente et sans client. But : vendre plus cher, profit. " +
				"Prix CAD, taxes non incluses, Canada entier. " +
				"Photo en ligne seulement avec les droits, sinon une image à nous. " +
				"Proximity et Azure ne sont pas dans cet écosystème.",
			Controls: "CMP vend. Le fetch lit le Scan App. Il ne le pilote pas.",
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
