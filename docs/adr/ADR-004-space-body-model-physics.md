# ADR-004 — Space / Body / MovementModel Physics Layers

## Status
Accepted

## Context
A monolithic physics object that owns position, collision, and movement logic becomes hard to test and reuse. Movement rules differ significantly between game types (top-down vs. platformer), and collision resolution must be independent of how an entity moves. Mixing these concerns forces every test to construct a full physics object even when only one aspect is under test.

## Decision
Physics is split into three distinct layers:

| Layer | Package | Responsibility |
|---|---|---|
| `Body` | `physics/body` | Owns position (x16/y16), shape, and collision callbacks (`OnTouch`, `OnBlock`). |
| `Space` | `physics/space` | Holds all active bodies, resolves collisions between them via `ResolveCollisions`. |
| `MovementModel` | `physics/movement` | Applies velocity and calls `Space.ResolveCollisions` each tick. Swappable via `MovementModelEnum` (`TopDown`, `Platform`). |

`Body` implements `body.Collidable`; `Space` implements `body.BodiesSpace`; `MovementModel` depends only on the `body.MovableCollidable` and `body.BodiesSpace` interfaces.

## Consequences
- Each layer is independently testable with mocks of the adjacent interfaces.
- Switching movement behaviour (e.g. top-down to platformer) requires only changing the `MovementModel`; the body and space are unchanged.
- Adding a new movement model means implementing `MovementModel` and registering it in `NewMovementModel` — no changes to `Body` or `Space`.
- The three-layer split adds indirection; understanding a full physics tick requires reading across three packages.
