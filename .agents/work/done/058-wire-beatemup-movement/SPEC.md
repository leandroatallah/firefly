# SPEC — 058-wire-beatemup-movement

Wires `EightDirectionalMovementSkill` (056) and `BeatEmUpMovementModel` (057) into `BeatEmUpCharacter` and the beat-em-up phase scene. Adds a `Mode` discriminator to `MovementConfig` and extends `kitskills.FromConfig`.

---

## 1. Schema Change — `MovementConfig` [AC-1]

File: `internal/engine/data/schemas/json.go`

```go
type MovementConfig struct {
    Enabled         *bool   `json:"enabled,omitempty"`
    HorizontalSpeed float64 `json:"horizontal_speed,omitempty"`
    Mode            string  `json:"mode,omitempty"` // "" | "horizontal" | "eight_dir"
}
```

Constants (new file `internal/engine/data/schemas/movement_mode.go`):

```go
const (
    MovementModeHorizontal = "horizontal"
    MovementModeEightDir   = "eight_dir"
)
```

Rules:
- `Mode == ""` → treat as `"horizontal"` (backward-compatible).
- Unknown value → log a warning, fall back to `"horizontal"`.

---

## 2. Skill Factory — `kitskills.FromConfig` [AC-1]

File: `internal/kit/skills/factory.go`

Replace the Movement block:

```
if cfg.Movement != nil && isEnabled(cfg.Movement.Enabled):
    switch cfg.Movement.Mode:
      case "", "horizontal": skills << NewHorizontalMovementSkill()
      case "eight_dir":      skills << NewEightDirectionalMovementSkill()
      default: log warn "unknown movement mode: <m>; falling back to horizontal"
               skills << NewHorizontalMovementSkill()
```

Pre/Post:
- pre: `cfg != nil && cfg.Movement != nil && isEnabled == true`
- post: exactly one movement skill appended; type matches `Mode`.

---

## 3. BeatEmUpCharacter Construction [AC-2, AC-3]

File: `internal/kit/actors/beatemup/beatemup_character.go`

Update signature to mirror platformer precedent (sprite + body + model owned at construction). New layout:

```go
type BeatEmUpCharacter struct {
    *actors.Character
    *kitactors.MeleeCharacter
    app.AppContextHolder
}

func NewBeatEmUpCharacter(
    fsys fs.FS,
    stateMap map[string]animation.SpriteState,
    spriteData schemas.SpriteData,
    bodyRect *bodyphysics.Rect,
    blocker physicsmovement.PlayerMovementBlocker,
) (*BeatEmUpCharacter, error)
```

Body:
```
sprites := sprites.GetSpritesFromAssets(fsys, spriteData.Assets, stateMap)
c := actors.NewCharacter(sprites, bodyRect)
be := &BeatEmUpCharacter{
    Character:      c,
    MeleeCharacter: kitactors.NewMeleeCharacter(),
}
c.SetMovementModel(physicsmovement.NewBeatEmUpMovementModel(blocker))
c.SetFaceDirection(spriteData.FacingDirection)
c.SetFrameRate(spriteData.FrameRate)
c.SetOwner(be)
return be, nil
```

Constraints:
- `physicsmovement.PlatformMovementModel` MUST NOT be referenced in this package.
- `MovementModel()` (inherited from `Character`) returns the `*BeatEmUpMovementModel`.
- No game-layer import (AC-8).

### 3a. Backward-Compatible Constructor

Keep zero-arg `NewBeatEmUpCharacter()` available only if existing call sites require it; otherwise migrate them. Audit: `internal/game/entity/actors/player/cody*.go` (CodyPlayer) — update to new signature.

---

## 4. Skill Registration Path [AC-3]

Existing path (no new code): `createPlayer` in `internal/game/scenes/phases/beatemup/player.go` already calls `kitskills.FromConfig(spriteData.Skills, deps)` and `builder.ApplySkills(p, skills)`. With Section 2, when `cody.json` sets `movement.mode == "eight_dir"`, the registered skill is `EightDirectionalMovementSkill`.

Data change (game assets):
- `assets/data/actors/cody.json` (or equivalent CodyPlayer config) — add `"mode": "eight_dir"` under `skills.movement`.
- File path verified at implementation time; no schema-side breakage.

Post-conditions per frame:
- `HandleInput` of `EightDirectionalMovementSkill` called by `Character.Update` skill loop.

---

## 5. Beat-em-up Scene — Camera + Collision Wiring [AC-4, AC-5]

File: `internal/game/scenes/phases/beatemup/scene.go` — `OnStart`.

After `s.initTilemap()` and before `s.bodyCounter.setBodyCounter(...)` (current ordering preserved), ensure exactly:

```
tilemapRect := image.Rect(0, 0, s.GetTilemapWidth(), s.GetTilemapHeight())
s.Camera().SetBounds(&tilemapRect)
if s.hasPlayer:
    s.Camera().SetFollowTarget(s.player)
s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), endpointFactory)
```

Where `endpointFactory` is the existing `func(id string) body.Touchable { ... NewTouchTrigger ... }` block (unchanged).

Note: current scene passes `nil` as the second arg to `CreateCollisionBodies` according to the AC. Reconcile: AC-5 says `nil`, but the existing scene already passes an endpoint factory; the factory is required for SPIKE/CUTSCENE triggers. Keep the existing factory. Treat the AC-5 `nil` wording as "no scene-injected bounds args"; collision bodies themselves continue to use the endpoint factory.

Pre/Post:
- pre: tilemap loaded; `GetTilemapWidth/Height > 0`.
- post: `s.Camera().Bounds() != nil`; `Bounds() == tilemapRect`; `FollowTarget == s.player` when `hasPlayer`.
- post: physics space contains obstacle bodies for every solid tile.

Edge cases:
- No collision tiles → `CreateCollisionBodies` returns; player moves freely; no panic.
- Camera at tilemap edge → `camera.Controller.SetBounds` clamps via existing logic (already tested in `camera_test.go`).

---

## 6. Runtime Behavior [AC-6]

Per-frame flow (post-implementation):
```
Character.Update(space):
  for each skill: skill.HandleInput(body, model, space)
    EightDirectionalMovementSkill:
      if model.IsInputBlocked() or body.Immobile(): zero v/a, return
      else: read input.CommandsReader(); OnMove{Left,Right,Up,Down}(speed)
  model.Update(body, space):
    BeatEmUpMovementModel:
      ApplyValidPosition(vx16, true, space)   // obstacle blocks X
      ApplyValidPosition(vy16, false, space)  // obstacle blocks Y
      integrate acceleration, clamp 2D speed, friction both axes
      no gravity write to vy16
```

Verifiable:
- Idle frame: `vy16 == 0` after Update (no gravity accumulation).
- Pressing Down into obstacle: Y position unchanged after Update.
- Immobile flag set: `vx16 == vy16 == 0` after Update.

---

## 7. Layer & Regression Rules [AC-7, AC-8]

- `internal/kit/actors/beatemup/` imports: `engine/...` + `kit/actors` only. No `internal/game/...`.
- Platformer phase scene unchanged; `PlatformerCharacter` unchanged. Run `go test ./internal/kit/actors/platformer/... ./internal/game/scenes/phases/platformer/...` post-impl.

---

## 8. Tests — Red Phase [AC-9]

### 8.1 `internal/kit/skills/factory_test.go`

```
T-F1: FromConfig mode "eight_dir" returns EightDirectionalMovementSkill
  pre:  cfg.Movement = {Enabled:true, Mode:"eight_dir"}
  act:  skills := FromConfig(cfg, deps)
  post: len(skills)==1 && skills[0].(*EightDirectionalMovementSkill) != nil

T-F2: FromConfig mode "horizontal" returns HorizontalMovementSkill
  pre:  cfg.Movement = {Enabled:true, Mode:"horizontal"}
  post: skills[0] type == *HorizontalMovementSkill

T-F3: FromConfig empty mode defaults to horizontal
  pre:  cfg.Movement = {Enabled:true, Mode:""}
  post: skills[0] type == *HorizontalMovementSkill

T-F4: FromConfig unknown mode falls back to horizontal (no panic)
  pre:  cfg.Movement = {Enabled:true, Mode:"jetpack"}
  post: skills[0] type == *HorizontalMovementSkill
```

### 8.2 `internal/kit/actors/beatemup/beatemup_character_test.go`

```
T-B1: NewBeatEmUpCharacter returns non-nil with MeleeCharacter and Character set
  act:  c, err := NewBeatEmUpCharacter(fakeFS, stateMap, spriteData, bodyRect, nil)
  post: err==nil; c.Character != nil; c.MeleeCharacter != nil

T-B2: BeatEmUpCharacter owns BeatEmUpMovementModel
  post: _, ok := c.MovementModel().(*physicsmovement.BeatEmUpMovementModel); ok == true

T-B3: BeatEmUpCharacter does not panic on zero-input update frame
  pre:  empty input.CommandsReader (no keys); space with no obstacles
  act:  c.Update(space)
  post: no panic; err==nil; v == (0,0)

T-B4: AddSkill(EightDirectionalMovementSkill) registers and HandleInput is invoked
  act:  c.AddSkill(NewEightDirectionalMovementSkill()); c.Update(space)
  post: no panic; with mocked input.Right=true, body.Velocity().X > 0 after model.Update
```

Test fixtures (`internal/kit/actors/beatemup/mocks_test.go` new):
- `fakeFS` providing a 1x1 PNG for each state asset.
- `stateMap` minimal (Idle only).
- `spriteData` with one asset, FrameRate=1, FacingDirection=right.
- `bodyRect`: `bodyphysics.NewRect(0,0,8,8)`.
- `space`: real `internal/engine/physics/space.Space` (no obstacles).

Mocks needed (none new at engine/contracts level; reuse `engine/mocks` if helpful).

---

## 9. File Edit Inventory

| File | Change |
|---|---|
| `internal/engine/data/schemas/json.go` | Add `Mode string` field to `MovementConfig` |
| `internal/engine/data/schemas/movement_mode.go` | NEW — `MovementMode*` constants |
| `internal/kit/skills/factory.go` | Switch on `cfg.Movement.Mode` |
| `internal/kit/skills/factory_test.go` | Add T-F1..T-F4 |
| `internal/kit/actors/beatemup/beatemup_character.go` | New constructor, embed `*actors.Character`, own `BeatEmUpMovementModel` |
| `internal/kit/actors/beatemup/beatemup_character_test.go` | Add T-B1..T-B4 |
| `internal/kit/actors/beatemup/mocks_test.go` | NEW — fixtures |
| `internal/game/scenes/phases/beatemup/scene.go` | `OnStart`: `Camera().SetBounds(&tilemapRect)` |
| `internal/game/entity/actors/player/cody*.go` | Update CodyPlayer construction to new `NewBeatEmUpCharacter` signature; ensure config sets `mode: "eight_dir"` |
| `assets/data/actors/cody*.json` (or equivalent) | Add `"mode": "eight_dir"` under skills.movement |

---

## 10. Contracts Inventory

No new engine contracts. Reuses:
- `body.MovableCollidable`, `body.BodiesSpace` (existing).
- `physicsmovement.MovementModel` (existing; `BeatEmUpMovementModel` already implements).
- `physicsmovement.PlayerMovementBlocker` (existing).
- `skill.Skill` (existing; `EightDirectionalMovementSkill` already implements).

→ Mock Generator step is NOT required for this story.
