# USER STORY — 056-eight-dir-movement-skill

**Branch:** `056-eight-dir-movement-skill`
**Bounded Context:** Kit (`internal/kit/skills/`)

---

## Story

As a kit developer,
I want an `EightDirectionalMovementSkill` that translates Up/Down/Left/Right input commands into body acceleration calls,
so that any actor using a compatible movement model can move in 8 directions without coupling the skill to a specific genre.

---

## Background

`HorizontalMovementSkill` (`internal/kit/skills/platform_move.go`) is the existing pattern: it wraps `skill.SkillBase`, reads `input.CommandsReader()`, respects `Immobile()` and `IsInputBlocked()`, and drives a `body.MovableCollidable` via `OnMoveLeft/Right`. Beat-em-up (and future top-down) actors need the same pattern extended to all four axes.

---

## Constraints (resolved in grilling session 2026-05-09)

- **Y axis is ground-plane depth**, not altitude. `OnMoveUp`/`OnMoveDown` drive `y16`. Altitude is never touched by this skill.
- **Skill is the sole input reader** for movement. `BeatEmUpMovementModel` (057) is passive; it never reads keys. This follows the platformer skill architecture, not the top-down embedded-input pattern.
- **Genre-agnostic**: no import of any beat-em-up, platformer, or top-down package.
- **Selected via `cfg.Movement.Mode == "eight_dir"`** in `kitskills.FromConfig` (see story 058).

---

## Acceptance Criteria

- AC-1: `EightDirectionalMovementSkill` lives in `internal/kit/skills/eight_dir_move.go`, embeds `skill.SkillBase`, and compiles as part of package `kitskills`.
- AC-2: `HandleInput` calls `body.OnMoveLeft(speed)`, `OnMoveRight(speed)`, `OnMoveUp(speed)`, `OnMoveDown(speed)` on the active directions read from `input.CommandsReader()` (Left/Right/Up/Down commands). `OnMoveUp`/`OnMoveDown` drive `y16` (ground-plane depth); altitude is never written.
- AC-3: When `body.Immobile()` is true, `HandleInput` zeroes X and Y velocity components and returns without calling any `OnMove*` method.
- AC-4: When the movement model signals `IsInputBlocked()`, `HandleInput` returns immediately (same guard as `HorizontalMovementSkill`).
- AC-5: `Update` delegates to `SkillBase.Update` and is otherwise a no-op.
- AC-6: The skill has no import of any movement model package, beat-em-up package, or platformer package — it accepts any `body.MovableCollidable`.
- AC-7: Table-driven unit tests cover: move left only, move right only, move up only, move down only, diagonal (left+up), immobile guard, input-blocked guard, no input (no `OnMove*` called).

---

## Behavioral Edge Cases

- Diagonal input (e.g., Left + Up simultaneously): both `OnMoveLeft` and `OnMoveUp` are called; speed normalization is the movement model's responsibility.
- Immobile body: both vx16 and vy16 are zeroed; `OnMove*` must not be called.
- Simultaneous Left + Right (or Up + Down): both calls fire; conflict resolution deferred to the model.
- Guard order (`IsInputBlocked` vs `Immobile`): defer to spec engineer — confirm ordering is safe for all callers.
