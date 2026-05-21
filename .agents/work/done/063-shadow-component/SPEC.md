# SPEC — 063-shadow-component

**Bounded Context:** Kit (`internal/kit/render/shadow/`)
**Layer rule:** kit may import engine; must NOT import `internal/game/`.

---

## 1. Package & Files [AC-1, AC-6]

New package: `internal/kit/render/shadow/`

- `internal/kit/render/shadow/shadow.go` — types, constants, `Draw` API.
- `internal/kit/render/shadow/shadow_test.go` — Red-Phase tests.

Imports allowed: `image/color`, `github.com/hajimehoshi/ebiten/v2`, `github.com/hajimehoshi/ebiten/v2/vector`, `internal/engine/contracts/body`, `internal/engine/render/camera`, `internal/engine/render/draworder` (for the `Altitudable` interface only — re-import or duplicate locally).

Forbidden imports: `internal/game/...`, `internal/kit/scenes/...`.

---

## 2. Constants [AC-3, AC-4]

```
const (
    ShadowAlpha          = 0.50  // 50% opacity
    ShadowBaseWidthRatio = 0.75  // oval width = body.Shape.Width() * ratio at altitude=0
    ShadowBaseHeight     = 4     // oval height in pixels at altitude=0 (flat oval)
    ShadowMinScale       = 0.30  // floor when shrinking with altitude
    ShadowAltitudeFalloff = 64.0 // px of altitude at which scale reaches min (linear)
)

// ShadowColor: semi-transparent black, alpha pre-multiplied at draw time.
var ShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: uint8(255 * ShadowAlpha)}
```

---

## 3. Types [AC-1]

```go
// AltitudeBody is the minimum body interface the shadow needs.
// Satisfied by any beatemup actor (Body + Altitude16).
type AltitudeBody interface {
    GetPositionMin() (x, y int)   // body top-left in pixels (ground-plane)
    GetShape() body.Shape
    Altitude() int                 // ground-plane altitude in pixels (>=0)
}

// Bounds is the computed oval rectangle (in world-pixel coords) for testing.
type Bounds struct {
    CenterX, CenterY float64 // ground-plane center
    Width, Height    float64 // oval bounding box dimensions
}
```

No struct state required — shadow is a stateless render utility. The "component" is the package-level `Draw` function operating on any `AltitudeBody`.

---

## 4. Public API [AC-1, AC-2, AC-3]

```go
// ComputeBounds returns the oval bounds for body b at its current altitude.
// Centered on (X + W/2, Y + H) — i.e. the body's foot midpoint on the ground plane.
// Width/Height shrink linearly with altitude, clamped at ShadowMinScale.
func ComputeBounds(b AltitudeBody) Bounds

// ScaleFor returns the linear scale factor in [ShadowMinScale, 1.0] for altitude px.
func ScaleFor(altitude int) float64

// Draw renders a single shadow for b onto screen via camera, ONLY when altitude > 0.
// Returns true if drawn, false if skipped (altitude == 0 or b == nil).
func Draw(screen *ebiten.Image, cam *camera.Controller, b AltitudeBody) bool

// DrawAll iterates bodies and draws shadows for every Altitudable with altitude>0.
// Non-altitude bodies and zero-altitude bodies are no-ops.
// Allocates nothing when no body is airborne.
func DrawAll(screen *ebiten.Image, cam *camera.Controller, bodies []body.Collidable)
```

---

## 5. Pseudocode

```
ScaleFor(alt):
    if alt <= 0: return 1.0
    t = float(alt) / ShadowAltitudeFalloff
    if t > 1.0: t = 1.0
    return 1.0 - t*(1.0 - ShadowMinScale)

ComputeBounds(b):
    x, y := b.GetPositionMin()
    w, h := b.GetShape().Width(), b.GetShape().Height()
    s := ScaleFor(b.Altitude())
    cx := float64(x) + float64(w)/2
    cy := float64(y) + float64(h)         // foot line == ground plane Y
    bw := float64(w) * ShadowBaseWidthRatio * s
    bh := float64(ShadowBaseHeight) * s
    return Bounds{cx, cy, bw, bh}

Draw(screen, cam, b):
    if b == nil: return false
    if b.Altitude() <= 0: return false
    bn := ComputeBounds(b)
    drawOvalToCamera(screen, cam, bn, ShadowColor)
    return true

DrawAll(screen, cam, bodies):
    for each c in bodies:
        a, ok := c.(AltitudeBody)
        if !ok: continue
        Draw(screen, cam, a)
```

### Oval rendering

Use `vector.DrawFilledCircle` with non-uniform GeoM scale via a 1x1 white image, OR render an oval into a small offscreen `*ebiten.Image` and `cam.Draw` it. Implementer's choice. Required behavior:

- Final on-screen pixels are camera-translated (use `cam.Draw` with an `*ebiten.Image` + `ebiten.DrawImageOptions`, NOT `vector.DrawFilledCircle` directly to `screen`, since the world→screen offset must apply).
- Alpha respected via `op.ColorScale` or the source image's pre-multiplied alpha.

---

## 6. Beatemup Scene Integration [AC-5, AC-7]

File: `internal/kit/scenes/phases/beatemup/scene.go`

Insert a shadow pass **before** the actor draw loop in both `DrawActors` and `fullDraw`:

```
// in DrawActors(screen):
shadow.DrawAll(screen, s.camera, s.space.Bodies())
for _, b := range draworder.SortByGroundYAltitude(s.space.Bodies()) { ... }

// in fullDraw(screen), AFTER tilemap draw, BEFORE the body draw loop:
shadow.DrawAll(screen, s.camera, space.Bodies())
for _, b := range draworder.SortByGroundYAltitude(space.Bodies()) { ... }
```

Platformer scene (`internal/kit/scenes/phases/platformer/scene.go`) is NOT modified — shadows are beat-em-up-only.

---

## 7. Pre/Post Conditions

| Function | Pre | Post |
|---|---|---|
| `ScaleFor(alt)` | `alt >= 0` | returns ∈ `[ShadowMinScale, 1.0]`; `alt==0 ⇒ 1.0`; monotonic non-increasing in alt |
| `ScaleFor(alt)` | `alt >= ShadowAltitudeFalloff` | returns `ShadowMinScale` exactly |
| `ComputeBounds(b)` | body has positive W,H | `CenterX == x + W/2`; `CenterY == y + H`; `Width > 0`; `Height > 0` |
| `Draw(...)` | `b.Altitude() == 0` | returns `false`, no GPU op |
| `Draw(...)` | `b.Altitude() > 0` | returns `true`, exactly 1 draw call onto screen |
| `DrawAll([])` | empty slice | no draw calls, no allocations |

---

## 8. Red Phase Tests [AC-8]

Package: `shadow_test` (or `shadow`). Use headless `ebiten.NewImage(w, h)`.

**T-S1: ScaleFor table-driven**
```
pre:  cases = [{0,1.0}, {32,0.65}, {64,0.30}, {128,0.30}, {-5,1.0}]
act:  got := ScaleFor(alt)
post: |got - want| < 1e-6
```

**T-S2: ComputeBounds at altitude=0**
```
pre:  body at (x=100, y=200), W=20, H=32, altitude=0
act:  b := ComputeBounds(body)
post: b.CenterX == 110; b.CenterY == 232; b.Width == 20*0.75; b.Height == 4
```

**T-S3: ComputeBounds shrinks with altitude**
```
pre:  body W=20, H=32, altitude=64
act:  b := ComputeBounds(body)
post: b.Width == 20*0.75*ShadowMinScale; b.Height == 4*ShadowMinScale
      b.CenterY unchanged from altitude=0 case
```

**T-S4: Draw skipped at altitude=0**
```
pre:  body altitude=0; recordingShadowSink (test-only, see below)
act:  drew := Draw(screen, cam, body)
post: drew == false; sink.Calls == 0
```

**T-S5: Draw fires when airborne**
```
pre:  body altitude=10
act:  drew := Draw(screen, cam, body)
post: drew == true; sink.Calls == 1
```

**T-S6: DrawAll empty/no-airborne is no-op**
```
pre:  bodies = [] OR bodies = [grounded, grounded]
act:  DrawAll(screen, cam, bodies)
post: sink.Calls == 0
```

**T-S7: DrawAll counts only airborne**
```
pre:  bodies = [grounded, airborne(alt=10), airborne(alt=50), nonAltitudable]
act:  DrawAll(screen, cam, bodies)
post: sink.Calls == 2
```

**T-S8: Scene integration — shadow drawn before actor sprite**
File: `internal/kit/scenes/phases/beatemup/scene_test.go`

```
pre:  NewForTest with one airborne body; record draw order via
      SetShadowDrawHandlerForTest + SetActorDrawHandlerForTest
act:  scene.DrawActors(screen)
post: order == [shadow(b), actor(b)]
```

This requires a new test hook in `scene.go`:
```go
// SetShadowDrawerForTest overrides the shadow draw pass.
func (s *BeatemupPhaseScene) SetShadowDrawerForTest(f func(*ebiten.Image, []body.Collidable))
```

### Recording sink

To assert draw counts without a GPU, the shadow package exposes a test seam:

```go
// drawSink is the package var that performs the actual oval draw.
// Production: ovalDrawerFn = drawOval.
// Tests: replace with a recorder.
var ovalDrawerFn func(screen *ebiten.Image, cam *camera.Controller, b Bounds, c color.Color) = drawOval

// SetOvalDrawerForTest replaces the drawer; returns a restore func.
func SetOvalDrawerForTest(f func(*ebiten.Image, *camera.Controller, Bounds, color.Color)) (restore func())
```

---

## 9. Mock / Contract Inventory

- No new contracts in `internal/engine/contracts/`. `AltitudeBody` is a kit-local structural interface (satisfied by existing beat-em-up actors via `GetPositionMin`, `GetShape`, `Altitude`).
- `Mock Generator` may be SKIPPED.

---

## 10. Behavioural Edge Mapping

| Edge | Handled by |
|---|---|
| Very large altitude | `ScaleFor` clamps at `ShadowMinScale` (T-S1 last row) |
| Altitude==0 landing frame | `Draw` early-return (T-S4) |
| No airborne bodies | `DrawAll` loop no-op (T-S6) |
| Headless test | `ebiten.NewImage` + sink seam |
