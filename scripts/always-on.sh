#!/usr/bin/env bash
# Keep Cashtro OS up. Closing the laptop does not stop the cloud VM.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p data bin logs
ADDR="${CASHTRO_ADDR:-:8080}"
DATA="${CASHTRO_DATA:-$ROOT/data/cashtro.json}"
LOG="${CASHTRO_LOG:-$ROOT/logs/cashtro.log}"
export CASHTRO_CLOSED="${CASHTRO_CLOSED:-1}"
export CASHTRO_PULSE_EVERY="${CASHTRO_PULSE_EVERY:-2m}"

echo "[always-on] building cashtro"
go build -o "$ROOT/bin/cashtro" ./cmd/cashtro

echo "[always-on] starting · addr=$ADDR · data=$DATA · closed=$CASHTRO_CLOSED · every=$CASHTRO_PULSE_EVERY"
echo "[always-on] desk can be closed. agentics keep building."
while true; do
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) boot" | tee -a "$LOG"
  set +e
  "$ROOT/bin/cashtro" -addr "$ADDR" -data "$DATA" >>"$LOG" 2>&1
  code=$?
  set -e
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) exit=$code · restart in 2s" | tee -a "$LOG"
  sleep 2
done
