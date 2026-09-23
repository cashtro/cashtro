package agents

// HustlePartner is The Hustle. One chain, twenty agents, five brains each.
// It reads, holds the street metal and the corporate metal, and acts only
// through a legal company. Epicenter Einstein is its friend and decides.
type HustlePartner struct {
	Name           string `json:"name"`
	Repo           string `json:"repo"`
	Brain          string `json:"brain"`
	Agents         int    `json:"agents"`
	Brains         int    `json:"brains"`
	BestFriend     string `json:"bestFriend"`
	KillerInstinct bool   `json:"killerInstinct"`
	LoadsInstinct  bool   `json:"loadsInstinctBrain"`
	OwnBrain       bool   `json:"ownBrain"`
	MainBrain      string `json:"mainBrain"`
	SameShape      bool   `json:"sameShape"`
	Niche          string `json:"niche"`
	LoadsVoltron   bool   `json:"loadsVoltronBrain"`
	Automated      bool   `json:"automated"`
	Weapons        bool   `json:"weapons"`
	Attacks        bool   `json:"attacks"`
	Illegal        bool   `json:"illegal"`
	Sent           bool   `json:"sent"`
	Paid           bool   `json:"paid"`
	Launched       bool   `json:"launched"`
	OpenedBy       string `json:"openedBy"`
	OrderedBy      string `json:"orderedBy"`
	DecidedBy      string `json:"decidedBy"`
}

// Hustle is the roster. The working chain lives in Evolu-Jeunes/Hustle/brain/hustle.ts.
func Hustle() HustlePartner {
	return HustlePartner{
		Name:           "The Hustle",
		Repo:           "Evolu-Jeunes/Hustle",
		Brain:          "hustle",
		Agents:         20,
		Brains:         100,
		BestFriend:     "instinct",
		KillerInstinct: true,
		LoadsInstinct:  false,
		OwnBrain:       true,
		MainBrain:      "instinct",
		SameShape:      true,
		Niche:          "métaux légaux",
		LoadsVoltron:   true,
		Automated:      true,
		Weapons:        false,
		Attacks:        false,
		Illegal:        false,
		Sent:           false,
		Paid:           false,
		Launched:       false,
		OpenedBy:       "voltron",
		OrderedBy:      "scrum",
		DecidedBy:      "instinct",
	}
}
