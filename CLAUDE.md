@.agents/constitution.md

---

See `AGENTS.md` for the primary workflow and skills reference.

## Context Map

Where to find persistent context across sessions:

| Need | Source |
|---|---|
| Why a design decision was made | `docs/adr/` |
| Ubiquitous language + non-negotiable standards | `.agents/constitution.md` |
| Active / backlog / done stories | `.agents/work/` |
| SDD pipeline (Story → Spec → TDD → Code) | `.agents/WORKFLOW.md` |
| Testing patterns + coverage priorities | `AGENTS.md` |
| Engine internals | `internal/engine/README.md` |
| Combat system | `internal/engine/combat/README.md` |

## Claude Code Notes

- Sub-agents are defined in `.claude/agents/` (generated from `.agents/agents/` via `scripts/sync-agents.sh`).
- Skills are in `.claude/skills/` (generated from `.agents/skills/` via `scripts/sync-skills.sh`).
- Do not edit `.claude/agents/` or `.claude/skills/` directly — edit the sources in `.agents/` and re-run the sync scripts.
