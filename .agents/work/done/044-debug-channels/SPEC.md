# SPEC — 044-debug-channels

## Branch

`044-debug-channels`

## Summary

Introduce a new bounded leaf package `internal/engine/debug/` providing a named-channel observability layer with two primitives: `Log` (every call) and `Watch` (change-only). The package is silently disabled by default and is enabled per-channel via `assets/data/debug.json`. The hot path on disabled channels is a single boolean guard — no string formatting, no map lookup, no allocation.

## Package Layout

```
internal/engine/debug/
├── debug.go           // public API: Init, InitFromReader, Log, Watch, Enabled, Reset
├── debug_test.go      // table-driven tests for all ACs
└── doc.go             // package overview (optional)
```

No interface is added under `internal/engine/contracts/`. Rationale: `debug` is a leaf utility package with package-level functions (similar to `log` or `data/config`). It is not injected; it is not mocked; it has no collaborators. Adding a contract would invert the design without benefit.

## Public API

```go
package debug

// Init loads the JSON config at path. Missing file silently disables all
// channels. Any I/O or parse error is treated as "no config" and is non-fatal.
// Init is idempotent and may be called more than once (later calls replace
// state); concurrent calls are safe.
func Init(path string)

// InitFromReader is the test-friendly seam used by Init internally and by
// unit tests to feed JSON without touching the filesystem. A nil reader
// disables all channels.
func InitFromReader(r io.Reader)

// Log writes a formatted line to stdout on every call when channel is
// enabled. No-op when channel is disabled or package is uninitialized.
func Log(channel, format string, args ...any)

// Watch writes a formatted line to stdout only when value differs from the
// previous call with the same (channel, key) pair. First call always logs.
func Watch(channel, key string, value any)

// Enabled reports whether channel is currently enabled. Exposed for tests
// and for callers wishing to guard expensive argument construction.
func Enabled(channel string) bool

// Reset clears all internal state (channels, watch cache). Test-only helper;
// safe to call from production code but normally unnecessary.
func Reset()
```

## Internal State

```go
//nolint:gochecknoglobals — package-level state mirrors data/config pattern
var (
    enabled   atomic.Bool                    // master "any channels active" flag
    mu        sync.RWMutex                   // guards channels and watchCache
    channels  map[string]bool                // channel name → enabled
    watchCache map[string]string             // "channel/key" → fmt.Sprint(lastValue)
)
```

- `enabled` is an `atomic.Bool` set to `true` only when at least one channel parsed from config is `true`. It is the single-bool fast-path guard required by AC-7.
- `channels` is the per-channel toggle map.
- `watchCache` keys are concatenated `channel + "/" + key` per AC-5 (per-channel scoping).
- Values stored as `string` via `fmt.Sprint(value)` so non-comparable types (slices, maps, structs containing slices) can be tracked without panics on `==`.

## JSON Schema

`assets/data/debug.json` is a flat object mapping channel name to bool:

```json
{
  "player_state": true,
  "physics": false,
  "enemy_state": true
}
```

- Unknown channels (those used by `Log`/`Watch` but absent from JSON) are treated as disabled.
- Unmarshal target: `map[string]bool`. Any value that is not a bool causes the file to be rejected (silent disable) — keeps the schema strict.
- Missing file, empty file, malformed JSON → all channels disabled, no error surfaced.

## Algorithms

### Init / InitFromReader

```
Init(path):
    f, err := os.Open(path)
    if err != nil:
        Reset()
        return
    defer f.Close()
    InitFromReader(f)

InitFromReader(r):
    Reset()
    if r == nil: return
    var m map[string]bool
    if json.NewDecoder(r).Decode(&m) != nil: return
    mu.Lock()
    channels = m
    anyOn := false
    for _, v := range m { if v { anyOn = true; break } }
    mu.Unlock()
    enabled.Store(anyOn)
```

### Log (fast path)

```
Log(channel, format, args...):
    if !enabled.Load(): return            // single bool guard, no map access
    mu.RLock()
    on := channels[channel]
    mu.RUnlock()
    if !on: return
    fmt.Printf(format+"\n", args...)
```

### Watch (change-detection)

```
Watch(channel, key, value):
    if !enabled.Load(): return            // single bool guard
    mu.RLock()
    on := channels[channel]
    mu.RUnlock()
    if !on: return
    cacheKey := channel + "/" + key
    newStr := fmt.Sprint(value)
    mu.Lock()
    prev, seen := watchCache[cacheKey]
    if seen && prev == newStr {
        mu.Unlock()
        return
    }
    watchCache[cacheKey] = newStr
    mu.Unlock()
    fmt.Printf("[%s] %s=%s\n", channel, key, newStr)
```

The output format `[<channel>] <key>=<value>` for `Watch` is chosen for grep-friendly per-channel filtering. `Log` passes through user-supplied format verbatim with a trailing newline.

### Reset

```
Reset():
    mu.Lock()
    channels = nil
    watchCache = make(map[string]string)
    mu.Unlock()
    enabled.Store(false)
```

## Engine Wiring

`internal/engine/app/engine.go` `NewGame(ctx *AppContext) *Game`:

```go
import "github.com/boilerplate/ebiten-template/internal/engine/debug"

func NewGame(ctx *AppContext) *Game {
    debug.Init("assets/data/debug.json")     // <-- new line, before any further setup
    return &Game{AppContext: ctx}
}
```

Rationale for placing `debug.Init` inside `NewGame` rather than `internal/game/app/setup.go`:
- AC-6 explicitly requires the engine entry point.
- The engine package owns the lifecycle; the game-side bootstrap may be replaced per project but the engine wiring should always activate debug.

A missing `assets/data/debug.json` is the default state for production / fresh checkouts and produces no error and no output, satisfying AC-1.

## Pre-conditions

- Go 1.25+ stdlib (`encoding/json`, `os`, `sync`, `sync/atomic`, `fmt`, `io`).
- `assets/` directory exists and is the current working directory's relative root (existing convention; `data/config` and i18n use the same path style).

## Post-conditions

- After `Init` returns, `Enabled(name)` reflects the JSON contents.
- After a successful `Init` followed by `Reset`, all channels are off and `watchCache` is empty.
- `Log` and `Watch` are safe to call from multiple goroutines (RWMutex on reads, full lock on writes).

## Integration Points (Bounded Context)

| Caller | Calls | Story |
|---|---|---|
| `internal/engine/app/engine.go::NewGame` | `debug.Init("assets/data/debug.json")` | this story |
| Future Actor / State code | `debug.Watch(...)`, `debug.Log(...)` | out of scope (per "Out of Scope" in USER_STORY) |

No call sites are added to `internal/game/` in this story.

## Acceptance Criteria Traceability

| AC | Spec Section |
|---|---|
| AC-1 No-op when uninitialized / missing file | `Init` algorithm: `os.Open` error → `Reset()` → `enabled=false`; `Log`/`Watch` fast-path returns on `!enabled.Load()`. |
| AC-2 JSON enables channels | `InitFromReader` decode + `channels[channel]` lookup in `Log`. |
| AC-3 `Log` prints every call | `Log` algorithm has no de-dup; prints unconditionally when channel on. |
| AC-4 `Watch` only on change | `watchCache[cacheKey] == newStr` guard; first call has `!seen` → logs and stores. |
| AC-5 `Watch` keys per-channel | `cacheKey = channel + "/" + key`; same `key` under different `channel` → distinct entries. |
| AC-6 Engine wires `Init` | `NewGame` change in `internal/engine/app/engine.go`. |
| AC-7 No overhead on disabled | `enabled.Load()` is a single atomic-bool read; if false, function returns before any map lookup, allocation, or formatting. |

## Red Phase (Failing Test Scenario)

`internal/engine/debug/debug_test.go` (new file). Tests capture stdout by swapping `os.Stdout` for an `os.Pipe` and reading the pipe in the assertion phase, then restoring. Each test calls `debug.Reset()` in setup.

Table-driven test list:

1. `TestLog_NoopWhenUninitialized` — without calling `Init`, `Log("any","x")` produces empty output. (AC-1)
2. `TestLog_NoopOnMissingFile` — `Init("/nonexistent.json")` then `Log("foo","x")` produces empty output. (AC-1)
3. `TestLog_PrintsWhenChannelEnabled` — `InitFromReader` with `{"player_state":true}`; three calls to `Log("player_state","pos=%d",i)` produce three lines containing `pos=0`, `pos=1`, `pos=2`. (AC-2, AC-3)
4. `TestLog_SilentWhenChannelDisabled` — config `{"physics":false}`; `Log("physics","x")` produces empty output. (AC-2)
5. `TestWatch_LogsOnceForRepeatedValue` — channel on; ten calls of `Watch("player_state","state","idle")` produce exactly one line. (AC-4)
6. `TestWatch_LogsAgainOnValueChange` — after the ten idle calls, one call with `"walk"` produces a second line containing `state=walk`. (AC-4)
7. `TestWatch_PerChannelScoping` — two channels enabled; `Watch("player_state","state","idle")` and `Watch("enemy_state","state","idle")` both log (two lines total). (AC-5)
8. `TestEnabled_ReflectsConfig` — after `InitFromReader({"a":true,"b":false})`, `Enabled("a")==true`, `Enabled("b")==false`, `Enabled("c")==false`.
9. `TestInit_MalformedJSON_DisablesAll` — feed `"{not json"`; `Enabled("anything")==false`; `Log` produces nothing. (AC-1)
10. `TestNoOverhead_DisabledFastPath` — sanity: with `enabled==false`, calling `Log` 10_000 times completes without allocating channels map (asserted via `testing.AllocsPerRun` ≤ small bound, e.g. 0 allocs for the disabled path). (AC-7)

These tests must initially fail because `internal/engine/debug/` does not exist.

An additional integration test in `internal/engine/app/app_test.go` will assert that `NewGame(ctx)` does not panic when `assets/data/debug.json` is absent, satisfying the wiring contract of AC-6. (Existing test already calls `NewGame`; the new assertion is implicit — no panic — and is satisfied by the silent-no-op `Init`.)

## Non-Goals

- No call sites in `internal/game/` (per USER_STORY "Out of Scope").
- No structured logging, log levels, or external sinks.
- No build tags or compile-time elision.
- No CLI flags or env-var activation.

## Pipeline Next Steps

No new contracts were introduced (the package is a leaf utility with package-level functions). Therefore the **Mock Generator is skipped**.

Recommended order:

1. **TDD Specialist** — author `internal/engine/debug/debug_test.go` covering ACs 1–7 (Red Phase scenarios above).
2. **Feature Implementer** — implement `internal/engine/debug/debug.go` plus the one-line wiring change in `internal/engine/app/engine.go` (Green Phase).
3. **Workflow Gatekeeper** — verify coverage delta and AC traceability.
