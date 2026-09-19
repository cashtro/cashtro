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

## Delivery catalog (Go)

This repo now includes a stdlib Go service that models the same line:
**idea → concept → production**. Seed data is the live mandate list above.
The board is a small HTML UI plus a JSON API. No frameworks.

```bash
go test ./...
go run ./cmd/cashtro
```

Then open `http://localhost:8080`.

| Method | Path | What it does |
| --- | --- | --- |
| `GET` | `/` | Delivery board |
| `GET` | `/health` | Liveness |
| `GET` | `/api/profile` | Public builder card |
| `GET` | `/api/ships` | Mandates on the line |
| `POST` | `/api/ships` | Park a new mandate in idea |
| `POST` | `/api/ships/{id}/advance` | Move one stage right |

Production is the end of the board. A further advance returns HTTP 409.
