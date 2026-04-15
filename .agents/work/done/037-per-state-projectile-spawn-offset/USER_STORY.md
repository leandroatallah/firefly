# US-037 — Per-State Projectile Spawn Offset

**Branch:** `037-per-state-projectile-spawn-offset`
**Bounded Context:** Engine (`internal/engine/combat/weapon/`) + Game (`internal/game/`)

## Story

As a game developer using this boilerplate,
I want each actor state to define its own projectile spawn offset,
so that the muzzle position changes naturally when the character is ducking, jumping, or shooting diagonally — like Metal Slug or Cuphead.

## Context

US-034 introduced a single static `spawnOffsetX16`/`spawnOffsetY16` on `ProjectileWeapon`. This works for a single-pose character, but real action games change the gun arm position depending on animation state. A crouching player shoots from a lower position; an upward shot comes from a raised arm; a falling shot may come from a different side.

The natural place to store this data is the entity JSON (`assets/entities/player/climber.json`), alongside the existing per-state sprite/collision data under `skills.shooting`. State names in the JSON match the strings passed to `actors.RegisterState` (e.g. `"duck"`, `"jump_shoot"`, `"idle"`).

At runtime, `ProjectileWeapon.Fire()` receives the actor's current `ActorStateEnum`. It looks up the matching offset in a `map[ActorStateEnum][2]int` (x16, y16). If no entry exists for the current state, it falls back to the default offset from US-034.

## Proposed JSON Structure

```json
"skills": {
  "shooting": {
    "enabled": true,
    "cooldown_frames": 15,
    "projectile_speed": 6,
    "projectile_range": 4,
    "directions": 4,
    "spawn_offset": { "x": 8, "y": 8 },
    "state_spawn_offsets": {
      "duck":       { "x": 6, "y": 12 },
      "jump_shoot": { "x": 8, "y": 5 },
      "fall_shoot": { "x": 8, "y": 6 }
    }
  }
}
```

`spawn_offset` is the existing default (backward-compatible). `state_spawn_offsets` maps state name → offset in pixels (converted to fp16 on load, not at fire time).

## Acceptance Criteria

- **AC1** — Entity JSON `skills.shooting` supports an optional `state_spawn_offsets` map: `map[string]{"x": int, "y": int}`. Absent map → behavior identical to US-034 (default offset only).
- **AC2** — At load time, state names in `state_spawn_offsets` are resolved to `ActorStateEnum` via `actors.GetStateEnum()`. Unknown names are logged as a warning and skipped (not fatal).
- **AC3** — `ProjectileWeapon` stores `stateOffsets map[ActorStateEnum][2]int` in addition to the existing `spawnOffsetX16`/`spawnOffsetY16`.
- **AC4** — `ProjectileWeapon.Fire()` accepts the current `ActorStateEnum` (new parameter or via a state provider interface). It looks up `stateOffsets[state]`; if found, those values override the default offset for this shot only.
- **AC5** — Facing-direction X-axis negation (from US-034 AC3) still applies to per-state offsets.
- **AC6** — Default offset of `(0, 0)` with no `state_spawn_offsets` entry produces identical behavior to US-034 (backward compatible).
- **AC7** — Unit tests cover: state with matching offset uses it; state without entry falls back to default; facing-left negates X for per-state offsets; unknown state name in JSON is silently skipped.

## Design Notes

- `Fire()` signature change: add `state actors.ActorStateEnum` parameter, or introduce a `StateProvider` interface (`CurrentState() actors.ActorStateEnum`) injected at construction — prefer the interface to avoid coupling weapon to the actor package.
- Pixel values in JSON (integers) are converted to fp16 on load: `x16 = x * 16`.
- Package: `internal/engine/combat/weapon/`
- JSON loader: wherever `climber.json` shooting skill config is parsed (likely `internal/game/` entity factory).
- Depends on US-034 (base offset infrastructure already in place).

## Notes

- No new VFX position concerns: muzzle flash already inherits the spawn position from `Fire()` (US-030), so it benefits automatically.
- This story is purely additive — no breaking changes to existing weapon callers that don't pass per-state offsets.
