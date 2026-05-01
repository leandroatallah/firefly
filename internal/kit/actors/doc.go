// Package kitactors provides genre-reusable character trait components
// for platformer games built on the Firefly engine.
//
// Traits are independently composable: a concrete game character can embed
// any combination (e.g., ShooterCharacter + MeleeCharacter for a brawler).
//
// Dependency rule (enforced by CI):
//   - kitactors MAY import internal/engine/...
//   - kitactors MUST NOT import internal/game/...
package kitactors
