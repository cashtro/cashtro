package agents

import "strings"

// Question is asked before an action, to fill a context gap.
type Question struct {
	ID       string `json:"id"`
	Ask      string `json:"ask"`
	Improves string `json:"improves"`
}

// QuestionsFor returns the specific questions that must be answered
// before this action can take a position.
func QuestionsFor(action string) []Question {
	switch action {
	case "scanapp", "chain":
		return []Question{
			{ID: "boutique", Ask: "Quelle boutique ou quel client reçoit cette chaîne maintenant, pas tous les sites?", Improves: "on ne déploie pas la chaîne au mauvais endroit"},
			{ID: "codes", Ask: "Quels codes on scanne dans cette action : EAN-13, QR, ou un code interne?", Improves: "le scanner est réglé sur le vrai format"},
			{ID: "photo", Ask: "La photo de fiche est générée par nous, ou prise dans quelle banque dont on a déjà la licence?", Improves: "on ne publie pas l'image de la base ni une photo grattée"},
			{ID: "site", Ask: "Quel repo Proximity publie la fiche, et ce site est-il déjà en ligne chez le client?", Improves: "un site live se lit, il ne se réécrit pas"},
			{ID: "kick", Ask: "Ce produit passe-t-il dans un live Kick d'Empire dans cette action, ou seulement sur le site?", Improves: "le live ne vend pas un item qui n'est pas dans cette action"},
			{ID: "prix", Ask: "Le prix affiché est-il en dollars canadiens, taxes incluses, pour un acheteur au Québec?", Improves: "la fiche respecte le prix qu'on a le droit d'afficher"},
		}
	case "fiche":
		return []Question{
			{ID: "acheteur", Ask: "Qui achète cette fiche : un particulier, un détaillant, ou le live?", Improves: "le texte et le prix visent le bon acheteur"},
			{ID: "droits", Ask: "Qui détient les droits de la photo, et où est la preuve?", Improves: "la photo peut sortir"},
		}
	case "assign":
		return []Question{
			{ID: "pourquoi", Ask: "Pourquoi ce cerveau pour ce repo, et quel cerveau on écarte?", Improves: "l'assignation n'est pas un réflexe"},
			{ID: "live", Ask: "Ce repo est-il un site client déjà en production?", Improves: "on ne programme pas un site live comme un brouillon"},
		}
	case "contradict":
		return []Question{
			{ID: "situation", Ask: "Que se passe-t-il concrètement maintenant, avant de choisir une option?", Improves: "le choix part du fait, pas du nom de l'option"},
			{ID: "intouchable", Ask: "Qu'est-ce qui ne doit pas changer : un site live, une règle, un budget?", Improves: "l'option optimisée ne casse pas ce qui est fixé"},
			{ID: "mieux", Ask: "Ici, « mieux » veut dire moins de temps, moins d'argent, ou moins de risque?", Improves: "le score mesure la bonne chose"},
		}
	default:
		return []Question{
			{ID: "situation", Ask: "Quelle est la situation concrète, là, avant cette action?", Improves: "on ne vise pas une position sans le fait"},
			{ID: "manque", Ask: "Quel contexte manque encore pour que cette action soit juste?", Improves: "la lacune est nommée avant le geste"},
			{ID: "intouchable", Ask: "Qu'est-ce que cette action ne doit pas changer?", Improves: "le geste reste dans sa limite"},
		}
	}
}

// OpenQuestions returns the questions that still have no answer.
func OpenQuestions(action string, answers map[string]string) []Question {
	var open []Question
	for _, q := range QuestionsFor(action) {
		if strings.TrimSpace(answers[q.ID]) == "" {
			open = append(open, q)
		}
	}
	return open
}

func askMessage(open []Question) string {
	if len(open) == 0 {
		return "contexte suffisant"
	}
	parts := make([]string, 0, len(open))
	for _, q := range open {
		parts = append(parts, q.Ask)
	}
	return "questions avant l'action : " + strings.Join(parts, " ")
}
