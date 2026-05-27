# SPEC — 072-debug-flag-overlay

**Branch:** `072-debug-flag-overlay`
**Bounded Context:** Engine — `internal/engine/debug/`, `internal/engine/ui/debugoverlay/`, `internal/engine/app/`

---

## 1. Debug Registry [AC-1, AC-2, AC-4, AC-10]

### File: `internal/engine/debug/registry.go` (new)

```go
package debug

// Entry is a single registered debug flag (name + pointer to backing bool).
type Entry struct {
    Name string
    Ptr  *bool
}

// Register stores ptr under name (last-write-wins). Safe for concurrent use.
func Register(name string, ptr *bool)

// List returns a name-sorted snapshot of all registered entries.
// Always returns a non-nil slice (possibly empty).
func List() []Entry
```

### Internal state (extend `debug.go`)

Add alongside existing `channels`/`watchCache`:

```
registry map[string]*bool  // guarded by existing `mu`
```

### Behaviour

```
Register(name, ptr):
  mu.Lock()
  if registry == nil: registry = map[string]*bool{}
  registry[name] = ptr        // overwrite on duplicate name (no panic)
  mu.Unlock()

List():
  mu.Lock()
  names := sortedKeys(registry)
  out := make([]Entry, 0, len(names))
  for n in names: out = append(out, Entry{n, registry[n]})
  mu.Unlock()
  return out                  // non-nil even when empty

Reset():                      // extend existing function
  mu.Lock()
  channels.Store(nil)
  watchCache = make(map[string]string)
  registry = nil              // NEW: clear registry
  mu.Unlock()
  enabled.Store(false)
```

### `Init` / `InitFromReader` integration [AC-2]

After `channels.Store(&m)` in `InitFromReader`, populate registry from the
*stored* map so pointers remain valid:

```
storedMap := channels.Load()  // *map[string]bool
for k := range *storedMap:
    ptr := mapValuePointer(storedMap, k)   // &(*storedMap)[k] not allowed in Go;
                                           // see Implementation Note below
    Register(k, ptr)
```

**Implementation note:** Go does not allow `&m[k]` on a map. Two options:
- (Preferred) Replace `var m map[string]bool` with a slice of `Entry`-like
  records OR keep `m` but also build `regMap := map[string]*bool` whose values
  point into a stable per-key `*bool` (one allocation per key). The new
  `regMap` is stored alongside `channels` and is what `Enabled` consults too.
- Simpler alternative: change `channels` from `map[string]bool` to
  `map[string]*bool`. `Enabled(channel)` becomes `p := (*m)[channel]; return p != nil && *p`.

The Spec mandates option 2 (simpler, no double-store):

```
// Replace internal type:
channels atomic.Pointer[map[string]*bool]

InitFromReader(r):
  Reset()
  if r == nil: return
  var raw map[string]bool
  if decode fails: return
  m := make(map[string]*bool, len(raw))
  anyOn := false
  for k, v := range raw:
      val := v
      m[k] = &val
      if v: anyOn = true
  mu.Lock()
  channels.Store(&m)
  for k, p := range m: registry[k] = p   // (init registry if nil)
  enabled.Store(anyOn)
  mu.Unlock()

Enabled(ch):
  if !enabled.Load(): return false
  m := channels.Load(); if m == nil: return false
  p := (*m)[ch]; return p != nil && *p
```

### Tests: `internal/engine/debug/registry_test.go` (new) [AC-10]

Table-driven `TestRegisterList`:

```
T-R1: empty registry
  pre:  Reset()
  act:  got := List()
  post: got != nil && len(got) == 0

T-R2: single entry round-trip
  pre:  Reset(); var b bool
  act:  Register("a", &b); got := List()
  post: len(got)==1; got[0].Name=="a"; got[0].Ptr==&b

T-R3: multiple entries sorted alphabetically
  pre:  Reset(); var a,b,c bool
  act:  Register("c", &c); Register("a", &a); Register("b", &b); got := List()
  post: [got[i].Name for i] == ["a","b","c"]

T-R4: duplicate name last-write-wins
  pre:  Reset(); var b1, b2 bool
  act:  Register("x", &b1); Register("x", &b2); got := List()
  post: len(got)==1; got[0].Ptr==&b2

T-R5: Reset clears registry
  pre:  Reset(); var b bool; Register("x", &b)
  act:  Reset(); got := List()
  post: got != nil && len(got) == 0

T-R6: Init populates registry from JSON
  pre:  Reset()
  act:  InitFromReader(strings.NewReader(`{"a":true,"b":false}`)); got := List()
  post: len(got)==2; got[0].Name=="a"; *got[0].Ptr==true;
        got[1].Name=="b"; *got[1].Ptr==false
```

Adjust existing `debug_test.go` only if internal type change forces it (most
tests use the public API and continue to pass).

---

## 2. Game Config Registration [AC-3]

### File: `internal/game/app/config.go` (modify)

After the existing `flag.BoolVar(...)` calls, before `return cfg`:

```go
debug.Register("cam_debug", &cfg.CamDebug)
debug.Register("collision_box", &cfg.CollisionBox)
```

Add import: `"github.com/boilerplate/ebiten-template/internal/engine/debug"`.

**Ordering:** Registration occurs at `NewConfig()` time. `flag.Parse()` is
called by the caller after `NewConfig` returns; toggling later through the
overlay still writes to the same memory because Register stores the pointer,
not the value. (Edge case in story acknowledged.)

---

## 3. DebugOverlay UI [AC-5, AC-6, AC-7, AC-9, AC-11]

### File: `internal/engine/ui/debugoverlay/overlay.go` (new)

```go
package debugoverlay

import (
    "github.com/boilerplate/ebiten-template/internal/engine/debug"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/text/v2"
)

type DebugOverlay struct {
    open    bool
    cursor  int
    face    *text.GoTextFace   // optional; nil-safe in Draw
    // Injectable seams for unit tests (default to inpututil wrappers).
    keyJustPressed func(ebiten.Key) bool
}

func New() *DebugOverlay
func (o *DebugOverlay) Open()
func (o *DebugOverlay) Close()
func (o *DebugOverlay) IsOpen() bool
func (o *DebugOverlay) SetFont(f *text.GoTextFace)
func (o *DebugOverlay) Update() bool
func (o *DebugOverlay) Draw(screen *ebiten.Image)
```

### Pseudocode

```
Update() bool:
  if !o.open: return false
  entries := debug.List()
  n := len(entries)
  if n > 0:
      if keyJustPressed(KeyArrowUp):
          o.cursor = (o.cursor - 1 + n) % n
      if keyJustPressed(KeyArrowDown):
          o.cursor = (o.cursor + 1) % n
      if keyJustPressed(KeySpace) || keyJustPressed(KeyEnter):
          if entries[o.cursor].Ptr != nil:
              *entries[o.cursor].Ptr = !*entries[o.cursor].Ptr
  else:
      o.cursor = 0
  if keyJustPressed(KeyF1):
      o.open = false
      return false
  // Clamp after possible registry mutation
  if n == 0: o.cursor = 0
  else if o.cursor >= n: o.cursor = n - 1
  return true

Draw(screen):
  if !o.open: return
  draw semi-transparent panel (e.g. RGBA{0,0,0,180}) covering a centered rect
  entries := debug.List()
  for i, e := range entries:
      mark := "[ ]"; if e.Ptr != nil && *e.Ptr: mark = "[x]"
      line := mark + " " + e.Name
      if i == o.cursor: highlight row (filled bg or different color)
      if face != nil: text.Draw(...)
```

### Import constraint check [AC-9]

`internal/engine/ui/debugoverlay/` imports only:
- `github.com/boilerplate/ebiten-template/internal/engine/debug`
- `github.com/hajimehoshi/ebiten/v2` + sub-packages

No imports of `internal/game/...` or `internal/kit/...`. Verified by `go list`
in Gatekeeper stage.

### Tests: `internal/engine/ui/debugoverlay/overlay_test.go` (new) [AC-11]

Use injected `keyJustPressed` stub (not `inpututil`) so tests are deterministic.

```
T-O1: Update returns false when closed
  pre:  o := New()             // open=false
  act:  got := o.Update()
  post: got == false

T-O2: F1 closes the overlay
  pre:  o.Open(); debug.Reset(); var a bool; debug.Register("a", &a)
        stub: F1 just pressed
  act:  got := o.Update()
  post: got == false; o.IsOpen() == false

T-O3: Up wraps from index 0 to last
  pre:  Reset(); register a,b,c (sorted); o.Open(); o.cursor=0; stub: Up
  act:  o.Update()
  post: o.cursor == 2

T-O4: Down wraps from last to 0
  pre:  Reset(); register a,b,c; o.Open(); o.cursor=2; stub: Down
  act:  o.Update()
  post: o.cursor == 0

T-O5: Space toggles pointed-to bool
  pre:  Reset(); var a bool = false; debug.Register("a", &a); o.Open(); cursor=0; stub: Space
  act:  o.Update()
  post: a == true

T-O6: Enter toggles pointed-to bool
  pre:  Reset(); var a bool = true; debug.Register("a", &a); o.Open(); cursor=0; stub: Enter
  act:  o.Update()
  post: a == false

T-O7: Empty registry — Update is safe
  pre:  Reset(); o.Open(); stub: Space
  act:  got := o.Update()
  post: got == true; no panic; o.cursor == 0

T-O8: Draw with empty registry no-ops (smoke)
  pre:  Reset(); o.Open(); screen := ebiten.NewImage(64,64)
  act:  o.Draw(screen)
  post: no panic
```

Table-driven where possible (T-O3..T-O6 share fixture).

---

## 4. Engine Wiring [AC-8]

### File: `internal/engine/app/engine.go` (modify)

```go
type Game struct {
    AppContext    *AppContext
    debugOverlay  *debugoverlay.DebugOverlay   // NEW
    debugVisible  bool                          // existing — keep for DebugPhysics
    debugFontFace *text.GoTextFace
}

NewGame(ctx):
    debug.Init("assets/data/debug.json")
    return &Game{
        AppContext:   ctx,
        debugOverlay: debugoverlay.New(),
    }
```

### Update pseudocode

```
Update():
  ctx.FrameCount++

  if IsKeyJustPressed(F1):
      if debugOverlay.IsOpen(): debugOverlay.Close()
      else:                     debugOverlay.Open()

  if debugOverlay.IsOpen():
      debugOverlay.Update()           // returns true; we still skip below
      return nil                      // skip Dialogue + Scene updates

  if ctx.DialogueManager != nil: ctx.DialogueManager.Update()
  ctx.SceneManager.Update()
  return nil
```

**Note on F1:** The current `debugVisible` toggle (for `DebugPhysics`) is
replaced by overlay open/close. `DebugPhysics` is no longer drawn by F1;
either remove the call site or gate it behind a separately-bound key. For
this story: **remove** the `if g.debugVisible { g.DebugPhysics(screen) }`
block from `Draw`. `DebugPhysics` function can stay (unused) and be addressed
in a follow-up story (see NOTES.md).

### Draw pseudocode

```
Draw(screen):
  ctx.SceneManager.Draw(screen)
  if ctx.DialogueManager != nil: ctx.DialogueManager.Draw(screen)
  debugOverlay.Draw(screen)   // no-op when closed
```

### Update existing tests

`internal/engine/app/app_test.go` `TestGameUpdateAndDrawIntegration` continues
to pass — overlay starts closed, scene/dialogue update normally.

Add `TestGame_OverlayOpenSuppressesSceneUpdate`:

```
T-G1: overlay open suppresses scene + dialogue update
  pre:  game with mock SceneManager, mock DialogueManager
        directly invoke game.debugOverlay.Open() (or expose test helper)
  act:  game.Update()
  post: SceneManager.UpdateCalled == false
        DialogueManager.UpdateCalled == false
        FrameCount incremented (still happens before the guard)

T-G2: overlay closed runs scene + dialogue update (regression of existing test)
  — existing TestGameUpdateAndDrawIntegration suffices
```

Expose `OverlayForTest() *debugoverlay.DebugOverlay` on `Game` (test-only
accessor) OR add `SetDebugOverlayOpenForTest(bool)` — TDD Specialist chooses.
Prefer a small exported method on `Game` named `DebugOverlay() *debugoverlay.DebugOverlay`
since it's harmless to expose for other engine tests.

---

## 5. Mock Inventory

No new contracts. Existing mocks reused:
- `internal/engine/mocks/MockSceneManager` (has `UpdateCalled`, `DrawCalled`).
- `internal/engine/mocks/MockDialogueManager` (must have `UpdateCalled`; verify
  in TDD stage, add field if missing).

Skip Mock Generator stage if `MockDialogueManager.UpdateCalled` already exists.

---

## 6. Pre/Post-conditions Summary

| ID | Pre | Post |
|---|---|---|
| AC-1 | `Reset()` | `Register("a",&b)` then `List()[0] == {Name:"a", Ptr:&b}` |
| AC-2 | `Reset()` | `InitFromReader(`{"a":true}`)` → `len(List())==1; *List()[0].Ptr==true` |
| AC-3 | `NewConfig()` called | `debug.List()` contains `cam_debug` and `collision_box` |
| AC-4 | flags registered | `Reset()` → `len(List())==0` |
| AC-5 | — | package `debugoverlay` exists with `DebugOverlay` + `Update`/`Draw` |
| AC-6 | overlay open, cursor on entry e | Space stub → `*e.Ptr` toggled; F1 stub → `IsOpen()==false`, Update returns false |
| AC-7 | overlay open with entries | Draw renders panel + lines `[x] name`/`[ ] name`, cursor row highlighted |
| AC-8 | overlay open | `SceneManager.Update` and `DialogueManager.Update` not called |
| AC-9 | — | `go list -deps internal/engine/ui/debugoverlay` shows no `internal/game/` or `internal/kit/` |
| AC-10 | — | tests T-R1..T-R6 pass |
| AC-11 | — | tests T-O2..T-O6 pass |

---

## 7. File Manifest

**New:**
- `internal/engine/debug/registry.go`
- `internal/engine/debug/registry_test.go`
- `internal/engine/ui/debugoverlay/overlay.go`
- `internal/engine/ui/debugoverlay/overlay_test.go`

**Modified:**
- `internal/engine/debug/debug.go` — switch `channels` to `map[string]*bool`; clear registry in `Reset`; populate registry in `InitFromReader`.
- `internal/game/app/config.go` — `debug.Register("cam_debug",...)`, `debug.Register("collision_box",...)`.
- `internal/engine/app/engine.go` — add `*debugoverlay.DebugOverlay`, F1 toggles overlay, skip Scene/Dialogue update when open, draw overlay.
- `internal/engine/app/app_test.go` — add `T-G1` overlay-open-skips-update test.
- (possibly) `internal/engine/mocks/...` — only if `MockDialogueManager.UpdateCalled` doesn't exist.
