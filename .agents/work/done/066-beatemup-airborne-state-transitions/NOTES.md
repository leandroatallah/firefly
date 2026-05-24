# NOTES — 066 Beat-Em-Up Airborne State Transitions

## Design Choices

- **Altitude axis vs `IsGoingUp`/`IsFalling`.** `*MovableBody.IsGoingUp()`/`IsFalling()` read `vy16` (ground-plane Y), which is **not** the altitude in beat-em-up. We deliberately reference `VAltitude16()` and `Altitude()` directly inside `beatemupMovementTransitions`. The platformer handler keeps using `IsGoingUp/IsFalling` because in that genre Y *is* altitude.
- **Apex branch first.** The `state == Jumping && VAltitude16() >= 0 && airborne → Falling` branch is intentionally checked before the generic `VAltitude16() < 0` branch. This makes the apex frame (`vAlt16==0`) deterministic: once `Jumping`, the next non-ascending frame moves to `Falling`.
- **Landing lock via `IsAnimationFinished()`.** Same pattern as the platformer handler; chosen for visual consistency. Reads from the current state instance (`LandingState` already implements the animation-complete check via the base state machinery).
- **Airborne guard against Walking.** `airborne` is defined as `Altitude() > 0` rather than `state ∈ {Jumping, Falling}` so the guard still holds if a contributor exotic state interleaves; the explicit case ordering above makes AC-5 hold regardless.
- **No new contracts.** The handler is a free function reading from `*actors.Character`. Nothing crosses bounded contexts.

## Risks & Quirks

- `IsAnimationFinished()` semantics depend on the registered state implementation. If `LandingState` ever loops, the lock becomes permanent. Risk is low because `LandingState` is a "play once then surrender" state by convention, but TDD should pin behaviour with a controllable fake.
- Buffered jump on landing frame (AC edge): the handler enters `Landing` first; `BeatEmUpJumpSkill` re-fires next frame. If both happen in the same frame in some future change, expect a Landing → Jumping transition on the very next handler call (R1 supersedes R4 only when `vAlt16 < 0`).
- Same-frame land (jump impulse with `alt==0`): a negative `vAlt16` is applied by the skill before the handler runs in the same frame, so R1 fires and we never spuriously enter Landing.
- `Hurted` is set by contributors *before* the movement handler runs; no extra guard needed.

## Future

- [ ] Crouch / duck airborne variants (Falling+Down).
- [ ] Air-attack states (Jumping+Attack) — needs a contributor, not this handler.
- [ ] Hard-landing variant when descent speed exceeds a threshold.
- [ ] Wire `OnJump`/`OnLand` VFX hooks for beat-em-up actors (parity with platformer).

## Playtest

**Standalone:** Yes — run `go run cmd/game/main.go`, load a beat-em-up phase using Cody. Press jump; expect three distinct sprite animations: jump (ascent), fall (descent), land (touchdown). Walking + jump combo should keep jump/fall sprites in the air and resume walk on landing once the land animation completes.
