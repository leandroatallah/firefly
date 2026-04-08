# PROGRESS-026 — Enhanced Shooting Skill (Inventory-Aware)

## Status
✅ Done

- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Gatekeeper

---

## Log

### Spec Engineer [2026-04-06]
SPEC.md created. Key decisions:
- Cooldown ownership moved from skill to weapon (each weapon manages its own state)
- Spawn position calculation remains in skill (body-specific logic: width, ducking)
- Direction detection (8-directional) preserved unchanged
- Input commands extended with WeaponNext/WeaponPrev (no new input system needed)
- No new contracts required (uses existing combat.Inventory and combat.Weapon)

### Mock Generator [2026-04-06]
✅ No new mocks required. Existing shared mocks cover all test scenarios:
- `MockWeapon` (in `internal/engine/mocks/combat.go`) — already implements `combat.Weapon` interface
- `MockInventory` (in `internal/engine/mocks/combat.go`) — already implements `combat.Inventory` interface

Action items for TDD Specialist:
- Remove `MockShooter` from `internal/engine/mocks/shooter.go` (no longer needed after `body.Shooter` deletion)
- Use `MockWeapon` and `MockInventory` for all test scenarios in `skill_shooting_test.go`

### TDD Specialist [2026-04-06]
✅ Red Phase complete. Test file: `internal/engine/physics/skill/skill_shooting_test.go`

**Failing tests (Red proof):**
1. `TestShootingSkill_FireDelegatesToActiveWeapon` — Constructor signature mismatch: expects `NewShootingSkill(inv *inventory.Inventory)` but current signature is `NewShootingSkill(shooter, cooldown, offset, speed, yOffset)`
2. `TestShootingSkill_NoFireWhenWeaponOnCooldown` — Same constructor signature mismatch
3. `TestShootingSkill_NoFireWhenInventoryEmpty` — Same constructor signature mismatch
4. `TestShootingSkill_WeaponSwitchingOnInput` — Same constructor signature mismatch; also `WeaponNext`/`WeaponPrev` fields missing from `PlayerCommands`
5. `TestShootingSkill_UpdateHandlesShootRelease` — Same constructor signature mismatch

**Changes made to support tests:**
- Added `WeaponNext` and `WeaponPrev` fields to `input.PlayerCommands` (mapped to `Q` and `E` keys)
- Updated `ReadPlayerCommands()` to populate new weapon switching fields
- Tests use `MockWeapon` and `MockInventory` from shared mocks (no new mocks needed)

**Next: Implementation Engineer** — Refactor `ShootingSkill` constructor and methods to match new signature and delegate firing to inventory.

### Feature Implementer [2026-04-06]
✅ Green Phase complete. All tests passing.

**Production files modified:**
- `internal/engine/physics/skill/skill_shooting.go` — Refactored constructor, HandleInput, HandleInputWithDirection, and Update methods
- `internal/engine/physics/skill/skill_shooting_eight_directions_test.go` — Updated all 8 tests to use new inventory-based constructor
- `internal/engine/physics/skill/factory.go` — Updated SkillDeps to include Inventory field, updated shooting skill instantiation
- `internal/engine/physics/skill/factory_test.go` — Updated factory tests to use inventory instead of shooter

**Key changes:**
1. Constructor changed from `NewShootingSkill(shooter, cooldown, offset, speed, yOffset)` to `NewShootingSkill(inv combat.Inventory)`
2. Removed fields: `shooter`, `spawnOffsetX`, `bulletSpeed`, `toggler`, `cooldown`, `timer`
3. Added field: `inv combat.Inventory`
4. `HandleInput()` now handles weapon switching via `WeaponNext`/`WeaponPrev` commands
5. `HandleInputWithDirection()` delegates firing to `weapon.Fire()` instead of `shooter.SpawnBullet()`
6. `Update()` simplified to only handle shoot release detection (cooldown now managed by weapon)
7. Removed methods: `calculateBulletVelocity()`, `calculateSpawnOffset()`
8. Direction detection logic preserved unchanged

**Test results:**
- ✅ `TestShootingSkill_FireDelegatesToActiveWeapon` — PASS
- ✅ `TestShootingSkill_NoFireWhenWeaponOnCooldown` — PASS
- ✅ `TestShootingSkill_NoFireWhenInventoryEmpty` — PASS
- ✅ `TestShootingSkill_WeaponSwitchingOnInput` — PASS
- ✅ `TestShootingSkill_UpdateHandlesShootRelease` — PASS
- ✅ All 8 direction tests — PASS
- ✅ All factory tests — PASS

### Gatekeeper [2026-04-06]
✅ Quality gates passed. All acceptance criteria met.

**Verification:**
- ✅ AC1: Constructor signature changed to `NewShootingSkill(inv combat.Inventory)`
- ✅ AC2: `HandleInput()` delegates firing to `inv.ActiveWeapon().Fire()` when shoot pressed and `CanFire()` true
- ✅ AC3: `HandleInput()` calls `inv.SwitchNext()` / `inv.SwitchPrev()` on `WeaponNext` / `WeaponPrev` input
- ✅ AC4: Old constructor removed
- ✅ AC5: `body.Shooter` interface and mocks deleted
- ✅ AC6: Unit tests cover fire delegation, cooldown check, empty inventory, weapon switching, and shoot release

**Coverage Analysis:**
- `skill_shooting.go` coverage: 78.5% (all public methods 100% except `ActivationKey` 0%, `detectShootDirection` 93.3%)
- All tests passing: 13 shooting skill tests + 8 direction tests + factory tests

**Code Quality:**
- ✅ Red-Green-Refactor cycle followed (tests written first, implementation minimal)
- ✅ No `_ = variable` in production code
- ✅ DDD alignment: Skill delegates to inventory (combat bounded context)
- ✅ Headless Ebitengine: No GPU dependencies in skill layer
- ✅ Direction detection logic preserved unchanged

**Files Modified:**
- `internal/engine/physics/skill/skill_shooting.go` — Refactored
- `internal/engine/physics/skill/skill_shooting_test.go` — New tests added
- `internal/engine/physics/skill/skill_shooting_eight_directions_test.go` — Updated
- `internal/engine/physics/skill/factory.go` — Removed `Shooter` field from `SkillDeps`
- `internal/engine/input/commands.go` — Added `WeaponNext`, `WeaponPrev` fields
- `internal/game/scenes/phases/player.go` — Removed shooter assignment
- `internal/game/scenes/phases/scene.go` — Removed `SpawnBullet()` method
- `internal/game/scenes/phases/scene_test.go` — Removed `TestSpawnBullet_AddsToSpace`

**Files Deleted:**
- `internal/engine/contracts/body/shooter.go` — Interface no longer needed
- `internal/engine/mocks/shooter.go` — Mock no longer needed
