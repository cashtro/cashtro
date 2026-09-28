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

// VoltronStart is one agentic chain Voltron starts. A website is not a start.
type VoltronStart struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Task      string `json:"task"`
	StartedBy string `json:"startedBy"`
	Acts      bool   `json:"acts"`
	Replies   bool   `json:"replies"`
}

// VoltronOS is the operating system of the agentics.
// The brain sits beside the three front-end stages. Every chain with tasks starts here.
type VoltronOS struct {
	Name    string         `json:"name"`
	Role    string         `json:"role"`
	Brain   string         `json:"brain"`
	Stages  []string       `json:"stages"`
	Acts    bool           `json:"acts"`
	Replies bool           `json:"replies"`
	Starts  []VoltronStart `json:"starts"`
}

// TheOS is Voltron. Pressing a start is an agentic doing its task.
func TheOS() VoltronOS {
	starts := make([]VoltronStart, 0, len(Layers())+4)
	for _, layer := range Layers() {
		starts = append(starts, VoltronStart{
			ID: layer.ID, Name: layer.Name, Kind: "layer",
			Task: "Voltron starts " + layer.Name + ". The chain acts from here.",
			StartedBy: "voltron", Acts: true, Replies: false,
		})
	}
	starts = append(starts,
		VoltronStart{ID: "nft-giant", Name: "Giant chain", Kind: "chain", Task: "Voltron opens the Giant chain. The brain is Giant's. The repo stays Evolu-Jeunes/Giant.", StartedBy: "voltron", Acts: true, Replies: false},
		VoltronStart{ID: "trading", Name: "Trading bot", Kind: "chain", Task: "Voltron starts the trading-bot chain. The Giant repo stays apart.", StartedBy: "voltron", Acts: true, Replies: false},
		VoltronStart{ID: "giant-trading", Name: "Synchronie Giant", Kind: "sync", Task: "Voltron starts the sync. The two repos work together and stay two repos.", StartedBy: "voltron", Acts: true, Replies: false},
		VoltronStart{ID: "maitre", Name: "Twenty-seat chain", Kind: "crew", Task: "Voltron starts the maître's chain. Einstein is the brain beside it.", StartedBy: "voltron", Acts: true, Replies: false},
	)
	return VoltronOS{
		Name:    "Voltron",
		Role:    "OS of the agentics",
		Brain:   "instinct",
		Stages:  []string{"idea", "concept", "production"},
		Acts:    true,
		Replies: false,
		Starts:  starts,
	}
}

// VoltronStartBy returns one chain Voltron can start.
func VoltronStartBy(id string) (VoltronStart, bool) {
	for _, start := range TheOS().Starts {
		if start.ID == id {
			return start, true
		}
	}
	return VoltronStart{}, false
}
