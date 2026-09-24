package agents

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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
	if len(book.Laws) != len(Loi().Laws)+2 || len(book.Passed) != 2 {
		t.Fatalf("congress laws %d passed %d", len(book.Laws), len(book.Passed))
	}
	if !strings.Contains(book.Passed[0].Rules[0], "The Eye") || !strings.Contains(book.Passed[1].Rules[0], "sanction nommée") {
		t.Fatalf("passed %+v", book.Passed)
	}
	houses := Houses()
	if len(houses) != 2 || houses[0].ID != "cia" || !houses[0].Gathers || houses[0].GivesTo[0] != "instinct" || houses[0].GivesTo[1] != "fbi" {
		t.Fatalf("cia %+v", houses[0])
	}
	if houses[1].ID != "fbi" || houses[1].Gathers || len(houses[0].Departments) != len(Chart().Departments) || len(houses[1].Departments) != len(houses[0].Departments) {
		t.Fatalf("fbi %+v", houses[1])
	}
	if houses[0].Opposition.Partner != "fbi-opposition" || houses[1].Opposition.Partner != "cia-opposition" || !houses[0].Opposition.BuildsTogether || !houses[1].Opposition.ScrutinizesOwn {
		t.Fatalf("opposition %+v %+v", houses[0].Opposition, houses[1].Opposition)
	}
	if _, _, okStory := book.Sanction("", "loi-25", "fait", "histoire"); okStory {
		t.Fatal("a sanction needs a name")
	}
	book, told, okStory := book.Sanction("brouillon-sorti", "loi-25", "un build a quitté le local", "Le geste est arrêté. L'histoire s'appelle brouillon-sorti.")
	if !okStory || !told.Stopped || !told.Traced || !told.Saved || told.Understood || len(book.Stories) != 1 {
		t.Fatalf("story %+v", told)
	}
	book, told, okStory = book.Train(told.Name, "", nil)
	if !okStory || told.Understood || told.Rounds != 1 {
		t.Fatalf("open training %+v", told)
	}
	book, told, okStory = book.Train(told.Name, "appliquer la loi", []string{"Quelle loi s'applique?", "Quel geste s'arrête?"})
	if !okStory || !told.Understood || told.Retrain || told.Rounds != 2 {
		t.Fatalf("trained %+v", told)
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
	body := FusionBody()
	if len(body.RunBy) != 2 || body.RunBy[0] != "instinct" || body.RunBy[1] != "voltron" || !body.Hustler || !body.Person || body.Chair != "hustler" || body.Decision != "hustler" {
		t.Fatalf("runners %+v", body)
	}
	if body.Smell || body.Eat {
		t.Fatal("this body does not smell and does not eat")
	}
	if body.Brain != "instinct" || body.Spy != "eye" || body.Enforcer != "fusion" || body.Cyber != "security" || body.CyberAttacks || !body.Connected {
		t.Fatalf("body %+v", body)
	}
	if !body.Graphic.Growing || !body.Graphic.Knows || len(body.Graphic.Files) != 3 {
		t.Fatalf("graphic %+v", body.Graphic)
	}
	for _, organ := range body.Organs {
		if organ.ID == "smell" {
			t.Fatal("smell is not an organ")
		}
	}
	if len(body.Relay) != 2 || body.Relay[0] != "hustler" || body.Relay[1] != "instinct" || body.RelayTo != "Einstein" || !body.Vision || body.Profile != "mogul, 50 Cent" || !body.Ground || !body.Contacts || !body.Swarm {
		t.Fatalf("relay %+v", body.Relay)
	}
	if len(body.Eyes) != len(body.Members) || len(body.Organs) != 11 {
		t.Fatalf("eyes %d members %d organs %d", len(body.Eyes), len(body.Members), len(body.Organs))
	}
	resume := FusionResume()
	if len(body.Veins) < 8 {
		t.Fatalf("veins %d", len(body.Veins))
	}
	for _, vein := range body.Veins {
		if vein.Copied || vein.Lesson == "" {
			t.Fatalf("vein copied %+v", vein)
		}
	}
	vault := OpenVault()
	if vault.Name != "The Vault" || !vault.Heart || !vault.Intact || len(vault.Brains) != len(Brains()) {
		t.Fatalf("vault %+v", vault.Name)
	}
	if len(vault.Section("core")) < 8 || len(vault.Section("money")) != 1 || len(vault.Section("keys")) != 2 {
		t.Fatalf("sections core %d money %d keys %d", len(vault.Section("core")), len(vault.Section("money")), len(vault.Section("keys")))
	}
	for _, entry := range vault.Section("keys") {
		if entry.Body != "" {
			t.Fatal("a key stored a value")
		}
	}
	if _, _, ok := vault.Add("keys", "slot", "stripe", "sk_test_"+"SHOULDNOT"); ok {
		t.Fatal("a secret must stay out of the vault")
	}
	for _, line := range []string{"The hustler is the chair", "relays everything to Einstein", "50 Cent", "Robert Greene", "The 50th Law", "The Vault is the heart", "one swarm chain", "CIA: The Eye", "FBI: The Fusion", "retrains the agent", "does not smell and does not eat", "docs/MAP.md"} {
		if !strings.Contains(resume, line) {
			t.Fatalf("resume missing %s", line)
		}
	}
	root, err := findFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(root), "docs", "FUSION.md"), []byte(resume), 0o644); err != nil {
		t.Fatal(err)
	}
	raw, err := json.MarshalIndent(OpenVault(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(root), "state", "vault.json"), append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
