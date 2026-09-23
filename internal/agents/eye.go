package agents

// TheEye is the defensive watch. It reads cashtro and Evolu-Jeunes,
// names a weakness, and writes a patch. It does not attack or copy a secret.
type Eye struct {
	Name          string   `json:"name"`
	Repo          string   `json:"repo"`
	Brain         string   `json:"brain"`
	Scope         []string `json:"scope"`
	Steps         []string `json:"steps"`
	Attacks       bool     `json:"attacks"`
	Flood         bool     `json:"flood"`
	CopiesSecrets bool     `json:"copiesSecrets"`
	AppliesPatch  bool     `json:"appliesPatch"`
	DecidedBy     string   `json:"decidedBy"`
}

// TheEye returns the watch roster. The working look lives in Evolu-Jeunes/Forge/brain/eye.ts.
func TheEye() Eye {
	return Eye{
		Name:          "The Eye",
		Repo:          "Evolu-Jeunes/Forge",
		Brain:         "eye",
		Scope:         []string{"cashtro", "Evolu-Jeunes"},
		Steps:         []string{"veille", "menace", "faille", "recherche", "rustine", "preuve", "garde", "sceau"},
		Attacks:       false,
		Flood:         false,
		CopiesSecrets: false,
		AppliesPatch:  false,
		DecidedBy:     "epicenter",
	}
}
