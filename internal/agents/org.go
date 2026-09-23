package agents

// Organization is the cooperative. Three chiefs, nine departments,
// the 15 real agentics as specialized employees. No extra processes.
// The central agency is another repo and another prompt. It is not here.
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
	Craft     []Craft  `json:"craft"`
	Members   []Member `json:"members"`
}

// Member is one of the 15 agentics, with a job and skills.
type Member struct {
	Agent  string   `json:"agent"`
	Title  string   `json:"title"`
	Skills []string `json:"skills"`
}

// Division is a profit or activity center. One chief, one function, one product.
type Division struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Chief      string   `json:"chief"`
	Department string   `json:"department"`
	Function   string   `json:"function"`
	Product    string   `json:"product"`
	Revenue    string   `json:"revenue"`
	Guard      string   `json:"guard"`
	Work       []string `json:"work,omitempty"`
}

// Chart returns the operating chart.
func Chart() Organization {
	org := Organization{
		Seats: []Seat{
			{ID: "ceo", Title: "CEO", Agent: "manager", Mandate: "Arbitrage, rentabilité, assignation. Stratège : une question et la position avant d'arbitrer. Rien ne sort sans son accord via comms."},
			{ID: "cto", Title: "CTO", Agent: "architect", Mandate: "Technique, chaînes, modèles. Stratège : deux façons, le coût, et la réponse adverse avant de construire."},
			{ID: "cmp", Title: "CMP", Agent: "comms", Mandate: "Marché. Vente, campagnes, porte de sortie. Stratège : une offre, la mesure, puis comms.allow."},
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
				ID: "wordpress", Name: "WordPress", ReportsTo: "cto", Chief: "operator",
				Mandate: "Epicenter et Voltron dirigent. Le commit reste sur les sites WordPress (PHP, ACF Pro, XAMPP) et sur Proximity, Next Proximity, Proximity App, et Api-Proximity. Lovable n'y touche pas.",
				Members: []Member{
					{Agent: "operator", Title: "confection WordPress", Skills: []string{"PHP", "ACF Pro", "XAMPP", "thèmes"}},
					{Agent: "memory", Title: "dossier de chaque site", Skills: []string{"contexte", "pages", "champs ACF"}},
				},
			},
			{
				ID: "azure", Name: "Azure", ReportsTo: "cto", Chief: "deploy",
				Mandate: "L'autre équipe de la même maison. Applications et pages qui ne sont pas des thèmes WordPress. Elle ne touche pas PHP, ACF Pro, ni XAMPP. Elle ne scanne pas et elle ne vend pas le marketplace.",
				Members: []Member{
					{Agent: "deploy", Title: "chef d'équipe Azure", Skills: []string{"Azure", "preview", "release"}},
					{Agent: "operator", Title: "confection hors thème", Skills: []string{"applications", "pages", "comptes hors WordPress"}},
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
			{
				ID: "flux", Name: "Flow stack", ReportsTo: "cto", Chief: "architect",
				Mandate: "La technologie qui fait tourner les workers, pas un résultat ajouté après. Chaque division a une fonction et un produit. Voltron tient le kernel. Epicenter lance le cycle.",
				Members: []Member{
					{Agent: "architect", Title: "chef du flux", Skills: []string{"division", "fonction", "produit"}},
					{Agent: "operator", Title: "worker confection", Skills: []string{"exécution locale"}},
					{Agent: "security", Title: "worker filtre", Skills: []string{"secrets", "Loi 25"}},
					{Agent: "reviewer", Title: "worker second filtre", Skills: []string{"option courte"}},
					{Agent: "investigator", Title: "worker trace", Skills: []string{"rayon"}},
					{Agent: "deploy", Title: "worker sortie", Skills: []string{"preview", "allow"}},
				},
			},
		},
		Divisions: []Division{
			{ID: "wordpress", Name: "WordPress Proximity", Chief: "operator", Department: "wordpress", Function: "operator.work", Product: "Sites WordPress, Proximity, Next Proximity, Proximity App, Api-Proximity", Revenue: "Sites WordPress, Proximity, Next Proximity, Proximity App, Api-Proximity.", Guard: "Le commit ne sort pas de ces dépôts. Security filtre, reviewer filtre encore, comms.allow avant de quitter le local. Lovable n'entre pas."},
			{ID: "proximity", Name: "Proximity Agency", Chief: "deploy", Department: "azure", Function: "manager.line", Product: "Comptes hors WordPress et hors produits Proximity nommés", Revenue: "L'autre équipe : comptes qui ne sont pas le WordPress ni les produits Proximity nommés.", Guard: "Ne pas committer ces dépôts sur le département WordPress, ni l'inverse."},
			{
				ID: "scanapp", Name: "Scan App", Chief: "comms", Department: "marche", Function: "operator.work", Product: "Stock du scanner interne",
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
					"Le stock vient de inventory-scanner-theta.vercel.app. Le repo est le marketplace, qui vend via MCP et Stripe. Empire vend le même stock en live. Proximity n'est pas dans cette chaîne.",
				},
			},
			{ID: "panda", Name: "Panda", Chief: "architect", Department: "ingenierie", Function: "architect.plan", Product: "Evolu-Jeunes/Panda", Revenue: "White-glove : IA, installation, cours, entreprises et particuliers.", Guard: "On vend le service. On ne déploie pas chez le client sans allow."},
			{ID: "nft-giant", Name: "NFT + Giant", Chief: "manager", Department: "controle", Function: "manager.run", Product: "Evolu-Jeunes/Giant brain", Revenue: "Art avec utilité, token Giant, lecture de marché.", Guard: "Epicenter décide en dernier. Pas de mint, pas d'ordre, pas de pont."},
			{ID: "ecole", Name: "École", Chief: "research", Department: "marche", Function: "research.ingest", Product: "Evolu-Jeunes/educonnexion", Revenue: "Cours tech (WordPress, Web3).", Guard: "Le cours décrit. Il ne promet pas un gain de trading."},
			{ID: "marketing", Name: "Marketing", Chief: "planner", Department: "marche", Function: "manager.crm", Product: "Evolu-Jeunes/CRM-Agents", Revenue: "L'agence: marketing, Panda, Proximity cloud. Elle relie chaque chaîne.", Guard: "Les leads passent dans la copie agents. La copie Fix2 n'y entre pas. Lovable n'entre pas dans WordPress. Epicenter décide avant un envoi."},
			{ID: "fix2", Name: "Fix2", Chief: "operator", Department: "operations", Function: "operator.work", Product: "Evolu-Jeunes/CRM", Revenue: "Homme à tout faire, résidentiel. Le CRM est le carnet Lovable.", Guard: "Pas le carnet des agents. Pas WordPress. Pas d'envoi, pas de paiement, pas de permis déposé. Epicenter décide."},
			{ID: "empire", Name: "Empire Media", Chief: "delivery", Department: "operations", Function: "delivery.advance", Product: "Evolu-Jeunes/EmpireMedia", Revenue: "Live sell du stock Scan App, sur toutes les plateformes live.", Guard: "Le live vend une fiche déjà scannée. Pas de stock inventé. Proximity ne publie pas ces fiches."},
			{ID: "trading", Name: "Trading", Chief: "investigator", Department: "controle", Function: "investigator.trace", Product: "Bots et Giant", Revenue: "Bots sur exchanges, Web3, liés à Giant.", Guard: "Aucun ordre et aucune vente de token sans réponse AMF écrite, puis comms.allow."},
			{ID: "fonds", Name: "Fonds", Chief: "planner", Department: "direction", Function: "manager.tax", Product: "Grand livre interne", Revenue: "Position, chèques en brouillon, TPS et TVQ.", Guard: "Pas une banque. Pas une déclaration. Rien ne sort sans allow."},
		},
	}
	for i := range org.Departments {
		org.Departments[i].Craft = craftFor(org.Departments[i].ID)
	}
	return org
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
