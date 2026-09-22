package vapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvDoesNotOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("VAPI_TEST_DOTENV=fromfile\nVAPI_TEST_KEEP=file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("VAPI_TEST_KEEP", "already")
	_ = os.Unsetenv("VAPI_TEST_DOTENV")
	LoadDotEnv(path)
	if os.Getenv("VAPI_TEST_DOTENV") != "fromfile" {
		t.Fatalf("dotenv = %q", os.Getenv("VAPI_TEST_DOTENV"))
	}
	if os.Getenv("VAPI_TEST_KEEP") != "already" {
		t.Fatalf("override = %q", os.Getenv("VAPI_TEST_KEEP"))
	}
}

func TestEnsureAssistantCreatesWhenMissing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/assistant":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/assistant":
			_, _ = w.Write([]byte(`{"id":"asst_new","name":"Cashtro Teal"}`))
		default:
			http.Error(w, r.URL.Path, http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	t.Setenv("VAPI_ASSISTANT_ID", "")
	c := &Client{Key: "secret", BaseURL: srv.URL, HTTP: srv.Client()}
	asst, err := c.EnsureAssistant(context.Background(), "be Teal", DefaultVoice())
	if err != nil || asst.ID != "asst_new" {
		t.Fatalf("ensure = %+v %v", asst, err)
	}
}

func TestClientChatInlineAssistant(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat" {
			http.Error(w, r.URL.Path, http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(`{"id":"chat-inline","output":[{"role":"assistant","content":"inline ok"}]}`))
	}))
	t.Cleanup(srv.Close)
	c := &Client{Key: "secret", BaseURL: srv.URL, HTTP: srv.Client()}
	out, err := c.Chat(context.Background(), ChatRequest{
		Input:     "hi",
		Assistant: TransientAssistant("be Teal", DefaultVoice()),
	})
	if err != nil || out.ID != "chat-inline" || out.Output[0].Content != "inline ok" {
		t.Fatalf("chat = %+v %v", out, err)
	}
}
