# SPEC — 019 Godoc for Contracts Interfaces

## Scope

Documentation-only changes to `internal/engine/contracts/`. No logic, no signatures, no behaviour changes.

## Files in scope

| File | Interface(s) |
|---|---|
| `body/body.go` | `Shape`, `Movable`, `Collidable`, `Obstacle`, `Drawable`, `Touchable`, `Alive`, `Body`, `BodiesSpace`, `Ownable` |
| `navigation/navigation.go` | `Scene`, `SceneFactory`, `SceneManager`, `Transition` |
| `vfx/vfx.go` | `Manager` |
| `animation/animation.go` | `SpriteState`, `FacingDirectionEnum` constants |
| `context/context.go` | `ContextProvider` |
| `scene/freeze.go` | `Freezable` |

## Rules

1. Every exported interface type must have a godoc comment (starts immediately above `type … interface`).
2. Every method inside those interfaces must have a one-line `// MethodName …` godoc comment.
3. Exported constants and type aliases in scope must also carry a godoc comment.
4. No line of production logic may be added, removed, or altered.

## Acceptance tests (no code tests required — documentation only)

- `go build ./internal/engine/contracts/...` exits 0.
- `go vet ./internal/engine/contracts/...` exits 0.
- Manual review: every interface and method listed above has a godoc comment.

## Out of scope

- Files already documented (`one_way_platform.go`, `shooter.go`, `sequences.go`).
- Any package outside `internal/engine/contracts/`.
- The pre-existing build failure in `internal/game/scenes/phases`.
