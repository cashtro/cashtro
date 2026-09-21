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
	t.Setenv("CASHTRO_CLOSED", "")
	k, err := agents.Boot()
	if err != nil {
		t.Fatal(err)
	}
	return New(k)
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

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/agents/architect/invoke", strings.NewReader(`{"capability":"architect.plan","payload":{"goal":"workflow pipeline for research notes"}}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"shape":"workflow"`) {
		t.Fatalf("architect = %d %s", res.Code, res.Body.String())
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
	if !strings.Contains(body, "Cashtro OS") || !strings.Contains(body, "idea → concept") || !strings.Contains(body, "Closed hours") {
		t.Fatalf("index missing OS shell copy")
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

func TestStages(t *testing.T) {
	h := handler(t)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/stages", nil))
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), `"idea"`) || !strings.Contains(string(body), `"production"`) {
		t.Fatalf("stages = %s", body)
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
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/ships", nil))
	if !strings.Contains(res.Body.String(), "closed-hours-flow") {
		t.Fatalf("ships missing closed-hours-flow: %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/watch/pulse", strings.NewReader(`{"note":"still here"}`)))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "still here") {
		t.Fatalf("pulse = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/watch/open", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"closed":false`) {
		t.Fatalf("open = %d %s", res.Code, res.Body.String())
	}
}
