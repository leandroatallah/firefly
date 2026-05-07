# PROGRESS — 053-altitude-engine-foundation

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [x] Mock Generator (handled inline by Feature Implementer — mocks_test.go files updated)
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Sonnet] [Story Architect] 2026-05-06 [FINISHED]: USER_STORY.md created for Phase 1 altitude engine foundation — Body/Movable contract extensions, physics body implementation, screen-coord mapping (ScreenY = Y - Altitude), Z-sort by Y, and beat-em-up actor scaffold.
- [Opus] [Spec Engineer] 2026-05-06 [FINISHED]: SPEC.md created. Key decisions: (1) Z-sort lives in a new engine-layer package `internal/engine/render/draworder/` (not on `BodiesSpace`) so collision/physics keep deterministic ID order while only the render path sorts by Y; (2) sort key is fp16 ground Y (`y16`) — altitude deliberately ignored so jumping entities stay in the same depth slot; (3) `sort.SliceStable` on a copy of the input slice (Bodies() cache must not be mutated); (4) `BeatEmUpCharacter` scaffold embeds the existing `kitactors.MeleeCharacter` to compose cleanly with future Phase 2 work; (5) `Position()` change is bit-identical when `altitude16 == 0`, preserving every existing 2D test.
- [Opus] [TDD Specialist] 2026-05-06 [FINISHED]: Red Phase landed across four files. Test files: `internal/engine/physics/body/body_test.go` (extended; adds compile-time `bodycontract.Body` assertion plus `TestBody_AltitudeAccessors`, `TestBody_Altitude16_StoredDirectly`, `TestBody_SetAltitude_UsesFp16`, `TestBody_Position_AltitudeMapsToScreenY`, `TestBody_Position_ZeroAltitude_IsBitIdentical`); `internal/engine/physics/body/body_movable_test.go` (extended; adds `bodycontract.Movable` assertion plus `TestMovableBody_AltitudeDynamics`, `TestMovableBody_AltitudeDefaultZero`); `internal/engine/render/draworder/draworder_test.go` (new; package-local `fakeCollidable` covers ascending order, stable equal-Y, altitude-ignored, no-mutation, empty/single); `internal/kit/actors/beatemup/beatemup_character_test.go` (new; `TestNewBeatEmUpCharacter_NotNil`). Red proof: `go test ./internal/engine/physics/body/ ./internal/engine/render/draworder/ ./internal/kit/actors/beatemup/` — physics/body fails to build because `SetAltitude/Altitude/Altitude16/SetAltitude16/VAltitude16/SetVAltitude16/AccelerationAltitude/SetAccelerationAltitude` are not defined on `*Body`/`*MovableBody`; draworder and beatemup fail with "no non-test Go files" because the production packages do not yet exist. Failures are missing-behavior, not missing-symbols-only — the position-mapping table (groundY/altitude → screen Y) and the altitude-ignored Z-sort case lock the observable contract.
- [Sonnet 4.6] [Feature Implementer] 2026-05-07 [FINISHED]: Green Phase complete. All tests pass (`go test ./...`). Production files modified/created:
  - `internal/engine/contracts/body/body.go` — added 4 methods to `Body` interface (Altitude, SetAltitude, Altitude16, SetAltitude16) and 4 methods to `Movable` interface (VAltitude16, SetVAltitude16, AccelerationAltitude, SetAccelerationAltitude)
  - `internal/engine/physics/body/body.go` — added `altitude16` field, altitude accessors, modified `Position()` to subtract altitude from Y (ScreenY = GroundY - Altitude)
  - `internal/engine/physics/body/body_movable.go` — added `vAltitude16` and `accAltitude` fields plus 4 accessors
  - `internal/engine/physics/body/obstacle.go` — added 4 forwarding altitude methods to `ObstacleRect` to resolve ambiguous selector (embeds both MovableBody and CollidableBody)
  - `internal/engine/entity/actors/character.go` — added 4 forwarding altitude methods to `Character` (same ambiguity)
  - `internal/engine/entity/items/item_base.go` — added 4 forwarding altitude methods to `BaseItem` (same ambiguity)
  - `internal/engine/render/draworder/draworder.go` — new package; `SortByGroundY` using `sort.SliceStable` on a copy, keyed on y16
  - `internal/game/scenes/phases/scene.go` — render loop now calls `draworder.SortByGroundY(space.Bodies())`
  - `internal/kit/actors/beatemup/doc.go` — new package doc
  - `internal/kit/actors/beatemup/beatemup_character.go` — new `BeatEmUpCharacter` scaffold embedding `MeleeCharacter`
  - Multiple `*_test.go` files — added altitude stub methods to test mocks that implement `body.Body` or `body.Movable`
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-05-07 [FINISHED]: All quality gates passed. One defect found and fixed: `internal/game/scenes/phases/scene.go` had the `draworder` import inserted out of alphabetical order (after `sequences`, should be after `render/camera`); fixed with `gofmt -w`. All 9 ACs verified against spec. TDD Red-Green-Refactor cycle confirmed (tests were written before production code; compile-time contract assertions in body_test.go and body_movable_test.go). Coverage: `internal/engine/physics/body/` = 92.4% (positive delta — new accessor and Position tests added coverage); `internal/engine/render/draworder/` = 100%; `internal/kit/actors/beatemup/` = 100%. Full suite: 50 packages, 0 failures. `golangci-lint` clean after gofmt fix. No `_ = variable` in production code. Layer rules honored (no engine->kit or engine->game imports introduced).
