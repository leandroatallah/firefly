# 017 — Sequences Command Test Coverage

**Branch:** `017-sequences-command-tests`
**Bounded Context:** Sequences

## Story

As a developer, I want the untested sequence commands in `internal/engine/sequences` to have test coverage, so that scripted event regressions are caught automatically.

## Context

`internal/engine/sequences` is at **61.1% coverage**. The following functions are at 0% or near-0%:

| Function | Coverage |
|---|---|
| `commands.go` — `Init`, `Update` (base parallel/wait commands) | 0.0% |
| `commands_actor.go` — `SetStateCommand.Update`, `FaceDirectionCommand.Update`, `WaitUntilGroundedCommand.Update`, `WaitUntilStateCommand.Update`, `WaitUntilAnimationDoneCommand.Update` | 0.0% |
| `commands_camera.go` — `CameraShakeCommand.Init/Update`, `CameraZoomCommand.Update` (partial) | 0–14% |
| `commands_vfx.go` — `SpawnParticlesCommand.Init`, `SpawnFloatingTextCommand.Init/Update` | 0–17% |
| `commands_music.go` — all three `Init` functions | 0.0% |
| `player.go` — `PlaySequence`, `endBlockingPhase` | 0–60% |

## Acceptance Criteria

- **AC1:** `commands.go` base commands (`ParallelCommand`, `WaitCommand`) are tested for `Init` and `Update` completion conditions.
- **AC2:** All zero-coverage actor commands (`SetState`, `FaceDirection`, `WaitUntilGrounded`, `WaitUntilState`, `WaitUntilAnimationDone`) have at least one passing test each.
- **AC3:** `commands_music.go` — all three `Init` paths are tested using a mock audio player.
- **AC4:** `player.go` — `PlaySequence` and `endBlockingPhase` are tested.
- **AC5:** `commands_camera.go` — `CameraShakeCommand` and `CameraZoomCommand` `Update` branches are covered.
- **AC6:** Coverage for `internal/engine/sequences` reaches **≥ 80%** (up from 61.1%).
- **AC7:** All tests use mock implementations at system boundaries (camera, audio, actor, vfx) — no real Ebitengine window.
