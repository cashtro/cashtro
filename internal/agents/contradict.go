package agents

import (
	"strconv"
	"strings"
)

// Option is one way to do the work. Lower cost, risk, and steps wins.
type Option struct {
	Name  string `json:"name"`
	Cost  int    `json:"cost"`
	Risk  int    `json:"risk"`
	Steps int    `json:"steps"`
}

// Proposal is what a department wants to do.
type Proposal struct {
	Subject string   `json:"subject"`
	Options []Option `json:"options"`
}

// Verdict is the contradiction. Ready is true only when one option is strictly leaner.
type Verdict struct {
	Attack string `json:"attack"`
	Best   string `json:"best,omitempty"`
	Ready  bool   `json:"ready"`
	Score  int    `json:"score,omitempty"`
}

func (o Option) score() int { return o.Cost + o.Risk + o.Steps }

// Contradict attacks the proposal and keeps only the most optimized option.
func Contradict(p Proposal) Verdict {
	if len(p.Options) == 0 {
		return Verdict{Attack: "rien à contredire. Pose au moins deux options.", Ready: false}
	}
	if len(p.Options) == 1 {
		return Verdict{
			Attack: "une seule option n'est pas un choix. Il en faut une autre, plus courte.",
			Ready:  false,
		}
	}
	best := 0
	tie := false
	for i := 1; i < len(p.Options); i++ {
		if p.Options[i].score() < p.Options[best].score() {
			best = i
			tie = false
		} else if p.Options[i].score() == p.Options[best].score() {
			tie = true
		}
	}
	if tie {
		return Verdict{
			Attack: "égalité. Aucune option n'est plus optimisée. Couper des étapes ou du coût.",
			Ready:  false,
		}
	}
	winner := p.Options[best]
	var rejected []string
	for i, opt := range p.Options {
		if i == best {
			continue
		}
		rejected = append(rejected, opt.Name+" ("+strconv.Itoa(opt.score())+")")
	}
	attack := "on rejette " + strings.Join(rejected, ", ") + ". " + winner.Name + " est plus optimisée (" + strconv.Itoa(winner.score()) + ")."
	return Verdict{Attack: attack, Best: winner.Name, Ready: true, Score: winner.score()}
}
