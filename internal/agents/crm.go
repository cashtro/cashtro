package agents

import (
	"fmt"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// CRMBook is the one Lovable CRM the brains may use.
// The copy under Proximity is not a second book.
type CRMBook struct {
	Repo     string   `json:"repo"`
	App      string   `json:"app"`
	Host     string   `json:"host"`
	Copy     string   `json:"copy"`
	Gathers  []string `json:"gathers"`
	Modules  []string `json:"modules"`
	Lines    []string `json:"lines"`
	Closed   []string `json:"closed"`
	Next     []string `json:"next"`
	Writes   bool     `json:"writes"`
	Sent     bool     `json:"sent"`
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

// Lines that may file a draft. WordPress, Giant, and trading stay out.
var crmLines = map[string]string{
	"marketing": "campagne, contact, lead",
	"panda":     "note client ou partenaire white-glove",
	"propres":   "note Pandora, pas un encaissement",
	"proximity": "nom du client seulement",
}

// CRM is the single book and the next steps.
func CRM() CRMBook {
	return CRMBook{
		Repo:    "Evolu-Jeunes/CRM",
		App:     "Itercore",
		Host:    "itercore.lovable.app",
		Copy:    "Evolu-Jeunes/Proximity/apps/CRM",
		Gathers: []string{"contact", "lead", "note"},
		Modules: []string{"contacts", "leads", "estimates", "operations", "subs", "voice", "portal"},
		Lines:   []string{"marketing", "panda", "propres", "proximity"},
		Closed:  []string{"wordpress", "nft-giant", "trading", "empire", "ecole", "scanapp", "marketplace", "fonds"},
		Next: []string{
			"Un seul carnet: Evolu-Jeunes/CRM. La copie dans Proximity ne compte pas.",
			"Les agents déposent un brouillon en mémoire. Ils n'écrivent pas dans Supabase d'ici.",
			"Marketing, Panda, Pandora et le nom d'un client Proximity peuvent croître dans ce carnet.",
			"Lovable n'édite pas un thème WordPress.",
			"Giant et les bots ne déposent pas un ordre ici.",
			"Un courriel, un SMS ou une invitation au portail attend comms.allow, puis Epicenter.",
		},
		Writes: false,
		Sent:   false,
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
	if line == "wordpress" {
		return CRMGrowth{}, fmt.Errorf("lovable n'entre pas dans wordpress")
	}
	why, ok := crmLines[line]
	if !ok {
		return CRMGrowth{}, fmt.Errorf("%s n'écrit pas dans ce CRM", line)
	}
	if line == "proximity" && kind == "lead" {
		return CRMGrowth{}, fmt.Errorf("proximity ne dépose qu'un contact ou une note")
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
		return kernel.Result{OK: true, Message: "un seul CRM Lovable", Data: CRM()}, nil
	}
	note, err := CRMGrow(k, payloadQuery(call, "line"), payloadQuery(call, "kind"), title, payloadQuery(call, "body"), payloadQuery(call, "send") == "true")
	if err != nil {
		return kernel.Result{OK: false, Message: err.Error()}, nil
	}
	return kernel.Result{OK: true, Message: "brouillon CRM, rien n'est envoyé", Data: note}, nil
}
