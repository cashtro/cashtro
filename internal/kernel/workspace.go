package kernel

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	maxMail        = 64
	maxNotes       = 200
	maxFacts       = 200
	maxTeamThreads = 32
	maxTeamMsgs    = 200

	// TeamChannelMicrosoft is the only conversation channel this OS speaks.
	TeamChannelMicrosoft = "microsoft-teams"
	// TeamKindConversation is the only Teams surface this OS speaks.
	TeamKindConversation = "conversation"
	// TeamTenantProximity is the only tenant the dock is scoped to.
	TeamTenantProximity = "Proximity"
)

var (
	// ErrTeamChannel is returned when the channel is not Microsoft Teams.
	ErrTeamChannel = errors.New("microsoft teams only")
	// ErrTeamKind is returned when the ask is not a conversation.
	ErrTeamKind = errors.New("conversation only")
	// ErrTeamTenant is returned when the tenant is not Proximity.
	ErrTeamTenant = errors.New("proximity tenant only")
	// ErrTeamThread is returned when a thread id is unknown.
	ErrTeamThread = errors.New("thread not found")
	// ErrTeamEmpty is returned when a say has no body.
	ErrTeamEmpty = errors.New("empty message")
)

// Mail is one async message between agentics.
type Mail struct {
	ID   int    `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Kind string `json:"kind"`
	Body string `json:"body"`
}

// Note is a sourced research record. Findings do not live in chat.
type Note struct {
	ID     int    `json:"id"`
	Agent  string `json:"agent"`
	Source string `json:"source"`
	URL    string `json:"url,omitempty"`
	Claim  string `json:"claim"`
	Quote  string `json:"quote,omitempty"`
}

// Fact is one episodic memory line.
type Fact struct {
	ID    int    `json:"id"`
	Topic string `json:"topic"`
	Text  string `json:"text"`
}

// Confirm is a human gate before outbound work.
type Confirm struct {
	ID      int    `json:"id"`
	Agent   string `json:"agent"`
	Cap     string `json:"capability"`
	Body    string `json:"body"`
	Status  string `json:"status"`
	Outcome string `json:"outcome,omitempty"`
}

// TeamThread is one Proximity Microsoft Teams conversation.
type TeamThread struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Tenant    string        `json:"tenant"`
	Channel   string        `json:"channel"`
	Kind      string        `json:"kind"`
	Dock      string        `json:"dock"`
	Messages  []TeamMessage `json:"messages"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// TeamMessage is one turn in a Teams conversation.
type TeamMessage struct {
	ID   int       `json:"id"`
	Role string    `json:"role"`
	Body string    `json:"body"`
	At   time.Time `json:"at"`
}

// TeamCard is the public Teams conversation bind card.
type TeamCard struct {
	Channel string `json:"channel"`
	Kind    string `json:"kind"`
	Tenant  string `json:"tenant"`
	Dock    string `json:"dock"`
	Hint    string `json:"hint"`
	Threads int    `json:"threads"`
}

// TeamOpen is the input for opening a conversation thread.
type TeamOpen struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Tenant  string `json:"tenant"`
	Channel string `json:"channel"`
	Kind    string `json:"kind"`
}

// InvokeCap routes a verb without the caller naming the owner.
func (k *Kernel) InvokeCap(ctx context.Context, cap string, call Call) (Result, error) {
	k.mu.RLock()
	c, ok := k.caps[cap]
	k.mu.RUnlock()
	if !ok {
		return Result{}, fmt.Errorf("%w: %s", ErrUnknownCapability, cap)
	}
	call.Capability = cap
	return k.Invoke(ctx, c.Agent, call)
}

// Post drops mail on an agent's desk.
func (k *Kernel) Post(from, to, kind, body string) (Mail, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, ok := k.procs[to]; !ok {
		return Mail{}, fmt.Errorf("%w: %s", ErrUnknownAgent, to)
	}
	k.mailSeq++
	m := Mail{ID: k.mailSeq, From: from, To: to, Kind: kind, Body: strings.TrimSpace(body)}
	k.mail = append(k.mail, m)
	if len(k.mail) > maxMail {
		k.mail = append([]Mail(nil), k.mail[len(k.mail)-maxMail:]...)
	}
	k.seq++
	k.events = append(k.events, Event{Seq: k.seq, At: k.now(), Source: from, Kind: "mail", Message: "to " + to + " · " + kind})
	return m, nil
}

// Inbox returns mail for one agent, oldest first.
func (k *Kernel) Inbox(to string) []Mail {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Mail, 0)
	for _, m := range k.mail {
		if to == "" || m.To == to {
			out = append(out, m)
		}
	}
	return out
}

// WriteNote stores a sourced finding.
func (k *Kernel) WriteNote(n Note) Note {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.noteSeq++
	n.ID = k.noteSeq
	n.Claim = strings.TrimSpace(n.Claim)
	n.Quote = clip(n.Quote, 500)
	k.notes = append(k.notes, n)
	if len(k.notes) > maxNotes {
		k.notes = append([]Note(nil), k.notes[len(k.notes)-maxNotes:]...)
	}
	k.seq++
	k.events = append(k.events, Event{Seq: k.seq, At: k.now(), Source: n.Agent, Kind: "note", Message: n.Claim, Data: map[string]any{"source": n.Source, "url": n.URL}})
	return n
}

// Notes returns research records, newest first.
func (k *Kernel) Notes() []Note {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Note, len(k.notes))
	copy(out, k.notes)
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out
}

// Remember writes an episodic fact.
func (k *Kernel) Remember(topic, text string) Fact {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.factSeq++
	f := Fact{ID: k.factSeq, Topic: strings.TrimSpace(topic), Text: strings.TrimSpace(text)}
	k.facts = append(k.facts, f)
	if len(k.facts) > maxFacts {
		k.facts = append([]Fact(nil), k.facts[len(k.facts)-maxFacts:]...)
	}
	return f
}

// Recall returns facts whose topic or text contains q.
func (k *Kernel) Recall(q string) []Fact {
	k.mu.RLock()
	defer k.mu.RUnlock()
	q = strings.ToLower(strings.TrimSpace(q))
	out := make([]Fact, 0)
	for _, f := range k.facts {
		if q == "" || strings.Contains(strings.ToLower(f.Topic), q) || strings.Contains(strings.ToLower(f.Text), q) {
			out = append(out, f)
		}
	}
	return out
}

// RequestConfirm parks outbound work until a human allows it.
func (k *Kernel) RequestConfirm(agent, cap, body string) Confirm {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.confSeq++
	c := Confirm{ID: k.confSeq, Agent: agent, Cap: cap, Body: strings.TrimSpace(body), Status: "pending"}
	k.confirms = append(k.confirms, c)
	k.seq++
	k.events = append(k.events, Event{Seq: k.seq, At: k.now(), Source: agent, Kind: "confirm.pending", Message: body})
	return c
}

// DecideConfirm allows or denies a pending confirm.
func (k *Kernel) DecideConfirm(id int, allow bool) (Confirm, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	for i := range k.confirms {
		if k.confirms[i].ID != id {
			continue
		}
		if k.confirms[i].Status != "pending" {
			return k.confirms[i], nil
		}
		if allow {
			k.confirms[i].Status = "allowed"
			k.confirms[i].Outcome = "human allowed · no outbound bind yet"
		} else {
			k.confirms[i].Status = "denied"
			k.confirms[i].Outcome = "human denied"
		}
		k.seq++
		k.events = append(k.events, Event{Seq: k.seq, At: k.now(), Source: "comms", Kind: "confirm." + k.confirms[i].Status, Message: k.confirms[i].Body})
		return k.confirms[i], nil
	}
	return Confirm{}, fmt.Errorf("confirm %d not found", id)
}

// Confirms returns the human gate queue.
func (k *Kernel) Confirms() []Confirm {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Confirm, len(k.confirms))
	copy(out, k.confirms)
	return out
}

// TeamCard returns the Proximity Microsoft Teams conversation dock card.
func (k *Kernel) TeamCard() TeamCard {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return TeamCard{
		Channel: TeamChannelMicrosoft,
		Kind:    TeamKindConversation,
		Tenant:  TeamTenantProximity,
		Dock:    "bottom",
		Hint:    "Proximity Microsoft Teams conversation AI · bottom dock · no Slack, mail, meetings, or files",
		Threads: len(k.threads),
	}
}

// OpenTeamThread starts a Proximity Microsoft Teams conversation.
func (k *Kernel) OpenTeamThread(in TeamOpen) (TeamThread, error) {
	tenant, err := NormalizeTeamTenant(in.Tenant)
	if err != nil {
		return TeamThread{}, err
	}
	channel, err := NormalizeTeamChannel(in.Channel)
	if err != nil {
		return TeamThread{}, err
	}
	kind, err := NormalizeTeamKind(in.Kind)
	if err != nil {
		return TeamThread{}, err
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		title = "Proximity desk"
	}

	k.mu.Lock()
	defer k.mu.Unlock()
	if k.threads == nil {
		k.threads = make(map[string]*TeamThread)
	}
	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = slugThread(title)
	}
	if existing, ok := k.threads[id]; ok {
		return cloneTeamThread(existing), nil
	}
	if len(k.threads) >= maxTeamThreads {
		return TeamThread{}, fmt.Errorf("too many threads")
	}
	now := k.now()
	t := &TeamThread{
		ID:        id,
		Title:     title,
		Tenant:    tenant,
		Channel:   channel,
		Kind:      kind,
		Dock:      "bottom",
		Messages:  nil,
		CreatedAt: now,
		UpdatedAt: now,
	}
	k.threads[id] = t
	k.seq++
	k.events = append(k.events, Event{Seq: k.seq, At: now, Source: "teams", Kind: "thread.open", Message: title, Data: map[string]any{"id": id, "tenant": tenant}})
	return cloneTeamThread(t), nil
}

// ListTeamThreads returns Proximity Teams conversations, newest first.
func (k *Kernel) ListTeamThreads() []TeamThread {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]TeamThread, 0, len(k.threads))
	for _, t := range k.threads {
		out = append(out, cloneTeamThread(t))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.After(out[j].UpdatedAt) })
	return out
}

// GetTeamThread returns one conversation by id.
func (k *Kernel) GetTeamThread(id string) (TeamThread, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	t, ok := k.threads[id]
	if !ok {
		return TeamThread{}, fmt.Errorf("%w: %s", ErrTeamThread, id)
	}
	return cloneTeamThread(t), nil
}

// AppendTeamMessage adds one turn to a conversation.
func (k *Kernel) AppendTeamMessage(id, role, body string) (TeamThread, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return TeamThread{}, ErrTeamEmpty
	}
	role = strings.TrimSpace(role)
	if role == "" {
		role = "user"
	}

	k.mu.Lock()
	defer k.mu.Unlock()
	t, ok := k.threads[id]
	if !ok {
		return TeamThread{}, fmt.Errorf("%w: %s", ErrTeamThread, id)
	}
	k.msgSeq++
	now := k.now()
	t.Messages = append(t.Messages, TeamMessage{ID: k.msgSeq, Role: role, Body: clip(body, 4000), At: now})
	if len(t.Messages) > maxTeamMsgs {
		t.Messages = append([]TeamMessage(nil), t.Messages[len(t.Messages)-maxTeamMsgs:]...)
	}
	t.UpdatedAt = now
	k.seq++
	k.events = append(k.events, Event{Seq: k.seq, At: now, Source: "teams", Kind: "thread." + role, Message: clip(body, 120), Data: map[string]any{"id": id}})
	return cloneTeamThread(t), nil
}

// NormalizeTeamChannel accepts Microsoft Teams aliases and rejects every other channel.
func NormalizeTeamChannel(s string) (string, error) {
	n := foldKey(s)
	switch n {
	case "", "microsoft-teams", "microsoftteams", "microsoft", "teams", "ms-teams", "msteams",
		"teams-microsoft", "teams-microsoft-com", "teams-microsfot":
		return TeamChannelMicrosoft, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrTeamChannel, strings.TrimSpace(s))
	}
}

// NormalizeTeamKind accepts conversation aliases and rejects meetings, files, and calls.
func NormalizeTeamKind(s string) (string, error) {
	n := foldKey(s)
	switch n {
	case "", "conversation", "chat", "convo", "thread", "dm", "1-1", "11":
		return TeamKindConversation, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrTeamKind, strings.TrimSpace(s))
	}
}

// NormalizeTeamTenant accepts Proximity aliases, including the common lroximity typo.
func NormalizeTeamTenant(s string) (string, error) {
	n := foldKey(s)
	switch n {
	case "", "proximity", "proximity-agency", "proximityagency", "lroximity", "proximityagency-ca":
		return TeamTenantProximity, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrTeamTenant, strings.TrimSpace(s))
	}
}

func foldKey(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.' || r == '-' || r == '_' || r == '/':
			b.WriteByte('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func slugThread(name string) string {
	n := foldKey(name)
	if n == "" {
		return "proximity-desk"
	}
	return n
}

func cloneTeamThread(t *TeamThread) TeamThread {
	out := *t
	if t.Messages != nil {
		out.Messages = append([]TeamMessage(nil), t.Messages...)
	}
	return out
}

func clip(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n])
}
