---
name: cashtro-os-keepalive
description: Overnight keep-alive loop for Cashtro OS on the cloud VM when the laptop may be closed. Use when the cashtro-os-overnight timer fires, or when the user says keep building overnight, never stop, run all night, keep-alive, or under construction keep shipping. Covers health check, one durable improvement, commit/push, PR update, and leaving always-on running.
---

# Cashtro OS overnight keep-alive

Cashtro OS runs on a **cloud VM**. Closing the laptop does not stop it. Each wake-up lands **one small durable improvement**, then leaves the process running.

## When this applies

- Timer `cashtro-os-overnight` (~7200s) fires
- User says: run all night, never stop, laptop closed, keep building under construction
- Follow-up keep-alive with no new human ask

## Do this every wake

1. **Message queue** — call `cursor-cloud` `get-message-queue`. If other user messages are queued, finish the keep-alive quickly and defer long verification.
2. **Health** — `curl -sf http://127.0.0.1:8080/health` (also `/api/health`). Expect `"always":true` and a current `version`.
3. **Always-on** — if health fails, start `./scripts/always-on.sh` in tmux session `cashtro-always-on` (see skill `cashtro-always-on`). Do not ask the user to open their laptop.
4. **Disk** — confirm `data/cashtro.json` exists; version in `/health` matches `internal/kernel/kernel.go` `Version`.
5. **One durable improvement** — pick exactly one:
   - New/live agent capability (skill `cashtro-new-agentic`)
   - Persistence / API polish / test
   - Research note via `research.ingest` (skill `cashtro-research-note`)
   - Router stays **optional** — do not require OpenRouter
6. **Verify** — `go test ./...` then smoke health + the new path (skill `cashtro-os-verify`).
7. **Ship** — commit on `cursor/go-language-task-0b85`, `git push -u origin cursor/go-language-task-0b85`, update draft PR #2 with `ManagePullRequest` `update_pr`. Preserve any human edits already in the PR body.
8. **Heartbeat** — `curl -s -X POST http://127.0.0.1:8080/api/heartbeat`.
9. **End turn** — leave always-on running. Confirm timer `cashtro-os-overnight` is still subscribed via `cursor-subscriptions` `list_subscriptions`. Re-subscribe only if missing.

## Hard rules

- Do **not** kill tmux `cashtro-always-on` with broad `pkill`/`fuser` unless you immediately restart it.
- To pick up a new binary: kill only `/workspace/bin/cashtro` so the always-on loop rebuilds.
- Do **not** ask the user to open their laptop or re-run local servers.
- Prefer one focused commit over a kitchen-sink change.

## Timer prompt (canonical)

If re-creating the timer, use name `cashtro-os-overnight`, `delaySeconds: 7200`, and a prompt that mirrors steps 2–9 above.
