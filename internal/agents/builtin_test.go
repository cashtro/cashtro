package agents

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBootLoadsAllAgentics(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	procs := k.Processes()
	if len(procs) != 15 {
		t.Fatalf("processes = %d, want 15", len(procs))
	}
	about := k.About()
	if about.Running != 15 || about.Live != 8 {
		t.Fatalf("about = %+v", about)
	}
	if len(k.Notes()) < 5 {
		t.Fatalf("research seeds = %d", len(k.Notes()))
	}
	if k.Catalog() == nil {
		t.Fatal("delivery did not attach catalog")
	}

	res, err := k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.about"})
	if err != nil || !res.OK {
		t.Fatalf("os.about: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "delivery", kernel.Call{Capability: "delivery.list"})
	if err != nil {
		t.Fatal(err)
	}
	ships, ok := res.Data.([]any)
	if !ok {
		// list returns []catalog.Ship
		if res.Data == nil {
			t.Fatal("empty list")
		}
	}
	_ = ships

	payload, _ := json.Marshal(map[string]string{"name": "North desk"})
	res, err = k.Invoke(context.Background(), "delivery", kernel.Call{
		Capability: "delivery.create",
		Payload:    payload,
	})
	if err != nil || !res.OK {
		t.Fatalf("create: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "router", kernel.Call{Capability: "model.status"})
	if err != nil || !res.OK {
		t.Fatalf("router status: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "router", kernel.Call{
		Capability: "model.chat",
		Payload:    []byte(`{"prompt":"ping"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("unbound chat should be ok=false: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"delivery"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("explorer: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "research", kernel.Call{Capability: "research.list"})
	if err != nil || !res.OK {
		t.Fatalf("research: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.status"})
	if err != nil || !res.OK {
		t.Fatalf("manager.status: %+v %v", res, err)
	}
	status := res.Data.(map[string]any)
	if status["version"] != "2" {
		t.Fatalf("board version = %v", status["version"])
	}
	if got := status["agents"]; got != 15 {
		t.Fatalf("board agents = %v, want 15", got)
	}
	if got := status["lines"]; got != 10 {
		t.Fatalf("board lines = %v, want 9", got)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.lines"})
	if err != nil || !res.OK {
		t.Fatalf("manager.lines: %+v %v", res, err)
	}
	lines, ok := res.Data.([]Line)
	if !ok || len(lines) != 10 {
		t.Fatalf("lines = %#v", res.Data)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.line",
		Payload:    []byte(`{"id":"proximity"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.line: %+v %v", res, err)
	}
	if res.Data.(Line).Brain != "Forgeron" {
		t.Fatalf("proximity brain = %+v", res.Data)
	}

	org := Chart()
	if len(org.Seats) != 3 || len(org.Departments) != 7 || len(org.Divisions) != 8 {
		t.Fatalf("org seats/depts/divs = %d %d %d", len(org.Seats), len(org.Departments), len(org.Divisions))
	}
	var scanWork int
	for _, d := range org.Divisions {
		if d.ID == "scanapp" {
			scanWork = len(d.Work)
		}
	}
	if scanWork < 8 {
		t.Fatalf("scanapp chain = %d steps", scanWork)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.chain",
		Payload:    []byte(`{"id":"scanapp"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("chain without context should ask: %+v %v", res, err)
	}
	if qs, ok := res.Data.([]Question); !ok || len(qs) != 1 || qs[0].ID != "catalogue" {
		t.Fatalf("questions = %#v", res.Data)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.chain",
		Payload:    []byte(`{"id":"scanapp","answers":{"catalogue":"à confirmer"}}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.chain: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.fiche",
		Payload:    []byte(`{"code":"123","name":"Noix","price":"4.00","stock":2,"photo":"database"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("database photo must not publish: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.fiche",
		Payload:    []byte(`{"code":"123","name":"Noix","price":"4.00","stock":2,"photo":"generated"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("fiche without buyer should ask: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.fiche",
		Payload:    []byte(`{"code":"123","name":"Noix","price":"4.00","stock":2,"photo":"generated","answers":{"acheteur":"particulier au Québec","droits":"générée par nous, pas une personne"}}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("generated photo should publish: %+v %v", res, err)
	}
	if len(org.Members()) != 15 {
		t.Fatalf("employees = %d, want 15", len(org.Members()))
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.contradict",
		Payload:    []byte(`{"subject":"site","answers":{"situation":"le thème Noix existe déjà","intouchable":"aucun site live","mieux":"moins de temps"},"options":[{"name":"refaire","cost":8,"steps":6},{"name":"réutiliser","cost":1,"steps":2}]}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.contradict: %+v %v", res, err)
	}
	if res.Data.(Verdict).Best != "réutiliser" {
		t.Fatalf("best = %+v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.org"})
	if err != nil || !res.OK {
		t.Fatalf("manager.org: %+v %v", res, err)
	}
	loi := Loi()
	if loi.Officer != "manager" || loi.Delegate != "security" {
		t.Fatalf("officer = %s delegate = %s", loi.Officer, loi.Delegate)
	}
	found := false
	for _, law := range loi.Laws {
		if law.ID == "loi-25" {
			found = true
		}
	}
	if !found || len(loi.Gates) < 8 {
		t.Fatalf("loi laws=%d gates=%d", len(loi.Laws), len(loi.Gates))
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.loi"})
	if err != nil || !res.OK {
		t.Fatalf("manager.loi: %+v %v", res, err)
	}

	if len(AgentIDs()) != 15 {
		t.Fatalf("agentics in ecosystems = %d, want 15 (%v)", len(AgentIDs()), AgentIDs())
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.ecosystems"})
	if err != nil || !res.OK {
		t.Fatalf("manager.ecosystems: %+v %v", res, err)
	}
	if links, ok := res.Data.([]Link); !ok || len(links) < 8 {
		t.Fatalf("links = %#v", res.Data)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.automate"})
	if err != nil || !res.OK {
		t.Fatalf("manager.automate: %+v %v", res, err)
	}
	gotAgents := res.Data.(map[string]any)["agents"].([]string)
	if len(gotAgents) != 15 {
		t.Fatalf("automate agents = %v", gotAgents)
	}
	if len(k.Inbox("operator")) == 0 || len(k.Inbox("security")) == 0 || len(k.Inbox("init")) == 0 {
		t.Fatal("ecosystem mail did not reach the agentics")
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.assign",
		Payload:    []byte(`{"brain":"Forgeron","repo":"Evolu-Jeunes/Btkavocat","answers":{"pourquoi":"Forgeron confectionne, on écarte le Hustler","live":"oui, preview seulement"}}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("manager.assign: %+v %v", res, err)
	}
}
