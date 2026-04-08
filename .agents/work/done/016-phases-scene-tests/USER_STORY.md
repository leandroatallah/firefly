# US-016 — PhasesScene Test Coverage

**Branch:** `016-phases-scene-tests`
**Bounded Context:** Game Logic

## Story

As a developer, I want `internal/game/scenes/phases` to have test coverage for its core lifecycle and goal tracking, so that regressions in level management are caught automatically.

## Context

`internal/game/scenes/phases` is at **0.0% coverage**. It is the foundation for all playable levels — it manages actor spawning, bullet lifecycle, goal tracking, scene transitions, and the freeze-frame effect. It has no test files at all.

## Acceptance Criteria

- **AC1:** `PhasesScene.OnStart()` is tested: verifies actors and items are initialized without panicking.
- **AC2:** Goal completion logic is tested: all goals met triggers scene transition.
- **AC3:** Goal partial completion is tested: not all goals met does not trigger transition.
- **AC4:** `PhasesScene.Update()` is tested with a mock `BodiesSpace` and mock `SceneManager`.
- **AC5:** Freeze-frame integration is tested: `FreezeController` pauses actor updates for the configured frame count.
- **AC6:** Bullet spawn and despawn lifecycle is tested (bullet goes out of bounds or hits target → removed).
- **AC7:** Coverage for `internal/game/scenes/phases` reaches **≥ 40%**.
- **AC8:** All tests use mock implementations from `internal/engine/mocks/` or package-local mocks — no real Ebitengine window or GPU calls.
