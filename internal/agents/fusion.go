package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

const fusionWorkerCount = 80

// FusionBureau is the internal compliance and incident-response service shared
// by the cashtro account and the Evolu-Jeunes organization. It has no public
// police authority and cannot take an outbound or destructive action.
type FusionBureau struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Mandate        string         `json:"mandate"`
	Scope          []string       `json:"scope"`
	Authority      string         `json:"authority"`
	OperatingRules []string       `json:"operatingRules"`
	HumanGate      string         `json:"humanGate"`
	WorkerCount    int            `json:"workerCount"`
	Workers        []FusionWorker `json:"workers"`
}

// FusionWorker is an accountable assignment in the bureau roster. Workers are
// records managed by one kernel process, not 80 duplicate goroutines.
type FusionWorker struct {
	ID        string   `json:"id"`
	Unit      string   `json:"unit"`
	Role      string   `json:"role"`
	LawID     string   `json:"lawId"`
	Scope     []string `json:"scope"`
	ReadOnly  bool     `json:"readOnly"`
	ReportsTo string   `json:"reportsTo"`
}

// FusionCheck is a traceable compliance finding. Empty evidence never clears a
// gate; the bureau holds the action for review instead.
type FusionCheck struct {
	LawID       string `json:"lawId"`
	WorkerID    string `json:"workerId"`
	EvidenceRef string `json:"evidenceRef,omitempty"`
	Decision    string `json:"decision"`
	Reason      string `json:"reason"`
	HumanGate   bool   `json:"humanGate"`
}

var fusionUnits = map[string]string{
	"loi-25": "Protection des renseignements",
	"pipeda": "Protection interprovinciale",
	"casl":   "Communications commerciales",
	"charte": "Services en français",
	"lpc":    "Protection du consommateur",
	"image":  "Image et droits",
	"racj":   "Concours publicitaires",
	"amf":    "Marchés et actifs numériques",
}

var fusionRoles = []string{
	"analyste", "enquêteur", "auditeur", "agent de liaison", "coordonnateur de remédiation",
}

// Fusion returns the deterministic 80-worker roster. Ten workers are assigned
// to each legal regime already represented by Loi().
func Fusion() FusionBureau {
	laws := Loi().Laws
	workers := make([]FusionWorker, 0, fusionWorkerCount)
	scope := []string{"cashtro", "Evolu-Jeunes"}
	for _, law := range laws {
		for n := 1; n <= 10; n++ {
			workers = append(workers, FusionWorker{
				ID:        fmt.Sprintf("fusion-%s-%02d", law.ID, n),
				Unit:      fusionUnits[law.ID],
				Role:      fusionRoles[(n-1)%len(fusionRoles)],
				LawID:     law.ID,
				Scope:     append([]string(nil), scope...),
				ReadOnly:  true,
				ReportsTo: "fusion",
			})
		}
	}
	return FusionBureau{
		ID:        "fusion",
		Name:      "The Fusion",
		Mandate:   "Observer, documenter et faire respecter les portes de conformité dans l'écosystème logiciel Evolian et Astro.",
		Scope:     scope,
		Authority: "Contrôle interne du logiciel seulement; aucune autorité policière publique.",
		OperatingRules: []string{
			"Observation et collecte de preuves en lecture seule par défaut.",
			"Aucune attaque, surveillance de personnes, arrestation, sanction ou action physique.",
			"Aucun secret ni renseignement personnel copié dans un rapport.",
			"Une remédiation modifiant un dépôt, un service ou une donnée attend une confirmation humaine.",
		},
		HumanGate:   "comms.allow",
		WorkerCount: len(workers),
		Workers:     workers,
	}
}

func fusionWorkerForLaw(lawID string) (FusionWorker, bool) {
	for _, worker := range Fusion().Workers {
		if worker.LawID == lawID {
			return worker, true
		}
	}
	return FusionWorker{}, false
}

type fusionAgent struct {
	k *kernel.Kernel
}

func (a *fusionAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "fusion", Name: "The Fusion", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role:         "internal-compliance",
		Summary:      "80-worker internal compliance and incident-response bureau for cashtro and Evolu-Jeunes.",
		Capabilities: []string{"fusion.status", "fusion.check", "fusion.remediate"},
		Autostart:    true,
	}
}

func (a *fusionAgent) Boot(_ context.Context, k *kernel.Kernel) error {
	a.k = k
	bureau := Fusion()
	k.Publish("fusion", "roster", "The Fusion roster loaded", map[string]any{
		"workers": bureau.WorkerCount,
		"scope":   bureau.Scope,
	})
	return nil
}

func (a *fusionAgent) Invoke(_ context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "fusion.status":
		return kernel.Result{OK: true, Message: "The Fusion ready", Data: Fusion()}, nil
	case "fusion.check":
		var in struct {
			LawID       string `json:"lawId"`
			EvidenceRef string `json:"evidenceRef"`
		}
		if err := json.Unmarshal(call.Payload, &in); err != nil {
			return kernel.Result{}, err
		}
		worker, ok := fusionWorkerForLaw(strings.TrimSpace(in.LawID))
		if !ok {
			return kernel.Result{OK: false, Message: "unknown law: " + in.LawID}, nil
		}
		evidenceRef := strings.TrimSpace(in.EvidenceRef)
		if found := secretNames(evidenceRef); len(found) > 0 {
			return kernel.Result{OK: false, Message: "secret refused; submit a non-sensitive reference only"}, nil
		}
		check := FusionCheck{
			LawID: in.LawID, WorkerID: worker.ID, EvidenceRef: evidenceRef,
			Decision: "hold", Reason: "preuve requise avant de franchir la porte", HumanGate: true,
		}
		if check.EvidenceRef != "" {
			check.Decision = "review"
			check.Reason = "preuve enregistrée; validation humaine requise"
		}
		a.k.Remember("fusion", check.LawID+" "+check.Decision)
		return kernel.Result{OK: false, Message: check.Reason, Data: check}, nil
	case "fusion.remediate":
		var in struct {
			Finding string `json:"finding"`
			Action  string `json:"action"`
		}
		if err := json.Unmarshal(call.Payload, &in); err != nil {
			return kernel.Result{}, err
		}
		body := strings.TrimSpace(in.Finding + " → " + in.Action)
		if body == "→" || body == "" {
			return kernel.Result{OK: false, Message: "finding and action required"}, nil
		}
		if found := secretNames(body); len(found) > 0 {
			return kernel.Result{OK: false, Message: "secret refused; submit a non-sensitive finding reference only"}, nil
		}
		confirm := a.k.RequestConfirm("fusion", "fusion.remediate", body)
		return kernel.Result{OK: false, Message: "remediation pending human confirmation", Data: confirm}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}
