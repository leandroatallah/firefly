# PROGRESS — 056-eight-dir-movement-skill

**Status:** Backlog

## Pipeline State

- [x] Story Architect
- [ ] Spec Engineer
- [ ] Mock Generator
- [ ] TDD Specialist
- [ ] Feature Implementer
- [ ] Workflow Gatekeeper

## Log

- [Sonnet 4.6] [Story Architect] 2026-05-09 [FINISHED]: USER_STORY.md created. Story introduces `EightDirectionalMovementSkill` in `internal/kit/skills/` as a genre-agnostic 8-direction input-to-body bridge, following the `HorizontalMovementSkill` pattern.
- [Sonnet 4.6] [Story Architect] 2026-05-09 [UPDATED]: Cross-cutting questions resolved in grilling session. Confirmed: Y axis is ground-plane depth (not altitude), skill is the sole input reader (passive model), skill is genre-agnostic and selected by `cfg.Movement.Mode == "eight_dir"`. Guard order (IsInputBlocked vs Immobile) remains open for spec engineer.
