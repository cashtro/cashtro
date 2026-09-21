#!/usr/bin/env bash
# Keep Ultron IDE up on the cloud VM — no Cursor token required for operators.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
mkdir -p data bin logs
ADDR="${ULTRON_ADDR:-:9090}"
DATA="${ULTRON_DATA:-$ROOT/data/ultron.json}"
CASHTRO="${CASHTRO_URL:-http://127.0.0.1:8080}"
LOG="${ULTRON_LOG:-$ROOT/logs/ultron.log}"
OWNER_PW="${ULTRON_OWNER_PASSWORD:-ultron-change-me}"

echo "[ultron-on] starting · addr=$ADDR · data=$DATA · cashtro=$CASHTRO"
while true; do
  echo "[ultron-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) build" | tee -a "$LOG"
  set +e
  go build -o "$ROOT/bin/ultron" ./cmd/ultron
  build=$?
  set -e
  if [ "$build" -ne 0 ]; then
    echo "[ultron-on] build failed · retry in 5s" | tee -a "$LOG"
    sleep 5
    continue
  fi
  if command -v fuser >/dev/null 2>&1; then
    fuser -k "${ADDR#:}/tcp" 2>/dev/null || true
  fi
  echo "[ultron-on] $(date -u +%Y-%m-%dT%H:%M:%SZ) boot" | tee -a "$LOG"
  set +e
  "$ROOT/bin/ultron" -addr "$ADDR" -data "$DATA" -cashtro "$CASHTRO" -boot-password "$OWNER_PW" >>"$LOG" 2>&1
  code=$?
  set -e
  echo "[ultron-on] exit=$code · restart in 2s" | tee -a "$LOG"
  sleep 2
done
