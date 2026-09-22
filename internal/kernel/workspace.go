package kernel

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

const (
	maxMail  = 64
	maxNotes = 200
	maxFacts = 200
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

func clip(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n])
}
