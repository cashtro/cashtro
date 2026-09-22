package agents

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/cashtro/cashtro/internal/kernel"
)

const (
	deptCorporation = "corporation"
	blockTopic      = "inquisitor.block."
)

// Department is one self-improving seat: a business line, a brain, or the whole corp.
type Department struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Kind       string   `json:"kind"`
	Brain      string   `json:"brain,omitempty"`
	Mandate    string   `json:"mandate"`
	Score      int      `json:"score"`
	Gaps       []string `json:"gaps"`
	Contradict string   `json:"contradict"`
	Optimize   string   `json:"optimize"`
	Blocked    bool     `json:"blocked"`
	Improved   int      `json:"improved"`
	Carrier    string   `json:"carrier"`
}

type inquisitorAgent struct {
	k     *kernel.Kernel
	mu    sync.Mutex
	depts map[string]*Department
}

func (a *inquisitorAgent) Spec() kernel.Spec {
	return stamp(kernel.Spec{
		ID: "inquisitor", Name: "L'Inquisiteur", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role:    "contradict",
		Summary: "Contradicts, tests, and blocks until each department self-improves toward the most optimized option.",
		Capabilities: []string{
			"inquisitor.status",
			"inquisitor.audit",
			"inquisitor.block",
			"inquisitor.release",
			"inquisitor.smarter",
			"inquisitor.optimize",
			"inquisitor.dissent",
			"inquisitor.ask",
			"inquisitor.answer",
		},
		Autostart: true,
	})
}

func (a *inquisitorAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	a.mu.Lock()
	a.depts = seedDepartments()
	a.mu.Unlock()
	// First corporation push: tick every department once. Do not freeze the desk.
	ticked := a.smarterAll(ctx)
	k.Publish("inquisitor", "boot", "L'Inquisiteur online · contradict · test · optimize · no freeze at boot", map[string]any{
		"departments": len(a.snapshot()),
		"ticked":      ticked,
		"blocked":     0,
	})
	k.Remember("inquisitor", "corporation step-up: every department owns its smarter loop. Dissent can block until that department self-improves.")
	return nil
}

func (a *inquisitorAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "inquisitor.status":
		return a.status(), nil
	case "inquisitor.audit":
		return a.audit(call)
	case "inquisitor.block":
		return a.block(call, true)
	case "inquisitor.release":
		return a.block(call, false)
	case "inquisitor.smarter":
		return a.smarter(ctx, call)
	case "inquisitor.optimize":
		return a.optimize(call)
	case "inquisitor.dissent":
		return a.dissent(call)
	case "inquisitor.ask":
		return a.ask(call)
	case "inquisitor.answer":
		return a.answer(call)
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func (a *inquisitorAgent) status() kernel.Result {
	depts := a.snapshot()
	blocked := 0
	for _, d := range depts {
		if d.Blocked {
			blocked++
		}
	}
	return kernel.Result{OK: true, Message: "inquisitor status", Data: map[string]any{
		"departments": depts,
		"blocked":     blocked,
		"asks":        a.k.Asks(),
		"rule":        "contradict, test, block until self-improve, always the most optimized option",
	}}
}

func (a *inquisitorAgent) audit(call kernel.Call) (kernel.Result, error) {
	id := deptID(call)
	a.parkAsk("audit", id, questionsFor("audit", id))
	targets := a.pick(id)
	if len(targets) == 0 {
		return kernel.Result{OK: false, Message: "unknown department: " + id}, nil
	}
	out := make([]Department, 0, len(targets))
	for _, d := range targets {
		d.Contradict = contradict(d)
		d.Optimize = optimizeLine(d)
		out = append(out, *d)
		_, _ = a.k.Post("inquisitor", d.Carrier, "audit", d.Contradict)
		a.k.Remember("inquisitor", "audit "+d.ID+" · "+d.Contradict)
	}
	return kernel.Result{OK: true, Message: "audited " + strconv.Itoa(len(out)), Data: map[string]any{
		"departments": out,
		"asks":        a.k.Asks(),
	}}, nil
}

func (a *inquisitorAgent) block(call kernel.Call, blocked bool) (kernel.Result, error) {
	id := deptID(call)
	if id == "" {
		return kernel.Result{OK: false, Message: "department required"}, nil
	}
	action := "block"
	if !blocked {
		action = "release"
	}
	a.parkAsk(action, id, questionsFor(action, id))
	d, ok := a.get(id)
	if !ok {
		return kernel.Result{OK: false, Message: "unknown department: " + id}, nil
	}
	if !blocked && d.Improved < 1 {
		return kernel.Result{OK: false, Message: id + " has not self-improved yet — smarter first", Data: d}, nil
	}
	a.mu.Lock()
	d.Blocked = blocked
	if blocked {
		d.Improved = 0
		d.Contradict = contradict(d)
	}
	a.mu.Unlock()
	setDeptBlock(a.k, id, blocked)
	verb := "released"
	if blocked {
		verb = "blocked until self-improve"
	}
	_, _ = a.k.Post("inquisitor", d.Carrier, action, id+" "+verb)
	a.k.Remember("inquisitor", id+" "+verb)
	return kernel.Result{OK: true, Message: id + " " + verb, Data: *d}, nil
}

func (a *inquisitorAgent) smarter(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	id := deptID(call)
	a.parkAsk("smarter", id, questionsFor("smarter", id))
	if id == "" || id == deptCorporation {
		n := a.smarterAll(ctx)
		return kernel.Result{OK: true, Message: "smarter corporation · " + strconv.Itoa(n) + " departments", Data: map[string]any{
			"departments": a.snapshot(),
			"ticked":      n,
			"asks":        a.k.Asks(),
		}}, nil
	}
	d, ok := a.get(id)
	if !ok {
		return kernel.Result{OK: false, Message: "unknown department: " + id}, nil
	}
	a.tick(ctx, d)
	return kernel.Result{OK: true, Message: "smarter " + d.ID + " · score " + strconv.Itoa(d.Score), Data: map[string]any{
		"department": *d,
		"ticked":     1,
		"asks":       a.k.Asks(),
	}}, nil
}

func (a *inquisitorAgent) optimize(call kernel.Call) (kernel.Result, error) {
	id := deptID(call)
	proposed := firstNonEmpty(payloadQuery(call, "option"), payloadQuery(call, "prompt"), payloadQuery(call, "action"))
	a.parkAsk("optimize", id, questionsFor("optimize", id))
	targets := a.pick(id)
	if len(targets) == 0 {
		return kernel.Result{OK: false, Message: "unknown department: " + id}, nil
	}
	out := make([]map[string]any, 0, len(targets))
	for _, d := range targets {
		best := optimizeLine(d)
		if proposed != "" {
			best = "Reject \"" + clipRunes(proposed, 80) + "\". " + best
		}
		d.Optimize = best
		d.Contradict = contradict(d)
		out = append(out, map[string]any{
			"department": d.ID,
			"reject":     d.Contradict,
			"choose":     best,
		})
		a.k.Remember("inquisitor", d.ID+" optimize → "+best)
	}
	return kernel.Result{OK: true, Message: "optimize " + strconv.Itoa(len(out)), Data: map[string]any{
		"options": out,
		"asks":    a.k.Asks(),
	}}, nil
}

func (a *inquisitorAgent) dissent(call kernel.Call) (kernel.Result, error) {
	proposed := firstNonEmpty(payloadQuery(call, "action"), payloadQuery(call, "prompt"), payloadQuery(call, "option"), "unspecified action")
	id := deptID(call)
	d, ok := a.get(firstNonEmpty(id, deptCorporation))
	if !ok {
		d, _ = a.get(deptCorporation)
	}
	a.parkAsk("dissent", d.ID, questionsFor("dissent", d.ID))
	rej := "No. \"" + clipRunes(proposed, 120) + "\" is the generic move. " + contradict(d)
	best := optimizeLine(d)
	return kernel.Result{OK: true, Message: "dissent", Data: map[string]any{
		"department": d.ID,
		"reject":     rej,
		"choose":     best,
		"asks":       a.k.Asks(),
	}}, nil
}

func (a *inquisitorAgent) ask(call kernel.Call) (kernel.Result, error) {
	action := firstNonEmpty(payloadQuery(call, "action"), payloadQuery(call, "prompt"), "action")
	id := deptID(call)
	qs := payloadList(call, "questions")
	if len(qs) == 0 {
		qs = questionsFor(action, id)
	}
	ask := a.k.RequestAsk("inquisitor", action, id, qs)
	return kernel.Result{OK: true, Message: "ask #" + strconv.Itoa(ask.ID), Data: ask}, nil
}

func (a *inquisitorAgent) answer(call kernel.Call) (kernel.Result, error) {
	id := payloadInt(call, "id")
	if id == 0 {
		id = payloadInt(call, "askId")
	}
	answers := payloadStringMap(call, "answers")
	if len(answers) == 0 {
		if note := firstNonEmpty(payloadQuery(call, "answer"), payloadQuery(call, "prompt"), payloadQuery(call, "text")); note != "" {
			answers = map[string]string{"note": note}
		}
	}
	ask, err := a.k.AnswerAsk(id, answers)
	if err != nil {
		return kernel.Result{}, err
	}
	a.k.Remember("ask", "answered #"+strconv.Itoa(id)+" · "+ask.Action)
	return kernel.Result{OK: true, Message: "answered ask #" + strconv.Itoa(id), Data: ask}, nil
}

func (a *inquisitorAgent) smarterAll(ctx context.Context) int {
	a.mu.Lock()
	ids := make([]string, 0, len(a.depts))
	for id := range a.depts {
		ids = append(ids, id)
	}
	a.mu.Unlock()
	n := 0
	for _, id := range ids {
		d, ok := a.get(id)
		if !ok {
			continue
		}
		a.tick(ctx, d)
		n++
	}
	return n
}

func (a *inquisitorAgent) tick(ctx context.Context, d *Department) {
	a.mu.Lock()
	filled := ""
	if len(d.Gaps) > 0 {
		filled = d.Gaps[0]
		d.Gaps = append(append([]string{}, d.Gaps[1:]...), nextGap(d.ID, filled))
	} else {
		filled = "no named gap — push a deeper specialist loop"
		d.Gaps = []string{nextGap(d.ID, filled)}
	}
	d.Score++
	d.Improved++
	d.Contradict = contradict(d)
	d.Optimize = optimizeLine(d)
	released := false
	if d.Blocked && d.Improved >= 1 {
		d.Blocked = false
		released = true
	}
	score := d.Score
	carrier := d.Carrier
	id := d.ID
	a.mu.Unlock()
	if released {
		setDeptBlock(a.k, id, false)
	}
	a.k.Remember("smarter", id+" filled · "+filled+" · score "+strconv.Itoa(score))
	_, _ = a.k.Post("inquisitor", carrier, "smarter", id+" · "+filled)
	if id == "teal" {
		_, _ = a.k.Invoke(ctx, "teal", kernel.Call{Capability: "teal.smarter"})
	}
}

func (a *inquisitorAgent) parkAsk(action, line string, qs []string) {
	if a.k == nil {
		return
	}
	a.k.RequestAsk("inquisitor", action, line, qs)
}

func (a *inquisitorAgent) get(id string) (*Department, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	d, ok := a.depts[id]
	return d, ok
}

func (a *inquisitorAgent) pick(id string) []*Department {
	a.mu.Lock()
	defer a.mu.Unlock()
	if id == "" || id == deptCorporation {
		out := make([]*Department, 0, len(a.depts))
		for _, d := range a.depts {
			out = append(out, d)
		}
		return out
	}
	d, ok := a.depts[id]
	if !ok {
		return nil
	}
	return []*Department{d}
}

func (a *inquisitorAgent) snapshot() []Department {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Department, 0, len(a.depts))
	for _, d := range a.depts {
		cp := *d
		cp.Gaps = append([]string(nil), d.Gaps...)
		out = append(out, cp)
	}
	return out
}

func seedDepartments() map[string]*Department {
	out := map[string]*Department{
		deptCorporation: {
			ID: deptCorporation, Name: "Corporation", Kind: "corporation",
			Mandate: "Cashtro OS · 5 brains · 10 lines. Push every seat until the chain returns a specialist result.",
			Gaps:    []string{"the chain is not specialized enough", "Hustler owns 6/10 lines", "no contradict seat before L'Inquisiteur"},
			Carrier: "manager",
		},
		"architecte": {
			ID: "architecte", Name: "L'Architecte", Kind: "brain", Brain: "Architecte",
			Mandate: "Pense · ML/DL. Trading + Giant. Designs must survive a contradict pass.",
			Gaps:    []string{"designs without a contradict pass", "trading methods not specialized per exchange"},
			Carrier: "architect",
		},
		"cartographe": {
			ID: "cartographe", Name: "Le Cartographe", Kind: "brain", Brain: "Cartographe",
			Mandate: "Voit · Graphify. École + maps. Measure must feed self-improve, not sit as a dashboard.",
			Gaps:    []string{"maps without dissent", "école curriculum not auto-updated from Graphify"},
			Carrier: "explorer",
		},
		"forgeron": {
			ID: "forgeron", Name: "Le Forgeron", Kind: "brain", Brain: "Forgeron",
			Mandate: "Construit 24/7. Proximity + Empire. Build loops must self-heal.",
			Gaps:    []string{"builds without self-critique", "usine WP/ACF has no per-site smarter loop"},
			Carrier: "delivery",
		},
		"orfevre": {
			ID: "orfevre", Name: "L'Orfèvre", Kind: "brain", Brain: "Orfèvre",
			Mandate: "Exécute 24/7. Deploy and polish. Execution without an optimize gate ships the generic option.",
			Gaps:    []string{"executes without an optimize gate", "Azure confection not scored per site"},
			Carrier: "deploy",
		},
		"hustler": {
			ID: "hustler", Name: "Le Hustler", Kind: "brain", Brain: "Hustler",
			Mandate: "Empire A à Z. Too many lines. Split specialist loops off the six-line pile.",
			Gaps:    []string{"owns 6/10 lines — over-concentrated", "no specialist sub-loop per line"},
			Carrier: "comms",
		},
	}
	for _, ln := range Lines() {
		out[ln.ID] = &Department{
			ID: ln.ID, Name: ln.Name, Kind: "line", Brain: ln.Brain,
			Mandate: ln.Mandate, Gaps: lineGaps(ln.ID), Carrier: lineCarrier(ln.ID),
		}
	}
	for _, d := range out {
		d.Contradict = contradict(d)
		d.Optimize = optimizeLine(d)
	}
	return out
}

func lineGaps(id string) []string {
	switch id {
	case "proximity":
		return []string{"usine WP/ACF not self-healing", "no per-site QA loop"}
	case "scanapp":
		return []string{"fiche photos still dirty", "no auto-sale loop"}
	case "panda":
		return []string{"white-glove install not productized", "Architecte offer not a repeatable kit"}
	case "nft-giant":
		return []string{"token utility unproven", "art drops not scored against Giant"}
	case "ecole":
		return []string{"curriculum not auto-updated", "WordPress teaching not feeding Proximity QA"}
	case "marketing":
		return []string{"campaigns not closed-loop with Graphify", "corporate CRM not a specialist desk"}
	case "empire":
		return []string{"live sell still depends on external platforms", "studio not a closed sales loop"}
	case "pandora":
		return []string{"intern desk not feeding spawn loop", "Teams scrape unbound — brain starves"}
	case "trading":
		return []string{"bots not specialized per method", "Giant bridge is narrative, not a control loop"}
	case "teal":
		return []string{"smarter only on teal, not the other departments", "voice lane not gated by department block"}
	default:
		return []string{"no specialist loop yet"}
	}
}

func lineCarrier(id string) string {
	switch id {
	case "proximity", "empire":
		return "delivery"
	case "scanapp":
		return "delivery"
	case "panda", "nft-giant", "trading":
		return "architect"
	case "ecole":
		return "research"
	case "marketing":
		return "comms"
	case "pandora", "teal":
		return "teal"
	default:
		return "manager"
	}
}

func contradict(d *Department) string {
	gap := "generic output"
	if len(d.Gaps) > 0 {
		gap = d.Gaps[0]
	}
	return "Stop treating " + d.Name + " as a shared brain. Current lacune: " + gap + ". Push this seat until it closes that gap itself."
}

func optimizeLine(d *Department) string {
	switch d.ID {
	case deptCorporation:
		return "Keep L'Inquisiteur live. Each of the 10 lines + 5 brains runs its own smarter loop. Hustler loses lines as those loops stand up."
	case "hustler":
		return "Peel Scan App, Panda, NFT, marketing, Pandora, Teal into specialist loops. Hustler keeps sell motion only."
	case "teal":
		return "teal.smarter stays local. Corporation smarter is inquisitor.smarter. Voice refuses blocked lines."
	case "proximity":
		return "Per-site smarter: confection, QA, Azure. Not one Forgeron dump."
	case "trading":
		return "One method, one bot, one contradict pass. Giant is a control signal, not a slogan."
	default:
		return "Specialize " + d.Name + " until the next result cannot be produced by another department."
	}
}

func nextGap(id, filled string) string {
	return "next lacune after \"" + clipRunes(filled, 60) + "\" on " + id + " — go one layer more specialist"
}

func questionsFor(action, dept string) []string {
	if dept == "" {
		dept = deptCorporation
	}
	switch action {
	case "block":
		return []string{
			"Quel symptôme concret justifie de bloquer " + dept + " ?",
			"Quelle lacune ce département doit combler avant reprise ?",
			"Quel résultat mesurable débloque ?",
		}
	case "release":
		return []string{
			"Qu'est-ce que " + dept + " a réellement amélioré ?",
			"Pourquoi ce n'est plus la version générique ?",
		}
	case "smarter":
		return []string{
			"Quelle lacune de " + dept + " est la plus urgente ?",
			"Comment ce département s'auto-améliore sans renvoyer au Hustler ?",
			"Quelle spécialisation manque encore après ce tick ?",
		}
	case "optimize":
		return []string{
			"Quelle option actuelle de " + dept + " est la moins spécialisée ?",
			"Quelle option plus étroite produit un meilleur résultat ?",
			"Que contredit L'Inquisiteur dans le plan tel quel ?",
		}
	case "audit":
		return []string{
			"Quel résultat de " + dept + " est trop générique aujourd'hui ?",
			"Quelle preuve manquerait pour dire que le département s'est auto-amélioré ?",
		}
	case "dissent":
		return []string{
			"Quelle action propose-t-on, exactement ?",
			"Pourquoi cette action n'est pas déjà l'option la plus optimisée ?",
			"Quel département devrait porter la version plus spécialisée ?",
		}
	default:
		return []string{
			"Quel contexte manque encore avant d'agir sur " + dept + " ?",
			"Comment cette action s'améliore si on la spécialise davantage ?",
		}
	}
}

func deptID(call kernel.Call) string {
	id := firstNonEmpty(payloadQuery(call, "department"), payloadQuery(call, "dept"), payloadQuery(call, "id"), payloadQuery(call, "line"))
	id = strings.ToLower(strings.TrimSpace(id))
	switch id {
	case "l'architecte", "architecte", "architect":
		return "architecte"
	case "le cartographe", "cartographe":
		return "cartographe"
	case "le forgeron", "forgeron":
		return "forgeron"
	case "l'orfèvre", "l'orfevre", "orfèvre", "orfevre":
		return "orfevre"
	case "le hustler", "hustler":
		return "hustler"
	case "corp", "os", "cashtro":
		return deptCorporation
	}
	return id
}

func setDeptBlock(k *kernel.Kernel, id string, blocked bool) {
	text := "running"
	if blocked {
		text = "blocked"
	}
	k.Remember(blockTopic+id, text)
}

// DepartmentBlocked is true when L'Inquisiteur has blocked this line, its brain, or the whole corporation.
func DepartmentBlocked(k *kernel.Kernel, id string) bool {
	if k == nil {
		return false
	}
	if factBlocked(k, deptCorporation) {
		return true
	}
	if id == "" {
		return false
	}
	if factBlocked(k, id) {
		return true
	}
	if ln, ok := LineByID(id); ok {
		return factBlocked(k, brainSlug(ln.Brain))
	}
	return false
}

func factBlocked(k *kernel.Kernel, id string) bool {
	facts := k.Recall(blockTopic + id)
	for i := len(facts) - 1; i >= 0; i-- {
		if facts[i].Topic == blockTopic+id {
			return facts[i].Text == "blocked"
		}
	}
	return false
}

func brainSlug(name string) string {
	switch strings.ToLower(name) {
	case "architecte":
		return "architecte"
	case "cartographe":
		return "cartographe"
	case "forgeron":
		return "forgeron"
	case "orfèvre", "orfevre":
		return "orfevre"
	case "hustler":
		return "hustler"
	default:
		return strings.ToLower(name)
	}
}
