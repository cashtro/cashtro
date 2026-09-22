package model

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBusUsesBothLanes(t *testing.T) {
	var ollamaModels, externalModels []string
	ollamaSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("ollama path %s", r.URL.Path)
		}
		var req ollamaChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Stream {
			t.Fatal("ollama stream must be off")
		}
		ollamaModels = append(ollamaModels, req.Model)
		_ = json.NewEncoder(w).Encode(ollamaChatResponse{Model: req.Model, Message: Message{Role: "assistant", Content: "local"}})
	}))
	defer ollamaSrv.Close()

	extSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatAPIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		externalModels = append(externalModels, req.Model)
		if req.Model == AlsoModel && (req.Reasoning == nil || req.Reasoning.Effort != "max") {
			t.Fatalf("glm reasoning = %+v", req.Reasoning)
		}
		if req.Model == DefaultModel && req.Reasoning != nil {
			t.Fatal("kimi must not send glm reasoning")
		}
		_ = json.NewEncoder(w).Encode(chatAPIResponse{
			Model: req.Model,
			Choices: []struct {
				Message Message `json:"message"`
			}{{Message: Message{Role: "assistant", Content: "cloud"}}},
		})
	}))
	defer extSrv.Close()

	bus := &Bus{
		Internal: &Ollama{BaseURL: ollamaSrv.URL, Model: DefaultOllamaModel, HTTP: ollamaSrv.Client(), live: true},
		External: &Client{BaseURL: extSrv.URL, APIKey: "test-key", Model: DefaultModel, Also: AlsoModel, HTTP: extSrv.Client()},
	}
	parts, err := bus.Chat(context.Background(), ChatRequest{
		Route:    "both",
		Messages: []Message{{Role: "user", Content: "ping"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 3 {
		t.Fatalf("parts = %d, want 3 (ollama + kimi + glm)", len(parts))
	}
	if parts[0].Route != "internal" || parts[0].Model != DefaultOllamaModel || parts[0].Content != "local" {
		t.Fatalf("internal = %+v", parts[0])
	}
	if parts[1].Route != "external" || parts[1].Model != DefaultModel {
		t.Fatalf("kimi = %+v", parts[1])
	}
	if parts[2].Route != "external" || parts[2].Model != AlsoModel {
		t.Fatalf("glm = %+v", parts[2])
	}
	if len(ollamaModels) != 1 || len(externalModels) != 2 {
		t.Fatalf("calls ollama=%v external=%v", ollamaModels, externalModels)
	}
}

func TestBusUnbound(t *testing.T) {
	bus := &Bus{}
	_, err := bus.Chat(context.Background(), ChatRequest{
		Route:    "both",
		Messages: []Message{{Role: "user", Content: "ping"}},
	})
	if !errors.Is(err, ErrUnbound) {
		t.Fatalf("err = %v", err)
	}
	card := bus.Card()
	if card.Internal.Provider != "ollama" || card.External.Model != DefaultModel || card.External.Also != AlsoModel {
		t.Fatalf("card = %+v", card)
	}
}
