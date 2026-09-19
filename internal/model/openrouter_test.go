package model

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCardUnbound(t *testing.T) {
	st := Card(nil)
	if st.Bound || st.Provider != "openrouter" {
		t.Fatalf("card = %+v", st)
	}
}

func TestChat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("auth = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Title") != "Cashtro OS" {
			t.Fatalf("title = %q", r.Header.Get("X-Title"))
		}
		var req chatAPIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Model != "openai/gpt-4o-mini" || req.Messages[0].Content != "ping" {
			t.Fatalf("req = %+v", req)
		}
		_ = json.NewEncoder(w).Encode(chatAPIResponse{
			Model: req.Model,
			Choices: []struct {
				Message Message `json:"message"`
			}{{Message: Message{Role: "assistant", Content: "pong"}}},
		})
	}))
	defer srv.Close()

	c := &Client{
		BaseURL: srv.URL,
		APIKey:  "test-key",
		Model:   DefaultModel,
		HTTP:    srv.Client(),
		Title:   "Cashtro OS",
	}
	got, err := c.Chat(context.Background(), ChatRequest{
		Messages: []Message{{Role: "user", Content: "ping"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Content != "pong" {
		t.Fatalf("content = %q", got.Content)
	}
}

func TestChatUnbound(t *testing.T) {
	var c *Client
	_, err := c.Chat(context.Background(), ChatRequest{Messages: []Message{{Role: "user", Content: "x"}}})
	if err == nil {
		t.Fatal("expected unbound error")
	}
}
