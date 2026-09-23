package agents

// OpsPartner is the operations chain inside Voltron. Three chains, twenty
// agents, five brains each. It is not a website. It does not send or charge.
type OpsPartner struct {
	Repo            string   `json:"repo"`
	Brain           string   `json:"brain"`
	Chains          []string `json:"chains"`
	Agents          int      `json:"agents"`
	Brains          int      `json:"brains"`
	LoadsInstinct   bool     `json:"loadsInstinctBrain"`
	OwnBrain        bool     `json:"ownBrain"`
	MainBrain       string   `json:"mainBrain"`
	SameShape       bool     `json:"sameShape"`
	Niche           string   `json:"niche"`
	LoadsVoltron    bool     `json:"loadsVoltronBrain"`
	ReportsTo       string   `json:"reportsTo"`
	InVoltron       bool     `json:"inVoltron"`
	Operative       bool     `json:"operative"`
	Website         bool     `json:"website"`
	OptionalThrough string   `json:"optionalThrough"`
	Automated       bool     `json:"automated"`
	Weapons         bool     `json:"weapons"`
	Sent            bool     `json:"sent"`
	Paid            bool     `json:"paid"`
	Dialed          bool     `json:"dialed"`
	OpenedBy        string   `json:"openedBy"`
	OrderedBy       string   `json:"orderedBy"`
	DecidedBy       string   `json:"decidedBy"`
}

// Ops is the roster. The working functions live in Evolu-Jeunes/OPS/brain/ops.ts.
func Ops() OpsPartner {
	return OpsPartner{
		Repo:            "Evolu-Jeunes/OPS",
		Brain:           "ops",
		Chains:          []string{"accueil", "terrain", "propriete"},
		Agents:          60,
		Brains:          300,
		LoadsInstinct:   false,
		OwnBrain:        true,
		MainBrain:       "instinct",
		SameShape:       true,
		Niche:           "opérations en brouillon",
		LoadsVoltron:    true,
		ReportsTo:       "voltron",
		InVoltron:       true,
		Operative:       true,
		Website:         false,
		OptionalThrough: "voltron",
		Automated:       true,
		Weapons:         false,
		Sent:            false,
		Paid:            false,
		Dialed:          false,
		OpenedBy:        "voltron",
		OrderedBy:       "scrum",
		DecidedBy:       "instinct",
	}
}
