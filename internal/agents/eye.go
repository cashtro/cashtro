package agents

// TheEye is the defensive watch. It reads cashtro and Evolu-Jeunes
// and names a weakness. Every change goes through Epicenter first.
// Instinct sits there. The watch does not attack or copy a secret.
type Eye struct {
	Name             string   `json:"name"`
	Repo             string   `json:"repo"`
	Brain            string   `json:"brain"`
	Gate             string   `json:"gate"`
	Scope            []string `json:"scope"`
	Steps            []string `json:"steps"`
	Attacks          bool     `json:"attacks"`
	Flood            bool     `json:"flood"`
	CopiesSecrets    bool     `json:"copiesSecrets"`
	AppliesPatch     bool     `json:"appliesPatch"`
	ThroughEpicenter bool     `json:"throughEpicenter"`
	DecidedBy        string   `json:"decidedBy"`
}

// TheEye returns the watch roster. The working look lives in Evolu-Jeunes/Forge/brain/eye.ts.
func TheEye() Eye {
	return Eye{
		Name:             "The Eye",
		Repo:             "Evolu-Jeunes/Forge",
		Brain:            "eye",
		Gate:             "cashtro/epicenter",
		Scope:            []string{"cashtro", "Evolu-Jeunes"},
		Steps:            []string{"veille", "menace", "faille", "recherche", "garde", "epicenter", "rustine", "preuve", "sceau"},
		Attacks:          false,
		Flood:            false,
		CopiesSecrets:    false,
		AppliesPatch:     false,
		ThroughEpicenter: true,
		DecidedBy:        "instinct",
	}
}

// EyeChange opens a draft only after Epicenter. The draft is not applied.
func EyeChange(via, seat string) (bool, string) {
	if via != "epicenter" && via != "cashtro/epicenter" {
		return false, "The Eye passe par Epicenter avant tout changement."
	}
	if seat != "instinct" {
		return false, "Instinct siège dans Epicenter. Sans lui, aucun changement."
	}
	return true, "brouillon ouvert. rien n'est appliqué"
}
