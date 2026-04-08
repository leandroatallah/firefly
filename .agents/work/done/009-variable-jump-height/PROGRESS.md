# PROGRESS-009 — Variable Jump Height

**Status:** ✅ Done
**Completed:** 2026-04-01
**Branch:** `009-variable-jump-height`

## Changes

| File | Change |
|---|---|
| `internal/engine/physics/skill/skill_platform_jump.go` | Added `jumpCutMultiplier`, `jumpCutPending`, `SetJumpCutMultiplier`, `applyJumpCut`; wired into `HandleInput`, `tryActivate`, `handleCoyoteAndJumpBuffering`, `Update` |
| `internal/engine/physics/skill/skill_platform_jump_test.go` | New file — 3 table-driven tests covering clamp, cut, no-cut, once-only |
| `internal/game/scenes/phases/player.go` | `jumpSkill.SetJumpCutMultiplier(0.4)` — activates the mechanic for the player |

## Quality Gates

- ✅ All 38 tests pass (`go test ./internal/engine/physics/skill/...`)
- ✅ New functions at 100% coverage (`SetJumpCutMultiplier`, `applyJumpCut`, `NewJumpSkill`)
- ✅ No global mutable state introduced
- ✅ Table-driven tests, no `ebiten.RunGame`, no `time.Sleep`
- ✅ No contract changes
- ✅ No new packages
