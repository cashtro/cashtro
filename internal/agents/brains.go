package agents

// Brain is one of the five minds the Manager steers.
type Brain struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Verb  string   `json:"verb"`
	Color string   `json:"color"`
	Subs  []string `json:"subs"`
}

// Brains is the Epicenter roster: Architecte, Cartographe, Forgeron, Orfèvre, Hustler.
func Brains() []Brain {
	return []Brain{
		{
			ID: "architecte", Name: "L'Architecte", Verb: "pense · ML/DL", Color: "#f0c14b",
			Subs: []string{"Meta-Cerveau", "Relativité", "Geomatria", "Marché", "Tisseur"},
		},
		{
			ID: "cartographe", Name: "Le Cartographe", Verb: "voit · Graphify", Color: "#6ea8ff",
			Subs: []string{"Planificateur", "Mappeur", "Graphifieur", "Moniteur", "Philosophe"},
		},
		{
			ID: "forgeron", Name: "Le Forgeron", Verb: "construit 24/7", Color: "#ff4d3d",
			Subs: []string{"Soudeur", "Mécanicien", "Tisserand", "Mineur", "Ambassadeur"},
		},
		{
			ID: "orfevre", Name: "L'Orfèvre", Verb: "exécute 24/7", Color: "#3dffa6",
			Subs: []string{"Contrôleur", "Livreur", "Polisseur", "Documenteur", "Diplomate"},
		},
		{
			ID: "hustler", Name: "Le Hustler", Verb: "empire A à Z", Color: "#c084fc",
			Subs: []string{"Stratège", "Traqueur", "Recruteur", "Empire", "Pontife"},
		},
	}
}

// BrainByName matches a line.Brain label such as "Forgeron".
func BrainByName(name string) (Brain, bool) {
	for _, b := range Brains() {
		if b.Name == name || b.ID == name {
			return b, true
		}
		switch name {
		case "Architecte", "L'Architecte":
			if b.ID == "architecte" {
				return b, true
			}
		case "Cartographe", "Le Cartographe":
			if b.ID == "cartographe" {
				return b, true
			}
		case "Forgeron", "Le Forgeron":
			if b.ID == "forgeron" {
				return b, true
			}
		case "Orfèvre", "L'Orfèvre", "Orfevre":
			if b.ID == "orfevre" {
				return b, true
			}
		case "Hustler", "Le Hustler":
			if b.ID == "hustler" {
				return b, true
			}
		}
	}
	return Brain{}, false
}
