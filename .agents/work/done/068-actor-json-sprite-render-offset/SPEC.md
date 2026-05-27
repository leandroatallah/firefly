# SPEC — 068-actor-json-sprite-render-offset

## 1. Schema [AC-1, AC-2, AC-7, AC-8]

**File:** `internal/engine/data/schemas/json.go` (package `schemas`)

Add new type:

```go
// SpriteOffset is a pixel-space draw-time nudge applied to an actor sprite for
// a given state. Has zero effect on physics, collision, or footprint. Positive
// X moves the sprite right in screen space; positive Y moves it down.
type SpriteOffset struct {
    X int `json:"x"`
    Y int `json:"y"`
}
```

Extend `AssetData` (add field; keep existing fields in current order):

```go
type AssetData struct {
    Path           string        `json:"path"`
    CollisionRects []ShapeRect   `json:"collision_rect"`
    FootprintRect  *ShapeRect    `json:"footprint_rect,omitempty"`
    Loop           *bool         `json:"loop,omitempty"`
    RenderOffset   *SpriteOffset `json:"render_offset,omitempty"`
}
```

Constraints:
- `schemas` package imports unchanged (no new deps).
- Absent `render_offset` key → `AssetData.RenderOffset == nil`.
- Present with `{x:0,y:0}` → non-nil pointer to zero struct.

## 2. Character render-offset registry [AC-3, AC-5]

**File:** `internal/engine/entity/actors/character.go` (package `actors`)

Add field on `Character`:

```go
renderOffsets map[ActorStateEnum]image.Point // per-state pixel-space draw nudge
```

Initialize lazily (nil-safe). Add methods:

```go
// SetRenderOffset registers a per-state pixel-space draw-time offset.
// Has no effect on body position, velocity, collision rects, or footprint.
func (c *Character) SetRenderOffset(state ActorStateEnum, dx, dy int)

// RenderOffset returns the registered pixel-space offset for the given state.
// Returns image.Point{} (zero) and ok=false when no offset is registered.
func (c *Character) RenderOffset(state ActorStateEnum) (image.Point, bool)
```

Pseudocode:

```
SetRenderOffset(state, dx, dy):
  if c.renderOffsets == nil: c.renderOffsets = map
  c.renderOffsets[state] = image.Pt(dx, dy)

RenderOffset(state):
  if c.renderOffsets == nil: return Point{}, false
  p, ok := c.renderOffsets[state]
  return p, ok
```

## 3. Apply offset at draw time [AC-3, AC-4, AC-6]

**File:** `internal/engine/entity/actors/character.go`, function `UpdateImageOptions()`.

Append after step 4 (`Translate(float64(x), float64(y))`) — i.e. as the final transform:

```
// 5. Apply per-state render offset (screen-space, post-flip, post-anchor).
if p, ok := c.renderOffsets[c.state.State()]; ok {
    c.imageOptions.GeoM.Translate(float64(p.X), float64(p.Y))
}
```

Rules:
- The offset is read every call (no caching) → mid-frame state transitions pick up the new state's offset.
- Facing direction does NOT mirror `X`. The offset is applied AFTER the `Scale(-1,1)` flip, so a left-facing actor with `X:-4` is still nudged -4 px on screen. See NOTES.md for rationale.
- `image.Point{}` (zero) and missing key both produce zero translation.

## 4. Builder wiring [AC-3, AC-8]

**File:** `internal/engine/entity/actors/builder/builder.go` (package `builder`)

Add helper:

```go
// ApplyRenderOffsets reads SpriteData.Assets[*].RenderOffset and registers each
// non-nil offset on the character keyed by the state enum from stateMap. Assets
// with nil RenderOffset are skipped. Unknown state keys are skipped silently
// (consistent with BuildStateMap's strict pre-check).
func ApplyRenderOffsets(
    character actors.ActorEntity,
    data schemas.SpriteData,
    stateMap map[string]animation.SpriteState,
)
```

Pseudocode:

```
for key, asset in data.Assets:
  if asset.RenderOffset == nil: continue
  st, ok := stateMap[key]; if !ok: continue
  enum, ok := st.(actors.ActorStateEnum); if !ok: continue
  character.GetCharacter().SetRenderOffset(enum, asset.RenderOffset.X, asset.RenderOffset.Y)
```

Call site: invoked by genre kit constructors (e.g. `BeatEmUpCharacter`) after `NewCharacter`. Engine builder does not auto-call; kit decides when. (Mirrors current pattern for `buildFootprints` in `internal/kit/actors/beatemup/beatemup_character.go`.)

**Kit integration (beatemup):** in `NewBeatEmUpCharacter` after `c := actors.NewCharacter(...)`:

```
builder.ApplyRenderOffsets(c, spriteData, stateMap)
```

(Adds import of `internal/engine/entity/actors/builder` if not already present. Builder already imports schemas and actors.)

## 5. Mock / contract inventory

- No new contracts in `internal/engine/contracts/`.
- No new shared mocks. `Mock Generator` stage → SKIP.

## 6. Pre-conditions

- `Character.renderOffsets` is nil on a fresh `NewCharacter`.
- `RenderOffset(s)` returns `(Point{}, false)` before any `SetRenderOffset`.
- `UpdateImageOptions()` is callable with `renderOffsets == nil` (no panic).

## 7. Post-conditions

- After `SetRenderOffset(state, dx, dy)`, `RenderOffset(state) == (image.Pt(dx,dy), true)`.
- `UpdateImageOptions()` with a registered offset for the current state produces an `imageOptions.GeoM` whose final-row translation equals `baselineTranslation + (dx, dy)`.
- `UpdateImageOptions()` with no registered offset produces an identical `GeoM` to the pre-feature behavior (zero regression).
- `Character.Position()`, `Velocity()`, `CollisionPosition()`, and (for beatemup) `Footprint()` are bit-for-bit unchanged regardless of registered offsets.

## 8. Red Phase test scenarios

### T-S1: AssetData unmarshals render_offset present
`internal/engine/data/schemas/json_test.go`
```
pre:  raw = {"path":"p","collision_rect":[],"render_offset":{"x":-4,"y":2}}
act:  json.Unmarshal(raw, &a)
post: a.RenderOffset != nil; a.RenderOffset.X == -4; a.RenderOffset.Y == 2
```

### T-S2: AssetData unmarshals render_offset absent
`internal/engine/data/schemas/json_test.go`
```
pre:  raw = {"path":"p","collision_rect":[]}
act:  json.Unmarshal(raw, &a)
post: a.RenderOffset == nil
```

### T-S3: AssetData unmarshals render_offset zero
```
pre:  raw = {"path":"p","collision_rect":[],"render_offset":{"x":0,"y":0}}
act:  json.Unmarshal(raw, &a)
post: a.RenderOffset != nil; *a.RenderOffset == SpriteOffset{}
```

### T-C1: Character SetRenderOffset / RenderOffset round-trip
`internal/engine/entity/actors/character_test.go` (or new file)
```
pre:  c := NewCharacter(emptySpriteMap, rect)
act:  c.SetRenderOffset(Idle, -4, 2)
post: p, ok := c.RenderOffset(Idle); ok==true; p == image.Pt(-4, 2)
      p2, ok2 := c.RenderOffset(Walking); ok2==false; p2 == image.Point{}
```

### T-C2: UpdateImageOptions applies offset as final translation [AC-3, AC-6]
Table-driven. For each row: build a Character with a known body rect and a sprite map containing one ebiten.NewImage(32,32) sprite for `Idle`. Capture `c.ImageOptions().GeoM` element (0,2) and (1,2) (translation row).
```
row "no offset registered":
  pre:  c.renderOffsets == nil
  act:  c.UpdateImageOptions()
  post: tx, ty == baseline_tx, baseline_ty   ← computed from current logic

row "offset (-4,0) at current state":
  pre:  c.SetRenderOffset(Idle, -4, 0); state == Idle
  act:  c.UpdateImageOptions()
  post: tx, ty == baseline_tx - 4, baseline_ty

row "offset (0, 3) at current state":
  pre:  c.SetRenderOffset(Idle, 0, 3); state == Idle
  act:  c.UpdateImageOptions()
  post: tx, ty == baseline_tx, baseline_ty + 3

row "offset registered for other state, current=Idle":
  pre:  c.SetRenderOffset(Walking, -10, -10); state == Idle
  act:  c.UpdateImageOptions()
  post: tx, ty == baseline_tx, baseline_ty   ← unaffected

row "offset (0,0) equivalent to no offset":
  pre:  c.SetRenderOffset(Idle, 0, 0); state == Idle
  act:  c.UpdateImageOptions()
  post: tx, ty == baseline_tx, baseline_ty
```

Helper: compute `baseline_tx, baseline_ty` by running `UpdateImageOptions()` first with no offsets registered, snapshotting the GeoM translation row, then resetting and applying the offset.

### T-C3: Facing-left does NOT mirror X [AC-9]
```
pre:  c.SetRenderOffset(Idle, -4, 0); c.SetFaceDirection(FaceDirectionLeft)
act:  c.UpdateImageOptions()
post: GeoM translation X == baseline_left_tx + (-4)   ← additive, not +4
```

### T-C4: Offset does not move body / collision / footprint [AC-5]
```
pre:  c := NewBeatEmUpCharacter(... assets with RenderOffset {x:-20, y:-20} on Idle ...)
      posBefore := c.Position(); collBefore := c.CollisionPosition(); fpBefore := c.Footprint()
act:  c.UpdateImageOptions()
post: c.Position() == posBefore
      c.CollisionPosition() deep-equals collBefore
      c.Footprint() == fpBefore
```

### T-B1: ApplyRenderOffsets registers per-state offsets [AC-3]
`internal/engine/entity/actors/builder/builder_test.go`
```
pre:  spriteData.Assets = {
        "idle":  {Path:"i.png", RenderOffset:&SpriteOffset{X:-4, Y:0}},
        "melee": {Path:"m.png", RenderOffset:&SpriteOffset{X:8, Y:-2}},
        "walking": {Path:"w.png"},  // no offset
      }
      stateMap built via BuildStateMap
      c := actors.NewCharacter(emptySpriteMap, rect)
act:  builder.ApplyRenderOffsets(c, spriteData, stateMap)
post: c.RenderOffset(Idle)    == (Pt(-4, 0), true)
      c.RenderOffset(Melee)   == (Pt( 8,-2), true)
      c.RenderOffset(Walking) == (Pt(0,0),  false)
```

### T-B2: ApplyRenderOffsets is a no-op when all RenderOffset are nil [AC-4, AC-8]
```
pre:  spriteData with two assets, both RenderOffset == nil
act:  builder.ApplyRenderOffsets(c, spriteData, stateMap)
post: for every state s in stateMap: c.RenderOffset(s) returns (Point{}, false)
```

## 9. Files to create / modify

Modify:
- `internal/engine/data/schemas/json.go` — add `SpriteOffset`; add `RenderOffset` to `AssetData`.
- `internal/engine/data/schemas/json_test.go` — T-S1, T-S2, T-S3.
- `internal/engine/entity/actors/character.go` — add `renderOffsets`, `SetRenderOffset`, `RenderOffset`, apply in `UpdateImageOptions`.
- `internal/engine/entity/actors/character_test.go` (new or existing) — T-C1…T-C4.
- `internal/engine/entity/actors/builder/builder.go` — add `ApplyRenderOffsets`.
- `internal/engine/entity/actors/builder/builder_test.go` — T-B1, T-B2.
- `internal/kit/actors/beatemup/beatemup_character.go` — call `builder.ApplyRenderOffsets` in `NewBeatEmUpCharacter`.

Do NOT modify:
- Any actor JSON file under `assets/`.
- `internal/engine/physics/**`, `internal/engine/render/sprites/**`.
