package agents

// TheEye is the defensive watch. It reads cashtro and Evolu-Jeunes
// and names a weakness. Every change goes through Instinct first.
// Instinct sits there. The watch does not attack or copy a secret.
type Eye struct {
	Name            string   `json:"name"`
	Repo            string   `json:"repo"`
	Brain           string   `json:"brain"`
	Gate            string   `json:"gate"`
	Scope           []string `json:"scope"`
	Steps           []string `json:"steps"`
	Attacks         bool     `json:"attacks"`
	Flood           bool     `json:"flood"`
	CopiesSecrets   bool     `json:"copiesSecrets"`
	AppliesPatch    bool     `json:"appliesPatch"`
	ThroughInstinct bool     `json:"throughInstinct"`
	DecidedBy       string   `json:"decidedBy"`
}

// TheEye returns the watch roster. The working look lives in Evolu-Jeunes/Forge/brain/eye.ts.
func TheEye() Eye {
	return Eye{
		Name:            "The Eye",
		Repo:            "Evolu-Jeunes/Forge",
		Brain:           "eye",
		Gate:            "cashtro/epicenter",
		Scope:           []string{"cashtro", "Evolu-Jeunes"},
		Steps:           []string{"veille", "menace", "faille", "recherche", "garde", "instinct", "rustine", "preuve", "sceau"},
		Attacks:         false,
		Flood:           false,
		CopiesSecrets:   false,
		AppliesPatch:    false,
		ThroughInstinct: true,
		DecidedBy:       "instinct",
	}
}

// EyeChange opens a draft only after Instinct. The draft is not applied.
func EyeChange(via, seat string) (bool, string) {
	if via != "instinct" && via != "cashtro/epicenter" {
		return false, "The Eye passe par Instinct avant tout changement."
	}
	if seat != "instinct" {
		return false, "Sans Instinct, aucun changement."
	}
	return true, "brouillon ouvert. rien n'est appliqué"
}
