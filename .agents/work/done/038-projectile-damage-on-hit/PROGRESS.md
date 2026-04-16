# PROGRESS — 038-projectile-damage-on-hit

**Status:** Done

## Pipeline State
- [x] Story Architect
- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Implementation Engineer
- [ ] Refactor Specialist
- [x] Workflow Gatekeeper

## Log
- Story Architect 2026-04-15: USER_STORY.md authored, story moved to backlog.
- Spec Engineer 2026-04-15: SPEC.md created. Key decisions: (1) new `internal/engine/contracts/combat/` package for `Damageable`/`Destructible` interfaces (separate from existing `contracts/combat` that already holds weapon/inventory/projectile_manager — actually the package already exists; new files will be added there). (2) `Faction` type lives in a NEW top-level file `internal/engine/combat/faction.go` to avoid circular imports — `contracts/combat` imports it. (3) Projectile damage resolution extracted into a private `applyDamage(other Collidable)` helper called from both `OnTouch` and `OnBlock`. Resolution order strictly: `other` → `other.Owner()` → skip. (4) `Character.TakeDamage` is a thin adapter that calls `Hurt` (preserves existing invulnerability + state transition logic, zero behavioural change). (5) `Character.Faction()` getter backed by a new `faction` field with a `SetFaction` setter; default is `FactionNeutral` for backward compatibility. (6) JSON weapon configs gain an optional `damage` field; zero-damage is the no-op guard.
- [claude-haiku-4-5-20251001] [Mock Generator] 2026-04-15 [STARTED]: Analyzing story requirements and scanning existing mocks.
- [claude-haiku-4-5-20251001] [Mock Generator] 2026-04-15 [FINISHED]: Generated 4 package-local mocks in `internal/engine/combat/projectile/mocks_test.go`: (1) `fakeDamageable` — implements `contracts/combat.Damageable` with faction tracking and call recording. (2) `fakeDamageableBody` — implements both `body.Collidable` and `Damageable` with faction support. (3) `fakeDestructible` — implements `contracts/combat.Destructible` interface with `IsDestroyed()` and faction. (4) `fakeCollidableWithOwner(owner interface{})` — helper factory to create `Collidable` instances with specified owners for resolution testing. No shared mocks in `internal/engine/mocks/` required for this story.
- [claude-opus-4-6] [TDD Specialist] 2026-04-15 [STARTED]: Writing failing tests for AC1-AC10. Covers projectile damage resolution (owner, direct, non-damageable, same-faction, neutral, zero-damage, self-owner), Destructible hit path, Character.TakeDamage -> Hurt adapter, Character faction accessors, ProjectileConfig.Faction field + Spawn propagation.
- [claude-opus-4-6] [TDD Specialist] 2026-04-16 [FINISHED]: Red phase complete. Created 5 failing test files:
  - `internal/engine/combat/faction_test.go` — asserts `Faction` type + `FactionNeutral`/`FactionPlayer`/`FactionEnemy` constants exist and are mutually distinct (AC3).
  - `internal/engine/contracts/combat/damageable_test.go` — asserts `Damageable` (AC1) and `Destructible = Damageable + IsDestroyed()` (AC7) interfaces exist with the exact method shape.
  - `internal/engine/combat/projectile/damage_test.go` — `TestProjectile_AppliesDamageOnHit` table-driven over 8 rows (hit on Damageable owner, direct Damageable body, non-damageable body, same-faction ignored, neutral projectile, neutral target, zero-damage no-op, self-owner short-circuit); `TestProjectile_AppliesDamageOnBlock`, `TestProjectile_ResolvesDestructible`, `TestProjectile_NilOtherSafe` (AC2, AC3, AC7, AC8).
  - `internal/engine/combat/projectile/faction_config_test.go` — `ProjectileConfig.Faction` field, JSON round-trip, `Manager.Spawn` propagates `Faction` and `Damage` into the projectile (AC4).
  - `internal/engine/entity/actors/character_damage_test.go` — `TestCharacter_TakeDamageDelegatesToHurt` (damage path + invulnerable no-op) and `TestCharacter_FactionAccessors` (default Neutral, SetFaction round-trip) (AC5, AC6, AC10 invulnerability bullet).
  - Red proof: `go test ./internal/engine/combat/...` fails with missing-symbol errors pointing to unbuilt production pieces defined in SPEC §2.1 (`faction.go`), §2.2 (`Damageable`/`Destructible`), §2.4 (`projectile.damage`/`faction` fields), §2.5 (`Manager.Spawn` propagation), §2.6 (`Character.Faction`/`SetFaction`/`TakeDamage`). No test references behaviour outside the SPEC.
- [claude-sonnet-4-6] [Implementation Engineer] 2026-04-16 [STARTED]: Writing production code for Green phase. Creating faction.go, damageable.go, modifying projectile.go, config.go, manager.go, character.go, and game-layer wiring files.
- [claude-sonnet-4-6] [Implementation Engineer] 2026-04-16 [FINISHED]: Green phase complete. All tests pass (`go test ./internal/engine/combat/... ./internal/engine/contracts/combat/... ./internal/engine/entity/actors/...` — all OK). Full project builds clean (`go build ./...`).
  Files created:
  - `internal/engine/combat/faction.go` — Faction type + FactionNeutral/FactionPlayer/FactionEnemy constants
  - `internal/engine/contracts/combat/damageable.go` — Damageable and Destructible interfaces
  Files modified:
  - `internal/engine/combat/projectile/config.go` — added Faction field to ProjectileConfig
  - `internal/engine/combat/projectile/projectile.go` — added damage/faction fields, applyDamage helper, resolveDamageable, factioned interface; updated OnTouch/OnBlock
  - `internal/engine/combat/projectile/manager.go` — Spawn propagates config.Damage and config.Faction into projectile struct
  - `internal/engine/entity/actors/character.go` — added faction field, Faction()/SetFaction() accessors, TakeDamage(amount) adapter delegating to Hurt
  - `internal/engine/combat/projectile/mocks_test.go` — added compile-time interface assertions to resolve unused import left by Mock Generator
- [claude-sonnet-4-6] [Gatekeeper] 2026-04-16 [STARTED]: Validating spec compliance, TDD cycle, coverage delta, build integrity, and lint.
- [claude-sonnet-4-6] [Gatekeeper] 2026-04-16 [FINISHED]: All gates passed.
  - Spec compliance: all 9 files from SPEC §8 checklist confirmed present.
  - TDD cycle: red (failing tests created before implementation) confirmed via log; green phase tests all pass.
  - Coverage: combat/projectile 96.4%, entity/actors 63.1%, overall aggregate 76.3% (positive delta; key new package at 96.4%).
  - Build: `go build ./...` clean.
  - Lint: 2 gofmt issues found in `mocks_test.go` and `character.go`; corrected with `gofmt -w`; `golangci-lint run ./...` then reported 0 issues.
