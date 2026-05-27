package debug

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

//nolint:gochecknoglobals
var (
	enabled    atomic.Bool
	mu         sync.Mutex // guards watchCache, registry, and initialization
	channels   atomic.Pointer[map[string]*bool]
	watchCache map[string]string
	registry   map[string]*bool
)

func init() {
	Reset()
}

// Init loads the JSON config at path. Missing file silently disables all
// channels. Any I/O or parse error is treated as "no config" and is non-fatal.
// Init is idempotent and may be called more than once (later calls replace
// state); concurrent calls are safe.
func Init(path string) {
	f, err := os.Open(path)
	if err != nil {
		clearChannels()
		return
	}
	defer f.Close()
	InitFromReader(f)
}

// clearChannels resets channels + watchCache + enabled flag, preserving the
// registry so externally-registered pointers (e.g. CLI debug flags) survive
// across Init calls.
func clearChannels() {
	mu.Lock()
	channels.Store(nil)
	watchCache = make(map[string]string)
	mu.Unlock()
	enabled.Store(false)
}

// InitFromReader is the test-friendly seam used by Init internally and by
// unit tests to feed JSON without touching the filesystem. A nil reader
// disables all channels.
func InitFromReader(r io.Reader) {
	clearChannels()
	if r == nil {
		return
	}

	var raw map[string]bool
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return
	}

	m := make(map[string]*bool, len(raw))
	anyOn := false
	for k, v := range raw {
		val := v
		m[k] = &val
		if v {
			anyOn = true
		}
	}

	mu.Lock()
	channels.Store(&m)
	mu.Unlock()
	enabled.Store(anyOn)
}

// Log writes a formatted line to stdout on every call when channel is
// enabled. No-op when channel is disabled or package is uninitialized.
func Log(channel, format string, args ...any) {
	if Enabled(channel) {
		doPrintf(format, args)
	}
}

//go:noinline
func doPrintf(format string, args []any) {
	fmt.Printf(format+"\n", args...)
}

// Watch writes a formatted line to stdout only when value differs from the
// previous call with the same (channel, key) pair. First call always logs.
func Watch(channel, key string, value any) {
	if Enabled(channel) {
		watchSlow(channel, key, value)
	}
}

func watchSlow(channel, key string, value any) {
	cacheKey := channel + "/" + key
	newStr := fmt.Sprint(value)

	mu.Lock()
	if watchCache == nil {
		watchCache = make(map[string]string)
	}
	prev, seen := watchCache[cacheKey]
	if seen && prev == newStr {
		mu.Unlock()
		return
	}
	watchCache[cacheKey] = newStr
	mu.Unlock()

	fmt.Printf("[%s] %s=%s\n", channel, key, newStr)
}

// Enabled reports whether channel is currently enabled. Exposed for tests
// and for callers wishing to guard expensive argument construction.
func Enabled(channel string) bool {
	m := channels.Load()
	if m == nil {
		return false
	}
	p := (*m)[channel]
	return p != nil && *p
}

// Reset clears all internal state (channels, watchCache, registry). Test-only
// helper; safe to call from production code but normally unnecessary.
func Reset() {
	mu.Lock()
	channels.Store(nil)
	watchCache = make(map[string]string)
	registry = nil
	mu.Unlock()
	enabled.Store(false)
}
