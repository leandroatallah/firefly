# USER STORY — 061-altitude-jump-ground-detection

**Branch:** `061-altitude-jump-ground-detection`
**Bounded Context:** Physics (`internal/engine/physics/movement/`)

---

## Story

As an engine developer,
I want `BeatEmUpMovementModel` to apply altitude-axis gravity and support jump input via `VAltitude16`,
so that beat-em-up actors can leave the ground, arc through the air, and land — with `Altitude <= 0` as the grounded condition.

---

## Background

`BeatEmUpMovementModel` (story 057) is altitude-silent: it never touches `altitude16` or `vAltitude16`. The body contracts and fixed-point storage for altitude already exist (story 053). This story activates altitude physics inside the model.

`VAltitude16` is the altitude-axis velocity on `body.Movable`. Gravity is applied using existing `config.Physics.UpwardGravity` / `config.Physics.DownwardGravity` constants (same semantics as `PlatformMovementModel.handleGravity`, but on the altitude axis, not the Y axis). The entity is grounded when `body.Altitude() <= 0`; on landing, altitude is clamped to 0 and `VAltitude16` is zeroed.

**Depends on:** 058 (BeatEmUpMovementModel wired into BeatEmUpCharacter).

---

## Acceptance Criteria

- AC-1: Each `Update` frame, `UpwardGravity` is added to `VAltitude16` when `VAltitude16 < 0` (rising); `DownwardGravity` is added when `VAltitude16 >= 0` (falling) — mirroring `PlatformMovementModel.handleGravity` on the altitude axis.
- AC-2: Gravity is not applied when `body.Altitude() <= 0` and `VAltitude16 >= 0` (entity is at rest on the ground).
- AC-3: `Altitude` is updated each frame by `fp16.From16(VAltitude16)` and stored via `body.SetAltitude`.
- AC-4: When `body.Altitude() <= 0` after integration, altitude is clamped to 0 and `VAltitude16` is set to 0 (landing).
- AC-5: A jump is initiated externally by setting `VAltitude16` to a negative value (upward impulse); the model does not own jump input — it is passive on this axis too.
- AC-6: `body.Freeze() == true` causes `Update` to return early before any altitude mutation.
- AC-7: Existing 2D bodies that never set altitude remain unaffected: `Altitude()` stays 0, `VAltitude16()` stays 0, no regression.
- AC-8: Table-driven unit tests cover: gravity accumulates while airborne (rising arc), gravity accumulates while falling, landing clamps altitude and zeroes VAltitude16, freeze guard, zero-altitude body is unaffected (2D regression).

---

## Behavioral Edge Cases

- `VAltitude16` set to large negative value (high jump): gravity accumulates each frame until peak, then fall; landing clamp fires correctly.
- Multiple consecutive landing frames: clamping is idempotent — altitude stays 0, velocity stays 0.
- Body with `Freeze() == true` mid-air: no altitude change that frame; state preserved for next frame.
- `MaxFallSpeed` equivalent for altitude axis: not required in this story — leave a clean integration point; do not hardcode a cap.
