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

The hitbox active-frame check lives in `MeleeWeapon.IsHitboxActive()` (`internal/kit/combat/weapon/melee.go`). `melee.State` reaches it only through the public `weaponIface`. The override mechanism must therefore be: `State.OnStart` extracts `HitboxFrames` from the incoming step's `AssetData` and calls a new `SetActiveFramesOverride(*[2]int)` method on the weapon; `MeleeWeapon` uses the override in `IsHitboxActive` when set, and falls back to `ComboStep.ActiveFrames` when nil; the override is cleared at the start of each new swing so it never leaks between steps or actors.

---

## Acceptance Criteria

- AC-1: `schemas.AssetData` gains an optional `HitboxFrames *HitboxFrameRange` field tagged `json:"hitbox_frames,omitempty"`; absent field parses without error and leaves the pointer nil.
- AC-2: `HitboxFrameRange` is a new struct in `internal/engine/data/schemas/` with fields `Start int` (`json:"start"`) and `End int` (`json:"end"`).
- AC-3: `weapon.MeleeWeapon` gains `SetActiveFramesOverride(override *[2]int)`; when non-nil, `IsHitboxActive` uses `override[0]`/`override[1]` instead of `ComboStep.ActiveFrames`; when nil, behaviour is unchanged.
- AC-4: `melee.State.OnStart` calls `weapon.SetActiveFramesOverride` with the `[2]int` derived from the incoming step's `AssetData.HitboxFrames`, or nil if `HitboxFrames` is absent; the override is set before `weapon.Fire` is called.
- AC-5: `melee.State` receives the asset map for the current actor (`map[string]AssetData`) and a state-name resolver (`func(stepIdx int) string`) at construction or via a setter, so `OnStart` can look up the correct `AssetData` by step index without importing `internal/game/`.
- AC-6: When `HitboxFrames` is nil for the current step, `ComboStep.ActiveFrames` governs hitbox activation (zero regression).
- AC-7: When `HitboxFrames` is present, its `[Start, End]` range replaces `ComboStep.ActiveFrames` for hitbox activation only; damage, startup frames, hitbox rect, and all other weapon fields are unaffected.
- AC-8: The override is cleared (set to nil) at the start of each swing in `MeleeWeapon.startSwing`, preventing stale values from leaking to subsequent combo steps or a different actor sharing the weapon instance.
- AC-9: Layer rules upheld: `internal/engine/data/schemas/` does not import `internal/kit/` or `internal/game/`; the override is consumed only inside `internal/kit/combat/`.
- AC-10: No existing actor JSON files require modification; all omit `hitbox_frames` and continue to behave identically.
- AC-11: Table-driven unit tests in `weapon` package cover: nil override uses `ComboStep.ActiveFrames`; non-nil override uses JSON range; frame exactly at `Start` and exactly at `End` are both active; `Start > End` means hitbox never activates; `Start == End` is a valid single-frame window.
- AC-12: Table-driven unit tests in `melee` package cover: `State.OnStart` sets override from `AssetData.HitboxFrames` when present; sets nil override when absent; override is cleared before each new swing fires.

---

## Behavioral Edge Cases

- `Start > End`: treat as invalid config; hitbox never activates for that step (weapon sees override but `swingFrame` never satisfies the range).
- `Start == End`: single-frame active window; must activate on exactly that frame.
- Multi-step combo: step 0 may have `hitbox_frames`, step 1 may not; each `OnStart` call re-reads its own step's `AssetData` independently.
- Non-melee states that declare `hitbox_frames`: field is silently ignored outside the melee state machine.
- State transition mid-combo: `OnStart` always sets or clears the override from the incoming step's `AssetData`; no stale value from the previous step survives.
- `MeleeWeapon` shared across actors: `startSwing` clears the override so whichever actor fires next installs its own value via `OnStart`.
