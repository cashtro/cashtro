package agents

import (
	"context"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestRepairCrewReadsTheWholeChain(t *testing.T) {
	crew := RepairCrew()
	if len(crew) != 20 || crew[0].ID != "maitre" || crew[19].ID != "einstein" {
		t.Fatalf("crew = %d lead %+v", len(crew), crew[0])
	}
	for i, agent := range crew {
		if agent.Order != i+1 || agent.Craft == "" || agent.Name == "" {
			t.Fatalf("seat %+v", agent)
		}
		if !OnChain(agent.ID) {
			t.Fatalf("%s is not on the chain", agent.ID)
		}
	}
	for _, ln := range Lines() {
		if OnChain(ln.ID) {
			t.Fatalf("project %s is on the chain", ln.ID)
		}
	}
	if OnChain("educonnexion") || OnChain("Evolu-Jeunes/educonnexion") || OnChain("chain-1") {
		t.Fatal("Éduconnexion is a website, not a chain")
	}

	var school, wordpress, giant, bot bool
	for _, ln := range Lines() {
		for _, repo := range ln.Repos {
			if repo == "Evolu-Jeunes/educonnexion" && ln.ID == "ecole" {
				school = true
			}
			if repo == "Evolu-Jeunes/educonnexion" && ln.ID == "wordpress" {
				wordpress = true
			}
			if repo == "Evolu-Jeunes/Giant" && ln.ID == "nft-giant" {
				giant = true
			}
			if repo == "Evolu-Jeunes/Giant" && ln.ID == "trading" {
				t.Fatal("the Giant repo is on the trading list")
			}
			if repo == "cashtro/trading_bot-main" && ln.ID == "trading" {
				bot = true
			}
		}
	}
	if !school || !wordpress {
		t.Fatal("Éduconnexion must stay on its website projects")
	}
	if !giant || !bot {
		t.Fatal("Giant and the trading bot each keep their repo")
	}
	sync := GiantSync()
	if sync.Mixed || !sync.Together || sync.GiantRepo == sync.BotRepo || sync.LiveOrder || sync.Trade != "blockchain" {
		t.Fatalf("sync %+v", sync)
	}
	if _, ok := Chain(sync.GiantChain); !ok {
		t.Fatal("Giant chain removed")
	}
	if _, ok := Chain(sync.BotChain); !ok {
		t.Fatal("trading bot chain removed")
	}

	report := SurveyChain()
	if report.Lead != "maitre" || report.Decided != "instinct" || report.Applied || len(report.Agents) != 20 || len(report.Repair) != 5 {
		t.Fatalf("report %+v", report)
	}
	if len(report.Breaks) != 0 {
		t.Fatalf("the chain must not file a project as a break: %+v", report.Breaks)
	}

	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	if k.About().Agents != 15 {
		t.Fatalf("boot agents = %d", k.About().Agents)
	}
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.repair"})
	if err != nil || !res.OK {
		t.Fatalf("repair: %+v %v", res, err)
	}
	got := res.Data.(ChainReport)
	if len(got.Agents) != 20 || got.Lead != "maitre" || len(got.Breaks) != 0 {
		t.Fatalf("invoke lead %s breaks %d", got.Lead, len(got.Breaks))
	}
}
