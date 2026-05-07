// Package beatemup provides genre-reusable character traits for beat-em-up
// (2.5D) games built on the Firefly engine.
//
// Phase 1 (this package) is a scaffold: it registers the package and a
// placeholder type so the new altitude axis can be exercised by future
// stories (gravity, jump, ground detection) without touching engine layout.
//
// Dependency rule (enforced by CI):
//   - beatemup MAY import internal/engine/...
//   - beatemup MUST NOT import internal/game/...
package beatemup
