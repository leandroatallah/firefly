# Caveman Commit

Terse git commit messages for the caveman skill.

## Rules

1. **Subject line ≤50 characters**
2. **Focus on "why" over "what"** - code shows what, commit shows why
3. **Imperative mood** - "fix bug" not "fixed bug" or "fixes bug"
4. **No fluff** - skip explanations unless critical

## Examples

### Bad
```
Update the user authentication logic to handle edge cases where the token might be expired or invalid and also add better error messages
```

### Good
```
auth: handle expired tokens gracefully

Add 401 redirect on token expiry. Prevents crash on stale session.
```

## Pattern

```
{scope}: {action} {object}

{why it matters}. {what changed briefly}.
```

## Scopes

- `auth` - authentication/authorization
- `ui` - user interface
- `api` - API endpoints
- `db` - database/schema
- `test` - test code
- `ci` - CI/CD
- `docs` - documentation
- `perf` - performance
- `fix` - bug fixes (no scope)

## Actions

- `add` - new feature/file
- `fix` - bug fix
- `rm` - remove code/file
- `refactor` - restructure
- `update` - modify existing
- `move` - relocate file/code
