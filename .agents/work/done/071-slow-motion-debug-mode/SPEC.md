# SPEC — 071-slow-motion-debug-mode

**Branch:** `071-slow-motion-debug-mode`
**Bounded Context:** Game Logic / Engine App / Config

---

## 1. Config Layer [AC-1, AC-9]

**File:** `internal/engine/data/config/config.go`

Add two fields to `AppConfig` (place them adjacent to other debug flags like `CamDebug`, `CollisionBox`):

```go
SlowMo       bool    // when true, lower effective TPS by SlowMoFactor
SlowMoFactor float64 // multiplier on ebiten.DefaultTPS, valid range (0,1]; clamped at apply time
```

**Layer constraint:** package `config` MUST NOT import `github.com/hajimehoshi/ebiten/v2`. No helper, constant, or test in this package may reference `ebiten.DefaultTPS` or `ebiten.SetTPS`.

---

## 2. Flag Registration [AC-2, AC-6]

**File:** `internal/game/app/config.go`, function `NewConfig()`.

Append after the existing `flag.IntVar(&cfg.TypingSoundCooldownFrames, ...)` line:

```go
flag.BoolVar(&cfg.SlowMo, "slow-mo", false, "Enable slow-motion debug mode (lowers effective TPS)")
flag.Float64Var(&cfg.SlowMoFactor, "slow-mo-factor", 0.25, "Slow-motion TPS multiplier (clamped to [0.05, 1.0])")
```

No struct-literal initialization of `SlowMo` / `SlowMoFactor` in the `&config.AppConfig{...}` block — `flag.*Var` sets defaults.

---

## 3. Engine Layer — Pure Helper [AC-5, AC-7, AC-8]

**File:** `internal/engine/app/engine.go` (or new sibling `slowmo.go` in same package — implementer choice).

Constants and helper, exported for testability:

```go
const (
    SlowMoMinFactor = 0.05
    SlowMoMaxFactor = 1.0
)

// EffectiveTPS computes the target TPS for slow-motion mode.
// Returns (targetTPS, shouldApply).
// shouldApply == false means SetTPS must NOT be called.
//   - slowMo==false                  → (0, false)
//   - clampedFactor >= 1.0 (no-op)   → (0, false)
//   - otherwise                       → (round(defaultTPS * clampedFactor), true)
// Factor is clamped to [SlowMoMinFactor, SlowMoMaxFactor] before evaluation.
func EffectiveTPS(slowMo bool, factor float64, defaultTPS int) (int, bool)
```

Pseudocode:

```
EffectiveTPS(slowMo, factor, defaultTPS):
  if !slowMo: return 0, false
  if factor < SlowMoMinFactor: factor = SlowMoMinFactor
  if factor > SlowMoMaxFactor: factor = SlowMoMaxFactor
  if factor == SlowMoMaxFactor: return 0, false        // 1.0 → no-op
  tps = int(math.Round(float64(defaultTPS) * factor))
  return tps, true
```

Import: add `"math"` to `internal/engine/app/engine.go` (or to `slowmo.go`).

---

## 4. Game.Update Wiring [AC-3, AC-4]

**File:** `internal/engine/app/engine.go`, struct `Game`.

Add field:

```go
slowMoApplied bool
```

In `Game.Update()`, **before** the existing `g.AppContext.FrameCount++` line, insert:

```go
if !g.slowMoApplied {
    g.slowMoApplied = true
    cfg := g.AppContext.Config
    if tps, ok := EffectiveTPS(cfg.SlowMo, cfg.SlowMoFactor, ebiten.DefaultTPS); ok {
        ebiten.SetTPS(tps)
    }
}
```

Post-conditions:
- `slowMoApplied` flips to `true` on the very first `Update()` regardless of `SlowMo` value.
- `ebiten.SetTPS` is called at most once per `Game` instance.
- When `cfg.SlowMo == false`, branch executes once and never calls `SetTPS`.

---

## 5. Test Plan (Red Phase)

### 5.1 Helper tests — `internal/engine/app/slowmo_test.go` [AC-5, AC-7, AC-8]

Package: `package app` (already imports `ebiten`; satisfies AC-9 because tests are NOT in `config` pkg).

Table-driven `TestEffectiveTPS`:

```
T-S1: slow-mo disabled
  pre:  slowMo=false, factor=0.25, defaultTPS=60
  post: tps=0, ok=false

T-S2: quarter speed default
  pre:  slowMo=true, factor=0.25, defaultTPS=60
  post: tps=15, ok=true

T-S3: half speed
  pre:  slowMo=true, factor=0.5, defaultTPS=60
  post: tps=30, ok=true

T-S4: factor 1.0 is no-op
  pre:  slowMo=true, factor=1.0, defaultTPS=60
  post: tps=0, ok=false

T-S5: factor 0.0 clamped to min
  pre:  slowMo=true, factor=0.0, defaultTPS=60
  post: tps=3, ok=true            // round(60 * 0.05) = 3

T-S6: factor 2.0 clamped to max → no-op
  pre:  slowMo=true, factor=2.0, defaultTPS=60
  post: tps=0, ok=false

T-S7: negative factor clamped to min
  pre:  slowMo=true, factor=-0.5, defaultTPS=60
  post: tps=3, ok=true

T-S8: rounding (factor=0.333…, defaultTPS=60)
  pre:  slowMo=true, factor=1.0/3.0, defaultTPS=60
  post: tps=20, ok=true
```

Use `ebiten.DefaultTPS` (60) as the `defaultTPS` arg in at least one case to lock the constant.

### 5.2 Game.Update guard test — extend `internal/engine/app/app_test.go` [AC-3, AC-4]

New test `TestGameUpdateSlowMoAppliedGuard`:

```
T-G1: slowMoApplied flips on first Update regardless of SlowMo
  pre:  cfg.SlowMo=false; game.slowMoApplied==false
  act:  game.Update(); game.Update()
  post: game.slowMoApplied==true (no panic, no error)

T-G2: SceneManager.Update still called when slow-mo path runs
  pre:  cfg.SlowMo=true, cfg.SlowMoFactor=0.25
  act:  game.Update()
  post: sm.UpdateCalled==true, FrameCount==1, game.slowMoApplied==true
```

Do NOT assert real `ebiten.CurrentTPS()` — Ebitengine does not guarantee read-back of `SetTPS` outside `RunGame`. Coverage of the SetTPS branch is satisfied by the helper tests plus the guard flag assertion. Document this in the test file comment.

### 5.3 Config preservation [AC-10]

- `TestSetAndGet` and `TestSetNil` in `internal/engine/data/config/config_test.go` MUST continue to pass with the new fields added (they don't reference the new fields).
- `TestGameUpdateAndDrawIntegration` MUST continue to pass; the new guard branch must be a no-op when `cfg.SlowMo==false`.

---

## 6. Resolving AC-8 / AC-9 Tension

AC-8 names `config_test.go` as the test location; AC-9 forbids `ebiten` import in the `config` package. The pure `EffectiveTPS` helper is the seam: it lives in `internal/engine/app/` and takes `defaultTPS int` as a parameter, so the `config` package never references `ebiten`. The table-driven test required by AC-8 lives in `internal/engine/app/slowmo_test.go` and covers every case enumerated in AC-8. This satisfies the **intent** of AC-8 (table-driven test of the clamp/no-op/TPS math) while strictly satisfying AC-9.

---

## 7. File Inventory

| Path | Change |
|---|---|
| `internal/engine/data/config/config.go` | +2 fields on `AppConfig` |
| `internal/game/app/config.go` | +2 `flag.*Var` calls in `NewConfig()` |
| `internal/engine/app/engine.go` | +`slowMoApplied bool` field; +`EffectiveTPS` helper + constants (or in new `slowmo.go`); +guard block in `Update()`; +`math` import |
| `internal/engine/app/slowmo_test.go` | NEW — table-driven `TestEffectiveTPS` |
| `internal/engine/app/app_test.go` | +`TestGameUpdateSlowMoAppliedGuard` |

## 8. Contracts

No new contracts. No mocks required. `MockSceneManager` and `MockDialogueManager` already cover `Game.Update` dependencies.

## 9. Out of Scope

See `NOTES.md`.
