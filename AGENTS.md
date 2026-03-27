# Agent Guidelines: Engine Test Coverage

This document provides specialized instructions for AI agents tasked with increasing test coverage for the `internal/engine` module.

## 🎯 Goal

Achieve **80%+ test coverage** across the codebase, prioritizing the engine's entity management and the game's level infrastructure.

## 🔝 Priorities

1. **Entity State Machine (`internal/engine/entity/actors`)**: 48.0% coverage. The `handleState` logic is the most critical and complex part of the engine and is currently under-tested.
2. **Level Management (`internal/game/scenes/phases`)**: 0.0% coverage. This is the foundation for all game levels.
3. **Player & Character Logic (`internal/game/entity/actors/player`)**: 52.2% coverage. Good progress here, but still needs more integration tests.
4. **Sequences (`internal/engine/sequences`)**: 60.8% coverage. Essential for cutscenes and scripted events.

## 🛠 Testing Strategy & Patterns

### 1. Table-Driven Tests

Prefer table-driven tests for logic with multiple input/output scenarios (e.g., movement, collisions, math).

```go
tests := []struct {
    name    string
    input   int
    want    int
}{
    {"Case A", 1, 2},
    {"Case B", 2, 4},
}
```

### 2. Mocking & Contracts

- Use the interfaces in `internal/engine/contracts/` to create mock implementations for testing.
- **Reusable Mocks**: If a mock is used in more than one test file across different packages, place it in `internal/engine/mocks`. This prevents code duplication.
- **Package-Specific Mocks**: If a mock is only relevant to a single package, define it within the `_test.go` file of that package or a `mocks_test.go` file in the same directory.
- Avoid using actual Ebitengine windows or GPU-dependent code in unit tests.
- Mock `BodiesSpace` to test `Actor` or `Item` updates in isolation.

### 3. Physics & Fixed-Point Arithmetic

- Always validate positions using `fp16.From16()` and `fp16.To16()` when checking `x16` and `y16` values.
- Test edge cases for collisions:
  - One pixel before collision.
  - Partial overlap.
  - Full overlap.
  - Multiple collidables in one space.
  - Fast movement (skipping over thin walls).

### 4. Scene Lifecycle

- Test that `OnStart()`, `Update()`, `Draw()`, and `OnFinish()` are called in the correct order.
- Validate `NavigateTo` and `NavigateBack` logic using a mock `SceneManager`.

### 5. Headless Ebitengine

- For tests that require an `ebiten.Image`, use `ebiten.NewImage(w, h)` in a headless environment.
- Avoid tests that depend on human interaction or specific frame timings (use `timing` package mocks).

### 6. Internationalization (i18n)

- The `I18nManager` loads translations from `assets/lang/{langCode}.json` files.
- Use `T(key, args...)` to retrieve translated strings with optional `fmt.Sprintf`-style formatting.
- When testing i18n-dependent code:
  - Create a mock `fs.FS` using `embed.FS` or `fstest.MapFS` for unit tests.
  - Test missing keys (should return the key itself as fallback).
  - Test formatting arguments: `T("key_with_%d", count)`.
  - Test missing language files (should return an error from `Load()`).

## 📋 Standard Workflow for Agents

1. **Analyze Coverage**: Use the coverage tools to identify gaps:
   - For a specific package: `go test ./internal/engine/[package] -coverprofile=coverage.out && go tool cover -func=coverage.out`
   - For a full project report: `bash scripts/test_coverage.sh`
2. **Consult Skills**: Review the relevant AI Agent Skills in `.agents/skills/` (e.g., `go-testing`, `mocking-strategies`, `fixed-point-physics`) for specialized patterns and requirements.
3. **Identify Gaps**: Read the source file and identify functions or branches with 0% coverage.
4. **Create Test File**: If it doesn't exist, create `[filename]_test.go`.
5. **Write Tests**: Follow the patterns in this document and the referenced skills. Ensure you test both "happy paths" and error/edge cases.
6. **Verify**: Run the test and check the new coverage percentage.

## ⚠️ Precautions

- **AI Agent Skills**: If you need to update agent instructions, modify the files in `.agents/skills/` and run `make sync-skills`.
- **Never commit changes**. Do not use `git commit` or attempt to stage/commit files. The user is responsible for all version control operations.
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
| `entity/actors` | 48.0% | `handleState` state machine, animation logic |
| `game/scenes/phases` | 0.0% | `PhasesScene` life cycle, goal tracking |
| `game/entity/actors/player` | 52.2% | Player input, physics integration, interactions |
| `sequences` | 60.8% | Command execution and completion conditions |
| `entity/items` | 52.3% | Item collection and state transitions |
| `scene` | 73.0% | Scene transitions and tilemap initialization |
| `game/physics/skill` | 78.0% | Game-specific power-ups (freeze, grow, star) |
| `data/i18n` | (new) | `I18nManager.Load()` and `T()` methods, error handling |
