#!/bin/bash
set -e

AGENTS_DIR=".agents/agents"
CLAUDE_DIR=".claude/agents"
KIRO_DIR=".kiro/agents"
QWEN_DIR=".qwen/agents"
GEMINI_DIR=".gemini/agents"
HEADER=""

# Create target directories
mkdir -p "$CLAUDE_DIR" "$KIRO_DIR" "$QWEN_DIR" "$GEMINI_DIR"

# Process each agent file
for agent_file in "$AGENTS_DIR"/*.md; do
    if [ ! -f "$agent_file" ]; then
        continue
    fi
    
    filename=$(basename "$agent_file")

    # Extract name and description
    name=$(awk '/^name:/ {sub(/^name: /, ""); print; exit}' "$agent_file")
    description=$(awk '/^description:/ {sub(/^description: /, ""); print; exit}' "$agent_file")

    # Convert name to lowercase identifier for Claude (required: lowercase, hyphens only)
    claude_name=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')

    # Convert name to lowercase identifier for Qwen (e.g., "Coverage Analyzer" -> "coverage-analyzer")
    qwen_name="$claude_name"
    
    # Extract capabilities
    capabilities=$(awk '/^capabilities:/,/^---$/ {
        if ($0 ~ /^  - /) {
            sub(/^  - /, "")
            print
        }
    }' "$agent_file")
    
    # Get content after second ---
    content=$(awk 'BEGIN {count=0} /^---$/ {count++; next} count>=2 {print}' "$agent_file")
    
    # Map to Claude format
    claude_tools=""
    while IFS= read -r cap; do
        case "$cap" in
            "execute_commands"|"run_shell_command") claude_tools="${claude_tools:+$claude_tools, }Bash" ;;
            "write_files") claude_tools="${claude_tools:+$claude_tools, }Write" ;;
            "read_files") claude_tools="${claude_tools:+$claude_tools, }Read" ;;
            "code_intelligence") claude_tools="${claude_tools:+$claude_tools, }Grep, Glob" ;;
            "delegate_subagents") claude_tools="${claude_tools:+$claude_tools, }Task" ;;
        esac
    done <<< "$capabilities"
    
    cat > "$CLAUDE_DIR/$filename" <<EOF
---
name: $claude_name
description: $description
tools: $claude_tools
---

$content
EOF
    
    # Map to Kiro format (.json)
    kiro_tools="["
    first=true
    while IFS= read -r cap; do
        tool=""
        case "$cap" in
            "execute_commands"|"run_shell_command") tool="execute_bash" ;;
            "write_files") tool="fs_write" ;;
            "read_files") tool="fs_read" ;;
            "code_intelligence") tool="code" ;;
            "delegate_subagents") tool="use_subagent" ;;
        esac
        [ -z "$tool" ] && continue
        if [ "$first" = false ]; then
            kiro_tools="${kiro_tools}, "
        fi
        first=false
        kiro_tools="${kiro_tools}\"$tool\""
    done <<< "$capabilities"
    kiro_tools="${kiro_tools}]"

    kiro_filename="${filename%.md}.json"
    cat > "$KIRO_DIR/$kiro_filename" <<EOF
{
  "name": "$name",
  "description": "$description",
  "tools": $kiro_tools,
  "prompt": $(printf '%s' "$content" | python3 -c "import json,sys; print(json.dumps(sys.stdin.read().strip()))")
}
EOF
    
    # Map to Qwen format
    qwen_tools=""
    while IFS= read -r cap; do
        case "$cap" in
            "execute_commands"|"run_shell_command") qwen_tools="${qwen_tools}  - run_shell_command\n" ;;
            "write_files") qwen_tools="${qwen_tools}  - write_file\n" ;;
            "read_files") qwen_tools="${qwen_tools}  - read_file\n" ;;
            "code_intelligence") qwen_tools="${qwen_tools}  - grep_search\n  - glob\n" ;;
            "delegate_subagents") qwen_tools="${qwen_tools}  - task\n" ;;
        esac
    done <<< "$capabilities"

    cat > "$QWEN_DIR/$filename" <<EOF
---
name: $qwen_name
description: $description
tools:
$(printf "%b" "$qwen_tools" | sed '/^$/d')
---

$content
EOF
    
    # Map to Gemini format
    gemini_name=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' | tr -cd 'a-z0-9_-')
    gemini_tools=""
    while IFS= read -r cap; do
        case "$cap" in
            "execute_commands"|"run_shell_command") gemini_tools="${gemini_tools}  - run_shell_command\n" ;;
            "write_files") gemini_tools="${gemini_tools}  - write_file\n" ;;
            "read_files") gemini_tools="${gemini_tools}  - read_file\n" ;;
            "code_intelligence") gemini_tools="${gemini_tools}  - grep_search\n  - glob\n" ;;
        esac
    done <<< "$capabilities"

    cat > "$GEMINI_DIR/$filename" <<EOF
---
name: $gemini_name
description: $description
kind: local
tools:
$(printf "%b" "$gemini_tools" | sed '/^$/d')
---

$content
EOF

    echo "Synced: $filename"
done

echo "✓ Agent sync complete"
