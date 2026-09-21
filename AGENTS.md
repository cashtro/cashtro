# AGENTS — cashtro/cashtro

This repo is the **manager**. The Go kernel (Cashtro OS / Voltron) stays.
The TypeScript control plane is the front door.

## Do

- Absorb the kernel. Do not rewrite it.
- Work on `cursor/<name>-b1d0` branches. Never force-push main.
- Record decisions in `docs/DECISIONS.md`. Progress in `docs/PROGRESS.md`.
- Hard budget $2 / run. Pause drain before new dispatch.
- Evolu-Jeunes private repos: do not guess names. See `docs/ACCESS_REQUIRED.md`.

## Do not

- Touch live client production (BTK, MD Clinic, Hypothèque, Éduconnexion, Proximity).
- Commit secrets. Values stay in env / vault. SecretRefs are names only.
- Spend over the cap. Do not add paid third-party services without Castro.

## Handshake

`agent.manifest.json` is how this repo joins its own network.
`pnpm onboard` drafts a manifest + AGENTS.md for a row in `inventory/repos.json`.

## Kernel seats

14 agentics on `:8080`. Live: init, delivery, research, explorer, memory, comms, planner.
Resident: operator, architect, deploy, reviewer, security, investigator, router.
