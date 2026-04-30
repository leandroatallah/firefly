# User Story 044 — Debug Channels

## Story

**As a** developer instrumenting engine and game logic during development,
**I want** a named-channel debug package in `internal/engine/debug/` that is off by default and enabled per-channel via a JSON config file,
**so that** I can observe runtime state (Actor position, current State, physics contacts, etc.) without introducing 60fps log spam or any overhead in production builds.

## Background

Debugging engine behaviour currently requires temporary `fmt.Println` calls scattered across packages. There is no structured, toggle-able mechanism for observing engine or game state at runtime. A lightweight debug package with named channels and change-only logging gives developers a permanent, low-noise observability layer that is silently disabled when the config file is absent.

## Acceptance Criteria

### AC-1: All Logging Is a No-Op When Uninitialized or Config Is Missing

**Given** `debug.Init(path)` has not been called, or the file at `path` does not exist,
**When** any call to `debug.Log(...)` or `debug.Watch(...)` is made,
**Then** no output is produced and no error is returned — the call is a silent no-op.

### AC-2: Channels Are Enabled via JSON Config

**Given** a file exists at `assets/data/debug.json` with the content:
```json
{ "player_state": true, "physics": false }
```
**And** `debug.Init("assets/data/debug.json")` is called at startup,
**When** `debug.Log("player_state", ...)` is called,
**Then** output is produced to stdout.
**When** `debug.Log("physics", ...)` is called,
**Then** no output is produced (channel is disabled).

### AC-3: `Log` Prints Every Call When Channel Is Enabled

**Given** the channel `"player_state"` is enabled,
**When** `debug.Log("player_state", "pos=%v", pos)` is called on three consecutive frames,
**Then** three lines of output are produced — one per call.

### AC-4: `Watch` Logs Only on Value Change

**Given** the channel `"player_state"` is enabled,
**When** `debug.Watch("player_state", "state", "idle")` is called ten consecutive frames,
**Then** exactly one line of output is produced (on the first call).
**When** `debug.Watch("player_state", "state", "walk")` is called on the eleventh frame,
**Then** a second line of output is produced (value changed from `"idle"` to `"walk"`).

### AC-5: `Watch` Keys Are Scoped Per Channel

**Given** channels `"player_state"` and `"enemy_state"` are both enabled,
**When** `debug.Watch("player_state", "state", "idle")` and `debug.Watch("enemy_state", "state", "idle")` are each called once,
**Then** two lines of output are produced — the same key/value pair is tracked independently per channel.

### AC-6: Engine Wires `debug.Init` at Startup

**Given** `internal/engine/app/engine.go` `NewGame()` is the engine entry point,
**When** `NewGame()` is called,
**Then** `debug.Init("assets/data/debug.json")` is called before any Scene or Actor is constructed — a missing file produces no error and no log output.

### AC-7: No Overhead on Disabled Channels

**Given** a channel is disabled (either absent from the config or set to `false`),
**When** `debug.Log(...)` or `debug.Watch(...)` is called for that channel,
**Then** no string formatting or map lookup beyond the channel-enabled check occurs — the hot path is a single boolean guard.

## API Contract

```go
package debug

// Init loads the JSON config at path. Missing file silently disables all channels.
func Init(path string)

// Log writes a formatted message to stdout on every call when the channel is enabled.
func Log(channel, format string, args ...any)

// Watch writes a formatted message to stdout only when value differs from the
// previous call with the same channel+key pair.
func Watch(channel, key string, value any)
```

## Out of Scope

- No call sites added to `internal/game/` or any Actor/State in this story.
- No structured logging (JSON output, log levels, external sinks).
- No build tags or compile-time removal of debug calls.
- No CLI flags — the JSON config file is the only activation mechanism.

## DDD Glossary (terms used in this story)

| Term | Meaning |
|---|---|
| Actor | An entity with a state machine (player, enemy) |
| State | A named node in an Actor's state machine |
| Body | A physics body with position and velocity |
| Phase | A playable game level Scene |
| Scene | A self-contained game state |
| Contract | A Go interface in `internal/engine/contracts/` |
