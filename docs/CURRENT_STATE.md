# CURRENT_STATE — cashtro/cashtro

Mapped 2026-09-21. Do not rewrite this kernel. Absorb it.

## What is here (alive)

Stdlib **Go 1.22** process OS. One binary: `cmd/cashtro` → HTTP `:8080`.

| Path | Role |
| --- | --- |
| `internal/kernel` | Process table, invoke bus, journal, mail, notes, memory, confirms |
| `internal/agents` | 15 registered agentics. 8 live, 7 resident (contract only) |
| `internal/catalog` | Delivery line: idea → concept → production (in-memory seed) |
| `internal/model` | Optional OpenRouter client. Boots unbound. |
| `internal/flow` | Owned n8n graphs, catch-hooks, run history |
| `internal/think` | Owned propose/critique dual. No Moonshot/GLM keys |
| `internal/server` | Desk HTML + REST (`/api/os`, `/api/agents`, `/api/ships`, …) |
| `.github/workflows/go.yml` | `go test` + `go vet` on main/PR |
| `.cursor/environment.json` | Cloud start: `go run ./cmd/cashtro` on `main` |

Seeded ships (products, not processes): BTK Avocats, MD Clinic,
Solution Hypothèque QC, Éduconnexion, Proximity (production);
ScanApp (concept); Cashtro delivery catalog (idea).

Live verbs today: `os.about`, delivery CRUD/advance, explorer search,
memory store/recall, comms confirm gate, planner backlog, research ingest,
`n8n.workflow` (owned graphs, no vendor keys),
`model.status` / `model.chat` (chat no-ops without a key).

## What is dead or missing

- No AGENTS.md, CLAUDE.md, `.mcp.json`, `agent.manifest.json`
- No TypeScript workspace, no Prisma, no queue, no cost ledger
- No GitHub inventory, no multi-repo registry
- Resident agentics (operator, reviewer, architect, deploy, security,
  investigator, unbound router) accept invoke and return “ready to bind”
- Kernel state is in-process. Restart wipes ships/notes/mail
- No auth on the desk. Fine for local. Not a fleet front door
- `scripts/voltron.sh` exists only on branch `cursor/voltron-overnight-b1d0`,
  not on `main`

## Absorb vs replace

| Piece | Decision |
| --- | --- |
| Go kernel + desk | **Absorb.** Keep as the local Voltron runtime / process OS |
| Delivery catalog seed | **Absorb** into registry Projects (kind=product) |
| 15 agentics | **Absorb** into registry Agents (`n8n` runtime=n8n, others http → kernel invoke) |
| OpenRouter bind | **Absorb** as MODEL_ROUTE |
| New control plane API | **Add.** Does not fork a second OS |
| Per-island agents in client repos | **Replace** with registry + manifest handshake (Phase 5) |

## What the manager can already control

Only this public repo. Evolu-Jeunes (claimed 48 repos) is invisible to
the current GitHub token. See `docs/ACCESS_REQUIRED.md`.
