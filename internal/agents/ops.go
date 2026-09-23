package agents

// OpsPartner is the operations conglomerate. Three chains, twenty agents,
// five brains each. It loads Tartaria and Voltron. It does not send or charge.
type OpsPartner struct {
	Repo          string   `json:"repo"`
	Source        string   `json:"source"`
	Brain         string   `json:"brain"`
	Chains        []string `json:"chains"`
	Agents        int      `json:"agents"`
	Brains        int      `json:"brains"`
	LoadsTartaria bool     `json:"loadsTartariaBrain"`
	LoadsVoltron  bool     `json:"loadsVoltronBrain"`
	Automated     bool     `json:"automated"`
	Weapons       bool     `json:"weapons"`
	Sent          bool     `json:"sent"`
	Paid          bool     `json:"paid"`
	Dialed        bool     `json:"dialed"`
	OpenedBy      string   `json:"openedBy"`
	OrderedBy     string   `json:"orderedBy"`
	DecidedBy     string   `json:"decidedBy"`
}

// Ops is the roster. The working functions live in Evolu-Jeunes/OPS/brain/ops.ts.
func Ops() OpsPartner {
	return OpsPartner{
		Repo:          "Evolu-Jeunes/OPS",
		Source:        "https://web-ops.shop/",
		Brain:         "ops",
		Chains:        []string{"accueil", "terrain", "propriete"},
		Agents:        60,
		Brains:        300,
		LoadsTartaria: true,
		LoadsVoltron:  true,
		Automated:     true,
		Weapons:       false,
		Sent:          false,
		Paid:          false,
		Dialed:        false,
		OpenedBy:      "voltron",
		OrderedBy:     "scrum",
		DecidedBy:     "tartaria",
	}
}
