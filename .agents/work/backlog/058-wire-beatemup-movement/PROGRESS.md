# PROGRESS — 058-wire-beatemup-movement

**Status:** Backlog

## Pipeline State

- [x] Story Architect
- [ ] Spec Engineer
- [ ] Mock Generator
- [ ] TDD Specialist
- [ ] Feature Implementer
- [ ] Workflow Gatekeeper

## Log

- [Sonnet 4.6] [Story Architect] 2026-05-09 [FINISHED]: USER_STORY.md created. Story wires `EightDirectionalMovementSkill` (056) and `BeatEmUpMovementModel` (057) into `BeatEmUpCharacter` and the beat-em-up phase scene. Three open questions flagged for spec phase: bounds source, skill registration API, and model ownership.
- [Sonnet 4.6] [Story Architect] 2026-05-09 [UPDATED]: All three original open questions resolved in grilling session. Bounds source: no minY/maxY; obstacle tiles handle vertical blocking. Skill registration: extend `kitskills.FromConfig` with `cfg.Movement.Mode` discriminator ("horizontal" / "eight_dir"). Model ownership: `BeatEmUpCharacter`-owned at construction, following platformer precedent. Camera bounds wiring via `Camera().SetBounds(tilemapRect)` added to scene ACs. No open questions remain; ready for spec phase.
