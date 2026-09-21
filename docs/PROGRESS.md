# PROGRESS

## Phase 0 — 2026-09-21

Done.

- Mapped `cashtro/cashtro` into `docs/CURRENT_STATE.md`
- Seeded `state/backlog.json` from README, catalog, kernel, mission
- `gh auth status` is the Cursor integration account, not Castro
- Visible: 1 public repo (`cashtro/cashtro`)
- Invisible: Evolu-Jeunes fleet → `docs/ACCESS_REQUIRED.md` (Option A)
- Filled mission blanks (OpenRouter, $2/run, local only) in `docs/MISSION.md`

Proof:

```
gh repo list cashtro --limit 200 --json name,visibility,updatedAt
→ [{"name":"cashtro","visibility":"PUBLIC",...}]

gh api orgs/Evolu-Jeunes/repos?type=all
→ []
```

## Phase 1 — 2026-09-21

Done.

```
pnpm --filter @cashtro/recon test
→ 3 pass (idempotent merge + catalog-in-report + schema reject)

pnpm recon
→ recon ok · 1 repos · 2 rows   # 2026-09-21T22:03:51Z

pnpm recon
→ recon ok · 1 repos · 2 rows   # second run identical row count
```

Second run stayed at 2 rows (`cashtro/cashtro` + `Evolu-Jeunes/*` limited).
`inventory/REPORT.md` now also lists 7 kernel-catalog ships as **not scanned**.

## Phase 2 — 2026-09-21

Done.

```
DATABASE_URL=file:../../data/control-plane.db pnpm exec prisma migrate deploy
→ Applying migration 20260921153000_init · All migrations applied

DATABASE_URL=... pnpm exec tsx prisma/seed.ts
→ seeded cashtro + Evolu-Jeunes denied project, 14 kernel agents, SecretRefs, backlog tasks

pnpm --filter @cashtro/sdk test
→ 1 pass
```

## Phase 3 — 2026-09-21

Done.

```
pnpm --filter @cashtro/api test
→ 3 pass
  registry seed + every mutating route writes an event
  pause refuses dispatch
  budget cap blocks an over-budget dispatch
```

`GET /docs` and `GET /docs/json` render OpenAPI 3.1.
Kill switch and $2 budget cap are enforced before dispatch.

## ScanApp concept work — 2026-09-21

Castro: there was no ScanApp work in the control plane. Fixed.

- Registry now seeds `scanapp` as a concept project from `state/catalog-ships.json`
- Capabilities: `scan.ingest`, `crm.upsert`, `bot.reply`
- Tasks attached: scope, manifest, find-repo (blocked on PAT), onboard
- Contract: `docs/SCANAPP.md`
- Draft handshake: `inventory/drafts/scanapp/agent.manifest.json`

GitHub access: Castro connected Cursor to cashtro + Evolu-Jeunes.
Re-checked 22:34Z. This agent is still account `cursor` (ghs_).
`/installation/repositories` = **selected**, only `cashtro/cashtro`.
Evolu-Jeunes exists (id 201155686), `public_repos=0`, GraphQL `totalCount=0`.
Need Cursor GitHub App installed on the **org** with All repositories.
See `docs/ACCESS_REQUIRED.md`.

## Phase 4 — 2026-09-21

Adapters + orchestrator. Depth 3 hard stop. HTTP → Voltron.

```
pnpm --filter @cashtro/adapters test
pnpm --filter @cashtro/api test   # includes e2e if :8080 is up
```

## Phase 5 — 2026-09-21

`agent.manifest.json` + `AGENTS.md` on cashtro. `pnpm onboard cashtro/cashtro`.
Does not invent Evolu-Jeunes repo paths.

## Phase 6 — 2026-09-21

`GET /costs`, pause drains queued jobs, `pnpm secret-scan` + CI job.
Budget cap still blocks before dispatch.

## Phase 7 — 2026-09-21

Four screens: Fastify `/ui` (phone) and Next.js `apps/web`.
Evolu-Jeunes still dark — recon cannot grow until Cursor App is on the org.

## Bind seats — 2026-09-21

operator / reviewer / architect / deploy / security / investigator
execute locally. Refuse live clients and non-local URLs. Router stays
resident until a key. Voltron desk + `scripts/voltron.sh` restored onto
this branch so a rebuild does not drop the overnight supervisor.

## Overnight build — 2026-09-21 23:20Z

- MCP registry: `pnpm mcp` / `.mcp.json` tools list_inventory, list_catalog,
  list_backlog, get_access
- GitHub webhook opens idempotent tasks
- `/ui` queue can dispatch; fleet can pause

## n8n fabric — 2026-09-21 23:40Z

Castro: n8n.io inside the same agents, 40× stronger, bigger architecture.

- `state/n8n-fleet.json` — 40 specialists, 14 seats
- Adapter: unbound or webhook error → Voltron `:8080`
- `GET /fleet` · `/ui/fleet` grid · MCP `list_n8n_fleet`
- Optional `docker-compose.n8n.yml` on `:5678`. No paid cloud.
