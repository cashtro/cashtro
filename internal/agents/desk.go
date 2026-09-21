package agents

import (
	"errors"
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// Plan is the user-facing contract of what Castro can ask here.
type Plan struct {
	Headline string     `json:"headline"`
	Promise  string     `json:"promise"`
	Loop     []string   `json:"loop"`
	Can      []PlanItem `json:"can"`
	Cannot   []PlanItem `json:"cannot"`
	Create   []PlanItem `json:"create"`
	Avoid    []PlanItem `json:"avoid"`
}

// PlanItem is one line on the assistant plan.
type PlanItem struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Owner  string `json:"owner,omitempty"`
}

func deskPlan() Plan {
	return Plan{
		Headline: "I am Castro's desk assistant in Cashtro OS.",
		Promise:  "Ask in your own words — FR or EN, messy is fine. I save the original, rewrite it into a shippable task, then I do every allowed job in this kernel. I do not become a second product.",
		Loop: []string{
			"You ask. Anything. One line or a whole mandate.",
			"Desk captures the raw text. It is never overwritten.",
			"Desk betters the ask: title, owner, verb, acceptance, and a clear rewrite.",
			"I execute what is allowed. Blocked asks stay on the desk with the reason.",
			"Each 'better' pass tightens the same request. Done parks it. Lessons stay.",
		},
		Can: []PlanItem{
			{Title: "Build Cashtro OS", Detail: "Write Go kernel, agents, APIs, tests, and the desk UI. New agentics register here.", Owner: "desk"},
			{Title: "Ship the delivery line", Detail: "Park mandates in idea, move them to concept, then production. Planner turns a goal into ships.", Owner: "delivery"},
			{Title: "Research and remember", Detail: "Ingest sourced notes, search the desk, store episodic facts, recall them later.", Owner: "research"},
			{Title: "Open a pull request", Detail: "Branch, commit, test, and open a draft PR for the work. I do not merge unless you say so.", Owner: "deploy"},
			{Title: "Run and prove it", Detail: "go test, the local OS on :8080, screenshots and walkthroughs before we call a ship done.", Owner: "reviewer"},
			{Title: "Drive a browser", Detail: "Click through the desk the way a shipper would. Operator is resident until bound; I can still verify the UI from this environment.", Owner: "operator"},
			{Title: "Draft outbound", Detail: "Mail, Slack, email drafts park behind the human gate. Nothing leaves until you allow it.", Owner: "comms"},
			{Title: "Plan and investigate", Detail: "Architecture, CI, CVE/SAST triage, and incident traces. Resident verbs stay on the desk; I still do the work in this repo.", Owner: "architect"},
		},
		Cannot: []PlanItem{
			{Title: "Illegal or harmful work", Detail: "No exploits, malware, unauthorized access, weapons, or anything involving minors. I refuse and park it as blocked."},
			{Title: "Send without your say-so", Detail: "No live email, Slack, tweets, or money moves until you allow the confirm. Reads are free. Writes wait."},
			{Title: "Merge or force-push on my own", Detail: "I open draft PRs. I do not merge, enable auto-merge, or rewrite published history unless you ask."},
			{Title: "Invent a second product", Detail: "New agentics boot into this kernel. They do not fork another app, repo, or brand."},
			{Title: "Leak private client code", Detail: "Evolu-Jeunes / Proximity client work stays private. This public OS does not absorb their source."},
			{Title: "Pretend a resident verb is live", Detail: "Operator, reviewer, architect, deploy, security, and investigator sit on the desk until a worker binds. I still help by writing the code and running what this environment can run."},
			{Title: "Think without a model key", Detail: "OpenRouter is optional. The OS boots without it. model.chat stays unbound until you export a key."},
			{Title: "Remember across chats by magic", Detail: "Supermemory is not connected here. The desk log, notes, and this repo are the memory. Capture every ask so the next pass starts smarter."},
		},
		Create: []PlanItem{
			{Title: "Agentics", Detail: "Live or resident processes with a spec, capabilities, boot, and invoke. Register them in this kernel.", Owner: "desk"},
			{Title: "Ships", Detail: "Mandates on the idea → concept → production board, with client, sector, and stack.", Owner: "delivery"},
			{Title: "Notes and facts", Detail: "Sourced research records and episodic memory lines.", Owner: "research"},
			{Title: "Requests", Detail: "Every Castro ask, original + improved, bettered on each pass.", Owner: "desk"},
			{Title: "Confirms", Detail: "Human gates for outbound drafts.", Owner: "comms"},
			{Title: "APIs and UI", Detail: "stdlib Go HTTP shell and the desk in the browser. No extra framework.", Owner: "desk"},
			{Title: "Tests and PRs", Detail: "go test ./..., go vet, GitHub draft pull requests.", Owner: "deploy"},
			{Title: "Playbooks", Detail: "AGENTS.md and README so the next assistant runs the same loop.", Owner: "desk"},
		},
		Avoid: []PlanItem{
			{Title: "A second Cashtro", Detail: "No parallel dashboard, SaaS, or agent runtime outside this OS."},
			{Title: "Client repos in this tree", Detail: "Do not copy BTK, MD Clinic, or other private mandates into cashtro/cashtro."},
			{Title: "Secret keys in git", Detail: "OpenRouter and every token stay in the environment, never in the repo."},
			{Title: "Fake walkthroughs", Detail: "Do not upload toy screenshots. Prove the desk with a real capture, better, and done."},
			{Title: "Silent scope creep", Detail: "If an ask is blocked or needs a bind, say so on the request. Do not quietly skip it."},
		},
	}
}

func deskInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "desk.plan":
		p := deskPlan()
		return kernel.Result{OK: true, Message: p.Headline, Data: p}, nil
	case "desk.list":
		reqs := k.Requests()
		return kernel.Result{OK: true, Message: "requests " + strconv.Itoa(len(reqs)), Data: reqs}, nil
	case "desk.capture":
		raw := payloadQuery(call, "raw")
		if raw == "" {
			raw = payloadQuery(call, "prompt")
		}
		if raw == "" {
			raw = payloadQuery(call, "body")
		}
		got, err := captureAsk(k, raw)
		if err != nil {
			if errors.Is(err, kernel.ErrInvalidRequest) {
				return kernel.Result{OK: false, Message: "raw required"}, nil
			}
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "captured #" + strconv.Itoa(got.ID) + " · pass " + strconv.Itoa(got.Pass), Data: got}, nil
	case "desk.better":
		id := payloadInt(call, "id")
		got, err := betterAsk(k, id)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "bettered #" + strconv.Itoa(got.ID) + " · pass " + strconv.Itoa(got.Pass), Data: got}, nil
	case "desk.done":
		id := payloadInt(call, "id")
		got, err := finishAsk(k, id, kernel.StatusDone)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "done #" + strconv.Itoa(got.ID), Data: got}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func seedDesk(k *kernel.Kernel) {
	k.Remember("assistant", "Castro wants one desk assistant who does every allowed task. Capture every request. Better it each pass. Stay in this kernel.")
	seeds := []string{
		"Standing rule: do all my tasks. Save every request. Better it each time. Speak FR/EN. Do not fork a second product.",
		"Ok now give me the list of thing i cna and cannot do with you here and what can we create ot not etc... fukk user friendly plan make sure all request get saved and bettered also eash time... i need an assistant who will do all my tasks",
	}
	for _, raw := range seeds {
		_, _ = captureAsk(k, raw)
	}
}

func captureAsk(k *kernel.Kernel, raw string) (kernel.Request, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return kernel.Request{}, kernel.ErrInvalidRequest
	}
	improved := betterOnce(kernel.Request{Raw: raw})
	got, err := k.SaveRequest(improved)
	if err != nil {
		return kernel.Request{}, err
	}
	k.Remember("request", got.Title)
	if got.Owner != "" && got.Owner != "desk" {
		_, _ = k.Post("desk", got.Owner, "request", got.Improved)
	}
	return got, nil
}

func betterAsk(k *kernel.Kernel, id int) (kernel.Request, error) {
	cur, err := latestOrID(k, id)
	if err != nil {
		return kernel.Request{}, err
	}
	next := betterOnce(cur)
	got, err := k.SaveRequest(next)
	if err != nil {
		return kernel.Request{}, err
	}
	k.Remember("request", "pass "+strconv.Itoa(got.Pass)+" · "+got.Title)
	return got, nil
}

func finishAsk(k *kernel.Kernel, id int, status string) (kernel.Request, error) {
	cur, err := latestOrID(k, id)
	if err != nil {
		return kernel.Request{}, err
	}
	return k.AdvanceRequest(cur.ID, status)
}

func latestOrID(k *kernel.Kernel, id int) (kernel.Request, error) {
	if id != 0 {
		return k.RequestByID(id)
	}
	all := k.Requests()
	if len(all) == 0 {
		return kernel.Request{}, kernel.ErrUnknownRequest
	}
	return all[0], nil
}

func betterOnce(prev kernel.Request) kernel.Request {
	r := prev
	r.Pass++
	raw := strings.TrimSpace(r.Raw)
	if blocked, why := refuseAsk(raw); blocked {
		r.Allowed = false
		r.Status = kernel.StatusBlocked
		r.Owner = "desk"
		r.Verb = "desk.capture"
		r.Reason = why
		r.Title = "Blocked ask"
		r.Improved = "This request cannot be executed. " + why + " Original kept on the desk."
		r.Criteria = []string{"raw saved", "classified cannot-do", "no execution"}
		r.Lessons = append(r.Lessons, "pass "+strconv.Itoa(r.Pass)+": refused and parked as blocked")
		return r
	}

	owner, verb, reason := routeAsk(raw)
	r.Allowed = true
	r.Owner = owner
	r.Verb = verb
	r.Reason = reason
	r.Title = titleFrom(raw, verb)
	r.Criteria = criteriaFor(raw, r.Pass)
	r.Improved = rewriteAsk(raw, r)
	switch {
	case r.Pass == 1:
		r.Status = kernel.StatusCaptured
	case r.Status == kernel.StatusCaptured, r.Status == "":
		r.Status = kernel.StatusClarified
	}
	r.Lessons = append(r.Lessons, lessonFor(r.Pass, owner, verb))
	return r
}

func refuseAsk(raw string) (bool, string) {
	low := strings.ToLower(raw)
	needles := []string{
		"exploit", "malware", "ransomware", "keylogger", "unauthorized access",
		"child sexual", "csam", "weaponize", "build a bomb", "credit card dump",
	}
	for _, n := range needles {
		if strings.Contains(low, n) {
			return true, "This is on the cannot-do list. I will not help with harmful or illegal work."
		}
	}
	return false, ""
}

func routeAsk(raw string) (owner, verb, reason string) {
	low := strings.ToLower(raw)
	switch {
	case hasAny(low, "send", "email", "slack", "tweet", "outbound"):
		return "comms", "comms.send", "Outbound waits on the human gate."
	case hasAny(low, "backlog", "mandate", "park it", "park a", "delivery line"):
		return "planner", "planner.backlog", "Goals become idea-stage ships."
	case hasAny(low, "advance", "move to concept", "production"):
		return "delivery", "delivery.advance", "The live line is idea → concept → production."
	case hasAny(low, "remember", "recall", "memory"):
		return "memory", "memory.store", "Facts live on the memory agentic."
	case hasAny(low, "research", "paper", "arxiv", "note", "find out"):
		return "research", "research.ingest", "Findings belong in the research library, not only in chat."
	case hasAny(low, "browse", "screenshot", "click through"):
		return "operator", "operator.browse", "Operator is resident; the assistant can still verify the UI."
	case hasAny(low, "deploy", "release", "ci"):
		return "deploy", "deploy.release", "Deploy is resident; PRs and tests still run from this environment."
	case hasAny(low, "cve", "sast", "security"):
		return "security", "security.triage", "Security stays on the desk until a worker binds."
	case hasAny(low, "incident", "blast radius", "failing check"):
		return "investigator", "investigator.trace", "Investigator is resident; traces still start from this repo."
	default:
		return "desk", "desk.capture", "Default owner is the desk. I execute allowed work in this kernel."
	}
}

func titleFrom(raw, verb string) string {
	low := strings.ToLower(raw)
	switch {
	case hasAny(low, "can and cannot", "cna and cannot", "what can we create", "user friendly plan"):
		return "Publish can/cannot plan and persist every request"
	case hasAny(low, "do all my tasks", "standing rule"):
		return "Standing rule: capture, better, and do every allowed task"
	case hasAny(low, "openrouter"):
		return "Keep OpenRouter optional"
	}
	cleaned := collapseSpace(raw)
	cleaned = strings.TrimRight(cleaned, ".!?")
	runes := []rune(cleaned)
	if len(runes) == 0 {
		return verb
	}
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

func rewriteAsk(raw string, r kernel.Request) string {
	clean := collapseSpace(raw)
	var b strings.Builder
	b.WriteString("Castro asked: ")
	b.WriteString(clean)
	b.WriteString("\n\nDo this: ")
	b.WriteString(r.Title)
	b.WriteString(". Route through ")
	b.WriteString(r.Owner)
	b.WriteString(" / ")
	b.WriteString(r.Verb)
	b.WriteString(". ")
	b.WriteString(r.Reason)
	b.WriteString(" Stay in Cashtro OS. Do not fork a second product.")
	if r.Pass >= 2 && len(r.Criteria) > 0 {
		b.WriteString("\n\nAcceptance:\n")
		for _, c := range r.Criteria {
			b.WriteString("- ")
			b.WriteString(c)
			b.WriteString("\n")
		}
	}
	if r.Pass >= 3 {
		b.WriteString("\nNext: execute ")
		b.WriteString(r.Verb)
		b.WriteString(", prove it, then desk.done.")
	}
	return strings.TrimSpace(b.String())
}

func criteriaFor(raw string, pass int) []string {
	low := strings.ToLower(raw)
	out := []string{"raw text saved on the desk", "ask classified can/cannot", "owner and verb assigned"}
	if hasAny(low, "plan", "can and cannot", "create") {
		out = append(out, "user-friendly can/cannot plan on the OS desk")
	}
	if hasAny(low, "saved", "better", "request") {
		out = append(out, "every new ask is captured and bettered on each pass")
	}
	if hasAny(low, "assistant", "all my tasks") {
		out = append(out, "assistant executes allowed work instead of only listing it")
	}
	if pass >= 2 {
		out = append(out, "acceptance criteria visible on the request card")
	}
	if pass >= 3 {
		out = append(out, "next verb named so the desk can invoke it")
	}
	if pass >= 4 {
		out = append(out, "constraints and cannot-do lines restated so scope cannot drift")
	}
	return out
}

func lessonFor(pass int, owner, verb string) string {
	switch pass {
	case 1:
		return "pass 1: saved original, classified, assigned " + owner + " / " + verb
	case 2:
		return "pass 2: added acceptance criteria so 'done' is testable"
	case 3:
		return "pass 3: named the next verb to execute"
	default:
		return "pass " + strconv.Itoa(pass) + ": tightened constraints from the cannot-do list"
	}
}

func hasAny(low string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(low, n) {
			return true
		}
	}
	return false
}

func collapseSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
