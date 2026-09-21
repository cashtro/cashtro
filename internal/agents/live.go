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
	if bus := k.Symbols(); bus != nil {
		for _, s := range bus.List() {
			blob := strings.ToLower(s.Name + " " + s.Summary + " " + string(s.On.Kind))
			if q == "" || strings.Contains(blob, q) {
				hits = append(hits, map[string]string{"kind": "symbol", "id": s.ID, "name": s.Name})
			}
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
		{Agent: "research", Source: "cashtro-symbols", URL: "", Claim: "Symbols is our internal Zapier. Triggers and kernel verbs. No subscription.", Quote: "A Symbol is trigger → steps. Connectors are the process table. Webhook catch-hooks live at /api/symbols/{id}/hook."},
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

type chooserAgent struct {
	k *kernel.Kernel
}

func (a *chooserAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "chooser", Name: "Chooser", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "desk", Summary: "You pick. Inbox jobs from Gmail plus the verbs this OS can run. Take or skip — nothing auto-runs.",
		Capabilities: []string{"chooser.list", "chooser.take", "chooser.skip", "chooser.scan"},
		Autostart:    true,
	}
}

func (a *chooserAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	seedChooser(k)
	desk := k.DeskCard()
	k.Publish("chooser", "desk", "you pick · "+strconv.Itoa(desk.PendingN)+" pending", map[string]any{
		"pending": desk.PendingN,
		"inbox":   desk.Scan.Inbox,
		"unread":  desk.Scan.Unread,
	})
	return nil
}

func (a *chooserAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	return chooserInvoke(a.k, call)
}

func chooserInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "chooser.list", "chooser.scan":
		d := k.DeskCard()
		return kernel.Result{OK: true, Message: "pending " + strconv.Itoa(d.PendingN), Data: d}, nil
	case "chooser.take":
		c, err := k.Take(payloadInt(call, "id"))
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "took " + c.Title, Data: c}, nil
	case "chooser.skip":
		c, err := k.Skip(payloadInt(call, "id"))
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "skipped " + c.Title, Data: c}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

type watchAgent struct {
	k *kernel.Kernel
}

func (a *watchAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "watch", Name: "Watch", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role: "night", Summary: "Keeps the desk flowing while things are closed. Pulse and persist. Outbound still needs a human.",
		Capabilities: []string{"watch.status", "watch.close", "watch.open", "watch.pulse"},
		Autostart:    true,
	}
}

func (a *watchAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	k.Publish("watch", "ready", "closed-hours watch online", nil)
	return nil
}

func (a *watchAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	return watchInvoke(a.k, call)
}

func watchInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	note := payloadQuery(call, "note")
	if note == "" {
		note = payloadQuery(call, "prompt")
	}
	switch call.Capability {
	case "watch.status":
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	case "watch.close":
		if note == "" {
			note = "things are closed"
		}
		k.SetClosed(true)
		k.RecordPulse(note)
		keepFlowing(k, note)
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	case "watch.open":
		if note == "" {
			note = "desk open"
		}
		k.SetClosed(false)
		k.RecordPulse(note)
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	case "watch.pulse":
		if note == "" {
			note = "pulse"
		}
		k.RecordPulse(note)
		if k.Closed() {
			keepFlowing(k, note)
		}
		w := k.WatchCard()
		return kernel.Result{OK: true, Message: w.Message, Data: w}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func keepFlowing(k *kernel.Kernel, note string) {
	if cat := k.Catalog(); cat != nil {
		if _, err := cat.Get("closed-hours-flow"); err != nil {
			_, _ = cat.Create(catalog.CreateShip{
				Name:   "Closed-hours flow",
				Client: "Cashtro OS",
				Sector: "ops",
				Stack:  []string{"Go"},
				Notes:  "Keep the work flowing through here while things are closed.",
			})
		}
	}
	k.Remember("watch", "closed hours · "+note)
	_, _ = k.Post("watch", "research", "pulse", "keep gathering while the desk is closed")
	_, _ = k.Post("watch", "planner", "pulse", "keep the line moving")
}
