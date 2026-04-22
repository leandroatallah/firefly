# PROGRESS — 040-player-melee-attack

**Status:** Done

## Pipeline State
- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log
# Format: - [Model] [Agent] [date]: Action/Decision
- [Opus 4.7] Spec Engineer 2026-04-19: STARTED — moved story from backlog to active; reviewed constitution, USER_STORY, existing Weapon/Damageable/Faction contracts, ProjectileWeapon, ShootingSkill, DashState, GroundedState, BodiesSpace.Query API.
- [Opus 4.7] Spec Engineer 2026-04-19: FINISHED — SPEC.md created. Key decisions: (1) `MeleeWeapon` implements existing `combat.Weapon` with no interface change; adds melee-only `ApplyHitbox(space)` / `IsHitboxActive()` surface. (2) `StateMeleeAttack.returnTo` captured at construction so air melee returns to `Falling` (AC8). (3) Per-swing `hitThisSwing` guard makes damage deterministic across the active-frame window. (4) Factory validates `active_frames` and `hitbox` at load time so content errors fail loud. (5) Melee is driven directly by `PlayerCommands.Melee` + state-machine transitions, NOT through `ShootingSkill`, because melee is a stateful actor node while `ShootingSkill` is a projectile pipeline.
- [Opus 4.7] TDD Specialist 2026-04-19: STARTED — reviewing SPEC §4 RED scenarios, existing weapon/factory/grounded-state test patterns, mocks, contracts (Weapon, Damageable, Factioned, BodiesSpace.Query).
- [Opus 4.7] TDD Specialist 2026-04-19: FINISHED — Red phase complete. Test files written:
  - `internal/engine/combat/weapon/melee_test.go` (new): `TestMeleeWeapon_Fire_HitboxActivation` (table-driven frame-by-frame), `TestMeleeWeapon_ApplyHitbox_FactionGating`, `TestMeleeWeapon_ApplyHitbox_SingleHitPerSwing`, `TestMeleeWeapon_Cooldown_PreventsRefire`, `TestMeleeWeapon_Fire_MirrorsHitboxWhenFacingLeft`. Package-local `fakeSpace` implements `body.BodiesSpace` to make `Query` behavior observable.
  - `internal/engine/combat/weapon/factory_test.go` (extended): `TestWeaponFactory_MeleeJSON` table with cases `melee weapon ok` / `melee missing hitbox` / `melee inverted active_frames` / `unknown type` regression; renamed pre-existing unknown-type test to avoid shadowing.
  - `internal/game/entity/actors/states/melee_state_test.go` (new): `TestMeleeAttackState_ReturnsToGrounded_WhenAnimationFinishes`, `TestMeleeAttackState_AirMelee_ReturnsToFalling`, `TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow`, `TestGroundedState_MeleePressed_TransitionsToMeleeAttack`, `TestGroundedState_DashPressed_TakesPrecedenceOverMelee`, `TestMeleeTrigger_BlockedDuringCooldown`.
  - `internal/game/entity/actors/states/mocks_test.go` (extended): added `MeleePressedFunc` / `MeleePressed()` to `MockInputSource` so the `GroundedInput` contract extension is enforced.
  - `internal/engine/input/commands_test.go`: RED-4 coverage gap review — `PlayerCommands.Melee` ↔ `KeyZ` mapping is already verified (cases at lines 26-33, 122-135); no additions needed.
  Red proof: `go test ./internal/engine/combat/weapon/... ./internal/game/entity/actors/states/...` fails at the build stage with behavior-gap errors only — undefined `weapon.MeleeWeapon`, `weapon.NewMeleeWeapon`, `gamestates.NewMeleeAttackState`, `gamestates.StateMeleeAttack`, `gamestates.TryMeleeFromFalling`. Each failure maps 1:1 to a SPEC §1 symbol that the Feature Implementer must build; no test fails for a trivial signature mismatch unrelated to observable behavior.
- [GPT-5.2] Feature Implementer 2026-04-19: STARTED — implementing `MeleeWeapon` swing/cooldown/hitbox logic, factory JSON parsing/validation for `type=melee`, and verifying `StateMeleeAttack` + grounded trigger wiring against red tests.
- [Sonnet 4.6] Feature Implementer 2026-04-19: FINISHED — all tests green. Production files verified in place:
  - `internal/engine/combat/weapon/melee.go`: `MeleeWeapon` struct with `NewMeleeWeapon`, `Fire`, `Update`, `CanFire`, `Cooldown`, `SetCooldown`, `SetOwner`, `IsHitboxActive`, `ApplyHitbox` (faction-gated, per-swing `hitThisSwing` guard), `hitboxRect` (mirrors offset when `FaceDirectionLeft`).
  - `internal/engine/combat/weapon/factory.go`: `parseMeleeWeapon` branch for `"type":"melee"` with validation of `active_frames` bounds and non-zero `hitbox` dimensions; pixel-to-fp16 conversion via `fp16.To16`.
  - `internal/game/entity/actors/states/melee_state.go`: `MeleeAttackState` with `SetAnimationFrames`, `OnStart` (fires weapon), `Update` (applies hitbox while active, transitions to `returnTo` when animation done), `TryMeleeFromFalling` helper.
  - `internal/game/entity/actors/states/grounded_state.go`: `MeleePressed()` checked after `DashPressed` (dash takes precedence), returns `StateMeleeAttack`.
  - `internal/game/entity/actors/states/grounded_input.go`: `MeleePressed() bool` added to `GroundedInput` interface.
  All tests pass: `go test ./internal/engine/combat/weapon/... ./internal/game/entity/actors/states/...` — ok. `internal/engine/audio` pre-existing concurrent map failure is unrelated to this story.
- [Sonnet 4.6] Workflow Gatekeeper 2026-04-21: STARTED — reading PROGRESS, USER_STORY, SPEC; verifying pipeline stages, spec compliance, TDD cycle, and code quality.
- [Sonnet 4.6] Workflow Gatekeeper 2026-04-21: FINISHED — all gates passed.
  Coverage delta (positive): `internal/engine/combat/weapon` 91.4%; `internal/game/entity/actors/states` 80.3%. Both exceed the 80% threshold.
  Lint fixes applied: gofmt on `melee.go` and `factory.go` (alignment whitespace); `//nolint:unparam` added inline to `newMeleeOwner` and `newPlayerOwner` test helpers.
  golangci-lint: 0 issues.
  Story folder moved to done/.
