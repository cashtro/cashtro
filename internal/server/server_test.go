package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/catalog"
)

func TestHealthAndProfile(t *testing.T) {
	h := New(catalog.New())

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/health", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("health status = %d", res.Code)
	}
	if !strings.Contains(res.Body.String(), `"service":"cashtro"`) {
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

func TestIndexHTML(t *testing.T) {
	h := New(catalog.New())
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("index status = %d", res.Code)
	}
	ct := res.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("content-type = %q", ct)
	}
	if !strings.Contains(res.Body.String(), "idea → concept") {
		t.Fatalf("index missing board copy")
	}
}

func TestShipLifecycleHTTP(t *testing.T) {
	h := New(catalog.New())

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
	h := New(catalog.New())
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ships", strings.NewReader(`{"name":"X","nope":true}`))
	h.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", res.Code, res.Body.String())
	}
}

func TestStages(t *testing.T) {
	h := New(catalog.New())
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/stages", nil))
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), `"idea"`) || !strings.Contains(string(body), `"production"`) {
		t.Fatalf("stages = %s", body)
	}
}
