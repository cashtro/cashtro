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
Mail, notes, memory, and human confirms are first-class. Delivery and
research are live. New agentics register into this kernel — they do not
become a second product.

```bash
go test ./...
go run ./cmd/cashtro
```

Open `http://localhost:8080`.

### Voltron (overnight, non-stop)

Voltron is this kernel assembled: every registered agentic under one
supervisor. On this machine it must stay up. The supervisor rebuilds
and restarts Cashtro OS if the process dies.

```bash
make voltron          # start and keep restarting
make voltron-status   # /health probe
```

Logs and heartbeat: `/tmp/voltron/voltron.log`, `/tmp/voltron/status.json`.

### Model route — both lanes

The router uses **both** lanes on `model.chat` with `"route":"both"` (the default).

| Lane | Provider | Models |
| --- | --- | --- |
| internal | Ollama | `llama3.2` (`OLLAMA_MODEL`) |
| external | OpenRouter | `moonshotai/kimi-k3` and `z-ai/glm-5.3` at reasoning `max` |

```bash
ollama serve
export OLLAMA_MODEL=llama3.2
export OPENROUTER_API_KEY=sk-or-...
export OPENROUTER_MODEL=moonshotai/kimi-k3
export OPENROUTER_ALSO_MODEL=z-ai/glm-5.3
go run ./cmd/cashtro
```

The kernel still boots with neither daemon nor key. The `router` stays
**resident** until Ollama answers or `OPENROUTER_API_KEY` is set.

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
| `GET` | `/api/teal` | Corporate Teal brain (Teams only) |
| `GET` | `/api/vapi` | Vapi (Vappy) voice bind + share link |
| `POST` | `/api/vapi/webhook` | Vapi end-of-call / transcript → Teal listen |
| `POST` | `/api/teal/scrape` | Graph scrape of Teams AI Teal (needs `TEAMS_TOKEN`) |
| `POST` | `/api/teal/ingest` | Ingest a Teams message (refuses Slack/mail) |
| `GET` | `/api/teal/interns` | Intern desk |
| `POST` | `/api/teal/question` | Answer intern / project questions |
| `POST` | `/api/teal/cursor` | Park a Cursor Cloud Agent launch (human gate) |
| `GET` | `/api/teal/card` | Teams Adaptive Card (Vapi + Cursor) |

### Teal · Vapi · Teams only

Corporate **AI Teal** lives on Microsoft Teams. Voice is **Vapi** (`https://vapi.ai/`, dashboard `https://dashboard.vapi.ai/`). The Cursor app launches from the Teams card (`https://cursor.com/agents`). New intel agentics spawn back into this kernel.

```bash
export VAPI_SHARE_URL=https://dashboard.vapi.ai/assistants/YOUR_ASSISTANT_ID
export VAPI_ASSISTANT_ID=YOUR_ASSISTANT_ID
# optional: export VAPI_API_KEY=...
# optional Graph scrape: export TEAMS_TOKEN=...
go run ./cmd/cashtro
# Point the Vapi assistant server URL at /api/vapi/webhook
```

## Control plane

TypeScript front door for the fleet. Does not replace this kernel.

```bash
pnpm recon            # inventory every reachable repo
pnpm db:migrate && pnpm db:seed
pnpm api              # http://127.0.0.1:8787/docs
```

See `docs/MISSION.md`, `inventory/REPORT.md`, `docs/ACCESS_REQUIRED.md`.
