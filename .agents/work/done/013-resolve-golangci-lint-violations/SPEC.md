# SPEC — 013 Resolve golangci-lint Violations

**Branch:** `013-resolve-golangci-lint-violations`
**Bounded Context:** Cross-cutting (`internal/engine/` + `internal/game/`)

---

## Fresh Linter Report (2026-04-04)

130 issues confirmed across 6 linters (same counts as story):

| Linter | Count |
|---|---|
| `unused` | 50 |
| `gochecknoglobals` | 44 |
| `staticcheck` | 23 |
| `ineffassign` | 7 |
| `gofmt` | 3 |
| `unparam` | 3 |

---

## Technical Requirements

### 1. `gofmt` — 3 files (AC4)

Format with `gofmt -w`:

- `internal/engine/contracts/body/body.go`
- `internal/engine/contracts/navigation/navigation.go`
- `internal/engine/contracts/vfx/vfx.go`

No logic changes; whitespace/import alignment only.

---

### 2. `ineffassign` — 7 violations (AC6)

All in `internal/engine/physics/movement/`:

**`movement_model_platform.go`** — 4 violations:
- Line 82: `vx16, vy16 := body.Velocity()` — `vx16` unused → `_, vy16 := body.Velocity()`
- Line 98: `vx16, vy16 := body.Velocity()` — `vy16` unused → `vx16, _ := body.Velocity()`
- Line 141: `vx16, vy16 := body.Velocity()` — `vx16` unused → `_, vy16 := body.Velocity()`
- Line 158: `vx16 = 0` — assigned but never read → remove the assignment

**`movement_models_test.go`** — 3 violations:
- Lines 289, 327, 334: `vx, vy = model.UpdateVerticalVelocity(...)` / `model.handleGravity(...)` — `vx` unused → `_, vy = ...`

---

### 3. `staticcheck` — 23 violations (AC5)

#### SA1019 — Deprecated APIs (must replace)

| File | Deprecated | Replacement |
|---|---|---|
| `internal/engine/app/engine.go:10` | `github.com/hajimehoshi/ebiten/v2/text` | `github.com/hajimehoshi/ebiten/v2/text/v2` |
| `internal/engine/render/vfx/vignette.go:132` | `.ReplacePixels(pixels)` | `.WritePixels(pixels)` |

#### S1008 — Simplify boolean return

- `internal/engine/physics/body/body_movable.go:195`: replace `if b.vy16 >= threshold { return true }; return false` with `return b.vy16 >= threshold`

#### ST1023 — Omit inferred type

- `internal/engine/physics/skill/skill_dash.go:69`: `var dirX int = 1` → `dirX := 1`

#### SA9003 — Empty branches (remove or add comment)

- `internal/engine/physics/movement/movement_funcs_test.go:88`: remove empty `if gotX == 0 && gotY == 0 {}` block
- `internal/engine/physics/skill/skill_test.go:79`: remove empty `if vx2 == 0 && vx != 0 {}` block

#### QF1001 — De Morgan's law (3 sites)

- `internal/engine/audio/loader.go:20`
- `internal/game/app/setup_audio.go:27`
- `internal/game/app/setup_audio.go:49`

Replace `!(A || B || C)` with `!A && !B && !C`.

#### QF1008 — Remove redundant embedded field selectors (14 sites)

Remove the intermediate embedded field name from selectors. Affected files:

- `internal/engine/entity/actors/character.go` (×2)
- `internal/engine/entity/actors/platformer/platformer.go`
- `internal/engine/entity/items/item_base.go` (×2)
- `internal/engine/physics/body/body_collidable.go` (×2)
- `internal/engine/render/vfx/text/floating_text.go`
- `internal/game/entity/actors/enemies/bat.go`
- `internal/game/entity/actors/enemies/wolf.go`
- `internal/game/entity/actors/player/climber.go`
- `internal/game/ui/speech/bubble.go` (×2)
- `internal/game/ui/speech/common.go`

---

### 4. `unparam` — 3 violations (AC3/AC6)

| File | Fix |
|---|---|
| `internal/engine/entity/actors/ducking_state_test.go:14` | Change `w` param to `_` or inline the constant |
| `internal/engine/physics/body/body_builder_test.go:25` | Change `state` param to `_` |
| `internal/game/scenes/phases/events.go:9` | Change `scene` param to `_` |

---

### 5. `unused` — 50 violations (AC3)

#### Remove dead production code

| File | Symbol |
|---|---|
| `internal/engine/entity/actors/actor_state_concrete.go:63` | field `count int` |
| `internal/engine/physics/movement/movement.go:13` | const `gravityForce` |
| `internal/engine/render/tilemap/tilemap_draw.go:111` | func `drawTileOpts` |
| `internal/engine/scene/scene_base.go:19` | field `space *space.Space` |
| `internal/engine/sequences/commands_sequence.go:17` | field `blockSequence *bool` |
| `internal/game/entity/items/item_power_base.go:62` | func `createPowerUpBase` |
| `internal/game/scenes/phases/scene.go:53` | field `hasEndpoints bool` |

#### Remove dead test code

| File | Symbol |
|---|---|
| `internal/engine/entity/actors/movement/mocks_test.go:190-191` | fields `moveLeftForce`, `moveRightForce` |
| `internal/engine/entity/actors/movement/mocks_test.go:271` | func `newMockActor` |
| `internal/engine/entity/actors/movement/state_chase_test.go:297` | type `mockActorWithSpace` + `Space()` method |
| `internal/engine/physics/body/body_test.go:197` | type `mockShape` + `Width()`, `Height()` methods |
| `internal/engine/physics/body/ownership_test.go:8` | interface `ownerGetter` |
| `internal/engine/render/camera/camera_test.go:13` | func `saveConfig` |
| `internal/engine/sequences/commands_vfx_test.go:81` | type `mockSceneManager` + `CurrentScene()` method |
| `internal/game/ui/hud/hud_test.go:17` | type `localMockActor` + all 32 methods |

---

### 6. `gochecknoglobals` — 44 violations (AC2)

No code changes. Add `//nolint:gochecknoglobals` with a one-line justification comment at each site.

#### Intentional state enum globals — annotate with: `// State enum: part of engine public API`

Files:
- `internal/engine/entity/actors/actor_state.go` (Idle, Walking, Jumping, Falling, Landing, Hurted, Dying, Dead, Exiting, IdleShooting, WalkingShooting, JumpingShooting, FallingShooting)
- `internal/engine/entity/actors/ducking_state.go` (Ducking)
- `internal/engine/entity/actors/movement/state_wander.go` (Wander)
- `internal/engine/entity/items/item_state.go` (Idle, Walking, Falling, Hurted)
- `internal/game/entity/actors/states/actor_state_concrete.go` (Dying, Dead, Exiting)
- `internal/game/entity/actors/states/dash_state.go` (StateDashing)
- `internal/game/entity/actors/states/grounded_state.go` (StateGrounded)
- `internal/game/entity/items/fall_platform.go` (Shaking, Break)
- `internal/game/scenes/phases/goal_type.go` (ReactEndpointType, SequenceGoalType, NoGoalType)

#### Intentional registry/singleton globals — annotate with: `// Singleton registry: intentional package-level state`

Files:
- `internal/engine/entity/actors/state_registry.go` (stateConstructors, stateEnums, nextEnumValue)
- `internal/engine/entity/actors/movement/registry.go` (movementStateConstructors, movementStateEnums, nextMovementEnumValue)
- `internal/engine/entity/items/state_registry.go` (stateConstructors, stateEnums, nextEnumValue)

#### Other intentional globals — annotate with inline justification

| File | Symbol | Justification comment |
|---|---|---|
| `internal/engine/audio/audio_test.go:12-13` | `audioManagerOnce`, `audioManager` | `// Test-level singleton: avoids re-initialising audio in each test` |
| `internal/engine/data/config/config.go:62` | `cfg` | `// Package-level config singleton: loaded once at startup` |
| `internal/engine/input/input.go:5` | `isKeyPressed` | `// Swappable function var: allows injection in tests` |
| `internal/engine/render/camera/camera.go:14` | `collisionBoxImage` | `// Lazily initialised debug image: allocated once per process` |
| `internal/engine/render/tilemap/tilemap_collisions.go:22` | `LayerNameMap` | `// Immutable lookup table: read-only after init` |

---

## Pre-conditions

- `.golangci.yml` is present and enables all 6 linters.
- No linter is disabled globally to silence these issues.
- `go build ./...` passes before starting.

## Post-conditions

- `golangci-lint run ./...` exits with code `0`.
- No linter is disabled in `.golangci.yml`.
- All `//nolint` directives include a justification comment on the preceding line.
- `go test ./...` still passes (no behaviour changes).

---

## Integration Points

- No new contracts needed; this story touches no interfaces.
- `internal/engine/contracts/` files are formatting-only changes.
- Removing `drawTileOpts` from `tilemap_draw.go` requires confirming it is unreferenced (grep confirms it is package-private and unused).
- Replacing `text` → `text/v2` in `engine.go` may require updating call-site API (font loading / draw calls); verify against Ebitengine v2.7 migration guide.

---

## Red Phase (Failing Test Scenario)

There are no new unit tests for this story — linter compliance is verified by the tool itself. The "failing test" is the linter exit code.

**Scenario:** `golangci-lint run ./...` currently exits `1` with 130 issues.

**Definition of Red:** CI / pre-push hook fails because `golangci-lint run ./...` returns non-zero.

**Definition of Green:** `golangci-lint run ./...` exits `0` with no issues reported.

The TDD Specialist's role here is to apply fixes file-by-file, re-running the linter after each group to confirm the issue count decreases monotonically to zero.

**Suggested fix order (lowest risk → highest risk):**
1. `gofmt` (mechanical, no logic)
2. `ineffassign` (mechanical blank-identifier substitutions)
3. `unparam` (mechanical blank-identifier substitutions)
4. `staticcheck` QF/S/ST codes (style, no behaviour change)
5. `unused` dead code removal (verify no external references first)
6. `staticcheck` SA1019 deprecated API replacements (API migration)
7. `gochecknoglobals` nolint annotations (annotation-only)
