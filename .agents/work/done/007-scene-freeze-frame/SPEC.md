# SPEC-007 — Scene-Level Freeze Frame

**Branch:** `007-scene-freeze-frame`
**Bounded Context:** Scene / Entity
**Package:** `internal/engine/scene/` + `internal/engine/contracts/scene/`

## Technical Requirements

### 1. New contract — `internal/engine/contracts/scene/freeze.go`

```go
package scene

type Freezable interface {
    FreezeFrame(durationFrames int)
    IsFrozen() bool
}
```

### 2. `FreezeController` — `internal/engine/scene/`

```go
type FreezeController struct { remaining int }

func (f *FreezeController) FreezeFrame(durationFrames int)
func (f *FreezeController) IsFrozen() bool
func (f *FreezeController) Tick()   // called once per frame by the Scene; decrements remaining
```

- `FreezeFrame(n)`: if `n <= 0`, no-op. Otherwise `f.remaining = n` (latest call wins / resets timer).
- `IsFrozen()`: returns `f.remaining > 0`.
- `Tick()`: if `remaining > 0`, decrement by 1.

### 3. Scene integration

The `Scene` (or `Phase`) embeds or holds a `*FreezeController`. In `Scene.Update()`:

```go
f.freeze.Tick()
if !f.freeze.IsFrozen() {
    // update actors and bodies
}
// Draw() is always called regardless
```

- `Actor.Update()` and `Body.Update()` are skipped while frozen.
- `Draw()` is never skipped.
- `Sequence` frame counter: the `sequences.Player.Update()` call is also gated behind `!IsFrozen()`.

### 4. Injectable via contract

Scene exposes `Freezable` (not the concrete `*FreezeController`) to game logic, so tests can inject a mock.

## Pre-conditions

- `FreezeFrame` called with `durationFrames > 0`.

## Post-conditions

- For exactly `durationFrames` ticks after the call, `IsFrozen() == true`.
- On tick `durationFrames + 1`, `IsFrozen() == false`.
- Calling `FreezeFrame(n)` while already frozen resets remaining to `n`.
- `durationFrames <= 0` → `IsFrozen()` remains false.

## Integration Points

- New contract file: `internal/engine/contracts/scene/freeze.go`.
- `FreezeController` lives in `internal/engine/scene/freeze.go`.
- `sequences.Player` (existing contract) update is gated by freeze.
- `body.Movable.SetFreeze` is a per-body freeze (SPEC-004); scene-level freeze is separate and additive.
- No changes to existing `navigation.Scene` interface required (freeze is an opt-in embed).

## Red Phase — Failing Test Scenario

File: `internal/engine/scene/freeze_test.go`

`TestFreezeController` (table-driven):

| case | action | frame | IsFrozen() |
|---|---|---|---|
| no freeze | — | 0 | false |
| freeze(3), frame 0 | FreezeFrame(3) | 0 | true |
| freeze(3), frame 1 | Tick() | 1 | true |
| freeze(3), frame 2 | Tick() | 2 | true |
| freeze(3), frame 3 | Tick() | 3 | false |
| freeze(0) | FreezeFrame(0) | 0 | false |
| freeze(-1) | FreezeFrame(-1) | 0 | false |
| reset mid-freeze | FreezeFrame(3), Tick(), FreezeFrame(5) | 1 | true (remaining=5) |

`TestFreezeControllerResetWins`:
- `FreezeFrame(2)`, `Tick()`, `FreezeFrame(4)` → after 4 more ticks `IsFrozen() == false`.

Test must fail (type does not exist yet) → implement → test passes.
