# Agent Guidelines

**Scope:** Coverage priorities, code style, story management scripts. Testing patterns → `.agents/TESTING.md`. Pipeline stages → `.agents/WORKFLOW.md`. Standards → `.agents/constitution.md`.

## 🎯 Goal

Achieve **80%+ test coverage** across the codebase, prioritizing the engine's entity management and the game's level infrastructure. The codebase follows a three-layer architecture (`engine`, `kit`, `game`): game may import engine and kit; kit may import engine; engine must not import kit or game.

## 🔝 Priorities

1. **Entity State Machine (`internal/engine/entity/actors`)**: 63.6% coverage. The `handleState` logic — including the new `StateContributor` hook (ADR-008) — is the most critical and complex part of the engine.
2. **Level Management (`internal/game/scenes/phases`)**: 18.4% coverage. This is the foundation for all game levels; still the lowest coverage across the codebase.
3. **Player & Character Logic (`internal/game/entity/actors/player`)**: 60.5% coverage. `WireStateContributors` and the dash/shoot contributors need dedicated branch coverage.
4. **Sequences (`internal/engine/sequences`)**: 86.4% coverage. Essential for cutscenes and scripted events; per-command `block_sequence` paths now covered.
5. **Composite Grounded State (`internal/game/entity/actors/states`)**: Sub-state machine (`GroundedState`, `DuckingState`, `DashState`). Each sub-state transition must be independently tested.
6. **Physics Skills (`internal/engine/physics/skill`)**: 79.5% coverage. `JumpSkill` jump-cut, `DashSkill` tween deceleration, and `ShootingSkill` direction detection.
7. **Scene Freeze (`internal/engine/scene`)**: `FreezeController` tick/reset logic needs full branch coverage.
8. **Physics Tween (`internal/engine/physics/tween`)**: `InOutSineTween` interpolation values and `Done()` boundary.
9. **Combat (`internal/engine/combat/...`)**: `weapon` 96.2%, `projectile` 89.9%, `inventory` 51.5%. Faction gating (`applyDamage` / `resolveDamageable`) and `EnemyShooting` gate chain are both new surface area.

## 🛠 Testing Strategy & Patterns

See **[`.agents/TESTING.md`](.agents/TESTING.md)** for all patterns and examples (table-driven tests, mocking, physics/fp16, scene lifecycle, headless Ebitengine, i18n).

## 📋 Feature Implementation Workflow

See **[`.agents/WORKFLOW.md`](.agents/WORKFLOW.md)** for the complete Spec-Driven Development (SDD) pipeline: Story Architect → Spec Engineer → Mock Generator → TDD Specialist → Feature Implementer → Gatekeeper. The Spec Engineer produces two files: `SPEC.md` (agent-optimized, token-efficient) and `NOTES.md` (human context: risks, rationale, investigation findings).

## 📋 Standard Workflow for Coverage Tasks

1. **Check Dashboard**: Run `bash scripts/story.sh` to see what is currently in progress across all worktrees.
2. **Analyze Coverage**: Use the coverage tools to identify gaps:
...
3. **Progress**: Mark `[/]` in `PROGRESS.md` when starting; `[x]` when done.
4. **Identify Gaps**: Read the source file and identify functions or branches with 0% coverage.
4. **Create Test File**: If it doesn't exist, create `[filename]_test.go`.
5. **Write Tests**: Follow the patterns in this document and the referenced skills. Ensure you test both "happy paths" and error/edge cases.
6. **Verify**: Run the test and check the new coverage percentage.

## ⚠️ Precautions

- **AI Agent Skills**: Skills are located in `.kiro/skills/` and synced from `.agents/skills/`. Modify source files in `.agents/skills/` and run `make sync-skills` to update.
- **Never commit changes**. Do not use `git commit` or attempt to stage/commit files. The user is responsible for all version control operations.
- **No AI attribution trailers**. Do not add `Co-Authored-By` or similar trailers to commit messages.
- **Do not modify production code** unless you find a bug that makes it untestable (e.g., global state that needs to be injected).
- **Keep tests fast**. Avoid long `time.Sleep` calls; use virtual time or frame counters.
- **No Flaky Tests**. Ensure tests are deterministic.
- **No `_ = variable` Pattern**. Do not use `_ = variable` to silence unused variable warnings. Use blank identifier in parameter lists instead: `func (t *T) Method(_ Type) {}`

## Code Style: Avoid `_ = variable` Pattern

**Do NOT do this in production code:**

```go
func (t *Transition) Update() {
    _ = t.active  // Bad: clutters code
}
```

**Do this instead:**

```go
func (t *Transition) Update() {}  // Clean: just remove unused field reference
// or
func (t *Transition) Draw(_ *ebiten.Image) {}  // Use blank in param list
```

**Acceptable in tests:** Using `_ = funcCall()` to verify a function doesn't panic without checking return value.

## 🔍 Key Packages to Target

| Package | Current Coverage | Focus Area |
| :--- | :--- | :--- |
| `entity/actors` | 63.6% | `handleState` state machine, `StateContributor` hook (ADR-008), animation logic |
| `game/scenes/phases` | 18.4% | `PhasesScene` life cycle, goal tracking |
| `game/entity/actors/player` | 60.5% | `WireStateContributors`, dash/shoot contributors, input integration |
| `sequences` | 86.4% | Command execution, `block_sequence` flag, one-time / interruptible logic |
| `entity/items` | 51.7% | Item collection and state transitions |
| `scene` | 73.2% | Scene transitions and tilemap initialization |
| `data/i18n` | 100.0% | `I18nManager.Load()` and `T()` methods, error handling |
| `combat/weapon` | 96.2% | `ProjectileWeapon`, per-state spawn offsets, `EnemyShooting` gate chain |
| `combat/projectile` | 89.9% | Lifetime, faction-gated damage, impact / despawn VFX |
| `combat/inventory` | 51.5% | Switch/add/ammo tracking |
| `physics/skill` | 79.5% | `JumpSkill`, `DashSkill`, `ShootingSkill`, `FromConfig` factory |
| `kit/actors/platformer` | 35.7% | `PlatformerCharacter` state callbacks, movement blockers, coin tracking |
| `kit/states` | 100.0% | `IdleSubState[E,I]` transition branches; maintain as kit grows |

## 🛠 Story Management Scripts

For **what to work on next** and cross-story dependencies, read the active epic's ROADMAP — see [`.agents/work/ROADMAP.md`](.agents/work/ROADMAP.md) for the epic index. Story IDs do not imply execution order.

Use these scripts to manage story lifecycle. Do **not** move folders manually.

| Script | Purpose |
| :--- | :--- |
| `bash scripts/story.sh new <id-slug>` | Create a new story in `backlog/` |
| `bash scripts/story.sh start <id-slug>` | Move story `backlog/` → `active/` (Spec Engineer) |
| `bash scripts/story.sh done <id-slug>` | Move story `active/` → `done/` (Gatekeeper) |
| `bash scripts/story.sh status` | List all stories by lane |
| `go run scripts/next-id.go` | Print the next available story ID |
| `go run scripts/kanban.go` | Generate `kanban.html` board |
