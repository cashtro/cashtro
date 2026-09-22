package vapi

import (
	"context"
	"encoding/json"
	"os"
	"strings"
)

const assistantName = "Cashtro Teal"

// Assistant is a Vapi assistant record.
type Assistant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TransientAssistant is a full assistant payload so chat/web work
// without a pre-created dashboard assistant.
func TransientAssistant(script string, v Voice) map[string]any {
	if strings.TrimSpace(script) == "" {
		script = "You are the Cashtro / Teal voice on Vapi. Talk with Castro. FR/EN. Keep answers short."
	}
	return map[string]any{
		"name":         assistantName,
		"firstMessage": "Hey Castro. Teal on the line.",
		"voice":        VoicePayload(v),
		"model": map[string]any{
			"provider": "openai",
			"model":    "gpt-4o-mini",
			"messages": []map[string]string{
				{"role": "system", "content": script},
			},
		},
	}
}

// CreateAssistant POSTs /assistant.
func (c *Client) CreateAssistant(ctx context.Context, spec map[string]any) (Assistant, error) {
	var out Assistant
	err := c.do(ctx, "POST", "/assistant", spec, &out)
	return out, err
}

// ListAssistants GETs /assistant.
func (c *Client) ListAssistants(ctx context.Context) ([]Assistant, error) {
	var raw json.RawMessage
	if err := c.do(ctx, "GET", "/assistant", nil, &raw); err != nil {
		return nil, err
	}
	var list []Assistant
	if err := json.Unmarshal(raw, &list); err == nil {
		return list, nil
	}
	var wrap struct {
		Results []Assistant `json:"results"`
		Data    []Assistant `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, err
	}
	if len(wrap.Results) > 0 {
		return wrap.Results, nil
	}
	return wrap.Data, nil
}

// EnsureAssistant returns the Cashtro Teal assistant, creating it if needed.
func (c *Client) EnsureAssistant(ctx context.Context, script string, v Voice) (Assistant, error) {
	if id := strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID")); id != "" {
		return Assistant{ID: id, Name: assistantName}, nil
	}
	list, err := c.ListAssistants(ctx)
	if err == nil {
		for _, a := range list {
			if a.Name == assistantName && a.ID != "" {
				return a, nil
			}
		}
	}
	return c.CreateAssistant(ctx, TransientAssistant(script, v))
}
