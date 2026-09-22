package fleet

import "testing"

func TestKernelAgentsMapped(t *testing.T) {
	want := []string{
		"init", "delivery", "router", "research", "explorer",
		"operator", "reviewer", "architect", "deploy", "security",
		"memory", "comms", "planner", "investigator",
	}
	got := map[string]int{}
	for _, e := range Edges {
		if e.Agent.Kind == "kernel" {
			got[e.Agent.ID]++
		}
	}
	for _, id := range want {
		if got[id] == 0 {
			t.Fatalf("kernel agent %s has no Worked edges", id)
		}
	}
}

func TestEveryProjectHasWork(t *testing.T) {
	want := []string{
		"cashtro", "evolu-jeunes", "btk-avocats", "md-clinic",
		"solution-hypotheque-qc", "educonnexion", "proximity",
		"scanapp", "cashtro-catalog", "axel", "happier", "wealth", "n8n-desk",
	}
	got := map[string]int{}
	for _, e := range Edges {
		got[e.Project.Slug]++
	}
	for _, slug := range want {
		if got[slug] == 0 {
			t.Fatalf("project %s has no agent that worked it", slug)
		}
	}
}

func TestCashtroHasKernelAndCloud(t *testing.T) {
	var kernel, cloud int
	for _, e := range Edges {
		if e.Project.Slug != "cashtro" {
			continue
		}
		switch e.Agent.Kind {
		case "kernel":
			kernel++
		case "cloud":
			cloud++
		}
	}
	if kernel < 14 {
		t.Fatalf("cashtro kernel edges = %d, want >= 14", kernel)
	}
	if cloud < 1 {
		t.Fatalf("cashtro cloud edges = %d, want >= 1", cloud)
	}
}

func TestNoEmptyIDs(t *testing.T) {
	if len(Edges) == 0 {
		t.Fatal("Edges is empty — fleet map did not init")
	}
	for i, e := range Edges {
		if e.Agent.ID == "" || e.Project.Slug == "" {
			t.Fatalf("edge %d missing agent or project", i)
		}
	}
}
