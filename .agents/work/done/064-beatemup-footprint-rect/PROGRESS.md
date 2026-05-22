# PROGRESS — 064-beatemup-footprint-rect

**Status:** Done

## Pipeline Status

| Stage | Agent | Status |
|---|---|---|
| Story Architect | story-architect | ✅ |
| Spec Engineer | spec-engineer | ✅ |
| Mock Generator | mock-generator | ⏭️ Skipped |
| TDD Specialist | tdd-specialist | ✅ |
| Feature Implementer | feature-implementer | ✅ |
| Workflow Gatekeeper | workflow-gatekeeper | ✅ |

## Log

- Story Architect 2026-05-20: USER_STORY.md created.
- [Opus 4.7] [Spec Engineer] 2026-05-22 [STARTED]: Reading USER_STORY, constitution, schemas, beatemup character, collision pipeline.
- [Opus 4.7] [Spec Engineer] 2026-05-22 [FINISHED]: SPEC.md + NOTES.md written. Key decisions: (1) `FootprintRect *ShapeRect` added to `AssetData` (optional pointer); (2) `BeatEmUpCharacter` stores a `map[ActorStateEnum]image.Rectangle` populated at construction from `spriteData.Assets` via `stateMap`; (3) override `CollisionPosition()` on `*BeatEmUpCharacter` so engine-level checks (`space.HasCollision`, `ApplyValidPosition` → `ResolveCollisions`) automatically use the footprint for beatemup actors only — keeps engine layer untouched and AC-6 isolated by construction; (4) zero-size footprint treated as absent (fallback per AC-3); (5) no mirroring on FaceLeft (parity with existing `collision_rect`). No new contracts introduced → Mock Generator SKIPPED.
- [Opus 4.7] [Mock Generator] 2026-05-22 [SKIPPED]: No new contract interfaces — `CollisionPosition()` override uses existing `body.Collidable` dispatch.
- [Opus 4.7] [TDD Specialist] 2026-05-22 [FINISHED]: Added failing tests in `internal/engine/data/schemas/json_test.go` (T-S1, T-S2) and `internal/kit/actors/beatemup/beatemup_character_test.go` (T-F1..T-F8, T-I1). Red proof: schema tests fail because `AssetData.FootprintRect` is undefined; beatemup tests fail because `BeatEmUpCharacter.Footprint` is undefined and `CollisionPosition` override is missing. Tests verify observable behavior via public APIs (`Footprint()`, `CollisionPosition()`, `space.HasCollision`) — no internal call counts. Integration test T-I1 wires through the real engine `space.HasCollision` so the override change is end-to-end observable.
- [Sonnet 4.6] [Feature Implementer] 2026-05-22 [FINISHED]: Production code written. All tests pass (Green).
  - Modified: `internal/engine/data/schemas/json.go` — added `FootprintRect *ShapeRect` field to `AssetData`.
  - Modified: `internal/kit/actors/beatemup/beatemup_character.go` — added `footprints` field, `buildFootprints` helper, `Footprint()` method, `CollisionPosition()` override. `buildFootprints` is called in `NewBeatEmUpCharacter` after construction. Both methods use explicit `c.CollidableBody.CollisionPosition()` selector to avoid infinite recursion.
  - Test result: `ok github.com/boilerplate/ebiten-template/internal/engine/data/schemas` and `ok github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup` — all 11 new tests (T-S1, T-S2, T-F1..T-F8, T-I1) pass alongside existing tests.
- [Sonnet 4.6] [Workflow Gatekeeper] 2026-05-22 [FINISHED]: All gates passed.
  - Red-Green-Refactor: TDD specialist confirmed red proof; feature implementer confirmed green. No refactor phase required (code is clean as written).
  - Spec alignment: AC-1 through AC-8 fully implemented and tested. All 11 spec test triples (T-S1, T-S2, T-F1..T-F8, T-I1) present and passing.
  - Coverage delta: schemas 100.0% (no change, already full); beatemup 50.0% → 75.0% (+25 pp). Positive delta confirmed.
  - golangci-lint: 0 issues on modified packages.
  - Layer compliance: `FootprintRect` consumed only inside `internal/kit/actors/beatemup/`; no new imports in schemas; `Footprint()` has no callers outside the package. AC-6 and AC-7 satisfied.
  - Standards: table-driven test used for T-I1 (two sub-cases: body overlaps/feet clear, feet overlap); no `_ = variable` in production code; DDD boundary preserved; no GPU code (headless-compatible).
