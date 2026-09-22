package ultron

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const maxBody = 1 << 20

// Handler returns the Ultron IDE HTTP shell.
func Handler(p *Plane) http.Handler {
	s := &api{p: p}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("GET /api/about", s.about)
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("POST /api/auth/logout", s.logout)
	mux.HandleFunc("POST /api/auth/password", s.changePassword)
	mux.HandleFunc("GET /api/me", s.me)
	mux.HandleFunc("GET /api/companies", s.companies)
	mux.HandleFunc("GET /api/companies/{id}", s.company)
	mux.HandleFunc("POST /api/companies", s.createCompany)
	mux.HandleFunc("GET /api/fleet/agents", s.agents)
	mux.HandleFunc("GET /api/fleet/workers", s.workers)
	mux.HandleFunc("POST /api/fleet/agents/{id}/invoke", s.invoke)
	mux.HandleFunc("GET /api/os/health", s.osHealth)
	mux.HandleFunc("GET /api/os/proxy", s.osProxyGET)
	mux.HandleFunc("POST /api/os/heartbeat", s.osHeartbeat)
	mux.HandleFunc("GET /api/users", s.users)
	mux.HandleFunc("POST /api/users", s.createUser)
	mux.HandleFunc("GET /api/audit", s.audit)
	return mux
}

type api struct{ p *Plane }

func (s *api) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(ideHTML)
}

func (s *api) health(w http.ResponseWriter, r *http.Request) {
	_, ok := s.p.PingCashtro()
	about := s.p.About(ok)
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"service":   "ultron",
		"version":   about.Version,
		"companies": about.Companies,
		"agents":    about.Agents,
		"workers":   about.Workers,
		"cashtroOk": about.CashtroOK,
		"cashtro":   about.CashtroURL,
		"always":    true,
		"data":      s.p.PersistPath(),
	})
}

func (s *api) about(w http.ResponseWriter, r *http.Request) {
	_, ok := s.p.PingCashtro()
	writeJSON(w, http.StatusOK, s.p.About(ok))
}

func (s *api) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	sess, user, err := s.p.Login(in.Email, in.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": sess.Token, "expiresAt": sess.ExpiresAt, "user": user})
}

func (s *api) logout(w http.ResponseWriter, r *http.Request) {
	s.p.Logout(bearer(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *api) changePassword(w http.ResponseWriter, r *http.Request) {
	u, ok := s.p.UserForToken(bearer(r))
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	var in struct {
		Current string `json:"current"`
		Next    string `json:"next"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.p.ChangePassword(u, in.Current, in.Next); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "password rotated · re-login"})
}

func (s *api) me(w http.ResponseWriter, r *http.Request) {
	u, ok := s.p.UserForToken(bearer(r))
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *api) companies(w http.ResponseWriter, r *http.Request) {
	u, ok, msg := s.p.Authorize(bearer(r), PermCompanyRead, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	writeJSON(w, http.StatusOK, s.p.ListCompanies(u))
}

func (s *api) company(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u, ok, msg := s.p.Authorize(bearer(r), PermCompanyRead, id)
	if !ok {
		writeAuth(w, msg)
		return
	}
	c, err := s.p.GetCompany(u, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"company": c,
		"agents":  s.p.ListAgents(u, id),
		"workers": s.p.ListWorkers(u, id),
	})
}

func (s *api) createCompany(w http.ResponseWriter, r *http.Request) {
	u, ok, msg := s.p.Authorize(bearer(r), PermCompanyWrite, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	var in Company
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c, err := s.p.UpsertCompany(u, in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (s *api) agents(w http.ResponseWriter, r *http.Request) {
	u, ok, msg := s.p.Authorize(bearer(r), PermFleetRead, r.URL.Query().Get("company"))
	if !ok {
		writeAuth(w, msg)
		return
	}
	writeJSON(w, http.StatusOK, s.p.ListAgents(u, r.URL.Query().Get("company")))
}

func (s *api) workers(w http.ResponseWriter, r *http.Request) {
	u, ok, msg := s.p.Authorize(bearer(r), PermFleetRead, r.URL.Query().Get("company"))
	if !ok {
		writeAuth(w, msg)
		return
	}
	writeJSON(w, http.StatusOK, s.p.ListWorkers(u, r.URL.Query().Get("company")))
}

func (s *api) invoke(w http.ResponseWriter, r *http.Request) {
	u, ok, msg := s.p.Authorize(bearer(r), PermFleetInvoke, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	var payload map[string]any
	_ = decodeJSON(r, &payload)
	if payload == nil {
		payload = map[string]any{}
	}
	res, err := s.p.InvokeAgent(u, r.PathValue("id"), payload)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *api) osHealth(w http.ResponseWriter, r *http.Request) {
	_, ok, msg := s.p.Authorize(bearer(r), PermOSBridge, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	data, up := s.p.PingCashtro()
	status := http.StatusOK
	if !up {
		status = http.StatusBadGateway
	}
	writeJSON(w, status, data)
}

func (s *api) osProxyGET(w http.ResponseWriter, r *http.Request) {
	_, ok, msg := s.p.Authorize(bearer(r), PermOSBridge, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/api/os"
	}
	status, raw, err := s.p.BridgeGET(path)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(raw)
}

func (s *api) osHeartbeat(w http.ResponseWriter, r *http.Request) {
	u, ok, msg := s.p.Authorize(bearer(r), PermOSBridge, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	status, raw, err := s.p.BridgePOST("/api/heartbeat", nil)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	s.p.Audit(u.Email, "os.heartbeat", "cashtro", status == 200, "")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(raw)
}

func (s *api) users(w http.ResponseWriter, r *http.Request) {
	_, ok, msg := s.p.Authorize(bearer(r), PermUsersRead, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	writeJSON(w, http.StatusOK, s.p.ListUsers())
}

func (s *api) createUser(w http.ResponseWriter, r *http.Request) {
	actor, ok, msg := s.p.Authorize(bearer(r), PermUsersWrite, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	var in struct {
		Email      string   `json:"email"`
		Name       string   `json:"name"`
		Role       Role     `json:"role"`
		CompanyIDs []string `json:"companyIds"`
		Password   string   `json:"password"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	u, err := s.p.CreateUser(actor, in.Email, in.Name, in.Role, in.CompanyIDs, in.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (s *api) audit(w http.ResponseWriter, r *http.Request) {
	_, ok, msg := s.p.Authorize(bearer(r), PermAuditRead, "")
	if !ok {
		writeAuth(w, msg)
		return
	}
	writeJSON(w, http.StatusOK, s.p.AuditLog())
}

func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	if t := r.Header.Get("X-Ultron-Token"); t != "" {
		return t
	}
	if c, err := r.Cookie("ultron_token"); err == nil {
		return c.Value
	}
	return ""
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBody))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeAuth(w http.ResponseWriter, msg string) {
	code := http.StatusUnauthorized
	if strings.Contains(msg, "forbidden") {
		code = http.StatusForbidden
	}
	writeJSON(w, code, map[string]string{"error": msg})
}

func writeErr(w http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case ErrForbidden:
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
	case ErrExists:
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case ErrInvalid, ErrBadLogin:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
