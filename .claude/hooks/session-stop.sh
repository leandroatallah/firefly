#!/usr/bin/env bash
set -euo pipefail

active_dir="/Users/leandroatallah/www/go/ebiten/firefly/.agents/work/active"

in_progress=""
for story_dir in "$active_dir"/*/; do
  [[ -d "$story_dir" ]] || continue
  progress="$story_dir/PROGRESS.md"
  [[ -f "$progress" ]] || continue
  if grep -q '\[/\]' "$progress"; then
    story=$(basename "$story_dir")
    in_progress+="  • $story"$'\n'
  fi
done

[[ -z "$in_progress" ]] && exit 0

msg="Stories with in-progress pipeline steps — log [FINISHED] in PROGRESS.md:"$'\n'"$in_progress"
jq -n --arg m "$msg" '{"systemMessage": $m}'
