#!/usr/bin/env bash
# Keep Cashtro OS up. Runs on the cloud VM — closing your laptop does not stop it.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p data bin logs
ADDR="${CASHTRO_ADDR:-:8080}"
DATA="${CASHTRO_DATA:-$ROOT/data/cashtro.json}"
LOG="${CASHTRO_LOG:-$ROOT/logs/cashtro.log}"

echo "[always-on] starting · addr=$ADDR · data=$DATA · log=$LOG"
echo "[always-on] laptop closed? fine. this is the cloud VM."
while true; do
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) build" | tee -a "$LOG"
  set +e
  go build -o "$ROOT/bin/cashtro" ./cmd/cashtro
  build=$?
  set -e
  if [ "$build" -ne 0 ]; then
    echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) build failed · retry in 5s" | tee -a "$LOG"
    sleep 5
    continue
  fi
  if command -v fuser >/dev/null 2>&1; then
    fuser -k "${ADDR#:}/tcp" 2>/dev/null || true
  fi
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) boot" | tee -a "$LOG"
  set +e
  "$ROOT/bin/cashtro" -addr "$ADDR" -data "$DATA" >>"$LOG" 2>&1
  code=$?
  set -e
  echo "[always-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) exit=$code · restart in 2s" | tee -a "$LOG"
  sleep 2
done
