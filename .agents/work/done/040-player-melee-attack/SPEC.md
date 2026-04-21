# SPEC — US-040 Player Melee Attack

**Branch:** `040-player-melee-attack`
**Bounded Contexts:**
- Engine: `internal/engine/combat/weapon/` (new `MeleeWeapon`, factory extension)
- Game: `internal/game/entity/actors/states/` (new `StateMeleeAttack`)

**Dependencies / Referenced Contracts:**
- `internal/engine/contracts/combat/weapon.go` — `Weapon`, `Faction`, `Factioned` (US-022, US-038)
- `internal/engine/contracts/combat/damageable.go` — `Damageable` (US-038)
- `internal/engine/contracts/body/body.go` — `BodiesSpace.Query(rect)`, `Collidable`, `MovableCollidable` (existing)
- `internal/engine/entity/actors/actor_state.go` — `ActorState`, `ActorStateEnum`, `BaseState`, `RegisterState` (existing)
- `internal/engine/input/commands.go` — `PlayerCommands.Melee` (mapped to `KeyZ`, already in place)

---

## 1. Technical Requirements

### 1.1 `MeleeWeapon` (new — `internal/engine/combat/weapon/melee.go`)

Implements the existing `combat.Weapon` interface (no interface changes). Since melee does not spawn projectiles, `Fire(...)` starts an active-frame swing window instead of emitting a projectile; damage resolution is performed via `ApplyHitbox(...)` each frame the hitbox is live (called by `StateMeleeAttack.Update()`).

```go
type MeleeWeapon struct {
    id              string
    damage          int
    cooldownFrames  int
    currentCooldown int
    activeFrames    [2]int // [startFrame, endFrame] inclusive, relative to Fire()
    hitboxW16       int
    hitboxH16       int
    hitboxOffX16    int // offset from owner origin, positive = forward
    hitboxOffY16    int
    owner           interface{}
    swingFrame      int  // -1 when idle; 0..∞ while active
    swinging        bool
    hitThisSwing    map[combat.Damageable]struct{} // prevent multi-hit per swing
}

// Weapon interface satisfaction
func (w *MeleeWeapon) ID() string
func (w *MeleeWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int)
func (w *MeleeWeapon) CanFire() bool
func (w *MeleeWeapon) Update() // decrements cooldown AND advances swingFrame
func (w *MeleeWeapon) Cooldown() int
func (w *MeleeWeapon) SetCooldown(frames int)
func (w *MeleeWeapon) SetOwner(owner interface{})

// Melee-specific surface (called from StateMeleeAttack.Update)
func (w *MeleeWeapon) IsHitboxActive() bool                 // true when startFrame <= swingFrame <= endFrame
func (w *MeleeWeapon) ApplyHitbox(space body.BodiesSpace)   // queries space, applies damage with faction gating
func (w *MeleeWeapon) Damage() int                          // getter for tests
func (w *MeleeWeapon) ActiveFrames() [2]int                 // getter for tests
```

Semantics:
- `Fire(...)` — if `!CanFire()` returns early. Otherwise sets `swinging=true`, `swingFrame=0`, resets `hitThisSwing`, stores direction (forward = owner's `FaceDirection`), sets `currentCooldown = cooldownFrames`.
- `Update()` — decrements `currentCooldown` when > 0; if `swinging`, increments `swingFrame`. When `swingFrame > activeFrames[1]` and no recovery needed, `swinging=false`.
- `IsHitboxActive()` — returns `swinging && activeFrames[0] <= swingFrame <= activeFrames[1]`.
- `ApplyHitbox(space)` — no-op if `!IsHitboxActive()`. Computes hitbox rect from owner (via `owner.(body.Collidable)` or `body.Positioned`): origin + `hitboxOff*16` flipped by `FaceDirection`, size `hitboxW16 x hitboxH16` (converted to pixel rect via `fp16.From16`). Calls `space.Query(rect)` and for each result:
  1. Skip if the target IS the owner (self-damage guard).
  2. Skip if target implements `combat.Factioned` AND owner implements `combat.Factioned` AND `target.Faction() == owner.Faction()` (AC4 — same-faction gate, reuses US-038 policy).
  3. If target implements `combat.Damageable` and not already in `hitThisSwing`, call `target.TakeDamage(w.damage)` and add to `hitThisSwing`.

Rationale:
- `Weapon.Fire` signature is preserved so that the Inventory / ShootingSkill pipeline remains type-safe. The `direction` and `faceDir` parameters are used to select hitbox mirroring; `state` is unused.
- Per-swing idempotent damage (`hitThisSwing`) prevents a single swing from dealing `N * damage` when an enemy overlaps the hitbox across multiple frames — consistent with single-hit melee semantics.

### 1.2 `weapon.Factory` extension (`internal/engine/combat/weapon/factory.go`)

Current behaviour: only `"type": "projectile"` is supported; other types return `fmt.Errorf("unsupported weapon type: %s", ...)`.

Change: route on `config.Type`:
- `"projectile"` → existing path.
- `"melee"` → new path; parses the schema in AC6 and returns `*MeleeWeapon`.
- otherwise → existing error.

New struct (private to `factory.go` or a shared anonymous struct in `NewWeaponFromJSON`):

```go
type meleeConfig struct {
    ID             string `json:"id"`
    Type           string `json:"type"`
    Damage         int    `json:"damage"`
    CooldownFrames int    `json:"cooldown_frames"`
    ActiveFrames   [2]int `json:"active_frames"` // [start, end]
    Hitbox         struct {
        Width   int `json:"width"`    // pixels
        Height  int `json:"height"`   // pixels
        OffsetX int `json:"offset_x"` // pixels
        OffsetY int `json:"offset_y"` // pixels
    } `json:"hitbox"`
}
```

Validation:
- `ActiveFrames[0] >= 0`, `ActiveFrames[1] >= ActiveFrames[0]` → else return `fmt.Errorf("invalid active_frames")`.
- `Hitbox.Width > 0 && Hitbox.Height > 0` → else return error.
- `CooldownFrames >= 0`.

On success the factory converts pixel values to fp16 (`fp16.To16`) and constructs the weapon.

### 1.3 `StateMeleeAttack` (new — `internal/game/entity/actors/states/melee_state.go`)

```go
// Package-level enum, registered at init time like StateDashing.
var StateMeleeAttack actors.ActorStateEnum

func init() {
    StateMeleeAttack = actors.RegisterState("melee_attack", func(b actors.BaseState) actors.ActorState {
        return &actors.IdleState{BaseState: b} // placeholder; MeleeAttackState constructed directly
    })
}

type MeleeAttackState struct {
    body     contractsbody.MovableCollidable
    space    contractsbody.BodiesSpace
    weapon   *weapon.MeleeWeapon   // concrete; we need IsHitboxActive/ApplyHitbox surface
    returnTo actors.ActorStateEnum // origin state (Grounded or Falling)
    started  bool
    finished bool
    tick     int
}

func NewMeleeAttackState(b contractsbody.MovableCollidable, space contractsbody.BodiesSpace,
                        w *weapon.MeleeWeapon, returnTo actors.ActorStateEnum) *MeleeAttackState
```

Lifecycle:
- `OnStart(currentCount)` — if `!w.CanFire()` set `finished=true` and return (this protects against re-entry while on cooldown; the caller that transitioned here is responsible for gating, but we keep the guard for safety). Otherwise call `w.Fire(x16, y16, body.FaceDirection(), body.ShootDirectionStraight, 0)`, set `started=true`, `tick=0`.
- `Update() actors.ActorStateEnum` — advance `tick`; if `w.IsHitboxActive()` call `w.ApplyHitbox(space)`. When `IsAnimationFinished()` returns true → set `finished=true`, return `returnTo` (AC2 / AC8 air-melee case).
- `OnFinish()` — no-op (weapon's internal `swinging` flag is cleared by its own `Update()`).

`IsAnimationFinished()` is inherited from `BaseState` and backed by animation data; for tests where no animation data exists, the state test double will set a deterministic frame count via a virtual clock.

### 1.4 Trigger wiring in `GroundedState` and `FallState`

**`GroundedState.Update()`** (file: `internal/game/entity/actors/states/grounded_state.go`):
- Extend `GroundedInput` contract with `MeleePressed() bool`.
- After `JumpPressed` / `DashPressed` checks, before sub-state dispatch:
  ```go
  if input.MeleePressed() {
      return StateMeleeAttack
  }
  ```
- The existing `DashPressed` short-circuit placement ensures melee does not fire while entering a dash; the dash state itself (`DashState.Update`) does not consume melee input because it does not route through `GroundedInput` — this satisfies AC5 "cannot be triggered while dashing".

**`FallState`** (file: `internal/engine/entity/actors/shooting_states.go` / `FallState`):
- Falling is currently an engine-level state. Air-melee trigger is added in the game layer via the Character's state driver. We introduce a small `FallingInput` consumer check in the game layer; concretely, the Character's `handleState` hook reads `PlayerCommands.Melee` when in `Falling` and transitions to `StateMeleeAttack` with `returnTo = Falling`.
- Concrete placement: add a `tryMeleeFromFalling(cmds input.PlayerCommands) (actors.ActorStateEnum, bool)` helper in `internal/game/entity/actors/states/melee_state.go` consumed by the Character wiring. This helper returns `(StateMeleeAttack, true)` when `cmds.Melee && !alreadyInMelee && weapon.CanFire()`.

Pre-condition for entering `StateMeleeAttack` from either source:
- `w.CanFire() == true` (cooldown elapsed).
- Current state is NOT `StateMeleeAttack`.
- Current state is NOT `StateDashing`.

### 1.5 State machine transitions

```
Grounded --[Melee pressed, weapon.CanFire]--> MeleeAttack --[IsAnimationFinished]--> Grounded
Falling  --[Melee pressed, weapon.CanFire]--> MeleeAttack --[IsAnimationFinished]--> Falling
Dashing  --[Melee pressed]--X (blocked; dash does not consume melee)
MeleeAttack --[Melee re-pressed]--X (guard: already swinging; weapon.CanFire == false)
```

The `returnTo` field captured at `NewMeleeAttackState(...)` construction is the single source of truth for where the state machine returns; this is the direct mechanism for AC8 "air melee transitions back to Falling (not Grounded)".

### 1.6 JSON schema (AC6)

Spec-locked schema for a melee entry in the player weapons config:

```json
{
  "id": "player_melee",
  "type": "melee",
  "damage": 1,
  "cooldown_frames": 20,
  "active_frames": [4, 10],
  "hitbox": {
    "width": 24,
    "height": 16,
    "offset_x": 12,
    "offset_y": 0
  }
}
```

Field contract:
- `id` — string, unique within Inventory.
- `type` — literal `"melee"`.
- `damage` — int ≥ 0.
- `cooldown_frames` — int ≥ 0.
- `active_frames` — `[startFrame, endFrame]`, both ints, inclusive; `start >= 0`, `end >= start`.
- `hitbox` — all values in pixels (consistent with projectile `spawn_offset` convention); converted to fp16 internally via `fp16.To16`.
- `offset_x` is mirrored when `FaceDirection == FaceDirectionLeft`.

---

## 2. Integration Points

- `combat.Inventory` — `MeleeWeapon` is added as an additional weapon via the existing inventory API; no changes to the Inventory contract.
- `ShootingSkill` — unchanged. Melee is driven by a distinct input (`PlayerCommands.Melee`) and owns its own state (`StateMeleeAttack`), so it does NOT route through `ShootingSkill.HandleInput`. Rationale: melee is stateful at the actor level (a state machine node), whereas `ShootingSkill` assumes weapons produce projectiles through an inventory pipeline.
- `Character.TakeDamage` / `Character.Faction()` — reused as the concrete `Damageable` + `Factioned` implementation for enemies and the player. No changes required.
- `BodiesSpace.Query(image.Rectangle)` — reused as-is. The melee weapon translates its fp16 hitbox to a pixel-space `image.Rectangle` for querying.

---

## 3. Pre- and Post-Conditions (per AC)

| AC | Pre-condition | Post-condition |
|---|---|---|
| AC1 | `MeleeWeapon` struct has the specified fields | `var _ combat.Weapon = (*MeleeWeapon)(nil)` compiles; all fields accessible via accessors or tests |
| AC2 | `StateMeleeAttack` registered via `actors.RegisterState` at init | Calling `Update()` on a `MeleeAttackState` returns `returnTo` on the first tick where `IsAnimationFinished()==true`, otherwise `StateMeleeAttack` |
| AC3 | Weapon configured with `active_frames = [a, b]` and `Fire()` has been called | `IsHitboxActive()` returns `false` for `swingFrame < a`; `true` for `a <= swingFrame <= b`; `false` for `swingFrame > b` |
| AC4 | Two `Damageable + Factioned` actors overlapping the hitbox: one same-faction as owner, one different | Only the different-faction actor receives `TakeDamage(w.damage)`; same-faction actor is untouched |
| AC5 | Character is in `StateGrounded` or `Falling`, not dashing, not already in `StateMeleeAttack`, weapon `CanFire()==true`, `PlayerCommands.Melee==true` | State transitions to `StateMeleeAttack`. Attempting trigger during `StateDashing` or `StateMeleeAttack` has no effect. |
| AC6 | JSON config matches schema in §1.6 | `weapon.NewWeaponFromJSON(data, nil)` returns a `*MeleeWeapon` with fields populated; invalid configs (missing `hitbox`, inverted `active_frames`) return a non-nil error |
| AC7 | `"type": "melee"` passed in config | Factory returns a `*MeleeWeapon` instance; unsupported types still return the existing error |
| AC8 | Full wiring (state + weapon + space with a `Damageable` enemy) | All six bullet conditions in the story pass |

---

## 4. Red Phase — Failing Test Scenarios (for TDD Specialist)

### RED-1 — `internal/engine/combat/weapon/melee_test.go` (new)

Table-driven test `TestMeleeWeapon_Fire_HitboxActivation`:
- Fixture: `MeleeWeapon{damage:1, activeFrames:[3,5], cooldown:20, hitbox 24x16 offset 12/0}`, owner is a stub `Collidable+Factioned(FactionPlayer)` positioned at `(100, 100)`.
- Cases (frame index = number of `Update()` calls after `Fire()`):
  | frame | expected IsHitboxActive |
  |---|---|
  | 0 | false |
  | 2 | false |
  | 3 | true |
  | 4 | true |
  | 5 | true |
  | 6 | false |

Test `TestMeleeWeapon_ApplyHitbox_FactionGating`:
- Space contains: (a) enemy at overlap with `FactionEnemy`, (b) ally at overlap with `FactionPlayer`, (c) enemy outside rect.
- Call `Fire()`, `Update()` three times to reach active window, `ApplyHitbox(space)`.
- Assert: enemy (a).`TakeDamage` called once with amount=1; ally (b) and far enemy (c) never called.

Test `TestMeleeWeapon_ApplyHitbox_SingleHitPerSwing`:
- Keep hitbox active for 3 frames, call `ApplyHitbox` each frame.
- Assert: `TakeDamage` called exactly once on the overlapping enemy.

Test `TestMeleeWeapon_Cooldown_PreventsRefire`:
- `Fire()` once; immediately `CanFire()` → false. Advance `Update()` `cooldownFrames-1` times → still false. One more → true.

Test `TestMeleeWeapon_Fire_MirrorsHitboxWhenFacingLeft`:
- Owner facing left, `offsetX=12`. Assert computed query rect is to the left of the owner origin, not the right.

### RED-2 — `internal/engine/combat/weapon/factory_test.go` (extend)

Add cases to the existing factory test table:
- `"melee weapon ok"` — the full AC6 JSON → expect `*MeleeWeapon` with fields matching, no error.
- `"melee missing hitbox"` → error containing "hitbox".
- `"melee inverted active_frames"` (e.g. `[10, 4]`) → error containing "active_frames".
- `"unknown type"` — regression; must still error.

### RED-3 — `internal/game/entity/actors/states/melee_state_test.go` (new)

Test `TestMeleeAttackState_ReturnsToGrounded_WhenAnimationFinishes`:
- Construct `MeleeAttackState` with `returnTo = StateGrounded`, stubbed `IsAnimationFinished()` that becomes true on tick N.
- Call `OnStart`, then `Update()` repeatedly. Assert: first N-1 calls return `StateMeleeAttack`; Nth returns `StateGrounded`.

Test `TestMeleeAttackState_AirMelee_ReturnsToFalling`:
- Same as above but `returnTo = actors.Falling`. Assert final transition goes to `Falling`, not `StateGrounded` (AC8 explicit).

Test `TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow`:
- Inject a real `MeleeWeapon` and a `bodiesSpace` containing a `Damageable+Factioned` enemy overlapping the hitbox.
- Run the state until end. Assert `TakeDamage` called exactly once during the active window.

Test `TestGroundedState_MeleePressed_TransitionsToMeleeAttack`:
- Extend the existing `mockGroundedInput` in `grounded_state_test.go` with `MeleePressed()`.
- Set `MeleePressed=true`, call `GroundedState.Update()`. Assert return is `StateMeleeAttack`.

Test `TestGroundedState_DashPressed_TakesPrecedenceOverMelee`:
- `DashPressed=true, MeleePressed=true` → return is `StateDashing` (AC5 negative case).

Test `TestMeleeTrigger_BlockedDuringCooldown`:
- Weapon on cooldown. `MeleePressed=true`. Assert no transition to `StateMeleeAttack`.

### RED-4 — input contract regression

Add a case to `internal/engine/input/commands_test.go` (if not already covered): `PlayerCommands.Melee` maps to `KeyZ` — existing coverage noted at lines 28, 33, 135; no change needed unless a gap surfaces.

---

## 5. Non-Goals (out of scope for US-040)

- Combo / multi-hit sequences (handled by US-041).
- Melee animation assets beyond the state registration hook — animation data can be stubbed in tests; real sprite sheets land with art pass.
- Knockback / hitstun on damaged targets (future story).
- Directional melee (up-stab, down-stab) — current spec is `ShootDirectionStraight` only.

---

## 6. Key Design Decisions

1. **`Weapon` interface is not modified.** `MeleeWeapon` satisfies it natively and adds a melee-only `ApplyHitbox(space)` surface consumed exclusively by `StateMeleeAttack`. This keeps `Inventory` / `ShootingSkill` untouched.
2. **`returnTo` captured at state construction.** This is the cleanest implementation of AC8 (air melee returns to `Falling`) without re-querying grounded/falling at end-of-state, which would race with gravity.
3. **Single-hit-per-swing guard.** Without a `hitThisSwing` set, a multi-frame active window would compound damage. The guard makes damage deterministic and matches the AC phrasing "applies `damage`" (singular).
4. **Factory validates early.** Invalid `active_frames` / missing `hitbox` fail at load time, not at swing time — this keeps tests fast and content errors loud.
