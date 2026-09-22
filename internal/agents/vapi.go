package agents

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/vapi"
)

// vapiAgent is the live voice lane. Talk, outbound calls, voice switch,
// and cross-bridge context for every business line. Not a second OS.
type vapiAgent struct {
	k       *kernel.Kernel
	client  *vapi.Client
	mu      sync.Mutex
	voice   vapi.Voice
	line    string
	chatID  string
	turns   []TalkTurn
	pending map[int]PendingCall
	seq     int
}

// TalkTurn is one desk / chat exchange with Castro.
type TalkTurn struct {
	ID     string `json:"id"`
	Line   string `json:"line,omitempty"`
	Voice  string `json:"voice"`
	User   string `json:"user"`
	Reply  string `json:"reply"`
	Live   bool   `json:"live"`
	ChatID string `json:"chatId,omitempty"`
}

// PendingCall is an outbound plan waiting on the human gate.
type PendingCall struct {
	ConfirmID int                `json:"confirmId"`
	To        string             `json:"to"`
	Line      string             `json:"line"`
	Voice     vapi.Voice         `json:"voice"`
	Prompt    string             `json:"prompt"`
	Live      *vapi.CallResponse `json:"live,omitempty"`
}

func (a *vapiAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "vapi", Name: "Vapi", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role:    "voice",
		Summary: "Vapi voice lane across every bridge project. Talk, automated outbound calls, changeable voices.",
		Capabilities: []string{
			"vapi.status",
			"vapi.voices",
			"vapi.voice",
			"vapi.talk",
			"vapi.call",
			"vapi.fire",
			"vapi.web",
			"vapi.bridge",
		},
		Autostart: true,
	}
}

func (a *vapiAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	a.client = vapi.FromEnv()
	a.voice = vapi.DefaultVoice()
	a.line = "pandora"
	a.pending = map[int]PendingCall{}
	k.Publish("vapi", "boot", "Vapi lane online · talk + calls · voice "+a.voice.Name, map[string]any{
		"voice":    a.voice,
		"line":     a.line,
		"keyBound": a.client.Bound(),
		"bridges":  lineIDs(),
	})
	k.Remember("vapi", "voice lane bound to infrastructure. Change voice with vapi.voice. Talk with vapi.talk. Calls park a confirm.")
	return nil
}

func (a *vapiAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "vapi.status":
		return a.status(), nil
	case "vapi.voices":
		return kernel.Result{OK: true, Message: "voices", Data: map[string]any{
			"current": a.currentVoice(),
			"voices":  vapi.Voices(),
			"change":  vapi.EnvCard()["changeVoice"],
		}}, nil
	case "vapi.voice":
		return a.setVoice(ctx, call)
	case "vapi.talk":
		return a.talk(ctx, call)
	case "vapi.call":
		return a.planCall(call)
	case "vapi.fire":
		return a.fire(ctx, call)
	case "vapi.web":
		return a.web(), nil
	case "vapi.bridge":
		return a.setBridge(call)
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func (a *vapiAgent) status() kernel.Result {
	a.mu.Lock()
	defer a.mu.Unlock()
	card := vapi.EnvCard()
	card["currentVoice"] = a.voice
	card["line"] = a.line
	card["turns"] = len(a.turns)
	card["pendingCalls"] = len(a.pending)
	card["bridges"] = lineIDs()
	card["keyBound"] = a.client.Bound()
	card["hint"] = a.hint()
	return kernel.Result{OK: true, Message: "vapi status", Data: card}
}

func (a *vapiAgent) hint() string {
	if a.client.Bound() {
		return "Live key bound. Talk hits api.vapi.ai/chat. Outbound /call still waits on the human gate. Voice patches the assistant when VAPI_ASSISTANT_ID is set."
	}
	return "No VAPI_API_KEY — talk stays local on the desk and still uses Pandora + line memory. Bind a key to automate live calls. Change voice anytime."
}

func (a *vapiAgent) currentVoice() vapi.Voice {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.voice
}

func (a *vapiAgent) setVoice(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	id := firstNonEmpty(payloadQuery(call, "id"), payloadQuery(call, "voice"), payloadQuery(call, "voiceId"), payloadQuery(call, "prompt"))
	v, ok := vapi.VoiceByID(id)
	if !ok {
		return kernel.Result{OK: false, Message: "unknown voice: " + id + " · GET /api/vapi/voices"}, nil
	}
	a.mu.Lock()
	a.voice = v
	a.mu.Unlock()
	a.k.Remember("vapi", "voice → "+v.Name+" ("+v.Provider+"/"+v.VoiceID+")")
	patched := false
	assistant := strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID"))
	if a.client.Bound() && assistant != "" {
		if err := a.client.PatchVoice(ctx, assistant, v); err != nil && !errors.Is(err, vapi.ErrUnbound) {
			return kernel.Result{OK: false, Message: err.Error(), Data: v}, nil
		} else if err == nil {
			patched = true
		}
	}
	return kernel.Result{OK: true, Message: "voice " + v.Name, Data: map[string]any{
		"voice":   v,
		"patched": patched,
		"change":  v.Change,
	}}, nil
}

func (a *vapiAgent) setBridge(call kernel.Call) (kernel.Result, error) {
	id := firstNonEmpty(payloadQuery(call, "line"), payloadQuery(call, "id"), payloadQuery(call, "project"), payloadQuery(call, "prompt"))
	ln, ok := LineByID(id)
	if !ok {
		return kernel.Result{OK: false, Message: "unknown line: " + id, Data: lineIDs()}, nil
	}
	a.mu.Lock()
	a.line = ln.ID
	a.mu.Unlock()
	a.k.Remember("vapi", "bridge → "+ln.ID)
	_, _ = a.k.Post("vapi", "manager", "bridge", ln.ID)
	return kernel.Result{OK: true, Message: "bridged " + ln.Name, Data: map[string]any{
		"line":    ln,
		"bridges": BridgesFor(ln.ID),
	}}, nil
}

func (a *vapiAgent) talk(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	text := firstNonEmpty(payloadQuery(call, "text"), payloadQuery(call, "input"), payloadQuery(call, "prompt"))
	if text == "" {
		return kernel.Result{OK: false, Message: "text required"}, nil
	}
	if line := payloadQuery(call, "line"); line != "" {
		if ln, ok := LineByID(line); ok {
			a.mu.Lock()
			a.line = ln.ID
			a.mu.Unlock()
		}
	}
	a.mu.Lock()
	lineID := a.line
	voice := a.voice
	prev := a.chatID
	a.mu.Unlock()
	script := BridgeScript(lineID, a.k.Recall("pandora"), a.k.Recall("vapi"))
	reply, chatID, live := "", prev, false
	if a.client.Bound() && strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID")) != "" {
		out, err := a.client.Chat(ctx, vapi.ChatRequest{
			AssistantID:    os.Getenv("VAPI_ASSISTANT_ID"),
			PreviousChatID: prev,
			Input:          text,
			AssistantOverride: map[string]any{
				"voice": vapi.VoicePayload(voice),
				"variableValues": map[string]string{
					"companyName": "Cashtro / Pandora PBTM",
					"line":        lineID,
					"script":      script,
				},
			},
		})
		if err == nil && len(out.Output) > 0 {
			reply = out.Output[0].Content
			chatID = out.ID
			live = true
		} else if err != nil && !errors.Is(err, vapi.ErrUnbound) {
			reply = "Vapi chat error: " + err.Error() + " — falling back to local brain. " + localTalk(text, script)
		}
	}
	if reply == "" {
		reply = localTalk(text, script)
	}
	turn := TalkTurn{ID: a.nextID("talk"), Line: lineID, Voice: voice.ID, User: text, Reply: reply, Live: live, ChatID: chatID}
	a.mu.Lock()
	a.turns = append(a.turns, turn)
	if len(a.turns) > maxTealLog {
		a.turns = append([]TalkTurn(nil), a.turns[len(a.turns)-maxTealLog:]...)
	}
	a.chatID = chatID
	a.mu.Unlock()
	a.k.Remember("talk", text+" → "+clipRunes(reply, 160))
	_, _ = a.k.Post("vapi", "teal", "talk", clipRunes(text, 80))
	return kernel.Result{OK: true, Message: "talk " + turn.ID, Data: turn}, nil
}

func (a *vapiAgent) planCall(call kernel.Call) (kernel.Result, error) {
	to := firstNonEmpty(payloadQuery(call, "to"), payloadQuery(call, "number"), payloadQuery(call, "customer"))
	prompt := firstNonEmpty(payloadQuery(call, "prompt"), payloadQuery(call, "text"), "automated call from Cashtro Teal")
	if to == "" {
		return kernel.Result{OK: false, Message: "customer number required (to)"}, nil
	}
	a.mu.Lock()
	lineID := firstNonEmpty(payloadQuery(call, "line"), a.line)
	voice := a.voice
	a.mu.Unlock()
	body := "Vapi outbound " + to + " · line " + lineID + " · voice " + voice.Name + " · " + prompt
	c := a.k.RequestConfirm("vapi", "vapi.call", body)
	plan := PendingCall{ConfirmID: c.ID, To: to, Line: lineID, Voice: voice, Prompt: prompt}
	a.mu.Lock()
	a.pending[c.ID] = plan
	a.mu.Unlock()
	a.k.Remember("call", body)
	return kernel.Result{OK: true, Message: "parked outbound #" + fmt.Sprintf("%d", c.ID), Data: map[string]any{
		"confirm": c,
		"plan":    plan,
		"hint":    "Allow the human gate, then vapi.fire (or POST /api/confirms/{id}/allow). Live dial needs VAPI_API_KEY + VAPI_PHONE_NUMBER_ID. Free Vapi numbers cannot outbound.",
	}}, nil
}

func (a *vapiAgent) fire(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	id := payloadInt(call, "id")
	if id == 0 {
		id = payloadInt(call, "confirmId")
	}
	a.mu.Lock()
	plan, ok := a.pending[id]
	a.mu.Unlock()
	if !ok {
		return kernel.Result{OK: false, Message: "no pending call #" + fmt.Sprintf("%d", id)}, nil
	}
	allowed := false
	for _, c := range a.k.Confirms() {
		if c.ID == id && c.Status == "allowed" {
			allowed = true
			break
		}
	}
	if !allowed {
		return kernel.Result{OK: false, Message: "confirm #" + fmt.Sprintf("%d", id) + " not allowed yet"}, nil
	}
	phoneID := strings.TrimSpace(os.Getenv("VAPI_PHONE_NUMBER_ID"))
	assistant := strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID"))
	if !a.client.Bound() || phoneID == "" || assistant == "" {
		return kernel.Result{OK: true, Message: "queued · bind VAPI_API_KEY + VAPI_PHONE_NUMBER_ID + VAPI_ASSISTANT_ID to dial", Data: plan}, nil
	}
	live, err := a.client.Call(ctx, vapi.CallRequest{
		AssistantID:   assistant,
		PhoneNumberID: phoneID,
		Customer:      map[string]any{"number": plan.To},
		AssistantOverride: map[string]any{
			"voice":          vapi.VoicePayload(plan.Voice),
			"firstMessage":   plan.Prompt,
			"variableValues": map[string]string{"line": plan.Line, "script": BridgeScript(plan.Line, a.k.Recall("pandora"), nil)},
		},
	})
	if err != nil {
		return kernel.Result{OK: false, Message: err.Error(), Data: plan}, nil
	}
	plan.Live = &live
	a.mu.Lock()
	a.pending[id] = plan
	a.mu.Unlock()
	a.k.Remember("call", "dialed "+plan.To+" · "+live.ID)
	return kernel.Result{OK: true, Message: "dialed " + live.ID, Data: plan}, nil
}

func (a *vapiAgent) web() kernel.Result {
	a.mu.Lock()
	voice := a.voice
	lineID := a.line
	a.mu.Unlock()
	assistant := strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID"))
	pub := strings.TrimSpace(os.Getenv("VAPI_PUBLIC_KEY"))
	talkURL := vapi.DashboardURL
	if assistant != "" {
		talkURL = vapi.DashboardURL + "assistants/" + assistant
	}
	return kernel.Result{OK: true, Message: "web talk", Data: map[string]any{
		"url":         talkURL,
		"publicKey":   pub,
		"assistantId": assistant,
		"voice":       voice,
		"line":        lineID,
		"script":      "@vapi-ai/web",
		"hint":        "Desk talk is POST /api/vapi/talk. Browser widget needs VAPI_PUBLIC_KEY. Dashboard Talk button: " + talkURL,
	}}
}

func (a *vapiAgent) nextID(prefix string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.seq++
	return fmt.Sprintf("%s-%d", prefix, a.seq)
}

func lineIDs() []string {
	out := make([]string, 0, len(Lines()))
	for _, ln := range Lines() {
		out = append(out, ln.ID)
	}
	return out
}

// BridgesFor returns the ecosystem hops that touch a line.
func BridgesFor(line string) []Link {
	out := make([]Link, 0)
	for _, ln := range Links() {
		if ln.From == line || ln.To == line {
			out = append(out, ln)
		}
	}
	return out
}

// BridgeScript is the voice system context for a line + memory.
func BridgeScript(lineID string, facts ...[]kernel.Fact) string {
	ln, ok := LineByID(lineID)
	if !ok {
		ln = Line{ID: lineID, Name: lineID, Mandate: "Cashtro OS"}
	}
	var b strings.Builder
	b.WriteString("You are the Cashtro / Teal voice on Vapi. Talk with Castro. FR/EN. ")
	b.WriteString("Current line: " + ln.Name + " (" + ln.ID + "). Brain: " + ln.Brain + ". " + ln.Mandate + " ")
	b.WriteString("Repos: " + strings.Join(ln.Repos, ", ") + ". ")
	for _, hop := range BridgesFor(ln.ID) {
		b.WriteString("Bridge " + hop.From + " → " + hop.To + " via " + hop.Via + ". ")
	}
	n := 0
	for _, pack := range facts {
		for _, f := range pack {
			if n == 6 {
				break
			}
			b.WriteString(f.Text + " ")
			n++
		}
	}
	b.WriteString("Keep answers short. Help interns, park ideas, never send outbound calls without a human confirm.")
	return strings.TrimSpace(b.String())
}

func localTalk(text, script string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "voice") || strings.Contains(lower, "voix"):
		return "Voice picker is GET /api/vapi/voices. Switch with POST /api/vapi/voice {\"id\":\"denise\"} (FR) or rachel/adam/nova. Same change lives in the Vapi dashboard under Assistants → Voice, or env VAPI_VOICE_ID."
	case strings.Contains(lower, "call") || strings.Contains(lower, "appel"):
		return "Outbound is POST /api/vapi/call {\"to\":\"+1...\"}. It parks a confirm. Allow it, then it dials when VAPI_API_KEY and VAPI_PHONE_NUMBER_ID are set. Free Vapi numbers cannot outbound."
	case strings.Contains(lower, "bridge") || strings.Contains(lower, "project"):
		return "Bridges: " + strings.Join(lineIDs(), ", ") + ". POST /api/vapi/bridge {\"line\":\"proximity\"} and I talk in that project's context."
	}
	return "Local Teal/Vapi (no live key yet). " + clipRunes(script, 280) + " You said: " + clipRunes(text, 80)
}
