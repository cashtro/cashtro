# Cashtro Builder

Scrum Master & Software Engineer · Founder · FR/EN · Remote worldwide

**I make teams ship: idea → concept → production.**

The delivery operation lives at [**Evolu-Jeunes**](https://github.com/Evolu-Jeunes) —
48 repos, 30+ platforms shipped for law firms, clinics, finance and
e-commerce. Client code stays private. The results are live:

- ⚖️ BTK Avocats — law firm platform (WordPress)
- 🏥 MD Clinic — medical clinic (WordPress/ACF)
- 🏦 Solution Hypothèque QC — mortgage services
- 🎓 Éduconnexion — education platform
- 🚀 Proximity — agency platform (Next.js/TypeScript + Azure APIs)
- 🤖 ScanApp, CRM, AI bots — TypeScript & Python builds

**Stack:** Go · PHP/WordPress · React/Next.js/TypeScript · Python · Azure

**Team:** 75+ developers trained on real client mandates

📩 alejandro@proximityagency.ca

---

## Cashtro OS

**Repo:** [github.com/cashtro/cashtro](https://github.com/cashtro/cashtro)

This repository **is** Cashtro OS. Under construction. The GitHub App
on this agent can only write this repo, so the kernel lives here — not
in a second empty `cashtro-os` clone.

Agents are processes. Capabilities are verbs. Mail, notes, memory, and
human confirms are first-class. Delivery, research, **Architect.plan**,
**Investigator.trace**, **Security.triage**, and **Reviewer.watch** are
live without a model. Watch keeps the desk flowing while things are
closed. Explorer searches pulses and mail. The desk shows the mailbox and memory. `/health` reports the last pulse. The disk image is `data/cashtro.json`. New
agentics register into this kernel — they do not fork a second product.

```bash
go test ./...
go run ./cmd/cashtro
# optional:
# go run ./cmd/cashtro -data data/cashtro.json
# docker build -t cashtro-os . && docker run -p 8080:8080 cashtro-os
```

Open `http://localhost:8080`.

### Do we need OpenRouter?

**No, not to boot.** The kernel, process table, delivery board, journal,
and research library run with zero model keys.

**Yes, if an agentic needs to think.** OpenRouter is the optional model bus:
one key, many models. Bind it when you want `model.chat` to execute.

```bash
export OPENROUTER_API_KEY=sk-or-...
export OPENROUTER_MODEL=openai/gpt-4o-mini   # optional
go run ./cmd/cashtro
```

Without a key the `router` process stays **resident** and tells you it is
unbound. The rest of the OS stays live.

| Method | Path | What it does |
| --- | --- | --- |
| `GET` | `/` | OS desk |
| `GET` | `/health` | Liveness + kernel counts |
| `GET` | `/api/os` | Manifesto and process counts |
| `GET` | `/api/agents` | Process table |
| `POST` | `/api/agents/{id}/invoke` | Run a capability |
| `GET` | `/api/model` | OpenRouter bind card |
| `GET` | `/api/notes` | Research library |
| `POST` | `/api/notes` | Ingest a sourced finding |
| `GET` | `/api/mail` | Agent mailbox |
| `GET` | `/api/memory` | Episodic recall |
| `GET` | `/api/confirms` | Human gate |
| `GET` | `/api/ships` | Delivery line |
| `POST` | `/api/ships` | Park a mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |
| `GET` | `/api/watch` | Closed-hours status and pulses |
| `POST` | `/api/watch/close` | Close the desk; work keeps flowing |
| `POST` | `/api/watch/pulse` | Heartbeat (parks work while closed) |
| `POST` | `/api/watch/open` | Open the desk |

### Always on (laptop closed)

This Cloud Agent runs on a remote VM. Closing your computer does **not**
stop Cashtro OS while the VM is up.

```bash
make always-on
# or: ./scripts/always-on.sh
# overnight:
# CASHTRO_CLOSED=1 make always-on
```

That loop rebuilds, starts `:8080`, and restarts on crash. State lands in
`data/cashtro.json` — notes, ships, pulses, **and the journal**. `Watch`
is a live system process: close the desk and
it pulses, persists, parks a closed-hours mandate, and keeps research /
delivery / memory moving. While the process is up it also heartbeats
every 30m (`CASHTRO_PULSE_EVERY`, or `0` to disable). Outbound
`comms.send` still waits at the human gate. Overnight wake-ups can also
be scheduled as Cursor timers on this agent conversation.

