# SPEC-014 — Centralize Player Input Commands

## Branch

`014-centralize-player-input-commands`

## Context

`ebiten.IsKeyPressed` and `inpututil.IsKeyJustPressed` are called directly in multiple engine/game/ui files:

**Skills & Movement:**
- `internal/engine/physics/skill/skill_shooting.go` — `HandleInput`, `Update` (shoot, directional)
- `internal/engine/physics/skill/skill_platform_jump.go` — `HandleInput` (jump activation, jump cut release)
- `internal/engine/physics/skill/skill_dash.go` — `HandleInput` (dash activation)
- `internal/engine/physics/skill/skill_platform_move.go` — `HandleInput` (horizontal movement)
- `internal/game/entity/actors/player/climber.go` — `Update` (duck, facing direction)

**UI:**
- `internal/engine/ui/menu/menu.go` — `Update` (navigate up/down, select, cancel)
- `internal/engine/ui/speech/dialogue.go` — `Update` (advance dialogue, skip typing)

This story introduces a `PlayerCommands` struct in `internal/engine/input` as the single source of logical input state, and threads it through all skill, movement, and UI consumers.

## Technical Requirements

### 1. `internal/engine/input/commands.go` (new file)

```go
type PlayerCommands struct {
    Up, Down, Left, Right bool  // directional
    Shoot                 bool  // X key
    Jump                  bool  // Space key
    Dash                  bool  // Shift key
    Confirm               bool  // Enter key (menu/dialogue)
    Cancel                bool  // Escape key (menu)
}

// ReadPlayerCommands returns the default keyboard mapping.
// Swappable via CommandsReader for game-layer overrides.
func ReadPlayerCommands() PlayerCommands
```

Expose a swappable function var (same pattern as `isKeyPressed`):

```go
//nolint:gochecknoglobals
var CommandsReader func() PlayerCommands = ReadPlayerCommands
```

`ReadPlayerCommands` maps:
- `Up` → `KeyUp || KeyW`
- `Down` → `KeyDown || KeyS`
- `Left` → `KeyLeft || KeyA`
- `Right` → `KeyRight || KeyD`
- `Shoot` → `KeyX`
- `Jump` → `KeySpace`
- `Dash` → `KeyShift`
- `Confirm` → `KeyEnter`
- `Cancel` → `KeyEscape`

### 2. `skill_shooting.go` changes

- `HandleInput` replaces its five `ebiten.IsKeyPressed` calls with `input.CommandsReader()`.
- `Update` replaces `ebiten.IsKeyPressed(ebiten.KeyX)` with `input.CommandsReader().Shoot`.
- `ActivationKey()` is **kept** (it is part of the `SkillBase` interface and unrelated to this story).
- No new fields on `ShootingSkill`.

### 3. `skill_platform_jump.go` changes

- `HandleInput` replaces `inpututil.IsKeyJustPressed(s.activationKey)` with `input.CommandsReader().Jump` (using `inpututil` for edge detection is replaced by the caller's responsibility to detect state changes).
- `HandleInput` replaces `inpututil.IsKeyJustReleased(s.activationKey)` with tracking `Jump` state changes in `Update`.
- `ActivationKey()` is **kept**.

### 4. `skill_dash.go` changes

- `HandleInput` replaces `inpututil.IsKeyJustPressed(d.activationKey)` with `input.CommandsReader().Dash`.
- `ActivationKey()` is **kept**.

### 5. `skill_platform_move.go` changes

- `HandleInput` replaces `input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft)` and `input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight)` with `input.CommandsReader().Left` and `input.CommandsReader().Right`.

### 6. `climber.go` changes

- `Update` replaces all `ebiten.IsKeyPressed` calls with a single `cmds := input.CommandsReader()` at the top of the method, then uses `cmds.Down`, `cmds.Left`, `cmds.Right`.

### 7. `internal/engine/ui/menu/menu.go` changes

- `Update` replaces `inpututil.IsKeyJustPressed` calls with state tracking via `input.CommandsReader()`:
  - `KeyW || KeyUp` → `cmds.Up`
  - `KeyS || KeyDown` → `cmds.Down`
  - `KeyEnter` → `cmds.Confirm`
  - `KeyEscape` → `cmds.Cancel`
- Caller becomes responsible for detecting state changes (edge detection).

### 8. `internal/engine/ui/speech/dialogue.go` changes

- `Update` replaces `inpututil.IsKeyJustPressed(ebiten.KeyEnter)` with state tracking via `input.CommandsReader().Confirm`.
- `shouldSkipTyping()` replaces `inpututil.IsKeyJustPressed(ebiten.KeyEnter)` with `input.CommandsReader().Confirm`.

## Pre-conditions

- `internal/engine/input` package already has the `isKeyPressed` swappable var pattern.
- `ShootingSkill.HandleInputWithDirection` already accepts `up, down, left, right bool` — no signature change needed there.

## Post-conditions

- Zero `ebiten.IsKeyPressed` calls remain in `skill_shooting.go`, `skill_platform_jump.go`, `skill_dash.go`, `skill_platform_move.go`, `climber.go`, `menu.go`, and `dialogue.go`.
- Zero `inpututil.IsKeyJustPressed` / `inpututil.IsKeyJustReleased` calls remain in `skill_platform_jump.go`, `skill_dash.go`, `menu.go`, and `dialogue.go` (replaced by state tracking via `CommandsReader`).
- Zero `input.IsSomeKeyPressed` calls remain in `skill_platform_move.go` (replaced by direct `CommandsReader` field access).
- `input.CommandsReader` can be replaced in tests and game setup without touching engine or UI code.
- All existing tests continue to pass.

## Integration Points (Bounded Context: Input)

| File | Change |
|---|---|
| `internal/engine/input/commands.go` | New: `PlayerCommands`, `ReadPlayerCommands`, `CommandsReader` |
| `internal/engine/physics/skill/skill_shooting.go` | Consume `input.CommandsReader()` for shoot + directional |
| `internal/engine/physics/skill/skill_platform_jump.go` | Consume `input.CommandsReader().Jump` |
| `internal/engine/physics/skill/skill_dash.go` | Consume `input.CommandsReader().Dash` |
| `internal/engine/physics/skill/skill_platform_move.go` | Consume `input.CommandsReader().Left/Right` |
| `internal/game/entity/actors/player/climber.go` | Consume `input.CommandsReader()` for duck + facing |
| `internal/engine/ui/menu/menu.go` | Consume `input.CommandsReader()` for navigation + confirm/cancel |
| `internal/engine/ui/speech/dialogue.go` | Consume `input.CommandsReader().Confirm` |

No new contract interface is required — `CommandsReader` is a plain function var, consistent with the existing `isKeyPressed` pattern in the same package.

## Red Phase

**File:** `internal/engine/input/commands_test.go`

**Scenario:** `ReadPlayerCommands` returns the correct struct when the injected key reader reports specific keys pressed.

```
Given isKeyPressed is stubbed to return true only for KeyX and KeyUp
When  ReadPlayerCommands() is called
Then  PlayerCommands{Up: true, Shoot: true, Down: false, Left: false, Right: false, Jump: false, Dash: false, Confirm: false, Cancel: false}
```

The test must fail before `commands.go` exists (compilation error / missing symbol), and pass once the implementation is in place.

Table cases to cover:
- No keys → all false
- KeyX only → Shoot true, rest false
- KeyUp + KeyW both → Up true (either key triggers it)
- KeyDown + KeyS → Down true
- KeyLeft + KeyA → Left true
- KeyRight + KeyD → Right true
- KeySpace → Jump true
- KeyShift → Dash true
- KeyEnter → Confirm true
- KeyEscape → Cancel true
- All keys → all true
