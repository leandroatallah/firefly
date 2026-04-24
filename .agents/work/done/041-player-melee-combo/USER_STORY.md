# US-041 — Player Melee Combo Chain

**Branch:** `041-player-melee-combo`
**Bounded Context:** Game (`internal/game/entity/actors/states/`)

## Story

As a player,
I want to chain up to 3 melee attacks in quick succession,
so that combat feels more dynamic and rewarding than a single strike.

## Context

US-040 establishes a single-hit melee attack with a swing hitbox and cooldown.
This story extends the melee system with a combo chain: pressing Z again within a combo window after each hit triggers the next attack in the sequence, each with its own hitbox parameters and animation.

## Acceptance Criteria

- **AC1** — Up to 3 combo hits can be chained by pressing Z within a configurable `combo_window_frames` after each hit completes.
- **AC2** — Each combo step has independent JSON-configurable parameters: `damage`, `active_frames`, `hitbox` (width, height, offset_x, offset_y).
- **AC3** — The combo chain resets to step 1 if:
  - The player does not press Z within `combo_window_frames` after the previous hit.
  - The player takes damage during any combo step.
  - The player begins a dash or jump mid-combo.
- **AC4** — After the final (3rd) hit, the chain always resets — no looping.
- **AC5** — Each combo step plays a distinct animation (e.g. `MeleeAttack1`, `MeleeAttack2`, `MeleeAttack3`).
- **AC6** — Combo configuration lives in the player JSON under a `combo_steps` array within the melee weapon block:
  ```json
  {
    "type": "melee",
    "combo_window_frames": 15,
    "combo_steps": [
      { "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 24, "height": 16, "offset_x": 12, "offset_y": 0 } },
      { "damage": 1, "active_frames": [3, 8],  "hitbox": { "width": 28, "height": 16, "offset_x": 14, "offset_y": -4 } },
      { "damage": 2, "active_frames": [5, 12], "hitbox": { "width": 32, "height": 20, "offset_x": 16, "offset_y": 0 } }
    ]
  }
  ```
- **AC7** — Unit tests verify:
  - Pressing Z within `combo_window_frames` advances to the next step.
  - Missing the window resets to step 1.
  - Taking damage mid-combo resets to step 1.
  - After step 3, chain resets.
  - Each step uses its own hitbox and damage values.

## Proposed Changes

- `internal/game/entity/actors/states/melee_state.go` — extend `StateMeleeAttack` with combo step tracking
- `internal/engine/combat/weapon/melee.go` — extend `MeleeWeapon` to hold `[]ComboStep` and `comboWindowFrames`
- `internal/engine/combat/weapon/factory.go` — parse `combo_steps` array from JSON
- Player JSON config — replace single-hit melee block with `combo_steps` array

## Dependencies

- US-040 — Player Melee Attack
