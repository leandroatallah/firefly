# NOTES — 058-wire-beatemup-movement

## Investigation findings

- `BeatEmUpCharacter` currently embeds only `*kitactors.MeleeCharacter` — it does **not** embed `*actors.Character` nor own a movement model. This is the gap. The story requires bringing it to parity with `PlatformerCharacter` (which embeds `*actors.Character` and sets `Owner`).
- The current beat-em-up scene already calls `s.Camera().SetFollowTarget(s.player)` and creates collision bodies with a (non-nil) endpoint factory. The new work is purely the `SetBounds(&tilemapRect)` call. The AC-5 phrasing "(nil)" likely refers to the bounds args being eliminated from the scene, not to the second argument of `CreateCollisionBodies`. Spec preserves the endpoint factory (SPIKE/CUTSCENE triggers depend on it).
- `physicsmovement.BeatEmUpMovementModel.Update` already does no gravity write and applies `ApplyValidPosition` on both axes, so obstacle-tile blocking is already correct at the model level. Wiring is the missing piece.
- `kitskills.FromConfig` already accepts `cfg.Movement.Enabled` — adding a `Mode` field is additive and backward-compatible.
- `EightDirectionalMovementSkill.HandleInput` already handles `Immobile` (zeroes velocity) and `IsInputBlocked`. No skill changes required for the immobile edge case.
- Camera `SetBounds` is already exercised in platformer/screen_flipper code and tested in `camera_test.go`; clamping behavior is engine-tested, so SPEC only requires verifying it is called.

## Rationale for design choices

- **Model owned at construction, not via `builder.ApplyPlatformerPhysics`**: follows platformer precedent. The platformer's `ApplyPlatformerPhysics` exists because old call sites built the character without a model. The story explicitly says "owned at construction" — so the constructor sets the model directly, and no `ApplyBeatEmUpPhysics` helper is introduced (avoids a parallel-but-different code path).
- **Mode discriminator on `MovementConfig` (not a new `MovementMode` block)**: minimal schema change, preserves all existing JSON. Empty/missing = backward-compatible.
- **No `ApplyBeatEmUpPhysics` builder helper**: explicitly out — model ownership is character-internal per the resolved constraints.

## Risks / mitigations

- Risk: existing `CodyPlayer` construction uses the old zero-arg `NewBeatEmUpCharacter()`. Mitigation: search and migrate all call sites in the same PR; failing build catches misses.
- Risk: cody config JSON missing the `mode` key — falls back to horizontal, breaking AC-6 silently in runtime. Mitigation: T-B4 enforces eight-dir behavior at the character unit level; integration is verified by adding `"mode": "eight_dir"` to the asset config in the same change set.
- Risk: `Camera().SetBounds(&tilemapRect)` clashes with future `screen_flipper` usage. Mitigation: out of scope; beat-em-up scene does not use screen flipper.

## Out of scope

- Cross-phase HP/inventory preservation (flagged in story).
- PlayArea object-layer support for configurable walkable strip.
- `TopDownMovementModel` wiring; only beat-em-up here.
- Refactoring `builder.ApplyPlatformerPhysics` to a generic helper.
