# Design Decision Summary — Explicit Shooting States

**Date:** 2026-04-02T12:29  
**Story:** USER STORY 011 — Refactor Shooting to Explicit Actor States

---

## Problem

The initial spec proposed moving `ShootingSkill` to the engine layer as a skill that modifies sprite lookup (composite sprite keys like "idle_shoot"). However, this approach had trade-offs:

**Option 1: Composite Sprite Keys (Initial Approach)**
- Sprite keys: "idle", "idle_shoot", "walk", "walk_shoot"
- State machine: Idle, Walking, Jumping (unchanged)
- Shooting: Concurrent skill that modifies sprite lookup

**Cons:**
- Sprite explosion: N states × M skills = N×M sprite variants
- Implicit coupling: sprite naming convention must match state+skill combo
- Harder to reason about: "What sprite am I showing?" requires checking state + all active skills
- Animation timing complexity: base state controls frame rate, but shooting might need different timing

---

## Solution: Explicit Shooting States (Chosen Approach)

**Option 2: Explicit Shooting States**
- State machine: Idle, Walking, IdleShooting, WalkingShooting, JumpingShooting, etc.
- Shooting: Triggers state transition (Idle → IdleShooting)
- Sprite mapping: 1:1 (IdleShooting → "idle_shoot" → idle_shoot.png)

**Pros:**
- **Explicit and clear:** Each visual state is a real state
- Sprite mapping is 1:1 (no composite keys)
- Animation timing per state (idle shooting can have different frame rate than idle)
- Easier to debug: state machine shows exactly what's happening
- Follows Cuphead's actual implementation (they likely use separate states)
- Matches existing architecture: `DashState` is a separate state, not a movement modifier

---

## Key Insight

User feedback: *"When the player is idle and shoot, it changes its sprite to a different. It isn't idle anymore."*

This is the critical point: **shooting changes the visual state**. It's not "idle with a modifier," it's "idle shooting" — a distinct state with its own:
- Sprite sheet
- Animation timing
- Possibly different hitboxes
- Different transition rules

---

## Implementation

1. **Register shooting state variants** in `actor_state.go`:
   - `IdleShooting`, `WalkingShooting`, `JumpingShooting`, `FallingShooting`

2. **ShootingSkill triggers state transitions**:
   - Shoot pressed + Idle → IdleShooting
   - Shoot released + IdleShooting → Idle
   - Movement input + IdleShooting → WalkingShooting

3. **Sprite system maps states to sprite sheets**:
   - `IdleShooting` → `"idle_shoot"` → `idle_shoot.png`
   - No composite keys, no implicit coupling

4. **State machine handles transitions**:
   - `Character.handleState()` includes logic for shooting state transitions
   - Shooting states follow the same movement rules as their base states

---

## Future Extensions

This design supports future shooting variants:
- `IdleShootingUp` (shooting straight up while idle)
- `WalkingShootingDiagonal` (shooting diagonally while walking)
- `DuckingShooting` (shooting while ducking)

Each variant is a new state with its own sprite sheet and transition rules.

---

## Conclusion

Explicit shooting states provide clarity, maintainability, and alignment with Cuphead's actual behavior. The state machine accurately reflects what the actor is doing, and the sprite system maps states to sprites without implicit coupling.
