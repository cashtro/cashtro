package agents

import (
	"context"
	"strings"
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
	book := OpenCongress()
	if len(book.Laws) != len(Loi().Laws) || len(book.Bills) < 2 {
		t.Fatalf("congress laws %d bills %d", len(book.Laws), len(book.Bills))
	}
	ok, msg := book.Enforce()
	if !ok || msg == "" {
		t.Fatal("laws already written must be enforced")
	}
	if ok, _ := FusionEnforce("attaque", nil); ok {
		t.Fatal("an attack must stay closed")
	}
	next, bill, filed := book.Introduce("Faut-il une loi pour le brouillon local?")
	if !filed || bill.Status != "introduced" || len(next.Laws) != len(book.Laws) {
		t.Fatalf("bill %+v", bill)
	}
	if _, _, passed := next.Pass(bill.ID, " "); passed {
		t.Fatal("an empty rule must stay a question")
	}
	if _, _, passed := next.Pass(bill.ID, "fraude"); passed {
		t.Fatal("a harmful rule must stay a question")
	}
	passed, law, okPass := next.Pass(bill.ID, "Un brouillon local reste local tant qu'instinct n'a pas décidé.")
	if !okPass || law.Rules[0] == "" || len(passed.Laws) != len(book.Laws)+1 {
		t.Fatalf("pass %+v", law)
	}
	if bureau := FusionBureau("loi-25"); len(bureau) < 10 {
		t.Fatalf("bureau %d", len(bureau))
	}

	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.fusion"})
	if err != nil || !res.OK {
		t.Fatalf("existing laws: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.fusion",
		Payload:    []byte(`{"question":"Faut-il une loi pour le brouillon local?"}`),
	})
	if err != nil || !res.OK || !strings.Contains(res.Message, "pas encore une loi") {
		t.Fatalf("bill: %+v %v", res, err)
	}
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.fusion",
		Payload:    []byte(`{"bill":"bill-3","rule":"Un brouillon local reste local tant qu'instinct n'a pas décidé."}`),
	})
	if err != nil || !res.OK || !strings.Contains(res.Message, "loi passée") {
		t.Fatalf("pass: %+v %v", res, err)
	}
	if len(k.Recall("congress")) == 0 || len(k.Recall("fusion")) == 0 {
		t.Fatal("the session did not grow memory")
	}
}
