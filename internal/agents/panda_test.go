package agents

import "testing"

func TestPandaScanFindsTheSchoolBreak(t *testing.T) {
	report := ReadPandaScan()
	if report.Repo != "Evolu-Jeunes/Panda" || report.Version != "v1" || !report.LocalSchool {
		t.Fatalf("report %+v", report)
	}
	var store, school bool
	for _, section := range report.Sections {
		if section.ID == "store" && section.Connected && section.File != "" {
			store = true
		}
		if section.ID == "school" && section.Connected {
			school = true
		}
	}
	if !store || school {
		t.Fatalf("sections %+v", report.Sections)
	}
	seen := map[string]bool{}
	for _, br := range report.Breaks {
		seen[br.ID] = true
		if br.Fact == "" {
			t.Fatalf("empty break %+v", br)
		}
	}
	for _, id := range []string{"school-route", "price-drawer", "store-line", "two-chains", "admin-demo"} {
		if !seen[id] {
			t.Fatalf("missing break %s in %+v", id, report.Breaks)
		}
	}
	for _, item := range report.Roadmap {
		if item.Built || item.Called {
			t.Fatalf("roadmap live %+v", item)
		}
	}
}

func TestSchoolV1ParksAClassWithoutACall(t *testing.T) {
	school := NewSchool()
	if _, err := school.Add("  ", "10"); err == nil {
		t.Fatal("empty name parked")
	}
	class, err := school.Add("Intro class", "120 CAD")
	if err != nil {
		t.Fatal(err)
	}
	if class.Charged || class.Called || class.Version != "v1" || class.Price != "120 CAD" {
		t.Fatalf("class %+v", class)
	}
	if len(school.List()) != 1 {
		t.Fatalf("list %d", len(school.List()))
	}
}
