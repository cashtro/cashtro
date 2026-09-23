package agents

import "testing"

func TestEveryProgramHasItsOwnBrain(t *testing.T) {
	main := MainBrain()
	if len(main) != 5 {
		t.Fatalf("main desks %d", len(main))
	}
	mainIDs := map[string]bool{}
	for _, desk := range main {
		mainIDs[desk.ID] = true
	}
	wantMain := []string{"architecte", "cartographe", "forgeron", "orfevre", "hustler"}
	for i, id := range wantMain {
		if main[i].ID != id {
			t.Fatalf("main desk %s", main[i].ID)
		}
	}
	seen := map[string]NicheBrain{}
	for _, brain := range NicheBrains() {
		if _, ok := seen[brain.ID]; ok {
			t.Fatalf("duplicate %s", brain.ID)
		}
		seen[brain.ID] = brain
		if brain.ID == "control" || brain.ID == "instinct" {
			t.Fatalf("the main brain was copied into a niche: %s", brain.ID)
		}
		if !brain.OwnBrain || brain.SameBrain || !brain.SameShape || !brain.FunctionsOnly || brain.LoadsInstinctBrain || brain.MainBrain != "instinct" || brain.DecidedBy != "instinct" {
			t.Fatalf("shape %+v", brain)
		}
		if len(brain.Desks) != 5 || brain.Workers != 900 || len(brain.Closed) == 0 {
			t.Fatalf("desk count %+v", brain.ID)
		}
		for _, desk := range brain.Desks {
			if mainIDs[desk.ID] || len(desk.SubBrains) != 5 || desk.Workers == 0 || len(desk.Directors) == 0 {
				t.Fatalf("desk %+v", desk)
			}
			bridges := 0
			for _, sub := range desk.SubBrains {
				if sub.Bridges {
					bridges++
				}
				if sub.Role == "" {
					t.Fatalf("empty role %s", sub.ID)
				}
			}
			if bridges != 1 {
				t.Fatalf("bridges %s %d", desk.ID, bridges)
			}
		}
	}
	for _, line := range Lines() {
		if line.ID == "control" {
			if _, ok := seen[line.ID]; ok {
				t.Fatal("control is the main brain")
			}
			continue
		}
		if _, ok := seen[line.ID]; !ok {
			t.Fatalf("line %s has no brain of its own", line.ID)
		}
	}
	for _, layer := range Layers() {
		if _, ok := seen[layer.ID]; !ok {
			t.Fatalf("chain %s has no brain of its own", layer.ID)
		}
	}
	eye, ok := seen["eye"]
	if !ok || eye.Desks[0].ID == "architecte" || eye.Desks[0].Name != "La Veille" {
		t.Fatalf("eye %+v", eye.Desks)
	}
}
