# User Story — 045-projectile-friendly-fire-prevention

## Title

Projectiles Must Not Collide With Other Projectiles

## As a...

Player engaged in combat

## I want...

Bullets I fire to pass through enemy bullets without being destroyed

## So that...

Combat feels physically accurate and gameplay is not unfairly disrupted by projectile-on-projectile interactions.

## Background

When a player's projectile (a bullet Body in the physics Space) travels toward an enemy, it may cross paths with an enemy projectile traveling in the opposite direction. Currently, the collision system treats any two projectile Bodies as valid collision pairs, causing both to be destroyed on contact. This is incorrect — projectiles are not Actors and should only interact with Actor Bodies (player or enemy Bodies), never with other projectile Bodies.

## Acceptance Criteria

- [ ] **Projectile vs. Projectile — no interaction**: When two projectile Bodies occupy the same region of the Space (regardless of owner — player or enemy), no collision event is generated between them, and neither projectile is destroyed.
- [ ] **Projectile vs. Actor — interaction preserved**: A projectile Body continues to register a valid collision against any Actor Body (player or enemy), triggering damage and projectile destruction as before.
- [ ] **No regression on melee hitboxes**: Melee attack hitboxes, which are distinct from projectile Bodies, are unaffected by this change and continue to interact with Actor Bodies correctly.
- [ ] **Deterministic behavior**: The filtering is applied at the collision detection or resolution layer — not as a side-effect of ordering or timing — so the result is the same regardless of frame rate or update order.
- [ ] **Collision behavior is opt-in per projectile type**: The collision filtering mechanism allows individual projectile Body types to declare whether they can be hit by other projectiles, so that future projectile types (e.g. an explosive rocket) can opt in to projectile-vs-projectile interaction without changing the default filtering logic.

## Out of Scope

- Friendly-fire between player projectiles and the player Actor (separate story if needed).
- Projectile-vs-environment (wall/terrain) collision behavior — unchanged.
- Implementing explosive projectiles or any other new projectile types — the extensibility hook is sufficient for this story.

## Edge Cases

- **Interceptable projectiles (future)**: Some projectile types — such as a rocket — may need to be destroyed when hit by another projectile (e.g. a player shoots down an incoming rocket). This is not the default behavior. The collision filtering design must make it straightforward to mark a specific projectile Body as interceptable, enabling projectile-vs-projectile collision for that type only, without affecting standard bullets. This can be modeled as a flag or collision category on the Body rather than a global rule change.

## Domain Notes

- A **projectile** is a Body managed within the physics Space, spawned by an Actor's attack sequence.
- An **Actor** is an entity with a state machine (player, enemy) whose Body is the valid target for projectile damage.
- Collision filtering should leverage the existing collision category or layer mechanism in `internal/engine/physics/` or `internal/engine/combat/`.
- The opt-in interceptability of a projectile Body should be expressible through the existing contracts in `internal/engine/contracts/` without requiring structural changes to the Space or Body types.
- Relevant packages to investigate: `internal/engine/physics/`, `internal/engine/combat/`, `internal/engine/contracts/`.
