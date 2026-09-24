package agents

import (
	"context"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestTheFusionCoversBothPeoples(t *testing.T) {
	force := Fusion()
	if force.Name != "The Fusion" || force.Agency != "fbi" || !force.BothAlways {
		t.Fatalf("name %+v", force.Name)
	}
	if len(force.Peoples) != 2 || force.Peoples[0] != "Evolu-Jeunes" || force.Peoples[1] != "cashtro" {
		t.Fatalf("peoples %+v", force.Peoples)
	}
	if len(force.Accounts) != 2 || force.Accounts[0] != "Evolu-Jeunes" || force.Accounts[1] != "cashtro" {
		t.Fatalf("accounts %+v", force.Accounts)
	}
	if force.Parent != "cia" || len(force.Under) != 2 || force.Under[0] != "cia" || force.Under[1] != "instinct" {
		t.Fatalf("command %+v", force.Under)
	}
	if len(force.Verbs) != 4 || force.Verbs[0] != "enforce" || force.Verbs[1] != "punish" || force.Verbs[2] != "trace" || force.Verbs[3] != "secure" {
		t.Fatalf("verbs %+v", force.Verbs)
	}
	if force.Agents < 80 || len(force.Workers) != force.Agents {
		t.Fatalf("agents %d workers %d", force.Agents, len(force.Workers))
	}
	if force.Repo != "cashtro/cashtro" || force.Brain != "fusion" {
		t.Fatalf("seat %+v", force)
	}
	if force.ImportsOther || force.MarkdownEdges || force.PathToChains || !force.Parallel {
		t.Fatal("The Fusion must not invent an import into a TypeScript brain")
	}
	if force.Weapons || force.Attacks || force.Illegal || force.CopiesSecrets || force.Surveillance || force.Impersonates || force.Sent || force.Deployed {
		t.Fatal("The Fusion must not attack, watch a person, or impersonate a state")
	}
	if force.OpenedBy != "voltron" || force.DecidedBy != "instinct" || force.MainBrain != "instinct" || force.LoadsInstinct || !force.OwnBrain {
		t.Fatalf("brain %+v", force)
	}
	seen := map[string]bool{}
	for _, w := range force.Workers {
		if seen[w.ID] {
			t.Fatalf("duplicate %s", w.ID)
		}
		seen[w.ID] = true
		if len(w.Peoples) != 2 || w.Peoples[0] != "Evolu-Jeunes" || w.Peoples[1] != "cashtro" {
			t.Fatalf("officer peoples %+v", w)
		}
		if w.Duty == "" || w.Law == "" {
			t.Fatalf("officer %+v", w)
		}
	}
	if !FusionCovers("The Fusion") || FusionCovers("Evolu-Jeunes") {
		t.Fatal("the name must mean both GitHubs, and only that name")
	}
	if ok, msg := FusionEnforce("cashtro", nil); ok || msg == "" {
		t.Fatal("enforcement without an answer must ask")
	}
	answers := map[string]string{
		"fait":  "un build sans version française",
		"loi":   "charte",
		"suite": "arrêter le geste",
	}
	if ok, msg := FusionEnforce("cashtro", answers); !ok || msg == "" {
		t.Fatalf("answered enforcement: %v %s", ok, msg)
	}
	if ok, _ := FusionEnforce("attaque", answers); ok {
		t.Fatal("an attack must stay closed")
	}
	if bureau := FusionBureau("loi-25"); len(bureau) < 10 {
		t.Fatalf("bureau %d", len(bureau))
	}

	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.fusion"})
	if err != nil || res.OK {
		t.Fatalf("fusion must ask first: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.fusion",
		Payload:    []byte(`{"law":"charte","answers":{"fait":"un build sans version française","loi":"charte","suite":"arrêter le geste"}}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("fusion after answers: %+v %v", res, err)
	}
	if len(k.Recall("fusion")) == 0 {
		t.Fatal("the answer did not grow memory")
	}
}
