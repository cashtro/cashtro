package agents

import "strings"

// ChainAgent is one seat on the repair crew.
// The crew is not a kernel process. Boot stays at fifteen.
type ChainAgent struct {
	Order int    `json:"order"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

// ChainBreak is one fault the crew can show from the live roster.
type ChainBreak struct {
	ID    string `json:"id"`
	Where string `json:"where"`
	Fact  string `json:"fact"`
}

// ChainFix is one step of the repair. Einstein decides before it is applied.
type ChainFix struct {
	Order int    `json:"order"`
	Break string `json:"break"`
	Do    string `json:"do"`
}

// ChainReport is the whole chaîne, read by the maître.
type ChainReport struct {
	Lead    string       `json:"lead"`
	Craft   string       `json:"craft"`
	Agents  []ChainAgent `json:"agents"`
	Breaks  []ChainBreak `json:"breaks"`
	Repair  []ChainFix   `json:"repair"`
	Decided string       `json:"decidedBy"`
	Applied bool         `json:"applied"`
}

// RepairCrew is twenty agents. The first is the highest seat for chains, workflows, and automation.
func RepairCrew() []ChainAgent {
	return []ChainAgent{
		{1, "maitre", "Maître de chaîne", "Chaînes agentiques, workflows, et automatisation. Il lit toute la chaîne et signe la réparation."},
		{2, "releveur", "Releveur", "Lister chaque ligne, chaque cerveau, et chaque couche Voltron."},
		{3, "compteur", "Compteur", "Compter les pas et refuser un trou dans l'ordre."},
		{4, "proprietaire", "Propriétaire", "Un dépôt n'a qu'une ligne propriétaire."},
		{5, "frontiere", "Frontière", "Un site avec son CRM reste hors de l'agentique."},
		{6, "nom", "Nom", "Un nom de dépôt n'a qu'un sens."},
		{7, "flux", "Flux", "Le pas N+1 suit le pas N. Un agent du pas est un des quinze."},
		{8, "voix", "Voix", "Vapi est le service client de Panda, et de personne d'autre."},
		{9, "carnet", "Carnet", "Le carnet des agents ne reçoit pas une fiche du site."},
		{10, "cerveau", "Cerveau", "Chaque cerveau se ferme sur Einstein."},
		{11, "voltron", "Voltron", "Chaque couche rapporte à Voltron et n'est pas un site."},
		{12, "automate", "Automate", "L'automatisation suit la chaîne. Elle ne la saute pas."},
		{13, "loi", "Loi", "Un geste qui rate la loi s'arrête avant le suivant."},
		{14, "brouillon", "Brouillon", "Un envoi, un appel, et un paiement restent un brouillon."},
		{15, "critique", "Critique", "Écrire la faute avec l'endroit où elle est."},
		{16, "plus-court", "Plus court", "Garder la réparation la plus courte qui referme la faute."},
		{17, "reparation", "Réparation", "Écrire les pas, dans l'ordre, pour toute la chaîne."},
		{18, "sceau", "Sceau", "Sceller le plan. Ne pas déployer."},
		{19, "relais", "Relais", "Remettre le plan au hustler, qui le relaie à Einstein."},
		{20, "einstein", "Einstein", "Trancher. Tant qu'il n'a pas tranché, le plan n'est pas appliqué."},
	}
}

// SurveyChain reads the live roster and names what is broken.
func SurveyChain() ChainReport {
	crew := RepairCrew()
	var breaks []ChainBreak
	known := map[string]bool{}
	for _, id := range AgentIDs() {
		known[id] = true
	}
	owner := map[string]string{}
	for _, ln := range Lines() {
		steps, ok := Chain(ln.ID)
		if !ok || len(steps) < 2 {
			breaks = append(breaks, ChainBreak{ID: "chaine-courte", Where: ln.ID, Fact: "La ligne n'a pas une chaîne d'au moins deux pas."})
		}
		for i, step := range steps {
			if step.Order != i+1 {
				breaks = append(breaks, ChainBreak{ID: "ordre", Where: ln.ID, Fact: "L'ordre des pas n'est pas 1, 2, 3."})
				break
			}
			if !known[step.Agent] {
				breaks = append(breaks, ChainBreak{ID: "agent-hors-table", Where: ln.ID + "/" + step.Agent, Fact: "Ce pas appelle un agent qui n'est pas dans les quinze."})
			}
			if ln.ID != "panda" && strings.Contains(strings.ToLower(step.Do), "vapi") {
				breaks = append(breaks, ChainBreak{ID: "vapi-hors-panda", Where: ln.ID, Fact: "Vapi est nommé hors du service client de Panda."})
			}
		}
		if ln.ID == "fix2" {
			breaks = append(breaks, ChainBreak{ID: "site-dans-la-chaine", Where: "fix2", Fact: "Fix Tout est un site avec un CRM. Il est encore une ligne de l'agentique."})
		}
		for _, repo := range ln.Repos {
			if prev, ok := owner[repo]; ok {
				breaks = append(breaks, ChainBreak{ID: "repo-deux-lignes", Where: repo, Fact: prev + " et " + ln.ID + " tiennent le même dépôt, donc deux chaînes."})
				continue
			}
			owner[repo] = ln.ID
		}
	}
	if copy := CRM().Copy; copy != "" {
		for _, ln := range Lines() {
			for _, repo := range ln.Repos {
				if repo == copy {
					breaks = append(breaks, ChainBreak{ID: "nom-double", Where: copy, Fact: "Ce nom est la copie du carnet des agents et un dépôt de " + ln.ID + "."})
				}
			}
		}
	}
	if ok, loose := VoltronHolds(); !ok {
		breaks = append(breaks, ChainBreak{ID: "voltron", Where: loose, Fact: "Une couche ne rapporte pas à Voltron."})
	}
	repair := []ChainFix{
		{1, "repo-deux-lignes", "Un dépôt a une ligne propriétaire. L'autre ligne peut le nommer, et sa chaîne ne commit pas dessus."},
		{2, "site-dans-la-chaine", "Retirer Fix Tout des lignes, de Chain, et de la division. Garder le CRM du site et le mur: aucune fiche ne passe dans le carnet."},
		{3, "nom-double", "Evolu-Jeunes/CRM ne veut dire que la copie du carnet des agents. Il quitte la liste de dépôts de Fix Tout."},
		{4, "vapi-hors-panda", "Vapi reste le service client de Panda sur db:panda. Toute autre chaîne qui le nomme est coupée."},
		{5, "agent-hors-table", "Un pas de chaîne n'appelle que l'un des quinze agents du kernel. Cette équipe de vingt reste hors de la table de boot."},
		{6, "ordre", "Récrire la chaîne pour que les pas se suivent à partir de 1, sans trou."},
		{7, "chaine-courte", "Chaque ligne garde au moins deux pas, et le dernier remet la décision à Einstein."},
		{8, "voltron", "Chaque couche rapporte à Voltron, reste opérante, et n'est pas un site."},
		{9, "site-dans-la-chaine", "Les agents branchés sur le site automatisent ce CRM-là. Ils ne prennent pas un ID dans l'agentique."},
		{10, "maitre", "Le maître relit toute la chaîne après ces pas. Einstein tranche. Rien n'est déployé par ce rapport."},
	}
	return ChainReport{
		Lead:    crew[0].ID,
		Craft:   crew[0].Craft,
		Agents:  crew,
		Breaks:  breaks,
		Repair:  repair,
		Decided: "instinct",
		Applied: false,
	}
}
