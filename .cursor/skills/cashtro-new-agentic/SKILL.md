---
name: cashtro-new-agentic
description: Add or promote a Cashtro OS agentic (process + capability) to ModeLive with deterministic invoke, optional REST surface, tests, and version bump. Use when adding an agent, making a resident agent live, binding a new capability, or the user says new agentic / live capability / under construction ship one verb.
---

# Ship one Cashtro agentic

New agentics register into **this** kernel. They do not become a second product.

## Pattern (follow existing live agentics)

1. **Spec** — in `internal/agents/builtin.go` `Builtins`, use `resident(Spec{..., Mode: ModeLive, Capabilities: []string{"agent.verb"}, Autostart: true}, invokeFn)`.
2. **Invoke** — implement `invokeFn` in `internal/agents/live.go`. Prefer deterministic Go (notes, memory, mail, confirms, catalog). Model calls go through `router` only when optional OpenRouter is bound.
3. **Side effects** — write durable evidence when useful: `WriteNote`, `Remember`, `Post`, `RequestConfirm`, then rely on `persist()` / autosave.
4. **HTTP (optional)** — if the desk needs a shortcut, add `HandleFunc` in `internal/server/server.go` mirroring `/api/browse`, `/api/deploy`, `/api/review`. Document in `README.md`.
5. **Tests** — extend `internal/agents/builtin_test.go` and `internal/server/server_test.go`. Keep `about.Live` / process count assertions truthful.
6. **Version** — bump `kernel.Version` in `internal/kernel/kernel.go` (semver patch for one capability).
7. **Research crumb** — optional overnight note via skill `cashtro-research-note`.
8. **Verify + ship** — `go test ./...`, restart always-on onto the new binary, commit/push, update PR.

## Live vs resident

| Mode | Meaning |
| --- | --- |
| `ModeLive` | Capability executes useful work today (dry-run OK) |
| `ModeResident` | Process on the desk; waiting for a worker/key bind |

`router` stays **resident** until `OPENROUTER_API_KEY` is set. Do not block OS work on a model key.

## Dry-run is fine

Operator/deploy/reviewer-style dry-runs that plan, gate, or score desk evidence are valid `ModeLive` work until a real worker binds.

## Do not

- Fork a second OS repo or parallel control plane
- Add OpenRouter as a hard boot dependency
- Register an agent without Autostart when it should appear on the desk at boot
