# US-010 — Player Shooting Skill (Cuphead-style)

**Branch:** `010-player-shooting-skill`
**Bounded Context:** Game Logic (`internal/game/`)

## Story

As a player,
I want to hold a shoot button and fire a continuous stream of bullets in the direction I am facing,
so that I can attack enemies with a rapid, directional ranged attack similar to Cuphead.

## Context

Shooting is a **concurrent skill** — it runs alongside the active `GroundedState` sub-states (idle, walking, aim-lock) without replacing them. When `AimLockHeld()` is true the Actor freezes horizontal movement and locks the facing direction; shooting can occur in any grounded sub-state as long as `ShootHeld()` is true.

Reference mechanics (https://github.com/setanarut/cuphead):
- Bullets spawn at a fixed offset from the Actor's center, alternating Y-offset via `OffsetToggler` (US-008).
- A cooldown (fire rate) limits how many bullets can be spawned per second.
- Bullets travel horizontally in the Actor's `FaceDirection` at a fixed speed.
- Bullets are destroyed on collision with an enemy or a solid tile.

## Acceptance Criteria

- **AC1** — `GroundedInput` gains a `ShootHeld() bool` method; all existing implementors must satisfy the updated interface.
- **AC2** — `GroundedState.Update()` calls a `ShootingSkill` on every tick where `ShootHeld()` is true, independent of the active sub-state.
- **AC3** — `ShootingSkill` enforces a configurable cooldown (in frames); it spawns at most one bullet per cooldown window.
- **AC4** — Each spawned bullet's Y-offset alternates between `+N` and `−N` pixels using `OffsetToggler` (already implemented in US-008).
- **AC5** — Bullets spawn at the Actor's current position plus a configurable horizontal spawn offset in the `FaceDirection`.
- **AC6** — A bullet `Body` travels at a fixed configurable speed in the `FaceDirection` each frame until it leaves the `BodiesSpace` or collides with a `Collidable` that is not its owner.
- **AC7** — Shooting is suppressed (no bullet spawned) while the Actor is in `StateDashing`.
- **AC8** — `ShootingSkill` is injected into `GroundedDeps`; it is never a global singleton.
- **AC9** — Unit tests cover: cooldown gating (no double-spawn within cooldown), alternating Y-offset over ≥4 shots, and suppression when dashing.

## Behavioral Edge Cases

- Holding shoot while transitioning from `SubStateWalking` → `SubStateAimLock` must not reset the cooldown counter or skip a shot.
- Releasing and immediately re-pressing shoot within the same cooldown window must not spawn an extra bullet.
- If the Actor's `FaceDirection` changes mid-cooldown (e.g. player turns around), the next bullet must use the new direction.
- A bullet that exits the `BodiesSpace` bounds must be queued for removal via `BodiesSpace.QueueForRemoval()`.

## Notes

- `ShootingSkill` lives in `internal/game/entity/actors/states/` alongside the other grounded sub-state files.
- Bullet entity lives in `internal/game/entity/actors/` (or a dedicated `bullets/` sub-package if the implementer prefers).
- Reuse `OffsetToggler` from US-008 — do not duplicate it.
- Fixed-point positions must use `x16`/`y16` conventions per the constitution.
