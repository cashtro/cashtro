// Teal is the corporate brain for Microsoft Teams AI Teal.
//
// Voice is Vapi (the "Vappy" assistant). Conversation surface is Teams only.
// Cursor Cloud Agents launch from Teams Adaptive Cards. New agentics spawn
// back into this kernel — they do not become a second product.
package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Evolu-Jeunes/Teal/internal/catalog"
	"github.com/Evolu-Jeunes/Teal/internal/kernel"
)

const (
	vapiProductURL   = "https://vapi.ai/"
	vapiDashboardURL = "https://dashboard.vapi.ai/"
	vapiDocsURL      = "https://docs.vapi.ai/"
	vapiGitHubURL    = "https://github.com/VapiAI"
	cursorAgentsURL  = "https://cursor.com/agents"
	tealTeamName     = "AI Teal"
	maxTealLog       = 80
)

// VapiCard is the voice bind the Teal brain publishes.
type VapiCard struct {
	Product   string `json:"product"`
	Dashboard string `json:"dashboard"`
	Docs      string `json:"docs"`
	GitHub    string `json:"github"`
	ShareURL  string `json:"shareUrl"`
	Assistant string `json:"assistantId,omitempty"`
	ServerURL string `json:"serverUrl"`
	KeyBound  bool   `json:"keyBound"`
	Hint      string `json:"hint"`
}

// TeamMsg is one scraped or ingested Teams message. Teams only.
type TeamMsg struct {
	ID        string `json:"id"`
	Team      string `json:"team"`
	Channel   string `json:"channel"`
	From      string `json:"from"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// Intern is one Teal desk seat (college stage / project help).
type Intern struct {
	ID     string `json:"id"`
	Seat   string `json:"seat"`
	Track  string `json:"track"`
	Status string `json:"status"`
	Mentor string `json:"mentor"`
	Help   string `json:"help,omitempty"`
}

// Idea is a generated or ingested brainstorm line.
type Idea struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Source string `json:"source"`
	Body   string `json:"body"`
}

// Question is an intern or call question the brain answered.
type Question struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	Text   string `json:"text"`
	Answer string `json:"answer"`
}

// CallListen is one Vapi end-of-call (or transcript) the brain absorbed.
type CallListen struct {
	ID         string `json:"id"`
	Assistant  string `json:"assistant,omitempty"`
	Summary    string `json:"summary"`
	Transcript string `json:"transcript,omitempty"`
}

// SpawnedIntel is a new agentic Teal minted from repeated work.
type SpawnedIntel struct {
	ID     string `json:"id"`
	Role   string `json:"role"`
	Reason string `json:"reason"`
	Mode   string `json:"mode"`
}

type tealAgent struct {
	k            *kernel.Kernel
	mu           sync.Mutex
	interns      []Intern
	ideas        []Idea
	questions    []Question
	calls        []CallListen
	spawned      []SpawnedIntel
	messages     []TeamMsg
	intelligence int
	seq          int
}

func (a *tealAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "teal", Name: "Teal", Kind: kernel.KindSystem, Mode: kernel.ModeLive,
		Role:    "corporate-brain",
		Summary: "Corporate Teal brain. Teams scrape + intern desk. Voice is the live vapi agent across every business line.",
		Capabilities: []string{
			"teal.status",
			"teal.vapi",
			"teal.scrape",
			"teal.ingest",
			"teal.idea",
			"teal.interns",
			"teal.assign",
			"teal.help",
			"teal.question",
			"teal.listen",
			"teal.spawn",
			"teal.smarter",
			"teal.cursor",
			"teal.card",
		},
		Autostart: true,
	}
}

func (a *tealAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	a.interns = []Intern{
		{ID: "intern-code", Seat: "programming stage", Track: "build", Status: "open", Mentor: "Forgeron", Help: "pair on Pandora / Proximity tickets"},
		{ID: "intern-ops", Seat: "ops stage", Track: "ops", Status: "open", Mentor: "Orfèvre", Help: "ship checklists and intern questions"},
		{ID: "intern-story", Seat: "story stage", Track: "ideas", Status: "open", Mentor: "Hustler", Help: "Pandora brainstorming and client narrative"},
	}
	seedPandoraBrain(k)
	card := vapiCard()
	k.Publish("teal", "boot", "AI Teal online · Teams desk · Vapi voice on every bridge", map[string]any{
		"team":     tealTeamName,
		"vapi":     card.ShareURL,
		"cursor":   cursorAgentsURL,
		"interns":  len(a.interns),
		"keyBound": card.KeyBound,
	})
	k.Remember("teal", "corporate brain owns Teams AI Teal ingest. Voice is the vapi kernel lane across every bridge project.")
	return nil
}

func (a *tealAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "teal.status":
		return a.status(), nil
	case "teal.vapi":
		return kernel.Result{OK: true, Message: "vapi bind", Data: vapiCard()}, nil
	case "teal.scrape":
		return a.scrape(ctx, call)
	case "teal.ingest":
		return a.ingest(call)
	case "teal.idea":
		return a.idea(call)
	case "teal.interns":
		a.mu.Lock()
		defer a.mu.Unlock()
		return kernel.Result{OK: true, Message: "intern desk", Data: append([]Intern(nil), a.interns...)}, nil
	case "teal.assign":
		return a.assign(call)
	case "teal.help":
		return a.help(call)
	case "teal.question":
		return a.question(call)
	case "teal.listen":
		return a.listen(call)
	case "teal.spawn":
		return a.spawnIntel(ctx, call)
	case "teal.smarter":
		return a.smarter(ctx)
	case "teal.cursor":
		return a.cursorLaunch(call)
	case "teal.card":
		return kernel.Result{OK: true, Message: "teams card", Data: a.teamsCard()}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func (a *tealAgent) status() kernel.Result {
	a.mu.Lock()
	defer a.mu.Unlock()
	card := vapiCard()
	return kernel.Result{OK: true, Message: "AI Teal status", Data: map[string]any{
		"team":         tealTeamName,
		"surface":      "microsoft-teams",
		"voice":        "vapi",
		"voiceLane":    "infrastructure",
		"vapi":         card,
		"cursor":       cursorAgentsURL,
		"interns":      append([]Intern(nil), a.interns...),
		"ideas":        len(a.ideas),
		"questions":    len(a.questions),
		"calls":        len(a.calls),
		"spawned":      append([]SpawnedIntel(nil), a.spawned...),
		"messages":     len(a.messages),
		"intelligence": a.intelligence,
		"teamsToken":   strings.TrimSpace(os.Getenv("TEAMS_TOKEN")) != "",
		"teamsWebhook": strings.TrimSpace(os.Getenv("TEAMS_WEBHOOK_URL")) != "",
		"pandora":      []string{"cashtro/Pandora", "Evolu-Jeunes/Pandora", "PBTM"},
	}}
}

func (a *tealAgent) scrape(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	team := payloadQuery(call, "team")
	if team == "" {
		team = tealTeamName
	}
	if !isTeamsOnly(payloadQuery(call, "channel"), payloadQuery(call, "surface")) {
		return kernel.Result{OK: false, Message: "microsoft teams only"}, nil
	}
	token := strings.TrimSpace(os.Getenv("TEAMS_TOKEN"))
	if token == "" {
		a.mu.Lock()
		local := append([]TeamMsg(nil), a.messages...)
		a.mu.Unlock()
		return kernel.Result{OK: true, Message: "teams scrape unbound · local teal thread", Data: map[string]any{
			"bound":    false,
			"team":     team,
			"hint":     "Set TEAMS_TOKEN (Graph) to pull AI Teal. Until then ingest Teams payloads via teal.ingest / POST /api/teal/ingest.",
			"messages": local,
		}}, nil
	}
	msgs, err := scrapeTeamsGraph(ctx, token, team, payloadQuery(call, "channel"))
	if err != nil {
		return kernel.Result{OK: false, Message: err.Error(), Data: map[string]any{"bound": true}}, nil
	}
	for _, m := range msgs {
		a.storeMsg(m)
		a.k.Remember("teams", m.From+": "+m.Text)
	}
	a.gain(len(msgs))
	return kernel.Result{OK: true, Message: "scraped " + strconv.Itoa(len(msgs)) + " from " + team, Data: map[string]any{
		"bound":    true,
		"team":     team,
		"messages": msgs,
	}}, nil
}

func (a *tealAgent) ingest(call kernel.Call) (kernel.Result, error) {
	if !isTeamsOnly(payloadQuery(call, "channel"), payloadQuery(call, "surface")) {
		return kernel.Result{OK: false, Message: "microsoft teams only"}, nil
	}
	text := firstNonEmpty(payloadQuery(call, "text"), payloadQuery(call, "body"), payloadQuery(call, "prompt"))
	if text == "" {
		return kernel.Result{OK: false, Message: "text required"}, nil
	}
	m := TeamMsg{
		ID:        a.nextID("msg"),
		Team:      firstNonEmpty(payloadQuery(call, "team"), tealTeamName),
		Channel:   firstNonEmpty(payloadQuery(call, "channel"), "microsoft-teams"),
		From:      firstNonEmpty(payloadQuery(call, "from"), "teal"),
		Text:      text,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	a.storeMsg(m)
	a.k.Remember("teams", m.From+": "+m.Text)
	a.k.WriteNote(kernel.Note{Agent: "teal", Source: "teams:" + m.Team, Claim: clipRunes(m.Text, 180), Quote: m.From})
	a.gain(1)
	return kernel.Result{OK: true, Message: "ingested teams message", Data: m}, nil
}

func (a *tealAgent) idea(call kernel.Call) (kernel.Result, error) {
	prompt := firstNonEmpty(payloadQuery(call, "prompt"), payloadQuery(call, "goal"), payloadQuery(call, "title"))
	if prompt == "" {
		prompt = "next Pandora move"
	}
	body := inventIdea(prompt, a.k.Recall("pandora"))
	idea := Idea{ID: a.nextID("idea"), Title: clipRunes(prompt, 80), Source: firstNonEmpty(payloadQuery(call, "source"), "teal"), Body: body}
	a.mu.Lock()
	a.ideas = append(a.ideas, idea)
	if len(a.ideas) > maxTealLog {
		a.ideas = append([]Idea(nil), a.ideas[len(a.ideas)-maxTealLog:]...)
	}
	a.mu.Unlock()
	a.k.Remember("idea", idea.Title+" · "+idea.Body)
	if cat := a.k.Catalog(); cat != nil {
		_, _ = cat.Create(catalog.CreateShip{Name: idea.Title, Notes: idea.Body, Sector: "pandora", Stack: []string{"Teams", "Vapi"}})
	}
	a.gain(1)
	return kernel.Result{OK: true, Message: "idea " + idea.ID, Data: idea}, nil
}

func (a *tealAgent) assign(call kernel.Call) (kernel.Result, error) {
	id := firstNonEmpty(payloadQuery(call, "id"), payloadQuery(call, "intern"))
	help := firstNonEmpty(payloadQuery(call, "help"), payloadQuery(call, "task"), payloadQuery(call, "prompt"))
	if id == "" || help == "" {
		return kernel.Result{OK: false, Message: "intern id and help/task required"}, nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := range a.interns {
		if a.interns[i].ID == id || strings.EqualFold(a.interns[i].Seat, id) || strings.EqualFold(a.interns[i].Track, id) {
			a.interns[i].Help = help
			a.interns[i].Status = "assigned"
			a.k.Remember("intern", a.interns[i].ID+" · "+help)
			_, _ = a.k.Post("teal", "planner", "intern", a.interns[i].ID+": "+help)
			return kernel.Result{OK: true, Message: "assigned " + a.interns[i].ID, Data: a.interns[i]}, nil
		}
	}
	return kernel.Result{OK: false, Message: "unknown intern: " + id}, nil
}

func (a *tealAgent) help(call kernel.Call) (kernel.Result, error) {
	ask := firstNonEmpty(payloadQuery(call, "prompt"), payloadQuery(call, "help"), payloadQuery(call, "task"))
	if ask == "" {
		ask = "how do I start"
	}
	answer := helpIntern(ask, a.k.Recall("pandora"))
	q := Question{ID: a.nextID("q"), From: firstNonEmpty(payloadQuery(call, "from"), "intern"), Text: ask, Answer: answer}
	a.mu.Lock()
	a.questions = append(a.questions, q)
	a.mu.Unlock()
	a.k.Remember("help", ask+" → "+answer)
	a.gain(1)
	return kernel.Result{OK: true, Message: "helped", Data: q}, nil
}

func (a *tealAgent) question(call kernel.Call) (kernel.Result, error) {
	return a.help(call)
}

func (a *tealAgent) listen(call kernel.Call) (kernel.Result, error) {
	parsed := parseVapiListen(call.Payload)
	if parsed.Transcript == "" && parsed.Summary == "" {
		parsed.Summary = firstNonEmpty(payloadQuery(call, "summary"), payloadQuery(call, "transcript"), payloadQuery(call, "prompt"))
		parsed.Transcript = payloadQuery(call, "transcript")
		parsed.Assistant = payloadQuery(call, "assistantId")
		parsed.ID = payloadQuery(call, "id")
	}
	if parsed.Summary == "" && parsed.Transcript == "" {
		return kernel.Result{OK: false, Message: "call transcript or summary required"}, nil
	}
	if parsed.ID == "" {
		parsed.ID = a.nextID("call")
	}
	if parsed.Summary == "" {
		parsed.Summary = clipRunes(parsed.Transcript, 240)
	}
	a.mu.Lock()
	a.calls = append(a.calls, parsed)
	if len(a.calls) > maxTealLog {
		a.calls = append([]CallListen(nil), a.calls[len(a.calls)-maxTealLog:]...)
	}
	a.mu.Unlock()
	blob := firstNonEmpty(parsed.Summary, parsed.Transcript)
	a.k.Remember("call", blob)
	a.k.WriteNote(kernel.Note{Agent: "teal", Source: "vapi", URL: vapiCard().ShareURL, Claim: clipRunes(blob, 180), Quote: parsed.Assistant})
	if looksLikeQuestion(blob) {
		_, _ = a.help(kernel.Call{Capability: "teal.help", Payload: mustJSON(map[string]string{"prompt": blob, "from": "call"})})
	}
	a.gain(2)
	return kernel.Result{OK: true, Message: "listened " + parsed.ID, Data: parsed}, nil
}

func (a *tealAgent) spawnIntel(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	role := firstNonEmpty(payloadQuery(call, "role"), payloadQuery(call, "name"), nextSpawnRole(a.spawned))
	reason := firstNonEmpty(payloadQuery(call, "reason"), payloadQuery(call, "prompt"), "pattern repeated on Teal desk")
	id := slugID(role)
	if id == "" {
		id = a.nextID("intel")
	}
	if _, err := a.k.Process(id); err == nil {
		id = id + "-" + strconv.Itoa(len(a.spawned)+1)
	}
	spec := kernel.Spec{
		ID: id, Name: titleID(id),
		Kind: kernel.KindUser, Mode: kernel.ModeResident, Role: role,
		Summary:      "Spawned by Teal from Teams/Vapi work: " + reason,
		Capabilities: []string{id + ".run"},
		Autostart:    true,
	}
	a.k.Register(resident(spec, nil))
	if err := a.k.Spawn(ctx, id); err != nil {
		return kernel.Result{}, err
	}
	got := SpawnedIntel{ID: id, Role: role, Reason: reason, Mode: string(kernel.ModeResident)}
	a.mu.Lock()
	a.spawned = append(a.spawned, got)
	a.mu.Unlock()
	a.k.Remember("spawn", id+" · "+reason)
	_, _ = a.k.Post("teal", "manager", "spawn", id+": "+reason)
	return kernel.Result{OK: true, Message: "spawned " + id, Data: got}, nil
}

func (a *tealAgent) smarter(ctx context.Context) (kernel.Result, error) {
	a.mu.Lock()
	score := a.intelligence + len(a.messages) + len(a.calls)*2 + len(a.questions)
	shouldSpawn := score >= 3 && len(a.spawned) < 8
	a.mu.Unlock()
	var spawned any
	if shouldSpawn {
		res, err := a.spawnIntel(ctx, kernel.Call{Capability: "teal.spawn", Payload: mustJSON(map[string]string{
			"reason": "teal.smarter · repeated Teams/Vapi work",
		})})
		if err != nil {
			return kernel.Result{}, err
		}
		spawned = res.Data
	}
	a.gain(1)
	a.k.Remember("smarter", "intelligence tick")
	return kernel.Result{OK: true, Message: "smarter", Data: map[string]any{
		"intelligence": a.intelligence,
		"spawned":      spawned,
		"score":        score,
	}}, nil
}

func (a *tealAgent) cursorLaunch(call kernel.Call) (kernel.Result, error) {
	prompt := firstNonEmpty(payloadQuery(call, "prompt"), payloadQuery(call, "goal"), "Teal brain · help interns on Pandora")
	link := cursorAgentsURL + "?q=" + url.QueryEscape(prompt)
	body := "Cursor Cloud Agent from Teams AI Teal: " + prompt + " · " + link
	c := a.k.RequestConfirm("teal", "teal.cursor", body)
	card := a.teamsCard()
	card["cursorPrompt"] = prompt
	card["cursorUrl"] = link
	return kernel.Result{OK: true, Message: "parked cursor launch #" + strconv.Itoa(c.ID), Data: map[string]any{
		"confirm": c,
		"card":    card,
		"url":     link,
	}}, nil
}

func (a *tealAgent) teamsCard() map[string]any {
	card := vapiCard()
	a.mu.Lock()
	intel := a.intelligence
	internN := len(a.interns)
	ideaN := len(a.ideas)
	a.mu.Unlock()
	return map[string]any{
		"type":    "AdaptiveCard",
		"version": "1.5",
		"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
		"body": []map[string]any{
			{"type": "TextBlock", "weight": "Bolder", "size": "Large", "text": "AI Teal · corporate brain"},
			{"type": "TextBlock", "wrap": true, "text": "Teams only. Vapi voice. Pandora memory. Intern desk. Cursor app."},
			{"type": "FactSet", "facts": []map[string]string{
				{"title": "Voice", "value": card.ShareURL},
				{"title": "Cursor", "value": cursorAgentsURL},
				{"title": "Interns", "value": strconv.Itoa(internN)},
				{"title": "Ideas", "value": strconv.Itoa(ideaN)},
				{"title": "Intel", "value": strconv.Itoa(intel)},
			}},
		},
		"actions": []map[string]any{
			{"type": "Action.OpenUrl", "title": "Talk on Vapi", "url": card.ShareURL},
			{"type": "Action.OpenUrl", "title": "Open Cursor", "url": cursorAgentsURL},
			{"type": "Action.OpenUrl", "title": "Vapi dashboard", "url": card.Dashboard},
		},
		"channel": "microsoft-teams",
		"team":    tealTeamName,
	}
}

func (a *tealAgent) storeMsg(m TeamMsg) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = append(a.messages, m)
	if len(a.messages) > maxTealLog {
		a.messages = append([]TeamMsg(nil), a.messages[len(a.messages)-maxTealLog:]...)
	}
}

func (a *tealAgent) nextID(prefix string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.seq++
	return fmt.Sprintf("%s-%d", prefix, a.seq)
}

func (a *tealAgent) gain(n int) {
	a.mu.Lock()
	a.intelligence += n
	a.mu.Unlock()
}

func vapiCard() VapiCard {
	assistant := strings.TrimSpace(os.Getenv("VAPI_ASSISTANT_ID"))
	share := strings.TrimSpace(os.Getenv("VAPI_SHARE_URL"))
	if share == "" && assistant != "" {
		share = vapiDashboardURL + "assistants/" + assistant
	}
	if share == "" {
		share = vapiProductURL
	}
	key := strings.TrimSpace(os.Getenv("VAPI_API_KEY")) != ""
	hint := "Vapi product link bound. Paste the AI Teal assistant share into VAPI_SHARE_URL (or VAPI_ASSISTANT_ID) to lock the exact voice."
	if key {
		hint = "Vapi key present. Server URL should POST /api/vapi/webhook (end-of-call-report + transcript)."
	}
	if assistant != "" {
		hint = "Vapi assistant " + assistant + " · " + share
	}
	return VapiCard{
		Product:   vapiProductURL,
		Dashboard: vapiDashboardURL,
		Docs:      vapiDocsURL,
		GitHub:    vapiGitHubURL,
		ShareURL:  share,
		Assistant: assistant,
		ServerURL: "/api/vapi/webhook",
		KeyBound:  key,
		Hint:      hint,
	}
}

func seedPandoraBrain(k *kernel.Kernel) {
	notes := []kernel.Note{
		{Agent: "teal", Source: "pandora", URL: "https://github.com/cashtro/Pandora", Claim: "Pandora is PBTM — Pandora Business Technology and Marketing — under Le Hustler.", Quote: "cashtro/Pandora + Evolu-Jeunes/Pandora. Empire line. Stripe on PBTM."},
		{Agent: "teal", Source: "pandora", URL: "https://github.com/Evolu-Jeunes/Pandora", Claim: "Pandora Brains is the product concept: a map of brains, not a second OS.", Quote: "Castro's Pandora Brains profile fits the brain concept."},
		{Agent: "teal", Source: "pandora", URL: "https://github.com/Evolu-Jeunes/Panda", Claim: "Panda white-glove and Pandora share the Hustler lane.", Quote: "Hustler sells. Architecte designs. Forgeron installs."},
		{Agent: "teal", Source: "teal", URL: cursorAgentsURL, Claim: "Corporate Teal manages interns from Teams. Vapi voice runs on Cashtro infrastructure across every bridge project.", Quote: "Ideas, intern seats, call listen, project help, self-spawned intel, Cursor cards."},
		{Agent: "teal", Source: "vapi", URL: vapiProductURL, Claim: "Vapi (Vappy) is the in-repo voice lane: talk, outbound calls, changeable voices, every business line.", Quote: "dashboard.vapi.ai · POST /api/vapi/talk · POST /api/vapi/call"},
	}
	for _, n := range notes {
		k.WriteNote(n)
		k.Remember("pandora", n.Claim)
	}
}

func inventIdea(prompt string, facts []kernel.Fact) string {
	var bits []string
	for i, f := range facts {
		if i == 4 {
			break
		}
		bits = append(bits, f.Text)
	}
	if len(bits) == 0 {
		return "Turn " + prompt + " into a Teams card + Vapi listen + intern ticket on the Pandora line."
	}
	return "From Pandora memory — " + prompt + ": " + strings.Join(bits, " · ")
}

func helpIntern(ask string, facts []kernel.Fact) string {
	lower := strings.ToLower(ask)
	switch {
	case strings.Contains(lower, "vapi") || strings.Contains(lower, "voice") || strings.Contains(lower, "call"):
		return "Voice is Vapi on Cashtro infrastructure: POST /api/vapi/talk, change voice POST /api/vapi/voice, outbound POST /api/vapi/call (human gate). Webhook /api/vapi/webhook still files transcripts into Teal. Dashboard " + vapiCard().ShareURL
	case strings.Contains(lower, "cursor"):
		return "Open the Cursor app from the Teams card: " + cursorAgentsURL + ". Park a launch with teal.cursor so a human confirms."
	case strings.Contains(lower, "intern") || strings.Contains(lower, "stage"):
		return "Intern desk is teal.interns. Assign with teal.assign. Help stays on Teams — no Slack, no mail."
	case strings.Contains(lower, "pandora") || strings.Contains(lower, "pbtm"):
		return "Pandora = PBTM under Hustler. Repos: cashtro/Pandora, Evolu-Jeunes/Pandora. Brainstorm lives in teal memory."
	}
	if len(facts) > 0 {
		return "From the brain: " + facts[0].Text + " — then file a Teams question so Teal gets smarter."
	}
	return "Ask on Teams AI Teal. Teal will file it, help the intern, and spawn intel if the pattern repeats."
}

func looksLikeQuestion(s string) bool {
	s = strings.ToLower(s)
	return strings.Contains(s, "?") || strings.Contains(s, "how ") || strings.Contains(s, "comment ") || strings.Contains(s, "help")
}

func nextSpawnRole(existing []SpawnedIntel) string {
	roles := []string{"intern-coach", "call-scribe", "idea-weaver", "project-steward", "question-router"}
	used := map[string]bool{}
	for _, s := range existing {
		used[s.Role] = true
	}
	for _, r := range roles {
		if !used[r] {
			return r
		}
	}
	return "teal-intel"
}

func slugID(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dash = false
			continue
		}
		if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func isTeamsOnly(channel, surface string) bool {
	blob := strings.ToLower(channel + " " + surface)
	if strings.Contains(blob, "slack") || strings.Contains(blob, "mail") || strings.Contains(blob, "discord") || strings.Contains(blob, "whatsapp") {
		return false
	}
	if strings.TrimSpace(blob) == "" {
		return true
	}
	return strings.Contains(blob, "team") || strings.Contains(blob, "msteams") || strings.Contains(blob, "teal")
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func clipRunes(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n])
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func titleID(id string) string {
	parts := strings.Split(id, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

// scrapeTeamsGraph pulls recent channel messages for a team named like "AI Teal".
func scrapeTeamsGraph(ctx context.Context, token, teamName, channelName string) ([]TeamMsg, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	teams, err := graphJSON(ctx, client, token, "https://graph.microsoft.com/v1.0/me/joinedTeams?$select=id,displayName")
	if err != nil {
		return nil, err
	}
	teamID, display := pickGraphName(teams, teamName, tealTeamName)
	if teamID == "" {
		return nil, fmt.Errorf("teams team %q not found", teamName)
	}
	if channelName == "" {
		channelName = "General"
	}
	chs, err := graphJSON(ctx, client, token, "https://graph.microsoft.com/v1.0/teams/"+url.PathEscape(teamID)+"/channels?$select=id,displayName")
	if err != nil {
		return nil, err
	}
	chID, chName := pickGraphName(chs, channelName, "General")
	if chID == "" {
		return nil, fmt.Errorf("teams channel %q not found", channelName)
	}
	msgs, err := graphJSON(ctx, client, token, "https://graph.microsoft.com/v1.0/teams/"+url.PathEscape(teamID)+"/channels/"+url.PathEscape(chID)+"/messages?$top=20")
	if err != nil {
		return nil, err
	}
	out := make([]TeamMsg, 0)
	for _, raw := range graphValues(msgs) {
		text := graphText(raw)
		if text == "" {
			continue
		}
		from, _ := raw["from"].(map[string]any)
		user, _ := from["user"].(map[string]any)
		name, _ := user["displayName"].(string)
		id, _ := raw["id"].(string)
		at, _ := raw["createdDateTime"].(string)
		out = append(out, TeamMsg{
			ID: id, Team: display, Channel: chName, From: name, Text: text, CreatedAt: at,
		})
	}
	return out, nil
}

func graphJSON(ctx context.Context, client *http.Client, token, rawURL string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("graph %s: %s", res.Status, clipRunes(string(body), 180))
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func graphValues(doc map[string]any) []map[string]any {
	raw, _ := doc["value"].([]any)
	out := make([]map[string]any, 0, len(raw))
	for _, v := range raw {
		if m, ok := v.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func pickGraphName(doc map[string]any, want, fallback string) (id, name string) {
	want = strings.ToLower(strings.TrimSpace(want))
	fb := strings.ToLower(fallback)
	var firstID, firstName string
	for _, m := range graphValues(doc) {
		n, _ := m["displayName"].(string)
		i, _ := m["id"].(string)
		if firstID == "" {
			firstID, firstName = i, n
		}
		ln := strings.ToLower(n)
		if ln == want || strings.Contains(ln, want) || strings.Contains(ln, "teal") || ln == fb {
			return i, n
		}
	}
	return firstID, firstName
}

func graphText(raw map[string]any) string {
	if s, ok := raw["body"].(map[string]any); ok {
		if t, ok := s["content"].(string); ok {
			return stripTags(t)
		}
	}
	if s, ok := raw["text"].(string); ok {
		return s
	}
	return ""
}

func stripTags(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		switch {
		case r == '<':
			in = true
		case r == '>':
			in = false
		case !in:
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(htmlUnescape(b.String()))
}

func htmlUnescape(s string) string {
	repl := []string{"&nbsp;", " ", "&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'"}
	return strings.NewReplacer(repl...).Replace(s)
}

func parseVapiListen(raw json.RawMessage) CallListen {
	var out CallListen
	if len(raw) == 0 {
		return out
	}
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		return out
	}
	msg, _ := env["message"].(map[string]any)
	if msg == nil {
		msg = env
	}
	if call, ok := msg["call"].(map[string]any); ok {
		out.ID, _ = call["id"].(string)
		out.Assistant, _ = call["assistantId"].(string)
	}
	if art, ok := msg["artifact"].(map[string]any); ok {
		out.Transcript, _ = art["transcript"].(string)
	}
	if an, ok := msg["analysis"].(map[string]any); ok {
		out.Summary, _ = an["summary"].(string)
	}
	if out.Transcript == "" {
		out.Transcript, _ = msg["transcript"].(string)
	}
	if out.Summary == "" {
		out.Summary, _ = msg["summary"].(string)
	}
	if out.Assistant == "" {
		out.Assistant, _ = msg["assistantId"].(string)
	}
	if out.ID == "" {
		out.ID, _ = msg["id"].(string)
	}
	return out
}

// postTeamsWebhook is used only after a human confirm when TEAMS_WEBHOOK_URL is set.
func postTeamsWebhook(ctx context.Context, payload any) error {
	hook := strings.TrimSpace(os.Getenv("TEAMS_WEBHOOK_URL"))
	if hook == "" {
		return fmt.Errorf("TEAMS_WEBHOOK_URL unset")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hook, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("teams webhook %s: %s", res.Status, b)
	}
	return nil
}
