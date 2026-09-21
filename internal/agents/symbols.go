package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/symbols"
)

type symbolsAgent struct {
	bus *symbols.Bus
	k   *kernel.Kernel
}

func (a *symbolsAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "symbols", Name: "Symbols", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "zap", Summary: "Internal Zapier. Trigger → kernel verbs. Zero subscription.",
		Capabilities: []string{"symbols.list", "symbols.create", "symbols.fire", "symbols.runs", "symbols.toggle", "symbols.connectors"},
		Autostart:    true,
	}
}

func (a *symbolsAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	if a.bus == nil {
		a.bus = symbols.New()
	}
	a.k = k
	k.AttachSymbols(a.bus)
	k.SetSink(func(ev kernel.Event) {
		go a.onEvent(ev)
	})
	k.Publish("symbols", "board", fmt.Sprintf("internal zapier online · %d symbols · $0", len(a.bus.List())), map[string]any{
		"symbols": len(a.bus.List()),
	})
	return nil
}

func (a *symbolsAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "symbols.list":
		return kernel.Result{OK: true, Message: "symbols " + fmt.Sprint(len(a.bus.List())), Data: a.bus.List()}, nil
	case "symbols.runs":
		return kernel.Result{OK: true, Message: "runs", Data: a.bus.Runs()}, nil
	case "symbols.connectors":
		return kernel.Result{OK: true, Message: "connectors", Data: a.k.Capabilities()}, nil
	case "symbols.toggle":
		id := payloadQuery(call, "id")
		s, err := a.bus.Toggle(id)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: toggleMsg(s), Data: s}, nil
	case "symbols.create":
		var in symbols.Create
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &in); err != nil {
				return kernel.Result{}, err
			}
		}
		s, err := a.bus.Create(in)
		if err != nil {
			return kernel.Result{}, err
		}
		a.k.Publish("symbols", "create", "composed "+s.Name, map[string]any{"id": s.ID})
		return kernel.Result{OK: true, Message: "composed " + s.Name, Data: s}, nil
	case "symbols.fire":
		id := payloadQuery(call, "id")
		input := symbols.FlattenJSON(call.Payload)
		delete(input, "id")
		run, err := a.fire(ctx, id, input)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: run.OK, Message: run.Message, Data: run}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func (a *symbolsAgent) onEvent(ev kernel.Event) {
	if ev.Source == "symbols" {
		return
	}
	for _, s := range a.bus.List() {
		if !s.Enabled || !symbols.MatchEvent(s.On, ev.Source, ev.Kind) {
			continue
		}
		input := map[string]string{
			"source":  ev.Source,
			"kind":    ev.Kind,
			"message": ev.Message,
			"body":    ev.Message,
			"claim":   ev.Message,
		}
		_, _ = a.fire(context.Background(), s.ID, input)
	}
}

func (a *symbolsAgent) fire(ctx context.Context, id string, input map[string]string) (symbols.Run, error) {
	if err := a.bus.Enter(id); err != nil {
		return symbols.Run{}, err
	}
	defer a.bus.Leave(id)

	sym, err := a.bus.Get(id)
	if err != nil {
		return symbols.Run{}, err
	}

	logs := make([]symbols.StepLog, 0, len(sym.Steps))
	ok := true
	msg := "fired " + sym.Name
	for _, step := range sym.Steps {
		fields := symbols.Render(step.Fields, input)
		raw, _ := json.Marshal(fields)
		res, err := a.k.InvokeCap(ctx, step.Cap, kernel.Call{Payload: raw})
		sl := symbols.StepLog{Cap: step.Cap}
		if err != nil {
			sl.OK = false
			sl.Message = err.Error()
			logs = append(logs, sl)
			ok = false
			msg = err.Error()
			break
		}
		sl.OK = res.OK
		sl.Message = res.Message
		logs = append(logs, sl)
		if !res.OK {
			ok = false
			msg = res.Message
			break
		}
	}
	run, err := a.bus.Finish(sym.ID, ok, msg, input, logs)
	if err != nil {
		return symbols.Run{}, err
	}
	kind := "fire"
	if !ok {
		kind = "fire.error"
	}
	a.k.Publish("symbols", kind, run.Message, map[string]any{"id": sym.ID, "ok": ok})
	return run, nil
}

func toggleMsg(s symbols.Symbol) string {
	if s.Enabled {
		return s.Name + " on"
	}
	return s.Name + " off"
}
