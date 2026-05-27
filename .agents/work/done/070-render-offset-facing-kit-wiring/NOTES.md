# NOTES — 070-render-offset-facing-kit-wiring

## Design Choices

- **`XFlipped *int` (pointer) not `XFlipped int`.** Distinguishes "not set" (use `X`) from "explicit 0" (left-facing nudge is zero). A plain int would force a sentinel like `math.MinInt`, which is fragile and surprises authors.
- **Resolve X at draw time, not at registration.** `SetRenderOffset` stores both values; `UpdateImageOptions` picks one each frame. Facing can flip mid-state (accel sign change), and the offset must follow immediately — caching the resolved X on `SetFaceDirection` would require invalidation hooks we don't have.
- **Reuse the local `fDirection` already resolved in `UpdateImageOptions`** (line ~328-336) rather than calling `c.FaceDirection()` again. The local already accounts for acceleration override, so the offset and the sprite flip stay in lockstep.
- **Storage type changed from `map[State]image.Point` to `map[State]renderOffset`.** This breaks the story-068 `SetRenderOffset(state, dx, dy)` signature; chose to migrate callers (`builder.ApplyRenderOffsets`, tests) rather than add a parallel method, because 068 is a single-call-site feature and the broader signature is the long-term API.
- **`RenderOffset(state)` kept returning `image.Point`** — resolved against current facing. Lets 068 tests survive with a `nil` flipped value and gives the platformer test a one-line assertion (`p == Pt(-2,0)`).
- **Platformer call site mirrors beatemup exactly** (line immediately after `actors.NewCharacter`, before the kit struct literal). Predictable for future kits.

## Risks & Quirks

- `SetRenderOffset` signature change is a breaking source-level change inside `internal/engine/`. Touch points are limited (one builder, the 068 unit tests) but any external caller would need updating.
- `XFlipped == &0` vs `nil`: tests must use a real `int` variable address; literal `&0` is not allowed in Go.
- `UpdateImageOptions`'s `fDirection` can be overridden by `accX` sign every frame — tests must zero acceleration before calling `SetFaceDirection` to assert a specific facing path.
- T-P1 needs a real sprite asset on the embedded `fs.FS` (or a fake `fs.FS`) because `sprites.GetSpritesFromAssets` decodes images. Check what 068's beatemup tests use; reuse the same fixture.
- Shooter kit (`shooter_character.go`) is intentionally untouched. If a future story turns it into a `SpriteData`-consuming constructor, it must also call `ApplyRenderOffsets`.

## Future

- [ ] Tune `cody-melee-0.png` via `x_flipped` in `assets/entities/player/cody.json`.
- [ ] Float64 sub-pixel offsets for fine-grained art tuning.
- [ ] Per-frame (not per-state) render offsets if animation authors need it.
- [ ] Wire `ApplyRenderOffsets` into any future genre kit constructors.
