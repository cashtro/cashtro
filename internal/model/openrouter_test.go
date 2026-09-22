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
		if req.Model != DefaultModel || req.Messages[0].Content != "ping" {
			t.Fatalf("req = %+v", req)
		}
		if req.Reasoning != nil {
			t.Fatalf("kimi route must not force glm reasoning: %+v", req.Reasoning)
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

func TestChatGLMMax(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatAPIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Model != AlsoModel {
			t.Fatalf("model = %s", req.Model)
		}
		if req.Reasoning == nil || req.Reasoning.Effort != GLMReasoningEffort {
			t.Fatalf("reasoning = %+v", req.Reasoning)
		}
		_ = json.NewEncoder(w).Encode(chatAPIResponse{
			Model: req.Model,
			Choices: []struct {
				Message Message `json:"message"`
			}{{Message: Message{Role: "assistant", Content: "ok"}}},
		})
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIKey: "test-key", Model: DefaultModel, Also: AlsoModel, HTTP: srv.Client()}
	got, err := c.Chat(context.Background(), ChatRequest{
		Model:    AlsoModel,
		Messages: []Message{{Role: "user", Content: "plan"}},
	})
	if err != nil || got.Content != "ok" || got.Model != AlsoModel {
		t.Fatalf("got %+v err %v", got, err)
	}
}

func TestChatUnbound(t *testing.T) {
	var c *Client
	_, err := c.Chat(context.Background(), ChatRequest{Messages: []Message{{Role: "user", Content: "x"}}})
	if err == nil {
		t.Fatal("expected unbound error")
	}
}
