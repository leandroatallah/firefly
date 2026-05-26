# Technical Specification — 067-actor-json-hitbox-active-frames

**Branch:** `067-actor-json-hitbox-active-frames`
**Bounded Context:** `internal/kit/combat/`, `internal/engine/data/schemas/`

---

## 1. Schema Layer [AC-1, AC-2, AC-9]

**File:** `internal/engine/data/schemas/json.go`

Add new struct:
```go
// HitboxFrameRange is an optional per-state override for the melee hitbox
// active-frame window. Start/End are inclusive swing-frame indices.
type HitboxFrameRange struct {
    Start int `json:"start"`
    End   int `json:"end"`
}
```

Extend `AssetData`:
```go
type AssetData struct {
    Path           string             `json:"path"`
    CollisionRects []ShapeRect        `json:"collision_rect"`
    FootprintRect  *ShapeRect         `json:"footprint_rect,omitempty"`
    Loop           *bool              `json:"loop,omitempty"`
    HitboxFrames   *HitboxFrameRange  `json:"hitbox_frames,omitempty"` // NEW
}
```

Constraints:
- Package imports unchanged. Must not import `internal/kit/` or `internal/game/`.
- Absent JSON field → pointer nil (zero-decoded).

---

## 2. Weapon Layer [AC-3, AC-7, AC-8]

**File:** `internal/kit/combat/weapon/melee.go`

Add field to `MeleeWeapon`:
```go
activeFramesOverride *[2]int
```

Add method:
```go
// SetActiveFramesOverride installs (or clears, when override == nil) a per-swing
// override of the active-frame window used by IsHitboxActive. Must be called
// before Fire for the override to take effect on the upcoming swing.
func (w *MeleeWeapon) SetActiveFramesOverride(override *[2]int) {
    w.activeFramesOverride = override
}
```

Modify `IsHitboxActive`:
```
IsHitboxActive:
  if !w.swinging: return false
  if w.activeFramesOverride != nil:
    return swingFrame >= override[0] && swingFrame <= override[1]
  step := w.steps[w.stepIndex]
  return swingFrame >= step.ActiveFrames[0] && swingFrame <= step.ActiveFrames[1]
```

Modify `startSwing` (AC-8):
```
startSwing:
  w.activeFramesOverride = nil   // clear stale override before each swing
  w.swinging = true
  w.swingFrame = 0
  w.hitThisSwing = make(...)
  w.currentCooldown = w.cooldownFrames
```

Note: `Update()` already uses `step.ActiveFrames[1]` to terminate the swing.
Do NOT change termination — the override governs hitbox activation only, never
swing duration (AC-7). If `override[1] > step.ActiveFrames[1]`, the swing ends
first and `IsHitboxActive()` returns false (swinging=false).

---

## 3. Melee State Layer [AC-4, AC-5]

**File:** `internal/kit/combat/melee/state.go`

### 3.1 Extend `weaponIface`

```go
type weaponIface interface {
    combat.Weapon
    IsHitboxActive() bool
    IsSwinging() bool
    IsInStartup() bool
    ApplyHitbox(space contractsbody.BodiesSpace)
    StepIndex() int
    ComboWindowRemaining() int
    ResetCombo()
    SetActiveFramesOverride(override *[2]int) // NEW
}
```

### 3.2 Add asset map + resolver to `State`

```go
import "github.com/boilerplate/ebiten-template/internal/engine/data/schemas"

type State struct {
    // ... existing fields ...
    assets         map[string]schemas.AssetData
    stepStateName  func(stepIdx int) string
}
```

### 3.3 Setter (kept separate from constructor to avoid signature churn) [AC-5]

```go
// SetHitboxFrameResolver wires the per-actor asset map and a function that
// maps a combo-step index to the asset key (typically the registered name of
// the per-step ActorStateEnum). When both are set, OnStart will install the
// AssetData.HitboxFrames override on the weapon for that step; otherwise it
// clears the override.
func (s *State) SetHitboxFrameResolver(
    assets map[string]schemas.AssetData,
    stepStateName func(stepIdx int) string,
) {
    s.assets = assets
    s.stepStateName = stepStateName
}
```

Both inputs are owner-supplied at install time. State must never import `internal/game/`.

### 3.4 Modify `OnStart` [AC-4]

Insert override resolution BEFORE `weapon.Fire`:

```
OnStart(currentCount):
  startCount=currentCount; frame=0
  set returnTo from owner airborne flags
  if owner.IsDucking: frame=animFrames; return
  stepUsed = weapon.StepIndex()

  // NEW: per-step hitbox-frames override
  var override *[2]int = nil
  if assets != nil && stepStateName != nil:
    name := stepStateName(stepUsed)
    if asset, ok := assets[name]; ok && asset.HitboxFrames != nil:
      hf := asset.HitboxFrames
      override = &[2]int{hf.Start, hf.End}
  weapon.SetActiveFramesOverride(override)

  weapon.Fire(x16, y16, faceDir, ShootDirectionStraight, 0)
  spawn vfx (unchanged)
```

Pre/post:
- pre: `weapon.SetActiveFramesOverride` is called exactly once per `OnStart`, before `Fire`.
- post: when `HitboxFrames` is nil/missing, override arg is nil; when present, override arg equals `&[2]int{Start,End}`.
- The override pointer is freshly allocated per `OnStart` (no shared backing array between actors).

---

## 4. Kit Bridge Layer [AC-5]

**File:** `internal/kit/states/melee_state.go`

Extend `meleeWeaponIface` to mirror `weaponIface` (add `SetActiveFramesOverride(*[2]int)`). No other changes required because `MeleeAttackState = meleeengine.State` is a type alias — the new `SetHitboxFrameResolver` setter is inherited automatically.

No new constructor signature. Existing `InstallMeleeAttackState` / `NewMeleeAttackState` continue to work.

---

## 5. Game-Layer Wiring (informational, not in scope for AC-10) [AC-10]

Out-of-spec for this story but the wiring sites that will call the new setter are:
- `internal/game/entity/actors/player/cody.go::SetMelee` — after `InstallState`, call `st.SetHitboxFrameResolver(p.GetSpriteData().Assets, stepStateNameFn)`.
- `internal/game/entity/actors/player/climber.go::SetMelee` — same.

`stepStateNameFn` is constructed locally as:
```go
stepStateNameFn := func(stepIdx int) string {
    return stepStates[stepIdx].String()
}
```

This story does NOT modify game files; spec calls these out so the next story can wire actors that opt in. No existing JSON requires editing (AC-10).

---

## 6. Mock / Contract Inventory

- `internal/kit/combat/melee/state_test.go` — existing tests use inline `weaponIface` doubles; new tests need `SetActiveFramesOverride(*[2]int)` on the test double (record last value).
- `internal/kit/combat/weapon/mocks_test.go` — no change.
- No new entries in `internal/engine/mocks/`.
- No new contract files in `internal/engine/contracts/`.

---

## 7. Red Phase — Weapon Package [AC-11]

**File:** `internal/kit/combat/weapon/melee_active_frames_override_test.go`

Single table-driven test:

```
T-W1: TestMeleeWeapon_IsHitboxActive_RespectsOverride
  table rows (override, stepActiveFrames=[3,5], swingFrame → wantActive):
    "nil override, frame in step window"    nil, frame=3 → true
    "nil override, frame outside step window" nil, frame=6 → false
    "override [1,2] uses override not step"  &{1,2}, frame=2 → true; frame=3 → false
    "override exact Start"                    &{4,7}, frame=4 → true
    "override exact End"                      &{4,7}, frame=7 → true
    "override Start>End never activates"      &{6,5}, frame=5 → false; frame=6 → false
    "override Start==End single-frame"        &{4,4}, frame=3 → false; frame=4 → true; frame=5 → false

  pre: weapon constructed with ComboStep{ActiveFrames=[3,5], ...}
       SetActiveFramesOverride(override); Fire(...); advance Update() to target swingFrame
  act: w.IsHitboxActive()
  post: equals wantActive
```

```
T-W2: TestMeleeWeapon_StartSwing_ClearsOverride
  pre: SetActiveFramesOverride(&[2]int{1,2}); Fire(...) (overrides applied to swing 1)
       After swing 1 ends and cooldown elapses, call Fire(...) again without
       calling SetActiveFramesOverride.
  act: advance to step.ActiveFrames[0] of swing 2; check IsHitboxActive()
  post: behaves as nil override (i.e., active at frame 3, inactive at frame 2).
        Equivalently: assert that after Fire, the public-observable activation
        window matches ComboStep.ActiveFrames, not the previously-installed override.
```

---

## 8. Red Phase — Melee Package [AC-12]

**File:** `internal/kit/combat/melee/state_override_test.go`

Add a recording test double for the weapon:

```
type recordingWeapon struct {
    // satisfies weaponIface
    lastOverride     *[2]int
    overrideCallCount int
    fireCalled       bool
    overrideSetBeforeFire bool
    // ... other no-op methods returning zero values ...
}
func (r *recordingWeapon) SetActiveFramesOverride(o *[2]int) {
    r.lastOverride = o
    r.overrideCallCount++
}
func (r *recordingWeapon) Fire(...) {
    if r.overrideCallCount > 0 { r.overrideSetBeforeFire = true }
    r.fireCalled = true
}
// IsSwinging/IsHitboxActive/IsInStartup/StepIndex/etc. return defaults.
```

Tests:

```
T-M1: TestState_OnStart_InstallsOverrideFromAssetData
  table rows (assetHitboxFrames, wantOverride):
    "present"  &HitboxFrameRange{2,4}, wantOverride=&[2]int{2,4}
    "absent"   nil,                     wantOverride=nil

  pre: assets := map[string]schemas.AssetData{
         "melee_attack_step_0": {HitboxFrames: row.assetHitboxFrames},
       }
       stepStateName := func(i int) string { return "melee_attack_step_0" }
       state := NewState(owner, space, recWeapon, nil, meleeAttackEnum, Idle, Falling)
       state.SetAnimationFrames(10)
       state.SetHitboxFrameResolver(assets, stepStateName)
       recWeapon.StepIndex returns 0
  act: state.OnStart(0)
  post:
    recWeapon.overrideCallCount == 1
    recWeapon.overrideSetBeforeFire == true
    (lastOverride == nil) == (wantOverride == nil)
    if both non-nil: *lastOverride == *wantOverride
```

```
T-M2: TestState_OnStart_ClearsOverrideWhenResolverMissing
  pre: SetHitboxFrameResolver never called (assets/stepStateName both nil)
  act: state.OnStart(0)
  post: recWeapon.overrideCallCount == 1 AND lastOverride == nil
        (override always set explicitly per OnStart to nil to clear stale state)
```

```
T-M3: TestState_OnStart_PerStepIndependentOverride
  pre: two-step combo. assets:
         "melee_attack_step_0": HitboxFrames={1,2}
         "melee_attack_step_1": no HitboxFrames
       stepStateName(0)="melee_attack_step_0"; stepStateName(1)="melee_attack_step_1"
       state.SetHitboxFrameResolver(assets, stepStateName)
  act: recWeapon.StepIndex=0; state.OnStart(0); record override1
       recWeapon.StepIndex=1; state.OnStart(0); record override2
  post: override1 == &[2]int{1,2}; override2 == nil
```

```
T-M4: TestState_OnStart_Ducking_DoesNotTouchOverride
  pre: owner.ducking = true; resolver installed with present HitboxFrames
  act: state.OnStart(0)
  post: recWeapon.overrideCallCount == 0; recWeapon.fireCalled == false
        (early-return path is unchanged)
```

---

## 9. Layer / Import Constraints Checklist [AC-9]

- `internal/engine/data/schemas/` imports unchanged (only `engine/contracts/animation`). MUST NOT import `internal/kit/` or `internal/game/`.
- `internal/kit/combat/weapon/melee.go` imports unchanged.
- `internal/kit/combat/melee/state.go` adds import `internal/engine/data/schemas` (allowed: kit may import engine).
- `internal/kit/combat/melee/` MUST NOT import `internal/game/`.
- `internal/kit/states/melee_state.go` adds only interface method to its local `meleeWeaponIface`.

---

## 10. Behavioral Invariants

- `SetActiveFramesOverride(nil)` is idempotent.
- Override governs hitbox activation only; `ComboStep.ActiveFrames[1]` continues to define swing termination in `Update()`.
- `OnStart` always calls `SetActiveFramesOverride` exactly once (with nil or a freshly-allocated pointer) before `Fire`, except in the early-return ducking branch.
- `startSwing` clears the override at the start of every swing, including step-to-step transitions inside a combo and shared-weapon hand-offs between actors.
