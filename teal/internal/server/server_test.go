package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Evolu-Jeunes/Teal/internal/agents"
	"github.com/Evolu-Jeunes/Teal/internal/catalog"
	"github.com/Evolu-Jeunes/Teal/internal/kernel"
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
	if !strings.Contains(res.Body.String(), `"os":"Teal"`) {
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
	if about.Agents != 17 || !strings.Contains(about.Manifesto, "Evolu-Jeunes") {
		t.Fatalf("about = %+v", about)
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/agents", nil))
	var procs []kernel.Process
	if err := json.Unmarshal(res.Body.Bytes(), &procs); err != nil {
		t.Fatal(err)
	}
	if len(procs) != 17 {
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
	if !strings.Contains(body, "Teal") || !strings.Contains(body, "Voice lane") {
		t.Fatalf("index missing OS shell copy")
	}
	if !strings.Contains(body, "TE") || !strings.Contains(body, "Crew") {
		t.Fatalf("index missing Teal desk")
	}
	if !strings.Contains(body, "Voice lane") || !strings.Contains(body, "/api/vapi/talk") {
		t.Fatalf("index missing Vapi lane")
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

func TestTealAndVapiHTTP(t *testing.T) {
	h := handler(t)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/vapi", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "vapi.ai") {
		t.Fatalf("vapi = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/teal", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "microsoft-teams") {
		t.Fatalf("teal = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/teal/ingest", strings.NewReader(`{"channel":"microsoft-teams","text":"stand-up on Pandora"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("ingest = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/vapi/webhook", strings.NewReader(`{"message":{"call":{"id":"hook1"},"analysis":{"summary":"Call about intern onboarding"}}}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "hook1") {
		t.Fatalf("webhook = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/teal/card", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "AdaptiveCard") {
		t.Fatalf("card = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/teal/assign", strings.NewReader(`{"id":"intern-code","help":"pair on Pandora landing"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "assigned") && !strings.Contains(res.Body.String(), "intern-code") {
		t.Fatalf("assign = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/teal/ingest", strings.NewReader(`{"channel":"slack","text":"x"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("slack ingest = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/vapi/voices", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "denise") {
		t.Fatalf("voices = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/vapi/voice", strings.NewReader(`{"id":"nova"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "Nova") {
		t.Fatalf("voice = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/vapi/bridge", strings.NewReader(`{"line":"trading"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "trading") {
		t.Fatalf("bridge = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/vapi/talk", strings.NewReader(`{"text":"bridge this to proximity"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "Local Teal/Vapi") && !strings.Contains(res.Body.String(), "Bridges") && !strings.Contains(res.Body.String(), "reply") {
		t.Fatalf("talk = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/vapi/call", strings.NewReader(`{"to":"+15555550199","prompt":"hi Castro"}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "confirm") {
		t.Fatalf("call = %d %s", res.Code, res.Body.String())
	}
	var parked map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &parked); err != nil {
		t.Fatal(err)
	}
	confirm, _ := parked["confirm"].(map[string]any)
	id := int(confirm["id"].(float64))

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/api/confirms/"+strconv.Itoa(id)+"/allow", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "queued") && !strings.Contains(res.Body.String(), "bind") && !strings.Contains(res.Body.String(), "fire") {
		t.Fatalf("allow fire = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/vapi/web", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "dashboard.vapi.ai") {
		t.Fatalf("web = %d %s", res.Code, res.Body.String())
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
