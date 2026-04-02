# SPEC 010 — Player Shooting Skill (Cuphead-style)

**Branch:** `010-player-shooting-skill`
**Bounded Context:** Game Logic (`internal/game/`)
**Story:** `USER_STORY.md`

---

## 1. Technical Requirements

### 1.1 Interface Change — `GroundedInput`

Add one method to the existing interface in `grounded_input.go`:

```go
type GroundedInput interface {
    HorizontalInput() int
    DuckHeld() bool
    HasCeilingClearance() bool
    JumpPressed() bool
    DashPressed() bool
    AimLockHeld() bool
    ShootHeld() bool   // NEW
}
```

All existing implementors (`MockInputSource` in `mocks_test.go`, and any real input adapter) must add `ShootHeld() bool`.

### 1.2 New Contract — `Shooter`

New file: `internal/engine/contracts/body/shooter.go`

```go
package body

// Shooter is the contract ShootingSkill depends on to spawn a bullet into the world.
type Shooter interface {
    // SpawnBullet adds a bullet Body to the BodiesSpace.
    // x16, y16 are fixed-point spawn position; speedX16 is signed fixed-point horizontal velocity.
    SpawnBullet(x16, y16, speedX16 int, owner interface{})
}
```

> Rationale: keeps `ShootingSkill` decoupled from any concrete bullet factory.

### 1.3 New Type — `ShootingSkill`

New file: `internal/game/entity/actors/states/shooting_skill.go`

```go
type ShootingConfig struct {
    CooldownFrames  int // frames between shots
    SpawnOffsetX16  int // horizontal offset from Actor center (fixed-point)
    BulletSpeedX16  int // bullet speed (fixed-point, always positive; direction applied at spawn)
    YOffset         int // half-amplitude for OffsetToggler
}

type ShootingSkill struct {
    cfg     ShootingConfig
    toggler *OffsetToggler
    shooter contractsbody.Shooter
    cooldown int // frames remaining until next shot is allowed
}

func NewShootingSkill(cfg ShootingConfig, shooter contractsbody.Shooter) *ShootingSkill

// Update is called every frame by GroundedState when ShootHeld() is true and the Actor
// is NOT in StateDashing. It decrements the cooldown and spawns a bullet when cooldown
// reaches zero.
//
// body is the Actor's Movable (provides FaceDirection and GetPosition16).
func (s *ShootingSkill) Update(body contractsbody.Movable)
```

Internal logic of `Update`:
1. If `s.cooldown > 0` → decrement and return (no spawn).
2. Compute spawn position:
   - `x16, y16 := body.GetPosition16()`
   - Apply `SpawnOffsetX16` in `body.FaceDirection()` direction.
   - Apply `s.toggler.Next()` to `y16`.
3. Compute `speedX16`: positive if facing right, negative if facing left.
4. Call `s.shooter.SpawnBullet(x16, y16, speedX16, body.Owner())`.
5. Reset `s.cooldown = s.cfg.CooldownFrames`.

### 1.4 `GroundedDeps` Change

```go
type GroundedDeps struct {
    Input    GroundedInput
    Shooting *ShootingSkill // NEW — nil disables shooting
    Body     contractsbody.Movable // NEW — needed by ShootingSkill.Update
}
```

### 1.5 `GroundedState.Update()` Change

After the dash/jump exit checks and before sub-state transition, add:

```go
if g.deps.Shooting != nil && input.ShootHeld() && g.activeKey != SubStateDashing {
    g.deps.Shooting.Update(g.deps.Body)
}
```

> `StateDashing` is a parent-level state, not a sub-state, so `g.activeKey` will never equal it. The suppression guard (AC7) is enforced by `GroundedState` not being active while `DashState` is active — no extra check needed inside `ShootingSkill` itself. The `g.activeKey != SubStateDashing` guard is therefore a no-op safety net; the real suppression is architectural.

### 1.6 New Type — `Bullet`

New file: `internal/game/entity/actors/states/bullet.go` (or `internal/game/entity/actors/bullets/bullet.go` — implementer's choice).

`Bullet` wraps a `contractsbody.Collidable` and implements `contractsbody.Touchable`:

```go
type Bullet struct {
    body    contractsbody.MovableCollidable
    space   contractsbody.BodiesSpace
    speedX16 int
}

// Update moves the bullet each frame and queues removal if out of bounds.
func (b *Bullet) Update()

// OnTouch queues the bullet for removal when it hits any Collidable that is not its owner.
func (b *Bullet) OnTouch(other contractsbody.Collidable)

// OnBlock satisfies Touchable; queues removal.
func (b *Bullet) OnBlock(other contractsbody.Collidable)
```

`Update` logic:
1. Apply `speedX16` via `body.SetVelocity(b.speedX16, 0)`.
2. Call `space.ResolveCollisions(b.body)`.
3. Check bounds via `space.GetTilemapDimensionsProvider()` — if bullet position is outside, call `space.QueueForRemoval(b.body)`.

---

## 2. Pre-conditions

- `OffsetToggler` exists in `offset_toggler.go` (US-008 ✅).
- `contractsbody.Movable` exposes `FaceDirection()`, `GetPosition16()`, and `Owner()`.
- `contractsbody.BodiesSpace` exposes `QueueForRemoval()`.
- `DashState` is a separate top-level state; `GroundedState` is inactive while dashing.

## 3. Post-conditions

- `GroundedInput` has `ShootHeld() bool`.
- `GroundedDeps` carries `Shooting *ShootingSkill` and `Body contractsbody.Movable`.
- `ShootingSkill.Update()` is called each frame `ShootHeld()` is true while grounded.
- At most one bullet spawns per `CooldownFrames` window regardless of input release/re-press.
- Bullet Y-offset alternates `+YOffset` / `−YOffset` across consecutive shots.
- Bullets are removed from `BodiesSpace` on out-of-bounds or collision with a non-owner `Collidable`.
- No global singletons; all dependencies injected.

---

## 4. Integration Points

| Point | Detail |
|---|---|
| `grounded_input.go` | Add `ShootHeld() bool` to `GroundedInput`; update `GroundedDeps` |
| `grounded_state.go` | Call `ShootingSkill.Update()` in `Update()` when `ShootHeld()` |
| `mocks_test.go` | Add `ShootHeldFunc func() bool` to `MockInputSource` |
| `internal/engine/contracts/body/shooter.go` | New `Shooter` interface |
| `shooting_skill.go` | New skill; depends on `Shooter` + `Movable` |
| `bullet.go` | New entity; implements `Touchable`, uses `BodiesSpace` |

---

## 5. Red Phase — Failing Test Scenarios

File: `internal/game/entity/actors/states/shooting_skill_test.go`

### Test 1 — Cooldown gating (no double-spawn within cooldown)

```
GIVEN ShootingSkill with CooldownFrames=3
WHEN Update is called 3 consecutive frames with ShootHeld=true
THEN SpawnBullet is called exactly once (on frame 1), not on frames 2 or 3
```

### Test 2 — Cooldown resets after window expires

```
GIVEN ShootingSkill with CooldownFrames=2
WHEN Update is called on frames 1, 2, 3 (all ShootHeld=true)
THEN SpawnBullet is called on frame 1 and frame 3 (total=2), not on frame 2
```

### Test 3 — Alternating Y-offset over ≥4 shots

```
GIVEN ShootingSkill with CooldownFrames=0, YOffset=4
WHEN Update is called 4 times (cooldown=0 so each call spawns)
THEN the y16 argument to SpawnBullet alternates: +4, -4, +4, -4
```

### Test 4 — Release and re-press within cooldown window does not spawn extra bullet

```
GIVEN ShootingSkill with CooldownFrames=3
WHEN frame 1: Update called (spawns), frame 2: not called, frame 3: Update called
THEN SpawnBullet called only once (cooldown not yet expired on frame 3)
```

### Test 5 — Suppression while dashing (AC7)

```
GIVEN GroundedState with ShootingSkill injected
AND  the Actor transitions to StateDashing
WHEN DashState.Update() runs for N frames
THEN ShootingSkill.Update() is never called (GroundedState is inactive)
```

> Test 5 is an integration-level test on `GroundedState`; tests 1–4 are pure unit tests on `ShootingSkill`.

---

## 6. Files to Create / Modify

| Action | Path |
|---|---|
| Modify | `internal/game/entity/actors/states/grounded_input.go` |
| Modify | `internal/game/entity/actors/states/grounded_state.go` |
| Modify | `internal/game/entity/actors/states/mocks_test.go` |
| Create | `internal/engine/contracts/body/shooter.go` |
| Create | `internal/game/entity/actors/states/shooting_skill.go` |
| Create | `internal/game/entity/actors/states/shooting_skill_test.go` |
| Create | `internal/game/entity/actors/states/bullet.go` (or `bullets/` sub-package) |
