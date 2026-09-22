package vapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVoiceCatalogAndLookup(t *testing.T) {
	got := Voices()
	if len(got) < 10 {
		t.Fatalf("voices = %d", len(got))
	}
	v, ok := VoiceByID("rachel")
	if !ok || v.Provider != "11labs" {
		t.Fatalf("rachel = %+v", v)
	}
	if _, ok := VoiceByID("denise"); !ok {
		t.Fatal("expected FR Denise")
	}
	if _, ok := VoiceByID("nope"); ok {
		t.Fatal("unknown voice")
	}
	p := VoicePayload(v)
	if p["provider"] != "11labs" || p["voiceId"] != "rachel" {
		t.Fatalf("payload = %#v", p)
	}
}

func TestUnboundClient(t *testing.T) {
	c := &Client{}
	if c.Bound() {
		t.Fatal("empty client bound")
	}
	_, err := c.Chat(nil, ChatRequest{Input: "hi"})
	if err != ErrUnbound {
		t.Fatalf("err = %v", err)
	}
}

func TestEnvCardHasChangeOptions(t *testing.T) {
	card := EnvCard()
	opts, _ := card["changeVoice"].([]string)
	if len(opts) < 4 {
		t.Fatalf("changeVoice = %#v", card["changeVoice"])
	}
	joined := strings.Join(opts, " ")
	if !strings.Contains(joined, "dashboard") || !strings.Contains(joined, "VAPI_VOICE_ID") {
		t.Fatalf("missing change paths: %s", joined)
	}
}

func TestClientChatCallPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/chat":
			_, _ = w.Write([]byte(`{"id":"chat-1","output":[{"role":"assistant","content":"hello Castro"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/call":
			_, _ = w.Write([]byte(`{"id":"call-1","status":"queued","type":"outboundPhoneCall"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/assistant/asst_1":
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, r.URL.Path, http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	c := &Client{Key: "secret", BaseURL: srv.URL, HTTP: srv.Client()}
	out, err := c.Chat(context.Background(), ChatRequest{AssistantID: "asst_1", Input: "hi"})
	if err != nil || out.ID != "chat-1" || out.Output[0].Content != "hello Castro" {
		t.Fatalf("chat = %+v %v", out, err)
	}
	call, err := c.Call(context.Background(), CallRequest{AssistantID: "asst_1", PhoneNumberID: "pn_1", Customer: map[string]any{"number": "+1"}})
	if err != nil || call.ID != "call-1" {
		t.Fatalf("call = %+v %v", call, err)
	}
	if err := c.PatchVoice(context.Background(), "asst_1", Voices()[0]); err != nil {
		t.Fatal(err)
	}
}
