# SPEC — 070-render-offset-facing-kit-wiring

> **Post-implementation revision (2026-05-27):** Auto-mirror semantics. The
> `XFlipped *int` override was removed entirely. `SpriteOffset` is now
> `{X, Y}` only. `Character.SetRenderOffset(state, dx, dy)` is 3-arg. When
> facing left, `RenderOffset()` and `UpdateImageOptions` step 5 negate `X`
> automatically (auto-mirror), since the offset is applied post-flip in
> screen space and must follow the mirrored sprite content. The original
> design preserved `X` across facings (068 regression guard) and offered
> per-asset opt-in via `x_flipped`. In practice every real use case wants
> the negated value, so YAGNI applied: removed the field, the JSON key, the
> override tests, and the explicit-zero path. The sections below describing
> `XFlipped`, `*int` plumbing, T-S1/T-S3, and the "X only, facing left →
> same X" T-C2 row are superseded.

---


## 1. Schema [AC-1, AC-2, AC-8]

**File:** `internal/engine/data/schemas/json.go` (package `schemas`)

Extend `SpriteOffset` (add field; keep `X`, `Y` in current order):

```go
type SpriteOffset struct {
    X        int  `json:"x"`
    Y        int  `json:"y"`
    XFlipped *int `json:"x_flipped,omitempty"`
}
```

Constraints:
- `schemas` package imports unchanged (no `kit`/`game` import).
- Absent `x_flipped` key → `SpriteOffset.XFlipped == nil`.
- Present with `"x_flipped": 0` → non-nil pointer to int 0 (explicit override).

## 2. Character render-offset registry [AC-4, AC-5]

**File:** `internal/engine/entity/actors/character.go` (package `actors`)

Replace the per-state storage to carry the optional flipped override. Define an internal value type (unexported) and update field + accessors:

```go
type renderOffset struct {
    X        int
    Y        int
    XFlipped *int // nil → use X for left-facing
}

// On Character (replaces the existing map[ActorStateEnum]image.Point):
renderOffsets map[ActorStateEnum]renderOffset
```

API:

```go
// SetRenderOffset registers a per-state pixel-space draw-time offset.
// dxFlipped is the optional left-facing X override; pass nil to reuse dx for both facings.
// dxFlipped == &0 is a valid explicit zero override (distinct from nil).
func (c *Character) SetRenderOffset(state ActorStateEnum, dx, dy int, dxFlipped *int)

// RenderOffset returns the resolved offset for the given state at the current facing.
// Kept for backward compatibility with story 068 tests: returns image.Point with X
// resolved against c.FaceDirection() at call time.
func (c *Character) RenderOffset(state ActorStateEnum) (image.Point, bool)
```

Pseudocode:

```
SetRenderOffset(state, dx, dy, dxFlipped):
  if c.renderOffsets == nil: c.renderOffsets = map
  c.renderOffsets[state] = renderOffset{X:dx, Y:dy, XFlipped:dxFlipped}

resolveOffsetX(o renderOffset, dir FacingDirection):
  if dir == FaceDirectionLeft && o.XFlipped != nil: return *o.XFlipped
  return o.X

RenderOffset(state):
  if c.renderOffsets == nil: return Point{}, false
  o, ok := c.renderOffsets[state]; if !ok: return Point{}, false
  return image.Pt(resolveOffsetX(o, c.FaceDirection()), o.Y), true
```

Migration: story 068 callers `SetRenderOffset(state, dx, dy)` → update signature to `SetRenderOffset(state, dx, dy, nil)`. Update existing tests accordingly (T-C1 etc remain valid with `nil`).

## 3. Apply offset at draw time [AC-3]

**File:** `internal/engine/entity/actors/character.go`, function `UpdateImageOptions()`.

Replace step 5 (current `Translate(p.X, p.Y)`) so X is facing-resolved:

```
// 5. Per-state render offset (screen-space, post-flip, post-anchor).
if c.renderOffsets != nil {
    if o, ok := c.renderOffsets[c.state.State()]; ok {
        x := o.X
        if fDirection == animation.FaceDirectionLeft && o.XFlipped != nil {
            x = *o.XFlipped
        }
        c.imageOptions.GeoM.Translate(float64(x), float64(o.Y))
    }
}
```

Rules:
- Lookup happens every call (no caching). `fDirection` is the local variable already resolved earlier in `UpdateImageOptions` (line ~328-336).
- `Y` is identical for both facings.
- `XFlipped == &0` produces zero X translation for left-facing — distinct from nil (which would fall back to `X`).

## 4. Builder wiring [AC-5, AC-8]

**File:** `internal/engine/entity/actors/builder/builder.go`

Update `ApplyRenderOffsets` to forward the `XFlipped` pointer:

```go
func ApplyRenderOffsets(
    character actors.ActorEntity,
    data schemas.SpriteData,
    stateMap map[string]animation.SpriteState,
) {
    for key, asset := range data.Assets {
        if asset.RenderOffset == nil { continue }
        st, ok := stateMap[key]; if !ok { continue }
        enum, ok := st.(actors.ActorStateEnum); if !ok { continue }
        character.GetCharacter().SetRenderOffset(
            enum,
            asset.RenderOffset.X,
            asset.RenderOffset.Y,
            asset.RenderOffset.XFlipped, // forward pointer as-is
        )
    }
}
```

Constraint: pointer is passed through unchanged so an explicit `0` survives.

## 5. Platformer kit wiring [AC-6, AC-7]

**File:** `internal/kit/actors/platformer/platformer.go`, function `NewPlatformerCharacter`.

Add import: `"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"`.

Insert immediately after `c := actors.NewCharacter(s, bodyRect)` (line 102) and BEFORE the `pf := &PlatformerCharacter{...}` construction (mirrors beatemup_character.go line 67):

```go
c := actors.NewCharacter(s, bodyRect)
builder.ApplyRenderOffsets(c, spriteData, stateMap)
pf := &PlatformerCharacter{ Character: c }
```

No other changes to platformer constructor; existing JSONs without `render_offset` produce a no-op (AC-7).

## 6. Mock / contract inventory

- No new contracts in `internal/engine/contracts/`.
- No new shared mocks. `Mock Generator` stage → SKIP.

## 7. Pre-conditions

- `Character.renderOffsets` is nil on fresh `NewCharacter` (unchanged).
- `RenderOffset(s)` returns `(Point{}, false)` before any `SetRenderOffset` (unchanged).
- `SpriteOffset{X:5}` (literal) parses with `XFlipped == nil`.
- `FaceDirection()` is callable on a fresh `Character` without panic (current behavior).

## 8. Post-conditions

- After `SetRenderOffset(s, dx, dy, nil)`: `RenderOffset(s)` returns `(Pt(dx,dy), true)` for both facings.
- After `SetRenderOffset(s, dx, dy, &dxf)`: `RenderOffset(s)` returns `(Pt(dx,dy), true)` when facing right; `(Pt(dxf,dy), true)` when facing left.
- `UpdateImageOptions()` translation row equals `baseline + (resolvedX, dy)` where `resolvedX = dxf if (left && flipped!=nil) else dx`.
- Body position, velocity, collision rects, footprint, altitude — bit-for-bit unchanged.
- `schemas` package import set unchanged (verified by `builder_layering_test.go` style check or `go list -deps`).

## 9. Red Phase test scenarios

### T-S1: SpriteOffset unmarshals x_flipped present [AC-1]
`internal/engine/data/schemas/json_test.go`
```
pre:  raw = {"x":-4,"y":2,"x_flipped":6}
act:  json.Unmarshal(raw, &o)
post: o.X==-4; o.Y==2; o.XFlipped!=nil; *o.XFlipped==6
```

### T-S2: SpriteOffset unmarshals x_flipped absent [AC-1, AC-2]
```
pre:  raw = {"x":-4,"y":2}
act:  json.Unmarshal(raw, &o)
post: o.XFlipped == nil
```

### T-S3: SpriteOffset unmarshals x_flipped zero (explicit override) [AC-1]
```
pre:  raw = {"x":-4,"y":2,"x_flipped":0}
act:  json.Unmarshal(raw, &o)
post: o.XFlipped!=nil; *o.XFlipped==0
```

### T-C1: SetRenderOffset / RenderOffset round-trip [AC-4]
`internal/engine/entity/actors/character_test.go` (table-driven)
```
row "nil flipped, facing right":
  pre:  c.SetRenderOffset(Idle, -4, 2, nil); c.SetFaceDirection(FaceDirectionRight)
  post: p,ok := c.RenderOffset(Idle); ok; p==Pt(-4,2)

row "nil flipped, facing left → falls back to X":
  pre:  c.SetRenderOffset(Idle, -4, 2, nil); c.SetFaceDirection(FaceDirectionLeft)
  post: p,ok := c.RenderOffset(Idle); ok; p==Pt(-4,2)

row "flipped set, facing right uses X":
  pre:  f := 6; c.SetRenderOffset(Idle, -4, 2, &f); c.SetFaceDirection(FaceDirectionRight)
  post: p,ok := c.RenderOffset(Idle); ok; p==Pt(-4,2)

row "flipped set, facing left uses XFlipped":
  pre:  f := 6; c.SetRenderOffset(Idle, -4, 2, &f); c.SetFaceDirection(FaceDirectionLeft)
  post: p,ok := c.RenderOffset(Idle); ok; p==Pt(6,2)

row "flipped=0 explicit override, facing left":
  pre:  f := 0; c.SetRenderOffset(Idle, -4, 2, &f); c.SetFaceDirection(FaceDirectionLeft)
  post: p,ok := c.RenderOffset(Idle); ok; p==Pt(0,2)

row "unregistered state":
  pre:  c.SetRenderOffset(Idle, -4, 2, nil)
  post: p,ok := c.RenderOffset(Walking); !ok; p==Point{}
```

### T-C2: UpdateImageOptions applies facing-resolved offset [AC-3, AC-9]
Table-driven. Build a Character with a 32x32 sprite for `Idle`, capture `GeoM` element (0,2) and (1,2). Compute baseline by calling `UpdateImageOptions` once with no offsets registered for the matching facing.
```
row "no offset":
  pre:  state=Idle, no offsets
  post: tx,ty == baseline_R_tx, baseline_R_ty

row "X only, facing right":
  pre:  SetRenderOffset(Idle, -4, 2, nil); face=Right
  post: tx,ty == baseline_R_tx - 4, baseline_R_ty + 2

row "X only, facing left → same X (068 regression)":
  pre:  SetRenderOffset(Idle, -4, 2, nil); face=Left
  post: tx,ty == baseline_L_tx - 4, baseline_L_ty + 2

row "X + XFlipped, facing right → uses X":
  pre:  f:=6; SetRenderOffset(Idle, -4, 2, &f); face=Right
  post: tx,ty == baseline_R_tx - 4, baseline_R_ty + 2

row "X + XFlipped, facing left → uses XFlipped":
  pre:  f:=6; SetRenderOffset(Idle, -4, 2, &f); face=Left
  post: tx,ty == baseline_L_tx + 6, baseline_L_ty + 2

row "XFlipped=0 explicit, facing left → 0":
  pre:  f:=0; SetRenderOffset(Idle, -4, 2, &f); face=Left
  post: tx,ty == baseline_L_tx + 0, baseline_L_ty + 2

row "Y identical for both facings":
  pre:  f:=6; SetRenderOffset(Idle, -4, 3, &f)
  post: ty(left) - baseline_L_ty == ty(right) - baseline_R_ty == 3
```

Note: facing is forced by `SetAcceleration(0,0)` followed by `SetFaceDirection(...)` — `UpdateImageOptions` only overrides `fDirection` when `accX != 0`.

### T-B1: ApplyRenderOffsets propagates XFlipped [AC-5, AC-10]
`internal/engine/entity/actors/builder/builder_test.go` (table-driven)
```
row "asset with XFlipped":
  pre:  six := 6; spriteData.Assets = {"idle": {Path:"i.png", RenderOffset:&SpriteOffset{X:-4, Y:2, XFlipped:&six}}}
  act:  builder.ApplyRenderOffsets(c, spriteData, stateMap)
  post: face=Right → c.RenderOffset(Idle) == (Pt(-4,2), true)
        face=Left  → c.RenderOffset(Idle) == (Pt( 6,2), true)

row "asset without XFlipped":
  pre:  spriteData.Assets = {"idle": {RenderOffset:&SpriteOffset{X:-4, Y:2}}}
  act:  builder.ApplyRenderOffsets(c, spriteData, stateMap)
  post: face=Right → (Pt(-4,2), true)
        face=Left  → (Pt(-4,2), true)   // falls back to X

row "asset with XFlipped=0 explicit":
  pre:  zero:=0; RenderOffset:&SpriteOffset{X:-4, Y:2, XFlipped:&zero}
  act:  ApplyRenderOffsets
  post: face=Left  → (Pt(0,2), true)
```

### T-B2: schemas layering unchanged [AC-8]
Extend `internal/engine/entity/actors/builder/builder_layering_test.go` (or add `internal/engine/data/schemas/layering_test.go`):
```
post: package "schemas" Imports() does not contain prefix "internal/kit/" or "internal/game/"
```

### T-P1: NewPlatformerCharacter applies render offsets [AC-6, AC-11]
`internal/kit/actors/platformer/platformer_test.go`
```
pre:  spriteData with Assets["idle"] = AssetData{
        Path: <test png>, CollisionRects: [...],
        RenderOffset: &SpriteOffset{X:-2, Y:0},
      }
      stateMap = BuildStateMap(spriteData)
act:  pf, err := NewPlatformerCharacter(fsys, stateMap, spriteData, bodyRect)
post: err == nil
      p, ok := pf.RenderOffset(actors.Idle); ok == true; p == image.Pt(-2, 0)
```

### T-P2: NewPlatformerCharacter without render_offset is a no-op [AC-7]
```
pre:  spriteData with no RenderOffset on any asset
act:  pf, err := NewPlatformerCharacter(...)
post: err == nil
      for every registered state s: _, ok := pf.RenderOffset(s); ok == false
```

## 10. Files to create / modify

Modify:
- `internal/engine/data/schemas/json.go` — add `XFlipped *int` to `SpriteOffset`.
- `internal/engine/data/schemas/json_test.go` — T-S1, T-S2, T-S3.
- `internal/engine/entity/actors/character.go` — change `renderOffsets` map value type, update `SetRenderOffset` signature, update `RenderOffset` to facing-resolve, update step 5 in `UpdateImageOptions`.
- `internal/engine/entity/actors/character_test.go` — T-C1, T-C2; migrate any existing 068 callers to pass `nil`.
- `internal/engine/entity/actors/builder/builder.go` — forward `XFlipped` in `ApplyRenderOffsets`.
- `internal/engine/entity/actors/builder/builder_test.go` — T-B1; migrate 068 expectations.
- `internal/engine/entity/actors/builder/builder_layering_test.go` (or new schemas layering test) — T-B2.
- `internal/kit/actors/platformer/platformer.go` — add `builder` import + `ApplyRenderOffsets` call.
- `internal/kit/actors/platformer/platformer_test.go` — T-P1, T-P2.

Do NOT modify:
- `internal/kit/actors/beatemup/beatemup_character.go` (already calls `ApplyRenderOffsets`).
- `internal/kit/actors/shooter_character.go` (not a constructor with `SpriteData`; out of scope).
- Any `assets/**/*.json` (no asset changes in this story; tuning is a follow-up).
- `internal/engine/physics/**`, `internal/engine/render/sprites/**`.
