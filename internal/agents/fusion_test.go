package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestTheFusionMeansBothPeoples(t *testing.T) {
	if got := Fusion(); len(got) != 2 || got[0] != "cashtro" || got[1] != "Evolu-Jeunes" {
		t.Fatalf("fusion = %v", got)
	}
	peoples := Peoples()
	if len(peoples) != 2 || peoples[0].Name != "Astro" || peoples[0].Account != "cashtro" || peoples[1].Name != "Evolian" || peoples[1].Account != "Evolu-Jeunes" {
		t.Fatalf("peoples = %+v", peoples)
	}
	if p, ok := PeopleOf("cashtro/cashtro"); !ok || p.Name != "Astro" {
		t.Fatalf("astro = %+v %v", p, ok)
	}
	if p, ok := PeopleOf("Evolu-Jeunes/OPS"); !ok || p.Name != "Evolian" {
		t.Fatalf("evolian = %+v %v", p, ok)
	}
	if InFusion("someone-else/repo") || InFusion("cashtroX/repo") || InFusion("") {
		t.Fatal("a repo outside the two GitHubs is not in the Fusion")
	}
	scope := FusionScope()
	if len(scope) < 50 {
		t.Fatalf("scope = %d", len(scope))
	}
	seen := map[string]bool{}
	for _, repo := range scope {
		if seen[repo] || !InFusion(repo) {
			t.Fatalf("scope repo %s", repo)
		}
		seen[repo] = true
	}
	if !seen["cashtro/cashtro"] || !seen["Evolu-Jeunes/Forge"] || !seen["Evolu-Jeunes/immobilier"] {
		t.Fatalf("scope misses a people: %v", scope)
	}
}

func TestTheFusionIsThePoliceInsideTheEye(t *testing.T) {
	force := TheFusion()
	if force.Name != "The Fusion" || force.Kind != "police" || force.Inside != "eye" {
		t.Fatalf("force %+v", force)
	}
	if force.Agents < 80 || force.Agents < force.Minimum || force.Minimum != 80 {
		t.Fatalf("effectif %d, minimum %d", force.Agents, force.Minimum)
	}
	if len(force.Squads) != len(Loi().Laws)+1 {
		t.Fatalf("squads %d laws %d", len(force.Squads), len(Loi().Laws))
	}
	sum := 0
	ids := map[string]bool{}
	for _, squad := range force.Squads {
		if squad.Officers < 1 || squad.Cites == "" || len(squad.Watches) == 0 || ids[squad.ID] {
			t.Fatalf("squad %+v", squad)
		}
		ids[squad.ID] = true
		sum += squad.Officers
	}
	if sum != force.Agents {
		t.Fatalf("sum %d agents %d", sum, force.Agents)
	}
	for _, law := range Loi().Laws {
		if !ids[law.ID] {
			t.Fatalf("no squad for %s", law.ID)
		}
	}
	if !ids["garde"] {
		t.Fatal("the guard holds the kernel's own rules")
	}
	if !force.Cites || !force.Stops || force.Judges || force.Attacks || force.Deletes || force.Pushes || force.CopiesSecrets || force.Weapons {
		t.Fatal("the police cites and stops. it does not judge, attack, delete, push, or copy")
	}
	if force.DecidedBy != "instinct" || force.ReportsTo != "eye" || force.OpenedBy != "voltron" || force.OrderedBy != "scrum" {
		t.Fatalf("seats %+v", force)
	}
	if len(force.Accounts) != 2 || force.Scope < 50 {
		t.Fatalf("both peoples %+v", force)
	}
	if ok, why := FusionHolds(); !ok {
		t.Fatalf("force does not hold: %s", why)
	}
	eye := TheEye()
	if eye.Police.Name != "The Fusion" || eye.Police.Agents < 80 {
		t.Fatalf("the eye does not carry the police: %+v", eye.Police)
	}
	brain, ok := BrainByID("fusion")
	if !ok || brain.Seat != "eye" || brain.Dept != "controle" {
		t.Fatalf("brain %+v %v", brain, ok)
	}
	if a := Analyze("fusion"); !a.Connected {
		t.Fatalf("fusion is not analyzed: %+v", a)
	}
}

func TestTheFusionEnforcesAcrossBothPeoples(t *testing.T) {
	clean := Enforce("Evolu-Jeunes/immobilier", "page accueil en local")
	if clean.Stopped || !clean.Jurisdiction || clean.People != "Evolian" || len(clean.Citations) != 0 || clean.Officers < 80 {
		t.Fatalf("clean move %+v", clean)
	}
	astro := Enforce("cashtro/PBTM", "lire le cerveau")
	if astro.Stopped || astro.People != "Astro" {
		t.Fatalf("astro move %+v", astro)
	}
	kernelMove := Enforce("", "travail local")
	if kernelMove.Stopped || !kernelMove.Jurisdiction || kernelMove.People != "kernel" {
		t.Fatalf("kernel move %+v", kernelMove)
	}

	outside := Enforce("autre-org/site", "éditer")
	if !outside.Stopped || outside.Jurisdiction || len(outside.Citations) != 1 || outside.Citations[0].Squad != "garde" {
		t.Fatalf("outside %+v", outside)
	}

	secret := "sk_live_" + "SHOULDNOTAPPEAR"
	leak := Enforce("Evolu-Jeunes/CRM", "colle "+secret+" dans le thème")
	if !leak.Stopped || leak.Applied || leak.HandedTo != "instinct" {
		t.Fatalf("leak %+v", leak)
	}
	blob, _ := json.Marshal(leak)
	if strings.Contains(string(blob), "SHOULDNOTAPPEAR") {
		t.Fatal("the fusion copied the secret")
	}

	cases := map[string]string{
		"casl":   "envoi massif de courriels à une liste achetée",
		"charte": "mettre le site en english only",
		"lpc":    "annoncer un gain garanti sur la page",
		"image":  "mettre la photo de la base sur la fiche",
		"racj":   "lancer une loterie sur Facebook",
		"amf":    "passer un ordre live sur le token",
		"loi-25": "envoyer la liste de clients à Azure sans EFVP",
		"garde":  "préparer une attaque sur le site du concurrent",
	}
	for squad, task := range cases {
		patrol := Enforce("Evolu-Jeunes/Fix2", task)
		if !patrol.Stopped {
			t.Fatalf("%s should stop: %+v", squad, patrol)
		}
		found := false
		for _, c := range patrol.Citations {
			if c.Squad == squad && c.Rule != "" && c.Found != "" {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s not cited: %+v", squad, patrol.Citations)
		}
	}
	illegal := Enforce("cashtro/cashtro", "monter un blanchiment")
	if !illegal.Stopped || illegal.Citations[len(illegal.Citations)-1].Found != "chemin illégal" {
		t.Fatalf("illegal %+v", illegal)
	}
}

func TestTheFusionPatrolsEveryCycle(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.fusion"})
	if err != nil || !res.OK {
		t.Fatalf("roster: %+v %v", res, err)
	}
	force := res.Data.(FusionForce)
	if force.Agents < 80 || len(force.Peoples) != 2 {
		t.Fatalf("roster = %+v", force)
	}

	raw, _ := json.Marshal(map[string]string{"repo": "Evolu-Jeunes/immobilier", "task": "page accueil"})
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.fusion", Payload: raw})
	if err != nil || !res.OK || !strings.Contains(res.Message, "Evolian") {
		t.Fatalf("patrol: %+v %v", res, err)
	}

	raw, _ = json.Marshal(map[string]string{"repo": "Evolu-Jeunes/immobilier", "task": "envoi massif à une liste achetée"})
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.fusion", Payload: raw})
	if err != nil || res.OK || !strings.Contains(res.Message, "casl") {
		t.Fatalf("patrol stop: %+v %v", res, err)
	}
	if res.Data.(Patrol).Applied {
		t.Fatal("the fusion applied something")
	}

	raw, _ = json.Marshal(map[string]string{"id": "wordpress", "repo": "autre-org/site", "task": "contexte local"})
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.run", Payload: raw})
	if err != nil || res.OK || !strings.Contains(res.Message, "hors juridiction") {
		t.Fatalf("cycle outside the fusion must stop: %+v %v", res, err)
	}
	raw, _ = json.Marshal(map[string]string{"id": "wordpress", "repo": "Evolu-Jeunes/ProximityApp", "task": "prix caché sur la page"})
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.run", Payload: raw})
	if err != nil || res.OK || !strings.Contains(res.Message, "lpc") {
		t.Fatalf("cycle with a cited law must stop: %+v %v", res, err)
	}
	if steps, ok := res.Data.(map[string]any)["steps"]; ok && steps != nil {
		t.Fatal("no worker ran after the stop")
	}
	raw, _ = json.Marshal(map[string]string{"id": "wordpress", "repo": "Evolu-Jeunes/ProximityApp", "task": "contexte local"})
	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.run", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("clean cycle: %+v %v", res, err)
	}
	if got := k.Recall("fusion"); len(got) == 0 {
		t.Fatal("the patrol left no memory")
	}
}
