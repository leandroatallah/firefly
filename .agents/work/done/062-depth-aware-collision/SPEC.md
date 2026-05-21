# SPEC — 062-depth-aware-collision

**Bounded Context:** Physics (`internal/engine/physics/space/`)
**Target File:** `internal/engine/physics/space/space.go`
**Test File:** `internal/engine/physics/space/space_test.go`

---

## 1. Scope

Gate `HasCollision` with a depth-lane check when BOTH bodies opt in. No changes to `internal/engine/contracts/`. No changes to `kit/` or `game/`. No new mocks (opt-in is an interface assertion against an interface declared inside `space/`).

---

## 2. Design Decisions

### 2.1 Opt-in mechanism: local interface assertion
- A new local interface `DepthLaneBody` declared in `space/space.go`. Bodies that satisfy it (downstream `kit`/`game` types) opt into the lane check. Plain 2D bodies that do not implement it use the existing path. Both `a` AND `b` must satisfy it for the gate to apply.
- Chosen over (a) a tag field on `Body` (would force a contract change), (b) config-level pair check (would import scene/config into `space/`).
- Layer-safe: interface is defined in `space/`; no upward imports needed.

### 2.2 LaneWidth source: per-body method on the opt-in interface
- `DepthLaneBody.LaneHalfWidth() int` returns the per-body half-width (in pixels) along the depth axis. The pair tolerance is the MAX of the two halves. This avoids any package-level constant and lets `kit`/`game` decide tuning without `space/` knowing about scene config.
- A package-level fallback constant `DefaultLaneHalfWidth = 8` is provided ONLY for documentation / downstream use; it is NOT consulted inside `HasCollision`.
- Edge case: `LaneHalfWidth() == 0` is valid (strict-lane mode) and must not panic.

### 2.3 GroundY source: from the opt-in interface, not `GetPosition16()`
- `DepthLaneBody.GroundY() int` returns the depth coordinate in pixels (the body's Y prior to altitude offset). Downstream impl can return `y` from `GetPosition16()` second component (after `fp16.From16`), but `space/` does not assume that mapping.
- Justification: `GetPosition16()` returns the SCREEN-Y of the rendered rect; in a 2.5D model the depth value semantically belongs to the lane interface. Keeping it on `DepthLaneBody` makes the contract explicit and decouples `space/` from `fp16` semantics for this check.

---

## 3. Types and Signatures [AC-4, AC-5, AC-6]

### File: `internal/engine/physics/space/depth_lane.go` (NEW)

```go
package space

// DefaultLaneHalfWidth is the recommended default half-width (in pixels) along
// the depth axis for 2.5D bodies. It is exported for downstream consumers; it
// is NOT consulted by HasCollision.
const DefaultLaneHalfWidth = 8

// DepthLaneBody is an opt-in marker interface for bodies that participate in
// 2.5D depth-lane collision gating. Bodies that do not satisfy this interface
// use the existing 2D bbox-only collision path.
type DepthLaneBody interface {
    // GroundY returns the body's depth coordinate in pixels (Y on the ground
    // plane, ignoring altitude). Two bodies are in the same lane when
    // abs(a.GroundY - b.GroundY) <= max(a.LaneHalfWidth, b.LaneHalfWidth).
    GroundY() int

    // LaneHalfWidth returns this body's depth tolerance in pixels. A value of
    // 0 means strict same-depth matching. Must be >= 0.
    LaneHalfWidth() int
}
```

### File: `internal/engine/physics/space/space.go` (MODIFY)

`HasCollision` is rewritten to apply the depth gate AFTER the existing bbox overlap check, only when both `a` and `b` implement `DepthLaneBody`. Existing signature unchanged.

```go
func HasCollision(a, b body.Collidable) bool
```

---

## 4. Algorithm [AC-1, AC-2, AC-3]

```
HasCollision(a, b):
  if a.ID()=="" or b.ID()=="" : return false
  if a.ID() == b.ID()         : return false

  rectsA := collisionRects(a)
  rectsB := collisionRects(b)

  bboxOverlap := any(rA.Overlaps(rB) for rA in rectsA, rB in rectsB)
  if !bboxOverlap: return false

  // Depth-lane gate: applies only when BOTH bodies opt in.
  da, okA := a.(DepthLaneBody)
  db, okB := b.(DepthLaneBody)
  if !(okA && okB):
      return true                       // legacy 2D path (AC-3)

  tol := max(da.LaneHalfWidth(), db.LaneHalfWidth())   // tol >= 0
  diff := abs(da.GroundY() - db.GroundY())
  return diff <= tol                    // AC-1, AC-2
```

Notes:
- `max` and `abs` are package-private helpers (or inlined). No external deps.
- Negative `LaneHalfWidth()` MUST NOT occur; treat as undefined behavior — no clamp, no panic (post-condition documented).

---

## 5. Pre / Post-conditions

| pre                                                                   | post                                  | AC  |
|-----------------------------------------------------------------------|---------------------------------------|-----|
| `a` or `b` has empty ID                                               | `false`                               | -   |
| `a.ID() == b.ID()`                                                    | `false`                               | -   |
| bboxes do not overlap                                                 | `false`                               | -   |
| bboxes overlap, neither opts in                                       | `true`                                | AC-3 |
| bboxes overlap, only `a` opts in                                      | `true` (legacy path)                  | AC-3 |
| bboxes overlap, both opt in, `abs(GroundY diff) > max(halfWidth)`     | `false`                               | AC-1 |
| bboxes overlap, both opt in, `abs(GroundY diff) <= max(halfWidth)`    | `true`                                | AC-2 |
| both opt in, `LaneHalfWidth() == 0` on both, `GroundY` equal          | `true`                                | edge |
| both opt in, `LaneHalfWidth() == 0` on both, `GroundY` differ by 1    | `false`                               | edge |

---

## 6. Red Phase Test Inventory [AC-7]

Append to `internal/engine/physics/space/space_test.go`. Introduce a new fixture `depthLaneCollidable` that embeds `testCollidable` and adds `groundY` + `laneHalfWidth` fields plus the two interface methods. Table-driven via `TestHasCollisionDepthLane`.

### Fixture (test-only)

```go
type depthLaneCollidable struct {
    *testCollidable
    groundY       int
    laneHalfWidth int
}

func (d *depthLaneCollidable) GroundY() int       { return d.groundY }
func (d *depthLaneCollidable) LaneHalfWidth() int { return d.laneHalfWidth }
```

### T-062-1: same-lane bbox overlap collides [AC-2]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=8}
      b=depthLane{rect=(5,5,15,15), groundY=104, halfW=8}
act:  HasCollision(a, b)
post: == true     // bbox overlaps AND |100-104|=4 <= max(8,8)
```

### T-062-2: different-lane bbox overlap does NOT collide [AC-1]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=8}
      b=depthLane{rect=(5,5,15,15), groundY=120, halfW=8}
act:  HasCollision(a, b)
post: == false    // |100-120|=20 > 8
```

### T-062-3: 2D-only bodies still collide on bbox overlap [AC-3]
```
pre:  a=testCollidable{rect=(0,0,10,10)}       // does NOT implement DepthLaneBody
      b=testCollidable{rect=(5,5,15,15)}
act:  HasCollision(a, b)
post: == true     // legacy path
```

### T-062-4: same-lane no bbox overlap does NOT collide [AC-7]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=8}
      b=depthLane{rect=(50,50,60,60), groundY=100, halfW=8}
act:  HasCollision(a, b)
post: == false    // bbox short-circuits before lane gate
```

### T-062-5: LaneHalfWidth=0 with equal GroundY collides [edge]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=0}
      b=depthLane{rect=(5,5,15,15), groundY=100, halfW=0}
act:  HasCollision(a, b)
post: == true     // strict lane, exact match
```

### T-062-6: LaneHalfWidth=0 with GroundY differ by 1 does NOT collide [edge]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=0}
      b=depthLane{rect=(5,5,15,15), groundY=101, halfW=0}
act:  HasCollision(a, b)
post: == false
```

### T-062-7: mixed scene — depth-opt-in vs plain 2D collides on bbox [AC-3, edge]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=8}
      b=testCollidable{rect=(5,5,15,15)}      // not opt-in
act:  HasCollision(a, b)
post: == true     // only one side opts in → legacy path
```

### T-062-8: asymmetric halfWidth uses MAX [AC-1, AC-2]
```
case A: a.halfW=2, b.halfW=10, |groundY diff|=8 → true   (max=10, 8<=10)
case B: a.halfW=2, b.halfW=10, |groundY diff|=11 → false (11>10)
```

### T-062-9: airborne same-lane — bbox check governs [edge]
```
pre:  a=depthLane{rect=(0,0,10,10), groundY=100, halfW=8}    // on ground
      b=depthLane{rect=(5,5,15,15), groundY=100, halfW=8}    // bbox already
                                                              // reflects screen-Y
                                                              // (Y - Altitude)
act:  HasCollision(a, b)
post: == true     // bbox overlaps AND same lane → collide
```

Note for T-062-9: `space/` does NOT know about altitude. Callers are expected to set `rect` to screen-Y (post-altitude) and `GroundY()` to depth-Y. This test documents that assumption without asserting altitude semantics.

### Table layout
All 9 cases collapse into one table-driven test `TestHasCollisionDepthLane` with fields `{name, aRect, aGroundY, aHalfW, aOptIn, bRect, bGroundY, bHalfW, bOptIn, want}`. `optIn=false` constructs a plain `testCollidable`; `optIn=true` constructs a `depthLaneCollidable`.

---

## 7. Layering Check [AC-6]

Files touched:
- `internal/engine/physics/space/space.go` — modify `HasCollision`.
- `internal/engine/physics/space/depth_lane.go` — NEW; `package space`. Imports: NONE (or stdlib only).
- `internal/engine/physics/space/space_test.go` — extend with fixture and table-driven test.

Forbidden imports in `space/` (verified by spec): no `internal/kit/`, no `internal/game/`. `DepthLaneBody` is defined inside `space/`; downstream packages implement it. Dependency direction is preserved (engine ← kit ← game).

---

## 8. AC → Test Map

| AC   | Test(s)                              |
|------|--------------------------------------|
| AC-1 | T-062-2, T-062-8B                    |
| AC-2 | T-062-1, T-062-8A                    |
| AC-3 | T-062-3, T-062-7                     |
| AC-4 | (spec) `LaneHalfWidth()` method, `DefaultLaneHalfWidth` const — no magic numbers in `HasCollision` |
| AC-5 | (spec) Section 2.1: opt-in via local interface `DepthLaneBody` |
| AC-6 | (spec+test) Section 7; verified by `go build ./internal/engine/physics/space/...` and import audit |
| AC-7 | T-062-1, T-062-2, T-062-3, T-062-4 (the four required cases) + T-062-5..9 (edge coverage) |
