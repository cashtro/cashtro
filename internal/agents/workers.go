package agents

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// Flow binds one division to the function that runs it and the product it may touch.
type Flow struct {
	Division string `json:"division"`
	Function string `json:"function"`
	Product  string `json:"product"`
	Agent    string `json:"agent"`
}

// Flows is the flow stack: division, function, product.
func Flows() []Flow {
	org := Chart()
	out := make([]Flow, 0, len(org.Divisions))
	for _, d := range org.Divisions {
		out = append(out, Flow{
			Division: d.ID,
			Function: d.Function,
			Product:  d.Product,
			Agent:    d.Chief,
		})
	}
	return out
}

func flowFor(line string) (Flow, bool) {
	for _, f := range Flows() {
		if f.Division == line {
			return f, true
		}
	}
	return Flow{}, false
}

var secretShapes = []struct {
	name string
	re   *regexp.Regexp
}{
	{name: "jeton stripe", re: regexp.MustCompile(`(?i)sk_(live|test)_[A-Za-z0-9]+`)},
	{name: "clé privée", re: regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	{name: "clé aws", re: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
}

func secretNames(text string) []string {
	var names []string
	for _, shape := range secretShapes {
		if shape.re.MatchString(text) {
			names = append(names, shape.name)
		}
	}
	return names
}

func initInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	if call.Capability == "os.cycle" {
		res, err := RunCycle(k, call)
		if err != nil {
			return res, err
		}
		k.Publish("init", "voltron", "cycle lancé avec epicenter", map[string]any{"line": payloadQuery(call, "id")})
		return res, nil
	}
	return kernel.Result{OK: true, Message: "about", Data: k.About()}, nil
}

func operatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	repo := payloadQuery(call, "repo")
	task := payloadQuery(call, "task")
	if task == "" {
		task = payloadQuery(call, "prompt")
	}
	if task == "" {
		task = "travail local"
	}
	if repo != "" {
		ln, ok := rosterOwner(repo)
		if !ok {
			return kernel.Result{OK: false, Message: "dépôt hors roster"}, nil
		}
		if ln.ID == "proximity" {
			return kernel.Result{OK: false, Message: "autre département: " + repo}, nil
		}
	}
	k.Remember("operator", strings.TrimSpace(repo+" "+task))
	return kernel.Result{OK: true, Message: "travail local enregistré", Data: map[string]any{
		"repo":     repo,
		"task":     task,
		"local":    true,
		"pushed":   false,
		"xampp":    repo == "" || strings.Contains(repo, "Evolu-Jeunes/") || strings.Contains(repo, "cashtro/"),
		"function": "operator.work",
	}}, nil
}

func reviewerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	var proposal Proposal
	if len(call.Payload) > 0 {
		_ = json.Unmarshal(call.Payload, &proposal)
	}
	if len(proposal.Options) < 2 {
		proposal.Options = []Option{
			{Name: "refaire", Cost: 8, Risk: 3, Steps: 6},
			{Name: "éditer la page demandée", Cost: 1, Risk: 1, Steps: 2},
		}
	}
	if proposal.Subject == "" {
		proposal.Subject = payloadQuery(call, "task")
	}
	verdict := Contradict(proposal)
	k.Remember("reviewer", verdict.Attack)
	return kernel.Result{OK: verdict.Ready, Message: verdict.Attack, Data: verdict}, nil
}

func architectInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	line := payloadQuery(call, "id")
	if line == "" {
		line = payloadQuery(call, "line")
	}
	if line == "" {
		line = "wordpress"
	}
	flow, ok := flowFor(line)
	if !ok {
		return kernel.Result{OK: false, Message: "division inconnue: " + line}, nil
	}
	verdict := Contradict(Proposal{
		Subject: line,
		Options: []Option{
			{Name: "nouveau produit", Cost: 5, Risk: 3, Steps: 4},
			{Name: "produit déjà nommé", Cost: 1, Risk: 1, Steps: 1},
		},
	})
	k.Remember("architect", flow.Division+" "+flow.Function+" "+flow.Product)
	return kernel.Result{OK: verdict.Ready, Message: flow.Function + " → " + flow.Product, Data: map[string]any{
		"flow":    flow,
		"verdict": verdict,
	}}, nil
}

func deployInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	repo := payloadQuery(call, "repo")
	task := payloadQuery(call, "task")
	id := payloadInt(call, "id")
	if id == 0 {
		c := k.RequestConfirm("deploy", "deploy.release", strings.TrimSpace(repo+" "+task))
		return kernel.Result{OK: false, Message: "rien n'est poussé. confirm en attente", Data: c}, nil
	}
	for _, c := range k.Confirms() {
		if c.ID == id && c.Status == "allowed" {
			k.Remember("deploy", "local "+repo)
			return kernel.Result{OK: true, Message: "sortie locale enregistrée. rien n'est poussé en production", Data: c}, nil
		}
	}
	return kernel.Result{OK: false, Message: "confirm refusé ou absent"}, nil
}

func securityInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	text := payloadQuery(call, "task") + " " + payloadQuery(call, "text") + " " + payloadQuery(call, "repo")
	if strings.TrimSpace(text) == "" && len(call.Payload) > 0 {
		text = string(call.Payload)
	}
	found := secretNames(text)
	if len(found) > 0 {
		k.Remember("security", "secret détecté: "+strings.Join(found, ", "))
		return kernel.Result{OK: false, Message: "interdit: " + strings.Join(found, ", ") + ". valeur non recopiée", Data: map[string]any{
			"findings": found,
			"copied":   false,
		}}, nil
	}
	k.Remember("security", "filtre passé")
	return kernel.Result{OK: true, Message: "filtre passé", Data: map[string]any{"findings": []string{}}}, nil
}

func investigatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	repo := payloadQuery(call, "repo")
	lineID := ""
	var siblings []string
	if repo != "" {
		if ln, ok := rosterOwner(repo); ok {
			lineID = ln.ID
			for _, name := range ln.Repos {
				if name != repo {
					siblings = append(siblings, name)
				}
				if len(siblings) == 8 {
					break
				}
			}
		}
	}
	if lineID == "" {
		lineID = payloadQuery(call, "id")
	}
	dept := ""
	for _, d := range Chart().Divisions {
		if d.ID == lineID {
			dept = d.Department
		}
	}
	k.Remember("investigator", lineID+" "+repo)
	return kernel.Result{OK: true, Message: "rayon nommé", Data: map[string]any{
		"repo":       repo,
		"line":       lineID,
		"department": dept,
		"siblings":   siblings,
	}}, nil
}

// RunCycle is the autonomous pass. Voltron and Epicenter both call it.
// Each worker runs. A secret stops the pass. A release stays behind comms.allow.
func RunCycle(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	line := payloadQuery(call, "id")
	if line == "" {
		line = payloadQuery(call, "line")
	}
	if line == "" {
		line = "wordpress"
	}
	steps, ok := Chain(line)
	if !ok {
		return kernel.Result{OK: false, Message: "chaîne inconnue: " + line}, nil
	}
	if open := AskSelf(line, nil).Open; len(open) > 0 {
		return kernel.Result{OK: false, Message: askMessage(open), Data: open}, nil
	}
	repo := payloadQuery(call, "repo")
	task := payloadQuery(call, "task")
	if task == "" {
		task = "travail local"
	}
	done := make([]map[string]any, 0, len(steps))
	for _, step := range steps {
		res := runStep(k, step.Agent, repo, task, line)
		done = append(done, map[string]any{
			"agent":   step.Agent,
			"order":   step.Order,
			"ok":      res.OK,
			"message": res.Message,
		})
		if !res.OK {
			return kernel.Result{OK: false, Message: res.Message, Data: map[string]any{
				"line": line, "directedBy": []string{"voltron", "epicenter"}, "steps": done,
			}}, nil
		}
	}
	k.Publish("manager", "epicenter", "cycle bouclé "+line, map[string]any{"steps": len(done)})
	return kernel.Result{OK: true, Message: "cycle autonome bouclé. rien n'est poussé", Data: map[string]any{
		"line": line, "directedBy": []string{"voltron", "epicenter"}, "steps": done,
	}}, nil
}

func runStep(k *kernel.Kernel, agent, repo, task, line string) kernel.Result {
	raw, _ := json.Marshal(map[string]string{"repo": repo, "task": task, "id": line, "line": line})
	call := kernel.Call{Payload: raw}
	switch agent {
	case "operator":
		call.Capability = "operator.work"
		res, _ := operatorInvoke(k, call)
		return res
	case "security":
		call.Capability = "security.triage"
		res, _ := securityInvoke(k, call)
		return res
	case "reviewer":
		call.Capability = "reviewer.watch"
		res, _ := reviewerInvoke(k, call)
		return res
	case "architect":
		call.Capability = "architect.plan"
		res, _ := architectInvoke(k, call)
		return res
	case "deploy":
		call.Capability = "deploy.release"
		res, _ := deployInvoke(k, call)
		return res
	case "investigator":
		call.Capability = "investigator.trace"
		res, _ := investigatorInvoke(k, call)
		return res
	case "comms":
		c := k.RequestConfirm("comms", "comms.allow", strings.TrimSpace(repo+" "+task))
		return kernel.Result{OK: true, Message: "en attente de comms.allow", Data: c}
	default:
		k.Remember(agent, line+" "+task)
		return kernel.Result{OK: true, Message: agent + " a fait son pas"}
	}
}
