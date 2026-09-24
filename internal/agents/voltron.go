package agents

// Layer is one agentic chain held inside Voltron.
// It reports to Voltron, stays operative there, and is not a website.
// An optional function runs only when Voltron calls it.
type Layer struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Repo            string `json:"repo"`
	ReportsTo       string `json:"reportsTo"`
	InVoltron       bool   `json:"inVoltron"`
	Operative       bool   `json:"operative"`
	Website         bool   `json:"website"`
	OptionalThrough string `json:"optionalThrough"`
}

// Layers is every agentic chain Voltron keeps on.
func Layers() []Layer {
	return []Layer{
		{ID: "ops", Name: "OPS", Repo: "Evolu-Jeunes/OPS", ReportsTo: "voltron", InVoltron: true, Operative: true, Website: false, OptionalThrough: "voltron"},
		{ID: "hustle", Name: "The Hustle", Repo: "Evolu-Jeunes/Hustle", ReportsTo: "voltron", InVoltron: true, Operative: true, Website: false, OptionalThrough: "voltron"},
		{ID: "eye", Name: "The Eye", Repo: "Evolu-Jeunes/Forge", ReportsTo: "voltron", InVoltron: true, Operative: true, Website: false, OptionalThrough: "voltron"},
		{ID: "forge", Name: "Forge", Repo: "Evolu-Jeunes/Forge", ReportsTo: "voltron", InVoltron: true, Operative: true, Website: false, OptionalThrough: "voltron"},
		{ID: "giant", Name: "Giant", Repo: "Evolu-Jeunes/Giant", ReportsTo: "voltron", InVoltron: true, Operative: true, Website: false, OptionalThrough: "voltron"},
		{ID: "fusion", Name: "The Fusion", Repo: "cashtro/cashtro", ReportsTo: "voltron", InVoltron: true, Operative: true, Website: false, OptionalThrough: "voltron"},
	}
}

// VoltronHolds reports whether every chain is inside Voltron and operative.
func VoltronHolds() (bool, string) {
	layers := Layers()
	if len(layers) == 0 {
		return false, "aucune chaîne"
	}
	for _, layer := range layers {
		if layer.ReportsTo != "voltron" || !layer.InVoltron || !layer.Operative || layer.Website || layer.OptionalThrough != "voltron" {
			return false, layer.ID
		}
	}
	return true, ""
}
