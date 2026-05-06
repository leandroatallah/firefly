// Package speech provides the dialogue orchestrator for genre-reusable UI.
//
// The orchestrator (Manager) composes engine subsystems — audio, config,
// input — into a multi-line typing flow with typing-sound scheduling and
// spelling-skip behaviour. Speech primitives (Speech, SpeechBase,
// SpeechFont) live in internal/engine/ui/speech and are not duplicated here.
//
// Dependency rule (enforced by CI):
//   - kit/ui/speech MAY import internal/engine/...
//   - kit/ui/speech MUST NOT import internal/game/...
//   - It implements internal/engine/contracts/dialogue.Manager.
package speech
