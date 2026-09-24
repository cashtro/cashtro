package agents

import (
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// FusionName is the one name for both peoples. Astro is the cashtro
// account. Evolian is the Evolu-Jeunes organization. Saying the name
// means both, every time. There is no Fusion of one.
const FusionName = "The Fusion"

// People is one of the two GitHubs under its Fusion name.
type People struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}

// Peoples are Astro and Evolian. The order is the order of the two GitHubs.
func Peoples() []People {
	return []People{
		{ID: "astro", Name: "Astro", Account: "cashtro"},
		{ID: "evolian", Name: "Evolian", Account: "Evolu-Jeunes"},
	}
}

// Fusion returns both accounts. It never returns one.
func Fusion() []string {
	peoples := Peoples()
	out := make([]string, len(peoples))
	for i, p := range peoples {
		out[i] = p.Account
	}
	return out
}

// PeopleOf names the people that owns a repo. A bare account name counts.
func PeopleOf(repo string) (People, bool) {
	for _, p := range Peoples() {
		if repo == p.Account || strings.HasPrefix(repo, p.Account+"/") {
			return p, true
		}
	}
	return People{}, false
}

// InFusion reports whether a repo belongs to Astro or Evolian.
func InFusion(repo string) bool {
	_, ok := PeopleOf(repo)
	return ok
}

// FusionScope is every repo on the roster under both peoples, each once.
func FusionScope() []string {
	seen := map[string]bool{}
	var out []string
	add := func(repo string) {
		if repo == "" || seen[repo] || !InFusion(repo) {
			return
		}
		seen[repo] = true
		out = append(out, repo)
	}
	for _, ln := range Lines() {
		for _, repo := range ln.Repos {
			add(repo)
		}
	}
	for _, layer := range Layers() {
		add(layer.Repo)
	}
	return out
}

// Squad is one unit of the force. One squad per law, plus the guard that
// holds the kernel's own rules. Ten officers each.
type Squad struct {
	ID       string   `json:"id"`
	Law      string   `json:"law"`
	Name     string   `json:"name"`
	Officers int      `json:"officers"`
	Watches  []string `json:"watches"`
	Cites    string   `json:"cites"`
}

// Citation is one rule the force cites on one move. It names the word
// that tripped the squad, never the value around it.
type Citation struct {
	Squad string `json:"squad"`
	Law   string `json:"law"`
	Found string `json:"found"`
	Rule  string `json:"rule"`
}

// Patrol is one pass of the force over one move. The force stops the
// move and hands it to Epicenter Einstein. It does not punish, delete,
// push, or copy.
type Patrol struct {
	Repo         string     `json:"repo"`
	People       string     `json:"people"`
	Jurisdiction bool       `json:"jurisdiction"`
	Officers     int        `json:"officers"`
	Citations    []Citation `json:"citations"`
	Stopped      bool       `json:"stopped"`
	Applied      bool       `json:"applied"`
	HandedTo     string     `json:"handedTo"`
	Message      string     `json:"message"`
}

// FusionForce is the police inside The Eye. The Eye watches. The Fusion
// enforces the laws in Loi and the kernel's own rules across both
// peoples. It cites the rule and stops the move. Epicenter Einstein judges.
type FusionForce struct {
	Name          string   `json:"name"`
	Kind          string   `json:"kind"`
	Inside        string   `json:"inside"`
	Peoples       []People `json:"peoples"`
	Accounts      []string `json:"accounts"`
	Agents        int      `json:"agents"`
	Minimum       int      `json:"minimum"`
	Squads        []Squad  `json:"squads"`
	Enforces      []string `json:"enforces"`
	Scope         int      `json:"scope"`
	Cites         bool     `json:"cites"`
	Stops         bool     `json:"stops"`
	Judges        bool     `json:"judges"`
	Attacks       bool     `json:"attacks"`
	Deletes       bool     `json:"deletes"`
	Pushes        bool     `json:"pushes"`
	CopiesSecrets bool     `json:"copiesSecrets"`
	Weapons       bool     `json:"weapons"`
	ReportsTo     string   `json:"reportsTo"`
	OpenedBy      string   `json:"openedBy"`
	OrderedBy     string   `json:"orderedBy"`
	DecidedBy     string   `json:"decidedBy"`
}

// FusionMinimum is the floor the person set for the force.
const FusionMinimum = 80

const fusionGuardID = "garde"

// FusionSquads builds one squad per law in Loi, then the guard.
func FusionSquads() []Squad {
	watches := map[string][]string{
		"loi-25": {"renseignements personnels", "données personnelles", "liste de clients", "courriels des clients", "numéro d'assurance sociale", "sans efvp"},
		"pipeda": {"hors québec sans", "transfert hors canada", "serveur américain"},
		"casl":   {"pourriel", "spam", "envoi massif", "liste achetée", "sans désabonnement", "sans consentement"},
		"charte": {"anglais seulement", "english only", "sans version française", "pas de français"},
		"lpc":    {"prix caché", "fausse urgence", "résultat garanti", "gain garanti", "revenu garanti"},
		"image":  {"image de la base", "photo de la base", "gratter une photo", "scraper une photo", "photo sans droits"},
		"racj":   {"loterie", "tirage en ligne", "concours sans règlement"},
		"amf":    {"ordre live", "passer un ordre", "vendre le token", "vente du token", "sollicitation du public", "mint"},
	}
	squads := make([]Squad, 0, len(Loi().Laws)+1)
	for _, law := range Loi().Laws {
		cite := ""
		if len(law.Rules) > 0 {
			cite = law.Rules[0]
		}
		squads = append(squads, Squad{
			ID:       law.ID,
			Law:      law.Name,
			Name:     "Escouade " + law.Name,
			Officers: 10,
			Watches:  watches[law.ID],
			Cites:    cite,
		})
	}
	squads = append(squads, Squad{
		ID:       fusionGuardID,
		Law:      "Règles du kernel",
		Name:     "La Garde",
		Officers: 10,
		Watches:  []string{"arme", "weapon", "attaque", "attack", "flood", "ddos", "exploit", "force push", "push en production", "déployer en production", "envoyer sans allow"},
		Cites:    "Aucun déploiement, aucun envoi, aucun ordre de trading sans comms.allow. Pas d'arme, pas d'attaque, pas de secret recopié.",
	})
	return squads
}

// TheFusion returns the force roster.
func TheFusion() FusionForce {
	squads := FusionSquads()
	agents := 0
	enforces := make([]string, 0, len(squads))
	for _, squad := range squads {
		agents += squad.Officers
		enforces = append(enforces, squad.ID)
	}
	return FusionForce{
		Name:          FusionName,
		Kind:          "police",
		Inside:        "eye",
		Peoples:       Peoples(),
		Accounts:      Fusion(),
		Agents:        agents,
		Minimum:       FusionMinimum,
		Squads:        squads,
		Enforces:      enforces,
		Scope:         len(FusionScope()),
		Cites:         true,
		Stops:         true,
		Judges:        false,
		Attacks:       false,
		Deletes:       false,
		Pushes:        false,
		CopiesSecrets: false,
		Weapons:       false,
		ReportsTo:     "eye",
		OpenedBy:      "voltron",
		OrderedBy:     "scrum",
		DecidedBy:     "instinct",
	}
}

// FusionHolds reports whether the force is staffed and stays a police, not an army.
func FusionHolds() (bool, string) {
	force := TheFusion()
	if force.Agents < force.Minimum {
		return false, "effectif sous " + itoa(force.Minimum)
	}
	if len(force.Accounts) != 2 || force.Accounts[0] != "cashtro" || force.Accounts[1] != "Evolu-Jeunes" {
		return false, "la Fusion veut dire les deux peuples"
	}
	if force.Attacks || force.Deletes || force.Pushes || force.CopiesSecrets || force.Weapons || force.Judges {
		return false, "la police cite et arrête. Elle ne frappe pas et ne juge pas"
	}
	if force.DecidedBy != "instinct" || force.Inside != "eye" {
		return false, "la Fusion siège dans The Eye et Epicenter Einstein juge"
	}
	if len(force.Squads) != len(Loi().Laws)+1 {
		return false, "une escouade par loi, plus la garde"
	}
	return true, ""
}

// Enforce runs the force over one move. A repo outside both peoples is
// outside the jurisdiction and the move stops there. Inside, every
// squad reads the task. A secret or a forbidden word is a citation. The
// force stops the move and hands it to Epicenter Einstein. Nothing is
// applied, deleted, pushed, or copied.
func Enforce(repo, task string) Patrol {
	force := TheFusion()
	patrol := Patrol{
		Repo:     repo,
		Officers: force.Agents,
		Applied:  false,
		HandedTo: "instinct",
	}
	if repo == "" {
		patrol.People = "kernel"
		patrol.Jurisdiction = true
	} else if people, ok := PeopleOf(repo); ok {
		patrol.People = people.Name
		patrol.Jurisdiction = true
	} else {
		patrol.Jurisdiction = false
		patrol.Citations = append(patrol.Citations, Citation{
			Squad: fusionGuardID,
			Law:   "Règles du kernel",
			Found: "dépôt hors des deux peuples",
			Rule:  "Le coup reste dans la Fusion : Astro et Evolian. Un dépôt hors des deux GitHub n'est pas touché.",
		})
		patrol.Stopped = true
		patrol.Message = "hors juridiction. la Fusion arrête le coup"
		return patrol
	}
	text := strings.ToLower(task)
	for _, squad := range force.Squads {
		for _, word := range squad.Watches {
			if strings.Contains(text, word) {
				patrol.Citations = append(patrol.Citations, Citation{Squad: squad.ID, Law: squad.Law, Found: word, Rule: squad.Cites})
				break
			}
		}
	}
	for _, name := range secretNames(task) {
		patrol.Citations = append(patrol.Citations, Citation{
			Squad: fusionGuardID,
			Law:   "Règles du kernel",
			Found: name,
			Rule:  "Un secret ne se recopie pas. La valeur reste où elle est.",
		})
	}
	if hustlerRefused(task) {
		patrol.Citations = append(patrol.Citations, Citation{
			Squad: fusionGuardID,
			Law:   "Règles du kernel",
			Found: "chemin illégal",
			Rule:  "Aucun chemin illégal. Le coup s'arrête et Epicenter Einstein est saisi.",
		})
	}
	patrol.Stopped = len(patrol.Citations) > 0
	if patrol.Stopped {
		laws := make([]string, 0, len(patrol.Citations))
		for _, c := range patrol.Citations {
			laws = append(laws, c.Squad)
		}
		patrol.Message = "la Fusion arrête le coup : " + strings.Join(laws, ", ") + ". Epicenter Einstein juge"
		return patrol
	}
	patrol.Message = "patrouille passée. " + patrol.People + " dans la Fusion"
	return patrol
}

// FusionInvoke is manager.fusion. Without a repo or a task it returns the
// roster. With a move it patrols, tells Contrôle, and hands a stop to
// Epicenter Einstein. Nothing is applied from here.
func FusionInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	repo := payloadQuery(call, "repo")
	task := payloadQuery(call, "task")
	if task == "" {
		task = payloadQuery(call, "text")
	}
	if repo == "" && task == "" {
		ok, why := FusionHolds()
		force := TheFusion()
		msg := "la Fusion : " + itoa(force.Agents) + " agents, Astro et Evolian"
		if !ok {
			msg = why
		}
		return kernel.Result{OK: ok, Message: msg, Data: force}, nil
	}
	patrol := Enforce(repo, task)
	k.Remember("fusion", strings.TrimSpace(repo+" "+patrol.Message))
	_, _ = k.Post("manager", "security", "fusion", patrol.Message)
	_, _ = k.Post("manager", "investigator", "fusion", patrol.Message)
	if patrol.Stopped {
		_, _ = k.Post("manager", "reviewer", "fusion", patrol.Message)
	}
	return kernel.Result{OK: !patrol.Stopped, Message: patrol.Message, Data: patrol}, nil
}
