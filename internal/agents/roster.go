package agents

import "github.com/cashtro/cashtro/internal/kernel"

// Builders are the public names processes wear on the desk.
//
// These are builder / creator names — Lautaro, Frida, Diego — not product
// brands and not philosopher surnames. Process IDs and capability verbs stay
// stable so the kernel can route. Casa Crypto, Prolifik, ScanApp, Proximity
// stay on the delivery line. A process never takes a product's name.
var Builders = map[string]string{
	"init":         "Lautaro",
	"architect":    "Frida",
	"operator":     "Diego",
	"delivery":     "Antoni",
	"research":     "Gabriela",
	"explorer":     "Marco",
	"reviewer":     "Juana",
	"deploy":       "Oscar",
	"security":     "Túpac",
	"memory":       "Luis",
	"comms":        "Violeta",
	"planner":      "Pablo",
	"investigator": "Caupolicán",
	"router":       "Hedy",
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
