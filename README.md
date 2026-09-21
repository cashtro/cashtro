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

This repo is **under construction**. It is the control plane for every
agentic we build here. Agents are processes. Capabilities are verbs.
Mail, notes, memory, human confirms, and desk requests are first-class.
You pick the next job. Symbols is the internal Zapier. Watch keeps work
flowing while the desk is closed. New agentics register into this kernel
— they do not become a second product.

```bash
go test ./...
go run ./cmd/cashtro
```

Open `http://localhost:8080`. Chrome: **Desk · You pick · Symbols · Line · Library**.

### Do we need Zapier?

**No.** Symbols is the internal Zapier. Triggers and kernel verbs live on
this OS. Webhook catch-hooks are `POST /api/symbols/{id}/hook`. Connectors
are the process table. There is no subscription and no second product.

### Do we need OpenRouter?

**No, not to boot.** The kernel, process table, delivery board, and journal
run with zero model keys.

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
| `GET` | `/api/inbox` | You-pick jobs (mail + verbs) |
| `POST` | `/api/inbox/{id}/take` | Take a job (mail parks on idea) |
| `POST` | `/api/inbox/{id}/skip` | Skip a job |
| `GET` | `/api/watch` | Closed-hours card |
| `POST` | `/api/watch/close` | Close desk · work still flowing |
| `POST` | `/api/watch/open` | Open desk |
| `POST` | `/api/watch/pulse` | Heartbeat |
| `GET` | `/api/plan` | Can / cannot / create |
| `GET` | `/api/requests` | Saved asks (original never overwritten) |
| `POST` | `/api/requests` | Capture & better an ask |
| `POST` | `/api/ships` | Park a mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |
| `GET` | `/api/symbols` | Internal Zapier table |
| `POST` | `/api/symbols` | Compose a symbol |
| `POST` | `/api/symbols/{id}/fire` | Run it now |
| `POST` | `/api/symbols/{id}/hook` | Catch-hook (webhook trigger) |
| `GET` | `/api/runs` | Symbol fire history |
| `GET` | `/api/connectors` | Kernel verbs you can wire |
