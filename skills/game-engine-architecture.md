# Game Engine Architecture

## Scene Lifecycle

Test that methods are called in the correct order:

```
OnStart() → Update() → Draw() → OnFinish()
```

Validate `NavigateTo` and `NavigateBack` using a mock `SceneManager`.

## Actor & State Machine

The `handleState` logic in `internal/engine/entity/actors` is the most critical area. Test all state transitions, including edge cases and invalid states.

## Headless Ebitengine

Use `ebiten.NewImage(w, h)` for tests requiring an `*ebiten.Image`. Avoid GPU-dependent or window-dependent code in unit tests.
