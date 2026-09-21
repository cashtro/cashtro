#!/usr/bin/env bash
# Run Cashtro OS + Ultron IDE together (company control plane over agentic kernel).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMUX_CONF="/exec-daemon/tmux.portal.conf"
tmux_cmd() {
  if [ -f "$TMUX_CONF" ]; then
    tmux -f "$TMUX_CONF" "$@"
  else
    tmux "$@"
  fi
}

ensure() {
  local name="$1"
  local script="$2"
  if tmux_cmd has-session -t "=$name" 2>/dev/null; then
    echo "[fleet-on] $name already running"
    return
  fi
  tmux_cmd new-session -d -s "$name" -c "$ROOT" -- "${SHELL:-bash}" -l
  tmux_cmd send-keys -t "$name:0.0" "$script" C-m
  echo "[fleet-on] started $name"
}

ensure cashtro-always-on "./scripts/always-on.sh"
sleep 1
ensure ultron-always-on "./scripts/ultron-on.sh"
echo "[fleet-on] Cashtro :8080 · Ultron IDE :9090 · laptop may be closed"
