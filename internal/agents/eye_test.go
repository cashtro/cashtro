package agents

import "testing"

func TestTheEyeWatchesOurRepos(t *testing.T) {
	eye := TheEye()
	if eye.Name != "The Eye" || eye.Repo != "Evolu-Jeunes/Forge" || eye.Brain != "eye" {
		t.Fatalf("name %+v", eye)
	}
	if len(eye.Scope) != 2 || eye.Scope[0] != "cashtro" || eye.Scope[1] != "Evolu-Jeunes" {
		t.Fatalf("scope %+v", eye.Scope)
	}
	if len(eye.Steps) != 9 || eye.Steps[0] != "veille" || eye.Steps[5] != "instinct" || eye.Steps[6] != "rustine" || eye.Steps[8] != "sceau" {
		t.Fatalf("chain %+v", eye.Steps)
	}
	if eye.Gate != "cashtro/epicenter" || !eye.ThroughInstinct {
		t.Fatalf("gate %+v", eye)
	}
	if eye.Attacks || eye.Flood || eye.CopiesSecrets || eye.AppliesPatch || eye.DecidedBy != "instinct" {
		t.Fatal("the eye must withhold secrets and wait for Instinct")
	}
	if ok, _ := EyeChange("eye", "instinct"); ok {
		t.Fatal("a change outside Instinct must stay closed")
	}
	if ok, _ := EyeChange("cashtro/epicenter", "voltron"); ok {
		t.Fatal("another seat must not open the change")
	}
	ok, msg := EyeChange("instinct", "instinct")
	if !ok || msg != "brouillon ouvert. rien n'est appliqué" {
		t.Fatalf("instinct change = %v %s", ok, msg)
	}
}
