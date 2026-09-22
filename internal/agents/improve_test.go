package agents

import (
	"strings"
	"testing"
)

func TestDepartmentFillsItsOwnGap(t *testing.T) {
	first, ok := Improve("ingenierie", "on fait le site")
	if !ok || first.Specialized || len(first.Lacunes) != 2 {
		t.Fatalf("first = %+v", first)
	}
	second, ok := Improve("ingenierie", first.Result)
	if !ok || !second.Specialized || len(second.Lacunes) != 0 {
		t.Fatalf("second = %+v", second)
	}
	if !strings.Contains(second.Result, "Deux façons") || !strings.Contains(second.Result, "Limite") {
		t.Fatalf("result = %s", second.Result)
	}
	for _, d := range Chart().Departments {
		if len(d.Craft) < 2 {
			t.Fatalf("%s has no craft", d.ID)
		}
		vague, ok := Improve(d.ID, "à faire")
		if !ok || vague.Specialized {
			t.Fatalf("%s did not see its gap: %+v", d.ID, vague)
		}
	}
}
