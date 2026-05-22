# Roadmap

Cross-story sequencing. **Story IDs do not imply execution order — read this file for current priority.**

Update this file when a story moves to `active/` or `done/`, or when sequencing changes.

---

## Active Sequence

| Order | Story                                | Depends on          | Notes                                                                                                             |
| ----- | ------------------------------------ | ------------------- | ----------------------------------------------------------------------------------------------------------------- |
| 1     | `061-altitude-jump-ground-detection` | 058 (done)          | Activates altitude gravity + landing in `BeatEmUpMovementModel`                                                   |
| 2     | `062-depth-aware-collision`          | 061 (logical prior) | Depth-lane gate in `HasCollision`; can start after 061                                                            |
| 3     | `063-shadow-component`               | 061                 | Shadow requires altitude to be live; independent of 062                                                           |
| 4     | `064-beatemup-footprint-rect`        | —                   | JSON `footprint_rect` schema + beat-em-up `Footprint()` accessor. Complements 062 (depth-lane gate); independent. |
| 5     | `065-beatemup-jump-skill`            | 061 (done)          | `BeatEmUpJumpSkill` kit skill: altitude-axis jump with coyote time, buffering, jump-cut.                          |

---

## Dependencies

```
058 (done)
  └─ 061-altitude-jump-ground-detection (done)
       ├─ 062-depth-aware-collision (done)
       ├─ 063-shadow-component (done)
       └─ 065-beatemup-jump-skill (done)

064-beatemup-footprint-rect (independent; feeds collision shapes used by 062)
```

---

## Sequencing Notes

- 061 is the prerequisite for 062, 063, and 065; they can proceed in parallel once 061 is done.
- 062 and 063 have no dependency on each other.
- 064 has no hard dependency; can ship in parallel with 061–063. Pairs naturally with 062 (062 = depth-lane gate mechanism; 064 = footprint shape fed into collision checks).
- 065 has no dependency on 062–064; can ship as soon as 061 is done.

---

## Parking Lot

Future stories not yet written. One line each.

- `060-html-report-generator` — in backlog.

---

## Recently Completed

- `059-thin-game-phase-scenes` — consolidate scene `OnStart` logic into kit layer (pure refactor).
- `058-wire-beatemup-movement` — wire `EightDirectionalMovementSkill` + `BeatEmUpMovementModel` into `BeatEmUpCharacter`; decouple engine from `*PlatformMovementModel`.
- `057-beatemup-movement-model` — `BeatEmUpMovementModel` with 8-way floor movement, friction, speed cap; altitude-silent.
- `056-eight-dir-movement-skill` — `EightDirectionalMovementSkill` kit skill.
- `055-kit-genre-phase-scenes` — genre routing, beat-em-up scene shell, altitude draw sort. Closed by 059.
- `054-` and earlier — see `.agents/work/done/`.
