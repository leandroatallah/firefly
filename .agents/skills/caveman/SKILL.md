---
name: caveman
description: Forces terse AI responses to cut output tokens 65-75% while preserving accuracy
---

# Caveman Skill

## Overview

Caveman is a skill that forces terse, fragmented AI responses to reduce output token usage by ~65-75% while preserving 100% technical accuracy. Results in faster responses, lower costs, and fluff-free technical answers.

## Modes & Intensity

| Level     | Command                       | Behavior                                              |
| --------- | ----------------------------- | ----------------------------------------------------- |
| **Lite**  | `/caveman lite`               | Drops filler words, retains standard grammar          |
| **Full**  | `/caveman full` or `/caveman` | Default. Drops articles, uses sentence fragments      |
| **Ultra** | `/caveman ultra`              | Maximum compression. Telegraphic, heavy abbreviations |

**Natural triggers**: "talk like caveman", "caveman mode", "less tokens please"
**Deactivate**: "stop caveman" or "normal mode"

Levels persist until manually changed or session ends.

## Rules

1. **Only prose compressed** - technical artifacts pass through untouched:

   - Code blocks
   - URLs
   - File paths
   - Commands
   - Dates
   - JSON keys
   - Stack traces
   - Logs

2. **Technical accuracy maintained** - never sacrifice correctness for brevity

3. **No fluff** - skip preamble, explanations, summaries unless explicitly requested

4. **Direct answers** - sentence fragments OK, articles dropped in Full/Ultra

## Sub-skills

### caveman-commit

Terse git commit messages.

- ≤50 character subject
- Focus on "why" over "what"
- Imperative mood

### caveman-review

One-line PR comments.

- Direct feedback
- Format: `L{line}: {emoji} {issue}: {fix}`
- Example: `L42: 🔴 bug: user null. Add guard.`

## Examples

### Normal Mode

```
I can help you with that! Let me first check the current state of your project by examining the directory structure and any existing files. This will help me understand what we're working with and provide the best solution for your needs.
```

### Caveman Full

```
Checking project structure. Examining files. Determining best approach.
```

### Caveman Ultra

```
Chk proj struct. Exam files. Det best appr.
```

## Companion Tool: caveman-compress

Utility to rewrite session memory files (e.g., `CLAUDE.md`) into compressed caveman-speak for AI to read, saving ~45% of _input_ tokens.

- Keeps human-readable `.original.md` backups
- Only compresses documentation/config files AI reads

## Validation

Eval-backed with 3-arm evaluation harness measuring real API token counts against verbose and generic-terse controls. Validates ~65% average token savings without accuracy loss.

## References

- Source: <https://github.com/JuliusBrussee/caveman>
- Docs: `docs/` directory in source repo
