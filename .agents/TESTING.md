# Testing Patterns

**Scope:** Patterns and examples. Non-negotiable rules live in `.agents/constitution.md §Tests`.

---

## 1. Table-Driven Tests

Prefer for any logic with multiple input/output scenarios.

```go
tests := []struct {
    name  string
    input int
    want  int
}{
    {"Case A", 1, 2},
    {"Case B", 2, 4},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := myFunc(tt.input)
        if got != tt.want {
            t.Errorf("got %d, want %d", got, tt.want)
        }
    })
}
```

## 2. Mocking

- Mock only at system boundaries using interfaces from `internal/engine/contracts/`.
- Shared mocks (used across multiple packages) → `internal/engine/mocks/`.
- Package-local mocks → `mocks_test.go` in the same package.
- Do not mock the database or internal packages directly — use contracts.
- Mock `BodiesSpace` to test `Actor` or `Item` updates in isolation.

## 3. Physics & Fixed-Point Arithmetic

- Always validate positions with `fp16.From16()` and `fp16.To16()` when checking `x16`/`y16` values.
- Scale factor is **16** (not 65536). See [ADR-007](../docs/adr/ADR-007-fp16-scale-factor.md).
  - 1 pixel = 16 units → use `<<4` for pixel-to-fp16
  - Velocity of 6 pixels/frame = 96 units/frame
- Collision edge cases to cover:
  - One pixel before collision
  - Partial overlap
  - Full overlap
  - Multiple collidables in one space
  - Fast movement (skipping over thin walls)

## 4. Scene Lifecycle

- Test `OnStart()`, `Update()`, `Draw()`, `OnFinish()` in correct order.
- Validate `NavigateTo` and `NavigateBack` using a mock `SceneManager`.

## 5. Headless Ebitengine

- For tests requiring `ebiten.Image`: use `ebiten.NewImage(w, h)` in a headless environment.
- Never use `ebiten.RunGame` or GPU-dependent calls in unit tests.
- Avoid tests dependent on human interaction or specific frame timings; use `timing` package mocks.
- No `time.Sleep` — use frame counters or virtual time.

## 6. Internationalization (i18n)

The `I18nManager` loads from `assets/lang/{langCode}.json`. Use `T(key, args...)`.

When testing i18n-dependent code:
- Create a mock `fs.FS` using `fstest.MapFS` for unit tests.
- Test missing keys (should return the key as fallback).
- Test formatting args: `T("key_with_%d", count)`.
- Test missing language files (should return error from `Load()`).

## 7. Determinism

- Tests must be deterministic and non-flaky.
- Never use real clocks or random seeds without explicit seeding.
- Physics and state machine tests: drive via frame counter, not wall time.
