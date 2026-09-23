package agents

// ForgePartner is the development conglomerate. It loads Tartaria and Voltron.
// Scrum takes every repository. The crew does not make a weapon.
type ForgePartner struct {
	Repo          string   `json:"repo"`
	Brain         string   `json:"brain"`
	LoadsTartaria bool     `json:"loadsTartariaBrain"`
	LoadsVoltron  bool     `json:"loadsVoltronBrain"`
	Weapons       bool     `json:"weapons"`
	Partners      []string `json:"partners"`
	Departments   []string `json:"departments"`
	Agents        int      `json:"agents"`
	Assistants    int      `json:"assistants"`
	Attacks       bool     `json:"attacks"`
	Flood         bool     `json:"flood"`
	Deployed      bool     `json:"deployed"`
	OpenedBy      string   `json:"openedBy"`
	OrderedBy     string   `json:"orderedBy"`
	DecidedBy     string   `json:"decidedBy"`
}

// Forge is the partner roster. The count matches Evolu-Jeunes/Forge/brain/forge.ts.
func Forge() ForgePartner {
	return ForgePartner{
		Repo:          "Evolu-Jeunes/Forge",
		Brain:         "forge",
		LoadsTartaria: true,
		LoadsVoltron:  true,
		Weapons:       false,
		Partners:      []string{"voltron", "scrum", "tartaria"},
		Departments:   []string{"giant", "web", "security", "scrum"},
		Agents:        100,
		Assistants:    100,
		Attacks:       false,
		Flood:         false,
		Deployed:      false,
		OpenedBy:      "voltron",
		OrderedBy:     "scrum",
		DecidedBy:     "tartaria",
	}
}
