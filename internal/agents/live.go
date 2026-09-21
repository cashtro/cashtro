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
	for _, req := range k.Requests() {
		blob := strings.ToLower(req.Raw + " " + req.Title + " " + req.Improved + " " + req.Status)
		if q == "" || strings.Contains(blob, q) {
			hits = append(hits, map[string]string{"kind": "request", "id": strconv.Itoa(req.ID), "name": req.Title})
		}
	}
	_, _ = k.Post("explorer", "research", "search", q)
	return kernel.Result{OK: true, Message: "explorer hit " + strconv.Itoa(len(hits)), Data: hits}, nil
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

type deskAgent struct {
	k *kernel.Kernel
}

func (a *deskAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "desk", Name: "Desk", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "assistant", Summary: "Castro's assistant. Captures every ask, betters it each pass, and publishes what we can and cannot do.",
		Capabilities: []string{"desk.plan", "desk.capture", "desk.better", "desk.list", "desk.done"},
		Autostart:    true,
	}
}

func (a *deskAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	seedDesk(k)
	k.Publish("desk", "plan", "assistant plan live · requests seeded", map[string]any{"requests": len(k.Requests())})
	return nil
}

func (a *deskAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	return deskInvoke(a.k, call)
}
