package debugoverlay

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/debug"
	"github.com/hajimehoshi/ebiten/v2"
)

// stubKeys returns a keyJustPressed func that reports true only for the keys
// passed in (one-shot). It allows simulating a single frame of input.
func stubKeys(keys ...ebiten.Key) func(ebiten.Key) bool {
	set := make(map[ebiten.Key]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	return func(k ebiten.Key) bool {
		return set[k]
	}
}

func TestDebugOverlay(t *testing.T) {
	t.Run("T-O1 Update returns false when closed", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		o := New()
		o.keyJustPressed = stubKeys()

		if got := o.Update(); got != false {
			t.Fatalf("Update() = %v, want false when closed", got)
		}
	})

	t.Run("T-O2 F1 closes the overlay", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var a bool
		debug.Register("a", &a)

		o := New()
		o.Open()
		o.keyJustPressed = stubKeys(ebiten.KeyF1)

		got := o.Update()
		if got != false {
			t.Fatalf("Update() = %v, want false after F1 closes overlay", got)
		}
		if o.IsOpen() {
			t.Fatalf("IsOpen() = true after F1, want false")
		}
	})

	t.Run("T-O3 Up wraps from index 0 to last", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var a, b, c bool
		debug.Register("a", &a)
		debug.Register("b", &b)
		debug.Register("c", &c)

		o := New()
		o.Open()
		o.cursor = 0
		o.keyJustPressed = stubKeys(ebiten.KeyArrowUp)

		o.Update()
		if o.cursor != 2 {
			t.Fatalf("cursor = %d, want 2 (Up wraps from 0 to last)", o.cursor)
		}
	})

	t.Run("T-O4 Down wraps from last to 0", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var a, b, c bool
		debug.Register("a", &a)
		debug.Register("b", &b)
		debug.Register("c", &c)

		o := New()
		o.Open()
		o.cursor = 2
		o.keyJustPressed = stubKeys(ebiten.KeyArrowDown)

		o.Update()
		if o.cursor != 0 {
			t.Fatalf("cursor = %d, want 0 (Down wraps from last to 0)", o.cursor)
		}
	})

	t.Run("T-O5 Space toggles pointed-to bool false to true", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		var a bool
		debug.Register("a", &a)

		o := New()
		o.Open()
		o.cursor = 0
		o.keyJustPressed = stubKeys(ebiten.KeySpace)

		o.Update()
		if a != true {
			t.Fatalf("a = %v, want true after Space toggle", a)
		}
	})

	t.Run("T-O6 Enter toggles pointed-to bool true to false", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		a := true
		debug.Register("a", &a)

		o := New()
		o.Open()
		o.cursor = 0
		o.keyJustPressed = stubKeys(ebiten.KeyEnter)

		o.Update()
		if a != false {
			t.Fatalf("a = %v, want false after Enter toggle", a)
		}
	})

	t.Run("T-O7 empty registry Update is safe", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		o := New()
		o.Open()
		o.keyJustPressed = stubKeys(ebiten.KeySpace)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Update() panicked with empty registry: %v", r)
			}
		}()

		got := o.Update()
		if got != true {
			t.Fatalf("Update() = %v, want true (overlay still open, no entries)", got)
		}
		if o.cursor != 0 {
			t.Fatalf("cursor = %d, want 0 with empty registry", o.cursor)
		}
	})

	t.Run("T-O8 Draw with empty registry no-ops smoke", func(t *testing.T) {
		debug.Reset()
		t.Cleanup(debug.Reset)

		o := New()
		o.Open()

		screen := ebiten.NewImage(64, 64)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Draw() panicked with empty registry: %v", r)
			}
		}()
		o.Draw(screen)
	})
}
