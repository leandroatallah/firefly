// Package kitstates contains genre-reusable concrete Actor sub-state
// implementations.
//
// Dependency rule (enforced by CI):
//   - kitstates MAY import internal/engine/...
//   - kitstates MUST NOT import internal/game/...
//
// Types here are parameterised on the caller's enum and input contract
// to avoid coupling to a specific game's state-machine vocabulary.
package kitstates
