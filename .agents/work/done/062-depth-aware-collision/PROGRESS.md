# PROGRESS — 062-depth-aware-collision

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [-] Mock Generator (skipped — no new contracts)
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Story Architect] 2026-05-20: USER_STORY.md created. Depth-lane gate added to `HasCollision`; BEATEMUP_PLAN Phase 3 collision item. Opt-in mechanism and LaneWidth source deferred to Spec Engineer.
- [Claude Opus 4.7] [Spec Engineer] 2026-05-20 [STARTED]: Drafting SPEC.md for depth-aware collision gating in `internal/engine/physics/space/space.go`.
- [Claude Opus 4.7] [Spec Engineer] 2026-05-20 [FINISHED]: SPEC.md and NOTES.md written. Key decisions: (1) Opt-in via local interface `DepthLaneBody` in `space/depth_lane.go` (no contract change, no kit/game import). (2) LaneWidth via per-body `LaneHalfWidth() int` method; pair tolerance = max of the two; `DefaultLaneHalfWidth=8` exported for downstream docs only. (3) GroundY via `DepthLaneBody.GroundY() int` (not `GetPosition16()`), because screen-Y diverges from ground-Y once altitude is active (story 061). No new contracts → Mock Generator can be skipped. 9 table-driven test cases enumerated covering AC-1..AC-7 + edge cases.
- [Claude Opus 4.7] [TDD Specialist] 2026-05-20 [FINISHED]: Added `depthLaneCollidable` test fixture and `TestHasCollisionDepthLane` (10 sub-cases across 9 IDs T-062-1..T-062-9 with -8A/-8B split) to `internal/engine/physics/space/space_test.go`. Created minimal `internal/engine/physics/space/depth_lane.go` stub declaring `DefaultLaneHalfWidth = 8` and the `DepthLaneBody` interface (`GroundY()`, `LaneHalfWidth()`), but left `HasCollision` unmodified. Red proof: `go test ./internal/engine/physics/space/... -run TestHasCollisionDepthLane` FAILS on T-062-2 (different-lane returns true), T-062-6 (halfW=0 off-by-one returns true), T-062-8B (asymmetric halfW exceeded returns true) — all three failures are the cases where the depth-lane gate must *reject* a bbox-overlap match, proving the missing behavior (not a missing symbol).
- [Feature Implementer] 2026-05-20 [FINISHED]: Modified `HasCollision` in `internal/engine/physics/space/space.go` to add depth-lane gate. Bbox overlap is now captured before returning; when both bodies implement `DepthLaneBody`, tolerance = max of the two `LaneHalfWidth()` values and `abs(da.GroundY()-db.GroundY()) <= tol` must hold. `depth_lane.go` stub required no changes. All 10 `TestHasCollisionDepthLane` sub-cases pass (including T-062-2, T-062-6, T-062-8B). Full `go test ./internal/engine/physics/space/...` passes. `go build ./internal/engine/...` succeeds.
- [Workflow Gatekeeper] 2026-05-20 [REJECTED]: golangci-lint reports one staticcheck violation (QF1001) in `internal/engine/physics/space/space.go` line 256: `if !(okA && okB)` must be rewritten as `if !okA || !okB` (De Morgan's law). All tests pass (10/10 sub-cases), coverage delta is positive (+1.0 pp: 90.1% → 91.1%), spec/story alignment is complete, and layer rules are upheld. The single blocker is the lint error. Backtrack to: Feature Implementer. Required fix: replace `if !(okA && okB) {` with `if !okA || !okB {` in `HasCollision`.
- [Feature Implementer] 2026-05-20 [LINT FIX]: Confirmed `if !okA || !okB {` already present at line 256 of `internal/engine/physics/space/space.go` (fix was applied prior to this session). `golangci-lint run ./...` exits clean (0 issues). `go test ./internal/engine/physics/space/...` passes. No source changes required.
- [Workflow Gatekeeper] 2026-05-20 [APPROVED]: All quality gates passed on re-validation. Tests: 10/10 sub-cases pass (T-062-1 through T-062-9 including -8A/-8B), full package suite green (41 passing, 1 skipped). Coverage: 91.1% (+1.0 pp from pre-story 90.1%). Lint: 0 issues. Layer rules: no forbidden imports in `internal/engine/physics/space/`. Story/Spec alignment: all 9 ACs satisfied. Red-Green-Refactor cycle confirmed in log. Promoting to done/.
