# 014 — Centralize Player Input Commands

## Problem

Key bindings are scattered across multiple skill, actor, and UI files:
- `skill_shooting.go` hardcodes `KeyX`, `KeyUp/W`, `KeyDown/S`, `KeyLeft/A`, `KeyRight/D` in `HandleInput` and `Update`
- `skill_platform_jump.go` hardcodes `KeySpace` in `HandleInput`
- `skill_dash.go` hardcodes `KeyShift` in `HandleInput`
- `skill_platform_move.go` hardcodes `KeyLeft/A`, `KeyRight/D` in `HandleInput`
- `climber.go` duplicates directional keys in `Update`
- `menu.go` hardcodes `KeyW/Up`, `KeyS/Down`, `KeyEnter`, `KeyEscape` in `Update`
- `dialogue.go` hardcodes `KeyEnter` in `Update` and `shouldSkipTyping()`

Any future key remapping requires hunting down multiple files and modifying engine/UI code.

## Solution

Introduce a `PlayerCommands` struct in `internal/engine/input` that represents the logical input state (up, down, left, right, shoot, jump, dash) — engine-level, key-agnostic from the consumer's perspective. It ships with a default `ReadPlayerCommands()` function that maps the current hardcoded keys, so nothing breaks out of the box.

The engine defines the *shape* of commands; the game overrides *how they're read*.

## How the game overrides it

In `internal/game/`, a game-specific input provider replaces `ReadPlayerCommands` — either by registering a custom reader function, or by the game wiring its own `PlayerCommands` builder and injecting it into the player/skill constructors. This keeps game-specific bindings (e.g. gamepad, alternative keyboard layout) entirely out of the engine.

## Affected files

- `internal/engine/input/` — add `commands.go` with `PlayerCommands` struct + default `ReadPlayerCommands()`
- `internal/engine/physics/skill/skill_shooting.go` — `HandleInput` and `Update` consume `PlayerCommands` instead of calling `ebiten.IsKeyPressed` directly
- `internal/engine/physics/skill/skill_platform_jump.go` — `HandleInput` consumes `PlayerCommands.Jump` instead of `inpututil.IsKeyJustPressed`
- `internal/engine/physics/skill/skill_dash.go` — `HandleInput` consumes `PlayerCommands.Dash` instead of `inpututil.IsKeyJustPressed`
- `internal/engine/physics/skill/skill_platform_move.go` — `HandleInput` consumes `PlayerCommands.Left/Right` instead of `input.IsSomeKeyPressed`
- `internal/game/entity/actors/player/climber.go` — `Update` reads `PlayerCommands` for duck and facing direction
- `internal/engine/ui/menu/menu.go` — `Update` consumes `PlayerCommands` for navigation and confirm/cancel
- `internal/engine/ui/speech/dialogue.go` — `Update` and `shouldSkipTyping()` consume `PlayerCommands.Confirm`
- `internal/game/` — optionally registers a custom commands reader at setup time if bindings differ from defaults

## Acceptance Criteria

- All hardcoded `ebiten.IsKeyPressed` calls for player actions are removed from `skill_shooting.go`, `skill_platform_jump.go`, `skill_dash.go`, `skill_platform_move.go`, `climber.go`, `menu.go`, and `dialogue.go`
- All `inpututil.IsKeyJustPressed` / `inpututil.IsKeyJustReleased` calls in jump, dash, menu, and dialogue are replaced by state tracking via `PlayerCommands`
- All `input.IsSomeKeyPressed` calls in `skill_platform_move.go` are replaced by direct `PlayerCommands` field access
- A single `ReadPlayerCommands()` in `internal/engine/input` serves as the default binding source
- The game layer can override bindings without modifying engine or UI code
