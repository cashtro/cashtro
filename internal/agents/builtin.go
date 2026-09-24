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
			Role: "kernel", Summary: "Voltron boots the OS and can start the same cycle Epicenter Einstein runs.",
			Capabilities: []string{"os.about", "os.cycle", "os.engine", "os.layers"}, Autostart: true,
		}, initInvoke),
		&deliveryAgent{cat: cat},
		&routerAgent{bus: router},
		&researchAgent{},
		resident(kernel.Spec{
			ID: "explorer", Name: "Explorer", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "search", Summary: "Searches processes, ships, and research notes on the desk.",
			Capabilities: []string{"explorer.search"}, Autostart: true,
		}, explorerInvoke),
		resident(kernel.Spec{
			ID: "operator", Name: "Operator", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "computer-use", Summary: "Travaille en local sur un dépôt du roster. Ne pousse rien.",
			Capabilities: []string{"operator.work", "operator.browse"}, Autostart: true,
		}, operatorInvoke),
		resident(kernel.Spec{
			ID: "reviewer", Name: "Reviewer", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "qa", Summary: "Filtre encore : garde l'option la plus courte.",
			Capabilities: []string{"reviewer.watch"}, Autostart: true,
		}, reviewerInvoke),
		resident(kernel.Spec{
			ID: "architect", Name: "Architect", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "design", Summary: "Relie la division, la fonction et le produit avant le cycle.",
			Capabilities: []string{"architect.plan"}, Autostart: true,
		}, architectInvoke),
		resident(kernel.Spec{
			ID: "deploy", Name: "Deploy", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "release", Summary: "Enregistre une sortie locale. La production attend comms.allow et ne part pas d'ici.",
			Capabilities: []string{"deploy.release"}, Autostart: true,
		}, deployInvoke),
		resident(kernel.Spec{
			ID: "security", Name: "Security", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "guard", Summary: "Filtre les secrets et la sortie d'un site déjà en ligne. Pas d'attaque.",
			Capabilities: []string{"security.triage"}, Autostart: true,
		}, securityInvoke),
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
			ID: "investigator", Name: "Investigator", Kind: kernel.KindUser, Mode: kernel.ModeLive,
			Role: "incident", Summary: "Nomme la ligne, le département et les dépôts voisins. Ne lit pas un secret.",
			Capabilities: []string{"investigator.trace"}, Autostart: true,
		}, investigatorInvoke),
		&managerAgent{},
	}
}

// managerAgent is the Epicenter Einstein project manager. It oversees all 57 repos
// across cashtro + Evolu-Jeunes, reads the project fiches, queries the
// Graphify map, and coordinates the five brains.
type managerAgent struct {
	k        *kernel.Kernel
	board    Board
	moves    []Move
	congress Congress
}

func (a *managerAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "manager", Name: "Manager", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role:    "instinct",
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
			"manager.fusion",
			"manager.contradict",
			"manager.ask",
			"manager.improve",
			"manager.chain",
			"manager.fiche",
			"manager.assign",
			"manager.graphify",
			"manager.ledger",
			"manager.watch",
			"manager.flow",
			"manager.run",
			"manager.engine",
			"manager.blueprint",
			"manager.funds",
			"manager.tax",
			"manager.check",
			"manager.lens",
			"manager.crm",
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
	k.Publish("manager", "instinct", "operating board loaded", map[string]any{
		"fiches":  board.Fiches,
		"evolu":   evolu,
		"cashtro": cashtro,
		"lines":   len(board.Lines),
		"links":   len(board.Links),
		"agents":  len(board.Agents),
		"gaps":    board.Gaps,
	})
	_, _ = a.runEcosystems()
	a.moves = AppendMove(nil, "boot", StrategyQuestions()[0].Ask, "Le tableau est relu. Chaque ligne est un nœud. Les trous sont publiés.")
	k.Publish("manager", "instinct", "move chain extended", map[string]any{
		"index": a.moves[0].Index,
		"hash":  a.moves[0].Hash,
	})
	return nil
}

func (a *managerAgent) note(action string, self SelfCheck) {
	q := StrategyQuestions()[0].Ask
	if len(self.Asked) > 0 {
		q = self.Asked[0].Ask
	}
	coup := self.Answered["coup"]
	if coup == "" {
		coup = "coup non joué"
	}
	a.moves = AppendMove(a.moves, action, q, coup)
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

	case "manager.fusion":
		if a.congress.House == "" {
			a.congress = OpenCongress()
		}
		if q := payloadQuery(call, "question"); q != "" {
			next, bill, ok := a.congress.Introduce(q)
			if !ok {
				return kernel.Result{OK: false, Message: "la question ne devient pas une loi"}, nil
			}
			a.congress = next
			a.k.Remember("congress", "question: "+q)
			return kernel.Result{OK: true, Message: "projet de loi déposé. pas encore une loi", Data: map[string]any{
				"bill": bill, "laws": len(a.congress.Laws),
			}}, nil
		}
		if name := payloadQuery(call, "sanction"); name != "" {
			next, told, ok := a.congress.Sanction(name, payloadQuery(call, "law"), payloadQuery(call, "fact"), payloadQuery(call, "story"))
			if !ok {
				return kernel.Result{OK: false, Message: "une sanction sans nom ni histoire n'est pas gardée"}, nil
			}
			a.congress = next
			a.k.Remember("sanction", told.Name+": "+told.Story)
			return kernel.Result{OK: true, Message: "sanction nommée gardée", Data: told}, nil
		}
		if id := payloadQuery(call, "bill"); id != "" {
			next, law, ok := a.congress.Pass(id, payloadQuery(call, "rule"))
			if !ok {
				return kernel.Result{OK: false, Message: "la loi n'est pas passée. la question reste ouverte"}, nil
			}
			a.congress = next
			a.k.Remember("congress", law.ID+": "+law.Rules[0])
			return kernel.Result{OK: true, Message: "loi passée", Data: map[string]any{
				"law": law, "laws": len(a.congress.Laws),
			}}, nil
		}
		ok, msg := a.congress.Enforce()
		a.k.Remember("fusion", msg)
		return kernel.Result{OK: ok, Message: msg, Data: map[string]any{
			"laws": a.congress.Laws, "bills": a.congress.Bills, "houses": Houses(), "body": FusionBody(), "memory": a.k.Recall("fusion"),
		}}, nil

	case "manager.watch":
		self := AskSelf("watch", nil)
		a.note("watch", self)
		if len(self.Open) > 0 {
			return kernel.Result{OK: false, Message: askMessage(self.Open), Data: self.Open}, nil
		}
		covered := 0
		for _, ln := range Lines() {
			if _, ok := Chain(ln.ID); ok {
				covered++
			}
		}
		_, _ = a.k.Post("manager", "memory", "watch", "tableau relu, aucun déploiement")
		return kernel.Result{OK: len(a.board.Gaps) == 0 && covered == len(Lines()), Message: "veille", Data: map[string]any{
			"gaps":   a.board.Gaps,
			"hors":   []string{"cashtro/agence"},
			"chains": covered,
			"lines":  len(Lines()),
		}}, nil

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

	case "manager.ask":
		action := payloadQuery(call, "action")
		if action == "" {
			action = payloadQuery(call, "id")
		}
		check := AskSelf(action, payloadAnswers(call))
		a.note("ask:"+action, check)
		for id, ans := range check.Answered {
			a.k.Remember(action, id+": "+ans)
		}
		return kernel.Result{OK: len(check.Open) == 0, Message: askMessage(check.Open), Data: check}, nil

	case "manager.ledger":
		return kernel.Result{OK: ChainIntact(a.moves), Message: "chaîne des coups", Data: a.moves}, nil

	case "manager.contradict":
		self := AskSelf("contradict", payloadAnswers(call))
		a.note("contradict", self)
		if len(self.Open) > 0 {
			return kernel.Result{OK: false, Message: askMessage(self.Open), Data: self.Open}, nil
		}
		var proposal Proposal
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &proposal); err != nil {
				return kernel.Result{}, err
			}
		}
		verdict := Contradict(proposal)
		_, _ = a.k.Post("manager", "reviewer", "contradict", verdict.Attack)
		_, _ = a.k.Post("manager", "explorer", "contradict", verdict.Attack)
		if verdict.Ready {
			a.k.Remember("optimisation", proposal.Subject+" → "+verdict.Best)
		}
		return kernel.Result{OK: verdict.Ready, Message: verdict.Attack, Data: verdict}, nil

	case "manager.improve":
		deptID := payloadQuery(call, "department")
		if deptID == "" {
			deptID = payloadQuery(call, "id")
		}
		draft := payloadQuery(call, "draft")
		improved, ok := Improve(deptID, draft)
		if !ok {
			return kernel.Result{OK: false, Message: "unknown department: " + deptID}, nil
		}
		dept, _ := DeptByID(deptID)
		_, _ = a.k.Post("manager", dept.Chief, "improve", improved.Result)
		if !improved.Specialized {
			a.k.Remember(deptID, "lacunes comblées: "+strings.Join(improved.Lacunes, ", "))
		}
		return kernel.Result{OK: true, Message: improved.Result, Data: improved}, nil

	case "manager.chain":
		id := payloadQuery(call, "id")
		if id == "" {
			id = payloadQuery(call, "line")
		}
		steps, ok := Chain(id)
		if !ok {
			return kernel.Result{OK: false, Message: "unknown chain: " + id}, nil
		}
		self := AskSelf(id, payloadAnswers(call))
		a.note("chain:"+id, self)
		if len(self.Open) > 0 {
			return kernel.Result{OK: false, Message: askMessage(self.Open), Data: self.Open}, nil
		}
		posted := make([]string, 0, len(steps))
		specialized := make([]Improvement, 0, len(steps))
		for _, step := range steps {
			dept, _ := DeptByAgent(step.Agent)
			improved, _ := Improve(dept.ID, step.Do)
			if !improved.Specialized {
				improved, _ = Improve(dept.ID, improved.Result)
			}
			if !improved.Specialized {
				return kernel.Result{OK: false, Message: "lacune non comblée: " + dept.ID, Data: improved}, nil
			}
			body := step.Line + " #" + strconv.Itoa(step.Order) + " · " + improved.Result
			if _, err := a.k.Post("manager", step.Agent, "chain", body); err != nil {
				return kernel.Result{}, err
			}
			posted = append(posted, step.Agent)
			specialized = append(specialized, improved)
		}
		verdict := Contradict(Proposal{
			Subject: id,
			Options: []Option{
				{Name: "rallonger la chaîne", Cost: 3, Risk: 2, Steps: len(steps) + 3},
				{Name: "garder cette chaîne", Cost: 1, Risk: 1, Steps: len(steps)},
			},
		})
		_, _ = a.k.Post("manager", "reviewer", "contradict", verdict.Attack)
		if !verdict.Ready || verdict.Best != "garder cette chaîne" {
			return kernel.Result{OK: false, Message: verdict.Attack, Data: verdict}, nil
		}
		return kernel.Result{OK: true, Message: id + " chain specialized", Data: map[string]any{
			"line": id, "steps": steps, "posted": posted, "specialized": specialized,
			"asked": self.Asked, "verdict": verdict,
		}}, nil

	case "manager.fiche":
		var fiche Fiche
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &fiche); err != nil {
				return kernel.Result{}, err
			}
		}
		ready, why := fiche.Ready()
		self := AskSelf("fiche", payloadAnswers(call))
		a.note("fiche", self)
		if ready {
			if len(self.Open) > 0 {
				return kernel.Result{OK: false, Message: askMessage(self.Open), Data: self.Open}, nil
			}
			_, _ = a.k.Post("manager", "comms", "fiche", fiche.Code+" prête")
			_, _ = a.k.Post("manager", "operator", "fiche", fiche.Code+" vers le marketplace et Empire")
		}
		return kernel.Result{OK: ready, Message: why, Data: fiche}, nil

	case "manager.assign":
		self := AskSelf("assign", payloadAnswers(call))
		a.note("assign", self)
		if len(self.Open) > 0 {
			return kernel.Result{OK: false, Message: askMessage(self.Open), Data: self.Open}, nil
		}
		brain := payloadQuery(call, "brain")
		repo := payloadQuery(call, "repo")
		if brain == "" || repo == "" {
			return kernel.Result{OK: false, Message: "brain and repo required"}, nil
		}
		_, _ = a.k.Post("manager", "planner", "assign", brain+":"+repo)
		a.k.Remember("instinct", "assigned "+repo+" to "+brain)
		return kernel.Result{OK: true, Message: "assigned " + repo + " → " + brain, Data: map[string]any{
			"brain": brain, "repo": repo,
		}}, nil

	case "manager.graphify":
		query := payloadQuery(call, "query")
		if query == "" {
			query = payloadQuery(call, "prompt")
		}
		_, _ = a.k.Post("manager", "explorer", "graphify", query)
		a.k.Remember("instinct", "graphify query: "+query)
		return kernel.Result{OK: true, Message: "graphify query: " + query, Data: map[string]any{
			"query": query,
			"nodes": 52527,
			"edges": 133783,
		}}, nil

	case "manager.flow":
		return kernel.Result{OK: true, Message: "division, fonction, produit", Data: Flows()}, nil

	case "manager.run":
		return RunCycle(a.k, call)

	case "manager.engine":
		return EngineRun(a.k, call)

	case "manager.blueprint":
		return kernel.Result{OK: true, Message: "blueprint", Data: MakeBlueprint(a.k)}, nil

	case "manager.funds":
		return FundsInvoke(a.k, call)

	case "manager.tax":
		return TaxInvoke(call)

	case "manager.check":
		return CheckInvoke(a.k, call)

	case "manager.lens":
		return kernel.Result{OK: true, Message: "lentille", Data: Lens(a.k)}, nil

	case "manager.crm":
		return CRMInvoke(a.k, call)

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
