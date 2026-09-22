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

func TestHealthAndProfile(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/health", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("health status = %d", res.Code)
	}
	if !strings.Contains(res.Body.String(), `"os":"Cashtro OS"`) || !strings.Contains(res.Body.String(), `"always":true`) {
		t.Fatalf("health body = %s", res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"uptimeSec"`) {
		t.Fatalf("health missing uptime: %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"always":true`) {
		t.Fatalf("api health = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/heartbeat", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"always":true`) {
		t.Fatalf("heartbeat = %d %s", res.Code, res.Body.String())
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
	if about.Agents != 14 || !strings.Contains(about.Manifesto, "under construction") {
		t.Fatalf("about = %+v", about)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/agents", nil))
	var procs []kernel.Process
	if err := json.Unmarshal(res.Body.Bytes(), &procs); err != nil {
		t.Fatal(err)
	}
	if len(procs) != 14 {
		t.Fatalf("agents = %d", len(procs))
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/trace?q=AOS", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "blast radius") {
		t.Fatalf("trace = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/security", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "security") {
		t.Fatalf("security = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/deploy", strings.NewReader(`{"target":"cashtro-os"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "dry-run") {
		t.Fatalf("deploy = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/review", strings.NewReader(`{"target":"cashtro-os"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "review") {
		t.Fatalf("review = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/browse", strings.NewReader(`{"url":"http://127.0.0.1:8080/"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "browse planned") {
		t.Fatalf("browse = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/plan", strings.NewReader(`{"goal":"Overnight Ultron company desk","items":["Giant pulse","Empire backlog"]}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "parked") {
		t.Fatalf("plan = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/architect", strings.NewReader(`{"goal":"Ultron company control plane","steps":["RBAC gate","Fleet bridge","Always-on"]}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "planned") {
		t.Fatalf("architect = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/memory", strings.NewReader(`{"topic":"overnight","text":"Cashtro remembers without a model key"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusCreated || !strings.Contains(res.Body.String(), "overnight") {
		t.Fatalf("memory store = %s", res.Body.String())
	}
	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/memory?q=overnight", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "remembers") {
		t.Fatalf("memory recall = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/search?q=delivery", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "explorer hit") {
		t.Fatalf("search = %s", res.Body.String())
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
	req = httptest.NewRequest(http.MethodPost, "/api/agents/explorer/invoke", strings.NewReader(`{"capability":"explorer.search"}`))
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
