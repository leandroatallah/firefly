# NOTES — 071-slow-motion-debug-mode

## Design Choices

- **Pure helper + thin wrapper.** `EffectiveTPS` is a pure function returning `(int, bool)`. The single side-effecting call to `ebiten.SetTPS` stays in `Game.Update`. This is the only way to satisfy AC-8 (table-driven test) and AC-9 (config pkg must not import ebiten) simultaneously — the helper lives in `internal/engine/app` (which already imports ebiten), and the `config` package stays pure data.
- **Guard placement before `FrameCount++`.** The story explicitly requires the first tick to already run at reduced speed (Behavioral Edge Cases). Calling `SetTPS` after `FrameCount++` would still be the first tick, but placing it at the top of `Update()` makes the precondition unambiguous and matches the "before first scene Update" requirement.
- **`slowMoApplied` always flips, even when disabled.** Cheaper than re-evaluating the branch every frame; one extra bool write on tick 1 vs. an unconditional `if cfg.SlowMo` check forever. Matches AC-4 ("zero overhead in Update").
- **No runtime toggle.** Story background explicitly rules out an F-key toggle; mid-session TPS change breaks replay/timing logic. Documented to prevent future drift.
- **Read-back of TPS not asserted in tests.** Ebitengine doesn't guarantee `ebiten.CurrentTPS()` reflects `SetTPS` outside an active `RunGame` loop. Helper-level math tests + guard flag assertion give full logical coverage without flake risk.

## Risks & Quirks

- `math.Round(60 * 0.05) == 3`. A factor below `1/60 ≈ 0.0167` would round to 0 and freeze the loop, but the `0.05` floor prevents this. Do not relax `SlowMoMinFactor` without re-checking the rounding floor.
- `ebiten.DefaultTPS` is currently 60; the helper takes it as a parameter so any future change in Ebitengine propagates automatically.
- The `flag` package's global registry means re-invoking `NewConfig()` in a test run would panic on duplicate flag registration. This is pre-existing behaviour; tests construct `*AppConfig` directly (see `TestGameUpdateAndDrawIntegration`) rather than calling `NewConfig()`.
- `SlowMo` composes with all other debug flags (fullscreen, no-sound, cam-debug) without interaction — they live in disjoint code paths.

## Future

- [ ] Per-scene slow-mo override (e.g., boss fight slow-motion replay).
- [ ] Slow-mo HUD overlay showing effective TPS.
- [ ] Time-scale field on `AppContext` for sub-tick scaling once dt-based physics lands.
- [ ] Runtime keybind once replay/timing systems can tolerate mid-session rate changes.
