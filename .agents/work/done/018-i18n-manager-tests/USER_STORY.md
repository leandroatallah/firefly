# US-018 — I18nManager Test Coverage

**Branch:** `018-i18n-manager-tests`
**Bounded Context:** i18n

## Story

As a developer, I want `internal/engine/data/i18n` to have full test coverage, so that translation loading and lookup regressions are caught automatically.

## Context

`internal/engine/data/i18n` has **no test files** (0% coverage). The package has a single file `i18n.go` with two public methods: `Load(langCode string) error` and `T(key string, args ...any) string`. It accepts an `fs.FS` at construction, making it straightforward to test with `testing/fstest.MapFS` — no real filesystem or assets needed.

## Acceptance Criteria

- **AC1:** `Load()` succeeds with a valid JSON language file and populates translations.
- **AC2:** `Load()` returns an error when the language file does not exist.
- **AC3:** `Load()` returns an error when the file contains invalid JSON.
- **AC4:** `T()` returns the translated string for a known key.
- **AC5:** `T()` returns the key itself as fallback when the key is missing.
- **AC6:** `T()` applies `fmt.Sprintf` formatting when `args` are provided.
- **AC7:** `T()` returns the raw value (no formatting) when no `args` are provided.
- **AC8:** Coverage for `internal/engine/data/i18n` reaches **100%**.
- **AC9:** Tests use `testing/fstest.MapFS` — no real filesystem access.
