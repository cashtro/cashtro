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
	maxMail     = 64
	maxNotes    = 200
	maxFacts    = 200
	maxRequests = 200

	StatusCaptured  = "captured"
	StatusClarified = "clarified"
	StatusDoing     = "doing"
	StatusDone      = "done"
	StatusBlocked   = "blocked"
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

// Request is one Castro ask on the desk. Raw is never overwritten.
type Request struct {
	ID        int       `json:"id"`
	Raw       string    `json:"raw"`
	Title     string    `json:"title"`
	Improved  string    `json:"improved"`
	Pass      int       `json:"pass"`
	Status    string    `json:"status"`
	Owner     string    `json:"owner,omitempty"`
	Verb      string    `json:"verb,omitempty"`
	Allowed   bool      `json:"allowed"`
	Reason    string    `json:"reason,omitempty"`
	Criteria  []string  `json:"criteria,omitempty"`
	Lessons   []string  `json:"lessons,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var (
	// ErrUnknownRequest is returned when a request id is not on the desk.
	ErrUnknownRequest = errors.New("request not found")
	// ErrInvalidRequest is returned when capture has no raw text.
	ErrInvalidRequest = errors.New("invalid request")
)

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
	if _, ok := k.procs[to]; !ok {
		k.mu.Unlock()
		return Mail{}, fmt.Errorf("%w: %s", ErrUnknownAgent, to)
	}
	k.mailSeq++
	m := Mail{ID: k.mailSeq, From: from, To: to, Kind: kind, Body: strings.TrimSpace(body)}
	k.mail = append(k.mail, m)
	if len(k.mail) > maxMail {
		k.mail = append([]Mail(nil), k.mail[len(k.mail)-maxMail:]...)
	}
	ev := k.appendEventLocked(from, "mail", "to "+to+" · "+kind, nil)
	k.mu.Unlock()
	k.fanout(ev)
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
	k.noteSeq++
	n.ID = k.noteSeq
	n.Claim = strings.TrimSpace(n.Claim)
	n.Quote = clip(n.Quote, 500)
	k.notes = append(k.notes, n)
	if len(k.notes) > maxNotes {
		k.notes = append([]Note(nil), k.notes[len(k.notes)-maxNotes:]...)
	}
	ev := k.appendEventLocked(n.Agent, "note", n.Claim, map[string]any{"source": n.Source, "url": n.URL})
	k.mu.Unlock()
	k.fanout(ev)
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
	k.confSeq++
	c := Confirm{ID: k.confSeq, Agent: agent, Cap: cap, Body: strings.TrimSpace(body), Status: "pending"}
	k.confirms = append(k.confirms, c)
	ev := k.appendEventLocked(agent, "confirm.pending", c.Body, nil)
	k.mu.Unlock()
	k.fanout(ev)
	return c
}

// DecideConfirm allows or denies a pending confirm.
func (k *Kernel) DecideConfirm(id int, allow bool) (Confirm, error) {
	k.mu.Lock()
	for i := range k.confirms {
		if k.confirms[i].ID != id {
			continue
		}
		if k.confirms[i].Status != "pending" {
			c := k.confirms[i]
			k.mu.Unlock()
			return c, nil
		}
		if allow {
			k.confirms[i].Status = "allowed"
			k.confirms[i].Outcome = "human allowed · no outbound bind yet"
		} else {
			k.confirms[i].Status = "denied"
			k.confirms[i].Outcome = "human denied"
		}
		ev := k.appendEventLocked("comms", "confirm."+k.confirms[i].Status, k.confirms[i].Body, nil)
		c := k.confirms[i]
		k.mu.Unlock()
		k.fanout(ev)
		return c, nil
	}
	k.mu.Unlock()
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

// SaveRequest creates or updates a desk request. Raw is set only on create.
func (k *Kernel) SaveRequest(r Request) (Request, error) {
	raw := strings.TrimSpace(r.Raw)
	if r.ID == 0 && raw == "" {
		return Request{}, fmt.Errorf("%w: raw is required", ErrInvalidRequest)
	}

	k.mu.Lock()
	defer k.mu.Unlock()
	now := k.now()

	if r.ID == 0 {
		k.reqSeq++
		r.ID = k.reqSeq
		r.Raw = raw
		r.CreatedAt = now
		if r.Status == "" {
			r.Status = StatusCaptured
		}
		r.UpdatedAt = now
		k.requests = append(k.requests, cloneRequest(r))
		if len(k.requests) > maxRequests {
			k.requests = append([]Request(nil), k.requests[len(k.requests)-maxRequests:]...)
		}
		k.seq++
		k.events = append(k.events, Event{Seq: k.seq, At: now, Source: "desk", Kind: "request.capture", Message: r.Title})
		return cloneRequest(r), nil
	}

	for i := range k.requests {
		if k.requests[i].ID != r.ID {
			continue
		}
		keep := k.requests[i].Raw
		created := k.requests[i].CreatedAt
		r.Raw = keep
		r.CreatedAt = created
		r.UpdatedAt = now
		k.requests[i] = cloneRequest(r)
		k.seq++
		k.events = append(k.events, Event{Seq: k.seq, At: now, Source: "desk", Kind: "request.better", Message: fmt.Sprintf("#%d pass %d", r.ID, r.Pass)})
		return cloneRequest(r), nil
	}
	return Request{}, fmt.Errorf("%w: %d", ErrUnknownRequest, r.ID)
}

// RequestByID returns one desk request.
func (k *Kernel) RequestByID(id int) (Request, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	for _, r := range k.requests {
		if r.ID == id {
			return cloneRequest(r), nil
		}
	}
	return Request{}, fmt.Errorf("%w: %d", ErrUnknownRequest, id)
}

// Requests returns desk asks, newest first.
func (k *Kernel) Requests() []Request {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Request, len(k.requests))
	for i, r := range k.requests {
		out[i] = cloneRequest(r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out
}

// AdvanceRequest sets a request status without touching the rewrite.
func (k *Kernel) AdvanceRequest(id int, status string) (Request, error) {
	status = strings.TrimSpace(status)
	switch status {
	case StatusCaptured, StatusClarified, StatusDoing, StatusDone, StatusBlocked:
	default:
		return Request{}, fmt.Errorf("%w: status %s", ErrInvalidRequest, status)
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	for i := range k.requests {
		if k.requests[i].ID != id {
			continue
		}
		k.requests[i].Status = status
		k.requests[i].UpdatedAt = k.now()
		k.seq++
		k.events = append(k.events, Event{Seq: k.seq, At: k.now(), Source: "desk", Kind: "request." + status, Message: k.requests[i].Title})
		return cloneRequest(k.requests[i]), nil
	}
	return Request{}, fmt.Errorf("%w: %d", ErrUnknownRequest, id)
}

func cloneRequest(r Request) Request {
	out := r
	if r.Criteria != nil {
		out.Criteria = append([]string(nil), r.Criteria...)
	}
	if r.Lessons != nil {
		out.Lessons = append([]string(nil), r.Lessons...)
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
