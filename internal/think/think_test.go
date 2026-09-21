package think

import (
	"strings"
	"testing"
)

func TestLocalOwnsTheLoop(t *testing.T) {
	got := Local("Advance ScanApp toward production. Keep confirms.")
	if got.Owner != Owner || got.Bound {
		t.Fatalf("dual = %+v", got)
	}
	if got.Propose.Seat != SeatPropose || got.Critique.Seat != SeatCritique {
		t.Fatalf("seats = %+v / %+v", got.Propose, got.Critique)
	}
	if !strings.Contains(got.Propose.Content, "human gate") {
		t.Fatalf("propose = %q", got.Propose.Content)
	}
	if got.Score < 70 {
		t.Fatalf("score = %d, want a gated proposal to score well", got.Score)
	}
	if !strings.Contains(got.Critique.Content, "SCORE:") {
		t.Fatalf("critique = %q", got.Critique.Content)
	}
	if strings.Contains(strings.ToLower(got.Propose.Content+got.Critique.Content), "moonshot") {
		t.Fatal("owned loop must not mention Moonshot")
	}
}

func TestLocalEmptyPromptStillRuns(t *testing.T) {
	got := Local("   ")
	if got.Prompt == "" || got.Score < 1 {
		t.Fatalf("empty prompt dual = %+v", got)
	}
}

func TestVendorLeakDropsScore(t *testing.T) {
	clean := Local("Park a retainer. Keep the human gate. Do not send outbound mail.")
	leaky := Local("Call moonshot and n8n.io then send outbound mail.")
	if leaky.Score >= clean.Score {
		t.Fatalf("leak score %d >= clean %d", leaky.Score, clean.Score)
	}
}
