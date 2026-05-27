package debug_test

import (
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/debug"
)

func TestRegisterList(t *testing.T) {
	t.Run("T-R1 empty registry returns non-nil empty slice", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		got := debug.List()
		if got == nil {
			t.Fatalf("List() returned nil, want non-nil empty slice")
		}
		if len(got) != 0 {
			t.Fatalf("List() len = %d, want 0", len(got))
		}
	})

	t.Run("T-R2 single entry round-trip", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var b bool
		debug.Register("a", &b)
		got := debug.List()
		if len(got) != 1 {
			t.Fatalf("len(List()) = %d, want 1", len(got))
		}
		if got[0].Name != "a" {
			t.Fatalf("got[0].Name = %q, want %q", got[0].Name, "a")
		}
		if got[0].Ptr != &b {
			t.Fatalf("got[0].Ptr = %p, want %p", got[0].Ptr, &b)
		}
	})

	t.Run("T-R3 multiple entries sorted alphabetically", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var a, b, c bool
		debug.Register("c", &c)
		debug.Register("a", &a)
		debug.Register("b", &b)

		got := debug.List()
		if len(got) != 3 {
			t.Fatalf("len(List()) = %d, want 3", len(got))
		}
		want := []string{"a", "b", "c"}
		for i, w := range want {
			if got[i].Name != w {
				t.Fatalf("got[%d].Name = %q, want %q", i, got[i].Name, w)
			}
		}
	})

	t.Run("T-R4 duplicate name last-write-wins", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var b1, b2 bool
		debug.Register("x", &b1)
		debug.Register("x", &b2)

		got := debug.List()
		if len(got) != 1 {
			t.Fatalf("len(List()) = %d, want 1", len(got))
		}
		if got[0].Ptr != &b2 {
			t.Fatalf("got[0].Ptr = %p, want %p (last-write-wins)", got[0].Ptr, &b2)
		}
	})

	t.Run("T-R5 Reset clears registry", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var b bool
		debug.Register("x", &b)

		debug.Reset()
		got := debug.List()
		if got == nil {
			t.Fatalf("List() returned nil after Reset, want non-nil empty slice")
		}
		if len(got) != 0 {
			t.Fatalf("len(List()) = %d after Reset, want 0", len(got))
		}
	})

	t.Run("T-R6 InitFromReader populates registry from JSON", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		debug.InitFromReader(strings.NewReader(`{"a":true,"b":false}`))

		got := debug.List()
		if len(got) != 2 {
			t.Fatalf("len(List()) = %d, want 2", len(got))
		}
		if got[0].Name != "a" {
			t.Fatalf("got[0].Name = %q, want %q", got[0].Name, "a")
		}
		if got[0].Ptr == nil || *got[0].Ptr != true {
			t.Fatalf("got[0].Ptr deref = %v, want true", got[0].Ptr)
		}
		if got[1].Name != "b" {
			t.Fatalf("got[1].Name = %q, want %q", got[1].Name, "b")
		}
		if got[1].Ptr == nil || *got[1].Ptr != false {
			t.Fatalf("got[1].Ptr deref = %v, want false", got[1].Ptr)
		}
	})
}
