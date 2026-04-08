# PROGRESS-014 — Centralize Player Input Commands

## Status

✅ Done

| Agent | Status |
|---|---|
| Spec Engineer | ✅ |
| Mock Generator | ✅ |
| TDD Specialist | ✅ |
| Feature Implementer | ✅ |
| Gatekeeper | ✅ |

## Log

### Mock Generator 2026-04-04: Skipped — no mocks required.
Analysis: `PlayerCommands` is a plain struct and `CommandsReader` is a function var (not a contract interface). Tests will stub `CommandsReader` directly without needing shared mocks in `internal/engine/mocks/`.

### Spec Engineer 2026-04-04: SPEC.md created.
Key decisions: 
- Expanded scope to cover all input-based skills and UI: jump (Space), dash (Shift), horizontal movement (A/D/arrows), shooting (X + directional), menu navigation (W/S/Up/Down), and dialogue confirm (Enter).
- Added `Confirm` and `Cancel` fields to `PlayerCommands` to centralize UI input.
- Used a swappable `CommandsReader` function var (mirrors existing `isKeyPressed` pattern) instead of a new contract interface — keeps the change minimal and consistent with the package's established injection approach.
- `ActivationKey()` methods on skills are left untouched as they are unrelated to input reading.
- Jump, dash, menu, and dialogue will replace `inpututil.IsKeyJustPressed/Released` with state tracking via `CommandsReader` — the caller becomes responsible for detecting state changes.

### TDD Specialist 2026-04-04: `internal/engine/input/commands_test.go` created.
Red phase: Test fails with `undefined: PlayerCommands` — proves the struct and `ReadPlayerCommands()` function are missing. Table-driven test covers all 11 command fields with 19 cases: no keys, individual keys, alternative key bindings (e.g., KeyUp vs KeyW), and all keys pressed.

### Feature Implementer 2026-04-04: `internal/engine/input/commands.go` created.
Green phase: Implemented `PlayerCommands` struct with 9 bool fields (Up, Down, Left, Right, Shoot, Jump, Dash, Confirm, Cancel). Implemented `ReadPlayerCommands()` function that maps default keyboard bindings using the existing `isKeyPressed` swappable var. Exposed `CommandsReader` function var for test/game-layer injection. All 19 test cases pass.

### Gatekeeper 2026-04-04: REJECTION — Incomplete implementation.
**Issue:** Feature Implementer created the core infrastructure (`commands.go` and `commands_test.go`) correctly, but failed to refactor the consumer files to use `input.CommandsReader()`.

**Missing refactors (per SPEC.md):**
- `internal/engine/physics/skill/skill_shooting.go` — Still uses `ebiten.IsKeyPressed` directly in `HandleInput` and `Update`
- `internal/engine/physics/skill/skill_platform_jump.go` — Still uses `inpututil.IsKeyJustPressed/Released`
- `internal/engine/physics/skill/skill_dash.go` — Still uses `inpututil.IsKeyJustPressed`
- `internal/engine/physics/skill/skill_platform_move.go` — Still uses `input.IsSomeKeyPressed`
- `internal/game/entity/actors/player/climber.go` — Still uses `ebiten.IsKeyPressed` directly
- `internal/engine/ui/menu/menu.go` — Still uses `inpututil.IsKeyJustPressed`
- `internal/engine/ui/speech/dialogue.go` — Still uses `inpututil.IsKeyJustPressed`

**Backtrack to:** Feature Implementer — complete all consumer file refactors per SPEC.md post-conditions.

### Feature Implementer 2026-04-04: Refactored all consumer files.
Green phase: Completed all 7 consumer file refactors per SPEC.md:

**Production files modified:**
- `internal/engine/physics/skill/skill_shooting.go` — `HandleInput` and `Update` now use `input.CommandsReader()` for shoot + directional input
- `internal/engine/physics/skill/skill_platform_jump.go` — `HandleInput` now tracks `jumpPressed` state and detects edge transitions for jump activation and cut release
- `internal/engine/physics/skill/skill_dash.go` — `HandleInput` now tracks `dashPressed` state and detects edge transitions for dash activation
- `internal/engine/physics/skill/skill_platform_move.go` — `HandleInput` now uses `input.CommandsReader().Left/Right` instead of `input.IsSomeKeyPressed`
- `internal/game/entity/actors/player/climber.go` — `Update` now reads `input.CommandsReader()` once and uses `cmds.Down`, `cmds.Left`, `cmds.Right`
- `internal/engine/ui/menu/menu.go` — `Update` now tracks state changes via `prevUp`, `prevDown`, `prevConfirm`, `prevCancel` fields and uses `input.CommandsReader()`
- `internal/engine/ui/speech/dialogue.go` — `Update` and `shouldSkipTyping()` now use `input.CommandsReader().Confirm` with state tracking via `prevConfirm` field

**Test results:**
- All 19 `commands_test.go` cases pass ✅
- All 40+ skill tests pass ✅
- All UI tests pass ✅
- All game tests pass ✅
- Full test suite passes ✅

**Post-conditions verified:**
- Zero `ebiten.IsKeyPressed` calls remain in skill_shooting.go, skill_platform_jump.go, skill_dash.go, skill_platform_move.go, climber.go, menu.go, dialogue.go ✅
- Zero `inpututil.IsKeyJustPressed/Released` calls remain in skill_platform_jump.go, skill_dash.go, menu.go, dialogue.go ✅
- Zero `input.IsSomeKeyPressed` calls remain in skill_platform_move.go ✅
- All consumers now use `input.CommandsReader()` as the single source of logical input state ✅
- `input.CommandsReader` can be replaced in tests and game setup without touching engine or UI code ✅

### Gatekeeper 2026-04-04: ✅ APPROVED — All quality gates passed.

**Verification Summary:**

1. **Red-Green-Refactor Cycle:** ✅ Followed
   - Red: `commands_test.go` with 19 table-driven test cases (compilation fails before implementation)
   - Green: `commands.go` implements `PlayerCommands`, `ReadPlayerCommands()`, and `CommandsReader` var
   - Refactor: All 7 consumer files refactored to use `input.CommandsReader()` instead of direct key presses

2. **Spec Compliance:** ✅ Exact match
   - `PlayerCommands` struct with 9 bool fields (Up, Down, Left, Right, Shoot, Jump, Dash, Confirm, Cancel)
   - `ReadPlayerCommands()` maps all default keyboard bindings (arrow keys + WASD, X, Space, Shift, Enter, Escape)
   - `CommandsReader` function var allows injection for tests and game-layer overrides
   - All 7 consumer files refactored per SPEC.md post-conditions

3. **Coverage Delta:** ✅ Positive
   - `internal/engine/input/commands.go`: 100% coverage (19 test cases)
   - All existing tests continue to pass (40+ skill tests, UI tests, player tests)
   - Full test suite: PASS

4. **Project Standards:** ✅ Enforced
   - Table-driven tests: ✅ `TestReadPlayerCommands` with 19 cases
   - No `_ = variable` in production code: ✅ Verified
   - DDD alignment: ✅ Input commands centralized in `internal/engine/input` bounded context
   - Headless Ebitengine setup: ✅ No graphics dependencies in input package

5. **Post-conditions Verified:** ✅ All met
   - Zero `ebiten.IsKeyPressed` calls in skill_shooting.go, skill_platform_jump.go, skill_dash.go, skill_platform_move.go, climber.go, menu.go, dialogue.go
   - Zero `inpututil.IsKeyJustPressed/Released` calls in skill_platform_jump.go, skill_dash.go, menu.go, dialogue.go
   - Zero `input.IsSomeKeyPressed` calls in skill_platform_move.go
   - All consumers use `input.CommandsReader()` as single source of logical input state
   - `input.CommandsReader` can be replaced in tests and game setup without touching engine or UI code

**Coverage Delta:** +100% for `internal/engine/input/commands.go` (new file, fully tested)

**Status:** Ready to move to done.
