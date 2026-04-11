# Caveman Review

One-line PR comments with direct, actionable feedback.

## Format

```
L{line}: {emoji} {issue}: {fix}
```

## Emojis

| Emoji | Meaning |
|-------|---------|
| 🔴 | Bug / critical error |
| 🟡 | Warning / potential issue |
| 🟢 | Suggestion / nice-to-have |
| 💡 | Architecture/design improvement |
| 🔒 | Security concern |
| ⚡ | Performance issue |

## Examples

```
L42: 🔴 bug: user null. Add guard.
L87: 🔒 leak: API key exposed. Use env var.
L15: 💡 complex logic. Extract function.
L103: ⚡ slow: O(n²) loop. Use map.
L64: 🟡 edge case: empty list. Handle.
```

## Rules

1. **One line per comment** - no paragraphs
2. **Direct** - skip pleasantries
3. **Actionable** - always include fix
4. **Specific** - reference exact line/issue
5. **No "nit"** - if worth mentioning, worth fixing

## Anti-patterns

```
L42: This looks good, but maybe consider adding a null check here? Just a thought!
L87: Nit: variable naming could be improved
```

```
L42: 🔴 bug: user null. Add guard.
L87: 🟢 naming: rename to `userData` for clarity.
```
