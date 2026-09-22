package agents

// Compliance is how the cooperative builds. It is an operating
// checklist, not a legal opinion. A Quebec lawyer confirms it before
// a live client site or a sale.
type Compliance struct {
	Disclaimer string `json:"disclaimer"`
	Officer    string `json:"officer"`
	Delegate   string `json:"delegate"`
	Laws       []Law  `json:"laws"`
	Gates      []Gate `json:"gates"`
}

// Law is one regime the build must respect.
type Law struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Where  string   `json:"where"`
	Source string   `json:"source"`
	Rules  []string `json:"rules"`
}

// Gate is one question Contrôle or the named owner must answer
// before a build is called done.
type Gate struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Owner    string `json:"owner"`
}

// Loi returns the Quebec and Canada working rules.
func Loi() Compliance {
	return Compliance{
		Disclaimer: "Checklist opérable, pas un avis juridique. Avant un site live ou une vente : avocat du Québec, et au besoin CAI, OQLF, RACJ ou AMF.",
		Officer:    "manager",
		Delegate:   "security",
		Laws: []Law{
			{
				ID: "loi-25", Name: "Loi 25", Where: "Québec",
				Source: "Loi sur la protection des renseignements personnels dans le secteur privé, RLRQ c. P-39.1",
				Rules: []string{
					"Le CEO est responsable par défaut (art. 3.1). La délégation à Contrôle se fait par écrit. Titre et coordonnées publiés.",
					"Politique de confidentialité, conservation, destruction, rôles du personnel, traitement des plaintes (art. 3.2).",
					"EFVP avant un système ou un service électronique qui touche des renseignements personnels, et avant une communication hors Québec (art. 3.3).",
					"Consentement manifeste, libre, éclairé, spécifique, distinct des conditions d'utilisation.",
					"Produit technologique offert au public : confidentialité au plus haut niveau par défaut (art. 9.1). Les témoins suivent le régime du consentement.",
					"Droits : accès, rectification, cessation de diffusion, portabilité.",
					"Incident de confidentialité : registre. Avis à la CAI et aux personnes s'il y a risque de préjudice sérieux.",
					"Décision automatisée : la personne est informée.",
					"Aucun renseignement personnel dans un prompt Ollama, Kimi ou GLM sans EFVP. Azure hors Québec = communication hors Québec.",
				},
			},
			{
				ID: "pipeda", Name: "LPRPDE (PIPEDA)", Where: "Canada, hors du champ couvert par la loi québécoise",
				Source: "Loi sur la protection des renseignements personnels et les documents électroniques",
				Rules: []string{
					"Activité commerciale fédérale, interprovinciale ou internationale : les principes de PIPEDA s'ajoutent.",
					"Au Québec, pour une entreprise provinciale, la Loi 25 est le régime principal. On ne choisit pas le plus faible.",
				},
			},
			{
				ID: "casl", Name: "LCAP (CASL)", Where: "Canada",
				Source: "Loi canadienne anti-pourriel",
				Rules: []string{
					"Message électronique commercial : consentement, identité de l'expéditeur, et mécanisme d'abonnement fonctionnel.",
					"Le CMP ne fait pas comms.allow sur une campagne qui rate un de ces trois points.",
				},
			},
			{
				ID: "charte", Name: "Charte de la langue française", Where: "Québec",
				Source: "Charte de la langue française, RLRQ c. C-11, et loi 96",
				Rules: []string{
					"Client au Québec : informé et servi en français.",
					"Site qui vend ou promeut au Québec : version française dans des conditions au moins aussi favorables.",
					"Contrat d'adhésion : version française remise ou accessible avant un choix exprès d'une autre langue. Le français n'est pas facturé.",
					"Factures, bons de commande et publications commerciales : en français.",
				},
			},
			{
				ID: "lpc", Name: "Protection du consommateur", Where: "Québec",
				Source: "Loi sur la protection du consommateur",
				Rules: []string{
					"Contrat à distance : divulgations exigées avant la conclusion, copie au consommateur.",
					"Deux versions linguistiques : l'interprétation la plus favorable au consommateur l'emporte.",
					"Pas de pratique interdite : prix caché, urgence fausse, résultat garanti.",
				},
			},
			{
				ID: "image", Name: "Image et photos", Where: "Québec",
				Source: "Code civil du Québec, art. 35 et suivants, et Loi 25",
				Rules: []string{
					"Photo d'une personne identifiable : consentement, ou une exception réelle. On ne vend pas cette image comme fiche.",
					"Photo de fiche : droit sur l'image détenu ou licencié. Pas de photo grattée.",
					"Donnée de scan qui identifie une personne : c'est un renseignement personnel. Pas de vente sans base légale.",
				},
			},
			{
				ID: "racj", Name: "Concours publicitaires", Where: "Québec",
				Source: "Loi sur les loteries et la Régie des alcools, des courses et des jeux",
				Rules: []string{
					"Un tirage promo n'est pas une loterie. Règlement, gratuité de participation sans achat obligatoire quand la loi l'exige, et règles RACJ avant publication.",
					"Aucun tirage en ligne tant que Contrôle n'a pas classé le dossier.",
				},
			},
			{
				ID: "amf", Name: "Valeurs mobilières", Where: "Québec et Canada",
				Source: "Autorité des marchés financiers, Lois sur les valeurs mobilières",
				Rules: []string{
					"Trading, token Giant, sollicitation du public : question d'inscription ou d'exemption avant tout ordre et avant toute vente de token.",
					"La division Trading reste fermée tant que cette question n'est pas répondue par écrit.",
				},
			},
		},
		Gates: []Gate{
			{ID: "renseignements", Question: "Ce build collecte, utilise ou communique des renseignements personnels?", Owner: "security"},
			{ID: "efvp", Question: "Si oui : EFVP faite et vue par le responsable avant le build?", Owner: "security"},
			{ID: "hors-quebec", Question: "Azure, analytics ou modèle (Kimi, GLM, Ollama distant) sortent-ils des renseignements du Québec?", Owner: "security"},
			{ID: "defaut", Question: "Confidentialité au plus haut par défaut, consentement séparé des conditions?", Owner: "architect"},
			{ID: "francais", Question: "Version française au moins aussi favorable pour un client au Québec?", Owner: "reviewer"},
			{ID: "casl", Question: "Courriel commercial : consentement, identité, désabonnement?", Owner: "comms"},
			{ID: "photos", Question: "Droits sur les photos réglés, aucune personne identifiable vendue sans consentement?", Owner: "reviewer"},
			{ID: "tirage", Question: "Si tirage : concours publicitaire classé, pas une loterie?", Owner: "security"},
			{ID: "trading", Question: "Si ordre ou token : réponse AMF écrite avant toute exécution?", Owner: "investigator"},
			{ID: "incident", Question: "Registre d'incidents et chemin d'avis à la CAI connus?", Owner: "security"},
		},
	}
}
