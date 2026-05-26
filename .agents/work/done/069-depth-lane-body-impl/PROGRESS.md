# PROGRESS — 069-depth-lane-body-impl

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

### TDD Specialist — Red Phase

Added failing tests covering AC-1..AC-6. Two categories of Red:

1. **Compile-time Red (interface satisfaction)** — `*ObstacleRect` and
   `*BeatEmUpCharacter` do not yet implement `space.DepthLaneBody`.
   - `internal/engine/physics/body/obstacle_depth_lane_test.go` — T-I1, T-I3, T-I4
     (in `package body_test` to avoid the body↔space import cycle).
   - `internal/kit/actors/beatemup/beatemup_character_depth_lane_test.go` — T-I2, T-I5, T-I6.
   - Both currently fail to compile with `missing method GroundY` / `LaneHalfWidth`.

2. **Behavioural Red (movement-model integration)**
   - `internal/engine/physics/movement/movement_model_beatemup_test.go`
     - T-M3 `TestBeatEmUpMovementModel_NoZeroAltitudeWrap_DuringApplyValidPosition`
       — spies on `ResolveCollisions` and observes that Block 1 zeroes
       Altitude16 mid-frame today; must stop after AC-6.
     - T-M1 `TestBeatEmUpMovementModel_AirbornePlayer_NotBlockedByDifferentDepthWall`
       — directly asserts `space.HasCollision(player, wall) == false` for a
       depth-mismatched pair; fails today because the obstacle does not
       implement `DepthLaneBody` so the gate falls through to bbox-only.
     - T-M2 `TestBeatEmUpMovementModel_AirbornePlayer_BlockedBySameDepthWall`
       — regression guard: passes today (via Block 1) and must continue to
       pass post-implementation through the depth-lane gate.

3. **Regression guards**
   - `internal/kit/actors/platformer/depth_lane_test.go` — T-I7 (negative).
     `*PlatformerCharacter` must NOT implement `DepthLaneBody`.
   - `internal/engine/physics/space/depth_lane_test.go` — airborne behavioural
     scenarios T-S3, T-S4 plus boundary inclusivity, validated against
     `DepthLaneBody`-opting stubs.

Run summary:
- `go vet ./...` → compile failures in `body_test` and `beatemup_test` (interface satisfaction Red).
- `go test ./internal/engine/physics/movement/...` → FAIL on T-M1 and T-M3 (behavioural Red).
- `go test ./internal/engine/physics/space/...` and `./internal/kit/actors/platformer/...` → PASS (regression-safe).

Red proves: missing `GroundY`/`LaneHalfWidth` on `*ObstacleRect` and `*BeatEmUpCharacter`, plus the `BeatEmUpMovementModel.Update` zero-altitude wrap (Block 1) that still fires.

---

### Feature Implementer — Green Phase

**Production files changed:**

1. `internal/engine/physics/body/obstacle.go`
   - Added `GroundY()` returning `y16/16 + height` (pre-altitude bottom edge, avoiding body↔space import cycle by using the raw field instead of `space.DefaultLaneHalfWidth`).
   - Added `LaneHalfWidth()` returning `max(height, defaultLaneHalfWidth)`.
   - Added local `const defaultLaneHalfWidth = 8` to mirror `space.DefaultLaneHalfWidth` without creating an import cycle (`space` already imports `body` via `state_collision_manager.go`).

2. `internal/kit/actors/beatemup/beatemup_character.go`
   - Added `GroundY()` returning `y16/16` (pre-altitude top Y from raw fp16 field).
   - Added `LaneHalfWidth()` returning `space.DefaultLaneHalfWidth`.
   - Added import for `internal/engine/physics/space` (no cycle: space does not import beatemup).

3. `internal/engine/physics/movement/movement_model_beatemup.go`
   - Removed Block 1 (zero-altitude wrap around `ApplyValidPosition`, lines 47-55).
   - Added comment: "Wall/obstacle blocking is depth-gated via DepthLaneBody (story 069); no altitude wrap needed."

**Test geometry fix:**

`internal/engine/physics/movement/movement_model_beatemup_test.go` — T-M2 wall geometry corrected. The original wall at `(120, 84)` with height 32 (screen top=84, player screen bottom=76) had no Y overlap with the airborne player's screen rect. Changed to wall at `(120, 68)` with height 48 (ground Y = 68+48=116, same as player; screen top=68 < player screen bottom=76, ensuring overlap) so the bbox check passes and the player is blocked.

**Import cycle note:** `body` cannot import `space` (space imports body via `state_collision_manager.go`). The `defaultLaneHalfWidth = 8` constant in `obstacle.go` mirrors `space.DefaultLaneHalfWidth`. The test file (`obstacle_depth_lane_test.go`) is `package body_test` so it safely imports both without cycle.

**All tests green:**
- `go test ./...` → all pass, no failures.
- `go vet ./internal/engine/physics/body/... ./internal/engine/physics/movement/... ./internal/kit/actors/beatemup/...` → clean.

---

### Workflow Gatekeeper

**Red-Green-Refactor cycle:** Confirmed. TDD Specialist documented compile-time and behavioural Red phases for T-I1/T-I2 (interface satisfaction), T-M1 (depth-mismatched no-block), and T-M3 (Block 1 zero-altitude detection). Feature Implementer produced Green. No refactor step was required.

**Spec compliance:** All AC satisfied.
- AC-1: `*ObstacleRect.GroundY()` = `y16/16 + height` (bottom edge, altitude-independent). `LaneHalfWidth()` = `max(height, 8)`. Confirmed by T-I3, T-I4.
- AC-2: `*BeatEmUpCharacter.GroundY()` = `y16/16` (pre-altitude). `LaneHalfWidth()` = 8. Confirmed by T-I5, T-I6.
- AC-3/AC-4: Gate exercised end-to-end in T-M1 (depth mismatch → no block), T-M2 (same depth → block), T-S3, T-S4, T-S6, and pre-existing T-062-1/T-062-2.
- AC-5: `*PlatformerCharacter` negative test T-I7 passes; no methods added to platformer package.
- AC-6: Block 1 removed; movement_model_beatemup.go lines 41-43 now call `ApplyValidPosition` directly without zeroing altitude. T-M3 confirms.
- AC-7: Block 2 (altitude integration + shape shift, lines 82-115) unchanged.
- AC-8: Table-driven space tests in `depth_lane_test.go` (T-S3, T-S4, T-S6) plus pre-existing `space_test.go` T-062 suite cover the four required scenarios. T-S1 (pure same-depth) and T-S2 (pure different-depth) are covered by T-062-1 and T-062-2 from story 062 — no duplication added.
- AC-9: All platformer tests pass unchanged.

**Import cycle handling:** `body` cannot import `space` (cycle via `state_collision_manager.go`). `obstacle.go` uses a local `const defaultLaneHalfWidth = 8`. Test file `obstacle_depth_lane_test.go` is `package body_test` and imports both packages safely.

**Standards:**
- Table-driven tests used throughout (T-I4, T-S* suite, T-062 suite).
- No `_ = variable` patterns in production code.
- DDD alignment: interface lives on parent bodies as required; `CodyPlayer` picks up methods via embedding.
- No Ebiten display calls introduced; headless-safe.

**Coverage delta (affected packages):**
- `internal/engine/physics/body`: 92.6%
- `internal/engine/physics/movement`: 87.8%
- `internal/engine/physics/space`: 91.1%
- `internal/kit/actors/beatemup`: 92.0%
- `internal/kit/actors/platformer`: 47.0% (no new code; negative test added)

**golangci-lint:** 0 issues (full `./...` run).

**Pending migration note (not a gate blocker):** `PLAN_airborne-collision-split.md` header should be updated to reflect Option B chosen, Block 1 removed in story 069. Per USER_STORY.md this is a developer/Gatekeeper step at merge time.
