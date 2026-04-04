# 023 — Projectile Manager

**Branch:** `023-projectile-manager`
**Bounded Context:** Engine (`internal/engine/combat/projectile/`)

## Story

As a game developer using this boilerplate,
I want a centralized projectile manager that handles spawn, update, draw, and despawn,
so that bullet lifecycle is no longer owned by individual scenes.

## Context

`PhasesScene.SpawnBullet()` and `PhasesScene.bullets` make bullet management scene-specific and non-reusable. Moving this to a standalone manager in `AppContext` allows any scene or weapon to spawn projectiles without coupling.

## Acceptance Criteria

- **AC1** — `projectile.Manager` struct with `Spawn()`, `Update()`, `Draw()`, `Clear()` methods.
- **AC2** — `Spawn()` accepts a `ProjectileConfig` (speed, damage, behavior) plus initial position and velocity in fp16.
- **AC3** — `Update()` moves all active projectiles and removes those that are out-of-bounds or have collided.
- **AC4** — `Draw()` renders all active projectiles to the provided `*ebiten.Image`.
- **AC5** — `Manager` is added to `app.AppContext` and initialized in `setup.go`.
- **AC6** — `PhasesScene.SpawnBullet()` and `PhasesScene.bullets` are removed; scene delegates to `AppContext.ProjectileManager`.
- **AC7** — Unit tests cover: spawn increases count, out-of-bounds despawn, `Clear()` resets state.

## Notes

- Object pooling is a follow-up optimization; this story uses a plain slice.
- `ProjectileConfig` lives in `internal/engine/combat/projectile/`.
- Collision detection reuses existing `body.BodiesSpace` interface.
