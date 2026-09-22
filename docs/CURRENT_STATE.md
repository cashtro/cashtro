# CURRENT_STATE — cashtro/cashtro

Mapped 2026-09-21. Do not rewrite this kernel. Absorb it.

## What is here (alive)

Stdlib **Go 1.22** process OS. One binary: `cmd/cashtro` → HTTP `:8080`.

| Path | Role |
| --- | --- |
| `internal/kernel` | Process table, invoke bus, journal, mail, notes, memory, confirms |
| `internal/agents` | 17 registered agentics. 10 live (incl. Teal + Vapi), 7 resident |
| `internal/vapi` | Voice catalog + Chat/Call client. Talk local without a key. |
| `internal/catalog` | Delivery line: idea → concept → production (in-memory seed) |
| `internal/model` | Optional OpenRouter client. Boots unbound. |
| `internal/server` | Desk HTML + REST (`/api/os`, `/api/agents`, `/api/ships`, …) |
| `.github/workflows/go.yml` | `go test` + `go vet` on main/PR |
| `teal/` | Fresh Evolu-Jeunes Teal product (`github.com/Evolu-Jeunes/Teal`) on `:8090` |

Seeded ships (products, not processes): BTK Avocats, MD Clinic,
Solution Hypothèque QC, Éduconnexion, Proximity (production);
ScanApp (concept); Cashtro delivery catalog (idea).

Live verbs today: `os.about`, delivery CRUD/advance, explorer search,
memory store/recall, comms confirm gate, planner backlog, research ingest,
`model.status` / `model.chat` (chat no-ops without a key), Teal corporate
brain (`teal.*`) — Teams scrape/ingest, intern desk, idea, self-spawn,
Cursor card. Vapi voice lane (`vapi.*`) — talk, voice switch, human-gated
outbound, bridges across every business line.

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
| 17 agentics (incl. Manager, Teal, Vapi) | **Absorb** into registry Agents (runtime=http → kernel invoke) |
| OpenRouter bind | **Absorb** as MODEL_ROUTE |
| New control plane API | **Add.** Does not fork a second OS |
| Per-island agents in client repos | **Replace** with registry + manifest handshake (Phase 5) |

## What the manager can already control

Only this public repo. Evolu-Jeunes (claimed 48 repos) is invisible to
the current GitHub token. See `docs/ACCESS_REQUIRED.md`.
