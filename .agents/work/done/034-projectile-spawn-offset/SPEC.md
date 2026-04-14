# SPEC — US-034 — Projectile Spawn Offset Configuration

**Branch:** `034-projectile-spawn-offset`

## Context

Currently, `ProjectileWeapon.Fire()` spawns projectiles and muzzle flash VFX at the caller-provided entity origin (`x16`, `y16`) without any configurable offset. This means projectiles visually originate from the entity's center rather than from the weapon sprite position. This story adds configurable fp16 offset fields so game code can align the spawn point with the weapon sprite.

## Technical Requirements

### Constructor Signature Change
**File:** `internal/engine/combat/weapon/weapon.go`

Add two parameters to the end of `NewProjectileWeapon`:

```go
func NewProjectileWeapon(
    id string,
    cooldownFrames int,
    projectileType string,
    projectileSpeed int,
    manager combat.ProjectileManager,
    muzzleEffectType string,
    spawnOffsetX16 int, // NEW - fp16 units
    spawnOffsetY16 int, // NEW - fp16 units
) *ProjectileWeapon
```

### Struct Modification

Add two fields to `ProjectileWeapon`:

```go
type ProjectileWeapon struct {
    // ... existing fields ...
    spawnOffsetX16 int // NEW
    spawnOffsetY16 int // NEW
}
```

### Fire Method Update

Compute the effective spawn position by applying the offset. Negate the X offset when facing left. Use the offset position for both VFX and projectile spawning:

```go
func (w *ProjectileWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection) {
    offsetX16 := w.spawnOffsetX16
    if faceDir == animation.FaceDirectionLeft {
        offsetX16 = -offsetX16
    }

    spawnX16 := x16 + offsetX16
    spawnY16 := y16 + w.spawnOffsetY16

    if w.vfxManager != nil && w.muzzleEffectType != "" {
        x := float64(spawnX16) / 16.0
        y := float64(spawnY16) / 16.0
        w.vfxManager.SpawnPuff(w.muzzleEffectType, x, y, 1, 0.0)
    }

    vx16, vy16 := w.calculateVelocity(direction, faceDir)
    w.manager.SpawnProjectile(w.projectileType, spawnX16, spawnY16, vx16, vy16, w.owner)
    w.currentCooldown = w.cooldownFrames
}
```

### Existing Constructor Call Sites

All existing `NewProjectileWeapon` calls (production and test) must be updated to pass two additional trailing arguments. For backward compatibility (no offset), pass `0, 0`.

**Test file updates** (`weapon_test.go`, `factory_test.go`): Every call to `NewProjectileWeapon` with 6 args becomes 8 args by appending `, 0, 0`.

## Pre-conditions

- `ProjectileWeapon` has a 6-parameter constructor (current state after US-030).
- `Fire()` spawns at the caller-provided position without offset.
- Muzzle flash VFX is implemented and spawns at the caller-provided position (US-030).

## Post-conditions

- Constructor accepts 8 parameters (2 new offset fields).
- `Fire()` applies offset to both projectile and VFX spawn positions.
- X offset is negated when facing left; Y offset is always additive.
- Offset of `(0, 0)` produces identical behavior to pre-change code.
- All existing tests pass with `0, 0` offset appended.

## Integration Points

- **Package:** `internal/engine/combat/weapon/`
- **Contract:** `combat.ProjectileManager` (unchanged -- `SpawnProjectile` receives offset-adjusted coordinates).
- **Shared mock:** `mocks.MockVFXManager` from `internal/engine/mocks/` (unchanged).
- **Local mock:** `mockProjectileManager` in `weapon/mocks_test.go` (unchanged).
- **Callers:** Any game code constructing `ProjectileWeapon` must supply offset values. Factory functions in `weapon/factory.go` may need updating.

## Red Phase

### Test File
`internal/engine/combat/weapon/weapon_test.go`

### Failing Test: Table-Driven Spawn Offset Scenarios

A single table-driven test covers AC1 through AC6. Each sub-test constructs a weapon with specific offset values, fires it, and asserts the spawn coordinates passed to both `SpawnProjectile` and `SpawnPuff`.

```go
func TestProjectileWeapon_Fire_SpawnOffset(t *testing.T) {
    tests := []struct {
        name           string
        spawnOffsetX16 int
        spawnOffsetY16 int
        fireX16        int
        fireY16        int
        faceDir        animation.FacingDirectionEnum
        wantSpawnX16   int
        wantSpawnY16   int
        wantVFXX       float64
        wantVFXY       float64
        hasVFX         bool
    }{
        {
            name:           "facing right applies positive X offset",
            spawnOffsetX16: 128,
            spawnOffsetY16: 64,
            fireX16:        320,
            fireY16:        480,
            faceDir:        animation.FaceDirectionRight,
            wantSpawnX16:   448,  // 320 + 128
            wantSpawnY16:   544,  // 480 + 64
            wantVFXX:       28.0, // 448 / 16
            wantVFXY:       34.0, // 544 / 16
            hasVFX:         true,
        },
        {
            name:           "facing left negates X offset",
            spawnOffsetX16: 128,
            spawnOffsetY16: 64,
            fireX16:        320,
            fireY16:        480,
            faceDir:        animation.FaceDirectionLeft,
            wantSpawnX16:   192,  // 320 - 128
            wantSpawnY16:   544,  // 480 + 64
            wantVFXX:       12.0, // 192 / 16
            wantVFXY:       34.0, // 544 / 16
            hasVFX:         true,
        },
        {
            name:           "zero offset preserves current behavior",
            spawnOffsetX16: 0,
            spawnOffsetY16: 0,
            fireX16:        320,
            fireY16:        480,
            faceDir:        animation.FaceDirectionRight,
            wantSpawnX16:   320,
            wantSpawnY16:   480,
            wantVFXX:       20.0,
            wantVFXY:       30.0,
            hasVFX:         true,
        },
        {
            name:           "zero offset without VFX manager",
            spawnOffsetX16: 0,
            spawnOffsetY16: 0,
            fireX16:        320,
            fireY16:        480,
            faceDir:        animation.FaceDirectionRight,
            wantSpawnX16:   320,
            wantSpawnY16:   480,
            hasVFX:         false,
        },
        {
            name:           "negative Y offset shifts spawn upward",
            spawnOffsetX16: 0,
            spawnOffsetY16: -32,
            fireX16:        160,
            fireY16:        320,
            faceDir:        animation.FaceDirectionRight,
            wantSpawnX16:   160,
            wantSpawnY16:   288,  // 320 - 32
            wantVFXX:       10.0, // 160 / 16
            wantVFXY:       18.0, // 288 / 16
            hasVFX:         true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var gotProjX16, gotProjY16 int
            mockProj := &mockProjectileManager{
                SpawnProjectileFunc: func(_ string, x16, y16, _, _ int, _ interface{}) {
                    gotProjX16 = x16
                    gotProjY16 = y16
                },
            }

            w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, mockProj, "muzzle_flash", tt.spawnOffsetX16, tt.spawnOffsetY16)

            var gotVFXX, gotVFXY float64
            var vfxCalled bool
            if tt.hasVFX {
                mockVFX := &mocks.MockVFXManager{
                    SpawnPuffFunc: func(_ string, x, y float64, _ int, _ float64) {
                        vfxCalled = true
                        gotVFXX = x
                        gotVFXY = y
                    },
                }
                w.SetVFXManager(mockVFX)
            }

            w.Fire(tt.fireX16, tt.fireY16, tt.faceDir, body.ShootDirectionStraight)

            if gotProjX16 != tt.wantSpawnX16 {
                t.Errorf("projectile x16: got %d, want %d", gotProjX16, tt.wantSpawnX16)
            }
            if gotProjY16 != tt.wantSpawnY16 {
                t.Errorf("projectile y16: got %d, want %d", gotProjY16, tt.wantSpawnY16)
            }

            if tt.hasVFX {
                if !vfxCalled {
                    t.Fatal("expected SpawnPuff to be called")
                }
                if gotVFXX != tt.wantVFXX {
                    t.Errorf("VFX x: got %v, want %v", gotVFXX, tt.wantVFXX)
                }
                if gotVFXY != tt.wantVFXY {
                    t.Errorf("VFX y: got %v, want %v", gotVFXY, tt.wantVFXY)
                }
            }
        })
    }
}
```

**Expected failure:** `NewProjectileWeapon` currently accepts 6 parameters; calling it with 8 will fail to compile. Once the constructor is updated, `Fire()` will not yet apply the offset, so position assertions will fail.

### No Interface Changes Required

- `combat.ProjectileManager` is unchanged (it already receives absolute coordinates).
- `vfx.Manager` is unchanged.
- No new contracts are needed.
