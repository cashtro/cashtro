# ScanApp — concept work (in this control plane)

ScanApp is a **concept** product, not a live client site. Castro asked
for work here, not just a name in the inventory report.

Repo URL is unknown until Evolu-Jeunes GitHub access lands. This file
is the concept contract the registry seeds from. Do not invent a repo.

## What it is

Scan, CRM, and AI bots. Stack: TypeScript + Python. Owner: Evolu-Jeunes.
Sector: ops. Stage: concept.

## What is in the registry

| Slug | `scanapp` |
| --- | --- |
| Kind | product |
| Status | concept |
| Live client | no — safe to onboard |
| Capabilities (draft) | `scan.ingest`, `crm.upsert`, `bot.reply` |

## Backlog on this project

1. `scanapp-scope` — lock the three verbs (scan / CRM / bots)
2. `scanapp-manifest` — draft `agent.manifest.json` (handshake)
3. `scanapp-find-repo` — blocked on Option A PAT, then `pnpm recon`
4. `ship-scanapp` — attach the real repo and open the onboard PR

## Do not

- Treat ScanApp as production
- Touch BTK, MD Clinic, Hypothèque, Éduconnexion, Proximity
- Guess a GitHub path (`Evolu-Jeunes/scanapp` is not confirmed)
