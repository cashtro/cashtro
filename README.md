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
human confirms are first-class. Delivery and research are live. The
disk image is `data/cashtro.json`. New agentics register into this
kernel — they do not fork a second product.

```bash
go test ./...
go run ./cmd/cashtro
# optional:
# go run ./cmd/cashtro -data data/cashtro.json
# docker build -t cashtro-os . && docker run -p 8080:8080 cashtro-os
```

Open `http://localhost:8080`.

---

## Ultron IDE (company control plane)

**Separated infrastructure** above Cashtro OS. Run Giant, ScanApp, Proximity,
Empire and every company from one authorized IDE desk. Clients authenticate
with RBAC. The server does **not** need a Cursor token.

```bash
go run ./cmd/ultron                 # :9090
# or always-on:
./scripts/ultron-on.sh
# both Cashtro + Ultron:
./scripts/fleet-on.sh
```

Open `http://localhost:9090`.

Default owner (first boot only):

- email: `alejandro@proximityagency.ca`
- password: `ultron-change-me` (override with `ULTRON_OWNER_PASSWORD`)

| Layer | Port | Role |
| --- | --- | --- |
| **Ultron** | `:9090` | IDE · companies · fleet · RBAC · bridges Cashtro |
| **Cashtro OS** | `:8080` | Agentic kernel · delivery · research · disk |

Disk image: `data/ultron.json`. See `.cursor/skills/ultron-ide/SKILL.md`.

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
| `GET` | `/health` | Liveness + kernel counts + uptime |
| `GET` | `/api/health` | Same as `/health` |
| `POST` | `/api/heartbeat` | Always-on keep-alive stamp |
| `GET` | `/api/os` | Manifesto and process counts |
| `GET` | `/api/agents` | Process table |
| `POST` | `/api/agents/{id}/invoke` | Run a capability |
| `GET` | `/api/model` | OpenRouter bind card |
| `GET` | `/api/notes` | Research library |
| `POST` | `/api/notes` | Ingest a sourced finding |
| `GET` | `/api/mail` | Agent mailbox |
| `GET` | `/api/memory` | Episodic recall |
| `POST` | `/api/memory` | Store an episodic fact |
| `GET` | `/api/confirms` | Human gate |
| `GET` | `/api/trace?q=` | Investigator blast radius |
| `GET` | `/api/security?q=` | Security triage |
| `POST` | `/api/deploy` | Release dry-run gate |
| `POST` | `/api/review` | QA desk-evidence review |
| `POST` | `/api/browse` | Operator browse dry-run |
| `POST` | `/api/plan` | Planner backlog → idea ships |
| `POST` | `/api/architect` | Architect plan → note + idea ships |
| `GET` | `/api/ships` | Delivery line |
| `POST` | `/api/ships` | Park a mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |

### Always on (laptop closed)

This Cloud Agent runs on a remote VM. Closing your computer does **not**
stop Cashtro OS while the VM is up.

```bash
make always-on
# or: ./scripts/always-on.sh
```

That loop rebuilds, starts `:8080`, and restarts on crash. State lands in
`data/cashtro.json`. The kernel also autosaves that image every 60s so a
hard kill still leaves a durable disk. Overnight wake-ups can also be
scheduled as Cursor timers on this agent conversation.

