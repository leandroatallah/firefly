# USER STORY — 071-slow-motion-debug-mode

**Branch:** `071-slow-motion-debug-mode`
**Bounded Context:** Game Logic (`internal/game/app/`, `internal/engine/data/config/`, `internal/engine/app/`)

---

## Story

As a developer,
I want to run the game at a configurable slow-motion speed via a CLI flag,
so that I can inspect combat timing, hitbox windows, physics trajectories, and animation frames frame-by-frame.

---

## Background

The game loop advances one tick per `Update()` call at Ebitengine's fixed TPS (60). There is no per-tick dt scaling today. Ebitengine exposes `ebiten.SetTPS(n)` to lower the logical tick rate without touching the renderer — using this as the slow-motion mechanism keeps physics and animation deterministic (they always advance by exactly one tick per `Update()`), while simply reducing how many ticks fire per real-world second. A factor of 0.25 lowers effective TPS from 60 to 15, making everything run at quarter speed without any floating-point accumulation or fixed-point breakage.

Integration uses both candidate points:
- `internal/game/app/config.go` — declares `-slow-mo` (bool) and `-slow-mo-factor` (float64) CLI flags, stores values in `AppConfig`. This is consistent with every other debug flag already in that file.
- `internal/engine/app/engine.go` (`Game.Update`) — reads `cfg.SlowMoFactor` once on first `Update()` after flag parse to call `ebiten.SetTPS`; when disabled the call is never made and default TPS is untouched.

No runtime keybinding is added: `engine.go` already has a pattern for F1 toggling debug overlays, but that toggle gates rendering, not game-loop timing. Changing TPS at runtime mid-session causes frame-count jumps that break replay and timing logic. CLI-flag-only scope is the correct and safe choice for this story.

---

## Acceptance Criteria

- AC-1: `AppConfig` in `internal/engine/data/config/config.go` gains two fields: `SlowMo bool` and `SlowMoFactor float64`.
- AC-2: `NewConfig()` in `internal/game/app/config.go` registers `-slow-mo` (bool, default `false`) and `-slow-mo-factor` (float64, default `0.25`) via `flag.BoolVar` / `flag.Float64Var` onto `cfg.SlowMo` and `cfg.SlowMoFactor`, following the existing flag registration pattern.
- AC-3: `Game.Update()` in `internal/engine/app/engine.go` calls `ebiten.SetTPS(int(math.Round(float64(ebiten.DefaultTPS) * cfg.SlowMoFactor)))` exactly once at startup when `cfg.SlowMo == true`; it uses a `slowMoApplied bool` field on `Game` to guard the one-time call.
- AC-4: When `cfg.SlowMo == false` (the default), `ebiten.SetTPS` is never called and the game runs at the default Ebitengine TPS with zero overhead in `Update()`.
- AC-5: `SlowMoFactor` values outside `(0, 1]` are clamped to `0.05` (min) and `1.0` (max) before the `SetTPS` call; a factor of exactly `1.0` is treated as a no-op (no `SetTPS` call, flag effectively disabled).
- AC-6: Running with `-slow-mo` and no `-slow-mo-factor` defaults to quarter-speed (factor `0.25`, effective TPS 15).
- AC-7: Running with `-slow-mo -slow-mo-factor=0.5` halves the TPS to 30; all game logic (physics, animations, hitbox windows) advances at exactly the same rate relative to ticks — only real-time wall-clock speed changes.
- AC-8: `config_test.go` (or a new file in `internal/engine/data/config/`) includes a table-driven test covering: `SlowMo=false` → no TPS override expected; `SlowMo=true, SlowMoFactor=0.25` → effective TPS = 15; `SlowMo=true, SlowMoFactor=1.0` → treated as no-op; `SlowMo=true, SlowMoFactor=0.0` → clamped to min, effective TPS = `int(60 * 0.05)` = 3; `SlowMo=true, SlowMoFactor=2.0` → clamped to 1.0, no-op.
- AC-9: Layer rules upheld — `internal/engine/data/config/` does not import `ebiten`; TPS clamping and `SetTPS` call remain in `internal/engine/app/engine.go`.
- AC-10: Existing tests in `internal/engine/app/app_test.go` and `internal/engine/data/config/config_test.go` continue to pass without modification.

---

## Behavioral Edge Cases

- `-slow-mo-factor=0` is not a valid speed; clamping to `0.05` prevents a TPS of zero which would freeze the game loop.
- `-slow-mo-factor` without `-slow-mo` has no effect; the field is stored but the `SetTPS` branch is never reached.
- The one-time `slowMoApplied` guard must fire before the first scene `Update()` so the very first tick already runs at reduced speed.
- `ebiten.DefaultTPS` is 60; the multiplication must use this constant, not a hardcoded literal, to remain correct if Ebitengine ever changes its default.
- Fullscreen, no-sound, and other debug flags are orthogonal; slow-mo composes with all of them without interaction.
