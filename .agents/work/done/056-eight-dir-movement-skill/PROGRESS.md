# PROGRESS — 056-eight-dir-movement-skill

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Sonnet 4.6] [Story Architect] 2026-05-09 [FINISHED]: USER_STORY.md created. Story introduces `EightDirectionalMovementSkill` in `internal/kit/skills/` as a genre-agnostic 8-direction input-to-body bridge, following the `HorizontalMovementSkill` pattern.
- [Sonnet 4.6] [Story Architect] 2026-05-09 [UPDATED]: Cross-cutting questions resolved in grilling session. Confirmed: Y axis is ground-plane depth (not altitude), skill is the sole input reader (passive model), skill is genre-agnostic and selected by `cfg.Movement.Mode == "eight_dir"`. Guard order (IsInputBlocked vs Immobile) remains open for spec engineer.
- [Opus 4.7] [Spec Engineer] 2026-05-13 [STARTED]: Drafting SPEC.md and NOTES.md for EightDirectionalMovementSkill.
- [Opus 4.7] [Spec Engineer] 2026-05-13 [FINISHED]: SPEC.md created. Key decisions: (1) Guard order is `IsInputBlocked` first, then `Immobile` (matches `HorizontalMovementSkill`); (2) `Immobile` zeroes both vx16 and vy16 (and both accel components) since this is an 8-dir skill, leaving altitude (`VAltitude16`) untouched; (3) no `axis` smoothing — skill calls `OnMoveLeft/Right/Up/Down` directly per pressed direction (no inertia field in MovementConfig for Y); (4) does not introduce new contracts — no mocks required.
- [Opus 4.7] [TDD Specialist] 2026-05-13 [FINISHED]: Wrote failing tests at `internal/kit/skills/eight_dir_move_test.go` (8-case table-driven `HandleInput` + `New`/`Update_NoOp`/`ActivationKey`) and appended `EightDirectionalMovementSkill_satisfies_Skill` sub-test to `internal/kit/skills/package_surface_test.go`. Added local `mockEightDirBody` wrapping `mockMovableCollidable` to record `OnMoveLeft/Right/Up/Down` call counts and last args. Red Phase verified: `go test ./internal/kit/skills/` fails to build with `undefined: NewEightDirectionalMovementSkill` (×4) — proving the missing behavior (skill type + constructor + HandleInput dispatch to 4-direction OnMove* calls under blocked/immobile guards) rather than just a missing symbol elsewhere.
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-05-13 [FINISHED]: All quality gates passed. Coverage delta: `internal/kit/skills` at 93.6% overall; `eight_dir_move.go` all functions at 100%. Red-Green-Refactor cycle confirmed (TDD log shows red phase). All 12 test scenarios (T-1..T-12) implemented and passing, including 8-case table-driven HandleInput. Implementation matches SPEC.md exactly: guard order IsInputBlocked->Immobile, both velocity and acceleration zeroed on Immobile, direct OnMove* dispatch with no inertia, no genre imports. golangci-lint: 0 issues.
