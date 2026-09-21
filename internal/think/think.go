// Package think is Cashtro's owned dual loop.
//
// Propose and critique run in this process. There is no Moonshot key,
// no GLM key, and no n8n webhook. OpenRouter stays the optional kernel
// model bus; this package never reads vendor seats.
package think

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	// Owner is who holds the loop. Not a vendor.
	Owner = "cashtro"
	// SeatPropose is the first-pass seat (the job Moonshot would have done).
	SeatPropose = "propose"
	// SeatCritique is the second-pass seat (the job GLM would have done).
	SeatCritique = "critique"
)

// Turn is one seat's output.
type Turn struct {
	Seat    string `json:"seat"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Score   int    `json:"score,omitempty"`
}

// Dual is propose then critique. Always local. Always ours.
type Dual struct {
	Owner    string `json:"owner"`
	Bound    bool   `json:"bound"`
	Prompt   string `json:"prompt"`
	Propose  Turn   `json:"propose"`
	Critique Turn   `json:"critique"`
	Score    int    `json:"score"`
	Message  string `json:"message"`
}

// Local runs the owned dual loop. Bound is always false — we do not
// phone Moonshot or Zhipu. The kernel can still think.
func Local(prompt string) Dual {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		prompt = "compound the live Evolu-Jeunes / Proximity line"
	}
	propose := propose(prompt)
	critique, score := critique(prompt, propose)
	return Dual{
		Owner:    Owner,
		Bound:    false,
		Prompt:   clip(prompt, 400),
		Propose:  Turn{Seat: SeatPropose, Role: "propose", Content: propose},
		Critique: Turn{Seat: SeatCritique, Role: "critique", Content: critique, Score: score},
		Score:    score,
		Message:  fmt.Sprintf("owned dual · score %d · no vendor keys", score),
	}
}

func propose(prompt string) string {
	focus := firstClause(prompt)
	return "Park one next move on " + focus +
		". Keep the human gate first-class. Do not send outbound mail. " +
		"Do not raise autonomy. Compound the existing Cashtro kernel — " +
		"do not fork a second product."
}

func critique(prompt, proposal string) (string, int) {
	score := 62
	notes := make([]string, 0, 6)
	blob := strings.ToLower(proposal + " " + prompt)

	if strings.Contains(blob, "human") || strings.Contains(blob, "confirm") || strings.Contains(blob, "gate") {
		score += 12
		notes = append(notes, "human gate present")
	} else {
		score -= 18
		notes = append(notes, "missing human gate")
	}
	if strings.Contains(blob, "autonomy") && (strings.Contains(blob, "never") || strings.Contains(blob, "do not") || strings.Contains(blob, "don't")) {
		score += 8
		notes = append(notes, "autonomy held")
	}
	if strings.Contains(blob, "outbound") || strings.Contains(blob, "send mail") {
		if strings.Contains(blob, "do not") || strings.Contains(blob, "never") || strings.Contains(blob, "don't") {
			score += 6
			notes = append(notes, "no outbound")
		} else {
			score -= 20
			notes = append(notes, "outbound risk")
		}
	}
	if strings.Contains(blob, "second product") || strings.Contains(blob, "kernel") {
		score += 6
		notes = append(notes, "stays on kernel")
	}
	if strings.Contains(blob, "moonshot") || strings.Contains(blob, "glm_api") || strings.Contains(blob, "n8n.io") {
		score -= 25
		notes = append(notes, "vendor leak")
	}
	if score < 1 {
		score = 1
	}
	if score > 100 {
		score = 100
	}
	if len(notes) == 0 {
		notes = append(notes, "neutral pass")
	}
	line := "Keep human confirms. Never raise own autonomy. Score " +
		fmt.Sprintf("%d", score) + ". " + strings.Join(notes, " · ") +
		fmt.Sprintf(". SCORE: %d", score)
	return line, score
}

func firstClause(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "the live board"
	}
	cut := strings.IndexAny(s, ".\n")
	if cut > 12 {
		s = s[:cut]
	}
	return clip(s, 80)
}

func clip(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	for i := n; i > 0; i-- {
		if unicode.IsSpace(r[i-1]) {
			return string(r[:i-1])
		}
	}
	return string(r[:n])
}
