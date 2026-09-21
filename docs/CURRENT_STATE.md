# CURRENT_STATE — cashtro/cashtro

Mapped 2026-09-21. Do not rewrite this kernel. Absorb it.

## What is here (alive)

Stdlib **Go 1.22** process OS. One binary: `cmd/cashtro` → HTTP `:8080`.

| Path | Role |
| --- | --- |
| `internal/kernel` | Process table, invoke bus, journal, mail, notes, memory, confirms |
| `internal/agents` | 14 registered agentics. 13 live local verbs, router resident until a key |
| `internal/catalog` | Delivery line: idea → concept → production (in-memory seed) |
| `internal/model` | Optional OpenRouter client. Boots unbound. |
| `internal/server` | Desk HTML + REST (`/api/os`, `/api/agents`, `/api/ships`, …) |
| `.github/workflows/go.yml` | `go test` + `go vet` on main/PR |
| `.cursor/environment.json` | Cloud start: `go run ./cmd/cashtro` on `main` |

Seeded ships (products, not processes): BTK Avocats, MD Clinic,
Solution Hypothèque QC, Éduconnexion, Proximity (production);
ScanApp (concept); Cashtro delivery catalog (idea).

Live verbs today: `os.about`, delivery CRUD/advance, explorer search,
memory store/recall, comms confirm gate, planner backlog, research ingest,
`model.status` / `model.chat` (chat no-ops without a key).

## What is dead or missing

- Router stays resident without OPENROUTER_API_KEY
- Kernel state is in-process. Restart wipes ships/notes/mail
- Evolu-Jeunes private fleet still invisible to this Cloud Agent token
- No auth on the desk. Fine for local.

## Absorb vs replace

| Piece | Decision |
| --- | --- |
| Go kernel + desk | **Absorb.** Keep as the local Voltron runtime / process OS |
| Delivery catalog seed | **Absorb** into registry Projects (kind=product) |
| 14 agentics | **Absorb** into registry Agents (runtime=http → kernel invoke) |
| OpenRouter bind | **Absorb** as MODEL_ROUTE |
| New control plane API | **Add.** Does not fork a second OS |
| Per-island agents in client repos | **Replace** with registry + manifest handshake (Phase 5) |

## What the manager can already control

Only this public repo. Evolu-Jeunes (claimed 48 repos) is invisible to
the current GitHub token. See `docs/ACCESS_REQUIRED.md`.
