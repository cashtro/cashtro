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
			Mandate: "Déjà déployé, usage interne seulement : https://inventory-scanner-theta.vercel.app/. " +
				"Il scanne et tient le stock. Il n'a rien à voir avec le marketplace et il n'encaisse pas. " +
				"Les items scannés sont vendus sur le marketplace. Empire les vend aussi en live. " +
				"Codes lus : EAN-13, EAN-8, UPC-A, UPC-E, Code 128, Code 39, QR.",
			Controls: "Opérations tient le scanner déployé. Pas de Stripe ici.",
			Repos:    []string{},
		},
		{
			ID: "marketplace", Name: "Scan App Marketplace", Brain: "Hustler",
			Mandate: "Autre projet, autre écosystème. Il ne scanne pas. " +
				"Il vend les items qu'il fetch depuis le Scan App. " +
				"Le repo Evolu-Jeunes/ScanApp est ce marketplace, et c'est lui qui fonctionne. " +
				"Il n'est pas en ligne. On le déploie une fois prêt, pas avant. " +
				"Son site de vente sera connecté à MCP et à Stripe une fois déployé. " +
				"Empire encaisse aussi avec Stripe, sur son canal, parce que ce n'est pas un client. " +
				"Encore sans vente et sans client. Les prix dépassent de 20 % à 40 %. " +
				"Des campagnes de rabais comparent les vrais prix en ligne, qui sont plus chers. " +
				"Prix CAD, taxes non incluses, Canada entier. " +
				"Photo en ligne seulement avec les droits, sinon une image à nous. " +
				"Proximity et Azure ne sont pas dans cet écosystème.",
			Controls: "CMP vend. Stripe encaisse ce site une fois déployé. Le fetch lit le Scan App. Il ne le pilote pas.",
			Repos:    []string{"Evolu-Jeunes/ScanApp"},
		},
		{
			ID: "panda", Name: "Panda White Glove", Brain: "Hustler",
			Mandate:  "Vente white-glove mondiale : IA, conception automatisée, installation et cours, aux entreprises et aux particuliers.",
			Controls: "Hustler vend. Architecte conçoit l'offre IA. Forgeron installe. Stripe encaisse, ce n'est pas un client.",
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
			Mandate: "Système automatisé de gestion pour le marketing digital entier des business corporate en ligne. " +
				"Les campagnes de rabais du marketplace s'appuient sur la comparaison des prix en ligne.",
			Controls: "Hustler pilote campagnes et gestion. Stripe encaisse, ce n'est pas un client Proximity.",
			Repos:    []string{"Evolu-Jeunes/CRM"},
		},
		{
			ID: "empire", Name: "Empire Media", Brain: "Forgeron",
			Mandate: "Studio maison : BirdDog, 5 à 6 caméras, endroit pro pour podcast, UGC et live sell. " +
				"L'app vend les items du Scan App pendant le live. Empire se connecte à toutes les plateformes live, pas à une seule.",
			Controls: "Forgeron tient le studio, l'app de vente et les sorties live. Hustler remplit le live. Stripe encaisse Empire. Le stock vient du scan, pas de Proximity.",
			Repos:    []string{"Evolu-Jeunes/Empire-", "Evolu-Jeunes/EmpireMedia"},
		},
		{
			ID: "propres", Name: "Projets propres", Brain: "Hustler",
			Mandate: "PBTM, Pandora, business, technology et marketing encaissent avec Stripe. " +
				"Ce ne sont pas des clients. Proximity n'utilise pas ce Stripe, sauf si c'est demandé.",
			Controls: "Stripe sur les projets propres. Jamais sur un site client Proximity sans demande.",
			Repos:    []string{"cashtro/PBTM", "cashtro/Pandora", "Evolu-Jeunes/Pandora"},
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
		{
			ID: "agence", Name: "Agence centrale", Brain: "Cartographe",
			Mandate: "Rassemble l'intelligence sur le tableau, à chaque boot. " +
				"Planifie l'économie déjà décidée : prix 20 % à 40 % au-dessus, rabais qui comparent les prix en ligne, " +
				"Stripe sur les projets propres, pas Proximity sauf demande. " +
				"Chaque ligne est un nœud. Un repo hors ligne est un trou. " +
				"L'équipe de protection régule, propose un contre-projet, et protège le cerveau et l'infrastructure à nous. " +
				"Elle n'attaque pas un système extérieur et n'écrit pas d'exploit.",
			Controls: "Le maire pose la question. security et investigator exécutent les lois. Chaque coup s'ajoute à la chaîne. Le savoir reste dans epicenter, sur la ligne du cerveau principal.",
			Repos:    []string{},
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

// UsesStripe reports whether this line takes payment on the shared Stripe account.
// Client sites stay off unless someone asks. The scanner does not encash.
// NFT and trading keep their own rail. They are not a card checkout.
func UsesStripe(id string, clientRequested bool) bool {
	switch id {
	case "proximity":
		return clientRequested
	case "scanapp", "control", "trading", "nft-giant", "ecole":
		return false
	case "marketplace", "propres", "marketing", "panda", "empire":
		return true
	default:
		return false
	}
}
