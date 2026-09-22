package agents

// Organization is the cooperative. Three chiefs, six departments,
// the 15 real agentics as specialized employees. No extra processes.
type Organization struct {
	Seats       []Seat       `json:"seats"`
	Departments []Department `json:"departments"`
	Divisions   []Division   `json:"divisions"`
}

// Seat is a chief of the cooperative.
type Seat struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Agent   string `json:"agent"`
	Mandate string `json:"mandate"`
}

// Department is a team with one chief and specialized employees.
type Department struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	ReportsTo string   `json:"reportsTo"`
	Chief     string   `json:"chief"`
	Mandate   string   `json:"mandate"`
	Members   []Member `json:"members"`
}

// Member is one of the 15 agentics, with a job and skills.
type Member struct {
	Agent  string   `json:"agent"`
	Title  string   `json:"title"`
	Skills []string `json:"skills"`
}

// Division is a profit or activity center. One chief, one revenue, one guard.
type Division struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Chief      string `json:"chief"`
	Department string `json:"department"`
	Revenue    string   `json:"revenue"`
	Guard      string   `json:"guard"`
	Work       []string `json:"work,omitempty"`
}

// Chart returns the operating chart.
func Chart() Organization {
	return Organization{
		Seats: []Seat{
			{ID: "ceo", Title: "CEO", Agent: "manager", Mandate: "Arbitrage, rentabilité, assignation. Rien ne sort sans son accord via comms."},
			{ID: "cto", Title: "CTO", Agent: "architect", Mandate: "Technique, chaînes, modèles. Décide comment c'est construit, pas ce qui est vendu."},
			{ID: "cmp", Title: "CMP", Agent: "comms", Mandate: "Marché. Vente, campagnes, porte de sortie. comms.allow est la caisse."},
		},
		Departments: []Department{
			{
				ID: "direction", Name: "Direction", ReportsTo: "ceo", Chief: "manager",
				Mandate: "Arbitre seulement après qu'Optimisation a contredit le choix. Une décision sans option rejetée n'est pas une décision.",
				Members: []Member{
					{Agent: "manager", Title: "CEO", Skills: []string{"arbitrage", "rentabilité", "assignation"}},
					{Agent: "init", Title: "secrétaire général", Skills: []string{"boot", "état du système"}},
				},
			},
			{
				ID: "ingenierie", Name: "Ingénierie", ReportsTo: "cto", Chief: "architect",
				Mandate: "Construit. Pose au moins deux façons de faire, avec le coût. Ne se note pas elle-même.",
				Members: []Member{
					{Agent: "architect", Title: "CTO", Skills: []string{"conception", "chaînes", "offres"}},
					{Agent: "router", Title: "employé spécialisé modèles", Skills: []string{"coût des modèles", "Ollama", "Kimi K3", "GLM max"}},
					{Agent: "memory", Title: "employé spécialisé dossier", Skills: []string{"mémoire", "fiches"}},
				},
			},
			{
				ID: "marche", Name: "Marché", ReportsTo: "cmp", Chief: "comms",
				Mandate: "Vend une offre à la fois. Chaque campagne a une variante plus courte. On mesure avant d'élargir.",
				Members: []Member{
					{Agent: "comms", Title: "CMP", Skills: []string{"vente", "campagnes", "comms.allow"}},
					{Agent: "planner", Title: "chef de département planification", Skills: []string{"backlog", "priorité rentable"}},
					{Agent: "research", Title: "employé spécialisé veille", Skills: []string{"école", "offre", "veille"}},
				},
			},
			{
				ID: "operations", Name: "Opérations", ReportsTo: "cto", Chief: "delivery",
				Mandate: "Livre. Ne déclare pas le travail fini. Optimisation teste le résultat contre un chemin plus court.",
				Members: []Member{
					{Agent: "delivery", Title: "chef de département livraison", Skills: []string{"idée", "concept", "production"}},
					{Agent: "operator", Title: "employé spécialisé exécution", Skills: []string{"confection", "browser"}},
					{Agent: "deploy", Title: "employé spécialisé release", Skills: []string{"preview", "release"}},
				},
			},
			{
				ID: "controle", Name: "Contrôle", ReportsTo: "ceo", Chief: "security",
				Mandate: "Dit si c'est permis et si ça peut casser. Ne choisit pas la meilleure option : c'est le métier d'Optimisation.",
				Members: []Member{
					{Agent: "security", Title: "chef de département risque", Skills: []string{"CVE", "secrets", "triage"}},
					{Agent: "investigator", Title: "employé spécialisé incidents", Skills: []string{"trace", "perte", "blast radius"}},
				},
			},
			{
				ID: "optimisation", Name: "Optimisation", ReportsTo: "ceo", Chief: "reviewer",
				Mandate: "Contredit et teste tout ce que les autres départements proposent. Ne s'arrête que sur l'option la plus optimisée : moins d'étapes, moins de coût, moins de risque.",
				Members: []Member{
					{Agent: "reviewer", Title: "chef de département, teste et contredit", Skills: []string{"contre-épreuve", "test", "option rejetée"}},
					{Agent: "explorer", Title: "employé spécialisé option plus courte", Skills: []string{"Graphify", "chemin existant", "moindre coût"}},
				},
			},
		},
		Divisions: []Division{
			{ID: "proximity", Name: "Proximity Agency", Chief: "operator", Department: "operations", Revenue: "Sites clients WordPress / PHP / ACF Pro.", Guard: "Les sites déjà en ligne se lisent. Preview seulement, puis comms.allow."},
			{
				ID: "scanapp", Name: "Scan App", Chief: "comms", Department: "marche",
				Revenue: "Le scan tient le stock vrai. La fiche vend ce stock. La promo et le média tournent autour de la fiche, pas autour de la photo seule.",
				Guard:   "On ne publie pas l'image brute de la base. Photo générée ou trouvée, avec les droits. Une personne identifiable n'est pas un produit.",
				Work: []string{
					"Scanner interchangeable (caméra, USB, ZXing, BarcodeDetector) : code, format, source.",
					"Le scan écrit le stock : ajout, retrait, vente, vérification. Produit = code unique, sku, prix, quantité, emplacement.",
					"La base garde souvent une image mauvaise ou vide. Ce n'est pas la fiche.",
					"Fiche de vente : nom, prix, code, lieu, description, à partir du produit scanné.",
					"Photos : en générer une ou en trouver une. Jamais l'image mal traitée de la base sur la fiche publique.",
					"Boutique : le stock disponible devient catalogue commandable. Le flux commande vendeur n'est pas fini dans le repo.",
					"Autour : vendre, promouvoir, marketing et média, automatisés sur ces fiches.",
					"Proximity publie la fiche sur les sites. Empire la vend en live. Ils ne réinventent pas le stock.",
				},
			},
			{ID: "panda", Name: "Panda", Chief: "architect", Department: "ingenierie", Revenue: "White-glove : IA, installation, cours, entreprises et particuliers.", Guard: "On vend le service. On ne déploie pas chez le client sans allow."},
			{ID: "nft-giant", Name: "NFT + Giant", Chief: "security", Department: "controle", Revenue: "Art avec utilité, token Giant, tirages de visibilité.", Guard: "Pas de mint ni de campagne sans revue risque et allow."},
			{ID: "ecole", Name: "École", Chief: "research", Department: "marche", Revenue: "Cours tech (WordPress, Web3).", Guard: "Le cours décrit. Il ne promet pas un gain de trading."},
			{ID: "marketing", Name: "Marketing", Chief: "planner", Department: "marche", Revenue: "Gestion marketing des business corporate.", Guard: "Une campagne à la fois. Mesurer avant d'élargir."},
			{ID: "empire", Name: "Empire Media", Chief: "delivery", Department: "operations", Revenue: "Live sell sur Kick : podcast, UGC, items du scan.", Guard: "Le live vend une fiche déjà scannée. Pas de stock inventé. Kick diffuse, il ne remplace pas le stock."},
			{ID: "trading", Name: "Trading", Chief: "investigator", Department: "controle", Revenue: "Bots sur exchanges, Web3, liés à Giant.", Guard: "Aucun ordre et aucune vente de token sans réponse AMF écrite, puis comms.allow."},
		},
	}
}

// Members returns each agentic once, in department order.
func (o Organization) Members() []Member {
	var out []Member
	seen := map[string]bool{}
	for _, d := range o.Departments {
		for _, m := range d.Members {
			if seen[m.Agent] {
				continue
			}
			seen[m.Agent] = true
			out = append(out, m)
		}
	}
	return out
}
