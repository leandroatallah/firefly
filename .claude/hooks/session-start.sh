#!/usr/bin/env bash
set -euo pipefail

cd /Users/leandroatallah/www/go/ebiten/firefly
WORK=".agents/work"
status=""
for lane in backlog active; do
  status+="
── $lane ──
"
  entries=$(ls "$WORK/$lane" 2>/dev/null | sed 's/^/  /' || true)
  status+="${entries:-(empty)}
"
done
status+="
── done (latest) ──
"
latest=$(ls "$WORK/done" 2>/dev/null | tail -1)
status+="  ${latest:-(empty)}
"
jq -n --arg msg "$status" '{"systemMessage": $msg}'
