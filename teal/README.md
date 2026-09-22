# Teal

Evolu-Jeunes voice brain. Fresh repo for Vapi talk, human-gated outbound
calls, intern desk, and every business-line bridge.

Cashtro OS (`cashtro/cashtro`) stays the manager. This product is
`github.com/Evolu-Jeunes/Teal`.

## Run

```bash
go test ./...
go run ./cmd/teal            # http://127.0.0.1:8090
```

Talk without a key. Change voice on the desk, `POST /api/vapi/voice`,
dashboard Assistants → Voice, or `VAPI_VOICE_ID`. Outbound
`POST /api/vapi/call` parks a confirm. Free Vapi numbers cannot outbound.
Do not buy numbers or fire live paid calls without Castro's allow.

```bash
export VAPI_ASSISTANT_ID=
export VAPI_API_KEY=
export VAPI_VOICE_ID=rachel
# optional Graph scrape
export TEAMS_TOKEN=
```

## Bridges

Proximity, Scan App, Panda, NFT/Giant, école, marketing, Empire,
trading, Pandora. `POST /api/vapi/bridge {"line":"proximity"}`.

## Split from cashtro

This directory is the Evolu-Jeunes/Teal module. Publish with a PAT that
can create org repos:

```bash
./scripts/publish-evolu-jeunes-teal.sh
```
