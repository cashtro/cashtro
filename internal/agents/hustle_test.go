package agents

import "testing"

func TestTheHustleReadsAndStaysLegal(t *testing.T) {
	chain := Hustle()
	if chain.Name != "The Hustle" || chain.Repo != "Evolu-Jeunes/Hustle" || chain.Brain != "hustle" {
		t.Fatalf("name %+v", chain)
	}
	if chain.Agents != 20 || chain.Brains != 100 {
		t.Fatalf("swarm %+v", chain)
	}
	if chain.BestFriend != "instinct" || !chain.KillerInstinct || !chain.LoadsInstinct || !chain.LoadsVoltron {
		t.Fatal("The Hustle must stay Epicenter Einstein's friend and load both brains")
	}
	if chain.Weapons || chain.Attacks || chain.Illegal || chain.Sent || chain.Paid || chain.Launched {
		t.Fatal("The Hustle must not attack, take an illegal path, charge, or launch")
	}
	if chain.OpenedBy != "voltron" || chain.OrderedBy != "scrum" || chain.DecidedBy != "instinct" {
		t.Fatalf("seats %+v", chain)
	}
	var found bool
	for _, brain := range Brains() {
		if brain.ID == "hustle" {
			found = true
		}
	}
	if !found {
		t.Fatal("The Hustle is not on the board")
	}
}
