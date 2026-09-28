package agents

// SeparatedChain is a project the agentic does not run.
type SeparatedChain struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Repo    string `json:"repo"`
	Owner   string `json:"owner"`
	With    string `json:"with"`
	Agentic bool   `json:"agentic"`
	Note    string `json:"note"`
}

// Chain1 is Éduconnexion. Another organism owns it.
// The work was done with Proximity agency. It is not an agentic chain.
func Chain1() SeparatedChain {
	return SeparatedChain{
		ID:      "chain-1",
		Name:    "Chain 1",
		Repo:    "Evolu-Jeunes/educonnexion",
		Owner:   "autre organisme",
		With:    "Proximity agency",
		Agentic: false,
		Note:    "Éduconnexion appartient à un autre organisme. Le projet est fait avec Proximity agency. Chain 1 n'est pas une chaîne agentique.",
	}
}
