package agents

import "strings"

// Craft is one mark a department must leave on its own result.
// If the mark is missing, the department fills it. That is the lacune.
type Craft struct {
	Mark string `json:"mark"`
	Fill string `json:"fill"`
}

// Improvement is a department closing its own gaps.
type Improvement struct {
	Department  string   `json:"department"`
	Lacunes     []string `json:"lacunes"`
	Result      string   `json:"result"`
	Specialized bool     `json:"specialized"`
}

func craftFor(id string) []Craft {
	switch id {
	case "direction":
		return []Craft{
			{Mark: "option rejetée", Fill: "Option rejetée : nommer celle qu'on ne fait pas, et pourquoi."},
			{Mark: "responsable", Fill: "Responsable : le chef qui tranche."},
		}
	case "ingenierie":
		return []Craft{
			{Mark: "deux façons", Fill: "Deux façons : la construction choisie et l'autre, avec le coût de chacune."},
			{Mark: "limite", Fill: "Limite : ce qu'on ne construit pas dans ce changement."},
		}
	case "marche":
		return []Craft{
			{Mark: "offre", Fill: "Offre : une seule, nommée, et pour qui."},
			{Mark: "mesure", Fill: "Mesure : le chiffre qu'on regarde avant d'élargir."},
		}
	case "operations":
		return []Craft{
			{Mark: "livré", Fill: "Livré : l'artifact concret, pas l'intention."},
			{Mark: "pas fini", Fill: "Pas fini : ce qui reste avant de dire que c'est en ligne."},
		}
	case "controle":
		return []Craft{
			{Mark: "règle", Fill: "Règle : permis ou interdit, citée, pas un avis."},
			{Mark: "ça casse", Fill: "Ça casse : le cas qui fait perdre de l'argent ou des données."},
		}
	case "wordpress":
		return []Craft{
			{Mark: "thème", Fill: "Thème : le dépôt WordPress, la page, et le champ ACF touchés."},
			{Mark: "xampp", Fill: "XAMPP : le changement est local tant que comms.allow n'est pas donné."},
		}
	case "azure":
		return []Craft{
			{Mark: "ressource azure", Fill: "Ressource Azure : le service ou le site branché, nommé."},
			{Mark: "site live", Fill: "Site live : lecture seule, preview, puis accord. Pas de réécriture directe."},
		}
	case "optimisation":
		return []Craft{
			{Mark: "rejet", Fill: "Rejet : l'option plus lourde, avec son score."},
			{Mark: "score", Fill: "Score : coût + risque + étapes. Le plus bas gagne."},
		}
	default:
		return nil
	}
}

// DeptByID returns one department.
func DeptByID(id string) (Department, bool) {
	for _, d := range Chart().Departments {
		if d.ID == id {
			return d, true
		}
	}
	return Department{}, false
}

// DeptByAgent returns the department that employs this agent.
func DeptByAgent(agent string) (Department, bool) {
	for _, d := range Chart().Departments {
		if d.Chief == agent {
			return d, true
		}
		for _, m := range d.Members {
			if m.Agent == agent {
				return d, true
			}
		}
	}
	return Department{}, false
}

func clip(draft string) string {
	s := strings.Join(strings.Fields(draft), " ")
	r := []rune(s)
	if len(r) > 90 {
		s = string(r[:90])
	}
	return s
}

// apply names the draft inside the missing mark, so the fill is about this
// action and not a generic sentence.
func (c Craft) apply(draft string) string {
	subject := clip(draft)
	if subject == "" {
		subject = "ce geste"
	}
	switch c.Mark {
	case "option rejetée":
		return "Option rejetée : ne pas rallonger « " + subject + " ». On ne fait pas le geste le plus long."
	case "responsable":
		return "Responsable : le chef du département tranche « " + subject + " »."
	case "deux façons":
		return "Deux façons : faire « " + subject + " » tel quel, ou le couper. On garde la plus courte."
	case "limite":
		return "Limite : « " + subject + " » ne réécrit pas un site déjà en ligne et ne déploie pas ce qui n'est pas prêt."
	case "offre":
		return "Offre : une seule, tirée de « " + subject + " », pour le client déjà nommé."
	case "mesure":
		return "Mesure : on regarde le résultat de « " + subject + " » avant d'élargir."
	case "livré":
		return "Livré : l'artifact de « " + subject + " », pas l'intention."
	case "pas fini":
		return "Pas fini : ce qui reste après « " + subject + " » avant de dire que c'est en ligne."
	case "règle":
		return "Règle : « " + subject + " » est permis ou interdit, cité, pas un avis."
	case "ça casse":
		return "Ça casse : le cas où « " + subject + " » perd de l'argent ou des données."
	case "ressource azure":
		return "Ressource Azure : le site ou le service nommé dans « " + subject + " »."
	case "site live":
		return "Site live : « " + subject + " » reste en lecture, preview, puis accord. Pas de réécriture directe."
	case "rejet":
		return "Rejet : l'option plus lourde que « " + subject + " », avec son score."
	case "score":
		return "Score : coût + risque + étapes de « " + subject + " ». Le plus bas gagne."
	default:
		return c.Fill
	}
}

// Improve makes the department fill whatever its specialty is missing.
// A second pass on the returned result finds no lacune.
func Improve(deptID, draft string) (Improvement, bool) {
	dept, ok := DeptByID(deptID)
	if !ok {
		return Improvement{}, false
	}
	out := strings.TrimSpace(draft)
	var gaps []string
	low := strings.ToLower(out)
	for _, c := range dept.Craft {
		if strings.Contains(low, strings.ToLower(c.Mark)) {
			continue
		}
		gaps = append(gaps, c.Mark)
		if out != "" {
			out += " "
		}
		out += c.apply(draft)
		low = strings.ToLower(out)
	}
	return Improvement{
		Department:  dept.ID,
		Lacunes:     gaps,
		Result:      out,
		Specialized: len(gaps) == 0,
	}, true
}
