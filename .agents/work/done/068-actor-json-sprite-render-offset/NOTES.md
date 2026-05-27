# NOTES — 068-actor-json-sprite-render-offset

## Design Choices

- **Storage on `Character`, not on the sprite map.** `sprites.Sprite` is shared across instances and is purely image+loop data. Per-state offsets live alongside other actor-instance config (state machine, footprints), keyed by `ActorStateEnum`, so they follow the same lifetime/wiring pattern as `buildFootprints` in `internal/kit/actors/beatemup/beatemup_character.go`.
- **Apply in `Character.UpdateImageOptions()` as the final translation.** The story said "kit-level sprite renderer"; in this codebase the actual draw-transform pipeline (flip, anchor, scale, world translate) is in engine's `Character`. Putting the offset there is the only way to satisfy AC-6 ("applied after all existing draw transformations"). The schema/render-time separation still holds: `schemas` does not depend on actors; physics/collision packages are untouched.
- **Builder helper in engine, called from kit.** `ApplyRenderOffsets` is genre-agnostic (offsets are universally per-state pixel nudges), so it lives in `internal/engine/entity/actors/builder` next to `BuildStateMap` and `SetCharacterBodies`. The beatemup kit invokes it explicitly; platformer/shooter kits can opt in later without an API change.
- **Facing-left does NOT mirror `X` (AC-9 decision).** The offset is applied AFTER `Scale(-1, 1)` and the post-flip recentering. Mirroring would force art authors to think about both poses; not mirroring matches what they see in their editor — "shift this artwork 4 px left on screen." Concrete: `cody-melee-0.png` is right-biased; a single `X:-4` corrects it regardless of facing direction, exactly as the author would tune by eye in-game. If a future story demands per-facing offsets, add `XFlipped int` (preserving current behavior when unset) rather than retroactively mirroring `X`.
- **No caching across frames.** Reading the map every `UpdateImageOptions()` call is a single map lookup; matches the edge case "no stale value from previous state."

## Risks & Quirks

- `{x:0,y:0}` produces a non-nil pointer but a zero `Translate`, so behavior is identical to nil — but `RenderOffset(state)` will return `ok=true` for explicit zeros. Tests must distinguish "registered with zero" from "not registered" by GeoM equality, not by `ok`.
- The translation row of `ebiten.GeoM` is accessed via `Element(0, 2)` / `Element(1, 2)` in tests. No public setter; tests must reset and recompute baselines per row.
- Large offsets are not clamped — an actor with `X:1000` will render off-screen but physics/collision remain valid. This is intentional (story Behavioral Edge Cases).
- The beatemup kit edit is a one-line addition; if a future story adds platformer/shooter actor JSONs with `render_offset`, they will silently no-op until their kit constructor also calls `ApplyRenderOffsets`. Easy to forget.

## Future

- [ ] Per-facing-direction offsets (`X`, `XFlipped`).
- [ ] Wire `ApplyRenderOffsets` into platformer and shooter kit constructors.
- [ ] Optional sub-pixel offsets (`float64`) if anim authors need finer control.

## Playtest

**Standalone:** Yes — add `"render_offset": {"x": -4, "y": 0}` to the `melee` asset in `assets/actors/player_cody.json` (or whichever beatemup actor JSON drives `cody-melee-0.png`), then `go run cmd/game/main.go`, trigger the melee attack while idle, and observe that the visible character no longer jumps leftward at the idle→melee transition. Body hitboxes (debug overlay) are unchanged.
