# SPEC 009 — Variable Jump Height

**Story:** `USER_STORY.md`
**Branch:** `009-variable-jump-height`

---

## Context

`JumpSkill` (`internal/engine/physics/skill/skill_platform_jump.go`) already owns jump activation. It calls `body.TryJump(force)` on key-press and ticks every frame via `Update`. It is the single correct place to add jump-cut logic — no other file needs to change.

The `Movable` contract already exposes `Velocity() / SetVelocity()`, `IsGoingUp()`, `JumpForceMultiplier()`. No new contract methods are required.

---

## Technical Requirements

### 1. `JumpSkill` — new fields

```
jumpCutMultiplier float64  // range (0.0, 1.0], default 1.0 (no cut)
jumpCutPending    bool     // true from TryJump until cut fires or body stops going up
```

### 2. `NewJumpSkill` — initialise defaults

```go
jumpCutMultiplier: 1.0,
```

### 3. `SetJumpCutMultiplier(m float64)` — public setter

- Clamp: if `m <= 0`, set to `0.1`; if `m > 1`, set to `1.0`.

### 4. `HandleInput` — detect key release

After the existing `IsKeyJustPressed` block, add:

```
if inpututil.IsKeyJustReleased(activationKey) && jumpCutPending {
    applyJumpCut(body)
}
```

### 5. `tryActivate` — set flag on successful jump

After `body.TryJump(force)`:

```go
s.jumpCutPending = true
```

### 6. `applyJumpCut(body)` — private helper

```
if body.IsGoingUp() {
    vx, vy := body.Velocity()
    body.SetVelocity(vx, int(float64(vy) * s.jumpCutMultiplier))
}
s.jumpCutPending = false
```

### 7. `Update` — clear flag when body stops going up

Inside the existing `Update`, after `SkillBase.Update`:

```
if jumpCutPending && !body.IsGoingUp() {
    jumpCutPending = false
}
```

---

## Pre-conditions

- `JumpSkill` has been activated this jump cycle (`jumpCutPending == true`).
- `body.IsGoingUp()` returns `true` (`vy16 < 0`).

## Post-conditions

- `vy16` is multiplied by `jumpCutMultiplier` exactly once.
- `jumpCutPending` is `false`.
- If `body.IsGoingUp()` is already `false` on release, `vy16` is unchanged.

---

## Integration Points

| Layer | File | Change |
|---|---|---|
| Physics / Skill | `internal/engine/physics/skill/skill_platform_jump.go` | Add fields, `SetJumpCutMultiplier`, `applyJumpCut`, wire into `HandleInput`, `tryActivate`, `Update` |
| Game Logic | `internal/game/entity/actors/player/climber.go` | Call `jumpSkill.SetJumpCutMultiplier(cfg)` during player setup (value from config or hardcoded constant) |

No contract changes. No new packages.

---

## Red Phase — Failing Test Scenario

**File:** `internal/engine/physics/skill/skill_platform_jump_test.go`

**Test name:** `TestJumpSkill_JumpCut`

**Scenario table:**

| name | jumpCutMultiplier | releaseWhileGoingUp | vy16Before | wantVy16After | wantPending |
|---|---|---|---|---|---|
| full hold — no cut | 1.0 | false | -320 | -320 | false |
| short press — cut applied | 0.5 | true | -320 | -160 | false |
| release while falling — no cut | 0.5 | true | +80 | +80 | false |
| cut applied only once | 0.5 | true (×2) | -320 | -160 | false |
| multiplier clamped below zero | -1.0 | true | -320 | -32 (0.1×) | false |

**Failing assertion (before implementation):**

```
// After simulating key-release while body.IsGoingUp() == true:
// got vy16 = -320, want -160
```

The test must use a stub `body` (local `mocks_test.go` in the skill package) and a stub input source injected via an interface — no `inpututil` calls in tests (system boundary). `JumpSkill` must accept an `InputSource` interface so the test can control key-press/release deterministically.

### `InputSource` interface (new, defined in the skill package)

```go
type InputSource interface {
    IsKeyJustPressed(key ebiten.Key) bool
    IsKeyJustReleased(key ebiten.Key) bool
}
```

`HandleInput` receives an `InputSource` instead of calling `inpututil` directly. The production default wraps `inpututil`. This keeps tests deterministic and non-flaky per the constitution.
