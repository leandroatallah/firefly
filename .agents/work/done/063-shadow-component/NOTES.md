# NOTES — 063-shadow-component

## Design Choices

- **Package over method.** Shadow is a package-level utility in `internal/kit/render/shadow/`, not a method on the actor. Reasons: (1) keeps actor surface lean — every beat-em-up actor already exposes `Altitude()`, `GetPositionMin()`, `GetShape()`; (2) lets the scene control draw order centrally (shadows always render before sprites, in a single pass); (3) avoids touching `BeatEmUpActorEntity` contract for a pure render concern.
- **No new contract.** `AltitudeBody` is a structural interface local to the shadow package. Adding it under `internal/engine/contracts/` would force `Mock Generator` work for a leaf render utility with no behavioural variation. Mock Generator is skipped.
- **Single sink seam (`ovalDrawerFn`).** Lets tests assert call counts without a GPU and without faking the camera. Production keeps a thin `drawOval` calling `cam.Draw` of a 1×1 white image scaled into an oval-approximated rect (or a small cached oval `*ebiten.Image`). Implementer chooses.
- **Linear falloff.** `ScaleFor` is a clamped linear interpolation. Simpler than quadratic and matches the genre cue ("higher = smaller shadow") without needing tuning hooks in this story.
- **Foot-line for CenterY.** `CenterY = y + H` puts the oval at the body's bottom edge (the ground-plane Y). This matches how the actor is positioned: the sprite is offset *up* by `Altitude` for rendering, but the body's `Y` itself is the ground-plane top — adding `H` gives the foot midpoint.

## Risks & Quirks

- **Oval rendering primitive.** Ebitengine's `vector` package draws directly to a target with no GeoM hook, so we cannot route an oval through `cam.Draw`. Implementer must either (a) generate a small oval `*ebiten.Image` and `cam.Draw` it with non-uniform scale, or (b) compute camera-relative coords and call `vector.DrawFilledCircle` after translating. Option (a) is preferred for consistency with the rest of the scene draw path.
- **Allocation pressure.** Constructing a new oval image per shadow per frame would allocate. Cache one 1×1 white image at package init; let GeoM scale it to oval bounds.
- **Z-order regression.** Adding `shadow.DrawAll` before `SortByGroundYAltitude(...)` means shadows are drawn in raw `Bodies()` order. This is fine because all shadows are translucent black with no overlap concerns; if two shadows overlap, alpha just adds — no visual artefact worth fixing here.
- **Beat-em-up only.** Platformer scene must remain untouched. AC-7 is enforced by code review — there is no shared draw helper to accidentally hit both.

## Future

- [ ] Per-actor shadow toggle (e.g., bosses or projectiles opting out).
- [ ] Configurable falloff/min-scale via beat-em-up scene options.
- [ ] Soft-edge oval texture instead of solid alpha for nicer visuals.
- [ ] Shadow color tied to floor brightness (per-tile sample).

## Playtest

**Standalone:** No — requires a beat-em-up phase with an airborne actor. Until the beat-em-up jump skill story (queued in backlog) lands, the only way to see a shadow in-game is to manually `SetAltitude(N)` on the player in a dev build. Verified via unit tests:

- `go test ./internal/kit/render/shadow/...`
- `go test ./internal/kit/scenes/phases/beatemup/...`

Once a jump input is wired up: `go run cmd/game/main.go`, enter a beat-em-up phase, jump, and observe the oval shadow staying on the ground while the sprite rises.
