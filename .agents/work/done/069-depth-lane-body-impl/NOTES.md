# NOTES — 069-depth-lane-body-impl

## Design Choices

- **`ObstacleRect.LaneHalfWidth = max(height, DefaultLaneHalfWidth)`** instead of `height/2`. The `HasCollision` gate compares raw `GroundY` values (character bottom Y vs obstacle bottom Y) and accepts when `|diff| <= max(halfA, halfB)`. Using the obstacle's full height as half-width means a character whose feet line falls anywhere from the obstacle's top edge to its bottom edge passes the gate — which matches the intuitive "the wall occupies a depth band equal to its rendered height" model. Half-height would have made the lane too narrow on the top side.
- **`BeatEmUpCharacter.GroundY` uses `y16 >> 16`** rather than `Position().Min.Y`. `Position()` subtracts altitude; we need pre-altitude depth. `y16` is the raw fp16 depth axis.
- **`DepthLaneBody` lives on the parent body, not on collision shapes.** `ResolveCollisions` hands `HasCollision` the parent (Owner) — that's where the interface must be satisfied. `CodyPlayer` and friends pick up the methods for free via embedded `*BeatEmUpCharacter` (no per-game re-declaration). PlatformerCharacter does not embed BeatEmUpCharacter, so it stays out — regression-safe.
- **Removing Block 1** (the zero-altitude wrap) is the whole point: with the gate now reliable in both directions (player-initiated and projectile-initiated checks), the workaround is redundant and was only ever a half-fix. The PLAN doc (`PLAN_airborne-collision-split.md`) called this Option B and originally rejected it on the assumption "depth gate separates depth lanes, not altitudes" — but the real-world wall heights make the lane wide enough to also cover the airborne case, so Option B is in fact sufficient and cleaner than Option C.

## Risks & Quirks

- **Wide obstacles eat the depth axis.** A 64-px-tall wall sets `LaneHalfWidth = 64`. Two such walls 60 px apart along Y will both collide with a character standing between them. In practice obstacles are arranged along the bottom of the play area so this is acceptable; flag if level layouts grow more vertical.
- **Owner re-declaration**: if any game-layer subclass overrides `GetPosition16` or wraps the body differently, double-check method promotion still surfaces `GroundY/LaneHalfWidth`. Today (`CodyPlayer` embeds `*BeatEmUpCharacter`) it works.
- **PLAN doc is now stale.** Migration note in the user story calls for updating `PLAN_airborne-collision-split.md`: Option B chosen, Block 1 removed, Option C retired. Gatekeeper or the developer should land that doc edit alongside the merge.

## Future

- [ ] Per-state lane width on `BeatEmUpCharacter` (forgiving hitboxes during attacks).
- [ ] Depth-lane debug visualizer (draw lane bands over obstacles).
- [ ] Story 064 (`beatemup-footprint-rect`) — orthogonal but will refine character footprint extents.
- [ ] Consider whether `MaxLaneHalfWidth` clamp is needed once level designers add tall background props.

## Playtest

**Standalone:** Yes — run `go run cmd/game/main.go`, enter a beat-em-up phase that contains background obstacles at a different depth than the player walkway. Jump near such an obstacle: the player should no longer be blocked in mid-air by it. Then walk into an obstacle at the player's own depth: the player must still be blocked (when shapes overlap on the depth axis). Compare against current `main` to observe the false-block regression that this story closes.
