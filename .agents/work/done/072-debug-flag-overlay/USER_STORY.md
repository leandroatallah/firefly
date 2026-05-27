# USER STORY — 072-debug-flag-overlay

**Branch:** `072-debug-flag-overlay`
**Bounded Context:** Engine (`internal/engine/debug/`, `internal/engine/ui/`, `internal/engine/app/`)

---

## Story

As a developer,
I want a runtime overlay toggled with F1 that lists all registered boolean debug flags and lets me toggle them with keyboard input,
so that I can enable or disable debug channels without restarting the game or editing files.

---

## Acceptance Criteria

- AC-1: `internal/engine/debug/` gains a registry: `Register(name string, ptr *bool)` stores the pointer under the given name; `List() []Entry` returns a name-sorted snapshot where `Entry` is `{Name string; Ptr *bool}`.
- AC-2: Existing `debug.Init` / `debug.InitFromReader` populate the registry by calling `Register` for each key loaded from JSON, binding the pointer stored in the channels map.
- AC-3: `internal/game/app/config.go` calls `debug.Register("cam_debug", &cfg.CamDebug)` and `debug.Register("collision_box", &cfg.CollisionBox)` after the `flag.BoolVar` calls so CLI-set values are preserved.
- AC-4: `internal/engine/debug/` gains `Reset()` clears the registry alongside the existing channel map (no registry leak between tests).
- AC-5: `internal/engine/ui/` gains `DebugOverlay` struct (new file `internal/engine/ui/debugoverlay/overlay.go`) with `Update() bool` and `Draw(*ebiten.Image)` methods.
- AC-6: `DebugOverlay.Update()` handles: Up/Down move cursor, Space or Enter toggle the flag under the cursor (writes through `*bool`), F1 closes the overlay; returns `true` while the overlay is open (signals caller to suppress scene update).
- AC-7: `DebugOverlay.Draw` renders: a semi-transparent background panel, then for each entry from `debug.List()` a line `[x] flag_name` or `[ ] flag_name`; the cursor row is highlighted.
- AC-8: `internal/engine/app/engine.go` wires the overlay: `Game` holds a `*ui/debugoverlay.DebugOverlay`; F1 in `Update()` shows/hides it; when overlay is open, `SceneManager.Update()` and `DialogueManager.Update()` are skipped.
- AC-9: `DebugOverlay` has no import of `internal/game/` or `internal/kit/`; it may import `internal/engine/debug/` and `internal/engine/`.
- AC-10: Table-driven unit tests for `Register` / `List`: empty registry returns empty slice; single entry round-trips name and pointer; multiple entries return alphabetically sorted slice.
- AC-11: Table-driven unit tests for `DebugOverlay.Update()`: Up wraps from top to bottom; Down wraps from bottom to top; Space toggles the pointed-to bool; F1 returns false (closed).

---

## Behavioral Edge Cases

- `Register` called with the same name twice: second call overwrites the pointer (last-write-wins); no panic.
- `List()` on an empty registry returns a non-nil empty slice (safe to range over).
- Overlay opened with zero registered flags: draws only the panel, no crash.
- Cursor position is clamped after a `Reset()` + re-registration cycle so it never exceeds `len(entries)-1`.
- `debug.Reset()` must clear the registry; calling `List()` after `Reset()` returns an empty slice.
- CLI flags (`CamDebug`, `CollisionBox`) are registered after `flag.Parse()` runs, so `Register` sees the parsed value; toggling in the overlay writes through the pointer and takes effect immediately.
- The overlay does not persist changes to `assets/data/debug.json` (in-memory only).
