package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/model"
)

const (
	proximityDeskID    = "proximity-desk"
	teamsSystemPrompt  = "You are Cashtro OS conversation AI for Proximity Agency. Channel is Microsoft Teams only. Conversation only — refuse Slack, email, meetings, files, and calls. Keep replies short. Outbound to real Teams still needs a human confirm."
	teamsLocalFallback = "Proximity desk · Microsoft Teams conversation only. I stay in this bottom dock. I will not open Slack, mail, meetings, or files. Outbound to real Teams still waits on a human confirm."
)

type teamsAgent struct {
	k *kernel.Kernel
}

func (a *teamsAgent) Spec() kernel.Spec {
	return kernel.Spec{
		ID: "teams", Name: "Teams", Kind: kernel.KindUser, Mode: kernel.ModeLive,
		Role: "conversation", Summary: "Proximity Microsoft Teams conversation AI on the bottom dock. No other channel.",
		Capabilities: []string{"teams.status", "teams.list", "teams.thread", "teams.say", "teams.send"},
		Autostart:    true,
	}
}

func (a *teamsAgent) Boot(ctx context.Context, k *kernel.Kernel) error {
	a.k = k
	th, err := k.OpenTeamThread(kernel.TeamOpen{
		ID:      proximityDeskID,
		Title:   "Proximity desk",
		Tenant:  kernel.TeamTenantProximity,
		Channel: kernel.TeamChannelMicrosoft,
		Kind:    kernel.TeamKindConversation,
	})
	if err != nil {
		return err
	}
	if len(th.Messages) == 0 {
		if _, err := k.AppendTeamMessage(th.ID, "assistant", teamsLocalFallback); err != nil {
			return err
		}
	}
	k.Publish("teams", "dock", "Proximity Microsoft Teams conversation · bottom", map[string]any{
		"tenant": kernel.TeamTenantProximity, "channel": kernel.TeamChannelMicrosoft, "kind": kernel.TeamKindConversation,
	})
	return nil
}

func (a *teamsAgent) Invoke(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	switch call.Capability {
	case "teams.status":
		return kernel.Result{OK: true, Message: "teams card", Data: a.k.TeamCard()}, nil
	case "teams.list":
		threads := a.k.ListTeamThreads()
		return kernel.Result{OK: true, Message: fmt.Sprintf("threads %d", len(threads)), Data: threads}, nil
	case "teams.thread":
		var in kernel.TeamOpen
		if err := decodeTeams(call, &in); err != nil {
			return kernel.Result{}, err
		}
		if in.ID != "" && in.Title == "" {
			th, err := a.k.GetTeamThread(in.ID)
			if err == nil {
				return kernel.Result{OK: true, Message: "thread " + th.ID, Data: th}, nil
			}
		}
		th, err := a.k.OpenTeamThread(in)
		if err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
		return kernel.Result{OK: true, Message: "thread " + th.ID, Data: th}, nil
	case "teams.say":
		return a.say(ctx, call)
	case "teams.send":
		body := payloadQuery(call, "body")
		if body == "" {
			body = payloadQuery(call, "prompt")
		}
		if body == "" {
			body = payloadQuery(call, "text")
		}
		if body == "" {
			return kernel.Result{OK: false, Message: kernel.ErrTeamEmpty.Error()}, nil
		}
		if _, err := kernel.NormalizeTeamChannel(payloadQuery(call, "channel")); err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
		if _, err := kernel.NormalizeTeamTenant(payloadQuery(call, "tenant")); err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
		c := a.k.RequestConfirm("teams", "teams.send", "Proximity Microsoft Teams · "+body)
		return kernel.Result{OK: true, Message: "parked confirm #" + itoa(c.ID), Data: c}, nil
	default:
		return kernel.Result{}, kernel.ErrUnknownCapability
	}
}

func (a *teamsAgent) say(ctx context.Context, call kernel.Call) (kernel.Result, error) {
	text := firstNonEmpty(payloadQuery(call, "text"), payloadQuery(call, "prompt"), payloadQuery(call, "body"))
	if text == "" {
		return kernel.Result{OK: false, Message: kernel.ErrTeamEmpty.Error()}, nil
	}
	if ch := payloadQuery(call, "channel"); ch != "" {
		if _, err := kernel.NormalizeTeamChannel(ch); err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
	}
	if tenant := payloadQuery(call, "tenant"); tenant != "" {
		if _, err := kernel.NormalizeTeamTenant(tenant); err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
	}
	if kind := payloadQuery(call, "kind"); kind != "" {
		if _, err := kernel.NormalizeTeamKind(kind); err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
	}

	id := payloadQuery(call, "thread")
	if id == "" {
		id = payloadQuery(call, "id")
	}
	if id == "" {
		id = proximityDeskID
	}
	if _, err := a.k.GetTeamThread(id); err != nil {
		if _, err := a.k.OpenTeamThread(kernel.TeamOpen{ID: id, Title: id}); err != nil {
			return kernel.Result{OK: false, Message: err.Error()}, nil
		}
	}

	if _, err := a.k.AppendTeamMessage(id, "user", text); err != nil {
		return kernel.Result{}, err
	}
	reply := a.reply(ctx, id, text)
	th, err := a.k.AppendTeamMessage(id, "assistant", reply)
	if err != nil {
		return kernel.Result{}, err
	}
	return kernel.Result{OK: true, Message: "said in " + th.ID, Data: th}, nil
}

func (a *teamsAgent) reply(ctx context.Context, threadID, text string) string {
	th, err := a.k.GetTeamThread(threadID)
	if err != nil {
		return teamsLocalFallback
	}
	msgs := []model.Message{{Role: "system", Content: teamsSystemPrompt}}
	for _, m := range th.Messages {
		role := m.Role
		if role != "user" && role != "assistant" && role != "system" {
			role = "user"
		}
		msgs = append(msgs, model.Message{Role: role, Content: m.Body})
	}
	raw, _ := json.Marshal(map[string]any{"messages": msgs, "prompt": text})
	res, err := a.k.Invoke(ctx, "router", kernel.Call{Capability: "model.chat", Payload: raw})
	if err != nil || !res.OK {
		return localTeamsReply(text)
	}
	if content := chatContent(res.Data); content != "" {
		return content
	}
	return localTeamsReply(text)
}

func localTeamsReply(text string) string {
	low := strings.ToLower(text)
	switch {
	case strings.Contains(low, "slack"), strings.Contains(low, "gmail"), strings.Contains(low, "mail"),
		strings.Contains(low, "discord"):
		return "This dock is Microsoft Teams only, Proximity tenant. I will not open Slack, mail, or another channel."
	case strings.Contains(low, "meeting"), strings.Contains(low, "calendar"), strings.Contains(low, "call"),
		strings.Contains(low, "file"):
		return "Conversation only. I will not schedule meetings, join calls, or open files. Stay in this Proximity Teams thread."
	case strings.Contains(low, "send"), strings.Contains(low, "post"), strings.Contains(low, "outbound"):
		return "I can draft here. teams.send parks a human confirm before anything leaves for real Microsoft Teams."
	default:
		return teamsLocalFallback
	}
}

func chatContent(data any) string {
	switch d := data.(type) {
	case model.ChatResponse:
		return strings.TrimSpace(d.Content)
	case *model.ChatResponse:
		if d != nil {
			return strings.TrimSpace(d.Content)
		}
	case map[string]any:
		if s, ok := d["content"].(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func decodeTeams(call kernel.Call, dst *kernel.TeamOpen) error {
	if len(call.Payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(call.Payload, dst); err != nil {
		dst.Title = strings.TrimSpace(string(call.Payload))
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
