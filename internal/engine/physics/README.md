# Physics Engine

The `physics` package provides a high-performance, deterministic 2D physics engine designed for both top-down and platformer gameplay. It prioritizes precision, stability, and developer flexibility through a decoupled architecture.

## Table of Contents

- [Core Architecture](#core-architecture)
- [Subpackages](#subpackages)
  - [Body: Entities & Shapes](#body-entities--shapes)
  - [Space: Collision Resolution](#space-collision-resolution)
  - [Movement: Physical Models](#movement-physical-models)
  - [Skill: Ability Integration](#skill-ability-integration)
- [Sophisticated Mechanics](#sophisticated-mechanics)
  - [Fixed-Point Arithmetic (fp16)](#fixed-point-arithmetic-fp16)
  - [Sub-Pixel Accumulation](#sub-pixel-accumulation)
  - [Continuous Collision Detection (Sweep)](#continuous-collision-detection-sweep)
- [Performance & Optimization](#performance--optimization)
- [Agent Quick-Reference](#agent-quick-reference)

---

## Core Architecture

The engine is built on a "Space-Body-Model" pattern:

1.  **Space**: A central registry that manages all physical entities and orchestrates collision detection.
2.  **Body**: Data-heavy structures representing physical properties (position, velocity, shapes).
3.  **Movement Model**: Logic-heavy components that apply forces and integrate movement based on specific game rules (e.g., Gravity, Friction).

---

## Subpackages

### Body: Entities & Shapes

[body](internal/engine/physics/body) manages physical representations.

- **Non-Obvious Pattern**: It uses an **Ownership Hierarchy**. A body can "own" other bodies. The [Ownership.LastOwner()](internal/engine/physics/body/ownership.go) method performs an iterative traversal with cycle detection to find the root entity (e.g., finding the `Character` that owns a sword's collision box).
- **Movable vs Collidable**: Logic is split into [MovableBody](internal/engine/physics/body/body_movable.go) (velocity/acceleration) and [CollidableBody](internal/engine/physics/body/body_collidable.go) (collision shapes/resolution).
- **Hitbox Resize Utilities**: `ResizeFixedBottom(rect, newHeight)` and `ResizeFixedTop(rect, newHeight)` return a new rect with the same bottom or top edge respectively, used by duck and dash states to shrink the hitbox without repositioning the actor.

### Space: Collision Resolution

[space](internal/engine/physics/space) is the "world" where physics happen.

- **Concurrency**: Uses `sync.RWMutex` for thread-safe access to bodies.
- **Optimization**: Maintains a sorted [bodiesCache](internal/engine/physics/space/space.go) to speed up `Bodies()` queries, invalidated only when bodies are added or removed.
- **Dynamic Collisions**: The [StateCollisionManager](internal/engine/physics/space/state_collision_manager.go) allows entities to swap collision shapes based on their animation state (e.g., a "crouch" state reducing the collision height).

### Movement: Physical Models

[movement](internal/engine/physics/movement) implements the "how" of movement.

- **Platformer Logic**: [PlatformMovementModel](internal/engine/physics/movement/movement_model_platform.go) implements asymmetric gravity (higher downward gravity for a "snappier" feel) and air control multipliers.
- **Top-Down Logic**: [TopDownMovementModel](internal/engine/physics/movement/movement_model_topdown.go) includes diagonal normalization to prevent the "diagonal speed boost" pitfall.
- **One-Way Platform Drop-Through**: `tryDropThrough` disables a `OneWayPlatform`'s solidity for the actor for a minimum of 2 frames when down + jump are pressed simultaneously. Input disambiguation: down alone → duck; down + jump → drop-through; jump alone → normal jump.

### Skill: Ability Integration

[skill](internal/engine/physics/skill) encapsulates discrete player abilities.

- **`JumpSkill`**: Handles jump activation and variable jump height. Calls `body.TryJump(force)` on key-press. On early key release while `body.IsGoingUp()`, multiplies upward velocity by `JumpCutMultiplier` (range `(0.0, 1.0]`, default `1.0`). The cut fires at most once per jump. Accepts an injectable `InputSource` interface for deterministic testing.
- **`DashSkill` / `DashState`**: Tween-based dash deceleration using `InOutSineTween`. Velocity follows an `InOutSine` curve from `DashSpeed` to `0` over a configurable frame duration. Gravity is suspended and the hitbox is reduced to duck height for the full dash. One air dash per jump; cooldown prevents immediate re-trigger.

---

## Sophisticated Mechanics

### Fixed-Point Arithmetic (fp16)

To ensure determinism and avoid floating-point jitter, the engine uses fixed-point math (scaled by 16).

- **Usage**: Internal positions and velocities are stored as `x16`, `y16`.
- **Conversion**: Use [fp16.To16()](internal/engine/utils/fp16/fp16.go) to convert pixels to internal units and `From16()` to convert back for rendering.

### Sub-Pixel Accumulation

The [CollidableBody.ApplyValidPosition](internal/engine/physics/body/body_collidable.go) method uses **Accumulators** (`accumulatorX16`, `accumulatorY16`).

- **Why**: If an object moves 0.4 pixels per frame, it shouldn't move for 2 frames and then jump 1 pixel on the 3rd. The accumulator stores this fractional movement until it reaches a full pixel threshold.

### Continuous Collision Detection (Sweep)

Instead of teleporting objects to their new position and checking for overlap (which causes "tunneling" through thin walls), the engine performs a **pixel-by-pixel sweep**:

```go
// Simplified logic from ApplyValidPosition
for i := 0; i < pixelSteps; i++ {
    b.SetPosition16(lastX16 + step16, lastY16)
    if _, blocking := space.ResolveCollisions(b); blocking {
        b.SetPosition16(lastX16, lastY16) // Revert
        break
    }
}
```

---

## Performance & Optimization

- **Two-Axis Separation**: Movement is resolved on X then Y independently. This allows "sliding" along walls (e.g., moving diagonally into a vertical wall still allows vertical movement).
- **Squared Magnitude**: When clamping max speed in top-down movement, the engine compares squared values (`velSq > maxSq`) to avoid the expensive `math.Sqrt` call.
- **Callback Pattern**: Collision triggers `OnTouch` and `OnBlock` callbacks on both involved bodies, allowing for decoupled event handling (e.g., a "HurtBox" body telling its owner to take damage).

---

## Agent Quick-Reference

### Critical Implementation Details

- **Determinism**: Never use `float64` for positions or velocities in production code. Use `fp16`.
- **Inertia**: If `HorizontalInertia` is set to `-1`, acceleration is treated as instant velocity. If `> 0`, it uses friction and acceleration curves.
- **Grounded State**: `PlatformMovementModel` applies a "sticking force" (`DownwardGravity - 1`) when grounded to prevent micro-bouncing on slopes/platforms.
- **Blocking**: A collision is only "blocking" if `other.IsObstructive()` is true. `OnTouch` triggers regardless of obstructiveness.

### Common Pitfalls

- **ID Requirement**: All bodies added to `Space` **must** have a unique ID.
- **Shape Support**: Currently, `SetPosition` and movement logic are optimized for **Rectangular Shapes**. Circular or complex polygons may require extending `Shape` and `ResolveCollisions`.
- **Frozen State**: If `body.Freeze()` is true, all physics updates for that body are skipped. Check this first if an object isn't moving.

### Example: Creating a Movable Collidable

```go
// 1. Define Shape
rect := body.NewRect(16, 16)

// 2. Create Body
baseBody := body.NewBody(rect)
movable := body.NewMovableBody(baseBody)
collidable := body.NewCollidableBody(baseBody)

// 3. Configure
collidable.SetID("player_01")
collidable.SetIsObstructive(true)
movable.SetMaxSpeed(4)

// 4. Register
space.AddBody(collidable)
```
