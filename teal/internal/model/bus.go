package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrUnbound is returned when the requested lane has no live model.
var ErrUnbound = errors.New("model route unbound")

// Bus is both model lanes: Ollama internal, OpenRouter external.
type Bus struct {
	Internal *Ollama
	External *Client
}

// BusFromEnv loads both lanes from the environment.
func BusFromEnv() *Bus {
	return &Bus{Internal: OllamaFromEnv(), External: FromEnv()}
}

// Live reports whether at least one lane can answer.
func (b *Bus) Live() bool {
	if b == nil {
		return false
	}
	return b.Internal.Bound() || Bound(b.External)
}

// Lane is one side of the route card.
type Lane struct {
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Bound     bool   `json:"bound"`
	Model     string `json:"model"`
	Also      string `json:"also,omitempty"`
	Reasoning string `json:"reasoning,omitempty"`
	BaseURL   string `json:"baseUrl"`
}

// RouteCard describes both lanes.
type RouteCard struct {
	Internal Lane   `json:"internal"`
	External Lane   `json:"external"`
	Hint     string `json:"hint"`
}

// Card returns the public status of both lanes.
func (b *Bus) Card() RouteCard {
	inModel, inURL := DefaultOllamaModel, DefaultOllamaURL
	inBound := false
	if b != nil && b.Internal != nil {
		if b.Internal.Model != "" {
			inModel = b.Internal.Model
		}
		if b.Internal.BaseURL != "" {
			inURL = b.Internal.BaseURL
		}
		inBound = b.Internal.Bound()
	}
	ext := Card(nil)
	if b != nil {
		ext = Card(b.External)
	}
	hint := "Both lanes: Ollama internal, Kimi K3 + GLM 5.3 max external."
	if !inBound && !ext.Bound {
		hint = "Neither lane is bound. Start Ollama or set OPENROUTER_API_KEY. Route stays Kimi K3 + GLM 5.3 max outside, Ollama inside."
	}
	return RouteCard{
		Internal: Lane{
			Name: "internal", Provider: "ollama", Bound: inBound,
			Model: inModel, BaseURL: inURL,
		},
		External: Lane{
			Name: "external", Provider: ext.Provider, Bound: ext.Bound,
			Model: ext.Model, Also: ext.Also, Reasoning: ext.Reasoning, BaseURL: ext.BaseURL,
		},
		Hint: hint,
	}
}

// Chat runs the requested lane. route=both calls Ollama, Kimi K3, and GLM max.
func (b *Bus) Chat(ctx context.Context, req ChatRequest) ([]ChatResponse, error) {
	if b == nil {
		return nil, ErrUnbound
	}
	route := strings.ToLower(strings.TrimSpace(req.Route))
	if route == "" {
		route = "both"
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages required")
	}
	var out []ChatResponse
	if route == "internal" || route == "both" {
		if b.Internal.Bound() {
			resp, err := b.Internal.Chat(ctx, req)
			if err != nil {
				return nil, err
			}
			out = append(out, resp)
		} else if route == "internal" {
			return nil, fmt.Errorf("%w: ollama", ErrUnbound)
		}
	}
	if route == "external" || route == "both" {
		if Bound(b.External) {
			parts, err := b.external(ctx, req)
			if err != nil {
				return nil, err
			}
			out = append(out, parts...)
		} else if route == "external" {
			return nil, fmt.Errorf("%w: openrouter", ErrUnbound)
		}
	}
	if len(out) == 0 {
		return nil, ErrUnbound
	}
	return out, nil
}

func (b *Bus) external(ctx context.Context, req ChatRequest) ([]ChatResponse, error) {
	models := []string{req.Model}
	if req.Model == "" || strings.EqualFold(req.Model, "both") {
		primary := b.External.Model
		if primary == "" {
			primary = DefaultModel
		}
		also := b.External.Also
		if also == "" {
			also = AlsoModel
		}
		models = []string{primary, also}
	}
	var out []ChatResponse
	for _, name := range models {
		resp, err := b.External.Chat(ctx, ChatRequest{Model: name, Messages: req.Messages})
		if err != nil {
			return nil, err
		}
		resp.Route = "external"
		out = append(out, resp)
	}
	return out, nil
}
