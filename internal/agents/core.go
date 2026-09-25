package agents

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// CoreInvoke creates the next kernel on this core.
// Child agents stay on the child. The parent process table does not grow.
func CoreInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	var body struct {
		ID             string        `json:"id"`
		Infrastructure string        `json:"infrastructure"`
		Agents         []kernel.Spec `json:"agents"`
	}
	if len(call.Payload) > 0 {
		if err := json.Unmarshal(call.Payload, &body); err != nil {
			return kernel.Result{OK: false, Message: "payload illisible"}, nil
		}
	}
	if strings.TrimSpace(body.Infrastructure) == "" {
		body.Infrastructure = trinityBuild()
	}
	link, err := k.Core().SpawnKernel(context.Background(), body.ID, body.Infrastructure, body.Agents)
	if err != nil {
		return kernel.Result{OK: false, Message: err.Error()}, nil
	}
	chain := k.Core().Chain()
	k.Publish("init", "core", "kernel chained", map[string]any{
		"id": link.ID, "index": link.Index, "agents": len(link.Agents),
	})
	return kernel.Result{OK: true, Message: "kernel chained", Data: map[string]any{
		"link":         link,
		"chain":        chain,
		"parentAgents": k.About().Agents,
	}}, nil
}
