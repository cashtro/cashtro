package kernel

import (
	"sort"
	"strings"
)

const (
	maxOffers = 64
	maxGens   = 80
)

// Lane is how an offer compounds on the wealth board.
type Lane string

const (
	LaneCreate Lane = "create"
	LaneBuild  Lane = "build"
	LaneGrow   Lane = "grow"
	LaneEvolve Lane = "evolve"
)

// Offer is one cash-facing mandate the forge keeps alive.
type Offer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Lane       Lane   `json:"lane"`
	Thesis     string `json:"thesis"`
	SKU        string `json:"sku,omitempty"`
	Score      int    `json:"score"`
	Hits       int    `json:"hits"`
	Generation int    `json:"generation"`
}

// Generation is one forge tick. Score is monotonic.
type Generation struct {
	N       int    `json:"n"`
	Score   int    `json:"score"`
	Lane    Lane   `json:"lane"`
	Offer   string `json:"offer"`
	Kimi    string `json:"kimi,omitempty"`
	GLM     string `json:"glm,omitempty"`
	Bound   bool   `json:"bound"`
	Message string `json:"message"`
}

// Ledger is the public wealth card.
type Ledger struct {
	Generation int          `json:"generation"`
	Score      int          `json:"score"`
	Offers     []Offer      `json:"offers"`
	History    []Generation `json:"history"`
}

// UpsertOffer creates or refreshes a wealth offer by id.
func (k *Kernel) UpsertOffer(in Offer) Offer {
	k.mu.Lock()
	defer k.mu.Unlock()
	in.ID = strings.TrimSpace(in.ID)
	in.Name = strings.TrimSpace(in.Name)
	in.Thesis = strings.TrimSpace(in.Thesis)
	in.SKU = strings.TrimSpace(in.SKU)
	if in.Lane == "" {
		in.Lane = LaneCreate
	}
	for i := range k.offers {
		if k.offers[i].ID == in.ID {
			if in.Name != "" {
				k.offers[i].Name = in.Name
			}
			if in.Thesis != "" {
				k.offers[i].Thesis = in.Thesis
			}
			if in.SKU != "" {
				k.offers[i].SKU = in.SKU
			}
			if in.Lane != "" {
				k.offers[i].Lane = in.Lane
			}
			return k.offers[i]
		}
	}
	k.offers = append(k.offers, in)
	if len(k.offers) > maxOffers {
		k.offers = append([]Offer(nil), k.offers[len(k.offers)-maxOffers:]...)
	}
	return in
}

// TouchOffer increments hits and stamps the generation that evolved it.
func (k *Kernel) TouchOffer(id string, generation, add int) (Offer, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	for i := range k.offers {
		if k.offers[i].ID != id {
			continue
		}
		k.offers[i].Hits++
		k.offers[i].Generation = generation
		if add < 1 {
			add = 1
		}
		k.offers[i].Score += add
		return k.offers[i], true
	}
	return Offer{}, false
}

// Offers returns the wealth board, highest score first.
func (k *Kernel) Offers() []Offer {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Offer, len(k.offers))
	copy(out, k.offers)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// CommitGeneration records one forge tick. Score never falls.
func (k *Kernel) CommitGeneration(g Generation) Generation {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.gen++
	g.N = k.gen
	if g.Score <= k.score {
		g.Score = k.score + 1
	}
	k.score = g.Score
	g.Kimi = clip(g.Kimi, 280)
	g.GLM = clip(g.GLM, 280)
	k.gens = append(k.gens, g)
	if len(k.gens) > maxGens {
		k.gens = append([]Generation(nil), k.gens[len(k.gens)-maxGens:]...)
	}
	k.seq++
	k.events = append(k.events, Event{
		Seq:     k.seq,
		At:      k.now(),
		Source:  "forge",
		Kind:    "evolve",
		Message: g.Message,
		Data:    map[string]any{"n": g.N, "score": g.Score, "offer": g.Offer, "bound": g.Bound},
	})
	return g
}

// History returns forge ticks, newest first.
func (k *Kernel) History() []Generation {
	k.mu.RLock()
	defer k.mu.RUnlock()
	out := make([]Generation, len(k.gens))
	copy(out, k.gens)
	sort.Slice(out, func(i, j int) bool { return out[i].N > out[j].N })
	return out
}

// Ledger returns generation, score, offers, and recent history.
func (k *Kernel) Ledger() Ledger {
	return Ledger{
		Generation: k.Generation(),
		Score:      k.Score(),
		Offers:     k.Offers(),
		History:    k.History(),
	}
}

// Generation is the current forge tick count.
func (k *Kernel) Generation() int {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.gen
}

// Score is the monotonic wealth score.
func (k *Kernel) Score() int {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.score
}
