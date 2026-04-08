# US-013 — Resolve golangci-lint Violations

**Branch:** `013-resolve-golangci-lint-violations`
**Bounded Context:** Cross-cutting (Engine + Game)

## Story

As a developer, I want the codebase to pass `golangci-lint run ./...` with zero issues, so that code quality gates in the pre-push hook and CI are meaningful and not ignored.

## Context

Running `golangci-lint run ./...` (with `.golangci.yml`) currently reports 130 issues across 6 linters. These must be resolved or explicitly suppressed with justification before the linter can serve as a reliable gate.

## Current Report Summary

| Linter | Count | Nature |
|---|---|---|
| `unused` | 50 | Dead code: unused fields, consts, funcs, test helpers |
| `gochecknoglobals` | 44 | Mix: intentional (state enums, registries) and accidental globals |
| `staticcheck` | 23 | Deprecated API calls, style suggestions (QF codes), empty branches |
| `ineffassign` | 7 | Assigned variables never read, mostly in physics |
| `gofmt` | 3 | Unformatted files in `contracts/` |
| `unparam` | 3 | Parameters always receiving the same value |

## Acceptance Criteria

- **AC1:** `golangci-lint run ./...` exits with code `0` (zero issues reported).
- **AC2:** Intentional globals (state enums, registries) are suppressed with `//nolint:gochecknoglobals` and a one-line comment explaining why.
- **AC3:** Dead code (`unused`) is either removed or, if kept intentionally, suppressed with justification.
- **AC4:** The 3 unformatted files in `internal/engine/contracts/` are formatted with `gofmt`.
- **AC5:** Deprecated API calls (`ReplacePixels`, `text` package) are replaced with their current equivalents.
- **AC6:** `ineffassign` violations in physics code are fixed (use `_` for intentionally discarded values).
- **AC7:** No linter is disabled in `.golangci.yml` solely to silence pre-existing issues without addressing them.

## Behavioral Edge Cases

- State enum globals (`Idle`, `Walking`, `Jumping`, etc.) are part of the engine's public API — they must not be removed, only annotated.
- Registry globals (`stateConstructors`, `stateEnums`) are intentional singletons — suppress, do not refactor.
- `gochecknoglobals` must remain enabled in `.golangci.yml`; suppression is per-site only.
- Test helpers flagged as `unused` should be removed only if confirmed unreferenced across all test files in the package.

## When Moving to Active

Before writing the spec, re-run the linter to get a fresh report (the `bullet.go` typecheck bug may be fixed by then):

```bash
golangci-lint run ./...
```

Use the updated output as the source of truth for the spec — do not rely on the report in this story.
