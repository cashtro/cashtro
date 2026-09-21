---
name: cashtro-os-verify
description: Smoke-test Cashtro OS health, process table, and key APIs after a change or keep-alive. Use when verifying the desk, before/after always-on restart, after go test, or when the user asks if the OS is up / healthy / live.
---

# Verify Cashtro OS

## Fast path

```bash
go test ./...
curl -sf http://127.0.0.1:8080/health
curl -sf http://127.0.0.1:8080/api/health
curl -sf http://127.0.0.1:8080/api/os
curl -s -X POST http://127.0.0.1:8080/api/heartbeat
```

Assert: `"status":"ok"`, `"always":true`, `version` matches `kernel.Version`, `live` count matches expectations (usually 13 with unbound router), `data` is `/workspace/data/cashtro.json` (or configured path).

## Capability smokes (as needed)

| Check | Command |
| --- | --- |
| Agents | `curl -s http://127.0.0.1:8080/api/agents` |
| Trace | `curl -s 'http://127.0.0.1:8080/api/trace?q=AOS'` |
| Security | `curl -s 'http://127.0.0.1:8080/api/security?q=AOS'` |
| Deploy | `curl -s -X POST http://127.0.0.1:8080/api/deploy -d '{"target":"cashtro-os"}'` |
| Review | `curl -s -X POST http://127.0.0.1:8080/api/review -d '{"target":"cashtro-os"}'` |
| Browse | `curl -s -X POST http://127.0.0.1:8080/api/browse -d '{"url":"http://127.0.0.1:8080/"}'` |
| Notes | `curl -s http://127.0.0.1:8080/api/notes` |
| Model | `curl -s http://127.0.0.1:8080/api/model` — unbound is OK |

## Always-on process

```bash
tmux -f /exec-daemon/tmux.portal.conf has-session -t '=cashtro-always-on'
pgrep -af '/workspace/bin/cashtro'
test -f /workspace/data/cashtro.json
```

If health fails, use skill `cashtro-always-on` before debugging application code.

## Output

Report a one-line verdict: version, live/resident counts, always-on up or restarted, tests pass/fail.
