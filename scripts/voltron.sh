#!/usr/bin/env bash
# Voltron is Cashtro OS assembled: kernel + every registered agentic.
# This supervisor keeps the desk listening and restarts it if it dies.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

STATE_DIR="${VOLTRON_STATE_DIR:-/tmp/voltron}"
LOG="${VOLTRON_LOG:-$STATE_DIR/voltron.log}"
STATUS="${VOLTRON_STATUS:-$STATE_DIR/status.json}"
PIDFILE="${VOLTRON_PIDFILE:-$STATE_DIR/voltron.pid}"
ADDR="${VOLTRON_ADDR:-:8080}"
BIN="${VOLTRON_BIN:-$ROOT/bin/cashtro}"
HEALTH_URL="${VOLTRON_HEALTH_URL:-http://127.0.0.1:8080/health}"

mkdir -p "$STATE_DIR" "$ROOT/bin"

log() {
  printf '[%s] %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$*" | tee -a "$LOG"
}

write_status() {
  local state="$1"
  local pid="${2:-0}"
  local note="${3:-}"
  python3 - "$STATUS" "$state" "$pid" "$note" <<'PY'
import json, sys, datetime
path, state, pid, note = sys.argv[1], sys.argv[2], int(sys.argv[3]), sys.argv[4]
payload = {
    "name": "voltron",
    "os": "Cashtro OS",
    "state": state,
    "pid": pid,
    "note": note,
    "url": "http://127.0.0.1:8080",
    "health": "http://127.0.0.1:8080/health",
    "updatedAt": datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
}
with open(path, "w", encoding="utf-8") as fh:
    json.dump(payload, fh, indent=2)
    fh.write("\n")
PY
}

build() {
  log "building cashtro → $BIN"
  go build -o "$BIN" ./cmd/cashtro
}

health_ok() {
  curl -fsS --max-time 2 "$HEALTH_URL" >/dev/null 2>&1
}

supervisor_running() {
  [[ -f "$PIDFILE" ]] || return 1
  local pid
  pid="$(cat "$PIDFILE" 2>/dev/null || true)"
  [[ -n "${pid:-}" ]] && kill -0 "$pid" 2>/dev/null
}

cmd_status() {
  local pid=""
  if supervisor_running; then
    pid="$(cat "$PIDFILE")"
  fi
  if health_ok; then
    echo "voltron: up pid=${pid:-unknown} $HEALTH_URL"
    curl -fsS --max-time 2 "$HEALTH_URL"
    echo
    return 0
  fi
  echo "voltron: down"
  [[ -f "$STATUS" ]] && cat "$STATUS"
  return 1
}

run_loop() {
  echo $$ >"$PIDFILE"
  trap 'write_status stopped 0 "supervisor exiting"; rm -f "$PIDFILE"; exit 0' INT TERM
  build
  local crashes=0
  while true; do
    write_status starting 0 "launching cashtro on $ADDR"
    log "starting cashtro on $ADDR"
    "$BIN" -addr "$ADDR" >>"$LOG" 2>&1 &
    local child=$!
    write_status running "$child" "cashtro listening on $ADDR"
    # Wait until the child dies. A healthy process stays here overnight.
    if wait "$child"; then
      log "cashtro exited 0"
    else
      log "cashtro exited $?"
    fi
    crashes=$((crashes + 1))
    write_status restarting 0 "exited; restart ${crashes} in 2s"
    log "restart ${crashes} in 2s"
    sleep 2
    build || log "rebuild failed; retrying previous binary"
  done
}

cmd_start() {
  if health_ok && supervisor_running; then
    log "already up"
    cmd_status
    return 0
  fi
  if supervisor_running; then
    log "supervisor alive but health failed; leaving it to restart"
    cmd_status || true
    return 0
  fi
  run_loop
}

case "${1:-start}" in
  start) cmd_start ;;
  status) cmd_status ;;
  *)
    echo "usage: $0 {start|status}" >&2
    exit 2
    ;;
esac
