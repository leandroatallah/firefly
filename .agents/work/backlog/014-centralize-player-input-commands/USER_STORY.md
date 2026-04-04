# 014 — Centralize Player Input Commands

## Problem

Key bindings are scattered across multiple files. `skill_shooting.go` hardcodes `KeyX`, `KeyUp/W`, `KeyDown/S`, `KeyLeft/A`, `KeyRight/D` directly inside `HandleInput` and `Update`. `climber.go` duplicates the same directional keys independently. Any future key remapping requires hunting down multiple files.

## Solution

Introduce a `PlayerCommands` struct in `internal/engine/input` that represents the logical input state (up, down, left, right, shoot) — engine-level, key-agnostic from the consumer's perspective. It ships with a default `ReadPlayerCommands()` function that maps the current hardcoded keys, so nothing breaks out of the box.

The engine defines the *shape* of commands; the game overrides *how they're read*.

## How the game overrides it

In `internal/game/`, a game-specific input provider replaces `ReadPlayerCommands` — either by registering a custom reader function, or by the game wiring its own `PlayerCommands` builder and injecting it into the player/skill constructors. This keeps game-specific bindings (e.g. gamepad, alternative keyboard layout) entirely out of the engine.

## Affected files

- `internal/engine/input/` — add `commands.go` with `PlayerCommands` struct + default `ReadPlayerCommands()`
- `internal/engine/physics/skill/skill_shooting.go` — `HandleInput` and `Update` consume `PlayerCommands` instead of calling `ebiten.IsKeyPressed` directly
- `internal/game/entity/actors/player/climber.go` — `Update` reads `PlayerCommands` for duck and facing direction
- `internal/game/` — optionally registers a custom commands reader at setup time if bindings differ from defaults

## Acceptance Criteria

- All hardcoded `ebiten.IsKeyPressed` calls for player actions are removed from `skill_shooting.go` and `climber.go`
- A single `ReadPlayerCommands()` in `internal/engine/input` serves as the default binding source
- The game layer can override bindings without modifying engine code
