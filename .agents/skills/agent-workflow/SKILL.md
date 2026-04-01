---
name: agent-workflow
description: A Spec-Driven Development (SDD) pipeline for implementing features and improvements with TDD.
---

# Chained Agent Workflow

A multi-agent pipeline following Spec-Driven Development (SDD). The spec is the source of truth — code is its expression.

## Entry Point

Any feature request, task, or improvement idea. The pipeline starts from a human-readable description, not from a coverage report.

## Story Folder Convention

Each story lives in a self-contained folder that moves through the pipeline:

```
.agents/work/
├── backlog/[ID]-[slug]/
│   ├── USER_STORY.md   ← written by Story Architect
│   └── PROGRESS.md     ← pipeline status tracker
├── active/[ID]-[slug]/
│   ├── USER_STORY.md
│   ├── SPEC.md         ← written by Spec Engineer
│   └── PROGRESS.md
└── done/[ID]-[slug]/
    ├── USER_STORY.md
    ├── SPEC.md
    └── PROGRESS.md
```

`PROGRESS.md` is the single source of truth for pipeline state. Every agent updates it before finishing.

## PROGRESS.md Format

```markdown
# PROGRESS — [ID]-[slug]

**Status:** 🔄 Active   (or ✅ Done)

## Pipeline Stages

| Stage | Status | Notes |
|---|---|---|
| Story Architect    | ✅ Complete | `USER_STORY.md` written |
| Spec Engineer      | ✅ Complete | `SPEC.md` written |
| Mock Generator     | ✅ Complete | No mocks required |
| TDD Specialist     | ✅ Complete | `path/to/foo_test.go` |
| Feature Implementer| ✅ Complete | `path/to/foo.go` |
| Gatekeeper         | ⬜ Pending  | |

## Log

- **[Agent] [date]**: What was done, decided, or why a backtrack occurred.
  Example: `Gatekeeper rejected: coverage dropped 2% — missing edge case in TestFooExit. Backtracking to TDD Specialist.`
```

The Log section is the agent's working memory across stateless sessions. Every agent **must** append an entry when it completes or when it rejects/backtracks. Keep entries concise — one or two lines max. This is the primary context source for the next agent starting a fresh session.

## Agents

**1. Story Architect**
Translates the feature request into a User Story with Acceptance Criteria using DDD ubiquitous language.
- Creates `backlog/[ID]-[slug]/USER_STORY.md` and `PROGRESS.md`.

**2. Spec Engineer**
Transforms the story into a Technical Specification: interface contracts, state machine transitions, pre/post-conditions.
- Moves folder from `backlog/` to `active/` using `bash scripts/story.sh start <id-slug>`.
- Writes `SPEC.md`, updates `PROGRESS.md`.

**3. Mock Generator**
Inspects `internal/engine/contracts/` and `internal/engine/mocks/`, generates or updates mocks required by the spec.
- Updates `PROGRESS.md` (or marks "skipped — no mocks required").

**4. TDD Specialist**
Writes failing `_test.go` files (Red phase) that exactly match the Spec's acceptance criteria.
- Updates `PROGRESS.md` with test file path(s).

**5. Feature Implementer**
Writes the minimum production code to make the failing tests pass (Green phase). Does not modify tests.
- Updates `PROGRESS.md` with production file path(s).

**6. Workflow Gatekeeper**
Validates spec compliance, TDD cycle, and code quality. Runs Coverage Analyzer to confirm a positive delta.
- Updates `PROGRESS.md` to `✅ Done`.
- Moves folder from `active/` to `done/` using `bash scripts/story.sh done <id-slug>`.

## Chain

```
[Feature Request]
      ↓
Story Architect → Spec Engineer → Mock Generator → TDD Specialist → Feature Implementer → Gatekeeper
                                                                                               ↓
                                                                                     Coverage Analyzer
                                                                                     (verification only)
```

## When to use the full pipeline

- New features or behaviours being added to the engine or game.
- Changes to existing state machines, physics, or scene lifecycle.
- Any work where the expected behaviour needs to be formally specified first.

## When you can skip to a specific agent

- **Known bug with a clear fix**: start at TDD Specialist (write a failing test that reproduces the bug, then fix it).
- **Coverage gap on already-understood code**: start at Mock Generator or TDD Specialist directly.
- **Refactor with no behaviour change**: skip to Feature Implementer, Gatekeeper validates no regressions.
