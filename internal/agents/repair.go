package agents

// ChainAgent is one seat on the agentic chain.
// The chain is this crew. A project, a website, and a repo are not seats.
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
		{1, "maitre", "Maître de chaîne", "La chaîne agentique, les workflows, et l'automatisation. Il tient cette chaîne à part des projets."},
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

// OnChain reports whether this id is a seat of the agentic chain.
// A project line, a website, and a repo are not on it.
func OnChain(id string) bool {
	for _, agent := range RepairCrew() {
		if agent.ID == id {
			return true
		}
	}
	return false
}

// SurveyChain is the agentic chain, built on its own.
// Giant and the trading bot keep their own chains and their own repos.
func SurveyChain() ChainReport {
	crew := RepairCrew()
	repair := []ChainFix{
		{1, "separee", "La chaîne des vingt sièges reste à part des projets."},
		{2, "sites", "Un site reste un site. Éduconnexion n'est pas une chaîne."},
		{3, "giant", "Giant est un écosystème. Son dépôt est Evolu-Jeunes/Giant. Sa chaîne reste."},
		{4, "bot", "Le bot de trading est cashtro/trading_bot-main. Sa chaîne reste. Les deux dépôts ne sont pas mélangés."},
		{5, "synchronie", "Une chaîne à part les fait travailler ensemble, y compris le trade blockchain. Un ordre live attend la réponse AMF écrite."},
	}
	return ChainReport{
		Lead:    crew[0].ID,
		Craft:   crew[0].Craft,
		Agents:  crew,
		Breaks:  nil,
		Repair:  repair,
		Decided: "instinct",
		Applied: false,
	}
}
