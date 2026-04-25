# SPEC — 042-melee-state-driven

**Branch:** `042-melee-state-driven` (already checked out)

**Bounded Context:** Game Logic (`internal/game/entity/actors/...`) with consumer surfaces in Engine (`internal/engine/contracts/combat`, `internal/engine/contracts/vfx`).

---

## 1. Goal

Move the entire melee weapon lifecycle (`Fire`, `Update`, `IsHitboxActive`, `ApplyHitbox`, VFX spawn) out of `ClimberPlayer.Update` and into `MeleeAttackState`. After this story, `ClimberPlayer.Update` no longer references the `melee` weapon. The Actor state machine becomes the sole authority that drives a melee swing — symmetric with how projectile firing is owned by shooting states.

This is a refactor: observable game behaviour is unchanged. Combo, buffering, anim-lock, and ducking guard semantics are preserved.

---

## 2. Ubiquitous Language Touchpoints

- `Actor` — `ClimberPlayer` (embeds `PlatformerCharacter`).
- `State` — `MeleeAttackState` and the per-step states from `MeleeAttackStepStates(n)`.
- `Body` / `Space` — `body.BodiesSpace` passed into `Update` and forwarded to `ApplyHitbox`.
- `Contract` — `combat.Weapon`, `vfx.Manager`, `body.BodiesSpace`, `body.Collidable`.

---

## 3. Contracts

### 3.1 Existing contracts used (no changes)

| Contract | File | Role in this story |
|---|---|---|
| `combat.Weapon` | `internal/engine/contracts/combat/weapon.go` | Base weapon surface (`Fire`, `CanFire`, `Update`, `Cooldown`, `SetCooldown`, `SetOwner`, `ID`). |
| `vfx.Manager` | `internal/engine/contracts/vfx/vfx.go` | `SpawnDirectionalPuff` is invoked from inside `MeleeAttackState.OnStart`. |
| `body.BodiesSpace` | `internal/engine/contracts/body/...` | Forwarded to `weapon.ApplyHitbox`. |
| `actors.StateContributor` | `internal/engine/entity/actors/state_contributor.go` | Existing `meleeContributor` keeps routing to per-step animation enums; unchanged. |

### 3.2 New contracts

**None at the engine layer.** All new typing is package-local to `internal/game/entity/actors/states` so the Bounded Context boundary is preserved.

A new **package-local** interface is introduced inside `melee_state.go` (or a sibling file in the same package) to keep `MeleeAttackState` decoupled from `vfx.Manager`:

```go
// meleeVFXSpawner is the minimum surface needed by MeleeAttackState
// to render the slash VFX. It is satisfied by vfx.Manager.
type meleeVFXSpawner interface {
    SpawnDirectionalPuff(typeKey string, x, y float64, faceRight bool, count int, randRange float64)
}
```

This is **not** added to `internal/engine/contracts/vfx/`. Rationale: the abstraction is a consumer-side narrowing of `vfx.Manager` for a single state; promoting it to engine contracts would expand the public surface for no consumer outside this state. Constitution §"Mock only at system boundaries" — package-local mocks belong in `mocks_test.go` (per the testing standard).

---

## 4. Technical Requirements

### 4.1 `MeleeAttackState` — full ownership of the swing

File: `internal/game/entity/actors/states/melee_state.go`.

**Constructor signature (revised):**

```go
func NewMeleeAttackState(
    owner    meleeOwnerIface,
    space    contractsbody.BodiesSpace,
    w        meleeWeaponIface,
    vfx      meleeVFXSpawner,        // NEW — nil-tolerant
    returnTo actors.ActorStateEnum,
) *MeleeAttackState
```

`vfx` may be `nil`; the state must guard with `if s.vfx != nil` before spawning, mirroring the current guard in `ClimberPlayer.spawnMeleeVFX`.

**`OnStart(currentCount int)` behaviour (revised):**

1. Reset frame counter (`s.frame = 0`).
2. Capture `stepUsed = weapon.StepIndex()` **before** firing (so the captured index reflects the step that is *about to* be swung; equivalent to today's pre-Fire capture order).
3. Call `weapon.Fire(x16, y16, faceDir, body.ShootDirectionStraight, 0)` using `owner.GetPosition16()` and `owner.FaceDirection()`.
4. Spawn the slash VFX at the offset `+12` px in the facing direction (or `-12` px when facing left), matching `ClimberPlayer.spawnMeleeVFX` exactly.

**Pre-condition:** the weapon must satisfy `CanFire() == true` at the moment of `OnStart`. Enforcement of this pre-condition is the **caller's** responsibility (i.e., `ClimberPlayer` only schedules the transition when `melee.CanFire()`). See §4.3.

**Post-condition (after `OnStart`):**

- `weapon.StepIndex()` may have advanced (because `Fire` resets internal swing state); however `s.stepUsed` reflects the index used for this swing.
- One VFX puff has been emitted at the correct offset/direction (when `vfx != nil`).

**`Update()` behaviour (unchanged from current code):**

1. `weapon.Update()`.
2. If `weapon.IsHitboxActive()` → `weapon.ApplyHitbox(space)`.
3. `frame++`. If `frame >= animFrames` → return `returnTo`. Else return `StateMeleeAttack`.

### 4.2 Real factory registration

`init()` in `melee_state.go` must register `StateMeleeAttack` with a factory that constructs a real `MeleeAttackState`, **not** an `actors.IdleState`.

Constraint: `actors.RegisterState` returns the enum at init time, before any `ClimberPlayer` exists. The factory receives only `BaseState` and cannot pull weapon/space/vfx out of thin air. Two acceptable resolutions — pick (B):

- **(A) Per-actor override.** Keep the placeholder factory and rely on the player to install the concrete `MeleeAttackState` via `Character.SetState` / state map at construction time.
- **(B) Inject a default state from the player builder.** Replace the registered placeholder by writing the real `MeleeAttackState` instance into the actor's state map during `SetMelee`, then have `Character.handleState` look up the per-actor instance instead of the global factory's instance.

This story uses **(B)**: the registered factory still returns a benign `IdleState` shell so anonymous lookups don't panic, but `ClimberPlayer.SetMelee` now also calls `actors.RegisterStateForActor(...)` (or equivalent — see §4.3) to install the real `MeleeAttackState` for *this* actor. If the codebase does not yet have a per-actor override hook, the Feature Implementer must add one (smallest viable change in `internal/engine/entity/actors`) — but only if AC-2 cannot otherwise be satisfied.

> **Note for Feature Implementer:** before adding any new engine API, verify whether `Character` already maps state enums to per-actor `ActorState` instances via the state map produced by `builder.PreparePlatformer`. If yes, `SetMelee` simply writes into that map. If no, prefer the smallest extension to that map's API.

### 4.3 `ClimberPlayer` refactor

File: `internal/game/entity/actors/player/climber.go`.

**Removals from `ClimberPlayer.Update`:**

- `p.melee.Update()` call.
- `p.melee.IsHitboxActive()` / `p.melee.ApplyHitbox(space)` block.
- `p.melee.Fire(...)` call.
- `p.spawnMeleeVFX(...)` call.

**What stays in `ClimberPlayer.Update` (combo / buffering / anim-lock orchestration is *not* part of this story — see story 041):**

- `meleeAnimWait` countdown.
- Dash/Jump combo-reset path (`ResetComboOnInterrupt`-equivalent).
- Buffered-press tracking (`meleeBuffered`, `meleeHeldPrev`).
- Movement lock when `isMeleeActive` (lines 134–149).
- The decision of *when* to begin a melee — **but** instead of calling `Fire` directly, `ClimberPlayer` now schedules the state transition (see below).

**How the state transition is scheduled:**

Currently `GroundedState.Update` returns `StateMeleeAttack` when `MeleePressed()` is true. This story keeps that path. Additionally, `ClimberPlayer` must drive the same transition when it would have called `Fire` directly today (combo continuation outside the natural press edge — e.g., buffered re-fire after a swing, mid-air melee from `Falling`). The mechanism:

- `ClimberPlayer` exposes an internal helper `tryEnterMeleeState()` that:
  1. Validates `melee != nil && melee.CanFire() && !melee.IsSwinging() && meleeAnimWait == 0 && !IsDucking()`.
  2. If the player is grounded **and** `melee.ComboWindowRemaining() > 0`, calls `melee.AdvanceCombo()`.
  3. Calls `p.SetNewState(StateMeleeAttack)` (or the platformer-character equivalent setter — match what `Hurt` uses).
- The state's `OnStart` performs the actual `Fire` + VFX.

Replacing the in-line `Fire` with `SetNewState(StateMeleeAttack)` is the load-bearing change. `meleeAnimWait` is still set inside `ClimberPlayer.Update` after the transition request, using `p.meleeStepAnimDuration(p.melee.StepIndex())` (read **after** the state is entered, i.e., after `AdvanceCombo`).

**Wiring change in `SetMelee`:**

```go
func (p *ClimberPlayer) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
    p.melee = w
    p.meleeVFX = vfxMgr
    if w != nil {
        w.SetOwner(p)
        stepStates := gamestates.MeleeAttackStepStates(len(w.Steps()))
        p.AddStateContributor(&meleeContributor{w: w, stepStates: stepStates})

        // NEW — install the real MeleeAttackState for this actor.
        gamestates.InstallMeleeAttackState(p.GetCharacter(), p, w, vfxMgr /*as meleeVFXSpawner*/)
    }
}
```

Where `gamestates.InstallMeleeAttackState` is a new package-level function in `melee_state.go` that constructs and registers the state on the character's state map. Its signature:

```go
func InstallMeleeAttackState(
    char  *actors.Character,
    owner meleeOwnerIface,
    w     meleeWeaponIface,
    vfx   meleeVFXSpawner,
)
```

It chooses `returnTo` based on grounded-vs-air at `OnStart`-time by reading `owner` (i.e., `MeleeAttackState.OnStart` resolves `returnTo` dynamically — see §4.4). To support that, `meleeOwnerIface` is extended:

```go
type meleeOwnerIface interface {
    contractsbody.Collidable
    FaceDirection() animation.FacingDirectionEnum
    IsGrounded() bool   // NEW — already exposed via PlatformerCharacter
}
```

If `PlatformerCharacter` exposes a different name for grounded check (e.g., `!IsFalling() && !IsGoingUp()`), use the existing method directly rather than introducing a new one.

### 4.4 Dynamic `returnTo`

The existing constructor takes a static `returnTo`. After this refactor, `returnTo` should be resolved per swing — at `OnStart` — from the owner's grounded state:

- Grounded → `StateGrounded`.
- Airborne → `actors.Falling`.

This subsumes the current air-melee path (`TryMeleeFromFalling`) so the same `MeleeAttackState` instance handles both grounded and aerial swings without duplication. The `returnTo` field becomes computed in `OnStart`, not constructor-injected. Constructor signature reverts to:

```go
func NewMeleeAttackState(
    owner meleeOwnerIface,
    space contractsbody.BodiesSpace,
    w     meleeWeaponIface,
    vfx   meleeVFXSpawner,
) *MeleeAttackState
```

### 4.5 Field comment on `melee`

```go
// melee is kept as a separate field (not in inventory) because MeleeWeapon
// exposes hitbox lifecycle methods (IsHitboxActive, ApplyHitbox) not present
// on the combat.Weapon interface.
melee *weapon.MeleeWeapon
```

### 4.6 Ducking guard

`GroundedState.Update` currently returns `StateMeleeAttack` on `MeleePressed()` regardless of ducking sub-state. The guard must live in **one** authoritative location. Choose `MeleeAttackState.OnStart`:

- If the owner reports `IsDucking()`, abort the swing: do not call `Fire`, do not spawn VFX, and arrange for `Update` on the next tick to return `returnTo` (set `frame = animFrames` so the very next `Update` resolves immediately). Rationale: keeps `GroundedState` stateless about weapon state; the same guard then naturally protects every entry path (buffered press, falling-air case, future paths).

This requires extending `meleeOwnerIface` with `IsDucking() bool`. `PlatformerCharacter` already has it.

---

## 5. State-Machine Transitions

```
GroundedState (any sub-state except Ducking)
    --MeleePressed && melee.CanFire()-->  StateMeleeAttack
                                              ↓ OnStart: Fire + VFX (capture stepUsed, dynamic returnTo)
                                              ↓ Update × animFrames: weapon.Update + ApplyHitbox
                                              ↓ frame == animFrames
                                          returnTo (StateGrounded | actors.Falling)

Falling (mid-air)
    --MeleePressed && melee.CanFire()-->  StateMeleeAttack  (returnTo computed = Falling)

GroundedState (Ducking sub-state)
    --MeleePressed-->                     StateMeleeAttack  (OnStart fast-aborts)
                                              ↓ next Update returns StateGrounded (no Fire, no VFX)
```

---

## 6. Pre / Post-Conditions Summary

| Boundary | Pre | Post |
|---|---|---|
| `MeleeAttackState.OnStart` | `weapon != nil`; `owner != nil`; `space != nil`. | Either: (a) `weapon.IsSwinging()==true`, exactly one VFX emitted, `stepUsed` recorded; OR (b) ducking-abort path: no Fire, no VFX, `frame >= animFrames`. |
| `MeleeAttackState.Update` | `OnStart` was called; `frame < animFrames` (else returns `returnTo`). | `weapon.Update` called once per tick; `ApplyHitbox` called only on ticks where `IsHitboxActive()` is true. |
| `ClimberPlayer.Update` | `space != nil`. | No call to `melee.Fire`, `melee.Update`, `melee.IsHitboxActive`, `melee.ApplyHitbox`, or `spawnMeleeVFX` appears anywhere on this method's call graph (statically — `grep` produces zero hits in `climber.go`). |
| `ClimberPlayer.SetMelee` | `w != nil ⇒ w.Steps()` non-empty. | The character's state map binds `StateMeleeAttack` to a fully-wired `MeleeAttackState`. |

---

## 7. Integration Points Within the Bounded Context

| File | Change |
|---|---|
| `internal/game/entity/actors/states/melee_state.go` | New `meleeVFXSpawner` iface; revised `NewMeleeAttackState` signature (adds `vfx`); revised `OnStart` (Fire + VFX + ducking abort + dynamic returnTo); new `InstallMeleeAttackState` helper; `meleeOwnerIface` gains `IsDucking()` and grounded check. |
| `internal/game/entity/actors/states/melee_state_test.go` | Tests updated by **TDD Specialist** to remove explicit `w.Fire(...)` calls before `OnStart` (now owned by the state) and to cover VFX spawn + ducking-abort + dynamic returnTo. Existing tests for `TryMeleeFromFalling` and `ResetComboOnInterrupt` remain. |
| `internal/game/entity/actors/states/grounded_state.go` | No change required (ducking abort handled in state OnStart). |
| `internal/game/entity/actors/player/climber.go` | Remove melee weapon drive from `Update`; replace with `tryEnterMeleeState` helper that calls `SetNewState(StateMeleeAttack)`; add field comment; extend `SetMelee` with `InstallMeleeAttackState`. |
| `internal/game/scenes/phases/player.go` | No change (still calls `climber.SetMelee(...)`). |

---

## 8. Red Phase Scenario (failing tests for TDD Specialist)

The TDD Specialist must produce these failing tests **before** any production code is written. All tests live in existing `_test.go` files in `internal/game/entity/actors/states/` and `internal/game/entity/actors/player/`.

### RED-1: `OnStart` fires the weapon (new behaviour)

```
Given: a fresh MeleeAttackState with a CanFire weapon (no prior w.Fire call).
When:  st.OnStart(0) is called.
Then:  weapon.IsSwinging() == true  AND  exactly 1 invocation was recorded
       on the weapon-fire spy.
```
Currently fails because `OnStart` does not call `Fire`.

### RED-2: `OnStart` spawns VFX with the correct facing offset

```
Given: a MeleeAttackState wired with a spy meleeVFXSpawner; owner faces right
       at GetPosition16() = (100*16, 50*16).
When:  st.OnStart(0).
Then:  spy recorded exactly 1 SpawnDirectionalPuff call with typeKey=="melee_slash",
       faceRight==true, x close to fp16.From16(100*16 + fp16.To16(12)),
       y close to fp16.From16(50*16).
And:   With owner facing left, faceRight==false and x is offset by -12 px.
```
Currently fails because `OnStart` does not spawn VFX (and the constructor does not accept a vfx spawner).

### RED-3: `OnStart` aborts when owner is ducking

```
Given: owner.IsDucking() == true; weapon.CanFire() == true.
When:  st.OnStart(0); next st.Update().
Then:  No call to weapon.Fire was recorded.
       No VFX spawn was recorded.
       Update() returns returnTo (StateGrounded for grounded owner).
```
Currently fails because there is no ducking guard in the state.

### RED-4: `returnTo` is computed dynamically from grounded state

```
Case A (grounded owner): owner.IsGrounded() == true.
   On final Update tick → StateGrounded.
Case B (airborne owner): owner.IsGrounded() == false.
   On final Update tick → actors.Falling.
```
The constructor no longer takes `returnTo` as a parameter; tests construct without it. Currently fails because the constructor still requires `returnTo`.

### RED-5: `ClimberPlayer.Update` makes no melee weapon calls

```
Given: a ClimberPlayer wired with a recording spy weapon and a recording space.
When:  Update(space) is invoked over a sequence of frames including a melee press.
Then:  spy.FireCount == 0
       spy.UpdateCount == 0
       spy.IsHitboxActiveCount == 0
       spy.ApplyHitboxCount == 0
And:   The climber transitioned to StateMeleeAttack (via SetNewState) on the press frame.
```
Currently fails because `ClimberPlayer.Update` calls all four methods directly.

### RED-6: hitbox is applied across the active window (regression)

Adapt the existing `TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow` to the new `OnStart`-fires contract: remove the explicit `w.Fire(...)` setup line and assert `len(enemy.damageCalls) == 1` after the full animation. Currently this test passes for the wrong reason (test fires manually). After removing the explicit Fire, today's production code under-fires (no Fire happens), so the test will fail until the new `OnStart` behaviour is in place.

### RED-7: Step capture order is preserved

Update `TestMeleeAttackState_UsesCurrentComboStep` to remove the explicit pre-`OnStart` `w.Fire(...)` calls. The state's `OnStart` must:
1. Read `StepIndex()` → save into `stepUsed`.
2. Call `Fire` (which advances/resets internal step bookkeeping).

Assert `StepUsed()` returns 0, 1, 2 across three consecutive swings driven only by the state.

---

## 9. Out of Scope

- Combo definition changes, damage/cooldown/hitbox tuning.
- Combo-window timing changes (story 041).
- Promoting `meleeVFXSpawner` to `internal/engine/contracts/vfx`.
- Generalizing the per-actor state-instance install pattern to other states (only `StateMeleeAttack` is touched).
- Renaming or removing `TryMeleeFromFalling` / `ResetComboOnInterrupt` helpers (kept for `ClimberPlayer.Update` orchestration).

---

## 10. Risks & Notes

- **Per-actor state install (§4.2 option B)** is the only place this spec touches engine internals. If `Character`'s state map already supports per-actor instance override, the change is one new package-level helper in `melee_state.go`; otherwise the Feature Implementer should add the smallest possible setter on `Character` (e.g., `SetStateInstance(enum, ActorState)`).
- **Order of ops in `OnStart`**: capture `stepUsed` *before* `Fire`. `Fire` may mutate `StepIndex` (it does not today — `AdvanceCombo` does — but the order is defensive and matches existing test expectations).
- **VFX nil tolerance** is required for unit tests that don't supply a manager.
- **`tryEnterMeleeState` helper** in `climber.go` is the single replacement for the inline Fire path. Buffered-press handling continues to call this helper.
