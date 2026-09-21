package agents

import "github.com/cashtro/cashtro/internal/kernel"

// Thinkers are the public names processes wear on the desk.
//
// Process IDs and capability verbs stay stable so the kernel can route.
// Product names (Casa Crypto, Prolifik, ScanApp, Proximity, …) stay on the
// delivery line. A process never takes a product's name.
var Thinkers = map[string]string{
	"init":         "Lovelace",
	"router":       "Turing",
	"delivery":     "Deming",
	"research":     "Hypatia",
	"explorer":     "Copernicus",
	"operator":     "Galileo",
	"reviewer":     "Kant",
	"architect":    "Leonardo",
	"deploy":       "Hopper",
	"security":     "Machiavelli",
	"memory":       "Locke",
	"comms":        "Voltaire",
	"planner":      "Confucius",
	"investigator": "Socrates",
}

// ProductNames must never appear as process display names.
var ProductNames = []string{
	"Casa Crypto",
	"Prolifik",
	"ScanApp",
	"Proximity",
	"BTK Avocats",
	"MD Clinic",
	"Solution Hypothèque QC",
	"Éduconnexion",
}

// Thinker returns the desk name for a process id.
func Thinker(id string) string {
	if name, ok := Thinkers[id]; ok {
		return name
	}
	return id
}

func wear(spec kernel.Spec) kernel.Spec {
	spec.Name = Thinker(spec.ID)
	return spec
}
