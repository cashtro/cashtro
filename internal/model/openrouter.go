// Package model is the optional brain bus for Cashtro OS.
//
// OpenRouter is not required to boot the kernel. Delivery, the process
// table, and the journal run without a key. Bind OPENROUTER_API_KEY when
// an agentic needs to think; one key reaches many models.
package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the OpenRouter chat API root.
	DefaultBaseURL = "https://openrouter.ai/api/v1"
	// DefaultModel is the primary brain: Moonshot Kimi K3.
	DefaultModel = "moonshotai/kimi-k3"
	// AlsoModel is the second brain: Z.ai GLM top tier.
	AlsoModel = "z-ai/glm-5.3"
	// GLMReasoningEffort is the max reasoning effort for the GLM route.
	GLMReasoningEffort = "max"
)

// Message is one chat turn.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is what agentics send to the router.
type ChatRequest struct {
	Route    string    `json:"route,omitempty"` // internal, external, or both
	Model    string    `json:"model,omitempty"`
	Messages []Message `json:"messages"`
}

// ChatResponse is the model reply we keep.
type ChatResponse struct {
	Route   string `json:"route,omitempty"`
	Model   string `json:"model"`
	Content string `json:"content"`
}

// Status is the public bind card.
type Status struct {
	Provider  string `json:"provider"`
	Bound     bool   `json:"bound"`
	Model     string `json:"model"`
	Also      string `json:"also"`
	Reasoning string `json:"reasoning,omitempty"`
	BaseURL   string `json:"baseUrl"`
	Hint      string `json:"hint"`
}

// Client talks to an OpenAI-compatible OpenRouter endpoint.
type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	Also    string
	HTTP    *http.Client
	Referer string
	Title   string
}

// FromEnv returns a client when OPENROUTER_API_KEY is set, otherwise nil.
func FromEnv() *Client {
	key := strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY"))
	if key == "" {
		return nil
	}
	modelName := strings.TrimSpace(os.Getenv("OPENROUTER_MODEL"))
	if modelName == "" {
		modelName = DefaultModel
	}
	also := strings.TrimSpace(os.Getenv("OPENROUTER_ALSO_MODEL"))
	if also == "" {
		also = AlsoModel
	}
	base := strings.TrimSpace(os.Getenv("OPENROUTER_BASE_URL"))
	if base == "" {
		base = DefaultBaseURL
	}
	return &Client{
		BaseURL: strings.TrimRight(base, "/"),
		APIKey:  key,
		Model:   modelName,
		Also:    also,
		HTTP:    &http.Client{Timeout: 45 * time.Second},
		Referer: "https://github.com/cashtro/cashtro",
		Title:   "Cashtro OS",
	}
}

// Bound reports whether a client can call a model.
func Bound(c *Client) bool {
	return c != nil && c.APIKey != ""
}

// Card returns the bind status for the desk.
func Card(c *Client) Status {
	if !Bound(c) {
		return Status{
			Provider:  "openrouter",
			Bound:     false,
			Model:     DefaultModel,
			Also:      AlsoModel,
			Reasoning: GLMReasoningEffort,
			BaseURL:   DefaultBaseURL,
			Hint:      "Kernel is up without a key. Route is Kimi K3 + GLM 5.3 max. Set OPENROUTER_API_KEY to bind.",
		}
	}
	also := c.Also
	if also == "" {
		also = AlsoModel
	}
	return Status{
		Provider:  "openrouter",
		Bound:     true,
		Model:     c.Model,
		Also:      also,
		Reasoning: GLMReasoningEffort,
		BaseURL:   c.BaseURL,
		Hint:      "Router live. Primary " + c.Model + ", also " + also + " (reasoning " + GLMReasoningEffort + ").",
	}
}

type reasoningOpt struct {
	Effort string `json:"effort,omitempty"`
}

type chatAPIRequest struct {
	Model     string        `json:"model"`
	Messages  []Message     `json:"messages"`
	Reasoning *reasoningOpt `json:"reasoning,omitempty"`
}

type chatAPIResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Chat sends a completion request.
func (c *Client) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if !Bound(c) {
		return ChatResponse{}, fmt.Errorf("openrouter unbound")
	}
	model := req.Model
	if model == "" {
		model = c.Model
	}
	if len(req.Messages) == 0 {
		return ChatResponse{}, fmt.Errorf("messages required")
	}
	apiReq := chatAPIRequest{Model: model, Messages: req.Messages}
	if glmMax(model, c) {
		apiReq.Reasoning = &reasoningOpt{Effort: GLMReasoningEffort}
	}
	body, err := json.Marshal(apiReq)
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Referer != "" {
		httpReq.Header.Set("HTTP-Referer", c.Referer)
	}
	if c.Title != "" {
		httpReq.Header.Set("X-Title", c.Title)
	}

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	res, err := httpClient.Do(httpReq)
	if err != nil {
		return ChatResponse{}, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return ChatResponse{}, err
	}
	var parsed chatAPIResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ChatResponse{}, fmt.Errorf("openrouter decode: %w", err)
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return ChatResponse{}, fmt.Errorf("openrouter: %s", parsed.Error.Message)
	}
	if res.StatusCode >= 300 {
		return ChatResponse{}, fmt.Errorf("openrouter http %d", res.StatusCode)
	}
	if len(parsed.Choices) == 0 {
		return ChatResponse{}, fmt.Errorf("openrouter: empty choices")
	}
	outModel := parsed.Model
	if outModel == "" {
		outModel = model
	}
	return ChatResponse{Route: "external", Model: outModel, Content: parsed.Choices[0].Message.Content}, nil
}

// glmMax reports whether this call should run the GLM route at max reasoning.
func glmMax(modelName string, c *Client) bool {
	also := AlsoModel
	if c != nil && c.Also != "" {
		also = c.Also
	}
	name := strings.ToLower(modelName)
	return name == strings.ToLower(also) || strings.Contains(name, "glm")
}
