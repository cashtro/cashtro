package agents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
)

func TestTealBrainTeamsOnlyAndVapi(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}

	res, err := k.Invoke(context.Background(), "teal", kernel.Call{Capability: "teal.status"})
	if err != nil || !res.OK {
		t.Fatalf("status: %+v %v", res, err)
	}
	st := res.Data.(map[string]any)
	if st["surface"] != "microsoft-teams" || st["voice"] != "vapi" {
		t.Fatalf("status = %#v", st)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{Capability: "teal.vapi"})
	if err != nil || !res.OK {
		t.Fatalf("vapi: %+v %v", res, err)
	}
	card := res.Data.(VapiCard)
	if !strings.Contains(card.ShareURL, "vapi.ai") {
		t.Fatalf("share = %q", card.ShareURL)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.ingest",
		Payload:    []byte(`{"channel":"slack","text":"nope"}`),
	})
	if err != nil || res.OK {
		t.Fatalf("slack should refuse: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.ingest",
		Payload:    []byte(`{"channel":"microsoft-teams","team":"AI Teal","from":"castro","text":"Pandora brainstorm: intern voice desk"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("ingest: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.listen",
		Payload:    []byte(`{"message":{"type":"end-of-call-report","call":{"id":"c1","assistantId":"asst_teal"},"artifact":{"transcript":"How do interns start on Pandora?"},"analysis":{"summary":"Intern asked how to start on Pandora."}}}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("listen: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.question",
		Payload:    []byte(`{"from":"intern-code","prompt":"how do I use vapi on teams"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("question: %+v %v", res, err)
	}
	q := res.Data.(Question)
	if !strings.Contains(strings.ToLower(q.Answer), "vapi") {
		t.Fatalf("answer = %q", q.Answer)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.assign",
		Payload:    []byte(`{"id":"intern-code","help":"pair on Pandora landing"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("assign: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.idea",
		Payload:    []byte(`{"prompt":"voice intern coach"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("idea: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{Capability: "teal.smarter"})
	if err != nil || !res.OK {
		t.Fatalf("smarter: %+v %v", res, err)
	}
	if _, err := k.Process("intern-coach"); err != nil {
		t.Fatalf("expected spawned intern-coach: %v", err)
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{
		Capability: "teal.cursor",
		Payload:    []byte(`{"prompt":"help interns build Pandora"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("cursor: %+v %v", res, err)
	}
	data := res.Data.(map[string]any)
	if !strings.Contains(data["url"].(string), "cursor.com/agents") {
		t.Fatalf("cursor url = %#v", data["url"])
	}
	if len(k.Confirms()) == 0 {
		t.Fatal("cursor launch must park a human confirm")
	}

	res, err = k.Invoke(context.Background(), "teal", kernel.Call{Capability: "teal.card"})
	if err != nil || !res.OK {
		t.Fatalf("card: %+v %v", res, err)
	}
	cardMap := res.Data.(map[string]any)
	if cardMap["team"] != tealTeamName || cardMap["channel"] != "microsoft-teams" {
		t.Fatalf("card = %#v", cardMap)
	}

	got := k.Recall("pandora")
	if len(got) < 3 {
		t.Fatalf("pandora memory = %d", len(got))
	}
}

func TestParseVapiListen(t *testing.T) {
	raw, _ := json.Marshal(map[string]any{
		"message": map[string]any{
			"call":     map[string]any{"id": "x", "assistantId": "a1"},
			"artifact": map[string]any{"transcript": "hello intern"},
			"analysis": map[string]any{"summary": "greeting"},
		},
	})
	got := parseVapiListen(raw)
	if got.ID != "x" || got.Assistant != "a1" || got.Transcript != "hello intern" || got.Summary != "greeting" {
		t.Fatalf("parsed = %+v", got)
	}
}

func TestIsTeamsOnly(t *testing.T) {
	if !isTeamsOnly("", "") || !isTeamsOnly("microsoft-teams", "") || !isTeamsOnly("AI Teal", "teams") {
		t.Fatal("expected teams allowed")
	}
	if isTeamsOnly("slack", "") || isTeamsOnly("mail", "outlook") {
		t.Fatal("expected slack/mail refused")
	}
}
