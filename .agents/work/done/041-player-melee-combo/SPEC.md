# SPEC — US-041 Player Melee Combo Chain

**Branch:** `041-player-melee-combo`

**Bounded Contexts:**
- Engine: `internal/engine/combat/weapon/` (extend `MeleeWeapon` + factory)
- Game: `internal/game/entity/actors/states/` (extend `MeleeAttackState` + `GroundedState` wiring)
- Assets: `assets/entities/player/climber.json` (melee config block; currently held in code at `internal/game/entity/actors/player/weapons.go`)

**Dependencies / Referenced Contracts:**
- `internal/engine/contracts/combat/weapon.go` — `Weapon`, `Faction`, `Factioned` (US-022, US-038).
- `internal/engine/contracts/combat/damageable.go` — `Damageable` (US-038).
- `internal/engine/contracts/body/body.go` — `BodiesSpace.Query(rect)`, `Collidable` (existing).
- `internal/engine/entity/actors/actor_state.go` — `ActorState`, `ActorStateEnum`, `RegisterState` (existing).
- `internal/engine/input/commands.go` — `PlayerCommands.Melee`, `PlayerCommands.Jump`, `PlayerCommands.Dash` (existing).
- Prior spec: `.agents/work/done/040-player-melee-attack/SPEC.md` — baseline single-hit melee.

---

## 1. Technical Requirements

### 1.1 `ComboStep` value type (new — `internal/engine/combat/weapon/melee.go`)

```go
// ComboStep is one step in a melee combo chain. All length/offset values are
// in fp16 (the factory converts from pixels via fp16.To16).
type ComboStep struct {
    Damage          int
    ActiveFrames    [2]int // [first, last] inclusive
    HitboxW16       int
    HitboxH16       int
    HitboxOffsetX16 int
    HitboxOffsetY16 int
}
```

Rationale: combo parameters are per-step and immutable after load. Exporting the struct (vs keeping it package-private) lets `gameplayer.NewPlayerMeleeWeapon` build a weapon in code for now, mirroring US-040's current approach (JSON wiring is partial — see §1.4).

### 1.2 `MeleeWeapon` changes (`internal/engine/combat/weapon/melee.go`)

`MeleeWeapon` is refactored so the "active step" drives hitbox/damage each swing.

```go
type MeleeWeapon struct {
    id                string
    cooldownFrames    int
    currentCooldown   int
    comboWindowFrames int         // AC1 — frames after last-hit to accept next input
    steps             []ComboStep // AC2 — len in [1, 3]

    owner interface{}

    // Combo runtime state
    stepIndex         int         // 0..len(steps)-1, current step being swung
    windowRemaining   int         // frames remaining in combo window; 0 when inactive
    swinging          bool
    swingFrame        int
    hitThisSwing      map[combat.Damageable]struct{}

    faceDir              animation.FacingDirectionEnum
    originX16, originY16 int
}
```

**Constructor (breaking change):**

```go
// NewMeleeWeapon constructs a combo-capable melee weapon.
// steps must be non-empty; comboWindowFrames >= 0.
func NewMeleeWeapon(id string, cooldownFrames, comboWindowFrames int, steps []ComboStep) *MeleeWeapon
```

The pre-US-041 constructor (US-040 signature) is removed. All callers (`gameplayer.NewPlayerMeleeWeapon`, test helpers in `melee_test.go` and `melee_state_test.go`) must be updated.

**New public methods:**

| Method | Semantics |
|---|---|
| `StepIndex() int` | current step index (0-based); always the step about to be — or just — swung. |
| `ComboWindowRemaining() int` | frames remaining in post-hit input window; 0 when outside window. |
| `ResetCombo()` | sets `stepIndex = 0`, `windowRemaining = 0`. Does **not** clear cooldown or interrupt an in-flight swing. |
| `Steps() []ComboStep` | getter for tests (returns a copy or read-only slice). |

**Preserved `combat.Weapon` surface:** `ID`, `Fire`, `CanFire`, `Update`, `Cooldown`, `SetCooldown`, `SetOwner`. Internals change; signatures are unchanged.

**Semantics:**

- `Fire(x16, y16, faceDir, _, _)`:
  1. If `!CanFire()` — return (cooldown guard, unchanged).
  2. If `windowRemaining == 0` and `stepIndex != 0`, reset `stepIndex = 0` (window expired but state not yet cleared; defensive).
  3. Record origin/face, `swinging = true`, `swingFrame = 0`, clear `hitThisSwing`, set `currentCooldown = cooldownFrames`, set `windowRemaining = 0` (disable window during an active swing).
- `Update()`:
  1. Decrement `currentCooldown` if > 0.
  2. If `swinging`: increment `swingFrame`. When `swingFrame > steps[stepIndex].ActiveFrames[1]`:
     - Set `swinging = false`.
     - If `stepIndex < len(steps)-1`: set `windowRemaining = comboWindowFrames` (open the window — AC1).
     - Else (`stepIndex == len(steps)-1`): call `ResetCombo()` (AC4 — last step always wraps).
  3. Else if `windowRemaining > 0`:
     - Decrement. When it reaches 0, call `ResetCombo()` (AC3 — window-miss reset).
- `IsHitboxActive()` — returns `swinging && steps[stepIndex].ActiveFrames[0] <= swingFrame <= steps[stepIndex].ActiveFrames[1]`.
- `ApplyHitbox(space)` — unchanged in shape, but reads damage/hitbox dims from `steps[stepIndex]` instead of weapon-level fields. Faction gate and single-hit-per-swing guard are preserved from US-040.

**Combo advance happens at `Fire()` time.** When the caller (state or climber update) detects that the player pressed melee **while `windowRemaining > 0`**, it must call `AdvanceCombo()` BEFORE `Fire()`:

```go
// AdvanceCombo bumps stepIndex if the combo window is open and the next step
// exists. Returns true if the index advanced. Does not call Fire.
func (w *MeleeWeapon) AdvanceCombo() bool
```

Invariants:
- If `windowRemaining == 0`, `AdvanceCombo()` returns `false` and does nothing. The next `Fire()` swings step 0.
- If `stepIndex == len(steps)-1`, `AdvanceCombo()` returns `false` (cannot advance past last step; `Update()` already reset on end-of-last-step).
- Otherwise, `stepIndex++`, `windowRemaining = 0`, return `true`.

**Why split `AdvanceCombo()` from `Fire()`:** the caller needs to know whether the press continued the combo or started a fresh chain, to disambiguate animation selection (AC5) and to respect the cooldown (AC1 implicitly requires that combo presses DURING cooldown — i.e., before the active window window closes — still land; in practice, the window only opens after swing recovery, so by design `CanFire` becomes true at (or before) window-open).

### 1.3 Interrupt / reset API (game-side reset triggers — AC3)

Reset is a one-liner: `weapon.ResetCombo()`. It is invoked from three external triggers:

1. **Window expiry** — internal to `MeleeWeapon.Update()` (§1.2); no caller action needed.
2. **Damage taken** — `ClimberPlayer.Hurt(damage)` must call `p.melee.ResetCombo()` before the dying/hurt transition (AC3 bullet 2).
3. **Dash / jump mid-combo** — when `GroundedState.Update()` returns `StateDashing` or `actors.Falling` as a result of `DashPressed()` / `JumpPressed()`, the player's update loop must call `p.melee.ResetCombo()`. Concretely, this is handled in `ClimberPlayer.Update()` after reading `cmds`: if `(cmds.Dash || cmds.Jump) && p.melee.ComboWindowRemaining() > 0 { p.melee.ResetCombo() }`.

Rationale for wiring reset in `ClimberPlayer` rather than in state transitions: the weapon lives on the player, not the state, and the combo window can elapse while the actor is in `StateGrounded` (between swings). Centralising the reset avoids sprinkling knowledge of `MeleeWeapon` across multiple states.

### 1.4 `MeleeAttackState` changes (`internal/game/entity/actors/states/melee_state.go`)

The state now selects the animation based on `stepIndex` at `OnStart` time.

```go
type MeleeAttackState struct {
    owner    meleeOwnerIface
    space    contractsbody.BodiesSpace
    weapon   meleeWeaponIface
    returnTo actors.ActorStateEnum
    animFrames int
    frame      int
    stepUsed   int // snapshot of weapon.StepIndex() at OnStart (for tests / animation hook)
}
```

Changes:
- `meleeWeaponIface` gains `StepIndex() int` and `ComboWindowRemaining() int` accessors (so tests and the animation hook don't cast to the concrete type).
- `OnStart(currentCount)` captures `s.stepUsed = s.weapon.StepIndex()` (taken **after** the climber's update has called `AdvanceCombo`+`Fire`, so it reflects the step actually being swung).
- `SetAnimationFrames(n)` is unchanged; the climber's animation controller is responsible for picking `MeleeAttack1` / `MeleeAttack2` / `MeleeAttack3` from `stepUsed` (AC5).
- `Update()` — unchanged logic; still returns `returnTo` when `frame >= animFrames`. It does NOT call `Fire()` anymore — the climber's update owns `Fire()` so that combo advance can happen before the state is entered.

**Important:** `OnStart` no longer calls `weapon.Fire(...)`. US-040's `OnStart` started the swing directly; US-041 moves this responsibility to `ClimberPlayer.Update()` so combo advance (§1.2) and first-press detection share a single code path. The state's job narrows to "play the animation and apply the hitbox each frame it's active". The melee_state test `TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow` will be updated accordingly — it must call `w.Fire(...)` explicitly before driving the state.

`TryMeleeFromFalling(w, meleePressed)` — unchanged signature; still returns `(StateMeleeAttack, true)` iff `meleePressed && w.CanFire()`. Air melee does NOT combo (out of scope — the window requires ground contact in current design; `GroundedState` is the sole combo-advance site). The helper remains for air-melee trigger parity with US-040.

### 1.5 `ClimberPlayer.Update()` changes (`internal/game/entity/actors/player/climber.go`)

Replace the single-swing logic with:

```go
func (p *ClimberPlayer) Update(space body.BodiesSpace) error {
    cmds := input.CommandsReader()

    if p.melee != nil {
        p.melee.Update()

        // Reset combo when the player commits to dash or jump mid-window.
        if (cmds.Dash || cmds.Jump) && p.melee.ComboWindowRemaining() > 0 {
            p.melee.ResetCombo()
        }

        meleePressed := cmds.Melee && !p.meleeHeldPrev
        if meleePressed && p.melee.CanFire() && !p.IsDucking() {
            // If window is open, advance combo BEFORE firing so the next swing
            // uses steps[stepIndex+1] (AC1).
            if p.melee.ComboWindowRemaining() > 0 {
                p.melee.AdvanceCombo()
            }
            x16, y16 := p.GetPosition16()
            p.melee.Fire(x16, y16, p.FaceDirection(), body.ShootDirectionStraight, 0)
            p.spawnMeleeVFX(x16, y16)
        }
        if p.melee.IsHitboxActive() {
            p.melee.ApplyHitbox(space)
        }
        p.meleeHeldPrev = cmds.Melee
    }

    // ...rest unchanged...
}
```

And in `Hurt`:

```go
func (p *ClimberPlayer) Hurt(damage int) {
    if p.melee != nil {
        p.melee.ResetCombo()
    }
    if p.State() == gamestates.Dying || p.State() == gamestates.Dead {
        return
    }
    p.SetNewStateFatal(gamestates.Dying)
}
```

### 1.6 Factory changes (`internal/engine/combat/weapon/factory.go`)

Extend `parseMeleeWeapon` to accept the combo schema (AC6). The legacy single-step schema is **removed** — the factory always expects a `combo_steps` array.

```go
func parseMeleeWeapon(data []byte) (*MeleeWeapon, error) {
    var config struct {
        ID                string `json:"id"`
        CooldownFrames    int    `json:"cooldown_frames"`
        ComboWindowFrames int    `json:"combo_window_frames"`
        ComboSteps        []struct {
            Damage       int    `json:"damage"`
            ActiveFrames [2]int `json:"active_frames"`
            Hitbox       *struct {
                Width   int `json:"width"`
                Height  int `json:"height"`
                OffsetX int `json:"offset_x"`
                OffsetY int `json:"offset_y"`
            } `json:"hitbox"`
        } `json:"combo_steps"`
    }
    // ...json.Unmarshal...
    // Validation:
    //   - len(ComboSteps) in [1, 3]                  → else "combo_steps must contain 1..3 entries"
    //   - ComboWindowFrames >= 0                     → else "invalid combo_window_frames"
    //   - CooldownFrames >= 0                        → existing
    //   - each step: Hitbox != nil                   → else "hitbox is required for melee combo step N"
    //   - each step: ActiveFrames[0] >= 0 && [1] >= [0] → else "invalid active_frames for combo step N"
    //   - each step: Hitbox.Width > 0 && Height > 0  → else "invalid hitbox dimensions for combo step N"
    // Convert each step's Hitbox via fp16.To16 and build []ComboStep.
    // Return NewMeleeWeapon(ID, CooldownFrames, ComboWindowFrames, steps).
}
```

**Legacy compatibility:** none. Any caller still using the old `{ "damage": N, "active_frames": [...], "hitbox": {...} }` shape at the top level will get `"combo_steps must contain 1..3 entries"` on load. Callers under game control are updated in §1.5/§1.7.

### 1.7 Player JSON config

The story example (AC6) shows the target schema. Two landing options:

- **Option A (code-only, minimal risk):** keep melee construction in `gameplayer.NewPlayerMeleeWeapon()` and build `[]ComboStep` in Go; defer JSON integration until a follow-up story. This mirrors US-040's final landing (melee is NOT in `climber.json` at present — it's built in Go).
- **Option B (full JSON, per story):** add a `melee` block to `assets/entities/player/climber.json` and wire a loader. Requires a new loader path because `climber.json` currently has no weapon entries (weapons are built imperatively in `NewClimberInventory`).

**Decision:** Option A for this story. The factory is still updated (§1.6) so the JSON contract is implemented and test-covered end-to-end, and any future JSON-driven player loader will work without further changes. Landing Option B requires a separate story for the climber data loader, orthogonal to combo mechanics.

`gameplayer.NewPlayerMeleeWeapon()` becomes:

```go
func NewPlayerMeleeWeapon() *weapon.MeleeWeapon {
    steps := []weapon.ComboStep{
        {Damage: 1, ActiveFrames: [2]int{4, 10}, HitboxW16: fp16.To16(24), HitboxH16: fp16.To16(16), HitboxOffsetX16: fp16.To16(12), HitboxOffsetY16: fp16.To16(0)},
        {Damage: 1, ActiveFrames: [2]int{3, 8},  HitboxW16: fp16.To16(28), HitboxH16: fp16.To16(16), HitboxOffsetX16: fp16.To16(14), HitboxOffsetY16: fp16.To16(-4)},
        {Damage: 2, ActiveFrames: [2]int{5, 12}, HitboxW16: fp16.To16(32), HitboxH16: fp16.To16(20), HitboxOffsetX16: fp16.To16(16), HitboxOffsetY16: fp16.To16(0)},
    }
    return weapon.NewMeleeWeapon("player_melee", 20, 15, steps)
}
```

### 1.8 State machine transitions

```
Grounded --[Melee pressed, CanFire, window==0]----------------> MeleeAttack(step=0)
Grounded --[Melee pressed, CanFire, window>0, step<last]------> MeleeAttack(step=stepIndex+1)
Grounded --[Melee pressed, CanFire, step==last (always ran)]--> (window already reset) step=0
MeleeAttack --[IsAnimationFinished]---------------------------> returnTo (Grounded or Falling)
Grounded/window>0 --[Dash or Jump pressed]--------------------> (weapon.ResetCombo) → Dashing/Falling
Grounded/window>0 --[window elapses]--------------------------> weapon.ResetCombo (internal)
Any state --[Hurt(damage)]------------------------------------> weapon.ResetCombo (in ClimberPlayer.Hurt)
```

### 1.9 Contracts introduced

No new files in `internal/engine/contracts/`. The combat `Weapon` contract is unchanged. All new API lives on the concrete `*MeleeWeapon` (exported methods: `AdvanceCombo`, `ResetCombo`, `StepIndex`, `ComboWindowRemaining`, `Steps`). The game-side `meleeWeaponIface` in `melee_state.go` is widened to include `StepIndex()` and `ComboWindowRemaining()`.

Because no new interfaces are introduced in `internal/engine/contracts/`, **no new mocks are required** in `internal/engine/mocks/`. The existing package-local `mocks_test.go` in `internal/game/entity/actors/states/` already holds the mocks that drive `GroundedState` tests and will simply gain a couple of fields (see §4.RED-3).

---

## 2. Integration Points

- **Inventory** (`internal/engine/combat/inventory/`) — unchanged. `MeleeWeapon` still satisfies `combat.Weapon`; combo state is internal.
- **`ShootingSkill`** — unchanged. Melee continues to bypass `ShootingSkill.HandleInput`.
- **Animation controller** — `stepIndex` informs which `MeleeAttack1/2/3` clip to play (AC5). The animation hook reads `MeleeAttackState.stepUsed` (or `weapon.StepIndex()` directly). Actual sprite assets / animation registration are out of scope here (defer with art pass); tests assert the step index surface is correct so the animation layer has a deterministic source of truth.
- **`BodiesSpace.Query`** — unchanged contract; per-step hitbox dimensions simply alter the rect.
- **`input.PlayerCommands`** — unchanged; existing `Melee`, `Jump`, `Dash` edges are consumed.

---

## 3. Pre- and Post-Conditions (per AC)

| AC | Pre-condition | Post-condition |
|---|---|---|
| AC1 | Weapon with `comboWindowFrames=15` and 3 steps; step 1 swing just ended; `meleePressed=true` on frame `<=15` after swing-end | `AdvanceCombo()` returns `true`; next `Fire()` swings with step 2 parameters. `StepIndex()==1` during that swing. |
| AC2 | Weapon built from `combo_steps` JSON with distinct damage/active_frames/hitbox per step | For each step k, during its active window, `ApplyHitbox` uses `steps[k].Damage` and a rect derived from `steps[k].Hitbox*` and `steps[k].ActiveFrames`. |
| AC3a | Step 1 swing ended; `comboWindowFrames` frames elapsed with no melee press | `StepIndex()==0` and `ComboWindowRemaining()==0`. The next `Fire()` swings step 1. |
| AC3b | Mid-combo (`windowRemaining>0`, `stepIndex>=0`); `ClimberPlayer.Hurt(n)` is called | `StepIndex()==0`, `ComboWindowRemaining()==0` immediately after `Hurt` returns. |
| AC3c | Mid-combo window open; `cmds.Dash` (or `cmds.Jump`) true during `ClimberPlayer.Update` | After `Update`, `StepIndex()==0` and `ComboWindowRemaining()==0`. The state machine has transitioned to `StateDashing` / `Falling`. |
| AC4 | Step 3 swing ends naturally (no interrupt) | Internal `Update()` invokes `ResetCombo` on the tick that leaves the active window, so `StepIndex()==0` and `ComboWindowRemaining()==0`. No fourth swing is ever possible. |
| AC5 | Each step has distinct animation clip registration | `MeleeAttackState.stepUsed == weapon.StepIndex()` captured at `OnStart`; the climber's animation controller can read this to select `MeleeAttack1/2/3`. |
| AC6 | AC6 JSON passed to `weapon.NewWeaponFromJSON` | Returns a `*MeleeWeapon` with `len(Steps())==3`, matching fields (dims converted via `fp16.To16`), `ComboWindowFrames==15`, no error. Invalid variants return a non-nil error. |
| AC7 | All tests below (§4) | Pass after Feature Implementer finishes Green phase. |

---

## 4. Red Phase — Failing Test Scenarios (for TDD Specialist)

### RED-1 — `internal/engine/combat/weapon/melee_test.go` (extend)

Existing tests must be updated to build the weapon with a single-step slice:

```go
func newTestMeleeWeapon(owner interface{}) *weapon.MeleeWeapon {
    steps := []weapon.ComboStep{{
        Damage: 1, ActiveFrames: [2]int{3, 5},
        HitboxW16: 24*16, HitboxH16: 16*16, HitboxOffsetX16: 12*16, HitboxOffsetY16: 0,
    }}
    w := weapon.NewMeleeWeapon("player_melee", 20, 0, steps)
    w.SetOwner(owner)
    return w
}
```

New tests (add to the same file):

- `TestMeleeWeapon_Combo_AdvancesWhenPressedWithinWindow`:
  - Build a 3-step weapon (distinct damage values 1/1/2, `comboWindowFrames=15`, `cooldownFrames=0` to ignore cooldown between steps).
  - Fire step 1; `Update()` past its active window; assert `ComboWindowRemaining()>0` and `StepIndex()==0`.
  - Call `AdvanceCombo()` → returns true; `StepIndex()==1`.
  - Fire + run step 2 to completion; `AdvanceCombo()` → true; `StepIndex()==2`.
  - Fire + run step 3 to completion; `StepIndex()==0` and `ComboWindowRemaining()==0` (AC4 wrap).

- `TestMeleeWeapon_Combo_ResetsOnWindowExpiry`:
  - Fire step 1; `Update()` past active window; `Update()` `comboWindowFrames` more times.
  - Assert `StepIndex()==0`, `ComboWindowRemaining()==0`.
  - `AdvanceCombo()` returns false.

- `TestMeleeWeapon_Combo_ResetsOnDemand`:
  - Fire step 1; advance to step 2; call `ResetCombo()`; assert `StepIndex()==0`, `ComboWindowRemaining()==0`.

- `TestMeleeWeapon_Combo_PerStepDamageAndHitbox`:
  - Build a 3-step weapon with different `Damage` and `HitboxW16` per step; same enemy placed to overlap **all** three step hitboxes.
  - For each step k in 0..2: Fire, advance to first active frame, `ApplyHitbox`; assert `enemy.damageCalls[k] == steps[k].Damage`.
  - After step k's swing ends, call `AdvanceCombo()` to prepare for step k+1 (skip after last).

- `TestMeleeWeapon_Combo_LastStepAlwaysResets`:
  - 3-step weapon; run through all three steps; assert post-step-3 `StepIndex()==0` without any reset call from the test.

### RED-2 — `internal/engine/combat/weapon/factory_test.go` (extend)

- `"melee combo ok"` — the full AC6 JSON → expect `*MeleeWeapon` with `len(Steps())==3`, `Steps()[0].Damage==1`, `Steps()[2].Damage==2`, `Steps()[1].HitboxOffsetY16 == fp16.To16(-4)`, no error.
- `"melee combo missing combo_steps"` → error containing `combo_steps`.
- `"melee combo too many steps"` (4 entries) → error containing `combo_steps`.
- `"melee combo step missing hitbox"` → error containing `hitbox` and step index.
- `"melee combo negative window"` (`combo_window_frames: -1`) → error containing `combo_window_frames`.
- `"melee combo inverted active_frames for step 2"` → error containing `active_frames` and step index.
- Pre-existing `"melee weapon ok"` (single-step schema) must be updated or replaced with the combo schema; the legacy shape is no longer supported.

### RED-3 — `internal/game/entity/actors/states/melee_state_test.go` (extend)

- Update `newTestMeleeWeaponForState` to use the new `NewMeleeWeapon(id, cooldown, window, steps)` signature with a one-step slice. Existing US-040 tests continue to pass (single-step combo behaviour is equivalent to the legacy single-hit).
- `TestMeleeAttackState_UsesCurrentComboStep`:
  - Build a 3-step weapon. Fire step 1 directly on the weapon; construct `MeleeAttackState`; call `OnStart(0)`; assert `state.StepUsed() == 0` (expose via a test-only accessor, e.g. `func (s *MeleeAttackState) StepUsed() int`). Then `AdvanceCombo`, `Fire`, re-construct `MeleeAttackState` → `StepUsed() == 1`.
  - (This also covers AC5's dependency: the animation layer reads `StepUsed()`.)

- Update `TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow`:
  - Since `OnStart` no longer calls `Fire`, the test must call `w.Fire(...)` explicitly before `st.OnStart(0)`. Assert damage still lands exactly once.

- `TestGroundedState_DashPressed_ResetsComboWindow` (new, in `grounded_state_test.go` or adjacent):
  - Using a real `*MeleeWeapon`: simulate open window (`Fire`, advance past active frames). Then call the climber-level reset path directly: construct a tiny harness that invokes the reset logic from §1.5 (or expose the reset via a helper `gamestates.ResetComboOnInterrupt(w, dashOrJumpPressed)` if cleaner). Assert `w.StepIndex()==0` and `w.ComboWindowRemaining()==0`.
  - Rationale: the full reset is wired in `ClimberPlayer.Update`, which is hard to unit-test without a physics harness; a tiny exported helper keeps the logic testable without touching the climber in this test.

- `TestClimberPlayer_Hurt_ResetsCombo` — acceptable placement options:
  1. If a climber test file exists (`internal/game/entity/actors/player/climber_test.go`) → add here.
  2. Otherwise, verify indirectly: `MeleeWeapon.ResetCombo()` coverage plus a code-review check that `Hurt` calls it. The TDD Specialist MUST prefer Option 1 (create the test file if needed).

### RED-4 — JSON-driven construction smoke

- `internal/engine/combat/weapon/factory_test.go` — add a test that constructs the weapon from the exact JSON in AC6 and drives a single full combo chain via `Fire` → `AdvanceCombo` → `Fire` → `AdvanceCombo` → `Fire`, asserting damage values pulled from the JSON.

---

## 5. Non-Goals (out of scope)

- Aerial combos. `TryMeleeFromFalling` still triggers a swing in the air, but combo advance only works on the ground (combo window management is handled in `ClimberPlayer.Update` only). An air-started swing lands step 0; pressing again mid-air does nothing unless back on the ground.
- Knockback scaling per combo step.
- Directional melee (up-stab / down-stab).
- Player JSON-driven weapon loading (Option B in §1.7) — factory is ready; player loader is not.
- Animation sprite work (`MeleeAttack1/2/3` sheets). Step index is exposed; asset registration is deferred.

---

## 6. Key Design Decisions

1. **`AdvanceCombo` is explicit and caller-driven.** Splitting combo advance from `Fire` avoids hidden state changes during `Fire` and makes the "press during window" edge testable in isolation. The climber update loop is the sole place where the press edge is converted into a combo advance.
2. **Window management lives on the weapon, reset triggers live on the player.** Expiry (time-based) is internal to `MeleeWeapon.Update()`; interrupt resets (damage/dash/jump) are external, invoked by `ClimberPlayer`. This split keeps the weapon ignorant of the actor's state machine while letting external reset policy evolve without weapon changes.
3. **`NewMeleeWeapon` signature breaks.** The legacy (`damage, cooldown, active, hitboxW, hitboxH, offX, offY`) signature cannot express multi-step hitboxes. A new constructor taking `[]ComboStep` is the minimum-complexity expression of the story; all call sites are updated in lockstep (factory, `gameplayer`, tests).
4. **Last-step wrap on `Update()` (not on next-press).** Resetting on the tick step-3 ends (AC4) is safer than deferring to the next press: it prevents stale `stepIndex` between combos and makes `AdvanceCombo()` semantics uniform (it never succeeds from the last step).
5. **State no longer calls `Fire`.** Moving `Fire` to the climber update normalises the entry paths: grounded first press, grounded combo press, and air press all share one call site, and the state's responsibility contracts to "drive hitbox + return to caller when animation ends".
6. **No new contracts.** Combo is an internal elaboration of an existing system. Introducing a contract now would be premature — the only consumer (climber) talks to the concrete type directly.
