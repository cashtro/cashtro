package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestEngineJoinsConstellations(t *testing.T) {
	if got := TaxQuebec(10000); got.TPS != 500 || got.TVQ != 998 || got.Total != 11498 {
		t.Fatalf("tax = %+v", got)
	}
	items := Constellations()
	if len(items) != 14 || emptyCount(items) != 0 {
		t.Fatalf("constellations = %d empty %d", len(items), emptyCount(items))
	}
	if OneAnchor(items) != OneAnchor(Constellations()) {
		t.Fatal("anchor moved")
	}

	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(map[string]any{"cents": 10000, "task": "brancher"})
	res, err := k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.engine", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("engine: %+v %v", res, err)
	}
	data := res.Data.(map[string]any)
	if data["empty"].(int) != 0 || data["anchor"] == "" {
		t.Fatalf("engine data = %+v", data)
	}
	res, err = k.Invoke(context.Background(), "init", kernel.Call{Capability: "os.engine", Payload: raw})
	if err != nil || !res.OK {
		t.Fatalf("os.engine: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{
		Capability: "manager.check",
		Payload:    []byte(`{"payee":"Proximity","cents":2500,"memo":"local"}`),
	})
	if err != nil || !res.OK || !strings.Contains(res.Message, "brouillon") {
		t.Fatalf("check: %+v %v", res, err)
	}
	check := res.Data.(Check)
	if check.Status != "brouillon" || check.Confirm == 0 {
		t.Fatalf("check = %+v", check)
	}

	res, err = k.Invoke(context.Background(), "manager", kernel.Call{Capability: "manager.lens"})
	if err != nil || !res.OK {
		t.Fatalf("lens: %+v %v", res, err)
	}
	lens := res.Data.(map[string]any)
	if lens["empty"].(int) != 0 || lens["constellations"].(int) != 14 || lens["brains"].(int) != 14 {
		t.Fatalf("lens = %+v", lens)
	}
	if data["analyzed"].(int) != len(Subjects()) || data["brains"].(int) != len(Brains()) {
		t.Fatalf("engine analysis = %+v", data)
	}
	if _, err := SaveBlueprint(k); err != nil {
		t.Fatal(err)
	}
}
