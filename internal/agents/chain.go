package agents

// ChainStep is one piece of work the main brain hands to an agent.
type ChainStep struct {
	Line  string `json:"line"`
	Order int    `json:"order"`
	Do    string `json:"do"`
	Agent string `json:"agent"`
}

// Fiche is a sales listing built from a scanned product.
// The database image is never the public photo.
type Fiche struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Stock   int    `json:"stock"`
	Photo   string `json:"photo"` // generated, licensed, or database
	Person  bool   `json:"person"`
	Consent bool   `json:"consent"`
}

// Ready reports whether this listing can be published or sold.
func (f Fiche) Ready() (bool, string) {
	if f.Code == "" || f.Stock < 1 {
		return false, "pas de stock scanné"
	}
	switch f.Photo {
	case "generated", "licensed":
	default:
		return false, "l'image de la base n'est pas la fiche"
	}
	if f.Person && !f.Consent {
		return false, "personne identifiable sans consentement"
	}
	if f.Name == "" || f.Price == "" {
		return false, "fiche incomplète"
	}
	return true, "fiche prête"
}

// Chain returns the ordered work for one business line.
func Chain(line string) ([]ChainStep, bool) {
	switch line {
	case "scanapp":
		return []ChainStep{
			{Line: "scanapp", Order: 1, Agent: "operator", Do: "Lire le stock sur https://inventory-scanner-theta.vercel.app/. Le repo ne fait pas le scan, il tient le marketplace."},
			{Line: "scanapp", Order: 2, Agent: "operator", Do: "Écrire le stock vrai : ajout, retrait, vente, vérification."},
			{Line: "scanapp", Order: 3, Agent: "delivery", Do: "Garder l'image de la base en interne. Elle ne va pas sur la fiche."},
			{Line: "scanapp", Order: 4, Agent: "comms", Do: "Monter la fiche depuis le produit scanné : nom, prix, code, lieu, description. Le prix dépasse de 20 % à 40 %. CAD, taxes non incluses, Canada."},
			{Line: "scanapp", Order: 5, Agent: "architect", Do: "Photo trouvée en ligne seulement avec les droits. Sinon une image à nous. Jamais une copie d'une photo sans droits."},
			{Line: "scanapp", Order: 6, Agent: "reviewer", Do: "Contrôle : droits photo, et consentement si une personne est dans le cadre."},
			{Line: "scanapp", Order: 7, Agent: "delivery", Do: "Le marketplace est un autre projet. Il fetch ces items pour les vendre. Il ne scanne pas."},
			{Line: "scanapp", Order: 8, Agent: "comms", Do: "Vendre, promouvoir, marketing et média sur la fiche prête. Campagnes de rabais qui comparent les prix en ligne, plus chers."},
			{Line: "scanapp", Order: 9, Agent: "operator", Do: "Vente sur le marketplace (MCP, Stripe) et en live via Empire. Proximity n'est pas dans cette chaîne."},
		}, true
	case "control":
		return []ChainStep{
			{Line: "control", Order: 1, Agent: "manager", Do: "Lire le tableau et nommer les trous. Ne pas trancher à la place d'Epicenter."},
			{Line: "control", Order: 2, Agent: "explorer", Do: "Lire Graphify avant le coup. Pas de deuxième carte."},
			{Line: "control", Order: 3, Agent: "memory", Do: "Ajouter un maillon à la chaîne. Ne pas réécrire le maillon d'avant. Ne pas mélanger le dépôt de l'agence."},
			{Line: "control", Order: 4, Agent: "init", Do: "Voltron relance le kernel s'il tombe. Il ne déploie pas."},
		}, true
	case "wordpress":
		return []ChainStep{
			{Line: "wordpress", Order: 1, Agent: "memory", Do: "Charger le contexte : site WordPress (pages, ACF) ou un produit nommé — Proximity, Next Proximity dans apps/web, Proximity App, Api-Proximity."},
			{Line: "wordpress", Order: 2, Agent: "operator", Do: "Le commit va seulement sur un site WordPress (PHP, ACF Pro, XAMPP) ou sur Proximity, Proximity App, l'API, ou Next. Pas Lovable. Pas Epicenter. Pas Voltron. Pas une autre marque."},
			{Line: "wordpress", Order: 3, Agent: "security", Do: "Filtre : secrets, Loi 25, site déjà en ligne. Interdit si ça sort de la demande."},
			{Line: "wordpress", Order: 4, Agent: "reviewer", Do: "Filtre encore : la plus courte des deux façons, et rien d'autre n'a bougé."},
			{Line: "wordpress", Order: 5, Agent: "comms", Do: "Accord avant de quitter le local. Sans allow, ça reste local."},
			{Line: "wordpress", Order: 6, Agent: "memory", Do: "Écrire le contexte à jour : pages, champs ACF, ou le produit touché, et ce qui reste."},
		}, true
	case "proximity":
		return []ChainStep{
			{Line: "proximity", Order: 1, Agent: "memory", Do: "Ces comptes ne sont pas le département WordPress. Proximity, Next, Proximity App et l'API n'y sont pas."},
			{Line: "proximity", Order: 2, Agent: "manager", Do: "Ne pas committer depuis cette ligne sur un thème, ni sur Proximity, Proximity App, Api-Proximity, ou Next."},
		}, true
	case "marketplace":
		return []ChainStep{
			{Line: "marketplace", Order: 1, Agent: "delivery", Do: "Fetch les items du scanner. Ne pas scanner. Ne pas déployer le site."},
			{Line: "marketplace", Order: 2, Agent: "comms", Do: "Préparer la fiche. Prix 20 % à 40 % au-dessus. Rabais contre les prix en ligne, plus chers. Le site reste hors ligne tant qu'il n'est pas prêt."},
		}, true
	case "panda":
		return []ChainStep{
			{Line: "panda", Order: 1, Agent: "architect", Do: "Décrire l'offre white-glove : IA, installation, cours."},
			{Line: "panda", Order: 2, Agent: "comms", Do: "Vendre le service. Ne pas installer chez le client sans accord."},
		}, true
	case "nft-giant":
		return []ChainStep{
			{Line: "nft-giant", Order: 1, Agent: "security", Do: "Revue avant un mint ou une campagne. Pas de mint sans cette revue."},
			{Line: "nft-giant", Order: 2, Agent: "investigator", Do: "Un tirage reste classé avant publication. Pas une loterie lancée d'ici."},
		}, true
	case "ecole":
		return []ChainStep{
			{Line: "ecole", Order: 1, Agent: "research", Do: "Décrire le cours tech. Ne pas promettre un gain."},
			{Line: "ecole", Order: 2, Agent: "planner", Do: "Ranger le programme. Le cours n'exécute pas un ordre."},
		}, true
	case "marketing":
		return []ChainStep{
			{Line: "marketing", Order: 1, Agent: "planner", Do: "Une campagne Pandora. Mesurer avant d'élargir."},
			{Line: "marketing", Order: 2, Agent: "explorer", Do: "Lire le contexte Pandora et CRM. Pas le compte d'un client Proximity."},
			{Line: "marketing", Order: 3, Agent: "comms", Do: "Facebook ou courriel seulement après consentement, identité, désabonnement, et allow. Les rabais comparent les prix en ligne."},
		}, true
	case "empire":
		return []ChainStep{
			{Line: "empire", Order: 1, Agent: "delivery", Do: "Live seulement avec une fiche déjà prête. Pas de stock inventé."},
			{Line: "empire", Order: 2, Agent: "comms", Do: "Toutes les plateformes live. Proximity ne publie pas ces fiches."},
		}, true
	case "propres":
		return []ChainStep{
			{Line: "propres", Order: 1, Agent: "comms", Do: "Stripe sur PBTM, Pandora et le marketing. Ce ne sont pas des clients."},
			{Line: "propres", Order: 2, Agent: "security", Do: "Refuser ce Stripe sur un site Proximity sans demande explicite."},
		}, true
	case "trading":
		return []ChainStep{
			{Line: "trading", Order: 1, Agent: "investigator", Do: "Aucun ordre live sans réponse AMF écrite."},
			{Line: "trading", Order: 2, Agent: "architect", Do: "Le modèle reste un modèle. Giant n'est pas un ordre."},
		}, true
	default:
		return nil, false
	}
}
