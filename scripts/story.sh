#!/usr/bin/env bash
# Usage:
#   scripts/story.sh new  <id-slug>   — create story in backlog
#   scripts/story.sh start <id-slug>  — backlog → active  (Spec Engineer)
#   scripts/story.sh done  <id-slug>  — active  → done    (Gatekeeper)
#   scripts/story.sh status           — list all lanes
#   scripts/story.sh dashboard        — real-time pipeline dashboard

set -euo pipefail
WORK=".agents/work"

cmd="${1:-dashboard}"
slug="${2:-}"

case "$cmd" in
  new)
    dir="$WORK/backlog/$slug"
    mkdir -p "$dir"
    cat > "$dir/USER_STORY.md" <<EOF
# $slug

## As a...

## I want...

## So that...

## Acceptance Criteria
- [ ]
EOF
    cat > "$dir/PROGRESS.md" <<EOF
# PROGRESS — $slug

**Status:** 📋 Backlog

## Pipeline State
- [ ] Story Architect
- [ ] Spec Engineer
- [ ] Mock Generator
- [ ] TDD Specialist
- [ ] Feature Implementer
- [ ] Workflow Gatekeeper

## Log
# Format: - [Model] [Agent] [date]: Action/Decision
EOF
    echo "✅ Created backlog/$slug"
    ;;

  start)
    mv "$WORK/backlog/$slug" "$WORK/active/$slug"
    sed -i '' 's/📋 Backlog/🔄 Active/' "$WORK/active/$slug/PROGRESS.md"
    echo "🔄 Moved to active: $slug"
    ;;

  done)
    mv "$WORK/active/$slug" "$WORK/done/$slug"
    sed -i '' 's/🔄 Active/✅ Done/' "$WORK/done/$slug/PROGRESS.md"
    echo "✅ Moved to done: $slug"
    ;;

  status)
    for lane in backlog active done; do
      echo "── $lane ──"
      ls "$WORK/$lane" 2>/dev/null | sed 's/^/  /' || echo "  (empty)"
    done
    ;;

  dashboard)
    go run scripts/sdd-dashboard.go
    ;;

  *)
    echo "Unknown command: $cmd"; exit 1 ;;
esac

