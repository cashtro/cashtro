# n8n fabric — 14 kernel + Steel, 58 specialists

n8n is **not** a second product. It is the workflow fabric on the
14 Voltron seats plus Steel on the control plane. Width 58. Depth still 3.
Each kernel seat has a self-improve specialist. Steel contradicts.

```mermaid
flowchart TB
  subgraph depth1["Depth 1 — control plane"]
    CP[orchestrate / GET /fleet / GET /corp]
  end
  subgraph depth2["Depth 2 — 58 specialists"]
    S1[14 seats · verbs + self-improve]
    S2[Steel · contradict / test / optimize]
  end
  subgraph depth3["Depth 3 — Voltron kernel"]
    K[POST /api/agents/:seat/invoke]
  end
  CP -->|resolveSpecialist| S1
  CP --> S2
  S1 -->|N8N_BASE_URL empty| K
  S2 -->|local $0| ST[contrarian adapter]
```

## Unbound tonight

`N8N_BASE_URL` empty → adapter posts nowhere and falls back to
`http://127.0.0.1:8080/api/agents/:seat/invoke`. Steel stays local and
never hits the kernel. The 58 workers exist in the registry (`GET /fleet`).

## Bind later (local only)

```bash
docker compose -f docker-compose.n8n.yml up -d
# .env
N8N_BASE_URL=http://127.0.0.1:5678
```

Do not buy n8n Cloud. Do not point a workflow at a live client site.

## Map

See `state/n8n-fleet.json`. Fifty-eight rows, each `{ seat, name, capability, webhook }`.
Steel is `n8n-55`–`n8n-58`. Self-improve is `n8n-41`–`n8n-54`.
