package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Evolu-Jeunes/Teal/internal/catalog"
	"github.com/Evolu-Jeunes/Teal/internal/kernel"
)

const maxBody = 1 << 20

// New returns the Cashtro OS HTTP shell.
func New(k *kernel.Kernel) http.Handler {
	s := &api{k: k, cat: k.Catalog()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /favicon.ico", s.favicon)
	mux.HandleFunc("GET /favicon.svg", s.favicon)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/os", s.osAbout)
	mux.HandleFunc("GET /api/agents", s.listAgents)
	mux.HandleFunc("GET /api/agents/{id}", s.getAgent)
	mux.HandleFunc("POST /api/agents/{id}/spawn", s.spawnAgent)
	mux.HandleFunc("POST /api/agents/{id}/invoke", s.invokeAgent)
	mux.HandleFunc("GET /api/capabilities", s.capabilities)
	mux.HandleFunc("GET /api/events", s.events)
	mux.HandleFunc("GET /api/model", s.modelStatus)
	mux.HandleFunc("GET /api/notes", s.notes)
	mux.HandleFunc("POST /api/notes", s.writeNote)
	mux.HandleFunc("GET /api/mail", s.mail)
	mux.HandleFunc("GET /api/memory", s.memory)
	mux.HandleFunc("GET /api/confirms", s.confirms)
	mux.HandleFunc("POST /api/confirms/{id}/allow", s.allowConfirm)
	mux.HandleFunc("POST /api/confirms/{id}/deny", s.denyConfirm)
	mux.HandleFunc("GET /api/profile", s.profile)
	mux.HandleFunc("GET /api/stages", s.stages)
	mux.HandleFunc("GET /api/ships", s.listShips)
	mux.HandleFunc("POST /api/ships", s.createShip)
	mux.HandleFunc("GET /api/ships/{id}", s.getShip)
	mux.HandleFunc("POST /api/ships/{id}/advance", s.advanceShip)
	mux.HandleFunc("GET /api/teal", s.tealStatus)
	mux.HandleFunc("GET /api/vapi", s.vapiStatus)
	mux.HandleFunc("GET /api/vapi/voices", s.vapiInvoke("vapi.voices"))
	mux.HandleFunc("POST /api/vapi/voice", s.vapiInvoke("vapi.voice"))
	mux.HandleFunc("POST /api/vapi/talk", s.vapiInvoke("vapi.talk"))
	mux.HandleFunc("POST /api/vapi/call", s.vapiInvoke("vapi.call"))
	mux.HandleFunc("POST /api/vapi/fire", s.vapiInvoke("vapi.fire"))
	mux.HandleFunc("GET /api/vapi/web", s.vapiInvoke("vapi.web"))
	mux.HandleFunc("POST /api/vapi/bridge", s.vapiInvoke("vapi.bridge"))
	mux.HandleFunc("POST /api/vapi/webhook", s.vapiWebhook)
	mux.HandleFunc("POST /api/teal/scrape", s.tealInvoke("teal.scrape"))
	mux.HandleFunc("POST /api/teal/ingest", s.tealInvoke("teal.ingest"))
	mux.HandleFunc("POST /api/teal/idea", s.tealInvoke("teal.idea"))
	mux.HandleFunc("GET /api/teal/interns", s.tealInterns)
	mux.HandleFunc("POST /api/teal/interns", s.tealInvoke("teal.assign"))
	mux.HandleFunc("POST /api/teal/assign", s.tealInvoke("teal.assign"))
	mux.HandleFunc("POST /api/teal/question", s.tealInvoke("teal.question"))
	mux.HandleFunc("POST /api/teal/help", s.tealInvoke("teal.help"))
	mux.HandleFunc("POST /api/teal/listen", s.tealInvoke("teal.listen"))
	mux.HandleFunc("POST /api/teal/spawn", s.tealInvoke("teal.spawn"))
	mux.HandleFunc("POST /api/teal/smarter", s.tealInvoke("teal.smarter"))
	mux.HandleFunc("POST /api/teal/cursor", s.tealInvoke("teal.cursor"))
	mux.HandleFunc("GET /api/teal/card", s.tealCard)
	return mux
}

type api struct {
	k   *kernel.Kernel
	cat *catalog.Catalog
}

func (s *api) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}

func (s *api) favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(faviconSVG)
}

func (s *api) health(w http.ResponseWriter, r *http.Request) {
	about := s.k.About()
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "teal",
		"os":      about.Name,
		"version": about.Version,
		"kernel":  about.Kernel,
		"agents":  about.Agents,
		"running": about.Running,
	})
}

func (s *api) osAbout(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.About())
}

func (s *api) listAgents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Processes())
}

func (s *api) getAgent(w http.ResponseWriter, r *http.Request) {
	p, err := s.k.Process(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *api) spawnAgent(w http.ResponseWriter, r *http.Request) {
	if err := s.k.Spawn(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	p, err := s.k.Process(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *api) invokeAgent(w http.ResponseWriter, r *http.Request) {
	var call kernel.Call
	if err := decodeJSON(r, &call); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	res, err := s.k.Invoke(r.Context(), r.PathValue("id"), call)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *api) capabilities(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Capabilities())
}

func (s *api) events(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Events())
}

func (s *api) modelStatus(w http.ResponseWriter, r *http.Request) {
	res, err := s.k.Invoke(r.Context(), "router", kernel.Call{Capability: "model.status"})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func (s *api) notes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Notes())
}

func (s *api) writeNote(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	res, err := s.k.Invoke(r.Context(), "research", kernel.Call{Capability: "research.ingest", Payload: raw})
	if err != nil {
		writeError(w, err)
		return
	}
	if !res.OK {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": res.Message})
		return
	}
	writeJSON(w, http.StatusCreated, res.Data)
}

func (s *api) mail(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Inbox(r.URL.Query().Get("to")))
}

func (s *api) memory(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Recall(r.URL.Query().Get("q")))
}

func (s *api) confirms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.k.Confirms())
}

func (s *api) allowConfirm(w http.ResponseWriter, r *http.Request) {
	s.decideConfirm(w, r, true)
}

func (s *api) denyConfirm(w http.ResponseWriter, r *http.Request) {
	s.decideConfirm(w, r, false)
}

func (s *api) decideConfirm(w http.ResponseWriter, r *http.Request, allow bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad confirm id"})
		return
	}
	c, err := s.k.DecideConfirm(id, allow)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if allow && c.Agent == "vapi" && c.Cap == "vapi.call" {
		raw, _ := json.Marshal(map[string]int{"id": c.ID, "confirmId": c.ID})
		res, ferr := s.k.Invoke(r.Context(), "vapi", kernel.Call{Capability: "vapi.fire", Payload: raw})
		if ferr == nil {
			writeJSON(w, http.StatusOK, map[string]any{"confirm": c, "fire": res})
			return
		}
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *api) profile(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cat.Profile())
}

func (s *api) stages(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, catalog.Stages())
}

func (s *api) listShips(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cat.List())
}

func (s *api) getShip(w http.ResponseWriter, r *http.Request) {
	ship, err := s.cat.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ship)
}

func (s *api) createShip(w http.ResponseWriter, r *http.Request) {
	var in catalog.CreateShip
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	raw, _ := json.Marshal(in)
	res, err := s.k.Invoke(r.Context(), "delivery", kernel.Call{Capability: "delivery.create", Payload: raw})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res.Data)
}

func (s *api) tealStatus(w http.ResponseWriter, r *http.Request) {
	res, err := s.k.Invoke(r.Context(), "teal", kernel.Call{Capability: "teal.status"})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func (s *api) vapiStatus(w http.ResponseWriter, r *http.Request) {
	res, err := s.k.Invoke(r.Context(), "vapi", kernel.Call{Capability: "vapi.status"})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func (s *api) vapiInvoke(cap string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if len(bytes.TrimSpace(raw)) == 0 {
			raw = []byte(`{}`)
		}
		res, err := s.k.Invoke(r.Context(), "vapi", kernel.Call{Capability: cap, Payload: raw})
		if err != nil {
			writeError(w, err)
			return
		}
		if !res.OK {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": res.Message})
			return
		}
		writeJSON(w, http.StatusOK, res.Data)
	}
}

func (s *api) tealCard(w http.ResponseWriter, r *http.Request) {
	res, err := s.k.Invoke(r.Context(), "teal", kernel.Call{Capability: "teal.card"})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func (s *api) tealInterns(w http.ResponseWriter, r *http.Request) {
	res, err := s.k.Invoke(r.Context(), "teal", kernel.Call{Capability: "teal.interns"})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func (s *api) vapiWebhook(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	res, err := s.k.Invoke(r.Context(), "teal", kernel.Call{Capability: "teal.listen", Payload: raw})
	if err != nil {
		writeError(w, err)
		return
	}
	if !res.OK {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": res.Message})
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func (s *api) tealInvoke(cap string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if len(bytes.TrimSpace(raw)) == 0 {
			raw = []byte(`{}`)
		}
		res, err := s.k.Invoke(r.Context(), "teal", kernel.Call{Capability: cap, Payload: raw})
		if err != nil {
			writeError(w, err)
			return
		}
		if !res.OK {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": res.Message})
			return
		}
		writeJSON(w, http.StatusOK, res.Data)
	}
}

func (s *api) advanceShip(w http.ResponseWriter, r *http.Request) {
	raw, _ := json.Marshal(map[string]string{"id": r.PathValue("id")})
	res, err := s.k.Invoke(r.Context(), "delivery", kernel.Call{Capability: "delivery.advance", Payload: raw})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res.Data)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBody))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, catalog.ErrNotFound), errors.Is(err, kernel.ErrUnknownAgent):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, catalog.ErrInvalid), errors.Is(err, kernel.ErrUnknownCapability):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, catalog.ErrDone), errors.Is(err, kernel.ErrNotRunning):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}
