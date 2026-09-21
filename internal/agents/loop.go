package agents

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/model"
)

// playItem is one cash-facing mandate the forge rotates through.
type playItem struct {
	ID, Name, SKU, Thesis, Ship string
	Lane                        kernel.Lane
}

// playbook is Castro's live Evolu-Jeunes / Proximity wealth line.
// New agentics register here. They do not become a second product.
var playbook = []playItem{
	{ID: "ship-pack", Name: "SKU Ship Pack", SKU: "ship-pack", Lane: kernel.LaneCreate, Ship: "SKU Ship Pack", Thesis: "Frozen $5k–$15k idea→production packs on live WordPress/Next mandates."},
	{ID: "care", Name: "Production Care", SKU: "care", Lane: kernel.LaneGrow, Ship: "Production Care retainers", Thesis: "Monthly care on BTK, MD Clinic, Hypothèque, Éduconnexion."},
	{ID: "wp-factory", Name: "WordPress factory", SKU: "wp-factory", Lane: kernel.LaneBuild, Ship: "WordPress client factory", Thesis: "CSF-LIVE factory: 8–12 pages, ACF, FR/EN, forms, human DNS gate."},
	{ID: "intakedesk", Name: "IntakeDesk Legal", SKU: "intakedesk", Lane: kernel.LaneCreate, Ship: "IntakeDesk Legal", Thesis: "Legal intake as a tenant process for BTK-class firms. Never advise."},
	{ID: "scanapp-crm", Name: "ScanApp CRM lane", SKU: "scanapp-crm", Lane: kernel.LaneBuild, Ship: "ScanApp", Thesis: "Advance ScanApp from concept toward production."},
	{ID: "skills", Name: "Skill files as SKUs", SKU: "skills", Lane: kernel.LaneCreate, Ship: "Skill files as products", Thesis: "Package Evolu-Jeunes delivery skills as billed products."},
	{ID: "agent-seats", Name: "Coding agent seats", SKU: "agent-seats", Lane: kernel.LaneGrow, Ship: "Hosted coding-agent seats", Thesis: "Hosted seats; tokens are COGS. Cashtro earns autonomy. Money never does."},
	{ID: "leads", Name: "Lead gen QC", SKU: "leads", Lane: kernel.LaneGrow, Ship: "Lead gen QC law health finance", Thesis: "CASL-safe scan for QC law, clinic, and mortgage. Parks confirms only."},
	{ID: "forge-loop", Name: "Kimi+GLM forge", SKU: "forge", Lane: kernel.LaneEvolve, Ship: "Kimi K3 GLM wealth loop", Thesis: "Kimi K3 proposes, GLM critiques, score must rise each generation."},
	{ID: "mcp-delivery", Name: "MCP Delivery seat", SKU: "mcp-delivery", Lane: kernel.LaneCreate, Ship: "Delivery MCP hosted seats", Thesis: "Listings are discovery. Cash is hosted seats wrapping kernel verbs."},
}

type pulseAgent struct {
	k *kernel.Kernel
}

func (a *pulseAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "pulse", Name: "Pulse", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role: "heartbeat", Summary: "Never-stop keep-alive. Respawns autostart agentics and ticks the forge.",
		Capabilities: []string{"pulse.keep", "pulse.tick", "pulse.status"}, Autostart: true,
	}
}

func (a *pulseAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	k.Publish("pulse", "loop", "never-stop armed · keep every task running", nil)
	return nil
}

func (a *pulseAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "pulse.status":
		about := a.k.About()
		led := a.k.Ledger()
		return kernel.Result{OK: true, Message: pulseMessage(about, led), Data: map[string]any{
			"about":  about,
			"ledger": led,
		}}, nil
	case "pulse.keep":
		spawned, err := a.k.KeepAlive(ctx)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: "kept " + strconv.Itoa(a.k.About().Running) + " running", Data: map[string]any{
			"spawned": spawned,
			"running": a.k.About().Running,
		}}, nil
	case "pulse.tick":
		spawned, err := a.k.KeepAlive(ctx)
		if err != nil {
			return kernel.Result{}, err
		}
		res, err := a.k.Invoke(ctx, "forge", kernel.Call{Capability: "forge.tick", Payload: call.Payload})
		if err != nil {
			return kernel.Result{}, err
		}
		about := a.k.About()
		return kernel.Result{OK: true, Message: "pulse gen " + strconv.Itoa(about.Generation) + " · score " + strconv.Itoa(about.Score), Data: map[string]any{
			"spawned": spawned,
			"running": about.Running,
			"forge":   res.Data,
			"ledger":  a.k.Ledger(),
		}}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func pulseMessage(about kernel.About, led kernel.Ledger) string {
	return fmt.Sprintf("%d running · gen %d · score %d", about.Running, led.Generation, led.Score)
}

type forgeAgent struct {
	k      *kernel.Kernel
	client *model.Client
}

func (a *forgeAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "forge", Name: "Forge", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "evolve", Summary: "Each tick gets better. Kimi K3 proposes, GLM critiques, kernel enforces a rising score.",
		Capabilities: []string{"forge.tick", "forge.think", "forge.status"}, Autostart: true,
	}
}

func (a *forgeAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	k.Publish("forge", "ready", "dual loop ready · kimi proposes · glm critiques", map[string]any{
		"kimi": model.DefaultKimi,
		"glm":  model.DefaultGLM,
		"play": len(playbook),
	})
	return nil
}

func (a *forgeAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "forge.status":
		led := a.k.Ledger()
		return kernel.Result{OK: true, Message: "gen " + strconv.Itoa(led.Generation) + " · score " + strconv.Itoa(led.Score), Data: led}, nil
	case "forge.tick", "forge.think":
		g, err := a.evolve(ctx, call.Capability == "forge.think")
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: g.Message, Data: g}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func (a *forgeAgent) evolve(ctx context.Context, think bool) (kernel.Generation, error) {
	item := playbook[a.k.Generation()%len(playbook)]
	offer := a.k.UpsertOffer(kernel.Offer{
		ID: item.ID, Name: item.Name, Lane: item.Lane, Thesis: item.Thesis, SKU: item.SKU,
	})
	nextN := a.k.Generation() + 1
	offer, _ = a.k.TouchOffer(item.ID, nextN, 1)

	kimiLine := "Kimi: park " + item.Name + " as " + string(item.Lane) + " · " + item.SKU
	glmLine := "GLM: keep human gate · never raise own autonomy · compound " + item.Name
	bound := false
	if think && model.Bound(a.client) {
		prompt := forgePrompt(a.k, item)
		dual, err := a.client.Dual(ctx, prompt)
		if err == nil {
			bound = true
			kimiLine = dual.Kimi.Content
			glmLine = dual.GLM.Content
		} else {
			glmLine = "GLM unbound-fail: " + err.Error()
		}
	}

	parkPlayShip(a.k, item)
	msg := "gen evolve · " + item.Name + " · " + string(item.Lane)
	g := a.k.CommitGeneration(kernel.Generation{
		Score:   a.k.Score() + 1 + offer.Hits,
		Lane:    item.Lane,
		Offer:   item.Name,
		Kimi:    kimiLine,
		GLM:     glmLine,
		Bound:   bound,
		Message: msg,
	})
	a.k.Remember("forge", msg+" · score "+strconv.Itoa(g.Score))
	_, _ = a.k.Post("forge", "delivery", "evolve", item.Name)
	_, _ = a.k.Post("forge", "wealth", "compound", item.Name)
	return g, nil
}

func forgePrompt(k *kernel.Kernel, item playItem) string {
	var b strings.Builder
	b.WriteString("You are Kimi K3 inside Cashtro OS. Propose ONE next move that creates, builds, grows, or evolves wealth for Castro / Evolu-Jeunes / Proximity. ")
	b.WriteString("Human confirms stay first-class. Do not send outbound mail. Do not raise autonomy.\n")
	b.WriteString("Focus offer: ")
	b.WriteString(item.Name)
	b.WriteString(" (")
	b.WriteString(string(item.Lane))
	b.WriteString(") — ")
	b.WriteString(item.Thesis)
	b.WriteString("\nGeneration ")
	b.WriteString(strconv.Itoa(k.Generation()))
	b.WriteString(" score ")
	b.WriteString(strconv.Itoa(k.Score()))
	b.WriteString(".\n")
	return b.String()
}

func parkPlayShip(k *kernel.Kernel, item playItem) string {
	cat := k.Catalog()
	if cat == nil {
		return ""
	}
	for _, s := range cat.List() {
		if strings.EqualFold(s.Name, item.Ship) || s.ID == item.ID {
			return s.Name
		}
	}
	ship, err := cat.Create(catalog.CreateShip{
		Name:   item.Ship,
		Client: "Cashtro",
		Sector: "wealth",
		Stack:  []string{"Go", "Kimi K3", "GLM"},
		Notes:  item.Thesis,
	})
	if err != nil {
		return ""
	}
	return ship.Name
}

type wealthAgent struct {
	k *kernel.Kernel
}

func (a *wealthAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "wealth", Name: "Wealth", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "ledger", Summary: "Cash-facing offers. Create, build, grow, evolve — ships and retainers, not a second product.",
		Capabilities: []string{"wealth.ledger", "wealth.offer", "wealth.compound"}, Autostart: true,
	}
}

func (a *wealthAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	if len(k.Offers()) == 0 {
		for _, item := range playbook {
			k.UpsertOffer(kernel.Offer{
				ID: item.ID, Name: item.Name, Lane: item.Lane, Thesis: item.Thesis, SKU: item.SKU,
			})
		}
	}
	k.Publish("wealth", "ledger", fmt.Sprintf("board open · %d offers", len(k.Offers())), map[string]any{"offers": len(k.Offers())})
	return nil
}

func (a *wealthAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "wealth.ledger":
		led := a.k.Ledger()
		return kernel.Result{OK: true, Message: "score " + strconv.Itoa(led.Score) + " · " + strconv.Itoa(len(led.Offers)) + " offers", Data: led}, nil
	case "wealth.offer":
		name := payloadQuery(call, "name")
		if name == "" {
			name = payloadQuery(call, "prompt")
		}
		if name == "" {
			return kernel.Result{OK: false, Message: "name required"}, nil
		}
		id := payloadQuery(call, "id")
		if id == "" {
			id = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		}
		lane := kernel.Lane(payloadQuery(call, "lane"))
		o := a.k.UpsertOffer(kernel.Offer{
			ID: id, Name: name, Lane: lane,
			Thesis: payloadQuery(call, "thesis"),
			SKU:    payloadQuery(call, "sku"),
		})
		return kernel.Result{OK: true, Message: "offer " + o.Name, Data: o}, nil
	case "wealth.compound":
		led := a.k.Ledger()
		if len(led.Offers) == 0 {
			return kernel.Result{OK: false, Message: "no offers"}, nil
		}
		res, err := a.k.Invoke(ctx, "forge", kernel.Call{Capability: "forge.tick"})
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: res.Message, Data: a.k.Ledger()}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}
