# SPEC — US-037 — Per-State Projectile Spawn Offset

**Branch:** `037-per-state-projectile-spawn-offset`
**Bounded Context:** Engine — `internal/engine/combat/weapon/`; Game — `internal/game/entity/actors/player/` and entity JSON loader.

## Context

US-034 added a single static `(spawnOffsetX16, spawnOffsetY16)` to `ProjectileWeapon` that is applied to every shot, regardless of the actor's current state. Real action games (Metal Slug, Cuphead) need the muzzle position to track the gun-arm pose: a ducking shot exits low, an upward shot exits high, etc.

This story extends `ProjectileWeapon` with an optional per-state offset table, sourced from entity JSON (`assets/entities/player/climber.json` under `skills.shooting.state_spawn_offsets`). At fire time, the weapon looks up the actor's current state and overrides the default offset if a matching entry exists. The X-axis facing-direction negation rule from US-034 still applies.

To keep the engine→actors dependency direction clean, the weapon contract uses an opaque integer state token (the underlying type of `actors.ActorStateEnum`). The game-side loader is responsible for resolving JSON state names (`"duck"`, `"jump_shoot"`, …) to enum values via `actors.GetStateEnum()` before constructing the weapon.

## Technical Requirements

### 1. Weapon contract change

**File:** `internal/engine/contracts/combat/weapon.go`

`Fire` gains a trailing `state int` parameter. Using `int` (not a domain type) keeps the contract free of an engine→actors import. Callers in skill code pass the actor's current state enum as an `int`.

```go
type Weapon interface {
    ID() string
    Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int)
    CanFire() bool
    Update()
    Cooldown() int
    SetCooldown(frames int)
}
```

### 2. ProjectileWeapon struct + constructor

**File:** `internal/engine/combat/weapon/weapon.go`

Add a per-state offset table keyed by the integer state token. Values are `[2]int` storing `[x16, y16]` already in fp16.

```go
type ProjectileWeapon struct {
    // ... existing fields ...
    spawnOffsetX16 int
    spawnOffsetY16 int
    stateOffsets   map[int][2]int // NEW — nil-safe lookup
}
```

The existing 8-arg constructor `NewProjectileWeapon(...)` is preserved (no breaking change to callers). A new setter exposes the table:

```go
// SetStateSpawnOffsets registers per-state spawn offsets. Values are fp16 (x16, y16).
// Passing a nil or empty map clears all per-state overrides.
func (w *ProjectileWeapon) SetStateSpawnOffsets(offsets map[int][2]int)
```

Rationale: a setter (instead of a constructor parameter) keeps US-034 call sites untouched, makes the feature additive, and matches the existing pattern used by `SetVFXManager` and `SetOwner`.

### 3. Fire method update

```go
func (w *ProjectileWeapon) Fire(
    x16, y16 int,
    faceDir animation.FacingDirectionEnum,
    direction body.ShootDirection,
    state int,
) {
    offsetX16 := w.spawnOffsetX16
    offsetY16 := w.spawnOffsetY16

    if w.stateOffsets != nil {
        if override, ok := w.stateOffsets[state]; ok {
            offsetX16 = override[0]
            offsetY16 = override[1]
        }
    }

    if faceDir == animation.FaceDirectionLeft {
        offsetX16 = -offsetX16
    }

    spawnX16 := x16 + offsetX16
    spawnY16 := y16 + offsetY16

    // ... unchanged VFX + projectile spawn + cooldown reset ...
}
```

Key invariants:
- Lookup uses `map[int]` zero-allocation indexing; nil map is safe (skipped via the `!= nil` guard).
- The X-axis negation for `FaceDirectionLeft` is applied **after** the per-state lookup, so the JSON value is always written as if facing right (consistent with US-034 AC3).
- Y is always additive — no facing transform.

### 4. Caller update — shooting skill

**File:** `internal/engine/physics/skill/skill_shooting.go`

The shooting skill already holds a `body.MovableCollidable` (`b`). To pass state, extend the actor-side abstraction with an interface assertion. Because the skill must not import `actors`, perform a type-assertion to a local interface:

```go
type stateProvider interface {
    State() int
}
```

Note: `actors.ActorStateEnum` is defined as `type ActorStateEnum int`. A method declared as `State() ActorStateEnum` satisfies `State() int` only if we widen via a small adapter. To keep this clean, prefer assertion to a local interface that returns the concrete enum type via reflection-free conversion:

```go
type actorStateReader interface {
    State() actors.ActorStateEnum
}

state := 0
if sr, ok := b.(actorStateReader); ok {
    state = int(sr.State())
}
weapon.Fire(x16, y16, b.FaceDirection(), direction, state)
```

If importing `actors` from `engine/physics/skill` is undesirable (verify package boundaries), the alternate is to declare the local interface as `State() int` and add a thin shim method `StateInt() int { return int(c.State()) }` on `Character`. **Decision:** import `actors` from `skill_shooting.go` is acceptable because the skill package already operates at the engine→game seam and the cost of a shim outweighs the import. If lint rejects it, fall back to the shim.

### 5. JSON schema extension

**File:** `assets/entities/player/climber.json`

Extend the existing `skills.shooting` block:

```json
"shooting": {
  "enabled": true,
  "cooldown_frames": 15,
  "projectile_speed": 6,
  "projectile_range": 4,
  "directions": 4,
  "spawn_offset": { "x": 8, "y": 8 },
  "state_spawn_offsets": {
    "duck":       { "x": 6, "y": 12 },
    "jump_shoot": { "x": 8, "y": 5 },
    "fall_shoot": { "x": 8, "y": 6 }
  }
}
```

Note: `spawn_offset` already exists conceptually from US-034 but is currently hardcoded in `internal/game/entity/actors/player/weapons.go`. This story does not require migrating the default offset to JSON (out of scope); only `state_spawn_offsets` must be wired through.

### 6. Game-side loader

**File:** `internal/game/entity/actors/player/weapons.go` (and/or the climber entity factory that parses `climber.json`).

Add a helper that:

1. Reads `skills.shooting.state_spawn_offsets` from the parsed JSON struct.
2. For each entry, calls `actors.GetStateEnum(name)`:
   - If `ok == false` → log a warning (`log.Printf("US-037: unknown state %q in state_spawn_offsets, skipping", name)`) and skip.
   - If `ok == true` → convert pixel ints to fp16 via `fp16.To16(x)` and `fp16.To16(y)`, store into `map[int][2]int{int(enum): {x16, y16}}`.
3. Calls `light.SetStateSpawnOffsets(table)` on each constructed `ProjectileWeapon`.

If the JSON map is absent or empty, `SetStateSpawnOffsets` is not called and behavior is byte-for-byte identical to US-034.

## State Machine Transitions

No new actor states are introduced. The states referenced by JSON (`duck`, `jump_shoot`, `fall_shoot`, etc.) are already registered in `internal/engine/entity/actors/actor_state.go` (lines 42–55) and `ducking_state.go`. Resolution is by string → enum at load time, frozen for the lifetime of the weapon.

## Pre-conditions

- US-034 is merged: `ProjectileWeapon` has `spawnOffsetX16`/`spawnOffsetY16` and the 8-arg constructor.
- `actors.GetStateEnum(name string) (ActorStateEnum, bool)` exists and is exported.
- The shooting skill calls `weapon.Fire(...)` from a single site (`skill_shooting.go:64`).
- `ProjectileWeapon` is the only implementation of `combat.Weapon` (verified via grep — `MockWeapon` is the only other implementer and lives in `internal/engine/mocks/combat.go`).

## Post-conditions

- `combat.Weapon.Fire` accepts a 5th `state int` parameter; all implementations (production + mocks) updated.
- `ProjectileWeapon.SetStateSpawnOffsets(map[int][2]int)` exists and is the sole mutator for the table.
- When `state` matches a registered offset, the override replaces the default for that single shot (not persisted).
- When `state` has no entry (or table is nil), default offset is used → identical to US-034.
- Facing-left negates the resolved X offset (per-state or default) consistently.
- Unknown JSON state names emit a warning and are skipped; loader returns no error.
- Default behavior (no `state_spawn_offsets` in JSON, or `(0,0)` default offset) reproduces pre-US-037 behavior bit-for-bit.

## Integration Points

- **Engine package:** `internal/engine/combat/weapon/` (struct, constructor unchanged, new setter, modified `Fire`).
- **Engine contract:** `internal/engine/contracts/combat/weapon.go` (interface signature change).
- **Engine caller:** `internal/engine/physics/skill/skill_shooting.go` (passes state to `Fire`).
- **Engine mock:** `internal/engine/mocks/combat.go` — `MockWeapon.Fire` signature must be updated.
- **Game loader:** `internal/game/entity/actors/player/weapons.go` (or wherever `climber.json` shooting block is consumed).
- **Asset:** `assets/entities/player/climber.json` (additive `state_spawn_offsets` block).
- **Cross-context dep:** `internal/engine/entity/actors` (`GetStateEnum`, `ActorStateEnum`) — used only by the game-side loader and the skill caller, never by the weapon package itself.

## Acceptance Criteria Mapping

| AC | Mechanism | Verified by |
|---|---|---|
| AC1 — JSON `state_spawn_offsets` is optional | Loader checks for nil/empty map; absent → `SetStateSpawnOffsets` not called | Loader unit test: parse JSON without the key, assert weapon's table is nil |
| AC2 — Unknown state names log warning + skip | Loader iterates entries, calls `GetStateEnum`; on `!ok`, log + continue | Loader unit test: JSON contains `"bogus_state"`, assert it is absent from table and no error returned |
| AC3 — `stateOffsets map[int][2]int` field | New struct field + `SetStateSpawnOffsets` setter | Weapon unit test: construct, set table, assert lookup behavior via `Fire` |
| AC4 — `Fire` accepts state and overrides default | Modified `Fire` signature + lookup logic | Weapon unit test (table-driven): state with entry → override coords; state without entry → default coords |
| AC5 — Facing-left negation applies to per-state offsets | Negation applied **after** lookup | Weapon unit test: same per-state offset, two rows for `FaceDirectionLeft` and `FaceDirectionRight`, assert X sign flip |
| AC6 — `(0,0)` default + no per-state map = US-034 behavior | `SetStateSpawnOffsets` not called → `stateOffsets == nil` → guard short-circuits | Weapon unit test: existing US-034 test row passes with new `state` arg = 0 |
| AC7 — Unit tests cover: match, fallback, facing-left negation, unknown JSON state | Two test files: `weapon_test.go` (Fire behavior) + loader test for JSON | See Red Phase below |

## Red Phase

Two failing tests drive this story.

### Test A — Weapon-level behavior

**File:** `internal/engine/combat/weapon/weapon_test.go`
**New test function:** `TestProjectileWeapon_Fire_StateSpawnOffset`

Table-driven, asserting projectile spawn coordinates passed to the mock `ProjectileManager`:

| Row | Default offset (x16,y16) | State table | `state` arg | `faceDir` | Fire pos (x16,y16) | Expected spawn (x16,y16) | Covers |
|---|---|---|---|---|---|---|---|
| `state with matching offset uses it` | (80, 160) | `{42: (96, 192)}` | 42 | Right | (320, 480) | (416, 672) | AC4 |
| `state without entry falls back to default` | (80, 160) | `{42: (96, 192)}` | 99 | Right | (320, 480) | (400, 640) | AC4 |
| `facing left negates per-state X` | (80, 160) | `{42: (96, 192)}` | 42 | Left | (320, 480) | (224, 672) | AC5 |
| `nil state table uses default offset` | (80, 160) | nil (setter not called) | 0 | Right | (320, 480) | (400, 640) | AC6 |
| `(0,0) default + no table = no offset` | (0, 0) | nil | 0 | Right | (320, 480) | (320, 480) | AC6 |

**Expected failure modes:**
1. `weapon.Fire(...)` does not accept a 5th argument → compile error.
2. `ProjectileWeapon` has no `SetStateSpawnOffsets` method → compile error.
3. After signature/setter additions but before lookup logic, override rows fail equality assertions.

### Test B — Loader-level behavior

**File:** TBD by Mock Generator / TDD Specialist — likely `internal/game/entity/actors/player/weapons_test.go` or wherever the climber JSON loader is unit-tested. If no such test file exists, the TDD Specialist creates one.
**New test function:** `TestLoadShootingSkill_StateSpawnOffsets`

Table-driven over inline JSON fixtures asserting the resulting `map[int][2]int` registered on the weapon (exposed via a test-only accessor `StateSpawnOffsets() map[int][2]int` on `ProjectileWeapon`, or asserted indirectly by firing the constructed weapon).

| Row | JSON `state_spawn_offsets` | Expected table | Covers |
|---|---|---|---|
| `absent key produces nil table` | (omitted) | nil / empty | AC1 |
| `valid entries are registered with fp16 conversion` | `{"duck":{"x":6,"y":12}}` | `{int(Ducking): (96, 192)}` | AC1, AC3 |
| `unknown state name is skipped with warning` | `{"bogus":{"x":1,"y":1},"duck":{"x":2,"y":3}}` | `{int(Ducking): (32, 48)}` (size 1) | AC2 |
| `multiple states load independently` | `{"duck":{"x":1,"y":1},"jump_shoot":{"x":2,"y":2}}` | size 2 | AC1 |

**Expected failure modes:**
1. Loader does not parse the new key → resulting table is nil or wrong size.
2. Loader does not call `actors.GetStateEnum` → keys are wrong.
3. fp16 conversion missing → values are pixel ints, not fp16.

### Mock updates required

- `internal/engine/mocks/combat.go` — `MockWeapon.Fire` signature.
- `internal/engine/combat/weapon/mocks_test.go` — verify no local mock implements `combat.Weapon`; if any does, update.

These are mechanical compile-driven updates the Mock Generator stage will handle.

## Out of Scope

- Migrating the default `spawn_offset` from hardcoded Go to JSON (separate story if desired).
- Per-frame interpolation of the offset across animation frames.
- Per-direction-of-shot offsets (e.g. different muzzle for `ShootDirectionUp` vs `ShootDirectionStraight`); current design keys solely on actor state, which already encodes shoot direction via states like `JumpingShooting`/`FallingShooting`.
- Reloading the JSON at runtime (offsets are loaded once at weapon construction).
