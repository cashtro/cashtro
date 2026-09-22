package agents

// Organization is the cooperative. Three chiefs, five departments,
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
	Revenue    string `json:"revenue"`
	Guard      string `json:"guard"`
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
				Members: []Member{
					{Agent: "manager", Title: "CEO", Skills: []string{"arbitrage", "rentabilité", "assignation"}},
					{Agent: "init", Title: "secrétaire général", Skills: []string{"boot", "état du système"}},
				},
			},
			{
				ID: "ingenierie", Name: "Ingénierie", ReportsTo: "cto", Chief: "architect",
				Members: []Member{
					{Agent: "architect", Title: "CTO", Skills: []string{"conception", "chaînes", "offres"}},
					{Agent: "router", Title: "employé spécialisé modèles", Skills: []string{"coût des modèles", "Ollama", "Kimi K3", "GLM max"}},
					{Agent: "explorer", Title: "employé spécialisé carte", Skills: []string{"Graphify", "recherche interne"}},
					{Agent: "memory", Title: "employé spécialisé dossier", Skills: []string{"mémoire", "fiches"}},
				},
			},
			{
				ID: "marche", Name: "Marché", ReportsTo: "cmp", Chief: "comms",
				Members: []Member{
					{Agent: "comms", Title: "CMP", Skills: []string{"vente", "campagnes", "comms.allow"}},
					{Agent: "planner", Title: "chef de département planification", Skills: []string{"backlog", "priorité rentable"}},
					{Agent: "research", Title: "employé spécialisé veille", Skills: []string{"école", "offre", "veille"}},
				},
			},
			{
				ID: "operations", Name: "Opérations", ReportsTo: "cto", Chief: "delivery",
				Members: []Member{
					{Agent: "delivery", Title: "chef de département livraison", Skills: []string{"idée", "concept", "production"}},
					{Agent: "operator", Title: "employé spécialisé exécution", Skills: []string{"confection", "browser"}},
					{Agent: "deploy", Title: "employé spécialisé release", Skills: []string{"preview", "release"}},
					{Agent: "reviewer", Title: "employé spécialisé qualité", Skills: []string{"revue", "artifacts"}},
				},
			},
			{
				ID: "controle", Name: "Contrôle", ReportsTo: "ceo", Chief: "security",
				Members: []Member{
					{Agent: "security", Title: "chef de département risque", Skills: []string{"CVE", "secrets", "triage"}},
					{Agent: "investigator", Title: "employé spécialisé incidents", Skills: []string{"trace", "perte", "blast radius"}},
				},
			},
		},
		Divisions: []Division{
			{ID: "proximity", Name: "Proximity Agency", Chief: "operator", Department: "operations", Revenue: "Sites clients WordPress / PHP / ACF Pro.", Guard: "Les sites déjà en ligne se lisent. Preview seulement, puis comms.allow."},
			{ID: "scanapp", Name: "Scan App", Chief: "comms", Department: "marche", Revenue: "Fiches vendues à partir de la donnée scan, avec de vraies photos.", Guard: "Pas de fiche publiée sans photo propre et sans allow."},
			{ID: "panda", Name: "Panda", Chief: "architect", Department: "ingenierie", Revenue: "White-glove : IA, installation, cours, entreprises et particuliers.", Guard: "On vend le service. On ne déploie pas chez le client sans allow."},
			{ID: "nft-giant", Name: "NFT + Giant", Chief: "security", Department: "controle", Revenue: "Art avec utilité, token Giant, tirages de visibilité.", Guard: "Pas de mint ni de campagne sans revue risque et allow."},
			{ID: "ecole", Name: "École", Chief: "research", Department: "marche", Revenue: "Cours tech (WordPress, Web3).", Guard: "Le cours décrit. Il ne promet pas un gain de trading."},
			{ID: "marketing", Name: "Marketing", Chief: "planner", Department: "marche", Revenue: "Gestion marketing des business corporate.", Guard: "Une campagne à la fois. Mesurer avant d'élargir."},
			{ID: "empire", Name: "Empire Media", Chief: "delivery", Department: "operations", Revenue: "Live sell des items, podcast, UGC.", Guard: "Le live vend ce qui est déjà une fiche. Pas de stock inventé."},
			{ID: "trading", Name: "Trading", Chief: "investigator", Department: "controle", Revenue: "Bots sur exchanges, Web3, liés à Giant.", Guard: "Aucun ordre live sans comms.allow. Le chef trace les pertes, il ne pousse pas le volume."},
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
