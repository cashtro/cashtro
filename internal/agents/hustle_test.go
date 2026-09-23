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
	if chain.BestFriend != "instinct" || !chain.KillerInstinct || chain.LoadsInstinct || !chain.OwnBrain || chain.MainBrain != "instinct" || !chain.LoadsVoltron {
		t.Fatal("The Hustle must keep its own brain and stay friends with the main brain")
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
	if !chain.Person || chain.Team != 20 || !chain.UsesAllChains || !chain.PurposeOpen || len(chain.Accounts) != 2 {
		t.Fatalf("person %+v", chain)
	}
	seat := Hustler()
	if !seat.Person || !seat.Internal || !seat.ConceptsHere || !seat.OthersRunWork || seat.Purpose != "" || !seat.PurposeOpen {
		t.Fatalf("seat %+v", seat)
	}
	if len(seat.Questions) != 6 || !seat.Questions[0].Open || seat.Questions[0].Answer != "" {
		t.Fatalf("questions %+v", seat.Questions)
	}
	if _, ok := HustlerAnswer("but", " "); ok {
		t.Fatal("an empty answer must stay open")
	}
	if _, ok := HustlerAnswer("but", "fraude"); ok {
		t.Fatal("an illegal answer must stay open")
	}
	closed, ok := HustlerAnswer("courant", "le concept que je vais te dire")
	if !ok || closed.Open || closed.Answer == "" {
		t.Fatalf("answer %+v %v", closed, ok)
	}
	if ok, _ := HustlerUse("forge", "hustle"); ok {
		t.Fatal("the team does not take a resource outside Voltron")
	}
	ok, msg := HustlerUse("cashtro", "voltron")
	if !ok || msg == "" {
		t.Fatalf("use %v %s", ok, msg)
	}
	if ok, _ := HustlerUse("autre", "voltron"); ok {
		t.Fatal("a resource outside the two GitHubs must stay closed")
	}
}
