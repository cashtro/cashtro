package agents

// ForgePartner is the development side brain. It works with Voltron and Epicenter.
// It does not load Epicenter's brain. Security audits these repositories and does not attack.
type ForgePartner struct {
	Repo       string   `json:"repo"`
	Brain      string   `json:"brain"`
	Loads      bool     `json:"loadsEpicenterBrain"`
	Partners   []string `json:"partners"`
	Departments []string `json:"departments"`
	Agents     int      `json:"agents"`
	Assistants int      `json:"assistants"`
	Attacks    bool     `json:"attacks"`
	Flood      bool     `json:"flood"`
	Deployed   bool     `json:"deployed"`
	OpenedBy   string   `json:"openedBy"`
	OrderedBy  string   `json:"orderedBy"`
	DecidedBy  string   `json:"decidedBy"`
}

// Forge is the partner roster. The count matches Evolu-Jeunes/Forge/brain/forge.ts.
func Forge() ForgePartner {
	return ForgePartner{
		Repo:        "Evolu-Jeunes/Forge",
		Brain:       "forge",
		Loads:       false,
		Partners:    []string{"voltron", "scrum", "epicenter"},
		Departments: []string{"giant", "web", "security", "scrum"},
		Agents:      100,
		Assistants:  100,
		Attacks:     false,
		Flood:       false,
		Deployed:    false,
		OpenedBy:    "voltron",
		OrderedBy:   "scrum",
		DecidedBy:   "epicenter",
	}
}
