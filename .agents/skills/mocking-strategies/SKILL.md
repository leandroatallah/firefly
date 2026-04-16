---
name: mocking-strategies
description: Best practices for creating and placing mocks, including i18n and physics isolation.
---

# Mocking Strategies

## Placement Rules

- **Shared mocks** (used across multiple packages): place in `internal/engine/mocks/`
- **Package-local mocks** (used in one package only): define in `_test.go` or `mocks_test.go`

## Key Interfaces to Mock

- `BodiesSpace` — isolate `Actor` and `Item` updates from physics
- `SceneManager` — test `NavigateTo` / `NavigateBack` logic
- `fs.FS` — use `fstest.MapFS` for i18n unit tests without real files

## i18n Mock Example

```go
fsys := fstest.MapFS{
    "assets/lang/en.json": &fstest.MapFile{
        Data: []byte(`{"hello": "Hello"}`),
    },
}
```

## Style Rule

Do NOT use `_ = variable` to silence unused warnings. Use blank identifier in parameter lists:

```go
// Bad
func (t *T) Draw(img *ebiten.Image) { _ = img }

// Good
func (t *T) Draw(_ *ebiten.Image) {}
```
