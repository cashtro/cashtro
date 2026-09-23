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
	if len(eye.Steps) != 8 || eye.Steps[0] != "veille" || eye.Steps[7] != "sceau" {
		t.Fatalf("chain %+v", eye.Steps)
	}
	if eye.Attacks || eye.Flood || eye.CopiesSecrets || eye.AppliesPatch || eye.DecidedBy != "tartaria" {
		t.Fatal("the eye must withhold secrets and wait for Tartaria")
	}
}
