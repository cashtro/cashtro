package server

import (
	"errors"
	"net/http"

	"github.com/cashtro/cashtro/internal/fleet"
	"github.com/cashtro/cashtro/internal/kernel"
)

type companySummary struct {
	fleet.Company
	OwnAgents int  `json:"ownAgents"`
	UsingMine bool `json:"usingMine"`
}

type boundAgent struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Mode         string   `json:"mode"`
	Summary      string   `json:"summary"`
	CashtroID    string   `json:"cashtroId,omitempty"`
	Capabilities []string `json:"capabilities"`
	Inherited    bool     `json:"inherited"`
	Source       string   `json:"source"`
}

type companyCard struct {
	fleet.Company
	UsingMine bool         `json:"usingMine"`
	Agents    []boundAgent `json:"agents"`
}

func (s *api) tenants() *fleet.Fleet {
	if f, ok := s.k.Tenants().(*fleet.Fleet); ok && f != nil {
		return f
	}
	return fleet.New()
}

func (s *api) listCompanies(w http.ResponseWriter, r *http.Request) {
	fl := s.tenants()
	out := make([]companySummary, 0)
	for _, c := range fl.List() {
		own, _ := fl.OwnAgents(c.ID)
		out = append(out, companySummary{Company: c, OwnAgents: len(own), UsingMine: fl.UsingMine(c.ID)})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *api) createCompany(w http.ResponseWriter, r *http.Request) {
	var in fleet.CreateCompany
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	co, err := s.tenants().Create(in)
	if err != nil {
		writeFleetError(w, err)
		return
	}
	s.k.Publish("fleet", "company.create", co.Name+" joined the ecosystem", map[string]any{"id": co.ID, "useMine": co.UseMine})
	writeJSON(w, http.StatusCreated, s.card(co.ID))
}

func (s *api) getCompany(w http.ResponseWriter, r *http.Request) {
	card := s.card(r.PathValue("id"))
	if card.ID == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "company not found"})
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (s *api) addCompanyAgent(w http.ResponseWriter, r *http.Request) {
	var in fleet.CreateAgentic
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	ag, err := s.tenants().AddAgentic(r.PathValue("id"), in)
	if err != nil {
		writeFleetError(w, err)
		return
	}
	s.k.Publish("fleet", "agentic.add", ag.Name+" bound to "+ag.CompanyID, map[string]any{"id": ag.ID, "cashtroId": ag.CashtroID})
	writeJSON(w, http.StatusCreated, ag)
}

func (s *api) invokeCompanyAgent(w http.ResponseWriter, r *http.Request) {
	var call kernel.Call
	if err := decodeJSON(r, &call); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	companyID := r.PathValue("id")
	agentID := r.PathValue("aid")
	fl := s.tenants()
	if _, err := fl.Get(companyID); err != nil {
		writeFleetError(w, err)
		return
	}

	target := ""
	own, _ := fl.OwnAgents(companyID)
	for _, a := range own {
		if a.ID == agentID {
			if a.CashtroID != "" {
				target = a.CashtroID
			} else {
				if call.Capability == "" && len(a.Capabilities) > 0 {
					call.Capability = a.Capabilities[0]
				}
				writeJSON(w, http.StatusOK, kernel.Result{
					OK:      true,
					Message: a.Name + " resident · bind a Cashtro id to execute through Castro's kernel",
					Data: map[string]any{
						"agent": a.ID, "companyId": companyID, "bridged": false,
						"note": "No own worker yet. Bind cashtroId or leave the company on Castro's mine.",
					},
				})
				return
			}
			break
		}
	}
	if target == "" {
		for _, id := range fleet.MineIDs {
			if id == agentID {
				target = id
				break
			}
		}
	}
	if target == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "agentic not on this company"})
		return
	}
	if call.Capability == "" {
		if p, err := s.k.Process(target); err == nil && len(p.Spec.Capabilities) > 0 {
			call.Capability = p.Spec.Capabilities[0]
		}
	}
	res, err := s.k.Invoke(r.Context(), target, call)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *api) card(id string) companyCard {
	fl := s.tenants()
	c, err := fl.Get(id)
	if err != nil {
		return companyCard{}
	}
	useMine := fl.UsingMine(id)
	agents := make([]boundAgent, 0)
	own, _ := fl.OwnAgents(id)
	seen := map[string]bool{}
	for _, a := range own {
		mode := a.Mode
		if a.CashtroID != "" {
			if p, err := s.k.Process(a.CashtroID); err == nil {
				mode = string(p.Spec.Mode)
			}
		}
		agents = append(agents, boundAgent{
			ID: a.ID, Name: a.Name, Role: a.Role, Mode: mode, Summary: a.Summary,
			CashtroID: a.CashtroID, Capabilities: a.Capabilities, Inherited: false, Source: "own",
		})
		if a.CashtroID != "" {
			seen[a.CashtroID] = true
		}
		seen[a.ID] = true
	}
	if useMine {
		for _, p := range s.k.Processes() {
			if p.Spec.Kind == kernel.KindSystem && id != "cashtro" {
				continue
			}
			if seen[p.Spec.ID] {
				continue
			}
			agents = append(agents, boundAgent{
				ID: p.Spec.ID, Name: p.Spec.Name, Role: p.Spec.Role, Mode: string(p.Spec.Mode),
				Summary: p.Spec.Summary, CashtroID: p.Spec.ID, Capabilities: p.Spec.Capabilities,
				Inherited: true, Source: "castro",
			})
		}
	}
	return companyCard{Company: c, UsingMine: useMine, Agents: agents}
}

func writeFleetError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, fleet.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, fleet.ErrExists):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}
