package agents

import (
	"context"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestBoundSeatsExecute(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	about := k.About()
	if about.Live < 13 {
		t.Fatalf("live=%d want >=13 (router may stay resident)", about.Live)
	}

	cases := []struct {
		id, cap, want string
		payload       string
	}{
		{"operator", "operator.browse", "operator probed", `{"url":"http://127.0.0.1:8080/health"}`},
		{"reviewer", "reviewer.watch", "reviewer listed", ""},
		{"architect", "architect.plan", "architect planned", `{"goal":"ScanApp concept"}`},
		{"deploy", "deploy.release", "deploy local-only", `{"target":"local"}`},
		{"security", "security.triage", "security flags", ""},
		{"investigator", "investigator.trace", "investigator traced", ""},
	}
	for _, tc := range cases {
		res, err := k.Invoke(context.Background(), tc.id, kernel.Call{Capability: tc.cap, Payload: []byte(tc.payload)})
		if err != nil {
			t.Fatalf("%s: %v", tc.id, err)
		}
		if strings.Contains(res.Message, "ready to bind") {
			t.Fatalf("%s still unbound: %s", tc.id, res.Message)
		}
		if !strings.Contains(res.Message, tc.want) {
			t.Fatalf("%s message=%q", tc.id, res.Message)
		}
	}

	refuse, err := k.Invoke(context.Background(), "architect", kernel.Call{
		Capability: "architect.plan",
		Payload:    []byte(`{"goal":"touch Proximity production"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if refuse.OK {
		t.Fatal("architect must refuse live client production")
	}

	bad, err := k.Invoke(context.Background(), "operator", kernel.Call{
		Capability: "operator.browse",
		Payload:    []byte(`{"url":"https://example.com"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if bad.OK {
		t.Fatal("operator must refuse non-local URLs")
	}
}
