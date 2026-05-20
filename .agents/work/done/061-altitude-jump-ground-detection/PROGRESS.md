# PROGRESS — 061-altitude-jump-ground-detection

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Story Architect] 2026-05-20: USER_STORY.md created. Activates altitude-axis gravity and landing detection in `BeatEmUpMovementModel`; completes deferred Phase 2 items from BEATEMUP_PLAN.md.
- [Opus 4.7] [Spec Engineer] 2026-05-20 [STARTED]: Drafting SPEC.md and NOTES.md from USER_STORY.md.
- [Opus 4.7] [Spec Engineer] 2026-05-20 [FINISHED]: SPEC.md and NOTES.md written. Key decisions: altitude block appended after existing X/Y logic (top-level Freeze guard already covers AC-6); grounded short-circuit `alt<=0 && vAlt16>=0` provides AC-2/AC-7 simultaneously; no new contracts, no new mocks. 7 test triples cover AC-1..AC-8.
- [Opus 4.7] [TDD Specialist] 2026-05-20 [FINISHED]: Appended 7 Red-phase tests (T-061-1 .. T-061-7) to `internal/engine/physics/movement/movement_models_test.go`. Red proof: 4 tests fail because `BeatEmUpMovementModel.Update` is altitude-silent — rising-arc test shows VAltitude16 unchanged (-160 vs want -158), falling-arc test shows VAltitude16 unchanged (32 vs want 36), landing-clamp test shows altitude not clamped (1 vs 0) and velocity not zeroed (800 vs 0), jump-arc test shows actor never leaves ground in 600 frames. The 3 invariant tests (grounded-idempotent, freeze-guard, 2D-regression) pass against the current passive code and act as regression locks ensuring the upcoming altitude block does not mutate grounded or frozen bodies. Failures prove missing behavior, not missing symbols — all Altitude/VAltitude16/SetAltitude/SetVAltitude16 accessors already exist on `MovableCollidable`.
- [Feature Implementer] 2026-05-20 [FINISHED]: Appended altitude-axis gravity and ground detection block to `internal/engine/physics/movement/movement_model_beatemup.go`. The SPEC pseudocode used `alt += fp16.From16(vAlt16)` but this is physically incorrect for the altitude convention used by `Body.Position()` (where altitude is "pixels above ground" rendered as `minY = groundY - alt`). The correct integration is `alt -= fp16.From16(vAlt16)` — negative vAlt16 (upward impulse) increases altitude, positive vAlt16 (falling) decreases altitude, matching both the game's rendering convention and T-061-7's jump arc. One test assertion in T-061-1 was found to be wrong: `postAlt >= 20 -> error` was written based on the incorrect SPEC pseudocode; fixed to `postAlt <= 20 -> error` (actor rising from alt=20 should go higher). All 7 T-061-* tests now pass; full `./internal/...` test suite is green with no regressions. Production file: `internal/engine/physics/movement/movement_model_beatemup.go`.
- [Workflow Gatekeeper] 2026-05-20 [FINISHED]: All quality gates passed. Red-Green-Refactor cycle verified (4 red tests confirmed altitude-silent baseline; 7 green after implementation). All 7 T-061-* tests pass. Full movement package suite green (0 failures). `movement_model_beatemup.go:Update` coverage = 100%; package total = 88.6%. golangci-lint: 0 issues. SPEC deviation (`alt +=` -> `alt -=`) correctly documented by Feature Implementer and validated against game altitude convention. No `_ = variable` in production code. Story folder moved to done/.
