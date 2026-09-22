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
	default:
		return nil, false
	}
}
