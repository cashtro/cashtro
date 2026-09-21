---
name: cashtro-agentic-map
description: Visual map of Cashtro OS / agentic control plane for Castro. Use when explaining the agentic OS, fleet hierarchy, process table, architecture, overnight progress, or multi-agent structure — never text-only for those topics. Prefer Mermaid and/or a Cursor canvas.
---

# Cashtro agentic OS map

When explaining the agentic OS, fleet hierarchy, progress, architecture, or multi-agent structure for Castro: **never text-only**. Open or update a canvas and/or include Mermaid.

## Default Mermaid (update counts from `/api/os`)

```mermaid
flowchart TB
  subgraph Kernel["Cashtro OS kernel"]
    Init["init · os.about / os.heartbeat"]
    Router["router · optional OpenRouter"]
  end
  subgraph Live["Live agentics"]
    Delivery["delivery"]
    Research["research"]
    Explorer["explorer"]
    Operator["operator"]
    Reviewer["reviewer"]
    Architect["architect"]
    Deploy["deploy"]
    Security["security"]
    Memory["memory"]
    Comms["comms"]
    Planner["planner"]
    Investigator["investigator"]
  end
  Disk["data/cashtro.json · autosave"]
  Desk["HTTP desk :8080"]
  Always["scripts/always-on.sh · cloud VM"]
  Kernel --> Live
  Live --> Disk
  Desk --> Kernel
  Always --> Desk
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
