# USER STORY — 065-beatemup-jump-skill

**Branch:** `065-beatemup-jump-skill`
**Bounded Context:** Kit (`internal/kit/skills/`)
**Depends on:** 061 (done)

---

## Story

As a kit developer,
I want a `BeatEmUpJumpSkill` that triggers altitude-axis jumps on `BeatEmUpMovementModel`,
so that beat-em-up actors can jump using the same UX guarantees (coyote time, buffering, jump-cut) as the platform jump skill.

---

## Background

Story 061 made `BeatEmUpMovementModel` passively apply gravity and landing on the altitude axis. AC-5 of 061 states the model is passive — jump impulse must come from an external source. No skill currently sets `VAltitude16`. The existing `JumpSkill` type-asserts `*PlatformMovementModel` and cannot be reused.

"Rising" on the altitude axis means `VAltitude16 < 0` (upward = decreasing altitude). "Grounded" means `body.Altitude() <= 0`.

---

## Acceptance Criteria

- AC-1: `BeatEmUpJumpSkill` lives in `internal/kit/skills/beatemup_jump.go`; `HandleInput` and `Update` type-assert `*physicsmovement.BeatEmUpMovementModel` and return early (no-op) if the model is any other type.
- AC-2: On `input.CommandsReader().Jump` leading-edge while grounded (`b.Altitude() <= 0`) or coyote counter > 0, the skill calls `b.SetVAltitude16(-force)` where `force = int(float64(cfg.Physics.JumpForce) * b.JumpForceMultiplier())`; `jumpCutPending` is set to true and `OnJump` callback is invoked if non-nil.
- AC-3: While airborne without a buffered jump, a new jump press is ignored (no double-jump).
- AC-4: Jump press while airborne sets `jumpBufferCounter = cfg.Physics.JumpBufferFrames`; when the actor next lands (`Altitude() <= 0`), the buffered jump fires immediately.
- AC-5: Coyote time — each `Update` frame while `Altitude() <= 0` the counter resets to `cfg.Physics.CoyoteTimeFrames`; while airborne the counter decrements each frame; a jump issued while counter > 0 is treated as grounded.
- AC-6: Jump-cut — on leading-edge release of Jump while `jumpCutPending` and `VAltitude16 < 0`, the skill applies `b.SetVAltitude16(int(float64(b.VAltitude16()) * s.jumpCutMultiplier))`; `jumpCutPending` is cleared once `VAltitude16 >= 0`.
- AC-7: `SetJumpCutMultiplier(m)` clamps `m` to `(0.1, 1.0]`; 1.0 is the default (no cut).
- AC-8: `model.IsInputBlocked() == true` causes `HandleInput` to return early without consuming input.
- AC-9: `b.Freeze() == true` — the base `SkillBase.Update` already exits early; no additional altitude mutation occurs.
- AC-10: `NewBeatEmUpJumpSkill` returns a ready skill; `OnJump` callback field is exported.
- AC-11: `factory.FromConfig` instantiates `BeatEmUpJumpSkill` when `cfg.Jump != nil && cfg.Movement.Mode == schemas.MovementModeEightDir`; the existing `JumpSkill` is instantiated for all other movement modes (no regression).
- AC-12: Table-driven unit tests in `internal/kit/skills/beatemup_jump_test.go` cover: jump from grounded, no double-jump while airborne, coyote-time jump, jump-buffer fires on landing, jump-cut multiplier applied to `VAltitude16`, no-op when model is `*PlatformMovementModel`, input-blocked guard.

---

## Behavioral Edge Cases

- `JumpForce` yields `force <= 0` after multiplier: jump is skipped silently; no `SetVAltitude16` call made.
- Jump pressed and released within same frame: leading-edge fires the jump; trailing-edge fires jump-cut in the same `HandleInput` call — guard against this by checking `jumpCutPending` is true before applying cut.
- `jumpBufferCounter` reaching 0 before landing: buffer expires; no jump fires on next landing.
- Multiple consecutive grounded frames: coyote counter stays at max; no double-reset side effects.
- Body mid-air with `Freeze() == true`: `SkillBase.Update` exits; coyote and buffer counters do not advance that frame.
- `OnJump` callback panics: that is the caller's responsibility; the skill does not recover.
