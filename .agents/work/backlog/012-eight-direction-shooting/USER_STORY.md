# USER STORY 012 — 8-Direction Shooting (Cuphead-style)

**Branch:** `012-eight-direction-shooting`
**Bounded Context:** Game Logic (`internal/game/`) + Physics (`internal/engine/physics/skill/`)

## Story

As a player,
I want to aim and shoot in 8 directions (horizontal, up, down, and 4 diagonals),
so that I can attack enemies positioned above, below, or diagonally from my character, just like in Cuphead.

## Context

US-011 established shooting as explicit actor states with horizontal-only shooting. This story extends that architecture to support 8-direction aiming.

In Cuphead, the player can shoot in 8 directions:
1. **Straight** (horizontal left/right)
2. **Up** (straight up)
3. **Down** (straight down, only while jumping/falling)
4. **Diagonal Up-Forward** (45° up + forward)
5. **Diagonal Down-Forward** (45° down + forward, only while jumping/falling)
6. **Diagonal Up-Back** (45° up + backward)
7. **Diagonal Down-Back** (45° down + backward, only while jumping/falling)
8. (Straight back is horizontal with opposite facing)

Direction is determined by directional input (arrow keys / D-pad) while holding the shoot button.

## Acceptance Criteria

- **AC1** — Input system detects 8 directional inputs: none (straight), up, down, up-forward, down-forward, up-back, down-back.
- **AC2** — New shooting state variants registered for each direction:
  - `IdleShootingUp`, `IdleShootingDiagonalUp`, `IdleShootingDiagonalDown`
  - `WalkingShootingUp`, `WalkingShootingDiagonalUp`, `WalkingShootingDiagonalDown`
  - `JumpingShootingUp`, `JumpingShootingDown`, `JumpingShootingDiagonalUp`, `JumpingShootingDiagonalDown`
  - `FallingShootingUp`, `FallingShootingDown`, `FallingShootingDiagonalUp`, `FallingShootingDiagonalDown`
- **AC3** — `ShootingSkill.HandleInput()` reads directional input and transitions to the appropriate directional shooting state.
- **AC4** — Bullet velocity is calculated based on direction:
  - Straight: `(±speedX, 0)`
  - Up: `(0, -speedY)`
  - Down: `(0, +speedY)` (only while airborne)
  - Diagonal: `(±speedX * 0.707, ±speedY * 0.707)` (normalized 45° vector)
- **AC5** — Bullet spawn offset is adjusted per direction (e.g., up shots spawn above the character's head).
- **AC6** — Shooting down is only allowed while jumping or falling (grounded down input triggers ducking, not shooting down).
- **AC7** — Directional state transitions: changing aim direction while shooting transitions between directional shooting states without releasing the shoot button.
- **AC8** — Sprite system maps directional shooting states to distinct sprite sheets (e.g., `"idle_shoot_up.png"`).
- **AC9** — Unit tests cover: all 8 directions, bullet velocity calculation, state transitions between directions, down-shooting restriction while grounded.

## Behavioral Edge Cases

- Holding shoot + pressing up mid-shot must transition from `IdleShooting` → `IdleShootingUp` without resetting cooldown.
- Releasing directional input while shooting must transition back to straight shooting (e.g., `IdleShootingUp` → `IdleShooting`).
- Diagonal input (e.g., up + forward) must take priority over straight up or straight forward.
- Shooting down while grounded must be ignored (ducking takes priority).
- Changing facing direction (left/right) while shooting up must not interrupt the shot.
- Diagonal bullets must travel at the same speed as straight bullets (normalized velocity).

## Notes

- Builds on US-011's explicit shooting state architecture.
- Directional input is read from the same input source as movement (arrow keys / D-pad).
- Diagonal velocity uses `0.707` (≈ 1/√2) to normalize the vector to the same speed as straight shots.
- Sprite sheets required for each directional variant (significant art asset requirement).
- `ShootingSkill` remains in `internal/engine/physics/skill/` but may need game-specific directional logic injected.
- Consider whether diagonal-back shooting is allowed while walking forward (Cuphead allows this).

## Success Criteria

- All 8 directions functional with correct bullet velocity.
- State transitions between directions are smooth (no cooldown reset).
- Sprite system maps all directional states to sprite sheets.
- Down-shooting restricted to airborne states.
- All tests pass with no regressions from US-011.
- Code coverage ≥74.6% (no delta loss).
