package agents

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// Move is one played action. The hash covers the previous move,
// so the progress chain is not rewritten. This is the board's log,
// not a coin and not a trade.
type Move struct {
	Index    int    `json:"index"`
	Prev     string `json:"prev"`
	Hash     string `json:"hash"`
	Action   string `json:"action"`
	Question string `json:"question"`
	Coup     string `json:"coup"`
}

func (m Move) seal() string {
	sum := sha256.Sum256([]byte(strconv.Itoa(m.Index) + "|" + m.Prev + "|" + m.Action + "|" + m.Question + "|" + m.Coup))
	return hex.EncodeToString(sum[:])
}

// AppendMove adds one move onto the chain.
func AppendMove(chain []Move, action, question, coup string) []Move {
	prev := ""
	idx := 0
	if n := len(chain); n > 0 {
		prev = chain[n-1].Hash
		idx = chain[n-1].Index + 1
	}
	m := Move{Index: idx, Prev: prev, Action: action, Question: question, Coup: coup}
	m.Hash = m.seal()
	return append(chain, m)
}

// ChainIntact reports whether every move still points at the one before it.
func ChainIntact(chain []Move) bool {
	if len(chain) == 0 {
		return false
	}
	prev := ""
	for i, m := range chain {
		if m.Index != i || m.Prev != prev || m.Hash != m.seal() {
			return false
		}
		prev = m.Hash
	}
	return true
}
