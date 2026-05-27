# PROGRESS — 072-debug-flag-overlay

**Status:** Done

## Pipeline State

- [x] Story Architect
- [x] Spec Engineer
- [ ] Mock Generator
- [x] TDD Specialist
- [x] Feature Implementer
- [x] Workflow Gatekeeper

## Log

- TDD Specialist: wrote failing Red tests.
  - `internal/engine/debug/registry_test.go` — T-R1..T-R6 prove `Register`/`List` behavior and `Reset`/`InitFromReader` integration are missing (undefined `debug.Register`, `debug.List`).
  - `internal/engine/ui/debugoverlay/overlay_test.go` — T-O1..T-O8 prove `DebugOverlay` (with injectable `keyJustPressed`) and its open/close/cursor/toggle semantics do not yet exist (undefined `New`).
  - `internal/engine/app/app_test.go` (T-G1 appended) — proves `Game.DebugOverlay()` accessor and the "skip Scene/Dialogue update when overlay open" guard are missing.
  - Red proof: `go test ./internal/engine/debug/... ./internal/engine/ui/debugoverlay/... ./internal/engine/app/...` fails at build/link with the exact missing-behavior symbols above.

- Feature Implementer: all tests Green.
  - `internal/engine/debug/debug.go` — switched `channels` to `atomic.Pointer[map[string]*bool]`, added `registry map[string]*bool` global, updated `InitFromReader`, `Enabled`, and `Reset`.
  - `internal/engine/debug/registry.go` — new file; `Register`, `List`, `Entry` type.
  - `internal/engine/ui/debugoverlay/overlay.go` — new file; `DebugOverlay` struct with `Open`/`Close`/`IsOpen`/`SetFont`/`Update`/`Draw`.
  - `internal/game/app/config.go` — added `debug.Register` calls for `cam_debug` and `collision_box`.
  - `internal/engine/app/engine.go` — added `debugOverlay` field, `DebugOverlay()` accessor, F1 toggle + overlay suppression logic in `Update`, `debugOverlay.Draw` in `Draw`, removed old `debugVisible` field and `DebugPhysics` conditional.
  - Verification: `go test ./internal/engine/debug/... ./internal/engine/ui/debugoverlay/... ./internal/engine/app/...` — all pass. `go build ./...` — clean.

- Workflow Gatekeeper: accepted after two lint fixes; all quality gates pass.
  - Fixes applied before acceptance:
    - Removed `_ = line` blank assignment from `internal/engine/ui/debugoverlay/overlay.go:103` (prohibited by project standard "No `_ = variable` in production code").
    - Removed explicit `bool` type annotations from `var a bool = false` / `var a bool = true` in `overlay_test.go` T-O5/T-O6 (staticcheck ST1023).
  - Coverage delta (new packages, all positive):
    - `internal/engine/debug`: 93.2% (new functions Register/List at 100%; overall package up from prior baseline)
    - `internal/engine/ui/debugoverlay`: 63.8% (new package; Update at 95.2%, Draw at 33.3% — draw path through font-less path is exercised by T-O8 smoke test)
    - `internal/engine/app`: 57.1% (NewGame, DebugOverlay accessor, Update, Draw, Layout all covered)
  - AC-9 import constraint confirmed: `go list -deps internal/engine/ui/debugoverlay` shows no `internal/game/` or `internal/kit/` entries.
  - `golangci-lint run ./...` — 0 issues after fixes.
  - All 28 tests pass across 3 packages (T-R1..T-R6, T-O1..T-O8, T-G1, existing regression suite).
