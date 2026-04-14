# Spec-Driven Development (SDD) Pipeline

A multi-agent workflow for implementing features with formal specification and TDD. The spec is the source of truth — code is its expression.

## Entry Point

Any feature request, task, or improvement idea. The pipeline starts from a human-readable description.

## Story Folder Convention

Each story moves through the pipeline as a self-contained folder:

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

The Log section is the agent's working memory across stateless sessions. Every agent **must** append an entry when it completes or when it rejects/backtracks.

## Pipeline Stages

### 1. Story Architect (`@story-architect`)
Translates the feature request into a User Story with Acceptance Criteria using DDD ubiquitous language.
- Creates `backlog/[ID]-[slug]/USER_STORY.md` and `PROGRESS.md`.

### 2. Spec Engineer (`@spec-engineer`)
Transforms the story into a Technical Specification: interface contracts, state machine transitions, pre/post-conditions.
- Moves folder from `backlog/` to `active/` using `bash scripts/story.sh start <id-slug>`.
- Writes `SPEC.md`, updates `PROGRESS.md`.

### 3. Mock Generator (`@mock-generator`)
Inspects `internal/engine/contracts/` and `internal/engine/mocks/`, generates or updates mocks required by the spec.
- Updates `PROGRESS.md` (or marks "skipped — no mocks required").

### 4. TDD Specialist (`@tdd-specialist`)
Writes failing `_test.go` files (Red phase) that exactly match the Spec's acceptance criteria.
- Updates `PROGRESS.md` with test file path(s).

### 5. Feature Implementer (`@feature-implementer`)
Writes the minimum production code to make the failing tests pass (Green phase). Does not modify tests.
- Updates `PROGRESS.md` with production file path(s).

### 6. Workflow Gatekeeper (`@workflow-gatekeeper`)
Validates spec compliance, TDD cycle, and code quality. Runs Coverage Analyzer to confirm a positive delta. Only if all checks pass, runs `golangci-lint run ./...` as the final gate.
- Updates `PROGRESS.md` to `✅ Done`.
- Moves folder from `active/` to `done/` using `bash scripts/story.sh done <id-slug>`.

## Pipeline Flow

```
[Feature Request]
      ↓
Story Architect → Spec Engineer → Mock Generator → TDD Specialist → Feature Implementer → Gatekeeper
                                                                                               ↓
                                                                                     Coverage Analyzer
                                                                                     (verification only)
```

## When to Use the Full Pipeline

- New features or behaviours being added to the engine or game.
- Changes to existing state machines, physics, or scene lifecycle.
- Any work where the expected behaviour needs to be formally specified first.

## When to Skip Ahead

- **Known bug with a clear fix**: start at TDD Specialist (write a failing test that reproduces the bug, then fix it).
- **Coverage gap on already-understood code**: start at Mock Generator or TDD Specialist directly.
- **Refactor with no behaviour change**: skip to Feature Implementer, Gatekeeper validates no regressions.

## Useful Scripts

```bash
bash scripts/story.sh new <id-slug>       # Create story in backlog/
bash scripts/story.sh start <id-slug>     # Move backlog/ → active/
bash scripts/story.sh done <id-slug>      # Move active/ → done/
bash scripts/story.sh status              # List all stories by lane
go run scripts/next-id.go                 # Print next available story ID
go run scripts/kanban.go                  # Generate kanban.html board
```

## Notes

- See `.agents/constitution.md` for non-negotiable project standards (code style, testing, mocks, coverage goals).
- See `AGENTS.md` for testing strategy, patterns, and coverage targets.
- No commits are made by agents. Version control is the developer's responsibility.
