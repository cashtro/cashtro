package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/agents"
	"github.com/cashtro/cashtro/internal/catalog"
	"github.com/cashtro/cashtro/internal/kernel"
)

func handler(t *testing.T) http.Handler {
	t.Helper()
	k, err := agents.Boot()
	if err != nil {
		t.Fatal(err)
	}
	return New(k)
}

func TestDeskHasStartUsingLink(t *testing.T) {
	h := handler(t)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("index status = %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, `id="use-link"`) || !strings.Contains(body, "Start using this desk") {
		t.Fatalf("desk is missing the start-using link: %s", body[:min(len(body), 400)])
	}
}

func TestHealthAndProfile(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/health", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("health status = %d", res.Code)
	}
	if !strings.Contains(res.Body.String(), `"os":"Cashtro OS"`) {
		t.Fatalf("health body = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/profile", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("profile status = %d", res.Code)
	}
	var profile catalog.Profile
	if err := json.Unmarshal(res.Body.Bytes(), &profile); err != nil {
		t.Fatal(err)
	}
	if profile.Name != "Castro" {
		t.Fatalf("profile name = %q", profile.Name)
	}
}

func TestOSAndAgents(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/os", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("os status = %d", res.Code)
	}
	var about kernel.About
	if err := json.Unmarshal(res.Body.Bytes(), &about); err != nil {
		t.Fatal(err)
	}
	if about.Agents != 15 || !strings.Contains(about.Manifesto, "under construction") {
		t.Fatalf("about = %+v", about)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/agents", nil))
	var procs []kernel.Process
	if err := json.Unmarshal(res.Body.Bytes(), &procs); err != nil {
		t.Fatal(err)
	}
	if len(procs) != 15 {
		t.Fatalf("agents = %d", len(procs))
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/notes", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "2606.01508") {
		t.Fatalf("notes = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/model", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"bound":false`) {
		t.Fatalf("model = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/agents/explorer/invoke", strings.NewReader(`{"capability":"explorer.search"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("invoke status = %d body=%s", res.Code, res.Body.String())
	}
}

func TestFavicon(t *testing.T) {
	h := handler(t)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/favicon.ico", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("favicon status = %d", res.Code)
	}
	if !strings.Contains(res.Header().Get("Content-Type"), "image/svg+xml") {
		t.Fatalf("favicon type = %q", res.Header().Get("Content-Type"))
	}
}

func TestIndexHTML(t *testing.T) {
	h := handler(t)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("index status = %d", res.Code)
	}
	ct := res.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("content-type = %q", ct)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Cashtro OS") || !strings.Contains(body, "idea → concept") {
		t.Fatalf("index missing OS shell copy")
	}
	if !strings.Contains(body, "GIANT") || !strings.Contains(body, "Ecosystem") || !strings.Contains(body, "Giant ecosystem") {
		t.Fatalf("index missing Giant ecosystem shell")
	}
	if !strings.Contains(body, "Add a company") || !strings.Contains(body, "Castro") {
		t.Fatalf("index missing evolvable fleet forms")
	}
	if !strings.Contains(body, "Close desk") || !strings.Contains(body, "Night shift now") {
		t.Fatalf("index missing closed-hours controls")
	}
	if !strings.Contains(body, `data-view="night"`) || !strings.Contains(body, "Overnight builds") {
		t.Fatalf("index missing night-shift view")
	}
}

func TestShipLifecycleHTTP(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/ships/missing", nil))
	if res.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d", res.Code)
	}

	body := `{"name":"North desk","client":"BTK Avocats","sector":"law","stack":["Go"]}`
	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ships", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", res.Code, res.Body.String())
	}
	var created catalog.Ship
	if err := json.Unmarshal(res.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID != "north-desk" || created.Stage != catalog.StageIdea {
		t.Fatalf("created = %+v", created)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/ships/north-desk/advance", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("advance status = %d", res.Code)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/ships/north-desk", nil))
	var got catalog.Ship
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Stage != catalog.StageConcept {
		t.Fatalf("stage = %s", got.Stage)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/ships", strings.NewReader(`{"name":""}`)))
	if res.Code != http.StatusBadRequest {
		t.Fatalf("empty name status = %d", res.Code)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/ships/proximity/advance", nil))
	if res.Code != http.StatusConflict {
		t.Fatalf("production advance status = %d", res.Code)
	}
}

func TestCreateRejectsUnknownField(t *testing.T) {
	h := handler(t)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ships", strings.NewReader(`{"name":"X","nope":true}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", res.Code, res.Body.String())
	}
}

func TestCompanyFleetHTTP(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/companies", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"Giant"`) {
		t.Fatalf("companies = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/companies/giant", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "Giant Conductor") {
		t.Fatalf("giant card = %s", res.Body.String())
	}
	if strings.Contains(res.Body.String(), `"inherited":true`) {
		t.Fatal("giant should not inherit Castro's while it has its own roster")
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/companies", strings.NewReader(`{"name":"Northwind","sector":"ops","notes":"new tenant"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("create company = %d %s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"usingMine":true`) || !strings.Contains(res.Body.String(), "Delivery") {
		t.Fatalf("empty company should use Castro's: %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/companies/northwind/agents", strings.NewReader(`{"name":"Northwind Clerk","role":"ops","cashtroId":"delivery","capabilities":["delivery.list"]}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusCreated || !strings.Contains(res.Body.String(), "northwind-clerk") {
		t.Fatalf("add agentic = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/companies/northwind/agents/northwind-clerk/invoke", strings.NewReader(`{"capability":"delivery.list"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"ok":true`) {
		t.Fatalf("invoke own bind = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/companies/northwind/agents/planner/invoke", strings.NewReader(`{"capability":"planner.backlog","payload":{"goal":"desk"}}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("invoke Castro mine = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/companies/scanapp/agents/scan-ingest/invoke", strings.NewReader(`{"capability":"scan.ingest"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "resident") {
		t.Fatalf("unbound own agentic = %d %s", res.Code, res.Body.String())
	}
}

func TestClosedHoursHTTP(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/watch", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"closed":false`) {
		t.Fatalf("watch open = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/watch/close", strings.NewReader(`{"note":"things are closed"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"closed":true`) {
		t.Fatalf("close = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/health", nil))
	if !strings.Contains(res.Body.String(), `"closed":true`) {
		t.Fatalf("health closed = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/watch/build", strings.NewReader(`{"note":"keep building"}`)))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"steps"`) {
		t.Fatalf("build = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/watch/open", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"closed":false`) {
		t.Fatalf("open = %s", res.Body.String())
	}
}

func TestStages(t *testing.T) {
	h := handler(t)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/stages", nil))
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), `"idea"`) || !strings.Contains(string(body), `"production"`) {
		t.Fatalf("stages = %s", body)
	}
}
