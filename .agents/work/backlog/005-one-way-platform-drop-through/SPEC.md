# SPEC 005 — One-Way Platform Drop-Through

**Branch:** `005-one-way-platform-drop-through`
**Bounded Context:** Physics / Entity
**Package:** `internal/engine/physics/movement/` + platform body

## Technical Requirements

### 1. New contract method on one-way platform body

Extend (or add) a `OneWayPlatform` interface in `internal/engine/contracts/body/`:

```go
type OneWayPlatform interface {
    Body
    IsOneWay() bool
    SetPassThrough(actor Collidable, frames int) // disables solidity for `actor` for `frames` frames
    IsPassThrough(actor Collidable) bool
}
```

- `SetPassThrough` registers the actor with a frame countdown.
- Each `Update()` tick on the platform decrements all countdowns; expired entries are removed.
- While `IsPassThrough(actor)` is true, `ResolveCollisions` skips this platform for that actor.

### 2. Drop-through trigger — movement handler or input handler

In the Actor's movement/input update (grounded phase):

```go
func tryDropThrough(actor MovableCollidable, platform OneWayPlatform, input InputSource)
```

- Condition: grounded on a `OneWayPlatform` AND `input.DropHeld()` (down + jump simultaneously).
- Action: `platform.SetPassThrough(actor, minFrames=2)`, do NOT set upward velocity.
- Vertical velocity is left at `0` or natural gravity — no jump impulse.

### 3. Collision resolution guard

In `BodiesSpace.ResolveCollisions` (or the platform's `OnTouch`):
- Before treating a one-way platform as solid from above, check `platform.IsPassThrough(actor)`.
- If true → skip collision response for this pair.

### Input disambiguation (vs. duck, SPEC-003)

- Down alone → duck (SPEC-003).
- Down + jump simultaneously → drop-through (this spec); duck is NOT entered.
- Jump alone on one-way platform → normal jump (no drop-through).

## Pre-conditions

- Actor is grounded on a `OneWayPlatform` (`IsOneWay() == true`).
- Down + jump input pressed in the same frame.

## Post-conditions

- Actor passes through the platform (no upward velocity applied).
- After actor's bottom edge clears the platform's top edge, `IsPassThrough` returns false and normal collision resumes.
- Solid (two-way) platforms are unaffected.

## Integration Points

- New `OneWayPlatform` interface in `internal/engine/contracts/body/` (new file `one_way_platform.go`).
- `BodiesSpace.ResolveCollisions` reads `OneWayPlatform.IsPassThrough`.
- `body.MovableCollidable` from existing contracts.
- Input disambiguation touches the same input path as SPEC-003 (duck); must be resolved before duck state is entered.

## Red Phase — Failing Test Scenario

File: `internal/engine/physics/movement/drop_through_test.go`

`TestDropThrough`:

| case | setup | expected |
|---|---|---|
| down+jump on one-way | grounded on OneWayPlatform, dropHeld=true | IsPassThrough==true, vy unchanged (no jump) |
| down alone on one-way | grounded, duckHeld=true, jumpHeld=false | IsPassThrough==false (duck, not drop) |
| jump alone on one-way | grounded, jumpHeld=true, downHeld=false | IsPassThrough==false, normal jump |
| down+jump on solid | grounded on solid body | IsPassThrough not applicable, no drop |
| pass-through expires | SetPassThrough(actor, 2), tick twice | IsPassThrough==false after 2 ticks |

Test must fail (interface and logic do not exist yet) → implement → test passes.
