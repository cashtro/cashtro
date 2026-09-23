package agents

import (
	"strconv"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// Tax is Quebec sales tax on one amount. TPS is 5%. TVQ is 9.975% of the
// same amount, not of the amount plus TPS. This is a calculator, not a filing.
type Tax struct {
	Cents int `json:"cents"`
	TPS   int `json:"tps"`
	TVQ   int `json:"tvq"`
	Total int `json:"total"`
}

// TaxQuebec rounds each tax to the nearest cent.
func TaxQuebec(cents int) Tax {
	if cents < 0 {
		cents = 0
	}
	tps := (cents*5 + 50) / 100
	tvq := (cents*9975 + 50000) / 100000
	return Tax{Cents: cents, TPS: tps, TVQ: tvq, Total: cents + tps + tvq}
}

func TaxInvoke(call kernel.Call) (kernel.Result, error) {
	cents := payloadInt(call, "cents")
	if cents == 0 {
		if raw := payloadQuery(call, "amount"); raw != "" {
			cents = parseCents(raw)
		}
	}
	if cents <= 0 {
		return kernel.Result{OK: false, Message: "montant requis"}, nil
	}
	tax := TaxQuebec(cents)
	return kernel.Result{OK: true, Message: "TPS et TVQ calculées", Data: tax}, nil
}

func parseCents(raw string) int {
	raw = strings.TrimSpace(strings.TrimPrefix(raw, "$"))
	raw = strings.ReplaceAll(raw, ",", ".")
	if raw == "" {
		return 0
	}
	if !strings.Contains(raw, ".") {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			return 0
		}
		return n * 100
	}
	parts := strings.SplitN(raw, ".", 2)
	dollars, err := strconv.Atoi(parts[0])
	if err != nil || dollars < 0 {
		return 0
	}
	frac := parts[1]
	if len(frac) > 2 {
		frac = frac[:2]
	}
	for len(frac) < 2 {
		frac += "0"
	}
	cents, err := strconv.Atoi(frac)
	if err != nil {
		return 0
	}
	return dollars*100 + cents
}

// Check is an internal draft. It is not a bank instrument.
type Check struct {
	Payee   string `json:"payee"`
	Cents   int    `json:"cents"`
	Memo    string `json:"memo"`
	Status  string `json:"status"`
	Confirm int    `json:"confirm"`
}

func CheckInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	payee := payloadQuery(call, "payee")
	memo := payloadQuery(call, "memo")
	cents := payloadInt(call, "cents")
	if cents == 0 {
		cents = parseCents(payloadQuery(call, "amount"))
	}
	if payee == "" || cents <= 0 {
		return kernel.Result{OK: false, Message: "bénéficiaire et montant requis"}, nil
	}
	if names := secretNames(payee + " " + memo); len(names) > 0 {
		return kernel.Result{OK: false, Message: "interdit: " + strings.Join(names, ", ")}, nil
	}
	body := payee + " · " + strconv.Itoa(cents) + " cents · " + memo
	confirm := k.RequestConfirm("comms", "comms.allow", body)
	check := Check{Payee: payee, Cents: cents, Memo: memo, Status: "brouillon", Confirm: confirm.ID}
	k.Remember("cheque", body)
	return kernel.Result{OK: true, Message: "chèque en brouillon", Data: check}, nil
}

func FundsInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	envelope := payloadQuery(call, "envelope")
	if envelope == "" {
		envelope = "caisse"
	}
	cents := payloadInt(call, "cents")
	if cents != 0 {
		k.Remember("fonds:"+envelope, strconv.Itoa(cents))
	}
	return kernel.Result{OK: true, Message: "position du grand livre", Data: FundsPosition(k)}, nil
}

// FundsPosition reads the envelopes held in memory. It does not move money.
func FundsPosition(k *kernel.Kernel) map[string]int {
	out := map[string]int{}
	for _, fact := range k.Recall("fonds:") {
		name := strings.TrimPrefix(fact.Topic, "fonds:")
		if name == "" || name == fact.Topic {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(fact.Text))
		if err != nil {
			continue
		}
		out[name] = n
	}
	if len(out) == 0 {
		out["caisse"] = 0
	}
	return out
}
