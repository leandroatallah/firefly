# SPEC — 069-depth-lane-body-impl

**Goal:** Make `space.DepthLaneBody` operational by implementing it on the two body types that participate in beat-em-up collisions (`*body.ObstacleRect`, `*beatemup.BeatEmUpCharacter`), then retire the zero-altitude workaround in `BeatEmUpMovementModel.Update`.

**Interface (already exists, no change):** `internal/engine/physics/space/depth_lane.go`
```
type DepthLaneBody interface {
    GroundY() int        // pre-altitude world Y of the ground footprint (bottom edge for obstacles)
    LaneHalfWidth() int  // half-extent of accepted depth tolerance, in px
}
const DefaultLaneHalfWidth = 8
```
The gate already runs in `space.HasCollision` (`space.go:269-293`). No change to `space/`.

---

## 1. ObstacleRect implementation [AC-1, AC-3, AC-4]

**File:** `internal/engine/physics/body/obstacle.go`

Add two methods on `*ObstacleRect`:

```
func (o *ObstacleRect) GroundY() int
    rect := o.Position()              // o.Altitude is always 0 for obstacles → ground == screen
    return rect.Max.Y                 // bottom edge in world space

func (o *ObstacleRect) LaneHalfWidth() int
    h := o.GetShape().(*Rect).Height()
    if h <= 0 { return space.DefaultLaneHalfWidth }
    return h / 2
```

Import: `"github.com/boilerplate/ebiten-template/internal/engine/physics/space"` (engine→engine, no cycle).

Rationale notes (one-liners):
- `GroundY = Position().Max.Y` is the obstacle's bottom edge in world coords; obstacles are not airborne so this is altitude-independent already.
- `LaneHalfWidth = height/2` centers the lane on `GroundY - height/2` (the obstacle's mid-Y). With the gate `|diff| <= max(halfA, halfB)`, a character at `GroundY = obstacle.GroundY - h/2 ± h/2` collides. Net effect: lane covers `[obstacle.top, obstacle.bottom]` in depth.
- Wait: the gate compares `GroundY` values directly. Character `GroundY` is its bottom (feet). Obstacle `GroundY` is its bottom. Their difference is the gap between feet line and obstacle bottom line. Lane width must cover the obstacle's full depth extent so a character whose feet line falls anywhere across the obstacle band collides. → `h/2` is correct because the obstacle's effective center for matching is its bottom; tolerance `h/2` is too small. **Use `h` (full height) as tolerance** to make the lane equal to the obstacle's full Y-extent (player feet anywhere from obstacle top to obstacle bottom triggers).
- **Final formula:** `LaneHalfWidth = max(h, DefaultLaneHalfWidth)`.

Override:
```
func (o *ObstacleRect) LaneHalfWidth() int
    h := o.GetShape().(*Rect).Height()
    if h < space.DefaultLaneHalfWidth { return space.DefaultLaneHalfWidth }
    return h
```

---

## 2. BeatEmUpCharacter implementation [AC-2, AC-3, AC-4]

**File:** `internal/kit/actors/beatemup/beatemup_character.go`

Add two methods on `*BeatEmUpCharacter`:

```
func (c *BeatEmUpCharacter) GroundY() int
    _, y16 := c.GetPosition16()      // ground Y in fp16, pre-altitude (altitude lives separately)
    return y16 >> 16                 // == fp16.From16(y16)
    // top-left feet position; bbox check still uses Position() rects, so consistency is on the comparison line

func (c *BeatEmUpCharacter) LaneHalfWidth() int
    return space.DefaultLaneHalfWidth // 8 px
```

Import: `"github.com/boilerplate/ebiten-template/internal/engine/physics/space"` (kit→engine, allowed).

**Owner indirection check:** `ResolveCollisions` passes the parent (e.g. `*CodyPlayer`) into `HasCollision`. `CodyPlayer` embeds `*BeatEmUpCharacter`, so the methods are promoted automatically and `space.DepthLaneBody` is satisfied through embedding. No re-declaration needed in `internal/game/...`. Verify by `var _ space.DepthLaneBody = (*BeatEmUpCharacter)(nil)` in tests.

---

## 3. PlatformerCharacter must NOT implement [AC-5]

**File:** `internal/kit/actors/platformer/*.go` — **no change**.

Negative compile-time guard in test (`internal/kit/actors/platformer/depth_lane_test.go` or existing test file):
```
// Compile-time: PlatformerCharacter does NOT implement DepthLaneBody.
// (No assertion needed — verified by absence of GroundY/LaneHalfWidth methods.)
```
Optional explicit assertion:
```
T-P5: _, ok := any((*PlatformerCharacter)(nil)).(space.DepthLaneBody); ok == false
```

---

## 4. Movement model — remove zero-altitude wrap [AC-6, AC-7]

**File:** `internal/engine/physics/movement/movement_model_beatemup.go`

Replace lines 41-55 (Block 1) with:

```
_, _, blockX := b.ApplyValidPosition(vx16, true, space)
_, _, blockY := b.ApplyValidPosition(vy16, false, space)
```

Keep Block 2 (altitude integration + shape shift, lines 93-127) **unchanged**.

Remove the now-stale comment about "zeroing altitude temporarily". Leave a one-line comment: `// Wall/obstacle blocking is depth-gated via DepthLaneBody (story 069); no altitude wrap needed.`

---

## 5. Test plan — Red Phase

### 5.1 Space gate behavioural tests [AC-3, AC-4, AC-8]

**File:** `internal/engine/physics/space/depth_lane_test.go` (new)

Table-driven. Each row constructs two fake `DepthLaneBody`-implementing bodies (use an existing test helper or a local stub) with controllable bbox + GroundY + LaneHalfWidth.

```
T-S1: same depth + bbox overlap → collide
  pre:  a.bbox=(0,0,16,16) a.GroundY=100 a.Half=8;
        b.bbox=(8,8,24,24) b.GroundY=100 b.Half=8
  act:  HasCollision(a,b)
  post: == true

T-S2: different depth (diff > max half) + bbox overlap → no collide
  pre:  a.GroundY=100 a.Half=8; b.GroundY=120 b.Half=8 (diff=20 > 8)
  act:  HasCollision(a,b)
  post: == false

T-S3: airborne player + same-depth wall → collide
  pre:  player at altitude=40, GroundY=100; wall GroundY=100, Half=h
        player Position() Y is 60 (100-40), wall Position() Y is 100 → bboxes still overlap on X axis
        (set X overlap explicitly; vertical screen-gap closed by wide wall rect)
  act:  HasCollision(player, wall)
  post: == true       // blocks even airborne

T-S4: airborne player + different-depth wall → no collide
  pre:  player altitude=40 GroundY=100; wall GroundY=160 Half=8 (diff=60 > 8)
  act:  HasCollision(player, wall)
  post: == false       // jump does not false-block on background wall

T-S5: zero-height obstacle uses DefaultLaneHalfWidth
  pre:  obstacle constructed via NewObstacleRect with 0-height Rect → obstacle.LaneHalfWidth()==8
  act:  obstacle.LaneHalfWidth()
  post: == 8

T-S6: boundary inclusive (diff == max half)
  pre:  a.GroundY=100 a.Half=8; b.GroundY=108 b.Half=8 (diff=8)
  act:  HasCollision(a,b)
  post: == true
```

### 5.2 Interface satisfaction tests [AC-1, AC-2, AC-5]

```
T-I1: var _ space.DepthLaneBody = (*body.ObstacleRect)(nil)            // body package
T-I2: var _ space.DepthLaneBody = (*beatemup.BeatEmUpCharacter)(nil)   // kit/actors/beatemup
T-I3: ObstacleRect.GroundY returns Position().Max.Y
  pre:  obs.SetPosition(10,20); obs.SetSize(32,16)
  act:  obs.GroundY()
  post: == 36     // 20 + 16
T-I4: ObstacleRect.LaneHalfWidth returns max(height, DefaultLaneHalfWidth)
  pre:  obs.SetSize(32,4)
  act:  obs.LaneHalfWidth()
  post: == 8      // DefaultLaneHalfWidth (4 < 8)
  pre:  obs.SetSize(32,32)
  act:  obs.LaneHalfWidth()
  post: == 32
T-I5: BeatEmUpCharacter.GroundY == y16>>16, altitude-independent
  pre:  c.SetPosition(0,150); c.SetAltitude(40)
  act:  c.GroundY()
  post: == 150
T-I6: BeatEmUpCharacter.LaneHalfWidth == space.DefaultLaneHalfWidth
  post: == 8
T-I7 (negative): _, ok := any((*platformer.PlatformerCharacter)(nil)).(space.DepthLaneBody); ok == false
```

### 5.3 Movement model regression [AC-6, AC-9]

**File:** `internal/engine/physics/movement/movement_model_beatemup_test.go`

```
T-M1: airborne player not blocked by depth-mismatched wall
  pre:  player altitude=40 GroundY=100, vx16=+1000
        space contains ObstacleRect at world (Position) y=160, height=16 → wall GroundY=176, Half=16
        depth diff = |100 - 176| = 76; max half = 16 → gate denies
  act:  model.Update(player, space)
  post: player x advanced by ~vx16 (no block)

T-M2: airborne player blocked by same-depth wall
  pre:  player altitude=40 GroundY=100, vx16=+1000
        ObstacleRect at world y=92, height=16 → wall GroundY=108, Half=16
        depth diff = 8 ≤ 16 → gate allows; bboxes overlap if vx step lands inside
  act:  model.Update(player, space)
  post: player x clamped (blockX==true path); position not advanced past wall

T-M3: Block 1 (zero-altitude wrap) removed
  source-text assertion or behavioural: player.Altitude unchanged across Update when starting >0
  pre:  player.Altitude16 = 40<<16 at frame start
  act:  model.Update (no velocity)
  post: Altitude16 never observed as 0 inside ApplyValidPosition path (test via space stub recording
        every CollisionPosition call's Y values — never see groundY for airborne body)
```

---

## 6. File inventory

| File | Change |
|---|---|
| `internal/engine/physics/body/obstacle.go` | + `GroundY()`, `LaneHalfWidth()` methods; import `space` |
| `internal/kit/actors/beatemup/beatemup_character.go` | + `GroundY()`, `LaneHalfWidth()` methods; import `space` |
| `internal/engine/physics/movement/movement_model_beatemup.go` | remove Block 1 (lines 41-55 zero-out wrap); update comment |
| `internal/engine/physics/space/depth_lane_test.go` | NEW — table tests T-S1..S6 |
| `internal/engine/physics/body/obstacle_depth_lane_test.go` | NEW — T-I1, T-I3, T-I4 |
| `internal/kit/actors/beatemup/beatemup_character_depth_lane_test.go` | NEW — T-I2, T-I5, T-I6 |
| `internal/kit/actors/platformer/depth_lane_test.go` | NEW — T-I7 negative assertion |
| `internal/engine/physics/movement/movement_model_beatemup_test.go` | + T-M1, T-M2, T-M3 |

No new contracts, no new mocks. `DepthLaneBody` already exists in `space/`.

---

## 7. Post-conditions (overall)

- `var _ space.DepthLaneBody = (*body.ObstacleRect)(nil)` compiles.
- `var _ space.DepthLaneBody = (*beatemup.BeatEmUpCharacter)(nil)` compiles.
- `var _ space.DepthLaneBody = (*platformer.PlatformerCharacter)(nil)` does NOT compile.
- `BeatEmUpMovementModel.Update` no longer reads/writes `Altitude16` before `ApplyValidPosition`.
- All existing platformer tests pass unchanged.
- Migration: update `.agents/work/epics/beatemup-mechanics/PLAN_airborne-collision-split.md` header to "Option B chosen; Block 1 removed in story 069" (Gatekeeper step, not in spec).
