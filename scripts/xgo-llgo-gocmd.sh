#!/usr/bin/env bash
# XGO_GOCMD wrapper: run llgo but drop xgo's -ldflags -X ... (unsupported by LLGo).
# Usage: export XGO_GOCMD=/path/to/xgo-llgo-gocmd.sh
set -euo pipefail
args=()
skip=0
for a in "$@"; do
  if [[ $skip -eq 1 ]]; then
    skip=0
    continue
  fi
  if [[ "$a" == "-ldflags" ]]; then
    skip=1
    continue
  fi
  args+=("$a")
done
exec llgo "${args[@]}"
