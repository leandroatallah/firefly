# US-007 — Scene-Level Freeze Frame

**Branch:** `007-scene-freeze-frame`
**Bounded Context:** Scene / Entity

## Story

As a game designer,
I want to trigger a brief full-scene freeze (hit-stop) on impactful events like parries or boss hits,
So that the game communicates weight and impact through a satisfying pause in gameplay.

## Acceptance Criteria

- AC1: A `FreezeFrame(durationFrames int)` method on the Scene (or injected via a contract) pauses all Actor and Body updates for the given number of frames.
- AC2: While frozen, `Actor.Update()` and `Body.Update()` are skipped; `Draw()` continues normally.
- AC3: The freeze resolves automatically after the specified duration — no manual reset required.
- AC4: Calling `FreezeFrame` while already frozen resets the timer to the new duration (latest call wins).
- AC5: The freeze state is accessible via a `IsFrozen() bool` query on the Scene.
- AC6: The freeze is triggered through a contract/interface, not by direct Scene reference, so it is injectable in tests.

## Edge Cases

- `durationFrames <= 0`: no freeze is applied.
- Freeze triggered during a Sequence: the Sequence's frame counter also pauses.

## Notes

- Depends on US-004 and US-006 being stable (freeze interacts with dash and state transitions).
- Lives in `internal/engine/scene/` with a contract in `internal/engine/contracts/`.
