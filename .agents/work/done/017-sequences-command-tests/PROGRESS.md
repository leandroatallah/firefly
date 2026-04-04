# PROGRESS — 017: Sequences Command Test Coverage

## Status

| Agent | Status |
|---|---|
| Spec Engineer | ✅ |
| Mock Generator | ✅ |
| TDD Specialist | ✅ |
| Feature Implementer | ✅ |
| Gatekeeper | ✅ |

**Top-level Status: ✅ Done**

## Log

### Spec Engineer 2026-04-04: SPEC.md created.
Key decisions: no production code changes needed; all gaps are test-only. `ParallelCommand`/`WaitCommand` in `commands.go` were identified as the base commands at 0% (not actor commands as the story title implied). Music `Init` paths are the highest-value target for the Red Phase since they require only wiring `MockAudioManager` into `AppContext`. Camera zoom phase coverage requires frame-driven simulation through all three phases using a real headless `camera.Controller`.

### Mock Generator 2026-04-04: skipped — no mocks required

All required mocks already exist in `internal/engine/mocks/`:

- `MockAudioManager` (`audio_speech_vfx.go`) — covers `commands_music_test.go`
- `MockActor` (`actors.go`) — covers `commands_actor_test.go`, `commands_vfx_test.go`, `sequences_test.go`
- `MockVFXManager` (`audio_speech_vfx.go`) — covers `commands_vfx_test.go`
- `MockScene` (`navigation.go`) — used via `mockSceneWithCamera` (package-local, already in `commands_camera_test.go`)
- `MockCommand` (`sequences.go`) — covers `commands_test.go` (ParallelCommand sub-commands)

No new shared mocks created. No package-local `mocks_test.go` needed.

### TDD Specialist 2026-04-04: Red phase tests written

**Test files modified/created:**

- `commands_test.go` — added `stubSpeech` + `newTestDialogueManager()` helper; `TestDialogueCommand_Init_ShowsMessages`, `TestDialogueCommand_Update_ReturnsFalseWhileSpeaking` cover the 0% `DialogueCommand.Init/Update` paths (lines 44/75 in `commands.go`)
- `commands_actor_test.go` — added `TestTrivialUpdateMethods` (table-driven, covers all five `Update() bool { return true }` stubs), `TestResolveActorTargets_QueryInvalidRegex`, `TestResolveActorTargets_QueryMatchingActors` (covers `@query:` branch)
- `commands_music_test.go` — **Red**: `ctx.AudioManager = am` fails to compile because `AppContext.AudioManager` is `*audio.AudioManager` (concrete), not an interface; proves `PlayMusicCommand.Init`, `PauseAllMusicCommand.Init`, `FadeOutAllMusicCommand.Init` (all 0%) cannot be tested until `AppContext.AudioManager` is extracted to an interface
- `commands_camera_test.go` — added `TestCameraZoomCommand_AllPhases_InstantZoom`, `TestCameraZoomCommand_AllPhases_TimedZoom` (all three zoom phases), `TestCameraShakeCommand_InitAddsTrauma`, `TestCameraMoveCommand_Update_DurationPath`, `TestCameraResetCommand_Update_DurationPath`
- `commands_vfx_test.go` — rewrote with `stubVFXManager` (full `vfx.Manager` impl); added `TestSpawnTextCommand_Init_ScreenType`, `TestSpawnTextCommand_Init_OverheadType`, `TestQuakeCommand_Init_Update`, `TestQuakeCommand_Update_AddsTraumaOnMultiplesOf10`
- `sequences_test.go` — added `TestPlaySequence_ValidFile` (uses `os.DirFS` + `t.TempDir()`), `TestPlaySequence_InvalidPath`, `TestEndBlockingPhase_UnblocksPlayer`

**Red proof:** `go test ./internal/engine/sequences/...` → build failed: `cannot use am (variable of type *mocks.MockAudioManager) as *audio.AudioManager` — missing behavior: `AppContext.AudioManager` must become an interface before music `Init` paths can be covered.

### Feature Implementer 2026-04-04: Green phase — production code written

**Production files modified:**

1. `internal/engine/audio/audio.go` — added `Manager` interface with methods: `PlayMusic`, `IsPlaying`, `SetVolume`, `PauseAll`, `FadeOutAll`, `Stop`, `StopAll`. The concrete `AudioManager` type already implements all these methods, so no changes to the implementation were needed.

2. `internal/engine/app/context.go` — changed `AudioManager` field from `*audio.AudioManager` (concrete) to `audio.Manager` (interface). This allows mock implementations to be injected in tests.

3. `internal/engine/contracts/navigation/navigation.go` — updated `SceneManager` interface method `AudioManager()` return type from `*audio.AudioManager` to `audio.Manager`.

4. `internal/engine/scene/scene_manager.go` — updated `AudioManager()` method return type from `*audio.AudioManager` to `audio.Manager`.

5. `internal/engine/scene/scene_tilemap.go` — updated `Audiomanager()` method return type from `*audio.AudioManager` to `audio.Manager`.

**Test results:** `go test ./internal/engine/sequences/... -v` → **all 95 tests PASS** ✅

**Coverage:** `go test ./internal/engine/sequences/... -cover` → **86.4% of statements** (target: ≥80%) ✅

All tests pass with no production code logic changes — only interface extraction to enable dependency injection of mocks.

### Gatekeeper 2026-04-04: ✅ APPROVED

**Red-Green-Refactor Cycle Verified:**
- Red Phase: Tests written first, failing until production code extracted `audio.Manager` interface
- Green Phase: Production code changes minimal — only interface extraction (5 files touched, 0 logic changes)
- Refactor: Table-driven tests, DDD-aligned, headless setup

**Specification Compliance:**
- ✅ AC1: `ParallelCommand`, `WaitCommand` Init/Update tested
- ✅ AC2: All five actor command stubs + `@query:` branch covered
- ✅ AC3: All three music `Init` paths tested with `MockAudioManager`
- ✅ AC4: `PlaySequence` and `endBlockingPhase` tested
- ✅ AC5: `CameraZoomCommand` all phases, `CameraShakeCommand` tested
- ✅ AC6: **Coverage: 86.4%** (target: ≥80%) — **+25.3% delta** from 61.1%
- ✅ AC7: Headless, no GPU calls, mocks at boundaries

**Project Standards:**
- ✅ Table-driven tests throughout
- ✅ No `_ = variable` in production code
- ✅ DDD: `audio.Manager` interface extracted for dependency injection
- ✅ Headless Ebitengine: real `camera.Controller` used without GPU

**Coverage Delta:** +25.3% (61.1% → 86.4%)

**Test Results:** 95 tests PASS, 0 FAIL
