package debug_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/debug"
)

// captureStdout swaps os.Stdout for a pipe, runs fn, and returns whatever was
// written to stdout while fn was executing.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	old := os.Stdout
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

	return <-done
}

func TestLog_NoopWhenUninitialized(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.Log("any", "x")
	})

	if got != "" {
		t.Fatalf("expected no output when uninitialized, got %q", got)
	}
}

func TestLog_NoopOnMissingFile(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.Init("/nonexistent-debug-config-xyz.json")
		debug.Log("foo", "x")
	})

	if got != "" {
		t.Fatalf("expected no output when config file missing, got %q", got)
	}
}

func TestLog_PrintsWhenChannelEnabled(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.InitFromReader(strings.NewReader(`{"player_state":true}`))
		for i := 0; i < 3; i++ {
			debug.Log("player_state", "pos=%d", i)
		}
	})

	cases := []struct {
		name   string
		needle string
	}{
		{"line 0", "pos=0"},
		{"line 1", "pos=1"},
		{"line 2", "pos=2"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(got, tc.needle) {
				t.Fatalf("expected output to contain %q, got %q", tc.needle, got)
			}
		})
	}

	if n := strings.Count(got, "\n"); n != 3 {
		t.Fatalf("expected 3 lines, got %d (output=%q)", n, got)
	}
}

func TestLog_SilentWhenChannelDisabled(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.InitFromReader(strings.NewReader(`{"physics":false}`))
		debug.Log("physics", "x")
	})

	if got != "" {
		t.Fatalf("expected no output when channel disabled, got %q", got)
	}
}

func TestWatch_LogsOnceForRepeatedValue(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.InitFromReader(strings.NewReader(`{"player_state":true}`))
		for i := 0; i < 10; i++ {
			debug.Watch("player_state", "state", "idle")
		}
	})

	if n := strings.Count(got, "\n"); n != 1 {
		t.Fatalf("expected exactly 1 line for repeated identical value, got %d (output=%q)", n, got)
	}
	if !strings.Contains(got, "idle") {
		t.Fatalf("expected output to contain %q, got %q", "idle", got)
	}
}

func TestWatch_LogsAgainOnValueChange(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.InitFromReader(strings.NewReader(`{"player_state":true}`))
		for i := 0; i < 10; i++ {
			debug.Watch("player_state", "state", "idle")
		}
		debug.Watch("player_state", "state", "walk")
	})

	if n := strings.Count(got, "\n"); n != 2 {
		t.Fatalf("expected 2 lines (one per distinct value), got %d (output=%q)", n, got)
	}
	if !strings.Contains(got, "state=walk") {
		t.Fatalf("expected output to contain %q, got %q", "state=walk", got)
	}
}

func TestWatch_PerChannelScoping(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.InitFromReader(strings.NewReader(`{"player_state":true,"physics":true}`))
		debug.Watch("player_state", "state", "idle")
		debug.Watch("physics", "state", "idle")
	})

	if n := strings.Count(got, "\n"); n != 2 {
		t.Fatalf("expected 2 lines (one per channel), got %d (output=%q)", n, got)
	}
}

func TestEnabled_ReflectsConfig(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	debug.InitFromReader(strings.NewReader(`{"a":true,"b":false}`))

	cases := []struct {
		name    string
		channel string
		want    bool
	}{
		{"explicitly true", "a", true},
		{"explicitly false", "b", false},
		{"absent channel", "c", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := debug.Enabled(tc.channel); got != tc.want {
				t.Fatalf("Enabled(%q) = %v, want %v", tc.channel, got, tc.want)
			}
		})
	}
}

func TestInit_MalformedJSON_DisablesAll(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	got := captureStdout(t, func() {
		debug.InitFromReader(strings.NewReader(`{not json`))
		debug.Log("anything", "x")
	})

	if got != "" {
		t.Fatalf("expected no output after malformed JSON, got %q", got)
	}
	if debug.Enabled("anything") {
		t.Fatalf("Enabled(\"anything\") = true after malformed JSON, want false")
	}
}

func TestNoOverhead_DisabledFastPath(t *testing.T) {
	debug.Reset()
	t.Cleanup(debug.Reset)

	// Master flag is false (no Init call). Fast path must not allocate.
	allocs := testing.AllocsPerRun(100, func() {
		for i := 0; i < 10_000; i++ {
			debug.Log("any", "msg=%d", 0)
		}
	})

	if allocs > 0 {
		t.Fatalf("expected 0 allocs on disabled fast path, got %v", allocs)
	}
}
