package agents

import "github.com/cashtro/cashtro/internal/kernel"

// Builders are the public names processes wear on the desk.
//
// Castro's call: names like Builder, Creator, Lautaro — then keep going
// in that voice. Not product brands. Not philosopher surnames.
// Process IDs and capability verbs stay stable so the kernel can route.
// Casa Crypto, Prolifik, ScanApp, Proximity stay on the delivery line.
var Builders = map[string]string{
	"init":         "Lautaro",
	"operator":     "Builder",
	"architect":    "Creator",
	"delivery":     "Maker",
	"deploy":       "Forger",
	"research":     "Scribe",
	"explorer":     "Pathfinder",
	"reviewer":     "Witness",
	"security":     "Sentinel",
	"memory":       "Keeper",
	"comms":        "Herald",
	"planner":      "Steward",
	"investigator": "Seeker",
	"router":       "Spark",
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

// Builder returns the desk name for a process id.
func Builder(id string) string {
	if name, ok := Builders[id]; ok {
		return name
	}
	return id
}

func wear(spec kernel.Spec) kernel.Spec {
	spec.Name = Builder(spec.ID)
	return spec
}
