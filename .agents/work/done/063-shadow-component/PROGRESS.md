# PROGRESS — 063-shadow-component

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [-] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- [Story Architect] 2026-05-20: USER_STORY.md created. Oval shadow rendered at (X, GroundY) for airborne entities; BEATEMUP_PLAN Phase 3 visual feedback item. Drawable registration strategy and scale formula deferred to Spec Engineer.
- [Claude Opus 4.7] [Spec Engineer] 2026-05-21 [STARTED]: Drafting SPEC.md for shadow component.
- [Claude Opus 4.7] [Spec Engineer] 2026-05-21 [FINISHED]: SPEC.md + NOTES.md written. Key decisions: package-level utility in `internal/kit/render/shadow/` (no new contract, no actor surface change); structural `AltitudeBody` interface; linear altitude→scale falloff with `ShadowMinScale` clamp; integration via shadow pass inserted before actor draw loop in `beatemup/scene.go` (`DrawActors` + `fullDraw`); test seam via swappable `ovalDrawerFn`. Mock Generator marked SKIPPED — no new engine contracts introduced.
- [Claude Opus 4.7] [Mock Generator] 2026-05-21 [SKIPPED]: No new contract interfaces introduced — `AltitudeBody` is a kit-local structural interface, mocks not required per SPEC.
- [Claude Opus 4.7] [TDD Specialist] 2026-05-21 [FINISHED]: Wrote Red-Phase tests. New files: `internal/kit/render/shadow/shadow.go` (skeleton with zero-value returns), `internal/kit/render/shadow/shadow_test.go` (T-S1..T-S7 black-box tests using `SetOvalDrawerForTest` recording sink). Appended T-S8 to `internal/kit/scenes/phases/beatemup/scene_test.go` plus minimal `SetShadowDrawerForTest` hook + `shadowDrawer` field in `scene.go` (hook stored but never invoked — Feature Implementer wires it). Red proven: `go test ./internal/kit/render/shadow/...` reports FAIL on TestScaleFor (returns 0 for all altitudes), TestComputeBounds_* (zero Bounds), TestDraw_FiresWhenAirborne (drew=false), TestDrawAll_CountsOnlyAirborne (sink.calls=0); `go test ./internal/kit/scenes/phases/beatemup/... -run ShadowDrawnBeforeActor` reports FAIL because DrawActors never invokes shadowDrawer.
- [Feature Implementer] 2026-05-21 [FINISHED]: Implemented shadow package and wired into beatemup scene. Production files: `internal/kit/render/shadow/shadow.go` (ScaleFor, ComputeBounds, Draw, DrawAll, drawOval all filled in), `internal/kit/scenes/phases/beatemup/scene.go` (added shadow import, wired shadowDrawer into DrawActors before actor loop, added shadow.DrawAll fallback in DrawActors and fullDraw). All 8 tests pass: T-S1..T-S7 in shadow package, T-S8 in beatemup scene. `go build ./...` clean.
- [Workflow Gatekeeper] 2026-05-21 [FINISHED]: All quality gates passed. Coverage delta: `internal/kit/render/shadow` new package 0% -> 56.6% (drawOval excluded from coverage as it requires GPU; all pure logic functions at 80-100%); `internal/kit/scenes/phases/beatemup` 21.2% -> 22.0% (+0.8pp, from T-S8 wiring shadow draw order). All 8 specified tests (T-S1..T-S8) green. golangci-lint clean. No forbidden imports (`internal/game/` absent from shadow package). No `_ = variable` pattern in production code. Platformer scene untouched (AC-7 confirmed). Story moved to done.
