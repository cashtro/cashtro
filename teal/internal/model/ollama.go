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
	// DefaultOllamaURL is the local Ollama daemon.
	DefaultOllamaURL = "http://127.0.0.1:11434"
	// DefaultOllamaModel is the internal model until OLLAMA_MODEL is set.
	DefaultOllamaModel = "llama3.2"
)

// Ollama is the internal model lane. No API key. The daemon must be up.
type Ollama struct {
	BaseURL string
	Model   string
	HTTP    *http.Client
	live    bool
}

// OllamaFromEnv configures the internal lane and probes the daemon once.
func OllamaFromEnv() *Ollama {
	base := strings.TrimSpace(os.Getenv("OLLAMA_HOST"))
	if base == "" {
		base = DefaultOllamaURL
	}
	name := strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
	if name == "" {
		name = DefaultOllamaModel
	}
	o := &Ollama{
		BaseURL: strings.TrimRight(base, "/"),
		Model:   name,
		HTTP:    &http.Client{Timeout: 120 * time.Second},
	}
	o.live = o.probe()
	return o
}

// Bound reports whether the daemon answered the probe.
func (o *Ollama) Bound() bool {
	return o != nil && o.live
}

func (o *Ollama) probe() bool {
	if o == nil || o.BaseURL == "" {
		return false
	}
	client := o.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.BaseURL+"/api/tags", nil)
	if err != nil {
		return false
	}
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 1<<16))
	return res.StatusCode < 300
}

type ollamaChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ollamaChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Error   string  `json:"error"`
}

// Chat calls the local Ollama /api/chat endpoint.
func (o *Ollama) Chat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if !o.Bound() {
		return ChatResponse{}, fmt.Errorf("ollama unbound")
	}
	name := req.Model
	if name == "" {
		name = o.Model
	}
	if len(req.Messages) == 0 {
		return ChatResponse{}, fmt.Errorf("messages required")
	}
	body, err := json.Marshal(ollamaChatRequest{Model: name, Messages: req.Messages, Stream: false})
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	client := o.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(httpReq)
	if err != nil {
		return ChatResponse{}, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return ChatResponse{}, err
	}
	var parsed ollamaChatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ChatResponse{}, fmt.Errorf("ollama decode: %w", err)
	}
	if parsed.Error != "" {
		return ChatResponse{}, fmt.Errorf("ollama: %s", parsed.Error)
	}
	if res.StatusCode >= 300 {
		return ChatResponse{}, fmt.Errorf("ollama http %d", res.StatusCode)
	}
	outModel := parsed.Model
	if outModel == "" {
		outModel = name
	}
	return ChatResponse{Model: outModel, Content: parsed.Message.Content, Route: "internal"}, nil
}
