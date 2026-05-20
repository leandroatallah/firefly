# SPEC — 061-altitude-jump-ground-detection

**Bounded Context:** Physics (`internal/engine/physics/movement/`)
**Target File:** `internal/engine/physics/movement/movement_model_beatemup.go`
**Test File:** `internal/engine/physics/movement/movement_models_test.go`

---

## 1. Scope

Activate altitude-axis gravity + landing inside `BeatEmUpMovementModel.Update`. No new public types, no new contracts. Jump input is external (sets `VAltitude16` to a negative value). The model is passive — it only integrates.

No new mocks. No changes to `internal/engine/contracts/`.

---

## 2. Altitude Integration Block [AC-1..AC-7]

Inserted into `BeatEmUpMovementModel.Update` AFTER the existing `b.SetVelocity(vx16, vy16)` final line. Freeze guard at top of `Update` already covers AC-6.

### Signature unchanged
```go
func (m *BeatEmUpMovementModel) Update(b body.MovableCollidable, space body.BodiesSpace) error
```

### Pseudocode (placed at end of Update, after final SetVelocity)

```
// Altitude axis (story 061)
vAlt16 := b.VAltitude16()
alt    := b.Altitude()           // pixels (= fp16.From16(altitude16))

grounded := alt <= 0 && vAlt16 >= 0
if !grounded:
    cfg := config.Get()
    if vAlt16 < 0:               // rising
        vAlt16 += cfg.Physics.UpwardGravity
    else:                        // falling (vAlt16 >= 0 AND airborne)
        vAlt16 += cfg.Physics.DownwardGravity

    // Integrate altitude (pixel-level via fp16.From16)
    alt += fp16.From16(vAlt16)

    if alt <= 0:                 // landing
        alt = 0
        vAlt16 = 0

    b.SetAltitude(alt)
    b.SetVAltitude16(vAlt16)
// else: resting on ground — no mutation (AC-2, AC-7)
```

### Post-conditions (per frame, when not frozen)
- `b.Freeze() == true`                    → no altitude/velocity mutation on altitude axis (AC-6).
- pre `vAlt16 < 0`                        → post `vAlt16' == vAlt16 + UpwardGravity` (then integrated) (AC-1).
- pre `vAlt16 >= 0` AND airborne          → post `vAlt16' == vAlt16 + DownwardGravity` (then integrated/clamped) (AC-1).
- pre `Altitude() <= 0 && vAlt16 >= 0`    → post `Altitude() == pre.Altitude()` AND `VAltitude16() == pre.VAltitude16()` (AC-2, AC-7).
- pre airborne, post `alt <= 0`           → `Altitude() == 0` AND `VAltitude16() == 0` (AC-4).
- jump impulse                             → set externally via `b.SetVAltitude16(-N)` before next Update (AC-5).
- 2D body never touching altitude         → `Altitude() == 0`, `VAltitude16() == 0` invariants hold (AC-2 → no-op branch) (AC-7).

### Imports already present
`config`, `fp16` (no new imports needed).

---

## 3. Red Phase Test Inventory [AC-8]

Append to `internal/engine/physics/movement/movement_models_test.go`. Each test sets `config.Set(&config.AppConfig{ScreenWidth:320, ScreenHeight:240, Physics: config.PhysicsConfig{UpwardGravity:2, DownwardGravity:4, SpeedMultiplier:1.0}})`.

Helper used by all tests: `newMockMovableCollidable()` (already exists). Use `sp := space.NewSpace()` and `model := NewBeatEmUpMovementModel(nil)`.

### T-061-1: rising arc — UpwardGravity accumulates [AC-1]
```
pre:  actor at (100,100), VAltitude16=-fp16.To16(10), Altitude=20, Freeze=false
act:  model.Update(actor, sp)
post: VAltitude16() == -fp16.To16(10) + UpwardGravity + integrate_carry
      (assertion: post VAltitude16 > pre VAltitude16 by exactly UpwardGravity == 2)
      Altitude() == 20 + fp16.From16(-fp16.To16(10)+2)  // = 20 + From16(-10*65536+2) ≈ 10
      (assertion: Altitude() < 20 AND Altitude() > 0)
```

### T-061-2: falling — DownwardGravity accumulates [AC-1]
```
pre:  actor at (100,100), VAltitude16=fp16.To16(2), Altitude=50
act:  model.Update(actor, sp)
post: VAltitude16() == fp16.To16(2) + DownwardGravity == fp16.To16(2)+4
      Altitude() < 50 AND Altitude() > 0
```

### T-061-3: landing clamps altitude and zeroes velocity [AC-4]
```
pre:  actor at (100,100), VAltitude16=fp16.To16(50), Altitude=1   // about to land
act:  model.Update(actor, sp)
post: Altitude() == 0
      VAltitude16() == 0
```

### T-061-4: idempotent grounded — no mutation [AC-2, AC-7]
```
pre:  actor at (100,100), VAltitude16=0, Altitude=0
act:  for i in 0..5: model.Update(actor, sp)
post: Altitude() == 0
      VAltitude16() == 0
```

### T-061-5: freeze guard skips altitude mutation [AC-6]
```
pre:  actor, VAltitude16=fp16.To16(5), Altitude=30, Freeze=true
act:  model.Update(actor, sp)
post: VAltitude16() == fp16.To16(5)   (unchanged)
      Altitude() == 30                 (unchanged)
```

### T-061-6: 2D regression — body never touching altitude stays at 0 [AC-7]
```
pre:  actor at (100,100), VAltitude16=0, Altitude=0, vx16=fp16.To16(2)
act:  for i in 0..30: model.Update(actor, sp)
post: Altitude() == 0
      VAltitude16() == 0
      (X position moved as before — no regression on 2D plane)
```

### T-061-7: external jump impulse → full rise/fall/land arc [AC-1, AC-3, AC-4, AC-5]
```
pre:  actor at (100,100), Altitude=0, VAltitude16=0
      Step A: actor.SetVAltitude16(-fp16.To16(8))   // external jump
act:  loop model.Update until Altitude()==0 again AND frame > 1
post: At some intermediate frame N: Altitude() > 0     (rose)
      At some frame M > N: VAltitude16() >= 0          (peaked/falling)
      At final frame: Altitude() == 0 AND VAltitude16() == 0   (landed)
      Loop terminates within 600 frames                 (deterministic, no runaway)
```

---

## 4. Out of Scope

- No `MaxFallSpeed`-equivalent cap on altitude axis (story leaves clean integration point; do not introduce a field).
- No jump input handling (external systems set `VAltitude16` directly).
- No state machine wiring (e.g., `Jumping` / `Falling` states) — that is a downstream story.
- No interaction with `body.AccelerationAltitude()` integration — reserved for future.

---

## 5. AC → Test Map

| AC  | Test(s) |
|-----|---------|
| AC-1 | T-061-1, T-061-2, T-061-7 |
| AC-2 | T-061-4 |
| AC-3 | T-061-1, T-061-2, T-061-7 |
| AC-4 | T-061-3, T-061-7 |
| AC-5 | T-061-7 |
| AC-6 | T-061-5 |
| AC-7 | T-061-4, T-061-6 |
| AC-8 | all T-061-* (table-driven where multiple scenarios share shape) |
