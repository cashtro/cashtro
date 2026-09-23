package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestWorkersRunAndStayLocal(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	about := k.About()
	if about.Live != 14 || about.Resident != 1 || about.Agents != 15 {
		t.Fatalf("about = %+v", about)
	}
	if len(Chart().Departments) != 9 {
		t.Fatalf("departments = %d", len(Chart().Departments))
	}
	flows := Flows()
	if len(flows) != 11 || flows[0].Function == "" || flows[0].Product == "" {
		t.Fatalf("flows = %+v", flows)
	}

	secret := "sk_test_" + "SHOULDNOTAPPEAR"
	raw, _ := json.Marshal(map[string]string{"task": "colle " + secret})
	res, err := k.Invoke(context.Background(), "security", kernel.Call{Capability: "security.triage", Payload: raw})
	if err != nil || res.OK {
		t.Fatalf("secret should stop security: %+v %v", res, err)
	}
	blob, _ := json.Marshal(res)
	if strings.Contains(string(blob), "SHOULDNOTAPPEAR") {
		t.Fatal("security copied the secret")
	}

	raw, _ = json.Marshal(map[string]string{"repo": "Evolu-Jeunes/Alaska", "task": "éditer"})
	res, err = k.Invoke(context.Background(), "operator", kernel.Call{Capability: "operator.work", Payload: raw})
	if err != nil || res.OK {
		t.Fatalf("other department must refuse: %+v %v", res, err)
	}

	raw, _ = json.Marshal(map[string]string{"repo": "Evolu-Jeunes/immobilier", "task": "page accueil"})
	res, err = k.Invoke(context.Background(), "operator", kernel.Call{Capability: "operator.work", Payload: raw})
	if err != nil || !res.OK || !strings.Contains(res.Message, "local") {
		t.Fatalf("local work: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "deploy", kernel.Call{Capability: "deploy.release", Payload: raw})
	if err != nil || res.OK || !strings.Contains(res.Message, "rien n'est poussé") {
		t.Fatalf("deploy: %+v %v", res, err)
	}

	raw, _ = json.Marshal(map[string]string{"id": "wordpress", "repo": "Evolu-Jeunes/ProximityApp", "task": "contexte local"})
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.run", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("manager.run: %+v %v", res, err)
	}
	if !strings.Contains(res.Message, "rien n'est poussé") {
		t.Fatalf("cycle message = %s", res.Message)
	}
	res, err = k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.cycle", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("os.cycle: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "architect", kernel.Call{
		Capability: "architect.plan",
		Payload:    []byte(`{"id":"wordpress"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("architect: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "investigator", kernel.Call{
		Capability: "investigator.trace",
		Payload:    []byte(`{"repo":"Evolu-Jeunes/immobilier"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("investigator: %+v %v", res, err)
	}
	data := res.Data.(map[string]any)
	if data["line"] != "wordpress" || data["department"] != "wordpress" {
		t.Fatalf("trace = %+v", data)
	}

	for _, id := range []string{"architecte", "forge", "eye", "giant", "accueil", "propriete"} {
		raw, _ = json.Marshal(map[string]string{"id": id, "task": "lire le cerveau"})
		res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.run", Payload: raw})
		if err != nil || !res.OK {
			t.Fatalf("cycle %s: %+v %v", id, res, err)
		}
		if !strings.Contains(res.Message, "rien n'est poussé") {
			t.Fatalf("cycle %s message = %s", id, res.Message)
		}
	}
}
