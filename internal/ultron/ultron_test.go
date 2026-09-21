package ultron

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testPlane(t *testing.T) *Plane {
	t.Helper()
	dir := t.TempDir()
	p := New(
		WithPersistPath(filepath.Join(dir, "ultron.json")),
		WithCashtroURL("http://127.0.0.1:9"), // intentionally down for unit tests
		WithBootOwnerPassword("ultron-change-me"),
		WithClock(func() time.Time { return time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC) }),
	)
	if err := p.Boot(); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestBootRBACAndCompanies(t *testing.T) {
	p := testPlane(t)
	about := p.About(false)
	if about.Companies < 6 || about.Agents < 8 || about.Users != 1 {
		t.Fatalf("about = %+v", about)
	}
	sess, user, err := p.Login("alejandro@proximityagency.ca", "ultron-change-me")
	if err != nil || user.Role != RoleOwner || sess.Token == "" {
		t.Fatalf("login: %+v %v", user, err)
	}
	if !Can(user.Role, PermUsersWrite) {
		t.Fatal("owner should manage users")
	}
	companies := p.ListCompanies(user)
	found := map[string]bool{}
	for _, c := range companies {
		found[c.ID] = true
	}
	for _, id := range []string{"giant", "scanapp", "proximity", "empire", "cashtro"} {
		if !found[id] {
			t.Fatalf("missing company %s", id)
		}
	}

	client, err := p.CreateUser(user, "client@giant.test", "Giant Client", RoleClient, []string{"giant"}, "client-pass-1")
	if err != nil {
		t.Fatal(err)
	}
	cs, _, err := p.Login("client@giant.test", "client-pass-1")
	if err != nil {
		t.Fatal(err)
	}
	cu, ok := p.UserForToken(cs.Token)
	if !ok || cu.ID != client.ID {
		t.Fatal("client session")
	}
	visible := p.ListCompanies(cu)
	if len(visible) != 1 || visible[0].ID != "giant" {
		t.Fatalf("client tenancy = %+v", visible)
	}
	if AllowsCompany(cu, "proximity") {
		t.Fatal("client must not see proximity")
	}
}

func TestHTTPLoginAndFleet(t *testing.T) {
	p := testPlane(t)
	h := Handler(p)

	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/health", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"service":"ultron"`) {
		t.Fatalf("health = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"alejandro@proximityagency.ca","password":"ultron-change-me"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("login = %d %s", res.Code, res.Body.String())
	}
	var login struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &login); err != nil || login.Token == "" {
		t.Fatalf("token parse: %v %s", err, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/companies", nil)
	req.Header.Set("Authorization", "Bearer "+login.Token)
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "Giant") {
		t.Fatalf("companies = %s", res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/fleet/agents/empire-planner/invoke", strings.NewReader(`{"goal":"expand"}`))
	req.Header.Set("Authorization", "Bearer "+login.Token)
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	// Cashtro down → bridge error is acceptable; local agents without cashtro still work.
	// empire-planner has CashtroID planner — expect bad gateway or error body.
	if res.Code != http.StatusOK && res.Code != http.StatusInternalServerError && res.Code != http.StatusBadGateway {
		t.Fatalf("invoke = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/fleet/agents/scan-ingest/invoke", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer "+login.Token)
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "dry-run") {
		t.Fatalf("local invoke = %d %s", res.Code, res.Body.String())
	}

	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "ULTRON") {
		t.Fatalf("ide = %d", res.Code)
	}
}

func TestPersistRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ultron.json")
	p := New(WithPersistPath(path), WithBootOwnerPassword("ultron-change-me"))
	if err := p.Boot(); err != nil {
		t.Fatal(err)
	}
	if err := SaveFile(path, p); err != nil {
		t.Fatal(err)
	}
	p2 := New(WithPersistPath(path))
	if err := LoadFile(path, p2); err != nil {
		t.Fatal(err)
	}
	if len(p2.ListUsers()) != 1 {
		t.Fatalf("users = %d", len(p2.ListUsers()))
	}
	_, _, err := p2.Login("alejandro@proximityagency.ca", "ultron-change-me")
	if err != nil {
		t.Fatal(err)
	}
}
