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

### Dual loop · Kimi K3 + GLM

Pulse never stops. **Forge** rotates Castro's wealth playbook
(create → build → grow → evolve). Each tick must raise the score.
OpenRouter stays optional:

- Unbound: the loop is deterministic. Processes stay up. Offers compound.
- Bound: `forge.think` sends the prompt to **Kimi K3** (`moonshotai/kimi-k3`)
  then **GLM-4.5** (`z-ai/glm-4.5`). The kernel still enforces a rising score
  and the human gate.

```bash
export OPENROUTER_API_KEY=sk-or-...
# optional seat overrides
export OPENROUTER_KIMI_MODEL=moonshotai/kimi-k3
export OPENROUTER_GLM_MODEL=z-ai/glm-4.5
go run ./cmd/cashtro
```

`GET /api/wealth` · `GET /api/forge` · `POST /api/pulse/tick`

Without a key the `router` process stays **resident** and tells you it is
unbound. Pulse, wealth, delivery, and research stay live.

| Method | Path | What it does |
| --- | --- | --- |
| `GET` | `/` | OS desk |
| `GET` | `/health` | Liveness + kernel counts |
| `GET` | `/api/os` | Manifesto, process counts, generation, score |
| `GET` | `/api/agents` | Process table |
| `POST` | `/api/agents/{id}/invoke` | Run a capability |
| `GET` | `/api/model` | OpenRouter bind card (Kimi + GLM seats) |
| `GET` | `/api/pulse` | Never-stop heartbeat |
| `POST` | `/api/pulse/tick` | Keep every autostart running + forge one generation |
| `GET` | `/api/forge` | Generation history |
| `GET` | `/api/wealth` | Offer ledger (create / build / grow / evolve) |
| `GET` | `/api/notes` | Research library |
| `POST` | `/api/notes` | Ingest a sourced finding |
| `GET` | `/api/mail` | Agent mailbox |
| `GET` | `/api/memory` | Episodic recall |
| `GET` | `/api/confirms` | Human gate |
| `GET` | `/api/ships` | Delivery line |
| `POST` | `/api/ships` | Park a mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |
