package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBuildersWearOnDesk(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for _, p := range k.Processes() {
		want, ok := Builders[p.Spec.ID]
		if !ok {
			t.Fatalf("process %q has no builder", p.Spec.ID)
		}
		if p.Spec.Name != want {
			t.Fatalf("%s name = %q, want builder %q", p.Spec.ID, p.Spec.Name, want)
		}
		seen[p.Spec.ID] = p.Spec.Name
	}
	if len(seen) != len(Builders) {
		t.Fatalf("desk builders = %d, roster = %d", len(seen), len(Builders))
	}
	if Builder("init") != "Lautaro" {
		t.Fatalf("init must be Lautaro, got %q", Builder("init"))
	}
	if Builder("operator") != "Builder" {
		t.Fatalf("operator must be Builder, got %q", Builder("operator"))
	}
	if Builder("architect") != "Creator" {
		t.Fatalf("architect must be Creator, got %q", Builder("architect"))
	}
}

func TestProductsStayProducts(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range k.Processes() {
		for _, product := range ProductNames {
			if strings.EqualFold(p.Spec.Name, product) || strings.EqualFold(p.Spec.ID, product) {
				t.Fatalf("process %s wore product name %q", p.Spec.ID, product)
			}
		}
	}
	cat := k.Catalog()
	if cat == nil {
		t.Fatal("delivery board missing")
	}
	found := map[string]bool{}
	for _, ship := range cat.List() {
		found[ship.Name] = true
		for _, builder := range Builders {
			if ship.Name == builder {
				t.Fatalf("product ship %q used a builder name", ship.Name)
			}
		}
	}
	for _, name := range []string{"ScanApp", "Proximity", "BTK Avocats", "Casa Crypto", "Prolifik"} {
		if !found[name] {
			t.Fatalf("expected product %q on the delivery line", name)
		}
	}
}

func TestExplorerFindsBuilderAndProduct(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"Lautaro"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("builder search: %+v %v", res, err)
	}
	if !hitKind(res.Data, "agent", "init") {
		t.Fatalf("Lautaro should find process init: %#v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"Builder"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("builder search: %+v %v", res, err)
	}
	if !hitKind(res.Data, "agent", "operator") {
		t.Fatalf("Builder should find process operator: %#v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"Creator"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("creator search: %+v %v", res, err)
	}
	if !hitKind(res.Data, "agent", "architect") {
		t.Fatalf("Creator should find process architect: %#v", res.Data)
	}

	res, err = k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"Proximity"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("product search: %+v %v", res, err)
	}
	if !hitKind(res.Data, "ship", "proximity") {
		t.Fatalf("Proximity should stay a ship: %#v", res.Data)
	}
}

func TestInvokeStillUsesProcessIDs(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]string{"name": "North desk"})
	res, err := k.Invoke(context.Background(), "delivery", kernel.Call{
		Capability: "delivery.create",
		Payload:    payload,
	})
	if err != nil || !res.OK {
		t.Fatalf("delivery id must still route: %+v %v", res, err)
	}
}

func hitKind(data any, kind, id string) bool {
	hits, ok := data.([]map[string]string)
	if !ok {
		raw, ok := data.([]any)
		if !ok {
			return false
		}
		for _, item := range raw {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if m["kind"] == kind && m["id"] == id {
				return true
			}
		}
		return false
	}
	for _, h := range hits {
		if h["kind"] == kind && h["id"] == id {
			return true
		}
	}
	return false
}
