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

GitHub access: Castro offered. Yes. Option A in `docs/ACCESS_REQUIRED.md`.
Do not paste the token in chat.

## Phase 4–7

Not started. Waiting on Evolu-Jeunes PAT, then ScanApp onboard.
