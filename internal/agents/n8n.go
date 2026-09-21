package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cashtro/cashtro/internal/flow"
	"github.com/cashtro/cashtro/internal/kernel"
)

type n8nAgent struct {
	eng *flow.Engine
	k   *kernel.Kernel
}

func (a *n8nAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "n8n", Name: "n8n", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "workflow", Summary: "Internal n8n. Node graphs we own. Zero vendor webhook or model keys.",
		Capabilities: []string{"n8n.list", "n8n.get", "n8n.create", "n8n.workflow", "n8n.runs", "n8n.toggle", "n8n.nodes"},
		Autostart:    true,
	}
}

func (a *n8nAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	if a.eng == nil {
		a.eng = flow.New()
	}
	a.k = k
	k.Publish("n8n", "board", fmt.Sprintf("owned n8n online · %d workflows · $0 keys", len(a.eng.List())), map[string]any{
		"workflows": len(a.eng.List()),
		"owner":     "cashtro",
		"vendor":    false,
	})
	return nil
}

func (a *n8nAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "n8n.list":
		return kernel.Result{OK: true, Message: "workflows " + strconv.Itoa(len(a.eng.List())), Data: a.eng.List()}, nil
	case "n8n.runs":
		return kernel.Result{OK: true, Message: "runs", Data: a.eng.Runs()}, nil
	case "n8n.nodes":
		return kernel.Result{OK: true, Message: "owned node catalog", Data: flow.Catalog()}, nil
	case "n8n.get":
		id := payloadQuery(call, "id")
		wf, err := a.eng.Get(id)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: true, Message: wf.Name, Data: wf}, nil
	case "n8n.toggle":
		id := payloadQuery(call, "id")
		wf, err := a.eng.Toggle(id)
		if err != nil {
			return kernel.Result{}, err
		}
		msg := wf.Name + " off"
		if wf.Enabled {
			msg = wf.Name + " on"
		}
		return kernel.Result{OK: true, Message: msg, Data: wf}, nil
	case "n8n.create":
		var in flow.Create
		if len(call.Payload) > 0 {
			if err := json.Unmarshal(call.Payload, &in); err != nil {
				return kernel.Result{}, err
			}
		}
		wf, err := a.eng.Create(in)
		if err != nil {
			return kernel.Result{}, err
		}
		a.k.Publish("n8n", "create", "composed "+wf.Name, map[string]any{"id": wf.ID})
		return kernel.Result{OK: true, Message: "composed " + wf.Name, Data: wf}, nil
	case "n8n.workflow":
		run, err := a.fire(ctx, call)
		if err != nil {
			return kernel.Result{}, err
		}
		return kernel.Result{OK: run.OK, Message: run.Message, Data: run}, nil
	default:
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
}

func (a *n8nAgent) fire(ctx context.Context, call kernel.Call) (flow.Run, error) {
	input := flow.FlattenJSON(call.Payload)
	id := input["id"]
	if id == "" {
		id = "wealth-dual"
	}
	delete(input, "id")
	if err := a.eng.Enter(id); err != nil {
		return flow.Run{}, err
	}
	defer a.eng.Leave(id)

	run, err := a.eng.Execute(ctx, id, input, a.invoke)
	if err != nil {
		return flow.Run{}, err
	}
	kind := "workflow"
	if !run.OK {
		kind = "workflow.error"
	}
	a.k.Publish("n8n", kind, run.Message, map[string]any{
		"id":     run.Workflow,
		"ok":     run.OK,
		"score":  run.Score,
		"owner":  run.Owner,
		"bound":  run.Bound,
		"vendor": false,
	})
	return run, nil
}

func (a *n8nAgent) invoke(ctx context.Context, cap string, fields map[string]string) (bool, string, error) {
	raw, _ := json.Marshal(fields)
	res, err := a.k.InvokeCap(ctx, cap, kernel.Call{Payload: raw})
	if err != nil {
		return false, "", err
	}
	return res.OK, res.Message, nil
}
