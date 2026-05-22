# SPEC — 064-beatemup-footprint-rect

## 1. Schema [AC-1, AC-7]

File: `internal/engine/data/schemas/json.go`

Extend `AssetData`:

```go
type AssetData struct {
    Path           string      `json:"path"`
    CollisionRects []ShapeRect `json:"collision_rect"`
    FootprintRect  *ShapeRect  `json:"footprint_rect,omitempty"` // NEW
    Loop           *bool       `json:"loop,omitempty"`
}
```

Constraints:
- Pointer type — absent field unmarshals to `nil`.
- No new imports in `schemas/`. No reference to `internal/kit/` or `internal/game/`.

## 2. BeatEmUp Footprint Storage [AC-2, AC-3, AC-6]

File: `internal/kit/actors/beatemup/beatemup_character.go`

Add per-state footprint map to `BeatEmUpCharacter`:

```go
type BeatEmUpCharacter struct {
    *actors.Character
    *kitactors.MeleeCharacter
    footprints map[actors.ActorStateEnum]image.Rectangle // local rect, NOT world-offset
}
```

Helper (package-private):

```go
func buildFootprints(
    assets map[string]schemas.AssetData,
    stateMap map[string]animation.SpriteState,
) map[actors.ActorStateEnum]image.Rectangle
```

Pseudocode:
```
buildFootprints:
  out := {}
  for key, asset := range assets:
    if asset.FootprintRect == nil: continue
    st, ok := stateMap[key]; if !ok: continue
    enumSt, ok := st.(actors.ActorStateEnum); if !ok: continue
    r := asset.FootprintRect
    if r.Width <= 0 || r.Height <= 0: continue   // zero-size → treat as absent
    out[enumSt] = image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
  return out
```

`NewBeatEmUpCharacter` populates `be.footprints = buildFootprints(spriteData.Assets, stateMap)` after embedding. Empty map is valid.

## 3. Footprint() Method [AC-2, AC-3, AC-8]

File: `internal/kit/actors/beatemup/beatemup_character.go`

```go
// Footprint returns the current state's footprint rectangle in WORLD coordinates.
// Falls back to the union of the actor's full collision rects when:
//   - the current state has no footprint_rect declared, OR
//   - the footprint map is nil/empty.
// If no collision rects exist either, returns the body Position().
func (c *BeatEmUpCharacter) Footprint() image.Rectangle
```

Pseudocode:
```
Footprint:
  st := c.State()
  local, hasFootprint := c.footprints[st]
  if hasFootprint:
    minX, minY := c.GetPositionMin()
    return local.Add(image.Pt(minX, minY))
  // fallback path (AC-3)
  rects := c.CollidableBody.CollisionPosition()   // explicit embedded selector
  if len(rects) == 0:
    return c.Position()
  return unionRects(rects)

unionRects(rs):
  u := rs[0]
  for i := 1..len(rs)-1: u = u.Union(rs[i])
  return u
```

Notes:
- Local rect stored unmirrored. Mirroring for `FaceDirectionLeft` is not applied (matches existing `CollidableBody` behavior, which also stores rects as-authored and only mirrors at draw time). See AC-8 below.
- World offset uses `GetPositionMin()` (top-left in screen-space), consistent with how `SetCollisionBodies` initialises collision rects relative to body origin.

## 4. CollisionPosition Override [AC-4, AC-5, AC-6]

File: `internal/kit/actors/beatemup/beatemup_character.go`

Override so engine-level checks (`space.HasCollision`, `CollidableBody.ApplyValidPosition` → `space.ResolveCollisions`) consume the footprint for beatemup actors only.

```go
// CollisionPosition shadows the embedded *CollidableBody method.
// When a footprint exists for the current state, it returns ONLY the footprint
// so movement-vs-world and actor-vs-actor checks use the feet area (AC-4).
// When absent, returns the embedded full collision rects (AC-3 fallback, AC-5
// attack-side resolution paths continue to see the body).
func (c *BeatEmUpCharacter) CollisionPosition() []image.Rectangle
```

Pseudocode:
```
CollisionPosition:
  st := c.State()
  if local, ok := c.footprints[st]:
    minX, minY := c.GetPositionMin()
    return []image.Rectangle{ local.Add(image.Pt(minX, minY)) }
  return c.CollidableBody.CollisionPosition()
```

Pre/post:
- pre: Method dispatch reaches `*BeatEmUpCharacter` via interface `body.Collidable`. Verified because `space.bodies` stores the registered `Touchable`/`Collidable`, and `NewBeatEmUpCharacter` is the concrete type registered upstream by builder.
- post: `space.HasCollision(beatemupActor, other)` evaluates AABB on the footprint world-rect.
- post: `b.ApplyValidPosition` step loop sees footprint via `ResolveCollisions` → `HasCollision`.

## 5. Attack Hitbox Path Unchanged [AC-5]

No code change. Verification only:
- `internal/kit/combat/melee/controller.go` allocates an independent hitbox body (not the actor's `CollidableBody`) — see line 202 region.
- `MeleeWeapon` uses `BodyRect`/`CollisionRects` for damage resolution. `Footprint()` is never called from combat.

Gatekeeper test: confirm no new caller of `Footprint()` outside `internal/kit/actors/beatemup/`.

## 6. Layer Compliance [AC-6, AC-7]

- `internal/engine/data/schemas/json.go` — only adds a new pointer field; imports unchanged.
- `internal/kit/actors/beatemup/` — already imports `schemas`, `actors`. No new external dep required.
- Platformer (`internal/kit/actors/platformer/`) and topdown actors do not implement or call `Footprint()`. `AssetData.FootprintRect` is silently nil for them and ignored.

## 7. Edge Cases (Decisions)

| Case | Decision |
|---|---|
| Zero width or height in `footprint_rect` | Treat as no-footprint; `Footprint()` falls back per AC-3 (see `buildFootprints` filter). |
| State with no entry in `footprints` map (e.g. airborne jump frame) | Fallback per AC-3. |
| `CollisionRects` empty AND no footprint | `Footprint()` returns `Position()` (body rect). `CollisionPosition()` falls back to embedded behavior (empty slice; engine `collisionRects` helper substitutes body `Position()`). |
| Multiple `CollisionRects` + single `footprint_rect` | Footprint is single; `CollisionPosition` returns slice of length 1 when footprint applies. |
| State change mid-frame | `c.State()` always reflects current state — `Footprint()` re-reads each call. No caching. |

## 8. Red Phase Test Triples

Test file: `internal/engine/data/schemas/json_test.go`

```
T-S1: AssetData unmarshals with footprint_rect present
  pre:  json = `{"path":"p","collision_rect":[],"footprint_rect":{"x":2,"y":4,"width":10,"height":3}}`
  act:  json.Unmarshal(raw, &a)
  post: a.FootprintRect != nil && a.FootprintRect.Rect() == (2,4,10,3)

T-S2: AssetData unmarshals when footprint_rect absent
  pre:  json = `{"path":"p","collision_rect":[]}`
  act:  json.Unmarshal(raw, &a)
  post: err == nil && a.FootprintRect == nil
```

Test file: `internal/kit/actors/beatemup/beatemup_character_test.go`

Test fixture helper extension: allow setting an `Assets` map with a footprint entry and a `stateMap` that maps the key to `actors.Idle`.

```
T-F1: Footprint() returns world-offset rect when state has footprint_rect
  pre:  Assets["idle"].FootprintRect = {X:2,Y:30,W:12,H:6};
        stateMap["idle"] = actors.Idle; actor at body min (100, 200); state == Idle
  act:  r := c.Footprint()
  post: r == image.Rect(102, 230, 114, 236)

T-F2: Footprint() falls back to collision rect union when no footprint for state
  pre:  Assets["idle"] has no FootprintRect; actor has 1 collision rect equal to Position();
        state == Idle
  act:  r := c.Footprint()
  post: r == c.Position()

T-F3: Footprint() falls back to Position() when no collisions AND no footprint
  pre:  Assets empty; no collision rects; body Position == image.Rect(0,0,8,8)
  act:  r := c.Footprint()
  post: r == image.Rect(0,0,8,8)

T-F4: CollisionPosition() returns ONLY footprint world-rect when footprint exists
  pre:  Same as T-F1
  act:  rs := c.CollisionPosition()
  post: len(rs) == 1 && rs[0] == image.Rect(102,230,114,236)

T-F5: CollisionPosition() falls back to embedded behavior when no footprint
  pre:  Same as T-F2
  act:  rs := c.CollisionPosition()
  post: rs == c.CollidableBody.CollisionPosition()  // deep-equal slice

T-F6: Zero-size footprint_rect treated as absent
  pre:  Assets["idle"].FootprintRect = {X:0,Y:0,W:0,H:0}; state==Idle
  act:  r := c.Footprint()
  post: r == c.Position()   // fallback path

T-F7: State change updates Footprint() target
  pre:  Assets has footprint for "idle" only; current state Idle → switch to Walking
        (Walking has no footprint entry)
  act:  r1 := c.Footprint(); c.SetNewStateFatal(actors.Walking); r2 := c.Footprint()
  post: r1 reflects idle footprint world-rect; r2 falls back per AC-3

T-F8: Facing-left does NOT mirror footprint (parity with collision_rect today)
  pre:  Assets["idle"].FootprintRect = {X:2,Y:30,W:12,H:6}; FaceDirection=Left;
        body min (100, 200); state == Idle
  act:  r := c.Footprint()
  post: r == image.Rect(102, 230, 114, 236)   // same as T-F1
```

Integration-style test (still unit-scope, no GPU):

```
T-I1: space.HasCollision uses footprint for beatemup actor
  pre:  obstacle body at (0,0,200,50);
        beatemup actor body at (0, 100, 32, 64) i.e. head far below obstacle base;
        footprint at local (4, 56, 24, 8) → world (4, 156, 24, 8) — well clear
  act:  space.HasCollision(actor, obstacle)
  post: false   // body bbox would overlap if footprint were not used? Adjust setup:
        // place obstacle at (0,140,200,200) and actor body (0,100, 32, 64);
        // body overlaps; footprint at world (4,156,28,164) overlaps too → true.
        // Two sub-cases:
        //   case A (body overlaps, feet do not) → false
        //   case B (feet overlap)               → true
```

Concrete sub-case data:
- A: obstacle rect (0,100,200,140); body rect (0,100,200,200); footprint world (4,180,28,188). Body overlaps obstacle on the head band; footprint does NOT. Expect `false`.
- B: Same obstacle; footprint world (4,110,28,135). Expect `true`.

## 9. Mock / Contract Inventory

- No new contract interface introduced.
- No new mocks required. (Mock Generator may be SKIPPED.)
- Existing `body.Collidable` interface dispatch is sufficient.

## 10. File Touch List

- MODIFY `internal/engine/data/schemas/json.go` — add `FootprintRect *ShapeRect`.
- MODIFY `internal/engine/data/schemas/json_test.go` — add T-S1, T-S2.
- MODIFY `internal/kit/actors/beatemup/beatemup_character.go` — add `footprints` field, `buildFootprints`, `Footprint()`, `CollisionPosition()` override.
- MODIFY `internal/kit/actors/beatemup/beatemup_character_test.go` — add T-F1..T-F8, T-I1.

No other files require modification.
