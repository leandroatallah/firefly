# Epic: Beat-em-up Mechanics

## Goal

Extend the current 2D platformer engine into a hybrid 2.5D beat-em-up system while maintaining 100% compatibility with existing 2D side-scrolling levels.

The coordinate model: **X** = horizontal, **Y** = ground/depth, **Altitude** = height above ground. Screen mapping: `ScreenX = X`, `ScreenY = Y - Altitude`. Existing 2D levels are unaffected because Altitude is always 0.

## Acceptance Criteria

- Player moves in 8 directions on the floor (X, Y axes)
- Player can jump and land using the Altitude axis (gravity + ground detection)
- Depth-aware collision: entities only collide when depth difference is within a lane threshold
- Shadow renders at `(X, GroundY)` as a visual anchor when airborne
- Footprint rect drives collision shape for beat-em-up entities
- All existing 2D levels pass regression tests

## Child Stories

| Story | Status | Notes |
|---|---|---|
| `053-altitude-engine-foundation` | done | Altitude contracts + physics body + Z-sort |
| `056-eight-dir-movement-skill` | done | `EightDirectionalMovementSkill` kit skill |
| `057-beatemup-movement-model` | done | `BeatEmUpMovementModel` — floor movement, altitude-silent |
| `058-wire-beatemup-movement` | done | Wire movement into `BeatEmUpCharacter`; decouple engine from platform model |
| `061-altitude-jump-ground-detection` | done | Altitude gravity + landing in `BeatEmUpMovementModel` |
| `062-depth-aware-collision` | done | Depth-lane gate in `HasCollision` |
| `063-shadow-component` | done | Shadow at `(X, GroundY)` |
| `064-beatemup-footprint-rect` | done | JSON `footprint_rect` schema + `Footprint()` accessor |
| `065-beatemup-jump-skill` | done | `BeatEmUpJumpSkill`: altitude jump with coyote time, buffering, jump-cut |
| `066-beatemup-airborne-state-transitions` | done | Airborne state transitions for `BeatEmUpCharacter` |
| `069-depth-lane-body-impl` | done | `DepthLaneBody` on `ObstacleRect` and `BeatEmUpCharacter` |
| `067-actor-json-hitbox-active-frames` | done | Hitbox active-frame data in actor JSON |
| `068-actor-json-sprite-render-offset` | done | Sprite render offset in actor JSON |
| `070-render-offset-facing-kit-wiring` | backlog | Per-facing `x_flipped` in `SpriteOffset`; wire `ApplyRenderOffsets` into platformer kit |

## Out of Scope

- Full scrum/agile tooling — workflow stays file-first
- Networked multiplayer
- Combo/chain combat system (separate epic)
