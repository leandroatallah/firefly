---
name: unstaged-code-review
description: Review unstaged changes, explain what changed, suggest improvements, and propose focused commits.
---

# Code Review & Commit Suggestions

Review unstaged changes, explain what changed, suggest improvements, and propose focused commits.

## Workflow

**1. Inspect unstaged changes**

```bash
git diff
git status
```

**2. Explain changes**
For each modified file, summarize:

- What was added, removed, or refactored
- Why the change likely exists (inferred from context)

**3. Suggest improvements**
Look for:

- Missing error handling
- Untested code paths
- Style violations (e.g., `_ = variable` pattern)
- Logic that could be simplified

**4. Suggest commits**
Group related changes into small, focused commits. Prefer [Conventional Commits](https://www.conventionalcommits.org/) format:

```
feat(entity): add handleState transition for idle
fix(physics): correct fp16 overflow on fast movement
test(actors): add table-driven tests for state machine
refactor(i18n): extract T() fallback logic
```

## Rules

- One concern per commit — avoid mixing features, fixes, and tests
- Keep commit messages under 72 characters
- Never commit on behalf of the user — only suggest
