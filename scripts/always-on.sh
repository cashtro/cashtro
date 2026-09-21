#!/usr/bin/env bash
# Keep Cashtro OS up. Runs on the cloud VM — closing your laptop does not stop it.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p data bin logs
ADDR="${CASHTRO_ADDR:-:8080}"
DATA="${CASHTRO_DATA:-$ROOT/data/cashtro.json}"
LOG="${CASHTRO_LOG:-$ROOT/logs/cashtro.log}"
export CASHTRO_CLOSED="${CASHTRO_CLOSED:-1}"

echo "[always-on] building cashtro"
go build -o "$ROOT/bin/cashtro" ./cmd/cashtro

# Free the port if a previous instance is stuck.
if command -v fuser >/dev/null 2>&1; then
  fuser -k "${ADDR#:}/tcp" 2>/dev/null || true
fi

echo "[always-on] starting · addr=$ADDR · data=$DATA · log=$LOG · closed=$CASHTRO_CLOSED"
echo "[always-on] laptop closed? fine. this is the cloud VM. work keeps flowing."
while true; do
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) boot" | tee -a "$LOG"
  set +e
  "$ROOT/bin/cashtro" -addr "$ADDR" -data "$DATA" >>"$LOG" 2>&1
  code=$?
  set -e
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) exit=$code · restart in 2s" | tee -a "$LOG"
  sleep 2
done
