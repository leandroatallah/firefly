# Roadmap

Cross-story sequencing. **Story IDs do not imply execution order — read this file for current priority.**

Update this file when a story moves to `active/` or `done/`, or when sequencing changes.

---

## Active Sequence

1. `070-render-offset-facing-kit-wiring` — per-facing `x_flipped` in `SpriteOffset`; wire `ApplyRenderOffsets` into platformer kit constructor.

---

## Dependencies

```
058 (done)
  └─ 061-altitude-jump-ground-detection (done)
       ├─ 062-depth-aware-collision (done)
       │    └─ 069-depth-lane-body-impl (done)
       ├─ 063-shadow-component (done)
       ├─ 064-beatemup-footprint-rect (done)
       ├─ 065-beatemup-jump-skill (done)
       └─ 066-beatemup-airborne-state-transitions (done)

067-actor-json-hitbox-active-frames (done)
068-actor-json-sprite-render-offset (done)
  └─ 070-render-offset-facing-kit-wiring (backlog)
```

---

## Parking Lot

Future stories not yet written. One line each.

- `060-html-report-generator` — in backlog.

---

## Recently Completed

- `068-actor-json-sprite-render-offset` — per-state sprite render offset in actor JSON.
- `067-actor-json-hitbox-active-frames` — hitbox active-frames override from JSON.
- `069-depth-lane-body-impl` — `DepthLaneBody` on `ObstacleRect` and `BeatEmUpCharacter`; correct depth-lane collision for airborne bodies.
- `066-beatemup-airborne-state-transitions` — airborne state transitions for `BeatEmUpCharacter`.
- `065-beatemup-jump-skill` — `BeatEmUpJumpSkill` kit skill: altitude-axis jump with coyote time, buffering, jump-cut.
- `064-beatemup-footprint-rect` — JSON `footprint_rect` schema + `Footprint()` accessor.
- `063-shadow-component` — shadow renders at `(X, GroundY)` as visual anchor when airborne.
- `062-depth-aware-collision` — depth-lane gate in `HasCollision`.
- `061-altitude-jump-ground-detection` — altitude gravity + landing in `BeatEmUpMovementModel`.
- `059-thin-game-phase-scenes` — consolidate scene `OnStart` logic into kit layer (pure refactor).
- `058-wire-beatemup-movement` — wire `EightDirectionalMovementSkill` + `BeatEmUpMovementModel` into `BeatEmUpCharacter`; decouple engine from `*PlatformMovementModel`.
- `057-beatemup-movement-model` — `BeatEmUpMovementModel` with 8-way floor movement, friction, speed cap; altitude-silent.
- `056-eight-dir-movement-skill` — `EightDirectionalMovementSkill` kit skill.
- `055-kit-genre-phase-scenes` — genre routing, beat-em-up scene shell, altitude draw sort. Closed by 059.
- `054-` and earlier — see `.agents/work/done/`.
