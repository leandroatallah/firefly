# SPEC — 017: Sequences Command Test Coverage

**Branch:** `017-sequences-command-tests`
**Bounded Context:** Sequences (`internal/engine/sequences/`)

---

## Goal

Raise `internal/engine/sequences` coverage from **61.1% → ≥ 80%** by adding tests for the zero/near-zero coverage paths identified in the story. No production code changes are required.

---

## Current Coverage Gaps (confirmed via `go tool cover`)

| File | Function | Coverage |
|---|---|---|
| `commands.go` | `ParallelCommand.Init/Update`, `WaitCommand.Init/Update` | 0% |
| `commands_actor.go` | `FollowActorCommand.Update`, `FollowPlayerCommand.Update`, `StopFollowingCommand.Update`, `RemoveActorCommand.Update`, `SetSpeedCommand.Update` | 0% |
| `commands_actor.go` | `resolveActorTargets` (`@query:` branch) | ~35% |
| `commands_camera.go` | `CameraZoomCommand.Update` (all phases), `CameraMoveCommand.Update` (duration path), `CameraResetCommand.Update` (duration path), `CameraShakeCommand.Init/Update` | 0–14% |
| `commands_music.go` | `PlayMusicCommand.Init`, `PauseAllMusicCommand.Init`, `FadeOutAllMusicCommand.Init` | 0% |
| `commands_vfx.go` | `SpawnTextCommand.Init` (screen + overhead paths), `QuakeCommand.Init/Update` | 0–17% |
| `player.go` | `PlaySequence`, `endBlockingPhase` (player-unblock branch) | 0–60% |

---

## Pre-conditions

- `mocks.MockAudioManager` already exists in `internal/engine/mocks/audio_speech_vfx.go` — use it for music command tests.
- `mocks.MockActor`, `mocks.MockVFXManager` already exist — use them for actor/vfx tests.
- `mockSceneWithCamera` is already defined in `commands_camera_test.go` (package-local) — reuse pattern for camera tests.
- No new shared mocks are needed; all new mocks go in `mocks_test.go` (package-local) if required.
- No production code changes.

## Post-conditions

- All new tests pass with `go test ./internal/engine/sequences/...`.
- Coverage reaches ≥ 80%.
- No `ebiten.RunGame` or GPU calls in any test.

---

## Integration Points

| Contract / Type | Used by |
|---|---|
| `contracts/sequences.Command` | All command tests |
| `mocks.MockAudioManager` | `commands_music_test.go` |
| `mocks.MockActor` | `commands_actor_test.go`, `commands_vfx_test.go` |
| `mocks.MockVFXManager` | `commands_vfx_test.go` |
| `camera.Controller` (real, headless) | `commands_camera_test.go` |
| `app.AppContext` | All tests via `setupTestAppContext()` |

---

## Technical Requirements

### 1. `commands.go` — `ParallelCommand` and `WaitCommand`

The coverage report shows `commands.go:44` (`Init`) and `commands.go:75` (`Update`) at 0%. These correspond to `ParallelCommand` and `WaitCommand` (base parallel/wait commands). Tests must cover:
- `ParallelCommand.Init` stores sub-commands; `Update` returns true only when all sub-commands are done.
- `WaitCommand.Init` resets state; `Update` returns true after the configured frame count.

Add to `commands_test.go`.

### 2. `commands_actor.go` — trivial `Update()` methods + `@query:` branch

The five `Update() bool { return true }` stubs on `FollowActorCommand`, `FollowPlayerCommand`, `StopFollowingCommand`, `RemoveActorCommand`, `SetSpeedCommand` are uncovered. One call each suffices.

The `@query:` branch in `resolveActorTargets` (invalid regex path + matching path) needs coverage.

Add to `commands_actor_test.go`.

### 3. `commands_music.go` — all three `Init` paths

Use `mocks.MockAudioManager` injected into `app.AppContext.AudioManager`:
- `PlayMusicCommand.Init`: (a) already playing + `Rewind=false` → no-op; (b) not playing → calls `PlayMusic`; (c) `Volume > 0` → calls `SetVolume`.
- `PauseAllMusicCommand.Init`: calls `PauseAll`.
- `FadeOutAllMusicCommand.Init`: calls `FadeOutAll`.

Add to `commands_music_test.go`.

### 4. `commands_camera.go` — `CameraZoomCommand.Update` phases + `CameraShakeCommand`

`CameraZoomCommand.Update` has three phases (zoom-in, wait, zoom-out). Tests must drive through all three using a real `camera.Controller` (headless):
- Instant zoom (Duration=0): phases 0→1→2 in one Update.
- Timed zoom: simulate frames through all three phases.
- `CameraShakeCommand.Init/Update`: verify `AddTrauma` is called and `Update` returns true.

`CameraMoveCommand.Update` duration path and `CameraResetCommand.Update` duration path also need frame-driven tests.

Add to `commands_camera_test.go`.

### 5. `commands_vfx.go` — `SpawnTextCommand.Init` + `QuakeCommand`

- `SpawnTextCommand.Init` with `Type="screen"`: calls `VFX.SpawnFloatingText`.
- `SpawnTextCommand.Init` with `Type="overhead"` and valid actor: calls `VFX.SpawnFloatingTextAbove`.
- `QuakeCommand.Init` + `Update`: verify trauma is added on frame multiples of 10; returns true after `Duration` frames.

`MockVFXManager` needs `SpawnFloatingText` and `SpawnFloatingTextAbove` — already implemented.

Add to `commands_vfx_test.go`.

### 6. `player.go` — `PlaySequence` + `endBlockingPhase` unblock branch

- `PlaySequence`: test with a valid temp JSON file (use `t.TempDir()`) — verifies the player starts playing.
- `PlaySequence` with invalid path: verifies player stays idle (no panic).
- `endBlockingPhase` with `BlockPlayerMovement=true` and a registered player: verifies `player.UnblockMovement()` is called. Use `mocks.MockActor` as the player.

Add to `sequences_test.go`.

---

## Red Phase (Failing Test Scenario)

The first test to write (drives the TDD cycle) is:

**File:** `commands_music_test.go`
**Test:** `TestPlayMusicCommand_Init_WithMockAudioManager`

```
Given an AppContext with a MockAudioManager
And a PlayMusicCommand{Path: "bgm.ogg", Loop: true}
When Init is called
Then MockAudioManager.PlayedPaths contains "bgm.ogg"
And MockAudioManager.LoopSettings["bgm.ogg"] == true
```

This test currently fails because `commands_music_test.go` only checks struct fields without a real `AppContext.AudioManager`, so `Init` is never actually called with a mock. Adding the mock-wired test will immediately exercise the 0% `Init` path.

---

## Test File Mapping

| Test file | New tests target |
|---|---|
| `commands_test.go` | `ParallelCommand`, `WaitCommand` |
| `commands_actor_test.go` | trivial `Update()` stubs, `@query:` branch |
| `commands_music_test.go` | all three `Init` paths via `MockAudioManager` |
| `commands_camera_test.go` | `CameraZoomCommand` all phases, `CameraShakeCommand`, `CameraMoveCommand`/`CameraResetCommand` duration paths |
| `commands_vfx_test.go` | `SpawnTextCommand.Init` both paths, `QuakeCommand` |
| `sequences_test.go` | `PlaySequence`, `endBlockingPhase` unblock branch |
