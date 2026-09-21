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
Mail, notes, memory, confirms, and **requests** are first-class. The
desk captures every ask and betters it each pass. Delivery and research
are live. New agentics register into this kernel — they do not become a
second product.

```bash
go test ./...
go run ./cmd/cashtro
```

Open `http://localhost:8080`.

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
| `GET` | `/api/plan` | What Castro can and cannot ask here |
| `GET` | `/api/requests` | Desk log — original + improved asks |
| `POST` | `/api/requests` | Capture an ask and better it (pass 1) |
| `POST` | `/api/requests/{id}/better` | Tighten the same ask again |
| `POST` | `/api/requests/{id}/done` | Park it as done |
| `GET` | `/api/ships` | Delivery line |
| `POST` | `/api/ships` | Park a mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |

### Desk assistant

Ask in your own words (FR/EN). The **Desk** agentic:

1. Saves the original text (never overwritten)
2. Rewrites it into a shippable task with owner, verb, and acceptance
3. Betters the same request on every later pass
4. Publishes a can / cannot / create / do-not-create plan on the OS desk

I execute every allowed task in this kernel. Blocked asks stay on the
desk with the reason. Open `AGENTS.md` for the operating contract.
