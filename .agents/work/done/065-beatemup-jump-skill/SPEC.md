# SPEC — 065-beatemup-jump-skill

## 1. Scope

Add `BeatEmUpJumpSkill` driving altitude-axis jumps on `*physicsmovement.BeatEmUpMovementModel`. Mirrors `JumpSkill` (platform) UX: coyote, buffer, jump-cut. Operates on `VAltitude16` (negative = up).

## 2. Files

- New: `internal/kit/skills/beatemup_jump.go`
- New: `internal/kit/skills/beatemup_jump_test.go`
- Edit: `internal/kit/skills/factory.go` — branch on `cfg.Movement.Mode`.
- Edit: `internal/engine/physics/movement/movement_model_beatemup.go` — add `IsInputBlocked()` method (pre-req for AC-8).

Package: `kitskills` (same as `JumpSkill`).

## 3. Engine Pre-Req [AC-8]

Add to `movement_model_beatemup.go`:
```go
func (m *BeatEmUpMovementModel) IsInputBlocked() bool {
    return m.playerMovementBlocker != nil && m.playerMovementBlocker.IsMovementBlocked()
}
```
`BeatEmUpMovementModel` already has `playerMovementBlocker` field. Satisfies `movement.InputBlocker` interface.

## 4. Type Definition [AC-1, AC-7, AC-10]

```go
type BeatEmUpJumpSkill struct {
    skill.SkillBase
    activationKey     ebiten.Key  // ebiten.KeySpace
    coyoteTimeCounter int
    jumpBufferCounter int
    jumpCutMultiplier float64     // default 1.0; clamped (0.1, 1.0]
    jumpCutPending    bool
    jumpPressed       bool        // previous frame edge tracker
    OnJump            func(body body.MovableCollidable)
}

func NewBeatEmUpJumpSkill() *BeatEmUpJumpSkill
func (s *BeatEmUpJumpSkill) SetJumpCutMultiplier(m float64)  // same clamp as JumpSkill
func (s *BeatEmUpJumpSkill) ActivationKey() ebiten.Key
func (s *BeatEmUpJumpSkill) HandleInput(b body.MovableCollidable, model physicsmovement.MovementModel, space body.BodiesSpace)
func (s *BeatEmUpJumpSkill) Update(b body.MovableCollidable, model physicsmovement.MovementModel)
```

Constructor sets `activationKey=ebiten.KeySpace`, `jumpCutMultiplier=1.0`, `SetState(skill.StateReady)`. Implements `skill.ActiveSkill`.

## 5. HandleInput pseudocode [AC-1, AC-2, AC-3, AC-4, AC-6, AC-8]

```
HandleInput(b, model, space):
  bm, _ := model.(*BeatEmUpMovementModel)
  if bm == nil: return                       # AC-1 no-op other models
  if bm.IsInputBlocked(): return             # AC-8
  pressed := input.CommandsReader().Jump
  if pressed && !s.jumpPressed:              # leading edge
    s.tryActivate(b)                         # AC-2, AC-3, AC-4
  if !pressed && s.jumpPressed && s.jumpCutPending:  # trailing edge
    s.applyJumpCut(b)                        # AC-6
  s.jumpPressed = pressed
```

```
tryActivate(b):
  cfg := config.Get()
  grounded := b.Altitude() <= 0
  if grounded || s.coyoteTimeCounter > 0:
    force := int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())
    if force <= 0: return                    # edge: skip silently
    b.SetVAltitude16(-force)                 # AC-2 (altitude up = negative)
    s.jumpCutPending = true
    if s.OnJump != nil: s.OnJump(b)
    s.coyoteTimeCounter = 0
    s.jumpBufferCounter = 0
  else:
    s.jumpBufferCounter = cfg.Physics.JumpBufferFrames   # AC-4
```

```
applyJumpCut(b):
  v := b.VAltitude16()
  if v < 0:                                  # still going up
    b.SetVAltitude16(int(float64(v) * s.jumpCutMultiplier))
  s.jumpCutPending = false
```

## 6. Update pseudocode [AC-3, AC-4, AC-5, AC-6, AC-9]

```
Update(b, model):
  s.SkillBase.Update(b, model)               # no-op base
  bm, _ := model.(*BeatEmUpMovementModel)
  if bm == nil: return                       # AC-1
  if b.Freeze(): return                      # AC-9 — no counter/altitude mutation

  # Clear jumpCutPending once the apex passes (AC-6)
  if s.jumpCutPending && b.VAltitude16() >= 0:
    s.jumpCutPending = false

  grounded := b.Altitude() <= 0
  wasOnGround := grounded                    # snapshot before buffer-fire
  cfg := config.Get()

  # Coyote (AC-5)
  if grounded:
    s.coyoteTimeCounter = cfg.Physics.CoyoteTimeFrames
  else if s.coyoteTimeCounter > 0:
    s.coyoteTimeCounter--

  # Buffer decay (AC-4)
  if s.jumpBufferCounter > 0:
    s.jumpBufferCounter--

  # Buffered jump fires immediately on landing (AC-4)
  if grounded && s.jumpBufferCounter > 0:
    force := int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())
    if force <= 0: return
    b.SetVAltitude16(-force)
    s.jumpCutPending = true
    if s.OnJump != nil: s.OnJump(b)
    s.jumpBufferCounter = 0
    s.coyoteTimeCounter = 0
  _ = wasOnGround  # not needed since checking grounded directly
```

Note: unlike `JumpSkill` which uses `model.OnGround()`, `BeatEmUpJumpSkill` uses `b.Altitude() <= 0` (AC story-defined ground predicate).

## 7. Factory wiring [AC-11]

In `factory.go` Jump branch — replace direct `NewJumpSkill()` with mode-based selection:

```go
if cfg.Jump != nil && isEnabled(cfg.Jump.Enabled) {
    isBeatEmUp := cfg.Movement != nil && cfg.Movement.Mode == schemas.MovementModeEightDir
    if isBeatEmUp {
        js := NewBeatEmUpJumpSkill()
        if cfg.Jump.JumpCutMultiplier > 0 { js.SetJumpCutMultiplier(cfg.Jump.JumpCutMultiplier) }
        if deps.OnJump != nil { js.OnJump = func(b body.MovableCollidable) { deps.OnJump(b) } }
        skills = append(skills, js)
    } else {
        // existing JumpSkill path unchanged
    }
}
```

Constraint: existing `JumpSkill` test suite must still pass — no regressions on `MovementModeHorizontal` / default.

## 8. Pre-Conditions

- `BeatEmUpMovementModel.IsInputBlocked()` exists.
- `body.MovableCollidable` exposes `Altitude()`, `VAltitude16()`, `SetVAltitude16()`, `JumpForceMultiplier()`, `Freeze()` — all already present.
- `schemas.MovementModeEightDir == "eight_dir"` exists.

## 9. Post-Conditions (checkable)

- `NewBeatEmUpJumpSkill()` returns non-nil with `jumpCutMultiplier == 1.0`, `state == StateReady`.
- After grounded leading-edge press: `b.VAltitude16() == -force`, `jumpCutPending == true`, `OnJump` fired exactly once.
- After airborne leading-edge press (no coyote, no ground): `b.VAltitude16()` unchanged, `jumpBufferCounter == cfg.JumpBufferFrames`.
- After Update while `b.Freeze()`: all counters unchanged, `VAltitude16` unchanged.
- Calling `HandleInput` with `*PlatformMovementModel`: no field mutation, no `b.SetVAltitude16` call.
- `SetJumpCutMultiplier(m)`: `m<=0 → 0.1`; `m>1 → 1.0`; else stored verbatim.

## 10. Red-Phase Test Triples [AC-12]

Test file: `internal/kit/skills/beatemup_jump_test.go`. Use `bodyphysics.NewObstacleRect`, `config.Set` per existing pattern. Real `BeatEmUpMovementModel` with `nil` blocker (so `IsInputBlocked()==false`) unless test name implies otherwise. Use a stub `PlayerMovementBlocker` for AC-8.

```
T1: grounded jump fires [AC-2]
  pre:  b.Altitude=0, b.VAltitude16=0, jumpPressed=false; cmds.Jump=true
  act:  HandleInput
  post: b.VAltitude16 == -cfg.JumpForce; jumpCutPending==true; OnJump called=1

T2: no double-jump while airborne [AC-3]
  pre:  b.Altitude=100, b.VAltitude16=-50, coyote=0, jumpPressed=false; cmds.Jump=true
  act:  HandleInput
  post: b.VAltitude16 == -50; jumpBufferCounter == cfg.JumpBufferFrames

T3: coyote jump [AC-5]
  pre:  b.Altitude=20, coyoteTimeCounter=3, jumpPressed=false; cmds.Jump=true
  act:  HandleInput
  post: b.VAltitude16 == -cfg.JumpForce; coyote==0

T4: jump buffered → fires on landing [AC-4]
  pre:  jumpBufferCounter=5, b.Altitude=0 (just landed), Freeze=false
  act:  Update
  post: b.VAltitude16 == -cfg.JumpForce; jumpBufferCounter==0; OnJump called=1

T5: jump-cut applies multiplier [AC-6]
  pre:  jumpCutPending=true, jumpPressed=true, b.VAltitude16=-320, jumpCutMultiplier=0.5; cmds.Jump=false
  act:  HandleInput
  post: b.VAltitude16 == -160; jumpCutPending==false

T6: no-op when model is *PlatformMovementModel [AC-1]
  pre:  model = NewPlatformMovementModel(nil); cmds.Jump=true; b.Altitude=0
  act:  HandleInput
  post: b.VAltitude16 unchanged (==0); no panic

T7: input-blocked guard [AC-8]
  pre:  blocker.IsMovementBlocked=true; cmds.Jump=true; b.Altitude=0
  act:  HandleInput
  post: b.VAltitude16 == 0; jumpPressed unchanged (false) — no state advanced

T8: SetJumpCutMultiplier clamp [AC-7]
  cases: 0.5→0.5; 1.0→1.0; 0.0→0.1; -1→0.1; 1.5→1.0

T9: coyote decrement while airborne [AC-5]
  pre:  b.Altitude=10, coyoteTimeCounter=2, Freeze=false
  act:  Update
  post: coyoteTimeCounter==1

T10: coyote reset while grounded [AC-5]
  pre:  b.Altitude=0, coyoteTimeCounter=0
  act:  Update
  post: coyoteTimeCounter == cfg.CoyoteTimeFrames

T11: Freeze blocks Update mutation [AC-9]
  pre:  b.Freeze=true, b.Altitude=10, coyoteTimeCounter=2, jumpBufferCounter=2, b.VAltitude16=-100
  act:  Update
  post: coyoteTimeCounter==2; jumpBufferCounter==2; b.VAltitude16==-100

T12: force<=0 skips silently [edge case]
  pre:  b.JumpForceMultiplier=0, b.Altitude=0; cmds.Jump=true
  act:  HandleInput
  post: b.VAltitude16==0; jumpCutPending==false; OnJump not called

T13: factory selects BeatEmUpJumpSkill for eight_dir [AC-11]
  cfg.Movement.Mode = "eight_dir", cfg.Jump != nil
  expect: FromConfig returns slice containing *BeatEmUpJumpSkill (not *JumpSkill)

T14: factory selects JumpSkill for default/horizontal [AC-11 regression]
  cfg.Movement.Mode = "" or "horizontal"
  expect: FromConfig returns slice containing *JumpSkill (no BeatEmUp variant)
```

Use `input.CommandsReader()`: tests must set `Jump` field via the existing test helper (mirror what `JumpSkill` tests do — check `platform_jump_test.go` for the pattern, likely `input.SetCommandsForTest` or direct construction).

## 11. Mocks / Contracts Inventory

- No new contracts in `internal/engine/contracts/`. Skill uses existing `body.MovableCollidable` and `physicsmovement.MovementModel` types.
- No new mocks required. Tests use concrete `BeatEmUpMovementModel`, `PlatformMovementModel`, and `NewObstacleRect` body.
- Local stub `PlayerMovementBlocker` for AC-8 — define inline in test file.

## 12. Out of Scope

- Configurable `MaxAltitudeFallSpeed` cap (handled by 061/model).
- Variable jump height beyond cut multiplier.
- Air-jump count / double-jump.
- Jump SFX hookup (caller's `OnJump`).
