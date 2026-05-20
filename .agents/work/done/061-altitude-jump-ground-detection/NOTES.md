# NOTES — 061-altitude-jump-ground-detection

## Investigation Findings

- Contracts already in place from story 053:
  - `body.Body.Altitude()` returns pixel int via `fp16.From16(altitude16)`.
  - `body.Body.SetAltitude(alt int)` stores `fp16.To16(alt)`.
  - `body.Movable.VAltitude16()` / `SetVAltitude16(v16 int)` on `MovableBody`.
  - `body.Movable.AccelerationAltitude()` / `SetAccelerationAltitude` (not used this story).
- `BeatEmUpMovementModel` (story 057, file `movement_model_beatemup.go`) currently never touches altitude — fully altitude-silent. Confirms "passive" baseline.
- `PlatformMovementModel.handleGravity` (file `movement_model_platform.go`) is the structural template for gravity semantics:
  ```
  if vy16 < 0:  vy16 += UpwardGravity
  else:         vy16 += DownwardGravity
  if vy16 > maxFallSpeed: vy16 = maxFallSpeed
  ```
  We mirror this on the altitude axis but **without** the fall-speed cap (story explicitly defers the cap).
- `config.Get().Physics.UpwardGravity` and `DownwardGravity` already exist and are used by tests (see `TestPlatformMovementModel_handleGravity`). No config additions needed.

## Design Choices

- **Where to insert**: at the very end of `Update`, AFTER the final `b.SetVelocity(vx16, vy16)`. This keeps the existing 2D plane (X/Y) logic intact and untouched, making the altitude axis a clean additive concern. It also means the existing `if b.Freeze()` early-return at the top of `Update` already covers AC-6 — no separate guard required.
- **Pixel-level integration vs fp16 integration**: the user story explicitly mandates `alt += fp16.From16(VAltitude16)`. This introduces some quantization (sub-pixel velocities truncate to 0 per frame), but matches the story contract literally and uses the existing `SetAltitude(int)` pixel API. A future refactor could move integration into `altitude16` space; out of scope here.
- **Grounded short-circuit**: `if alt <= 0 && vAlt16 >= 0` skips the entire block (no gravity, no integration, no clamp). This is what guarantees AC-7 (2D-only bodies never get an unsolicited altitude update) — for a body that has never set altitude, both values are 0 and the branch returns immediately, leaving everything 0.
- **Landing clamp idempotency**: when `alt <= 0` AND `vAlt16 >= 0` we don't even enter the gravity block, so the clamp can't be re-applied indefinitely. The clamp inside the block is only reached when airborne (entered with `alt > 0` OR `vAlt16 < 0`).
- **Jump as external impulse (AC-5)**: model owns no input. Callers (state machine, input handler) call `b.SetVAltitude16(-N)` before the next Update tick. This mirrors how horizontal acceleration is set externally and consumed passively.

## Risks

- **Quantization**: `fp16.From16(VAltitude16)` truncates. A small jump impulse (e.g. `-fp16.To16(0.5)` → -32768) may integrate to `From16(-32766) = -1` then `From16(-32764) = 0`, causing very short jumps to round-trip in 1 frame. Acceptable for now; flagged for future story.
- **Body type bridging**: `body.MovableCollidable` exposes both `Altitude()` and `VAltitude16()` — confirmed via interface composition (`Movable` + `Body`). No type assertions needed.
- **2D regression coverage**: T-061-6 must run multiple frames AND also exercise X velocity to prove the 2D path is unaffected. Already present in existing diagonal/friction tests; the new test focuses on the altitude invariants only.

## Future Hooks (NOT in this story)

- `MaxAltitudeFallSpeed` config field + clamp (mirrors `MaxFallSpeed`).
- Integrate `AccelerationAltitude()` into `VAltitude16` similarly to how `Acceleration()` feeds `Velocity()`.
- State-machine wiring: emit `Jumping`/`Falling`/`Landed` events on transitions.
- Z-axis collision (landing on raised platforms with altitude > 0).
