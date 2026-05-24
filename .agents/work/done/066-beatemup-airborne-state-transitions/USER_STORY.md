# Story 066 — Beat-Em-Up Airborne State Transitions

**Branch:** `066-beatemup-airborne-state-transitions`
**Bounded Context:** Kit

## Story

As a beat-em-up actor, I want my state machine to reflect Jumping, Falling, and Landing states during altitude-axis flight, so that the correct sprite animation plays while I am airborne.

## Context

`BeatEmUpCharacter` registers `beatemupMovementTransitions` as its `MovementTransitionHandler`. That handler only knows about `Walking` and `Idle` (ground-plane velocity). It has no awareness of the altitude axis (`VAltitude16()` / `Altitude()`), so actors stay in `Idle`/`Walking` while airborne. Sprite assets for `jump`, `fall`, and `land` exist in `cody.json` but are never reached.

The platformer genre layer solves the equivalent problem in `platformerMovementTransitions` (`internal/kit/actors/platformer/platformer.go`). The beat-em-up layer needs an analogous extension to `beatemupMovementTransitions` in `internal/kit/actors/beatemup/beatemup_character.go`.

## Acceptance Criteria

- AC-1: When `VAltitude16() < 0` (ascending), `beatemupMovementTransitions` sets state to `Jumping`, regardless of ground-plane velocity.
- AC-2: When `VAltitude16() > 0 && Altitude() > 0` (descending, airborne), the handler sets state to `Falling`.
- AC-3: When `Altitude()` returns to `0` after the actor was in `Falling`, the handler sets state to `Landing`.
- AC-4: While in `Landing`, the handler does not leave that state until `IsAnimationFinished()` returns `true`, then transitions to `Idle` or `Walking` (matching ground-plane movement).
- AC-5: While in `Jumping` or `Falling`, the handler does not apply `Walking`/`Idle` transitions.
- AC-6: While in `Jumping`, if the apex is passed (`VAltitude16() >= 0`) and the actor is still airborne, the handler transitions to `Falling`.
- AC-7: All new branching in the handler is covered by table-driven unit tests in `internal/kit/actors/beatemup/`.

## Edge Cases

- Actor lands on the same frame the jump started (altitude never exceeds 0) — must not enter Landing spuriously.
- `Hurted` state during flight — `pollStateContributors` already short-circuits; the movement handler must not override an active `Hurted` state (it runs after contributors, so this is satisfied by the existing `handleState` call order).
- Buffered jump fires on the landing frame (`BeatEmUpJumpSkill`) — `Landing` state should be entered first; the skill re-triggers jump on the next frame.
- Ground-plane movement while airborne — `vx`/`vy` changes must not cause a `Walking` transition during flight.
