# NOTES — 064-beatemup-footprint-rect

## Design Choices

- **`CollisionPosition()` override instead of an engine-level "footprint" concept.** The cleanest seam is method shadowing on `*BeatEmUpCharacter`. Engine collision (`space.HasCollision`, `ApplyValidPosition` → `ResolveCollisions`) already dispatches through the `body.Collidable` interface, so a beatemup-only override automatically scopes the behavior to this genre without touching `engine/`. This keeps schemas additive (a single optional pointer) and isolates AC-6 by construction.
- **Footprint stored as local rect, world-offset at call time.** Mirrors how `SetCollisionBodies` originally treats `CollisionRects`, but we deliberately keep this map on `BeatEmUpCharacter` (not as a real `CollidableBody`) because feet rects do not need ID, blocking flags, or `OnTouch` hooks — they only feed AABB. Less surface area than wiring them as full bodies.
- **Zero-size footprint = absent.** Treating `(0,0,0,0)` as a no-op is simpler than a "skip axis" rule, avoids ambiguity in `ApplyValidPosition`'s per-axis loop, and matches authoring intent (an empty JSON object likely means "I didn't fill this in").
- **No mirroring on `FaceDirectionLeft`.** Today's `CollisionRects` are also not mirrored — they are authored from sprite-frame origin and the mirror is applied at draw time only. Mirroring `footprint_rect` would silently desync from `collision_rect` and break existing levels. AC-8 asks for *parity* with the existing rect, which is "no mirror in the rect data".

## Risks & Quirks

- Method shadowing relies on callers dispatching through `*BeatEmUpCharacter` (or any interface that selects the override). The `space.Space` stores bodies as `body.Collidable`. When `NewBeatEmUpCharacter` registers itself in space, the registered interface value must hold `*BeatEmUpCharacter`, not the embedded `*CollidableBody`. Verify in `T-I1` that `space.HasCollision` observes the override path (case A returning `false` is the smoke test).
- Story 062 (`062-depth-aware-collision`) introduces `DepthLaneBody` and a depth gate inside `HasCollision`. The footprint override composes additively: bbox uses footprint, then depth-lane gate filters further. If 062 lands first, no integration drift expected. If 064 lands first, the depth gate path is simply unreached for beatemup actors that don't implement `DepthLaneBody`.
- `Footprint()` allocates `image.Rectangle` per call (cheap value type) but `CollisionPosition()` allocates a 1-elem slice on every collision check. Acceptable for movement (a handful per frame per actor); revisit if profiling shows pressure.
- Fallback path uses `unionRects` over `CollidableBody.CollisionPosition()`. If a state has multiple non-overlapping collision rects, the union may be larger than any single one — but this only happens when `footprint_rect` is absent, i.e., legacy behavior.

## Future

- [ ] Mirror `footprint_rect` on face-left if authoring data requires it (out of scope; current AC mandates parity with `collision_rect`).
- [ ] Genre-agnostic `Footprintable` contract in `internal/engine/contracts/body/` if topdown or others ever need feet-only collision.
- [ ] Debug overlay: draw footprint rect in a distinct color in `DrawCollisionBox`.
- [ ] Story 062 integration test: depth-lane gate + footprint composed.
- [ ] Per-axis footprint behavior for narrow corridors (Y-only footprint).

## Playtest

**Standalone:** No — visible behavior requires (a) authored `footprint_rect` entries in an existing beatemup actor JSON and (b) an obstacle whose top band overlaps the actor's head but not feet. The repo does not yet ship such a configured level. Verified via unit tests (T-F1..T-F8, T-I1) only.

To exercise manually after a level is authored: `go run cmd/game/main.go`, enter a beatemup scene, walk up to a wall — actor should stop when feet (not head) touch the base.
