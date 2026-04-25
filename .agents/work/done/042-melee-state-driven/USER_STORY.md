# User Story 042 — Melee State-Driven Lifecycle

## Story

**As a** game developer,
**I want** the melee weapon lifecycle (fire, hitbox activation, hitbox application) to be owned entirely by `MeleeAttackState`,
**so that** the Actor state machine is the single authoritative driver of melee behaviour, consistent with how projectile weapons are fired through states.

---

## Background

`ClimberPlayer.Update()` currently drives the melee weapon directly:

```
if meleePressed && p.melee.CanFire() && !p.IsDucking() {
    p.melee.Fire(...)
    p.spawnMeleeVFX(...)
}
if p.melee.IsHitboxActive() {
    p.melee.ApplyHitbox(space)
}
```

`GroundedState.Update()` already returns `StateMeleeAttack` when `MeleePressed()` is true, so the state transition is owned by the state machine. However, the actual weapon drive (`weapon.Update`, `Fire`, `IsHitboxActive`, `ApplyHitbox`) bypasses the state and lives in the Actor's top-level `Update()`. This splits the melee lifecycle across two layers, making the state machine incomplete.

`MeleeAttackState` has a stub implementation that already calls `Fire` on `OnStart` and drives `Update`/`IsHitboxActive`/`ApplyHitbox` per frame — but it is registered with a placeholder `actors.IdleState` factory and is not yet wired into `ClimberPlayer`.

---

## Acceptance Criteria

### AC-1: MeleeAttackState owns the full melee lifecycle

`MeleeAttackState.OnStart` fires the weapon (`weapon.Fire`) and spawns the VFX.
`MeleeAttackState.Update` calls `weapon.Update()`, checks `IsHitboxActive()`, and calls `ApplyHitbox(space)` each frame until the animation completes, then returns `returnTo`.
`ClimberPlayer.Update()` no longer calls `melee.Update()`, `melee.Fire()`, `melee.IsHitboxActive()`, or `melee.ApplyHitbox()`.

### AC-2: MeleeAttackState is registered with its real constructor

`StateMeleeAttack` is registered via `actors.RegisterState` using a factory that constructs a fully wired `MeleeAttackState`, not `actors.IdleState`.

### AC-3: VFX spawn is delegated to MeleeAttackState

`MeleeAttackState.OnStart` calls `spawnMeleeVFX` (or an equivalent VFX spawn path) so the slash effect fires at the correct position and facing direction when the state begins — not from `ClimberPlayer.Update()`.

### AC-4: `ClimberPlayer.melee` field carries a code comment

The `melee *weapon.MeleeWeapon` field on `ClimberPlayer` carries a one-line comment explaining why it is a separate field and not stored in the inventory:

```go
// melee is kept as a separate field (not in inventory) because MeleeWeapon
// exposes hitbox lifecycle methods (IsHitboxActive, ApplyHitbox) not present
// on the combat.Weapon interface.
melee *weapon.MeleeWeapon
```

### AC-5: Behaviour is preserved

After the refactor, the observable game behaviour is unchanged:
- A melee swing fires exactly once per `MeleePressed()` edge (no double-fire).
- The hitbox is active for the same window defined by the weapon parameters ([4..10] frames).
- `IsDucking()` guard: a melee swing cannot begin while the player is ducking (this guard may live in `GroundedState` sub-state transition or `MeleeAttackState.OnStart`).
- The slash VFX appears at the correct offset from the player's position, flipped for the correct facing direction.

### AC-6: Test coverage

Unit tests cover:
- `MeleeAttackState.OnStart` fires the weapon and triggers VFX.
- `MeleeAttackState.Update` drives the weapon each frame and calls `ApplyHitbox` when the hitbox is active.
- `MeleeAttackState.Update` returns `returnTo` once `animFrames` is exhausted.
- `ClimberPlayer.Update` does not call any melee weapon methods when `MeleeAttackState` is active (or equivalently, the melee block is absent from `ClimberPlayer.Update`).

---

## Scope

**In scope:**
- Remove the melee drive block from `ClimberPlayer.Update()`.
- Wire `MeleeAttackState` as a fully functional state constructor in `melee_state.go`.
- Pass VFX spawn capability into `MeleeAttackState` (via interface or callback).
- Add the explanatory comment on `ClimberPlayer.melee`.

**Out of scope:**
- Changes to the melee weapon parameters (damage, cooldown, hitbox size).
- Any new combo or chaining logic (see story 041).
- Changes to `GroundedState` beyond what is needed to ensure the ducking guard is correctly enforced.

---

## Relevant Files

| File | Role |
|---|---|
| `internal/game/entity/actors/player/climber.go` | Remove melee drive; add field comment |
| `internal/game/entity/actors/states/melee_state.go` | Own full lifecycle; real factory registration; VFX spawn |
| `internal/game/entity/actors/states/grounded_state.go` | Verify ducking guard is enforced before returning `StateMeleeAttack` |
| `internal/game/entity/actors/player/weapons.go` | `NewPlayerMeleeWeapon` factory (read-only reference) |
| `internal/game/entity/actors/phases/player.go` | Injection site for `MeleeAttackState` wiring |
