package kernel

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cashtro/cashtro/internal/catalog"
)

const maxChoices = 200

// Choice kinds Castro picks from on the desk.
const (
	KindMail = "mail"
	KindVerb = "verb"
	KindScan = "scan"
)

// Choice statuses.
const (
	StatusPending = "pending"
	StatusTaken   = "taken"
	StatusSkipped = "skipped"
)

// Choice is one job on the human chooser desk.
type Choice struct {
	ID      int       `json:"id"`
	Key     string    `json:"key"`
	Kind    string    `json:"kind"`
	Status  string    `json:"status"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Client  string    `json:"client,omitempty"`
	Sector  string    `json:"sector,omitempty"`
	Source  string    `json:"source,omitempty"`
	Urgency string    `json:"urgency,omitempty"`
	Stack   []string  `json:"stack,omitempty"`
	ShipID  string    `json:"shipId,omitempty"`
	Created time.Time `json:"createdAt"`
	Decided time.Time `json:"decidedAt,omitempty"`
}

// Scan is the last human-inbox pass. Device mail apps are out of reach;
// this is the connected Gmail mailbox the Samsung Gmail app would show.
type Scan struct {
	At      time.Time `json:"at"`
	Account string    `json:"account"`
	Device  string    `json:"device"`
	Inbox   int       `json:"inbox"`
	Unread  int       `json:"unread"`
	Starred int       `json:"starred"`
	Note    string    `json:"note"`
}

// Desk is the chooser card: scan stats plus jobs Castro can take or skip.
type Desk struct {
	Scan     Scan     `json:"scan"`
	Pending  []Choice `json:"pending"`
	Taken    []Choice `json:"taken"`
	Skipped  []Choice `json:"skipped"`
	Mail     []Choice `json:"mail"`
	Verbs    []Choice `json:"verbs"`
	PendingN int      `json:"pendingCount"`
}

// RecordScan stores the last inbox pass.
func (k *Kernel) RecordScan(s Scan) Scan {
	k.mu.Lock()
	if s.At.IsZero() {
		s.At = k.now()
	}
	s.Account = strings.TrimSpace(s.Account)
	s.Device = strings.TrimSpace(s.Device)
	s.Note = strings.TrimSpace(s.Note)
	k.scan = s
	k.seq++
	k.events = append(k.events, Event{
		Seq:     k.seq,
		At:      k.now(),
		Source:  "chooser",
		Kind:    "scan",
		Message: fmt.Sprintf("scan · %d inbox · %d unread", s.Inbox, s.Unread),
	})
	out := s
	k.mu.Unlock()
	k.persist()
	return out
}

// LastScan returns the last inbox pass.
func (k *Kernel) LastScan() Scan {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.scan
}

// Offer parks a job on the chooser if the key is new.
func (k *Kernel) Offer(c Choice) Choice {
	k.mu.Lock()
	c.Key = strings.TrimSpace(c.Key)
	c.Title = strings.TrimSpace(c.Title)
	c.Summary = strings.TrimSpace(c.Summary)
	if c.Key == "" {
		c.Key = slug(c.Kind + "-" + c.Title)
	}
	for _, existing := range k.choices {
		if existing.Key == c.Key {
			k.mu.Unlock()
			return existing
		}
	}
	k.choiceSeq++
	c.ID = k.choiceSeq
	if c.Status == "" {
		c.Status = StatusPending
	}
	if c.Kind == "" {
		c.Kind = KindMail
	}
	if c.Created.IsZero() {
		c.Created = k.now()
	}
	k.choices = append(k.choices, c)
	if len(k.choices) > maxChoices {
		k.choices = append([]Choice(nil), k.choices[len(k.choices)-maxChoices:]...)
	}
	k.seq++
	k.events = append(k.events, Event{
		Seq:     k.seq,
		At:      k.now(),
		Source:  "chooser",
		Kind:    "offer",
		Message: c.Title,
		Data:    map[string]any{"key": c.Key, "kind": c.Kind},
	})
	out := c
	k.mu.Unlock()
	k.persist()
	return out
}

// Choices returns desk jobs, newest first. Empty status returns all.
func (k *Kernel) Choices(status string) []Choice {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Choice, 0, len(k.choices))
	for _, c := range k.choices {
		if status == "" || c.Status == status {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out
}

// DeskCard is the chooser payload for the OS shell.
func (k *Kernel) DeskCard() Desk {
	all := k.Choices("")
	d := Desk{Scan: k.LastScan()}
	for _, c := range all {
		switch c.Status {
		case StatusTaken:
			d.Taken = append(d.Taken, c)
		case StatusSkipped:
			d.Skipped = append(d.Skipped, c)
		default:
			d.Pending = append(d.Pending, c)
		}
		switch c.Kind {
		case KindVerb:
			d.Verbs = append(d.Verbs, c)
		case KindMail, KindScan:
			d.Mail = append(d.Mail, c)
		}
	}
	d.PendingN = len(d.Pending)
	return d
}

// Take marks a pending job taken and, for mail, parks it on the delivery line.
func (k *Kernel) Take(id int) (Choice, error) {
	c, err := k.decide(id, StatusTaken)
	if err != nil {
		return Choice{}, err
	}
	if (c.Kind == KindMail || c.Kind == KindScan) && k.Catalog() != nil && c.Title != "" {
		ship, err := k.Catalog().Create(catalog.CreateShip{
			Name:   c.Title,
			Client: c.Client,
			Sector: firstNonEmpty(c.Sector, "ops"),
			Notes:  c.Summary,
			Stack:  c.Stack,
		})
		if err == nil {
			k.mu.Lock()
			for i := range k.choices {
				if k.choices[i].ID == c.ID {
					k.choices[i].ShipID = ship.ID
					c = k.choices[i]
					break
				}
			}
			k.mu.Unlock()
			k.persist()
		}
	}
	k.Remember("chooser", "took "+c.Title)
	return c, nil
}

// Skip marks a pending job skipped.
func (k *Kernel) Skip(id int) (Choice, error) {
	c, err := k.decide(id, StatusSkipped)
	if err != nil {
		return Choice{}, err
	}
	k.Remember("chooser", "skipped "+c.Title)
	return c, nil
}

func (k *Kernel) decide(id int, status string) (Choice, error) {
	k.mu.Lock()
	var out Choice
	found := false
	for i := range k.choices {
		if k.choices[i].ID != id {
			continue
		}
		found = true
		if k.choices[i].Status != StatusPending {
			k.mu.Unlock()
			return Choice{}, fmt.Errorf("%w: %s", ErrAlreadyDecided, k.choices[i].Key)
		}
		k.choices[i].Status = status
		k.choices[i].Decided = k.now()
		k.seq++
		k.events = append(k.events, Event{
			Seq:     k.seq,
			At:      k.now(),
			Source:  "chooser",
			Kind:    status,
			Message: k.choices[i].Title,
			Data:    map[string]any{"id": id, "key": k.choices[i].Key},
		})
		out = k.choices[i]
		break
	}
	k.mu.Unlock()
	if !found {
		return Choice{}, fmt.Errorf("%w: %d", ErrUnknownChoice, id)
	}
	k.persist()
	return out, nil
}

func firstNonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return strings.TrimSpace(v)
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
			continue
		}
		if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "job"
	}
	return out
}
