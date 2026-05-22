# Architecture Map

**Dependency rule:** `game` → `kit` → `engine`. No reverse imports.

---

## Contracts (`internal/engine/contracts/`)

Interfaces only — no implementations. All cross-layer dependencies go through contracts.

| Package | Key Interfaces | Implemented By |
|---|---|---|
| `contracts/body` | `Movable`, `Collidable`, `MovableCollidable`, `BodiesSpace` | `physics/body`, `physics/space` |
| `contracts/scene` | `SceneManager`, `Navigable`, `Freezable` | `engine/scene` |
| `contracts/animation` | `Animator` | `engine/render/sprites` |
| `contracts/context` | `Context` | `engine/app` |
| `contracts/vfx` | `VFX` spawners | `engine/render/vfx` |
| `contracts/sequences` | `SequencePlayer` | `engine/sequences` |

---

## Engine Layer (`internal/engine/`)

Core, game-agnostic systems. Must not import `kit` or `game`.

| Package | Role |
|---|---|
| `entity/actors` | `Character` — base actor with state machine, `StateContributor` hook (ADR-008) |
| `entity/items` | `Item` — base collectible/interactive entity |
| `physics/body` | `Body` — owns position (x16/y16), shape, collision callbacks |
| `physics/space` | `Space` — holds bodies, resolves collisions |
| `physics/movement` | `MovementModel` — applies velocity + calls Space each tick (Platformer / TopDown) |
| `physics/skill` | `JumpSkill`, `DashSkill`, `HorizontalMovementSkill`, `ShootingSkill` |
| `scene` | `SceneManager`, `FreezeController`, scene lifecycle |
| `sequences` | `SequencePlayer`, scriptable commands (actor, camera, music, VFX) |
| `input` | `HorizontalAxis` — last-pressed-wins directional input |
| `data/i18n` | `I18nManager` — loads `assets/lang/{langCode}.json`, provides `T(key, args...)` |
| `render` | Camera, particles, sprites, tilemap, VFX |
| `mocks` | Shared test mocks for contracts |

---

## Kit Layer (`internal/kit/`)

Genre-reusable concrete implementations. Imports `engine`, must not import `game`.

| Package | Role |
|---|---|
| `kit/skills` | Concrete skill implementations wrapping `engine/physics/skill` |
| `kit/combat/weapon` | `ProjectileWeapon`, `EnemyShooting` |
| `kit/combat/projectile` | Projectile lifecycle, faction-gated damage |
| `kit/combat/inventory` | Weapon switch/add/ammo tracking |
| `kit/combat/melee` | `melee.State`, melee controller, combo steps |
| `kit/actors/platformer` | `PlatformerCharacter` — trait composition for platformer actors |
| `kit/states` | `IdleSubState[E,I]` and other genre-reusable states |
| `kit/ui/speech` | `speech.Manager` — dialogue implementation |

---

## Game Layer (`internal/game/`)

Project-specific code. Wires kit + engine to art, levels, and rules.

| Package | Role |
|---|---|
| `game/entity/actors/player` | `ClimberPlayer`, `WireStateContributors` (dash + shooting contributors) |
| `game/entity/actors/enemies` | Game-specific enemy actors |
| `game/entity/actors/states` | `GroundedState` sub-state machine (idle, walk, duck, aim-lock) |
| `game/scenes/phases` | `PhasesScene` — level lifecycle, goal tracking |
| `game/app/setup` | Wires all systems, registers scenes |

---

## Dataflow (one frame)

```
[Input]
  → HorizontalAxis / button state
  → Actor.Update()
      → StateContributors polled (skill-driven states)
      → handleState() → state transition
      → MovementModel.Update()
          → Body velocity applied
          → Space.ResolveCollisions()
              → OnTouch / OnBlock callbacks
  → Sequences.Update() (scripted events)
  → Scene.Draw()
      → Camera transform
      → Sprites / Tilemap / VFX / Particles
```
