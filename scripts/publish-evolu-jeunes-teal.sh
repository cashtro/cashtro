#!/usr/bin/env bash
# Publish teal/ as github.com/Evolu-Jeunes/Teal.
# Needs a PAT that can create Evolu-Jeunes repos (this Cloud Agent token cannot).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
ORG="Evolu-Jeunes"
NAME="Teal"
REMOTE="https://github.com/${ORG}/${NAME}.git"

if ! gh api "repos/${ORG}/${NAME}" >/dev/null 2>&1; then
  echo "creating ${ORG}/${NAME} …"
  gh repo create "${ORG}/${NAME}" --private --description "Evolu-Jeunes Teal voice brain — Vapi talk, intern desk, every bridge" || {
    echo "cannot create ${ORG}/${NAME}. Need Option A PAT (docs/ACCESS_REQUIRED.md)."
    exit 1
  }
fi

git subtree split --prefix=teal -b evolu-jeunes-teal
git push -u "$REMOTE" evolu-jeunes-teal:main
echo "published ${REMOTE}"
