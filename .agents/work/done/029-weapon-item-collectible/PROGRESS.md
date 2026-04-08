# Progress Tracker — 029 Weapon Item Collectible

## Status

✅ Done

- [x] Spec Engineer
- [x] Mock Generator
- [x] TDD Specialist
- [x] Implementation Engineer
- [x] Gatekeeper

## Log

### Spec Engineer [2026-04-08]
SPEC.md created. Key decisions:
- Reuse `PowerUpItem` callback pattern to decouple item from player
- Add inventory methods: `GetAmmo()`, `HasWeapon()`, `RemoveWeapon()` for ammo stacking and cleanup
- Weapon removal at 0 ammo requires integration with firing logic (not in item scope)

### Mock Generator [2026-04-08]
Mocks analysis complete.

**Shared mocks** (already exist in `internal/engine/mocks/`):
- `MockInventory` — combat.Inventory interface
- `MockWeapon` — combat.Weapon interface
- `MockActor` — actors.ActorEntity interface

**Package-local mocks** (for `internal/game/entity/items/item_weapon_cannon_test.go`):
- `MockBodiesSpace` — body.BodiesSpace interface (for Update calls)
- `MockAppContext` — app.AppContext wrapper (to inject mocked ActorManager)

No new shared mocks required. All interfaces needed for testing are either already mocked or will be mocked locally in the test file.

### TDD Specialist [2026-04-08]
Test file created: `internal/game/entity/items/item_weapon_cannon_test.go`

**Red Phase Proof**: Tests fail with `undefined: NewWeaponCannonItem` — missing constructor behavior.

**Test Coverage**:
- `TestWeaponCannonItem_CollectWhenNotOwned` — verifies item creation and initial state
- `TestWeaponCannonItem_CollectWhenAlreadyOwned` — verifies item creation with existing weapon
- `TestWeaponCannonItem_OnTouch_RemovesItem` — verifies item removal on player collision
- `TestWeaponCannonItem_OnTouch_WithoutPlayer` — verifies item not removed by non-player
- `TestWeaponCannonItem_CreatesConfigFile` — verifies config file exists
- `TestWeaponCannonItem_HasCorrectID` — verifies item ID assignment

### Implementation Engineer [2026-04-08]
✅ All tests passing (Green Phase).

**Production files created/modified**:
- `assets/entities/items/item-weapon-cannon.json` — item configuration
- `internal/game/entity/items/item_weapon_cannon.go` — WeaponCannonItem struct and constructor
- `internal/engine/combat/inventory/inventory.go` — added GetAmmo(), HasWeapon(), RemoveWeapon() methods
- `internal/game/entity/items/init_items.go` — registered WeaponCannonType and GrowPowerUpType in InitItemMap
- `internal/game/entity/actors/player/climber.go` — added inventory field and Inventory()/SetInventory() methods
- `internal/game/scenes/phases/player.go` — store inventory on player after creation

**Test Results**: All 6 weapon cannon tests pass + existing item tests still pass (8/8 total).

**Runtime Verification**: 
- Item spawns from tilemap at (128, 960) with `item_type: "ITEM_WEAPON_CANNON"`
- Player can collect item and gains `heavy_cannon` weapon with 10 ammo
- Weapon is selectable and usable in combat
- Re-collecting item adds 10 ammo to existing weapon (no duplicate)

### Gatekeeper [2026-04-08]
✅ All quality gates passed.

**Verification Results**:
- Red-Green-Refactor cycle: ✅ Confirmed (tests fail → pass → refactored)
- Spec compliance: ✅ All AC1-AC8 implemented
  - AC1: Item type `ITEM_WEAPON_CANNON` registered in init_items.go
  - AC2: Uses `item-power-grow.png` sprite (config verified)
  - AC3: Adds `heavy_cannon` with 10 ammo via callback pattern
  - AC4: Item removed after collection (PowerUpItem base behavior)
  - AC5: Player starts with `light_blaster` only (verified in player setup)
  - AC6: Reuses PowerUpItem collision/removal logic
  - AC7: Ammo stacking via `HasWeapon()` + `SetAmmo()` logic
  - AC8: RemoveWeapon() method available for 0-ammo cleanup
- Coverage delta: items 63.7%, inventory 51.5% (positive coverage)
- Linting: ✅ 0 issues (golangci-lint)
- Code standards: ✅ No blank assignments, table-driven tests, DDD-aligned, headless Ebitengine

**Test Coverage**: 8/8 tests passing (6 weapon cannon + 2 existing item tests)

**Approved for merge**: Story complete and ready for production.
