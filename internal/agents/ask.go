package agents

import "strings"

// Question is asked before an action, to fill a context gap.
type Question struct {
	ID       string `json:"id"`
	Ask      string `json:"ask"`
	Improves string `json:"improves"`
}

// StrategyQuestions are asked before every move. One question, the position,
// the power move, and the reply it opens. Like a chess clock.
func StrategyQuestions() []Question {
	return []Question{
		{ID: "position", Ask: "Quelle est la position sur le tableau, avant ce coup?", Improves: "on ne joue pas sans voir la pièce"},
		{ID: "coup", Ask: "Quel est le coup de pouvoir, et quelle réponse adverse il ouvre?", Improves: "un coup sans réponse adverse n'est pas un coup"},
	}
}

func strategyKnown() map[string]string {
	return map[string]string{
		"position": "Le tableau est lu avant le coup. Les sites clients déjà en ligne ne bougent pas. Le marketplace n'est pas déployé. Le scanner reste interne.",
		"coup":     "Le coup de pouvoir est le plus court : moins d'étapes, moins de coût, moins de risque. La réponse adverse est nommée avant de jouer. Une seule option n'est pas un coup.",
	}
}

// QuestionsFor returns the specific questions that must be answered
// before this action can take a position. The two strategy questions
// come first, on every action.
func QuestionsFor(action string) []Question {
	return append(StrategyQuestions(), questionsFor(action)...)
}

func questionsFor(action string) []Question {
	switch action {
	case "scanapp", "chain":
		return []Question{
			{ID: "catalogue", Ask: "Empire vend tout le stock du scanner déployé, ou seulement les produits déjà mis en fiche par le marketplace?", Improves: "le live ne vend pas un item hors fiche"},
		}
	case "marketplace":
		return []Question{
			{ID: "marge", Ask: "De combien le prix de vente dépasse le prix de référence, et ce prix de référence est lequel : le coût, ou le prix vu en ligne?", Improves: "vendre plus cher sans inventer un prix"},
			{ID: "url", Ask: "Le site du marketplace est déjà en ligne à quelle adresse? Le scanner inventory-scanner-theta.vercel.app n'est pas ce site.", Improves: "Stripe se branche sur le site qui vend"},
			{ID: "stripe", Ask: "Stripe encaisse seulement sur le site du marketplace, ou aussi pendant le live Empire?", Improves: "un seul chemin d'argent, ou deux"},
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

// Known is what was already said. The system uses it instead of asking again.
func Known(action string) map[string]string {
	switch action {
	case "scanapp", "chain":
		return map[string]string{
			"catalogue": "Le Scan App est interne et déjà déployé. Il ne vend pas. Le marketplace vend les items scannés. Empire vend aussi ces produits en live, s'ils ont une fiche prête.",
		}
	case "marketplace":
		return map[string]string{
			"marge":  "Les prix dépassent de 20 % à 40 %. Des rabais seront lancés parce que les vrais prix en ligne sont plus chers. La comparaison en ligne sert la campagne marketing.",
			"url":    "Pas en ligne. Déployer une fois prêt. Le scanner déployé n'est pas ce site.",
			"stripe": "Stripe encaisse les projets propres : marketplace une fois déployé, Empire, PBTM, Pandora, business, technology et marketing. Pas Proximity, sauf demande. Pas le scanner, qui est interne.",
		}
	default:
		return nil
	}
}

// SelfCheck is the system questioning itself before it moves.
type SelfCheck struct {
	Action   string            `json:"action"`
	Asked    []Question        `json:"asked"`
	Answered map[string]string `json:"answered"`
	Open     []Question        `json:"open"`
}

// AskSelf poses the questions and answers them from what is already known.
// The strategy answers are always filled, so the chess question is asked
// and the move can still proceed.
func AskSelf(action string, extra map[string]string) SelfCheck {
	answered := map[string]string{}
	for k, v := range strategyKnown() {
		answered[k] = v
	}
	for k, v := range Known(action) {
		answered[k] = v
	}
	for k, v := range extra {
		if strings.TrimSpace(v) != "" {
			answered[k] = strings.TrimSpace(v)
		}
	}
	return SelfCheck{
		Action:   action,
		Asked:    QuestionsFor(action),
		Answered: answered,
		Open:     OpenQuestions(action, answered),
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
