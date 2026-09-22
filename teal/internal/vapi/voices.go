// Package vapi is the in-repo Vapi (Vappy) voice lane.
//
// Talk and outbound calls run through Cashtro OS. Voice is a catalog you can
// switch on the desk, via env, per-call override, or the Vapi dashboard.
// Live HTTP only fires when VAPI_API_KEY is set — outbound still parks a
// human confirm first.
package vapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	APIBase      = "https://api.vapi.ai"
	ProductURL   = "https://vapi.ai/"
	DashboardURL = "https://dashboard.vapi.ai/"
	DocsURL      = "https://docs.vapi.ai/"
	GitHubURL    = "https://github.com/VapiAI"
	VoiceDocs    = "https://docs.vapi.ai/assistants/voice-formatting-plan"
	CallDocs     = "https://docs.vapi.ai/calls/outbound-calling"
	ChatDocs     = "https://docs.vapi.ai/chat/quickstart"
)

// ErrUnbound means no VAPI_API_KEY — local talk still works.
var ErrUnbound = errors.New("vapi unbound")

// Voice is one switchable TTS option.
type Voice struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	VoiceID  string `json:"voiceId"`
	Lang     string `json:"lang"`
	Change   string `json:"change"`
}

// Voices is the in-repo picker. Castro can change voice here, in env,
// per call, or in the Vapi dashboard Voice dropdown.
func Voices() []Voice {
	change := "POST /api/vapi/voice · desk picker · dashboard.vapi.ai Assistants → Voice · env VAPI_VOICE_ID · per-call assistantOverrides.voice"
	return []Voice{
		{ID: "elliot", Name: "Elliot", Provider: "vapi", VoiceID: "Elliot", Lang: "en", Change: change},
		{ID: "cole", Name: "Cole", Provider: "vapi", VoiceID: "Cole", Lang: "en", Change: change},
		{ID: "rachel", Name: "Rachel", Provider: "11labs", VoiceID: "rachel", Lang: "en", Change: change},
		{ID: "adam", Name: "Adam", Provider: "11labs", VoiceID: "adam", Lang: "en", Change: change},
		{ID: "bella", Name: "Bella", Provider: "11labs", VoiceID: "bella", Lang: "en", Change: change},
		{ID: "josh", Name: "Josh", Provider: "11labs", VoiceID: "josh", Lang: "en", Change: change},
		{ID: "alloy", Name: "Alloy", Provider: "openai", VoiceID: "alloy", Lang: "en", Change: change},
		{ID: "nova", Name: "Nova", Provider: "openai", VoiceID: "nova", Lang: "en", Change: change},
		{ID: "onyx", Name: "Onyx", Provider: "openai", VoiceID: "onyx", Lang: "en", Change: change},
		{ID: "shimmer", Name: "Shimmer", Provider: "openai", VoiceID: "shimmer", Lang: "en", Change: change},
		{ID: "asteria", Name: "Aura Asteria", Provider: "deepgram", VoiceID: "aura-asteria-en", Lang: "en", Change: change},
		{ID: "luna", Name: "Aura Luna", Provider: "deepgram", VoiceID: "aura-luna-en", Lang: "en", Change: change},
		{ID: "orion", Name: "Aura Orion", Provider: "deepgram", VoiceID: "aura-orion-en", Lang: "en", Change: change},
		{ID: "sonic", Name: "Cartesia Sonic", Provider: "cartesia", VoiceID: "sonic", Lang: "en", Change: change},
		{ID: "aria", Name: "Azure Aria", Provider: "azure", VoiceID: "en-US-AriaNeural", Lang: "en", Change: change},
		{ID: "denise", Name: "Denise FR", Provider: "azure", VoiceID: "fr-FR-DeniseNeural", Lang: "fr", Change: change},
		{ID: "henri", Name: "Henri FR", Provider: "azure", VoiceID: "fr-FR-HenriNeural", Lang: "fr", Change: change},
	}
}

// VoiceByID returns a catalog voice (id, name, or provider:voiceId).
func VoiceByID(id string) (Voice, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	for _, v := range Voices() {
		if strings.ToLower(v.ID) == id || strings.ToLower(v.Name) == id || strings.ToLower(v.VoiceID) == id || strings.ToLower(v.Provider+":"+v.VoiceID) == id {
			return v, true
		}
	}
	return Voice{}, false
}

// DefaultVoice is Rachel unless VAPI_VOICE_ID is set.
func DefaultVoice() Voice {
	if v, ok := VoiceByID(os.Getenv("VAPI_VOICE_ID")); ok {
		return v
	}
	if p := strings.TrimSpace(os.Getenv("VAPI_VOICE_PROVIDER")); p != "" {
		id := strings.TrimSpace(os.Getenv("VAPI_VOICE_ID"))
		if id != "" {
			return Voice{ID: strings.ToLower(id), Name: id, Provider: p, VoiceID: id, Lang: "en", Change: "env VAPI_VOICE_PROVIDER + VAPI_VOICE_ID"}
		}
	}
	v, _ := VoiceByID("rachel")
	return v
}

// VoicePayload is the Vapi assistant.voice object.
func VoicePayload(v Voice) map[string]string {
	return map[string]string{"provider": v.Provider, "voiceId": v.VoiceID}
}

// Client talks to api.vapi.ai when a private key is set.
type Client struct {
	Key     string
	BaseURL string
	HTTP    *http.Client
}

// FromEnv builds a client. Empty key is valid — talk stays local.
func FromEnv() *Client {
	return &Client{
		Key:     strings.TrimSpace(os.Getenv("VAPI_API_KEY")),
		BaseURL: APIBase,
		HTTP:    &http.Client{Timeout: 20 * time.Second},
	}
}

// Bound reports a live private key.
func (c *Client) Bound() bool {
	return c != nil && strings.TrimSpace(c.Key) != ""
}

// ChatRequest is POST /chat.
type ChatRequest struct {
	AssistantID       string         `json:"assistantId,omitempty"`
	PreviousChatID    string         `json:"previousChatId,omitempty"`
	Input             string         `json:"input"`
	AssistantOverride map[string]any `json:"assistantOverrides,omitempty"`
}

// ChatResponse is the Vapi chat envelope we care about.
type ChatResponse struct {
	ID     string `json:"id"`
	Output []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"output"`
}

// CallRequest is POST /call (phone outbound).
type CallRequest struct {
	AssistantID       string         `json:"assistantId,omitempty"`
	PhoneNumberID     string         `json:"phoneNumberId,omitempty"`
	Customer          map[string]any `json:"customer,omitempty"`
	AssistantOverride map[string]any `json:"assistantOverrides,omitempty"`
}

// CallResponse is a created call.
type CallResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

func (c *Client) Chat(ctx context.Context, in ChatRequest) (ChatResponse, error) {
	var out ChatResponse
	err := c.do(ctx, http.MethodPost, "/chat", in, &out)
	return out, err
}

func (c *Client) Call(ctx context.Context, in CallRequest) (CallResponse, error) {
	var out CallResponse
	err := c.do(ctx, http.MethodPost, "/call", in, &out)
	return out, err
}

func (c *Client) PatchVoice(ctx context.Context, assistantID string, v Voice) error {
	if strings.TrimSpace(assistantID) == "" {
		return fmt.Errorf("assistant id required to patch voice")
	}
	body := map[string]any{"voice": VoicePayload(v)}
	return c.do(ctx, http.MethodPatch, "/assistant/"+assistantID, body, nil)
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	if !c.Bound() {
		return ErrUnbound
	}
	var buf bytes.Buffer
	if in != nil {
		if err := json.NewEncoder(&buf).Encode(in); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.BaseURL, "/")+path, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode >= 300 {
		return fmt.Errorf("vapi %s %s: %s", method, path, clip(string(raw), 240))
	}
	if out == nil || len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func clip(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n])
}

// EnvCard is the bind card the desk publishes.
func EnvCard() map[string]any {
	assistant := strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID"))
	share := strings.TrimSpace(os.Getenv("VAPI_SHARE_URL"))
	if share == "" && assistant != "" {
		share = DashboardURL + "assistants/" + assistant
	}
	if share == "" {
		share = ProductURL
	}
	v := DefaultVoice()
	return map[string]any{
		"product":       ProductURL,
		"dashboard":     DashboardURL,
		"docs":          DocsURL,
		"github":        GitHubURL,
		"voiceDocs":     VoiceDocs,
		"callDocs":      CallDocs,
		"chatDocs":      ChatDocs,
		"shareUrl":      share,
		"assistantId":   assistant,
		"phoneNumberId": strings.TrimSpace(os.Getenv("VAPI_PHONE_NUMBER_ID")),
		"publicKey":     strings.TrimSpace(os.Getenv("VAPI_PUBLIC_KEY")),
		"keyBound":      strings.TrimSpace(os.Getenv("VAPI_API_KEY")) != "",
		"serverUrl":     "/api/vapi/webhook",
		"voice":         v,
		"changeVoice":   []string{"POST /api/vapi/voice", "desk Voice picker", DashboardURL + " → Assistants → Voice", "env VAPI_VOICE_ID", "per-call assistantOverrides.voice"},
		"talk":          "/api/vapi/talk",
		"call":          "/api/vapi/call",
		"web":           "/api/vapi/web",
	}
}
