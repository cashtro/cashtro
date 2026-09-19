package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/cashtro/cashtro/internal/catalog"
)

const maxBody = 1 << 20

// New returns the Cashtro HTTP handler.
func New(cat *catalog.Catalog) http.Handler {
	s := &api{cat: cat}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/profile", s.profile)
	mux.HandleFunc("GET /api/stages", s.stages)
	mux.HandleFunc("GET /api/ships", s.listShips)
	mux.HandleFunc("POST /api/ships", s.createShip)
	mux.HandleFunc("GET /api/ships/{id}", s.getShip)
	mux.HandleFunc("POST /api/ships/{id}/advance", s.advanceShip)
	return mux
}

type api struct {
	cat *catalog.Catalog
}

func (s *api) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}

func (s *api) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "cashtro",
	})
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
	ship, err := s.cat.Create(in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, ship)
}

func (s *api) advanceShip(w http.ResponseWriter, r *http.Request) {
	ship, err := s.cat.Advance(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ship)
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
	case errors.Is(err, catalog.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, catalog.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, catalog.ErrDone):
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
