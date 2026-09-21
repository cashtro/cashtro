---
name: cashtro-agentic-map
description: Visual map of Cashtro OS / agentic control plane for Castro. Use when explaining the agentic OS, fleet hierarchy, process table, architecture, overnight progress, or multi-agent structure — never text-only for those topics. Prefer Mermaid and/or a Cursor canvas.
---

# Cashtro agentic OS map

When explaining the agentic OS, fleet hierarchy, progress, architecture, or multi-agent structure for Castro: **never text-only**. Open or update a canvas and/or include Mermaid.

## Default Mermaid (update counts from `/api/about` + `/api/os`)

```mermaid
flowchart TB
  subgraph Ultron["Ultron IDE :9090"]
    Auth["RBAC · owner/admin/operator/client/viewer"]
    Cos["Giant · ScanApp · Proximity · Empire · …"]
    Fleet["agents + workers"]
  end
  subgraph Kernel["Cashtro OS :8080"]
    Init["init · os.about / os.heartbeat"]
    Live["14 agentics · delivery · research"]
  end
  DiskU["data/ultron.json"]
  DiskC["data/cashtro.json"]
  Ultron -->|HTTP bridge| Kernel
  Auth --> Cos --> Fleet
  Ultron --> DiskU
  Kernel --> DiskC
```

## What to show

- Kernel vs live vs resident (router unbound = resident)
- Delivery line: idea → concept → production
- Durability: disk image + always-on + overnight timer
- One sentence of progress (version + last capability), not a dashboard dump

## Canvas

Prefer `canvases/agentic-os-plan.canvas.tsx` when a Cursor canvas is available. Keep any `~/.cursor/admiral-fleet` STATUS/PLAN/ROSTER aligned if that tree exists in the environment.

## Do not

- Explain fleet/OS architecture as a long prose-only reply
- Overlay unrelated marketing/dashboard chrome on the map
