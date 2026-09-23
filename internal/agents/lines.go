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
			Mandate:  "Synchronise les deux GitHub, les fiches, les lignes et les liens. Voltron le tient allumé. Graphify est la carte. C'est le tableau opérable, pas un neuvième produit.",
			Controls: "Manager lit state/operating.json au boot. Explorer lit le graphe. Memory ajoute un maillon.",
			Repos:    []string{"cashtro/cashtro", "cashtro/epicenter"},
		},
		{
			ID: "wordpress", Name: "WordPress Proximity", Brain: "Forgeron",
			Mandate: "Epicenter et Voltron dirigent ce département. Le commit, lui, ne part que d'ici. " +
				"Sites WordPress Proximity : PHP, ACF Pro, thèmes, XAMPP. " +
				"Et les produits à connaître et à committer : Proximity, Next Proximity (apps/web), Proximity App, Api-Proximity. " +
				"Proximity App est le WordPress de l'agence. Lovable n'entre pas.",
			Controls: "Operator confectionne. Security filtre. Reviewer filtre encore. Memory garde le contexte. Rien en ligne sans comms.allow. Aucun commit sur Epicenter, Voltron, ou une autre marque.",
			Repos: []string{
				"Evolu-Jeunes/immobilier", "Evolu-Jeunes/hypotheque", "Evolu-Jeunes/Hypotheque.ca",
				"Evolu-Jeunes/MD-clinic",
				"Evolu-Jeunes/AsselinCPA", "Evolu-Jeunes/GroulxGroulx", "Evolu-Jeunes/Corps-Art",
				"Evolu-Jeunes/AstroPoulet", "Evolu-Jeunes/Pollo", "Evolu-Jeunes/Decoland",
				"Evolu-Jeunes/bubbles", "Evolu-Jeunes/clic", "Evolu-Jeunes/multiservices",
				"Evolu-Jeunes/lanordique", "Evolu-Jeunes/Newnordique",
				"Evolu-Jeunes/h20", "Evolu-Jeunes/h20landing", "Evolu-Jeunes/H2oH2o", "cashtro/H2OriginTest",
				"Evolu-Jeunes/Proximity", "Evolu-Jeunes/ProximityApp", "Evolu-Jeunes/Api-Proximity",
				"Evolu-Jeunes/Proxy", "Evolu-Jeunes/Plugin",
				"Evolu-Jeunes/educonnexion", "Evolu-Jeunes/sigma", "Evolu-Jeunes/sigmaNew",
			},
		},
		{
			ID: "proximity", Name: "Proximity Agency", Brain: "Forgeron",
			Mandate: "L'autre département de la même maison. " +
				"Comptes qui ne sont pas des sites WordPress, ni Proximity, ni Next Proximity, ni Proximity App, ni l'API. " +
				"Ce département ne commit pas sur PHP, ACF Pro, ni XAMPP.",
			Controls: "L'autre équipe empile ces comptes. Le kernel ne les mélange pas avec le département WordPress.",
			Repos: []string{
				"Evolu-Jeunes/Proximity-Agentic",
				"Evolu-Jeunes/MB-Health", "Evolu-Jeunes/Btkavocat",
				"Evolu-Jeunes/Neuro-equilibre", "Evolu-Jeunes/CPA", "Evolu-Jeunes/Expertise",
				"Evolu-Jeunes/Alaska", "Evolu-Jeunes/Al-Soudani", "Evolu-Jeunes/MtlEvolution",
				"Evolu-Jeunes/noix", "Evolu-Jeunes/Noix_landing", "Evolu-Jeunes/nordiqueweb",
				"Evolu-Jeunes/Nhl", "Evolu-Jeunes/demo-repository",
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
			Mandate: "Giant a son propre cerveau crypto. Voltron ouvre le cycle, Giant le fait tourner, Epicenter décide en dernier. " +
				"NFT pour l'art utile, token GNT, lecture de marché, gems en note, modèles de bots. Aucun ordre live.",
			Controls: "Le cerveau de Giant tient les agents crypto. Epicenter est le dernier mot. Pas de mint, pas d'ordre, pas de pont.",
			Repos:    []string{"Evolu-Jeunes/Nft", "Evolu-Jeunes/Giant"},
		},
		{
			ID: "ecole", Name: "École tech", Brain: "Cartographe",
			Mandate:  "Enseigner tout ce qui est tech. Éduconnexion est le côté école.",
			Controls: "Cartographe tient le programme. Le thème Éduconnexion est au département WordPress.",
			Repos:    []string{"Evolu-Jeunes/educonnexion"},
		},
		{
			ID: "marketing", Name: "Marketing digital corporate", Brain: "Hustler",
			Mandate: "Le seul CRM Lovable est Evolu-Jeunes/CRM (Itercore) : contacts, leads, notes. " +
				"Pandora mesure une campagne et la range dans ce carnet. Panda et le nom d'un client Proximity peuvent y déposer une note. " +
				"La copie dans Proximity n'est pas un second CRM. Lovable n'édite pas les thèmes.",
			Controls: "Explorer lit. Memory dépose un brouillon. Epicenter décide avant un envoi. Stripe encaisse, ce n'est pas un client Proximity.",
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
				"Evolu-Jeunes/Giant",
			},
		},
		{
			ID: "fonds", Name: "Fonds", Brain: "Manager",
			Mandate: "Grand livre interne, chèques en brouillon, TPS et TVQ. " +
				"Ce n'est pas une banque et ce n'est pas une déclaration. Rien ne sort sans comms.allow.",
			Controls: "Planner tient la position. Security filtre. Reviewer garde le calcul le plus court. Le chèque reste un brouillon.",
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
	case "scanapp", "control", "trading", "nft-giant", "ecole", "fonds":
		return false
	case "marketplace", "propres", "marketing", "panda", "empire":
		return true
	default:
		return false
	}
}
