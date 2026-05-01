# PROGRESS — 045-projectile-friendly-fire-prevention

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Sonnet 4.6] [Story Architect] 2026-04-30 [STARTED]: Writing USER_STORY.md for projectile-friendly-fire-prevention bug.
- [Sonnet 4.6] [Story Architect] 2026-04-30 [FINISHED]: USER_STORY.md created with 4 acceptance criteria, out-of-scope section, and domain notes.
- [Sonnet 4.6] [Story Architect] 2026-04-30 [FINISHED]: USER_STORY.md amended — added Edge Cases section and a 5th acceptance criterion covering opt-in interceptability for future explosive projectile types (e.g. rockets).
- [Opus 4.7] [Spec Engineer] 2026-04-30 [STARTED]: Authoring SPEC.md — defining contract `body.Projectile` with `Interceptable()` flag and collision filtering at projectile-level via existing `Touchable.OnTouch`/`OnBlock`.
- [Opus 4.7] [Spec Engineer] 2026-04-30 [FINISHED]: SPEC.md written. Key decisions: (1) filtering lives in `projectile.OnTouch`/`OnBlock`, not in Space — preserves bounded contexts; (2) new contract `internal/engine/contracts/body/projectile.go` with single-method `Projectile { Interceptable() bool }`, sibling to `Passthrough`; (3) introduce unexported `projectileBody` wrapper in `combat/projectile/` so the trait is discoverable directly on the body; (4) `ProjectileConfig.Interceptable` defaults to false — no behaviour change for existing bullets; (5) interceptability is a target-property — interceptable rocket gets shot down, default bullet survives. Mocks not required (skip Mock Generator); next step is TDD Specialist.
- [Opus 4.7] [TDD Specialist] 2026-04-30 [STARTED]: Authoring failing tests for AC1–AC5 in `internal/engine/combat/projectile/`.
- [Opus 4.7] [TDD Specialist] 2026-04-30 [FINISHED]: Wrote three new test files — `friendly_fire_test.go` (unit tests for AC1–AC5: ignore-other-default-projectile on OnTouch/OnBlock, actor-still-hit, melee non-regression, ordering determinism, interceptable target IS hit, interceptable shooter still ignores default bullet), `friendly_fire_integration_test.go` (AC1 integration with real `space.Space` — two overlapping default projectiles both survive a `Manager.Update` frame), and `interceptable_config_test.go` (AC5 config — zero-value `ProjectileConfig{}.Interceptable==false` and `Manager.Spawn` produces a body satisfying `body.Projectile`). Added `mockProjectileBody` test double satisfying the future `body.Projectile` contract. Red proven via `go test ./internal/engine/combat/projectile/...`: build fails on missing symbols `body.Projectile` and `ProjectileConfig.Interceptable` — exactly the surface the Feature Implementer must introduce. Once the symbols exist, the table-driven assertions become the behavior gate (no `QueueForRemoval`, no `SpawnPuff`, no `TakeDamage` between two default projectiles).
- [Sonnet 4.6] [Feature Implementer] 2026-04-30 [STARTED]: Implementing production code to satisfy all failing tests in the Green phase.
- [Sonnet 4.6] [Feature Implementer] 2026-04-30 [FINISHED]: All tests pass (`go test ./internal/engine/combat/projectile/...` → ok). Production files modified/created:
  - `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/contracts/body/projectile.go` (created — `Projectile` interface with `Interceptable() bool`)
  - `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/combat/projectile/config.go` (modified — added `Interceptable bool` field to `ProjectileConfig`)
  - `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/combat/projectile/projectile.go` (modified — added `interceptable` field and `Interceptable()` method on `*projectile`; added `projectileBody` wrapper type; added `isProjectile` helper; updated `OnTouch`/`OnBlock` to skip non-interceptable projectile targets)
  - `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/combat/projectile/manager.go` (modified — `Spawn` wires `config.Interceptable` and wraps body as `projectileBody`)
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-04-30 [STARTED]: Validating implementation against SPEC.md, running tests and coverage.
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-04-30 [FINISHED]: All quality gates passed. TDD Red phase confirmed (build failed on missing symbols before implementation). All tests pass. Coverage: 90.4% (delta: positive from build-failure baseline; well above 80% target). golangci-lint: 0 issues. No standards violations found. Story moved to done.

## Test Files

- `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/combat/projectile/friendly_fire_test.go`
- `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/combat/projectile/friendly_fire_integration_test.go`
- `/Users/leandroatallah/www/go/ebiten/045-projectile-friendly-fire-prevention/internal/engine/combat/projectile/interceptable_config_test.go`
