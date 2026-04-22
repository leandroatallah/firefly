# User Story 042 — Weapon Hitbox Debug Visualization

## Story

**As a** developer debugging combat interactions in a Phase,
**I want** the `--collision-box` debug flag to also render projectile collision boxes and the active melee hitbox,
**so that** I can visually verify weapon hit detection aligns with the game world without needing a separate tool or flag.

## Background

The Phase Scene already draws collision boxes for Actors and obstacles when `--collision-box` is enabled. Projectiles fired from the player's weapon and the melee swing hitbox are invisible under this flag, making it difficult to diagnose missed or incorrect hits during combat.

## Acceptance Criteria

### AC-1: Projectile Collision Boxes Are Visible Under `--collision-box`

**Given** the game is running with `--collision-box` enabled,
**And** one or more projectiles are active in the current Phase,
**When** the Phase renders each frame,
**Then** each active projectile's collision box is drawn using the existing non-obstructive (green) debug style — the same visual as a non-obstructive Actor collision box.

### AC-2: Active Melee Hitbox Is Visible Under `--collision-box`

**Given** the game is running with `--collision-box` enabled,
**And** the player's Actor is performing a melee attack (the hitbox is active during the swing),
**When** the Phase renders each frame during the active swing frames,
**Then** the melee hitbox rectangle is drawn in a distinct orange color to differentiate it from Actor and projectile boxes.

### AC-3: Melee Hitbox Disappears When Swing Is Inactive

**Given** the game is running with `--collision-box` enabled,
**And** the player's Actor is NOT currently performing a melee swing (no active hitbox),
**When** the Phase renders each frame,
**Then** no melee hitbox rectangle is drawn — the overlay is frame-accurate to the hitbox active window.

### AC-4: No Debug Rendering Without the Flag

**Given** the game is running WITHOUT `--collision-box`,
**When** projectiles are active or a melee swing is in progress,
**Then** no projectile collision boxes and no melee hitbox rectangles are drawn on screen.

### AC-5: Debug Rendering Does Not Affect Game Logic

**Given** the `--collision-box` flag is enabled or disabled,
**When** any combat interaction occurs (projectile hit, melee hit),
**Then** hit detection, damage application, and Actor State transitions are identical in both cases — debug rendering is purely visual with no side effects on game logic.

## Visual Contract

| Shape | Color | Condition |
|---|---|---|
| Actor collision boxes | Red (obstructive) / Green (non-obstructive) | Existing behavior, unchanged |
| Projectile collision boxes | Green (non-obstructive style) | Each active projectile, every frame |
| Melee hitbox | Orange | Only during active swing frames |

## Out of Scope

- No new debug flags or CLI options are introduced; this extends `--collision-box` only.
- Enemy or NPC melee hitboxes are not in scope for this story.
- Rendering performance under high projectile counts is not a concern for this debug-only path.

## DDD Glossary (terms used in this story)

| Term | Meaning |
|---|---|
| Actor | The player entity with a state machine |
| Phase | The playable game level Scene |
| Scene | The self-contained game state being drawn |
| State | A named node in the Actor's state machine (e.g., melee swing active) |
| Body | The physics body carrying position; projectiles are Bodies |
