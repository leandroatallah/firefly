# SPEC — 057-beatemup-movement-model

**Branch:** `057-beatemup-movement-model`
**Bounded Context:** Physics (`internal/engine/physics/movement/`)
**Contract:** `MovementModel` (existing, in `movement_model.go`)

---

## 1. File Layout

| File | Status | Purpose |
|---|---|---|
| `internal/engine/physics/movement/movement_model_beatemup.go` | NEW | `BeatEmUpMovementModel` implementation |
| `internal/engine/physics/movement/movement_model.go` | EDIT | Add `BeatEmUp` enum + factory case |
| `internal/engine/physics/movement/movement_models_test.go` | EDIT (TDD) | Add `TestBeatEmUpMovementModel_*` tests |

No new contracts. No new mocks. Mock Generator stage may be **skipped**.

---

## 2. Public API [AC-1]

```go
package movement

type BeatEmUpMovementModel struct {
    playerMovementBlocker PlayerMovementBlocker // unused for AC scope; kept for symmetry with NewMovementModel factory signature
    isScripted            bool
}

func NewBeatEmUpMovementModel(playerMovementBlocker PlayerMovementBlocker) *BeatEmUpMovementModel

func (m *BeatEmUpMovementModel) Update(b body.MovableCollidable, space body.BodiesSpace) error
func (m *BeatEmUpMovementModel) SetIsScripted(isScripted bool)
```

Implements `MovementModel` (`Update`, `SetIsScripted`).

---

## 3. Enum & Factory [AC-1]

Edit `movement_model.go`:

```go
const (
    TopDown MovementModelEnum = iota
    Platform
    BeatEmUp           // NEW
)

// String() map gains:  BeatEmUp: "BeatEmUp"
```

Add factory case:

```go
case BeatEmUp:
    return NewBeatEmUpMovementModel(playerMovementBlocker), nil
```

---

## 4. Update() pseudocode [AC-2..AC-8]

```
Update(b, space):
    if b.Freeze(): return nil                                       // AC-8

    vx16, vy16 = b.Velocity()

    // Apply previous-frame velocity to position with collision resolution.
    _, _, _ = b.ApplyValidPosition(vx16, true,  space)              // AC-6 (X)
    _, _, _ = b.ApplyValidPosition(vy16, false, space)              // AC-6 (Y)
    vx16, vy16 = b.Velocity()

    clampToPlayArea(b, space)                                       // AC-7 (return value ignored)

    // Integrate acceleration set externally by the skill (passive model).
    accX, accY := b.Acceleration()
    scaledAccX, scaledAccY := smoothDiagonalMovement(accX, accY)    // AC-3
    vx16 = increaseVelocity(vx16, scaledAccX)
    vy16 = increaseVelocity(vy16, scaledAccY)

    // 2D speed cap (same math as TopDown).                          // AC-4
    speedMax16 := fp16.To16(b.MaxSpeed())
    if mult := config.Get().Physics.SpeedMultiplier; mult != 0 {
        speedMax16 = int(float64(speedMax16) * mult)
    }
    velSq := int64(vx16)*int64(vx16) + int64(vy16)*int64(vy16)      // NOTE: fixes TopDown bug — see NOTES.md
    maxSq := int64(speedMax16) * int64(speedMax16)
    if velSq > maxSq:
        scale = float64(speedMax16) / math.Sqrt(float64(velSq))
        vx16 = int(float64(vx16) * scale)
        vy16 = int(float64(vy16) * scale)

    b.CheckMovementDirectionX()
    b.SetAcceleration(0, 0)

    // Friction both axes.                                           // AC-5
    vx16 = reduceVelocity(vx16)
    vy16 = reduceVelocity(vy16)
    b.SetVelocity(vx16, vy16)
    return nil
```

**Post-conditions after Update (non-freeze, no input, idle body):**
- `accX, accY == 0, 0` (acceleration reset)
- `vy16` does not grow downward across frames (no gravity term anywhere). [AC-2]
- `vx16, vy16` monotonically approach 0 (friction). [AC-5]

---

## 5. Constraints & Non-Goals

- **No Y-axis gravity term.** No `UpwardGravity`, `DownwardGravity`, `MaxFallSpeed`, `handleGravity`, `onGround`, `CheckGround` are referenced.
- **No altitude write.** Update must not call any altitude setter. Altitude is owned by a future jump skill; this model is altitude-agnostic. [Constraint from USER_STORY]
- **Passive.** No `InputHandler` method. No `ebiten` import. The skill sets `Acceleration` before `Update`.
- **No new bounds args.** Constructor takes only `PlayerMovementBlocker` (kept solely for factory parity; field is currently unused by `Update`). Walkable strip enforced by Tiled obstacle tiles + `clampToPlayArea`.

---

## 6. Red Phase Test Scenarios [AC-9]

All tests live in `movement_models_test.go` and reuse `newMockMovableCollidable()` + `space.NewSpace()` already in the file. Tests are table-driven where multiple scenarios share a setup.

### T-BE1: Constructor [AC-1]
```
pre:  blocker := &mockPlayerMovementBlocker{}
act:  m := NewBeatEmUpMovementModel(blocker)
post: m != nil; m.playerMovementBlocker == blocker; m.isScripted == false
```

### T-BE2: SetIsScripted [AC-1]
```
pre:  m := NewBeatEmUpMovementModel(nil)
act:  m.SetIsScripted(true); then m.SetIsScripted(false)
post: m.isScripted == true after first call; == false after second
```

### T-BE3: Freeze guard [AC-8]
```
pre:  actor.SetPosition(100, 100); actor.SetVelocity(fp16.To16(3), fp16.To16(3)); actor.SetFreeze(true)
act:  m.Update(actor, sp)
post: err == nil; actor.Position().Min == (100, 100) (unchanged); velocity unchanged
```

### T-BE4: No Y gravity when idle [AC-2]
```
pre:  actor at (100,100); velocity (0,0); acceleration (0,0); no Freeze
act:  loop m.Update(actor, sp) for 60 frames
post: vy16 == 0 across every frame; actor.Position().Min.Y == 100
```

### T-BE5: Diagonal speed normalization [AC-3, AC-4]
Table-driven, two cases share setup (fresh actor per row):
```
maxSpeed = 5; SpeedMultiplier = 1.0

row "cardinal-X":
  pre:  acceleration = (fp16.To16(2), 0)
  act:  m.Update(actor, sp)  ×N frames until speed-cap reached (N=60)
  post: |vx16| ≈ fp16.To16(5) ± friction; vy16 == 0

row "diagonal":
  pre:  acceleration = (fp16.To16(2), fp16.To16(2)) each frame
  act:  m.Update(actor, sp)  ×60 frames (re-set accel each frame before Update)
  post: sqrt(vx16² + vy16²) ≤ fp16.To16(5) * 1.05  // within 5% of cap, not 1.41×
        and abs(|vx16| - |vy16|) small (symmetric)
```

### T-BE6: X obstacle collision respected [AC-6]
```
pre:  actor at (100,100) size 10x10, velocity (fp16.To16(20), 0)
      obstacle ObstacleRect at (120,100) size 10x10, IsObstructive=true, added to sp
act:  m.Update(actor, sp)
post: actor.Position().Min.X < 120   // did not pass through obstacle
      actor.Position().Min.X >= 100  // moved forward at least a little or was stopped
```

### T-BE7: Y obstacle collision respected [AC-6]
```
pre:  actor at (100,100) size 10x10, velocity (0, fp16.To16(20))
      obstacle at (100,120) size 10x10, IsObstructive=true
act:  m.Update(actor, sp)
post: actor.Position().Min.Y < 120
```

### T-BE8: Friction applied each frame [AC-5]
```
pre:  actor velocity (fp16.To16(4), fp16.To16(4)); acceleration (0,0)
act:  m.Update(actor, sp)
post: vx16After < fp16.To16(4); vy16After < fp16.To16(4)   // friction reduced both
      AND vx16After > 0; vy16After > 0                     // not zero in one step
```

### T-BE9: Acceleration reset after Update
```
pre:  acceleration = (fp16.To16(2), fp16.To16(2))
act:  m.Update(actor, sp)
post: actor.Acceleration() == (0, 0)
```

### T-BE10: clampToPlayArea engaged [AC-7]
```
pre:  config.ScreenWidth=320, ScreenHeight=240
      actor at (-10, -10) size 10x10
      sp = space.NewSpace() (no tilemap provider → uses screen size)
act:  m.Update(actor, sp)
post: actor.Position().Min.X == 0; actor.Position().Min.Y == 0
```

### T-BE11: Factory wiring [AC-1]
```
pre:  blocker := &mockPlayerMovementBlocker{}
act:  model, err := NewMovementModel(BeatEmUp, blocker)
post: err == nil; model != nil; type-asserts to *BeatEmUpMovementModel
```

### T-BE12: Enum String [AC-1]
```
act:  BeatEmUp.String()
post: == "BeatEmUp"
```

---

## 7. Acceptance Cross-Check

Each AC has at least one test:

- AC-1 → T-BE1, T-BE11, T-BE12
- AC-2 → T-BE4
- AC-3 → T-BE5
- AC-4 → T-BE5
- AC-5 → T-BE8
- AC-6 → T-BE6, T-BE7
- AC-7 → T-BE10
- AC-8 → T-BE3
- AC-9 → all of the above (table-driven where applicable)

---

## 8. Mock / Helper Inventory

- `mockPlayerMovementBlocker` — already in `movement_models_test.go`. Reuse.
- `newMockMovableCollidable()` — already in `movement_models_test.go`. Reuse.
- `dimsProvider` — already in `movement_test.go`. Reuse if tilemap-provider variant needed (not required for current ACs).
- Obstacle bodies for T-BE6/T-BE7: construct inline via `bodyphysics.NewObstacleRect(...)`, `SetIsObstructive(true)`, `AddCollisionBodies()`, `sp.AddBody(...)` — pattern lifted from `TestPlatformMovementModel_CheckGround`.

---

## 9. Out-of-Scope (deferred)

- Wiring into a beat-em-up scene → story 058.
- `EightDirectionalMovementSkill` that sets acceleration → story 056.
- Altitude axis / jump → future story; this spec only guarantees the model does not touch altitude.
