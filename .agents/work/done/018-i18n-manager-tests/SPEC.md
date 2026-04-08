# SPEC-018 — I18nManager Test Coverage

**Branch:** `018-i18n-manager-tests`  
**Bounded Context:** i18n (`internal/engine/data/i18n/`)  
**Date:** 2026-04-04

## Overview

Add comprehensive test coverage for `internal/engine/data/i18n` package. The package currently has zero test files. This spec defines the test scenarios that must pass to satisfy all acceptance criteria.

## Technical Requirements

### Existing Implementation

The `I18nManager` type in `internal/engine/data/i18n/i18n.go`:

- **Constructor:** `NewI18nManager(assets fs.FS) *I18nManager` — accepts a filesystem abstraction
- **Load method:** `Load(langCode string) error` — loads JSON translation files from `assets/lang/{langCode}.json`
- **Translate method:** `T(key string, args ...any) string` — retrieves translations with optional `fmt.Sprintf` formatting

### Test Infrastructure

- Use `testing/fstest.MapFS` for all filesystem mocking (no real filesystem access)
- Table-driven tests for all scenarios
- No GPU-dependent calls; no `ebiten.RunGame`
- Deterministic, non-flaky tests

### No New Contracts Required

The `I18nManager` is a standalone data service. No new interfaces need to be defined in `internal/engine/contracts/`.

## Pre-conditions

- `internal/engine/data/i18n/i18n.go` exists with `I18nManager` implementation
- No test file exists yet (`i18n_test.go` will be created)
- Test file will use only standard library and `testing/fstest`

## Post-conditions

- `internal/engine/data/i18n/i18n_test.go` created with full test coverage
- All acceptance criteria verified by passing tests
- Coverage for `internal/engine/data/i18n` reaches 100%
- All tests are deterministic and use `MapFS` for filesystem mocking

## Red Phase Scenario

The failing test suite will verify:

1. **Load Success (AC1):** `Load("en")` with valid JSON file populates `m.translations` correctly
2. **Load File Not Found (AC2):** `Load("xx")` returns error when file does not exist
3. **Load Invalid JSON (AC3):** `Load("invalid")` returns error when JSON is malformed
4. **T() Known Key (AC4):** `T("greeting")` returns translated value for existing key
5. **T() Missing Key Fallback (AC5):** `T("unknown_key")` returns the key itself
6. **T() With Formatting (AC6):** `T("welcome", "Alice")` applies `fmt.Sprintf` when args provided
7. **T() No Formatting (AC7):** `T("simple")` returns raw value when no args provided
8. **Coverage (AC8):** All code paths in `i18n.go` are exercised
9. **MapFS Only (AC9):** No real filesystem calls; all tests use `testing/fstest.MapFS`

## Integration Points

- **Bounded Context:** i18n
- **Package:** `internal/engine/data/i18n`
- **Dependencies:** Standard library only (`encoding/json`, `fmt`, `io`, `io/fs`, `testing`, `testing/fstest`)
- **No external contracts:** `I18nManager` is self-contained

## Test File Structure

`internal/engine/data/i18n/i18n_test.go` will contain:

- Table-driven test for `Load()` covering success, file-not-found, and invalid JSON cases
- Table-driven test for `T()` covering key lookup, fallback, formatting, and no-args cases
- Helper function to create `MapFS` with test data
- All tests use `*testing.T` parameter; no subtests required for this scope

## Notes

- The implementation is straightforward; tests are the primary deliverable
- No refactoring of `i18n.go` is required
- Tests must be minimal and focused on the acceptance criteria
