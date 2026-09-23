package agents

import (
	"fmt"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// CRMBook is the agent growth book. Fix2 does not write here.
// The copy under Proximity is not a second book.
// Voice is Vapi. A call is not placed from this process.
type CRMBook struct {
	Repo    string   `json:"repo"`
	App     string   `json:"app"`
	Host    string   `json:"host"`
	Copy    string   `json:"copy"`
	Voice   string   `json:"voice"`
	Gathers []string `json:"gathers"`
	Modules []string `json:"modules"`
	Lines   []string `json:"lines"`
	Parked  []string `json:"parked"`
	Next    []string `json:"next"`
	Writes  bool     `json:"writes"`
	Sent    bool     `json:"sent"`
	Dialed  bool     `json:"dialed"`
}

// CRMGrowth is one draft note. It is remembered. It is not written to Supabase and it is not sent.
type CRMGrowth struct {
	Line       string `json:"line"`
	Kind       string `json:"kind"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Status     string `json:"status"`
	Sent       bool   `json:"sent"`
	DecidedBy  string `json:"decidedBy"`
	Remembered bool   `json:"remembered"`
}

var crmKinds = map[string]bool{"contact": true, "lead": true, "note": true}

// What each line files. Every line is in the book. A theme edit, an order, and a send stay parked.
var crmLines = map[string]string{
	"control":     "le tableau, pas un produit",
	"wordpress":   "le nom du site et du client, pas une édition du thème",
	"proximity":   "le compte hors thème",
	"scanapp":     "une fiche de stock, pas une publication",
	"marketplace": "une vente notée, pas un checkout d'ici",
	"panda":       "note client ou partenaire",
	"nft-giant":   "une note d'art ou de token, pas un ordre",
	"ecole":       "une note de cours",
	"marketing":   "campagne, contact, lead",
	"empire":      "une note de live",
	"propres":     "note Pandora, pas un encaissement",
	"trading":     "une note de modèle, pas un ordre",
	"fonds":       "une note de livre, pas un virement",
}

// CRM is the agent book. Fix2 is absent on purpose.
func CRM() CRMBook {
	lines := make([]string, 0, len(Lines()))
	for _, ln := range Lines() {
		if ln.ID == "fix2" {
			continue
		}
		lines = append(lines, ln.ID)
	}
	return CRMBook{
		Repo:    "Evolu-Jeunes/CRM",
		App:     "Itercore",
		Host:    "itercore.lovable.app",
		Copy:    "Evolu-Jeunes/Proximity/apps/CRM",
		Voice:   "vapi",
		Gathers: []string{"contact", "lead", "note"},
		Modules: []string{"contacts", "leads", "notes", "voix"},
		Lines:   lines,
		Parked:  []string{"envoi", "appel", "écriture Supabase", "édition d'un thème", "ordre", "fiche Fix2"},
		Next: []string{
			"Les agents grandissent dans ce carnet. Fix2 n'y écrit pas.",
			"Vapi ouvre un projet à la fois. Chaque projet a sa propre base.",
			"La voix de ce carnet est Vapi. Aucun appel ne part d'ici.",
			"La copie dans Proximity n'est pas un second carnet.",
			"Un brouillon reste en mémoire. Supabase n'est pas écrit d'ici.",
			"Une ligne WordPress peut noter le client. Elle n'édite pas le thème.",
			"Un envoi, un appel ou un ordre attend Epicenter.",
		},
		Writes: false,
		Sent:   false,
		Dialed: false,
	}
}

// CRMGrow files one draft on the book. send is always refused.
func CRMGrow(k *kernel.Kernel, line, kind, title, body string, send bool) (CRMGrowth, error) {
	line = strings.TrimSpace(line)
	kind = strings.TrimSpace(kind)
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if send {
		return CRMGrowth{}, fmt.Errorf("envoi refusé")
	}
	if line == "fix2" {
		return CRMGrowth{}, fmt.Errorf("Fix2 tient son propre carnet")
	}
	why, ok := crmLines[line]
	if !ok {
		if _, known := LineByID(line); !known {
			return CRMGrowth{}, fmt.Errorf("%s n'est pas une ligne du conglomerat", line)
		}
		why = "note du conglomerat"
	}
	if !crmKinds[kind] {
		return CRMGrowth{}, fmt.Errorf("kind requis: contact, lead, note")
	}
	if title == "" {
		return CRMGrowth{}, fmt.Errorf("titre requis")
	}
	if names := secretNames(title + "\n" + body); len(names) > 0 {
		return CRMGrowth{}, fmt.Errorf("secret filtré: %s", strings.Join(names, ", "))
	}
	note := CRMGrowth{
		Line:      line,
		Kind:      kind,
		Title:     title,
		Body:      body,
		Status:    "brouillon",
		Sent:      false,
		DecidedBy: "epicenter",
	}
	if k != nil {
		k.Remember("crm:"+line, kind+" · "+title+" · "+why)
		note.Remembered = true
		k.Publish("memory", "epicenter", "crm brouillon "+line, map[string]any{
			"kind": kind, "title": title, "sent": false, "decidedBy": "epicenter",
		})
	}
	return note, nil
}

// CRMInvoke is manager.crm. An empty title returns the book. A title files a draft.
func CRMInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	title := payloadQuery(call, "title")
	if title == "" {
		return kernel.Result{OK: true, Message: "carnet des agents, voix Vapi, pas Fix2", Data: CRM()}, nil
	}
	note, err := CRMGrow(k, payloadQuery(call, "line"), payloadQuery(call, "kind"), title, payloadQuery(call, "body"), payloadQuery(call, "send") == "true")
	if err != nil {
		return kernel.Result{OK: false, Message: err.Error()}, nil
	}
	return kernel.Result{OK: true, Message: "brouillon CRM, rien n'est envoyé", Data: note}, nil
}
