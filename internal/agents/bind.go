package agents

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cashtro/cashtro/internal/kernel"
)

var liveClientTokens = []string{
	"btk", "md clinic", "hypothe", "educonnexion", "proximity",
	"azure", "vercel", "production site",
}

func operatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	targets := []string{
		"http://127.0.0.1:8080/health",
		"http://127.0.0.1:8787/health",
	}
	asked := payloadQuery(call, "url")
	if asked != "" {
		if !allowedLocalURL(asked) {
			return kernel.Result{OK: false, Message: "operator refuses non-local URL", Data: map[string]any{"url": asked}}, nil
		}
		targets = []string{asked}
	}
	client := &http.Client{Timeout: 2 * time.Second}
	probes := make([]map[string]any, 0, len(targets))
	ok := true
	for _, u := range targets {
		row := map[string]any{"url": u}
		res, err := client.Get(u)
		if err != nil {
			ok = false
			row["error"] = err.Error()
			probes = append(probes, row)
			continue
		}
		_ = res.Body.Close()
		row["status"] = res.StatusCode
		if res.StatusCode >= 400 {
			ok = false
		}
		probes = append(probes, row)
	}
	k.Publish("operator", "browse", "probed local desk", map[string]any{"n": len(probes)})
	return kernel.Result{OK: ok, Message: "operator probed " + itoa(len(probes)), Data: probes}, nil
}

func reviewerInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	roots := []string{"/opt/cursor/artifacts", "inventory"}
	found := []string{}
	for _, root := range roots {
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}
			if len(found) >= 24 {
				return filepath.SkipAll
			}
			found = append(found, path)
			return nil
		})
	}
	k.Publish("reviewer", "watch", "listed artifacts", map[string]any{"n": len(found)})
	return kernel.Result{OK: true, Message: "reviewer listed " + itoa(len(found)), Data: found}, nil
}

func architectInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	goal := payloadQuery(call, "goal")
	if goal == "" {
		goal = payloadQuery(call, "title")
	}
	if goal == "" {
		goal = "control plane + ScanApp concept"
	}
	if mentionsLiveClient(goal) {
		return kernel.Result{OK: false, Message: "architect will not plan live client production", Data: map[string]any{"goal": goal}}, nil
	}
	plan := []string{
		"Keep the Go kernel. Do not rewrite Voltron.",
		"Registry + Fastify API is the front door.",
		"Onboard ScanApp after Evolu-Jeunes access. No guessed repo path.",
		"Depth limit 3. Hard budget $2/run.",
		"Live WP/Azure clients stay untouched.",
	}
	k.Publish("architect", "plan", goal, map[string]any{"steps": len(plan)})
	return kernel.Result{OK: true, Message: "architect planned " + goal, Data: map[string]any{"goal": goal, "steps": plan}}, nil
}

func deployInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	target := strings.ToLower(payloadQuery(call, "target") + " " + payloadQuery(call, "title"))
	if mentionsLiveClient(target) || strings.Contains(target, "prod") {
		return kernel.Result{OK: false, Message: "deploy refuses production / live clients", Data: map[string]any{"target": target}}, nil
	}
	_, err := os.Stat("bin/cashtro")
	data := map[string]any{
		"localOnly": true,
		"binary":    err == nil,
		"note":      "local rebuild only. no Azure/Vercel/WP promote.",
	}
	k.Publish("deploy", "release", "local check", data)
	return kernel.Result{OK: true, Message: "deploy local-only", Data: data}, nil
}

func securityInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	flags := []string{}
	if _, err := os.Stat("docs/ACCESS_REQUIRED.md"); err == nil {
		flags = append(flags, "evolu-jeunes-access-denied")
	}
	if _, err := os.Stat(".env"); err == nil {
		flags = append(flags, "dotenv-present-do-not-commit")
	}
	k.Publish("security", "triage", "local board", map[string]any{"flags": flags})
	return kernel.Result{OK: true, Message: "security flags " + itoa(len(flags)), Data: flags}, nil
}

func investigatorInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	ev := k.Events()
	if len(ev) > 12 {
		ev = ev[len(ev)-12:]
	}
	about := k.About()
	k.Publish("investigator", "trace", "kernel journal", map[string]any{"events": len(ev)})
	return kernel.Result{OK: true, Message: "investigator traced " + itoa(len(ev)) + " events", Data: map[string]any{
		"about":  about,
		"events": ev,
	}}, nil
}

func allowedLocalURL(u string) bool {
	u = strings.ToLower(strings.TrimSpace(u))
	return strings.HasPrefix(u, "http://127.0.0.1:") || strings.HasPrefix(u, "http://localhost:")
}

func mentionsLiveClient(s string) bool {
	low := strings.ToLower(s)
	for _, tok := range liveClientTokens {
		if strings.Contains(low, tok) {
			return true
		}
	}
	return false
}

func itoa(n int) string {
	return strings.TrimSpace(strings.ReplaceAll(jsonNumber(n), ".0", ""))
}

func jsonNumber(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}
