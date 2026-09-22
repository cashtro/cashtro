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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/model"
)

// Boot builds a kernel with every Cashtro agentic registered and started.
func Boot(opts ...kernel.Option) (*kernel.Kernel, error) {
	k := kernel.New(opts...)
	cat := catalog.New()
	for _, agent := range Builtins(cat, model.BusFromEnv()) {
		k.Register(agent)
	}
	if err := k.Boot(context.Background()); err != nil {
		return nil, err
	}
	return k, nil
}

// Builtins is the process image of Cashtro OS.
func Builtins(cat *catalog.Catalog, router *model.Bus) []kernel.Agent {
	return []kernel.Agent{
		resident(kernel.Spec{
			ID: "init", Name: "Init", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
			Role: "kernel", Summary: "Boots the OS and publishes the manifesto.",
			Capabilities: []string{"os.about"}, Autostart: true,
		}, func(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
			return kernel.Result{OK: true, Message: "about", Data: k.About()}, nil
		}),
		&deliveryAgent{cat: cat},
		&routerAgent{bus: router},
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
		&managerAgent{},
	}
}

// managerAgent is the Epicenter project manager. It oversees all 57 repos
// across cashtro + Evolu-Jeunes, reads the project fiches, queries the
// Graphify map, and coordinates the five brains.
type managerAgent struct {
	k     *kernel.Kernel
	board Board
}

func (a *managerAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "manager", Name: "Manager", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role:    "epicenter",
		Summary: "CEO of the cooperative. Loads the operating board and coordinates departments.",
		Capabilities: []string{
			"manager.status",
			"manager.projects",
			"manager.lines",
			"manager.line",
			"manager.ecosystems",
			"manager.automate",
			"manager.org",
			"manager.loi",
			"manager.chain",
			"manager.fiche",
			"manager.assign",
			"manager.graphify",
		},
		Autostart: true,
	}
}

func (a *managerAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	board, err := LoadBoard()
	if err != nil {
		return err
	}
	a.board = board
	evolu, cashtro := board.counts()
	k.Publish("manager", "epicenter", "operating board loaded", map[string]any{
		"fiches":  board.Fiches,
		"evolu":   evolu,
		"cashtro": cashtro,
		"lines":   len(board.Lines),
		"links":   len(board.Links),
		"agents":  len(board.Agents),
		"gaps":    board.Gaps,
	})
	_, _ = a.runEcosystems()
	return nil
}

func (a *managerAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "manager.status":
		evolu, cashtro := a.board.counts()
		return kernel.Result{OK: true, Message: "operating board", Data: map[string]any{
			"version": a.board.Version,
			"fiches":  a.board.Fiches,
			"repos":   len(a.board.Repos),
			"evolu":   evolu,
			"cashtro": cashtro,
			"lines":   len(a.board.Lines),
			"links":   len(a.board.Links),
			"agents":  len(a.board.Agents),
			"gaps":    a.board.Gaps,
			"graph":   a.board.Graph,
			"rules":   a.board.Rules,
		}}, nil

	case "manager.lines":
		return kernel.Result{OK: true, Message: "business lines under the main brain", Data: Lines()}, nil

	case "manager.line":
		id := payloadQuery(call, "id")
		if id == "" {
			id = payloadQuery(call, "line")
		}
		ln, ok := LineByID(id)
		if !ok {
			return kernel.Result{OK: false, Message: "unknown line: " + id}, nil
		}
		return kernel.Result{OK: true, Message: ln.Name, Data: ln}, nil

	case "manager.ecosystems":
		return kernel.Result{OK: true, Message: "ecosystem connections", Data: Links()}, nil

	case "manager.org":
		return kernel.Result{OK: true, Message: "cooperative", Data: a.board.Org}, nil

	case "manager.loi":
		return kernel.Result{OK: true, Message: "quebec and canada gates", Data: a.board.Compliance}, nil

	case "manager.automate":
		posted, err := a.runEcosystems()
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "ecosystems connected", Data: map[string]any{
			"links":  Links(),
			"agents": posted,
		}}, nil

	case "manager.projects":
		evolu, cashtro := a.board.counts()
		return kernel.Result{OK: true, Message: "projects on the operating board", Data: map[string]any{
			"evolu_jeunes": evolu,
			"cashtro":      cashtro,
			"total":        len(a.board.Repos),
			"repos":        a.board.Repos,
		}}, nil

	case "manager.chain":
		id := payloadQuery(call, "id")
		if id == "" {
			id = payloadQuery(call, "line")
		}
		steps, ok := Chain(id)
		if !ok {
			return kernel.Result{OK: false, Message: "unknown chain: " + id}, nil
		}
		posted := make([]string, 0, len(steps))
		for _, step := range steps {
			body := step.Line + " #" + strconv.Itoa(step.Order) + " · " + step.Do
			if _, err := a.k.Post("manager", step.Agent, "chain", body); err != nil {
				return kernel.Result{}, err
			}
			posted = append(posted, step.Agent)
		}
		return kernel.Result{OK: true, Message: id + " chain posted", Data: map[string]any{
			"line": id, "steps": steps, "posted": posted,
		}}, nil

	case "manager.fiche":
		var fiche Fiche
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &fiche); err != nil {
				return kernel.Result{}, err
			}
		}
		ready, why := fiche.Ready()
		if ready {
			_, _ = a.k.Post("manager", "comms", "fiche", fiche.Code+" prête")
			_, _ = a.k.Post("manager", "operator", "fiche", fiche.Code+" vers Proximity et Empire")
		}
		return kernel.Result{OK: ready, Message: why, Data: fiche}, nil

	case "manager.assign":
		brain := payloadQuery(call, "brain")
		repo := payloadQuery(call, "repo")
		if brain == "" || repo == "" {
			return kernel.Result{OK: false, Message: "brain and repo required"}, nil
		}
		_, _ = a.k.Post("manager", "planner", "assign", brain+":"+repo)
		a.k.Remember("epicenter", "assigned "+repo+" to "+brain)
		return kernel.Result{OK: true, Message: "assigned " + repo + " → " + brain, Data: map[string]any{
			"brain": brain, "repo": repo,
		}}, nil

	case "manager.graphify":
		query := payloadQuery(call, "query")
		if query == "" {
			query = payloadQuery(call, "prompt")
		}
		_, _ = a.k.Post("manager", "explorer", "graphify", query)
		a.k.Remember("epicenter", "graphify query: "+query)
		return kernel.Result{OK: true, Message: "graphify query: " + query, Data: map[string]any{
			"query": query,
			"nodes": 39593,
			"edges": 99285,
		}}, nil

	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

// runEcosystems posts every ecosystem link to the agentics that carry it.
func (a *managerAgent) runEcosystems() ([]string, error) {
	seen := map[string]bool{}
	var posted []string
	for _, ln := range Links() {
		body := ln.From + " → " + ln.To + " · " + ln.Via
		for _, id := range ln.Agents {
			if _, err := a.k.Post("manager", id, "ecosystem", body); err != nil {
				return nil, err
			}
			if !seen[id] {
				seen[id] = true
				posted = append(posted, id)
			}
		}
		a.k.Remember("ecosystem", body)
	}
	_, _ = a.k.Post("manager", "explorer", "graphify", "ecosystem connections")
	return posted, nil
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
	bus *model.Bus
}

func (a *routerAgent) Spec() kernel.Spec {
	mode := kernel.ModeResident
	summary := "Both lanes: Ollama internal, Kimi K3 + GLM 5.3 max external."
	if a.bus.Live() {
		mode = kernel.ModeLive
		summary = "Router live on both lanes: Ollama internal, Kimi K3 + GLM 5.3 max external."
	}
	return kernel.Spec{
		ID: "router", Name: "Router", Kind: kernel.KindSystem, Mode: mode,
		Role: "model", Summary: summary,
		Capabilities: []string{"model.status", "model.chat"},
		Autostart:    true,
	}
}

func (a *routerAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	st := a.bus.Card()
	k.Publish("router", "bind", st.Hint, map[string]any{
		"internal": st.Internal.Bound,
		"external": st.External.Bound,
		"ollama":   st.Internal.Model,
		"kimi":     st.External.Model,
		"glm":      st.External.Also,
	})
	return nil
}

func (a *routerAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "model.status":
		st := a.bus.Card()
		return kernel.Result{OK: true, Message: st.Hint, Data: st}, nil
	case "model.chat":
		var in struct {
			Route    string          `json:"route"`
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
		out, err := a.bus.Chat(ctx, model.ChatRequest{Route: in.Route, Model: in.Model, Messages: msgs})
		if err != nil {
			if errors.Is(err, model.ErrUnbound) {
				st := a.bus.Card()
				return kernel.Result{OK: false, Message: st.Hint, Data: st}, nil
			}
			return kernel.Result{}, err
		}
		names := make([]string, 0, len(out))
		for _, part := range out {
			names = append(names, part.Route+":"+part.Model)
		}
		return kernel.Result{OK: true, Message: "routed " + strings.Join(names, ", "), Data: out}, nil
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
