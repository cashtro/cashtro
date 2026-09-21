---
name: cashtro-research-note
description: Ingest sourced research into Cashtro OS notes and memory (not chat). Use when gathering AOS research, seeding construction notes, Treg/Exa findings, arXiv/XKernel/12-factor insights, or the user says keep gathering info / research library / write a note.
---

# Cashtro research library

Findings live in the OS disk as **notes** and **facts**, not only in conversation.

## Ingest (preferred)

```bash
curl -s -X POST http://127.0.0.1:8080/api/notes \
  -H 'Content-Type: application/json' \
  -d '{
    "source":"short-id",
    "url":"https://...",
    "claim":"One sentence finding.",
    "quote":"Optional short quote ≤500 chars."
  }'
```

Or invoke `research` / `research.ingest` with the same JSON payload.

## Seed file (fresh boots)

Add durable seeds in `internal/agents/live.go` `seedResearch` so empty disks get the core library. Existing `data/cashtro.json` wins after `LoadFile` — still POST to the live API when updating a running image.

## Sources that already matter

- arXiv AOS papers (`arxiv:…`)
- XKernel, 12-factor agents
- Treg/Exa when signed in (catalog_search → catalog_get → call); if token expired, still land open-source notes
- Construction crumbs (`always-on`, overnight capabilities)

## Rules

- One claim per note; sourced `source` + `url` when possible
- Mirror important claims into memory via the research agent (ingest already remembers)
- Do not dump whole papers into quotes — clip
- OpenRouter is optional for research; deterministic ingest does not need a model

## After ingest

Confirm with `GET /api/notes` and `GET /api/memory?q=…`. Persist is automatic when the kernel has a data path.
