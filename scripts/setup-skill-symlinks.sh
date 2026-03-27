#!/bin/bash
set -e

SKILLS_DIR=".agents/skills"
CLAUDE_DIR=".claude/skills"
KIRO_DIR=".kiro/skills"
QWEN_DIR=".qwen/skills"

# Create target directories
mkdir -p "$CLAUDE_DIR" "$KIRO_DIR" "$QWEN_DIR"

# Remove existing symlinks/directories
rm -rf "$CLAUDE_DIR"/* "$KIRO_DIR"/* "$QWEN_DIR"/* 2>/dev/null || true

# Create symlinks for each skill
for skill_dir in "$SKILLS_DIR"/*; do
    if [ ! -d "$skill_dir" ]; then
        continue
    fi
    
    skill_name=$(basename "$skill_dir")
    
    # Create symlinks (relative paths for portability)
    ln -s "../../$SKILLS_DIR/$skill_name" "$CLAUDE_DIR/$skill_name"
    ln -s "../../$SKILLS_DIR/$skill_name" "$KIRO_DIR/$skill_name"
    ln -s "../../$SKILLS_DIR/$skill_name" "$QWEN_DIR/$skill_name"
    
    echo "Linked: $skill_name"
done

echo "✓ Skill symlinks created"
