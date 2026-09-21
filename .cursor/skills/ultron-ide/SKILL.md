---
name: ultron-ide
description: Ultron personal IDE OS — multi-company control plane with RBAC above Cashtro. Use when running companies (Giant, ScanApp, Proximity, Empire), IDE dashboard, client auth, RBAC, no Cursor token on server, or bridging to Cashtro OS from Ultron.
---

# Ultron IDE (above Cashtro)

Ultron is a **separated** control plane. Cashtro OS stays the agentic kernel on `:8080`. Ultron is the company IDE on `:9090` with auth + RBAC. Operators and clients use the server — **no Cursor token required**.

## Boot

```bash
# Cashtro kernel
./scripts/always-on.sh          # :8080

# Ultron IDE
./scripts/ultron-on.sh          # :9090

# Both (tmux)
./scripts/fleet-on.sh
```

Default owner (first boot): `alejandro@proximityagency.ca` / `ultron-change-me`  
Override with `ULTRON_OWNER_PASSWORD`. Change after first login in production.

Disk: `data/ultron.json` (mode 0600).

## RBAC

| Role | Scope |
| --- | --- |
| owner | all perms |
| admin | companies, fleet, bridge, users read |
| operator | read companies, invoke fleet, OS bridge |
| client | own `companyIds` only |
| viewer | read own companies |

## API (Bearer token)

- `POST /api/auth/login` `{email,password}`
- `GET /api/companies` · `POST /api/companies`
- `GET /api/fleet/agents` · `POST /api/fleet/agents/{id}/invoke`
- `GET /api/os/health` · `POST /api/os/heartbeat` · `GET /api/os/proxy?path=/api/os`
- `GET|POST /api/users` · `GET /api/audit`

## Companies seeded

Giant · ScanApp · Proximity · Empire · Evolu-Jeunes · Cashtro

Agents with `cashtroId` bridge invokes into Cashtro OS. Others dry-run on Ultron until a worker binds.

## Map

Use skill `cashtro-agentic-map` / canvas `canvases/agentic-os-plan.canvas.tsx` when explaining hierarchy.
