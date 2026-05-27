# NOTES — 072-debug-flag-overlay

## Design Choices

- **Registry lives in `internal/engine/debug/`, not a new package.** Keeps the
  flag/state surface in one place and lets `InitFromReader` directly populate
  pointers without a circular dependency.
- **Switch `channels` internal type to `map[string]*bool`.** Go disallows
  taking the address of a map value, so to have stable pointers we either
  store pointers in the map or maintain a parallel map. Storing pointers
  keeps everything coherent: `Enabled()`, `List()`, and the overlay all see
  the same memory. Existing `debug_test.go` uses only the public API so no
  test churn beyond the new registry tests.
- **CLI debug flags (`cam_debug`, `collision_box`) registered at
  `NewConfig()` time** — `Register` stores the pointer, so later
  `flag.Parse()` writes are visible through the registry. No ordering bug.
- **`DebugOverlay` lives in its own sub-package
  (`internal/engine/ui/debugoverlay/`)** to satisfy AC-5 ("new file ...") and
  avoid cluttering `ui/menu` with an unrelated widget. Tests inject a
  `keyJustPressed func(ebiten.Key) bool` seam so they don't depend on
  Ebitengine's real input loop.
- **Overlay-open short-circuits Scene and Dialogue updates** in `Game.Update`.
  Frame counter still increments — it's a global tick, not gameplay time.
  Slow-mo TPS work in story 071 is unaffected.
- **Existing `DebugPhysics` overlay is unplugged from F1.** F1 now owns the
  flag overlay. Re-binding `DebugPhysics` to another key (or migrating it
  into the new overlay) is out of scope.

## Risks & Quirks

- **Map-pointer rewrite of `channels`.** Switching from `map[string]bool` to
  `map[string]*bool` changes the behaviour of an aliased map only if external
  callers reach into internal state — none do (verified by grep on
  `debug.channels`). The Watch/Log fast paths remain allocation-free since
  the disabled fast path returns before any map lookup.
- **`Init` is the only public way to bind JSON-driven channel pointers.**
  Callers that bypass `Init` and only call `Register` for their own pointers
  still work; the registry merges JSON-loaded entries with explicit
  `Register` calls (last-write-wins per name).
- **Cursor clamping after `Reset` + re-register.** Edge case in the story:
  `Update()` clamps cursor to `len(entries)-1` (or 0 when empty) every frame
  so external mutations are tolerated.
- **No persistence** — overlay edits are in-memory only (story note). If a
  user toggles `collision_box` and restarts, the JSON value re-applies.

## Future

- [ ] Re-bind or absorb `DebugPhysics` into the new overlay.
- [ ] Numeric/enum debug values (currently bool-only).
- [ ] Persist overlay state back to `assets/data/debug.json` on demand.
- [ ] Mouse / touch input for the overlay.
- [ ] Scrollable list when entries exceed screen height.

## Playtest

**Standalone:** Yes — run `go run cmd/game/main.go`, press **F1** to open the
overlay. Use Up/Down to move the cursor, Space or Enter to toggle the
selected flag. Press F1 again to close. With `assets/data/debug.json`
containing `{"player_state":true, "physics":false}`, both keys appear in the
list along with `cam_debug` and `collision_box` from the CLI flags.
