# USER STORY — 067-actor-json-hitbox-active-frames

**Branch:** `067-actor-json-hitbox-active-frames`
**Bounded Context:** Kit (`internal/kit/combat/`, `internal/engine/data/schemas/`)

---

## Story

As a beat-em-up game developer,
I want to declare optional `hitbox_frames` per state in an actor's JSON sprite config,
so that each actor's melee hitbox active window matches its own animation timing without changing the shared weapon definition.

---

## Background

`weapon.ComboStep.ActiveFrames [2]int` is defined in Go code and shared across all actors that reference the same weapon. `CodyPlayer` and `Climber` share `NewPlayerMeleeWeapon()`. When the two actors have different animation speeds or art timing, one actor's hitbox fires at the wrong frames. This story introduces a JSON-side override scoped to a single actor state so each actor can independently tune its active window.

---

## Acceptance Criteria

- AC-1: `schemas.AssetData` gains an optional `HitboxFrames *HitboxFrameRange` field tagged `json:"hitbox_frames,omitempty"`; absent field parses without error and leaves the pointer nil.
- AC-2: `HitboxFrameRange` is a new struct in `internal/engine/data/schemas/` with fields `Start int` (`json:"start"`) and `End int` (`json:"end"`).
- AC-3: The beat-em-up melee state machine reads `HitboxFrames` from the current state's `AssetData` each time it evaluates the active window.
- AC-4: When `HitboxFrames` is nil for the current state, the weapon's `ComboStep.ActiveFrames` is used unchanged (zero regression).
- AC-5: When `HitboxFrames` is present, its `[Start, End]` range replaces `ComboStep.ActiveFrames` for hitbox activation logic only; damage, startup frames, and all other weapon fields are unaffected.
- AC-6: Layer rules upheld: `internal/engine/data/schemas/` does not import `internal/kit/` or `internal/game/`; the override is consumed only inside `internal/kit/combat/`.
- AC-7: No existing actor JSON files require modification; all omit `hitbox_frames` and continue to behave identically.
- AC-8: Table-driven unit tests cover: state with `hitbox_frames` activates hitbox within the JSON range; state without `hitbox_frames` activates within `ComboStep.ActiveFrames`; frame exactly at `Start` and exactly at `End` are both treated as active.

---

## Behavioral Edge Cases

- `Start > End`: treat as invalid config; hitbox never activates for that state (Spec Engineer must document the validation path).
- `Start == End`: single-frame active window; must be supported.
- State transition mid-combo: the override is re-read from the incoming state's `AssetData` on each state entry; no stale value from the previous step.
- Non-melee states that happen to declare `hitbox_frames`: field is silently ignored outside the melee state machine.
