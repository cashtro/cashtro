// Package agents loads the Cashtro OS agentics into the kernel.
//
// Live agentics execute work now. Resident agentics are first-class
// processes with a published contract — they boot, appear on the desk,
// and accept invokes — waiting for a worker bind. New agentics we create
// in this workspace register here. They do not become a second product.
package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/model"
)

// Boot builds a kernel with every Cashtro agentic registered and started.
func Boot(opts ...kernel.Option) (*kernel.Kernel, error) {
	k := kernel.New(opts...)
	cat := catalog.New()
	for _, agent := range Builtins(cat, model.FromEnv()) {
		k.Register(agent)
	}
	if err := k.Boot(context.Background()); err != nil {
		return nil, err
	}
	return k, nil
}

// Builtins is the process image of Cashtro OS.
func Builtins(cat *catalog.Catalog, router *model.Client) []kernel.Agent {
	return []kernel.Agent{
		resident(kernel.Spec{
			ID: "init", Name: "Init", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
			Role: "kernel", Summary: "Boots the OS and publishes the manifesto.",
			Capabilities: []string{"os.about"}, Autostart: true,
		}, func(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
			return kernel.Result{OK: true, Message: "about", Data: k.About()}, nil
		}),
		&deliveryAgent{cat: cat},
		&routerAgent{client: router},
		&researchAgent{},
		resident(kernel.Spec{
			ID: "explorer", Name: "Explorer", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "search", Summary: "Searches processes, ships, and research notes on the desk.",
			Capabilities: []string{"explorer.search"}, Autostart: true,
		}, explorerInvoke),
		resident(kernel.Spec{
			ID: "operator", Name: "Operator", Kind: kernel.KindUser, Mode: kernel.ModeResident,
			Role: "computer-use", Summary: "Drives the browser and desktop the way a shipper would.",
			Capabilities: []string{"operator.browse"}, Autostart: true,
		}, nil),
		resident(kernel.Spec{
			ID: "reviewer", Name: "Reviewer", Kind: kernel.KindUser, Mode: kernel.ModeResident,
			Role: "qa", Summary: "Reads walkthrough video and screenshot artifacts before we call a ship done.",
			Capabilities: []string{"reviewer.watch"}, Autostart: true,
		}, nil),
		resident(kernel.Spec{
			ID: "architect", Name: "Architect", Kind: kernel.KindUser, Mode: kernel.ModeResident,
			Role: "design", Summary: "Shapes AI-powered apps, agents, and workflows before they hit the line.",
			Capabilities: []string{"architect.plan"}, Autostart: true,
		}, nil),
		resident(kernel.Spec{
			ID: "deploy", Name: "Deploy", Kind: kernel.KindUser, Mode: kernel.ModeResident,
			Role: "release", Summary: "Owns CI, preview, and production promotion.",
			Capabilities: []string{"deploy.release"}, Autostart: true,
		}, nil),
		resident(kernel.Spec{
			ID: "security", Name: "Security", Kind: kernel.KindUser, Mode: kernel.ModeResident,
			Role: "guard", Summary: "Triage CVE and SAST findings on the board before they ship.",
			Capabilities: []string{"security.triage"}, Autostart: true,
		}, nil),
		resident(kernel.Spec{
			ID: "memory", Name: "Memory", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "recall", Summary: "Episodic facts. Store and recall without a model.",
			Capabilities: []string{"memory.store", "memory.recall"}, Autostart: true,
		}, memoryInvoke),
		resident(kernel.Spec{
			ID: "comms", Name: "Comms", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "signal", Summary: "Outbound only after a human confirm. No send until allow.",
			Capabilities: []string{"comms.send", "comms.pending", "comms.allow", "comms.deny"}, Autostart: true,
		}, commsInvoke),
		resident(kernel.Spec{
			ID: "planner", Name: "Planner", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "backlog", Summary: "Turns a goal into idea-stage ships on the delivery line.",
			Capabilities: []string{"planner.backlog"}, Autostart: true,
		}, plannerInvoke),
		resident(kernel.Spec{
			ID: "investigator", Name: "Investigator", Kind: kernel.KindUser, Mode: kernel.ModeResident,
			Role: "incident", Summary: "Traces a failing check or a live incident back to the blast radius.",
			Capabilities: []string{"investigator.trace"}, Autostart: true,
		}, nil),
		&pulseAgent{},
		&forgeAgent{client: router},
		&wealthAgent{},
	}
}

type deliveryAgent struct {
	cat *catalog.Catalog
}

func (a *deliveryAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "delivery", Name: "Delivery", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "ship", Summary: "The live line: idea → concept → production.",
		Capabilities: []string{"delivery.list", "delivery.create", "delivery.advance", "delivery.profile"},
		Autostart:    true,
	}
}

func (a *deliveryAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	k.AttachCatalog(a.cat)
	ships := a.cat.List()
	var idea, concept, prod int
	for _, s := range ships {
		switch s.Stage {
		case catalog.StageIdea:
			idea++
		case catalog.StageConcept:
			concept++
		case catalog.StageProduction:
			prod++
		}
	}
	k.Publish("delivery", "board", fmt.Sprintf("board attached · %d idea · %d concept · %d production", idea, concept, prod), map[string]any{
		"idea": idea, "concept": concept, "production": prod,
	})
	return nil
}

func (a *deliveryAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "delivery.list":
		return kernel.Result{OK: true, Message: "listed ships", Data: a.cat.List()}, nil
	case "delivery.profile":
		return kernel.Result{OK: true, Message: "profile", Data: a.cat.Profile()}, nil
	case "delivery.create":
		var in catalog.CreateShip
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &in); err != nil {
				return kernel.Result{}, err
			}
		}
		ship, err := a.cat.Create(in)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "parked " + ship.Name + " in idea", Data: ship}, nil
	case "delivery.advance":
		var in struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(call.Payload, &in); err != nil {
			return kernel.Result{}, err
		}
		ship, err := a.cat.Advance(in.ID)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: ship.Name + " → " + string(ship.Stage), Data: ship}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

type routerAgent struct {
	client *model.Client
}

func (a *routerAgent) Spec() kernel.Spec {
	mode := kernel.ModeResident
	summary := "Optional OpenRouter bus. Kimi K3 proposes, GLM critiques. Kernel boots without a key."
	if model.Bound(a.client) {
		mode = kernel.ModeLive
		summary = "OpenRouter bound. Kimi K3 proposes, GLM critiques, score must rise."
	}
	return kernel.Spec{
		ID: "router", Name: "Router", Kind: kernel.KindSystem, Mode: mode,
		Role: "model", Summary: summary,
		Capabilities: []string{"model.status", "model.chat", "model.dual"},
		Autostart:    true,
	}
}

func (a *routerAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	st := model.Card(a.client)
	k.Publish("router", "bind", st.Hint, map[string]any{"bound": st.Bound, "model": st.Model})
	return nil
}

func (a *routerAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "model.status":
		st := model.Card(a.client)
		return kernel.Result{OK: true, Message: st.Hint, Data: st}, nil
	case "model.chat":
		if !model.Bound(a.client) {
			st := model.Card(a.client)
			return kernel.Result{OK: false, Message: st.Hint, Data: st}, nil
		}
		var in struct {
			Prompt   string          `json:"prompt"`
			Model    string          `json:"model"`
			Messages []model.Message `json:"messages"`
		}
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &in); err != nil {
				return kernel.Result{}, err
			}
		}
		msgs := in.Messages
		if len(msgs) == 0 && in.Prompt != "" {
			msgs = []model.Message{{Role: "user", Content: in.Prompt}}
		}
		out, err := a.client.Chat(ctx, model.ChatRequest{Model: in.Model, Messages: msgs})
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "routed " + out.Model, Data: out}, nil
	case "model.dual":
		if !model.Bound(a.client) {
			st := model.Card(a.client)
			return kernel.Result{OK: false, Message: st.Hint, Data: st}, nil
		}
		prompt := payloadQuery(call, "prompt")
		if prompt == "" {
			prompt = payloadQuery(call, "goal")
		}
		if prompt == "" {
			prompt = "Propose one Cashtro OS move that creates, builds, grows, or evolves wealth. Keep the human gate."
		}
		out, err := a.client.Dual(ctx, prompt)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "kimi+glm dual", Data: out}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

type residentAgent struct {
	spec kernel.Spec
	fn   func(k *kernel.Kernel, call kernel.Call) (kernel.Result, error)
	k    *kernel.Kernel
}

func resident(spec kernel.Spec, fn func(k *kernel.Kernel, call kernel.Call) (kernel.Result, error)) *residentAgent {
	return &residentAgent{spec: spec, fn: fn}
}

func (a *residentAgent) Spec() kernel.Spec { return a.spec }

func (a *residentAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	return nil
}

func (a *residentAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	if a.fn != nil {
		return a.fn(a.k, call)
	}
	return kernel.Result{
		OK:      true,
		Message: a.spec.Name + " resident · " + call.Capability + " ready to bind",
		Data: map[string]any{
			"agent":      a.spec.ID,
			"capability": call.Capability,
			"mode":       a.spec.Mode,
			"role":       a.spec.Role,
			"note":       "Process is on the desk. Bind a worker to execute this verb.",
		},
	}, nil
}
