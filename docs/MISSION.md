# MISSION — Agentic Control Plane

Filled 2026-09-21 by the lead engineer (autonomy). Placeholders the brief
left blank are recorded here and in `docs/DECISIONS.md`.

```
CONTROL_PLANE_REPO   = cashtro/cashtro
GITHUB_ORGS          = Evolu-Jeunes
GITHUB_USERS         = cashtro
MODEL_ROUTE          = both: Ollama internal (llama3.2) + OpenRouter external (moonshotai/kimi-k3 and z-ai/glm-5.3 reasoning max)
HARD_BUDGET_PER_RUN  = $2.00
DEPLOY_TARGET        = local only
DB                   = SQLite (dev) → Postgres (prod), same Prisma schema
```

This repo is the manager. The Go kernel (Cashtro OS / Voltron) stays.
The TypeScript control plane is the front door other agents call.

## Autonomy

Decide layout, libraries, tests, branches, PRs, refactors, retries.
Stop only for: spend over $2/run, delete/force-push main, wider-than-read
secrets, paid third-party services, live client production sites.

## Phases in this brief

- **0** Bootstrap & access audit — this folder + `state/backlog.json`
- **1** Recon CLI — `pnpm recon`
- **2** Registry — Prisma + `packages/sdk`
- **3** Control plane API — `apps/api`
- 4–7 are out of scope until recon exists and Castro reads REPORT.md

## Stack

TypeScript / Node 20+. REST + OpenAPI 3.1 from code. BullMQ when
`REDIS_URL` is set; SQLite job table otherwise (no paid Redis).
Prisma. Cheap models for recon/classify; expensive only for plan/review.
Registry exposed later as MCP. No secrets in git.
