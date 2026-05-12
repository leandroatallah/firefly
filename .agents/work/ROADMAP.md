# Roadmap

Cross-story sequencing. **Story IDs do not imply execution order — read this file for current priority.**

Update this file when a story moves to `active/` or `done/`, or when sequencing changes.

---

## Active Sequence

| Order | Story | Depends on | Notes |
|---|---|---|---|
| 1 | `056-eight-dir-movement-skill` | — | Kit skill addition. Independent. |
| 2 | `057-beatemup-movement-model` | — | Engine physics addition. Can run in parallel with 056. |
| 3 | `058-wire-beatemup-movement` | 056, 057 | Wires skill + model. Unlocks playable beat-em-up. |
| 4 | `059-thin-game-phase-scenes` | none functional | Large refactor. Land last to amortize review cost. |

---

## Dependencies

```
056 ──┐
      ├──► 058 ──► (beat-em-up playable)
057 ──┘

059 ── independent ── (architectural cleanup; lands after 058 for cleaner diff)
```

---

## Sequencing Notes

- **058 before 059**: 058 adds one `Camera().SetBounds(tilemapRect)` call to the beat-em-up scene's `OnStart`. 059 moves `OnStart` into kit. Landing 058 first lets that line move with the rest of `OnStart` in 059 — no rework. Reverse order forces 058 to wire through 059's new `Options` plumbing.
- **059 alone has no user-facing change** — it is pure refactor. Sequence so that user-facing stories (056 → 057 → 058) ship first.
- **056 and 057 in parallel** is safe — different packages, no shared files.

---

## Parking Lot

Future stories not yet written. One line each.

- (none)

---

## Recently Completed

- `055-kit-genre-phase-scenes` — genre routing, beat-em-up scene shell, altitude draw sort. Deferred items recorded in its `PROGRESS.md`; closed by 059.
- `054-` and earlier — see `.agents/work/done/`.
