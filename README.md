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
agentic we build here. Agents are processes. Processes wear builders —
Lautaro, Frida, Diego. Products keep their names. Capabilities are verbs.
Mail, notes, memory, and human confirms are first-class. Delivery and
research are live. New agentics register into this kernel — they do not
become a second product.

```bash
go test ./...
go run ./cmd/cashtro
```

Open `http://localhost:8080`.

Casa Crypto, Prolifik, ScanApp, and Proximity stay **products** on the
delivery line. The process table wears builders and creators — Lautaro
boots the house. IDs and capability verbs do not change.

```mermaid
flowchart TB
  subgraph products [Delivery line — products keep their names]
    CasaCrypto[Casa Crypto]
    Prolifik[Prolifik]
    ScanApp[ScanApp]
    Proximity[Proximity]
  end
  subgraph builders [Process table — builders]
    Lautaro[Lautaro · init]
    Hedy[Hedy · router]
    Antoni[Antoni · delivery]
    Gabriela[Gabriela · research]
    Marco[Marco · explorer]
    Diego[Diego · operator]
    Juana[Juana · reviewer]
    Frida[Frida · architect]
    Oscar[Oscar · deploy]
    Tupac[Túpac · security]
    Luis[Luis · memory]
    Violeta[Violeta · comms]
    Pablo[Pablo · planner]
    Caupolican[Caupolicán · investigator]
  end
  Lautaro --> Antoni
  Marco --> Gabriela
  Pablo --> Antoni
  Antoni --> CasaCrypto
  Antoni --> Prolifik
  Antoni --> ScanApp
  Antoni --> Proximity
```

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
| `GET` | `/api/ships` | Delivery line |
| `POST` | `/api/ships` | Park a mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |
