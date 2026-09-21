---
name: cashtro-always-on
description: Keep Cashtro OS running on the cloud VM even when the laptop is closed. Use when health is down, always-on must start or restart, the user says never stop / laptop closed / always on, or when verifying tmux cashtro-always-on and :8080.
---

# Cashtro always-on (cloud VM)

Closing a laptop does **not** stop Cashtro OS while this cloud VM and the always-on loop are up.

## Check

```bash
curl -sf http://127.0.0.1:8080/health
tmux -f /exec-daemon/tmux.portal.conf ls
pgrep -af 'bin/cashtro|always-on.sh'
```

Healthy response includes `"status":"ok"`, `"always":true`, current `version`, and `data` pointing at `data/cashtro.json`.

## Start (if down)

```bash
SESSION_NAME="cashtro-always-on"
tmux -f /exec-daemon/tmux.portal.conf has-session -t "=$SESSION_NAME" 2>/dev/null \
  || tmux -f /exec-daemon/tmux.portal.conf new-session -d -s "$SESSION_NAME" -c /workspace -- "${SHELL:-bash}" -l
tmux -f /exec-daemon/tmux.portal.conf send-keys -t "$SESSION_NAME:0.0" './scripts/always-on.sh' C-m
```

Or: `make always-on` inside that tmux session.

`scripts/always-on.sh` rebuilds `bin/cashtro`, binds `:8080`, writes `logs/cashtro.log`, and restarts on exit.

## Restart onto new code

Kill **only** the OS binary so the loop rebuilds:

```bash
pkill -f '/workspace/bin/cashtro' || true
# wait for /health; confirm version bumped
```

Do **not** destroy the `cashtro-always-on` session unless it is wedged. If wedged, recreate the session and start the script again.

## Durability

- Disk image: `data/cashtro.json` (gitignored)
- Kernel autosave every ~60s while persist path is set
- Graceful stop calls `Kernel.Close()` for a final flush
- Overnight agent wake-ups: skill `cashtro-os-keepalive`

## Do not

- Require the user to keep their computer open
- Run the OS only in a one-shot foreground shell for overnight work
- Wipe `data/cashtro.json` unless the user explicitly asks to reset state
