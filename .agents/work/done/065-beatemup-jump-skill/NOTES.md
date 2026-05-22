# NOTES — 065-beatemup-jump-skill

## Design Choices

- **Ground predicate is `b.Altitude() <= 0`, not `model.OnGround()`.** The story defines grounding on the altitude axis (story 061). `BeatEmUpMovementModel` does not implement `Grounded`; using the body's altitude keeps the model passive as 061 AC-5 mandates.
- **`IsInputBlocked` added to `BeatEmUpMovementModel`.** Mirrors `PlatformMovementModel` and satisfies the `movement.InputBlocker` interface. Tiny scope creep but unavoidable to honor AC-8 cleanly.
- **Factory branches on `cfg.Movement.Mode == MovementModeEightDir`** rather than introducing a new `cfg.Jump.Variant` field. Keeps config schema stable; mode is the canonical signal of the physics flavor.
- **Apex detection uses `VAltitude16() >= 0`** instead of `IsGoingUp()`. `IsGoingUp` is defined on the vy16 axis and is unrelated to altitude motion.
- **`Update` is the only place buffer fires.** `HandleInput` only seeds the buffer; this keeps the leading-edge press path simple and matches `JumpSkill` semantics.
- **`Freeze()` check lives in the skill `Update`,** not in `SkillBase`. The story line "the base SkillBase.Update already exits early" is inaccurate — base is a no-op. We add an explicit guard so frozen actors do not bleed counters or get buffered re-fires.

## Risks & Quirks

- `int(float64(v) * 0.5)` truncates toward zero in Go; for `v=-320` this gives `-160` exactly, but odd inputs may drift one fp16 sub-pixel. Existing `JumpSkill_JumpCut` test accepts this; we follow the same convention.
- If a future story adds gravity-only-when-airborne logic that flips altitude sign sources, the `VAltitude16 >= 0` apex check may need re-evaluation.
- `input.CommandsReader()` is a global. Tests must reset commands between cases to avoid bleed.
- Factory regression: `JumpSkill` tests rely on the prior unconditional branch. Adding the mode switch must default the non-eight_dir paths to today's behavior.

## Future

- [ ] Move `cfg.Jump.Variant` into schema to decouple jump flavor from movement mode.
- [ ] Configurable jump-cut release window (e.g., only allow cut within first N frames).
- [ ] Hold-to-glide as a separate skill on top of altitude axis.

## Playtest

**Standalone:** No — beat-em-up movement is wired into the demo scene but no actor binds `BeatEmUpJumpSkill` yet. After this story:
1. Run `go test ./internal/kit/skills/... ./internal/engine/physics/movement/...` — all green.
2. Visible behavior requires a follow-up story to attach `BeatEmUpJumpSkill` to the player actor in the beat-em-up phase config (factory wiring is in place, but the active phase still uses the horizontal kit). Until then, verified via unit tests only.
