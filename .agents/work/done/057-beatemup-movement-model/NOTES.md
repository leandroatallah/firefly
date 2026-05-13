# NOTES — 057-beatemup-movement-model

Human-readable companion to `SPEC.md`. Captures investigation findings, rationale, and risks.

---

## Investigation Findings

### Code structure of existing models

- `TopDownMovementModel.Update` (movement_model_topdown.go) is the closest template: freeze guard → apply X/Y velocity to position with collision → `clampToPlayArea` → integrate acceleration (with `smoothDiagonalMovement`) → 2D speed cap → reset acceleration → friction on both axes.
- `PlatformMovementModel.Update` carries gravity, `onGround`, dash, and ground-stick logic — none of which apply to beat-em-up ground-plane.
- `clampToPlayArea` already clamps to either tilemap dimensions (when `TilemapDimensionsProvider` is set on the space) or screen size — exactly the "walkable strip via Tiled obstacle tiles + tilemap edge clamp" model the story prescribes. No changes needed.

### Helpers reused

- `smoothDiagonalMovement(accX, accY)` — normalizes diagonal acceleration by `1/sqrt(2)`.
- `increaseVelocity` / `reduceVelocity` — single-axis velocity and friction.
- `clampToPlayArea(body, space)` — tilemap/screen-edge clamp. Returns `isOnGround` (we ignore — it is a platformer concept).
- `body.ApplyValidPosition(distance16, isXAxis, space)` — collision-aware position application.

### Bug noticed in TopDown (NOT fixed in this story)

`TopDownMovementModel.Update` line 63 computes `velSq := int64(vx16) + int64(vy16)*int64(vy16)` — missing `*int64(vx16)` on the first term. The beat-em-up spec writes the correct form `vx16*vx16 + vy16*vy16`. This is intentional and noted here; correcting TopDown is out of scope.

### Passive model — no InputHandler

The story resolution mandates the model is passive. `playerMovementBlocker` is still accepted in the constructor so the `NewMovementModel` factory keeps a uniform signature across all three models, and a future input-aware variant or skill can use it. The field is presently unused in `Update`.

---

## Design Rationale

### Why keep `playerMovementBlocker` parameter when unused

Two reasons:
1. Factory `NewMovementModel(model, blocker)` already passes a single positional arg to each constructor. Dropping it for BeatEmUp would force a special-case factory.
2. Cheap forward compatibility: when 056 lands and an Eight-Directional skill needs blocker awareness, it can be wired without an API break.

If reviewers prefer a strict minimum API, the parameter can be dropped with a corresponding factory tweak — flag for grilling.

### Altitude-gravity integration point [USER_STORY open item]

The model does not read or write any altitude field. It only operates on `Velocity()` (x,y) and `Acceleration()` (x,y) — exactly the contract `MovableCollidable` exposes. A future jump/altitude skill will own altitude state on the body itself (e.g., a separate `Altitude()/SetAltitude()` pair) and run its own integration. The movement model never needs to know about altitude.

This is the cleanest possible integration point: do nothing. Hardcoding `altitude = 0` would be the wrong move (it would clobber any altitude that a jump skill is mid-applying). By being silent on altitude, BeatEmUp composes with any future altitude system.

### Why the 2D speed cap (not per-axis clamp)

Per-axis clamping permits diagonal speed = sqrt(2) × cardinal speed. The 2D vector cap mirrors TopDown's approach and is the standard fix.

### Friction on both axes

Symmetric to TopDown. The Y axis is ground-plane depth, not altitude, so applying horizontal-style friction to Y is semantically correct (deceleration when no input).

---

## Risks

| Risk | Mitigation |
|---|---|
| Test T-BE5 (diagonal cap) becomes flaky if exact speed boundary depends on per-frame friction noise. | Use a generous tolerance (`× 1.05`) and run for fixed 60 frames; assert magnitude, not exact equality. |
| `ApplyValidPosition` semantics may surprise — for very large velocity, it can teleport into the obstacle if collision query misses. | Use velocity `fp16.To16(20)` over a 20-pixel gap so the body firmly overlaps the obstacle; the collision resolver clips into a valid position. Asserts use `<` not `==`. |
| Forgetting to update `MovementModelEnum.String()` map causes silent "" return. | Spec explicitly lists the map entry. |
| Factory caller (story 058) may pass `nil` blocker. | Already permitted — `NewTopDownMovementModel(nil)` is exercised in existing tests. |

---

## Out-of-Scope

- Wiring into `beatemup_phase` scene (story 058).
- `EightDirectionalMovementSkill` (story 056).
- Jump / altitude axis (future).
- Fixing the `vx16 + vy16*vy16` typo in TopDown's speed-cap (separate cleanup).
- Camera bounds, playfield strip configuration — handled by Tiled obstacle tiles per resolved constraint.

---

## Open Questions for Reviewer

1. **Drop `playerMovementBlocker` from `NewBeatEmUpMovementModel`?** Current spec keeps it for factory uniformity (unused). If a strict-minimum API is preferred, drop both the field and add a factory special-case.
