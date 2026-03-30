---
name: agent-workflow
description: A Spec-Driven Development (SDD) pipeline for implementing features and improvements with TDD.
---

# Chained Agent Workflow

A multi-agent pipeline following Spec-Driven Development (SDD). The spec is the source of truth — code is its expression.

## Entry Point

Any feature request, task, or improvement idea. The pipeline starts from a human-readable description, not from a coverage report.

## Agents

**1. Story Architect**
Translates the feature request into a User Story with Acceptance Criteria using DDD ubiquitous language. Writes to `.agents/work/backlog/USER_STORY_[ID].md`.

**2. Spec Engineer**
Transforms the story into a Technical Specification: interface contracts, state machine transitions, pre/post-conditions. Writes to `.agents/work/active/SPEC_[ID].md`.

**3. Mock Generator**
Inspects `internal/engine/contracts/` and `internal/engine/mocks/`, generates or updates mocks required by the spec. Decides shared vs. package-local placement.

**4. TDD Specialist**
Writes failing `_test.go` files (Red phase) that exactly match the Spec's acceptance criteria. Tests verify observable behavior through public interfaces.

**5. Feature Implementer**
Writes the minimum production code to make the failing tests pass (Green phase). Does not modify tests.

**6. Workflow Gatekeeper**
Validates spec compliance, TDD cycle, and code quality. Runs Coverage Analyzer to confirm a positive delta. Moves the story to `.agents/work/done/` on success, or triggers a backtrack on failure.

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
