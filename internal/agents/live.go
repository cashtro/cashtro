package agents

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
)

func explorerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	q := strings.ToLower(payloadQuery(call, "query"))
	hits := make([]map[string]string, 0)
	for _, p := range k.Processes() {
		blob := strings.ToLower(p.Spec.Name + " " + p.Spec.Role + " " + p.Spec.Summary + " " + strings.Join(p.Spec.Capabilities, " "))
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "agent", "id": p.Spec.ID, "name": p.Spec.Name})
		}
	}
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			blob := strings.ToLower(s.Name + " " + s.Client + " " + s.Notes + " " + string(s.Stage))
			if q == "" || strings.Contains(blob, q) {
				hits = append(hits, map[string]string{"kind": "ship", "id": s.ID, "name": s.Name})
			}
		}
	}
	for _, n := range k.Notes() {
		blob := strings.ToLower(n.Claim + " " + n.Source + " " + n.Quote)
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "note", "id": n.Source, "name": n.Claim})
		}
	}
	_, _ = k.Post("explorer", "research", "search", q)
	return kernel.Result{OK: true, Message: "explorer hit " + strconv.Itoa(len(hits)), Data: hits}, nil
}

func investigatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	q := strings.ToLower(payloadQuery(call, "query"))
	if q == "" {
		q = strings.ToLower(payloadQuery(call, "prompt"))
	}
	blast := map[string]any{
		"query":    q,
		"ships":    []map[string]string{},
		"notes":    []map[string]string{},
		"events":   []map[string]string{},
		"mail":     []map[string]string{},
		"confirms": []map[string]string{},
		"facts":    []map[string]string{},
	}
	ships := blast["ships"].([]map[string]string)
	notes := blast["notes"].([]map[string]string)
	events := blast["events"].([]map[string]string)
	mail := blast["mail"].([]map[string]string)
	confirms := blast["confirms"].([]map[string]string)
	facts := blast["facts"].([]map[string]string)

	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			blob := strings.ToLower(s.Name + " " + s.Client + " " + s.Notes + " " + s.Sector + " " + string(s.Stage))
			if q == "" || strings.Contains(blob, q) {
				ships = append(ships, map[string]string{"id": s.ID, "stage": string(s.Stage), "name": s.Name})
			}
		}
	}
	for _, n := range k.Notes() {
		blob := strings.ToLower(n.Claim + " " + n.Source + " " + n.Quote)
		if q == "" || strings.Contains(blob, q) {
			notes = append(notes, map[string]string{"source": n.Source, "claim": n.Claim, "url": n.URL})
		}
	}
	for _, ev := range k.Events() {
		blob := strings.ToLower(ev.Source + " " + ev.Kind + " " + ev.Message)
		if q == "" || strings.Contains(blob, q) {
			events = append(events, map[string]string{"kind": ev.Kind, "source": ev.Source, "message": ev.Message})
		}
	}
	for _, m := range k.Inbox("") {
		blob := strings.ToLower(m.From + " " + m.To + " " + m.Kind + " " + m.Body)
		if q == "" || strings.Contains(blob, q) {
			mail = append(mail, map[string]string{"from": m.From, "to": m.To, "body": m.Body})
		}
	}
	for _, c := range k.Confirms() {
		blob := strings.ToLower(c.Agent + " " + c.Cap + " " + c.Body + " " + c.Status)
		if q == "" || strings.Contains(blob, q) {
			confirms = append(confirms, map[string]string{"id": strconv.Itoa(c.ID), "status": c.Status, "body": c.Body})
		}
	}
	for _, f := range k.Recall(q) {
		facts = append(facts, map[string]string{"topic": f.Topic, "text": f.Text})
	}

	blast["ships"] = ships
	blast["notes"] = notes
	blast["events"] = events
	blast["mail"] = mail
	blast["confirms"] = confirms
	blast["facts"] = facts
	total := len(ships) + len(notes) + len(events) + len(mail) + len(confirms) + len(facts)
	_, _ = k.Post("investigator", "memory", "trace", q)
	k.Remember("incident", "trace "+q+" · "+strconv.Itoa(total)+" hits")
	return kernel.Result{OK: true, Message: "blast radius " + strconv.Itoa(total), Data: blast}, nil
}

func memoryInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "memory.store":
		topic := payloadQuery(call, "topic")
		text := payloadQuery(call, "text")
		if text == "" {
			text = payloadQuery(call, "prompt")
		}
		if topic == "" {
			topic = "working"
		}
		f := k.Remember(topic, text)
		return kernel.Result{OK: true, Message: "remembered " + f.Topic, Data: f}, nil
	case "memory.recall":
		q := payloadQuery(call, "query")
		if q == "" {
			q = payloadQuery(call, "topic")
		}
		got := k.Recall(q)
		return kernel.Result{OK: true, Message: "recalled " + strconv.Itoa(len(got)), Data: got}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func commsInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "comms.send":
		body := payloadQuery(call, "body")
		if body == "" {
			body = payloadQuery(call, "prompt")
		}
		if body == "" {
			body = "outbound draft"
		}
		c := k.RequestConfirm("comms", "comms.send", body)
		return kernel.Result{OK: true, Message: "parked confirm #" + strconv.Itoa(c.ID), Data: c}, nil
	case "comms.pending":
		return kernel.Result{OK: true, Message: "confirms", Data: k.Confirms()}, nil
	case "comms.allow":
		id := payloadInt(call, "id")
		c, err := k.DecideConfirm(id, true)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: c.Outcome, Data: c}, nil
	case "comms.deny":
		id := payloadInt(call, "id")
		c, err := k.DecideConfirm(id, false)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: c.Outcome, Data: c}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func operatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	url := payloadQuery(call, "url")
	if url == "" {
		url = payloadQuery(call, "prompt")
	}
	if url == "" {
		url = "http://127.0.0.1:8080/"
	}
	steps := []string{
		"Open " + url,
		"Read /health and /api/os",
		"Walk process table and delivery board",
		"Capture evidence for reviewer.watch",
	}
	note := k.WriteNote(kernel.Note{
		Agent:  "operator",
		Source: "operator.browse",
		Claim:  "Browse plan · " + url,
		Quote:  strings.Join(steps, " → "),
	})
	k.Remember("operator", "browse plan · "+url)
	_, _ = k.Post("operator", "reviewer", "browse", url)
	return kernel.Result{
		OK:      true,
		Message: "browse planned · " + url,
		Data: map[string]any{
			"url":   url,
			"steps": steps,
			"note":  note,
			"mode":  "dry-run",
			"bind":  "computer-use worker not attached · plan only",
		},
	}, nil
}

func reviewerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	target := payloadQuery(call, "target")
	if target == "" {
		target = payloadQuery(call, "prompt")
	}
	if target == "" {
		target = "desk"
	}
	var idea, concept, prod int
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			switch s.Stage {
			case catalog.StageIdea:
				idea++
			case catalog.StageConcept:
				concept++
			case catalog.StageProduction:
				prod++
			}
		}
	}
	notes := len(k.Notes())
	events := len(k.Events())
	pending := 0
	for _, c := range k.Confirms() {
		if c.Status == "pending" {
			pending++
		}
	}
	checks := []map[string]any{
		{"name": "production evidence", "ok": prod > 0, "detail": strconv.Itoa(prod) + " ships"},
		{"name": "research library", "ok": notes >= 3, "detail": strconv.Itoa(notes) + " notes"},
		{"name": "journal trail", "ok": events > 0, "detail": strconv.Itoa(events) + " events"},
		{"name": "no pending outbound", "ok": pending == 0, "detail": strconv.Itoa(pending)},
		{"name": "wip visible", "ok": idea+concept >= 0, "detail": "idea=" + strconv.Itoa(idea) + " concept=" + strconv.Itoa(concept)},
	}
	verdict := "pass"
	for _, c := range checks {
		if ok, _ := c["ok"].(bool); !ok {
			verdict = "fail"
			break
		}
	}
	note := k.WriteNote(kernel.Note{
		Agent:  "reviewer",
		Source: "reviewer.watch",
		Claim:  "QA " + verdict + " · " + target,
		Quote:  "production=" + strconv.Itoa(prod) + " notes=" + strconv.Itoa(notes) + " events=" + strconv.Itoa(events) + " pending=" + strconv.Itoa(pending),
	})
	k.Remember("qa", verdict+" · "+target)
	_, _ = k.Post("reviewer", "deploy", "qa", verdict+" · "+target)
	return kernel.Result{
		OK:      true,
		Message: "review " + verdict + " · " + target,
		Data: map[string]any{
			"target":  target,
			"verdict": verdict,
			"checks":  checks,
			"note":    note,
			"mode":    "desk-evidence",
		},
	}, nil
}

func deployInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	target := payloadQuery(call, "target")
	if target == "" {
		target = payloadQuery(call, "prompt")
	}
	if target == "" {
		target = "cashtro-os"
	}
	var idea, concept, prod int
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			switch s.Stage {
			case catalog.StageIdea:
				idea++
			case catalog.StageConcept:
				concept++
			case catalog.StageProduction:
				prod++
			}
		}
	}
	pending := 0
	for _, c := range k.Confirms() {
		if c.Status == "pending" {
			pending++
		}
	}
	checks := []map[string]any{
		{"name": "disk image", "ok": k.PersistPath() != "", "detail": k.PersistPath()},
		{"name": "production ships", "ok": prod > 0, "detail": strconv.Itoa(prod)},
		{"name": "open idea work", "ok": true, "detail": strconv.Itoa(idea)},
		{"name": "concept WIP", "ok": true, "detail": strconv.Itoa(concept)},
		{"name": "pending confirms", "ok": pending == 0, "detail": strconv.Itoa(pending)},
		{"name": "research notes", "ok": len(k.Notes()) > 0, "detail": strconv.Itoa(len(k.Notes()))},
	}
	ready := true
	for _, c := range checks {
		if ok, _ := c["ok"].(bool); !ok {
			ready = false
			break
		}
	}
	status := "blocked"
	if ready {
		status = "ready"
	}
	note := k.WriteNote(kernel.Note{
		Agent:  "deploy",
		Source: "deploy.release",
		Claim:  "Release dry-run " + status + " · " + target,
		Quote:  "idea=" + strconv.Itoa(idea) + " concept=" + strconv.Itoa(concept) + " production=" + strconv.Itoa(prod) + " pending=" + strconv.Itoa(pending),
	})
	k.Remember("release", status+" · "+target)
	_, _ = k.Post("deploy", "security", "release", target)
	confirm := kernel.Confirm{}
	if ready {
		confirm = k.RequestConfirm("deploy", "deploy.release", "promote "+target+" (dry-run only · no outbound)")
	}
	return kernel.Result{
		OK:      true,
		Message: "deploy " + status + " · " + target,
		Data: map[string]any{
			"target":  target,
			"status":  status,
			"checks":  checks,
			"note":    note,
			"confirm": confirm,
			"mode":    "dry-run",
		},
	}, nil
}

func securityInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	q := strings.ToLower(payloadQuery(call, "query"))
	if q == "" {
		q = strings.ToLower(payloadQuery(call, "prompt"))
	}
	keywords := []string{"cve", "secret", "password", "token", "exploit", "breach", "rce", "sqli", "xss", "critical", "fail", "error"}
	if q != "" {
		keywords = append(keywords, q)
	}
	findings := make([]map[string]string, 0)
	match := func(blob, kind, id, name string) {
		low := strings.ToLower(blob)
		for _, kw := range keywords {
			if kw != "" && strings.Contains(low, kw) {
				findings = append(findings, map[string]string{
					"kind": kind, "id": id, "name": name, "hit": kw,
				})
				return
			}
		}
	}
	if cat := k.Catalog(); cat != nil {
		for _, s := range cat.List() {
			match(s.Name+" "+s.Notes+" "+s.Sector+" "+string(s.Stage), "ship", s.ID, s.Name)
		}
	}
	for _, n := range k.Notes() {
		match(n.Claim+" "+n.Source+" "+n.Quote, "note", n.Source, n.Claim)
	}
	for _, ev := range k.Events() {
		match(ev.Source+" "+ev.Kind+" "+ev.Message, "event", ev.Kind, ev.Message)
	}
	severity := "clear"
	if len(findings) > 0 {
		severity = "review"
		k.RequestConfirm("security", "security.triage", "triage "+strconv.Itoa(len(findings))+" hits · query="+q)
	}
	k.Remember("security", "triage "+severity+" · "+strconv.Itoa(len(findings))+" hits")
	_, _ = k.Post("security", "investigator", "triage", q)
	return kernel.Result{
		OK:      true,
		Message: "security " + severity + " · " + strconv.Itoa(len(findings)) + " hits",
		Data: map[string]any{
			"severity": severity,
			"query":    q,
			"findings": findings,
		},
	}, nil
}

func architectInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	goal := payloadQuery(call, "goal")
	if goal == "" {
		goal = payloadQuery(call, "prompt")
	}
	if goal == "" {
		goal = "unnamed system"
	}
	steps := payloadList(call, "steps")
	if len(steps) == 0 {
		steps = []string{
			"Map the control plane and process table",
			"Park the mandate on the delivery line",
			"Prove it with tests and a desk walkthrough",
			"Ship only after human confirm on outbound",
		}
	}
	var b strings.Builder
	b.WriteString("Plan for: ")
	b.WriteString(goal)
	b.WriteString("\n")
	for i, step := range steps {
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(". ")
		b.WriteString(step)
		b.WriteString("\n")
	}
	note := k.WriteNote(kernel.Note{
		Agent:  "architect",
		Source: "architect.plan",
		Claim:  "Plan: " + goal,
		Quote:  strings.TrimSpace(b.String()),
	})
	k.Remember("plan", goal)
	_, _ = k.Post("architect", "planner", "plan", goal)
	_, _ = k.Post("architect", "delivery", "plan", goal)

	created := make([]catalog.Ship, 0)
	if cat := k.Catalog(); cat != nil {
		for _, step := range steps {
			ship, err := cat.Create(catalog.CreateShip{
				Name:   step,
				Client: "Cashtro",
				Sector: "architecture",
				Stack:  []string{"Go"},
				Notes:  "from architect.plan · " + goal,
			})
			if err != nil {
				return kernel.Result{}, err
			}
			created = append(created, ship)
		}
	}
	return kernel.Result{
		OK:      true,
		Message: "planned " + strconv.Itoa(len(steps)) + " steps for " + goal,
		Data: map[string]any{
			"goal":  goal,
			"steps": steps,
			"note":  note,
			"ships": created,
		},
	}, nil
}

func plannerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	goal := payloadQuery(call, "goal")
	if goal == "" {
		goal = payloadQuery(call, "prompt")
	}
	if goal == "" {
		goal = "unnamed mandate"
	}
	cat := k.Catalog()
	if cat == nil {
		return kernel.Result{OK: false, Message: "delivery board not attached"}, nil
	}
	items := payloadList(call, "items")
	if len(items) == 0 {
		items = []string{goal}
	}
	created := make([]catalog.Ship, 0, len(items))
	for _, item := range items {
		ship, err := cat.Create(catalog.CreateShip{Name: item, Notes: goal, Sector: "research", Stack: []string{"Go"}})
		if err != nil {
			return kernel.Result{}, err
		}
		created = append(created, ship)
	}
	_, _ = k.Post("planner", "delivery", "backlog", goal)
	k.Remember("plan", goal)
	return kernel.Result{OK: true, Message: "parked " + strconv.Itoa(len(created)) + " on idea", Data: created}, nil
}

func researchInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "research.list":
		return kernel.Result{OK: true, Message: "notes " + strconv.Itoa(len(k.Notes())), Data: k.Notes()}, nil
	case "research.ingest", "note.write":
		n := kernel.Note{
			Agent:  "research",
			Source: payloadQuery(call, "source"),
			URL:    payloadQuery(call, "url"),
			Claim:  payloadQuery(call, "claim"),
			Quote:  payloadQuery(call, "quote"),
		}
		if n.Claim == "" {
			n.Claim = payloadQuery(call, "prompt")
		}
		if n.Source == "" {
			n.Source = "desk"
		}
		if n.Claim == "" {
			return kernel.Result{OK: false, Message: "claim required"}, nil
		}
		got := k.WriteNote(n)
		k.Remember("research", got.Claim)
		return kernel.Result{OK: true, Message: "ingested " + got.Source, Data: got}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func seedResearch(k *kernel.Kernel) {
	seeds := []kernel.Note{
		{Agent: "research", Source: "arxiv:2606.01508", URL: "https://arxiv.org/abs/2606.01508", Claim: "An AOS is a control plane over a classical OS, not a Linux replacement.", Quote: "AOS mediates tool invocations the way a classical OS mediates syscalls. First-class entities include agent identity, goals, capabilities, context, and execution records."},
		{Agent: "research", Source: "arxiv:2608.03214", URL: "https://arxiv.org/abs/2608.03214", Claim: "Split governance from runtime: intent/policy/audit vs lifecycle/routing/memory.", Quote: "Control & Governance plane owns authority and human oversight. Runtime & Coordination owns agent lifecycle and tool routing."},
		{Agent: "research", Source: "XKernel", URL: "https://github.com/JosephBerm/XKernel", Claim: "Treat agents as first-class processes with capability tokens and typed IPC.", Quote: "Unix treats processes; Kubernetes treats containers; an agent OS treats agents."},
		{Agent: "research", Source: "12-factor agents", URL: "https://github.com/humanlayer/12-factor-agents", Claim: "Own the loop in deterministic code. The model only fills structured next steps.", Quote: "Human confirm sits between selection and invocation. OpenRouter stays optional."},
		{Agent: "research", Source: "treg", URL: "https://treg.to", Claim: "Exa publication search is available as research.ingest input at $0.007/call when Treg is signed in.", Quote: "catalog_search → catalog_get → call. Token was expired this pass; notes still landed from open sources."},
		{Agent: "research", Source: "always-on", URL: "https://github.com/cashtro/cashtro", Claim: "Always-on means the cloud VM keeps the kernel looping — closing a laptop does not stop Cashtro OS.", Quote: "scripts/always-on.sh rebuilds and restarts :8080. Autosave flushes data/cashtro.json so hard kills still leave a durable image."},
		{Agent: "research", Source: "planner-api", URL: "https://github.com/cashtro/cashtro", Claim: "Planner is desk-reachable: POST /api/plan parks goal items as idea-stage ships on the delivery line.", Quote: "Ultron and overnight keep-alives can backlog company work without OpenRouter."},
		{Agent: "research", Source: "architect-api", URL: "https://github.com/cashtro/cashtro", Claim: "Architect is desk-reachable: POST /api/architect turns a goal into a plan note, memory, and idea-stage ships.", Quote: "Deterministic design loop — no model key required."},
		{Agent: "research", Source: "memory-api", URL: "https://github.com/cashtro/cashtro", Claim: "Memory is desk-writable: POST /api/memory stores episodic facts; GET /api/memory recalls them.", Quote: "Working memory survives always-on restarts via data/cashtro.json."},
		{Agent: "research", Source: "search-api", URL: "https://github.com/cashtro/cashtro", Claim: "Explorer is desk-reachable: GET /api/search?q= hits agents, ships, and research notes.", Quote: "Ultron and keep-alives can find work without OpenRouter."},
		{Agent: "research", Source: "comms-api", URL: "https://github.com/cashtro/cashtro", Claim: "Comms is desk-reachable: POST /api/comms parks an outbound draft behind a human confirm.", Quote: "No send until allow — Ultron clients inherit the gate."},
	}
	for _, n := range seeds {
		k.WriteNote(n)
		k.Remember("research", n.Claim)
	}
}

func payloadQuery(call kernel.Call, key string) string {
	if len(call.Payload) == 0 {
		return ""
	}
	var obj map[string]any
	if err := json.Unmarshal(call.Payload, &obj); err != nil {
		return strings.TrimSpace(string(call.Payload))
	}
	if v, ok := obj[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func payloadList(call kernel.Call, key string) []string {
	if len(call.Payload) == 0 {
		return nil
	}
	var obj map[string]any
	if err := json.Unmarshal(call.Payload, &obj); err != nil {
		return nil
	}
	raw, ok := obj[key]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	case string:
		parts := strings.Split(v, "\n")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if s := strings.TrimSpace(p); s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func payloadInt(call kernel.Call, key string) int {
	if len(call.Payload) == 0 {
		return 0
	}
	var obj map[string]any
	if err := json.Unmarshal(call.Payload, &obj); err != nil {
		return 0
	}
	switch v := obj[key].(type) {
	case float64:
		return int(v)
	}
	return 0
}

type researchAgent struct {
	k *kernel.Kernel
}

func (a *researchAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "research", Name: "Research", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "library", Summary: "First-class notes. Gathers AOS papers and desk findings while the OS is under construction.",
		Capabilities: []string{"research.list", "research.ingest", "note.write"},
		Autostart:    true,
	}
}

func (a *researchAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	seedResearch(k)
	k.Publish("research", "library", "seeded construction notes", map[string]any{"notes": len(k.Notes())})
	return nil
}

func (a *researchAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	return researchInvoke(a.k, call)
}
