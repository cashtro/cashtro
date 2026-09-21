package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestThinkersWearOnDesk(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for _, p := range k.Processes() {
		want, ok := Thinkers[p.Spec.ID]
		if !ok {
			t.Fatalf("process %q has no thinker", p.Spec.ID)
		}
		if p.Spec.Name != want {
			t.Fatalf("%s name = %q, want thinker %q", p.Spec.ID, p.Spec.Name, want)
		}
		seen[p.Spec.ID] = p.Spec.Name
	}
	if len(seen) != len(Thinkers) {
		t.Fatalf("desk thinkers = %d, roster = %d", len(seen), len(Thinkers))
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
		for _, thinker := range Thinkers {
			if ship.Name == thinker {
				t.Fatalf("product ship %q used a thinker name", ship.Name)
			}
		}
	}
	for _, name := range []string{"ScanApp", "Proximity", "BTK Avocats", "Casa Crypto", "Prolifik"} {
		if !found[name] {
			t.Fatalf("expected product %q on the delivery line", name)
		}
	}
}

func TestExplorerFindsThinkerAndProduct(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "explorer", kernel.Call{
		Capability: "explorer.search",
		Payload:    []byte(`{"query":"Hypatia"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("thinker search: %+v %v", res, err)
	}
	if !hitKind(res.Data, "agent", "research") {
		t.Fatalf("Hypatia should find process research: %#v", res.Data)
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
