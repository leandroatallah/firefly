package phaseoverlay

import (
	"testing"

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

func threeEntries() []Entry {
	return []Entry{
		{ID: 1, Name: "Intro"},
		{ID: 7, Name: "Area 1 - Stage 1"},
		{ID: 10, Name: "Area 1 - Stage 4"},
	}
}

func TestPhaseOverlay(t *testing.T) {
	t.Run("T-P1 Update returns false when closed", func(t *testing.T) {
		o := New()
		o.keyJustPressed = stubKeys()

		if got := o.Update(); got != false {
			t.Fatalf("Update() = %v, want false when closed", got)
		}
	})

	t.Run("T-P2 F2 closes the overlay", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.keyJustPressed = stubKeys(ebiten.KeyF2)

		if got := o.Update(); got != false {
			t.Fatalf("Update() = %v, want false after F2 closes overlay", got)
		}
		if o.IsOpen() {
			t.Fatalf("IsOpen() = true after F2, want false")
		}
	})

	t.Run("T-P3 Escape closes the overlay", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.keyJustPressed = stubKeys(ebiten.KeyEscape)

		if got := o.Update(); got != false {
			t.Fatalf("Update() = %v, want false after Escape", got)
		}
		if o.IsOpen() {
			t.Fatalf("IsOpen() = true after Escape, want false")
		}
	})

	t.Run("T-P4 Up wraps from index 0 to last", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.cursor = 0
		o.keyJustPressed = stubKeys(ebiten.KeyArrowUp)

		o.Update()
		if o.cursor != 2 {
			t.Fatalf("cursor = %d, want 2 (Up wraps from 0 to last)", o.cursor)
		}
	})

	t.Run("T-P5 Down wraps from last to 0", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.cursor = 2
		o.keyJustPressed = stubKeys(ebiten.KeyArrowDown)

		o.Update()
		if o.cursor != 0 {
			t.Fatalf("cursor = %d, want 0 (Down wraps from last to 0)", o.cursor)
		}
	})

	t.Run("T-P6 Enter selects pointed-to phase ID and closes", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.cursor = 1

		var got int
		var called bool
		o.SetOnSelect(func(id int) {
			got = id
			called = true
		})
		o.keyJustPressed = stubKeys(ebiten.KeyEnter)

		if alive := o.Update(); alive != false {
			t.Fatalf("Update() = %v, want false after Enter selects and closes", alive)
		}
		if !called {
			t.Fatalf("onSelect was not called")
		}
		if got != 7 {
			t.Fatalf("selected ID = %d, want 7 (entries[1].ID)", got)
		}
		if o.IsOpen() {
			t.Fatalf("IsOpen() = true after Enter, want false")
		}
	})

	t.Run("T-P7 Enter with nil onSelect is safe", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.cursor = 0
		o.keyJustPressed = stubKeys(ebiten.KeyEnter)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Update() panicked with nil onSelect: %v", r)
			}
		}()

		if alive := o.Update(); alive != false {
			t.Fatalf("Update() = %v, want false (overlay closes on Enter)", alive)
		}
	})

	t.Run("T-P8 empty entries Update is safe", func(t *testing.T) {
		o := New()
		o.Open()
		o.keyJustPressed = stubKeys(ebiten.KeyEnter)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Update() panicked with no entries: %v", r)
			}
		}()

		if alive := o.Update(); alive != true {
			t.Fatalf("Update() = %v, want true (overlay open, no entries)", alive)
		}
		if o.cursor != 0 {
			t.Fatalf("cursor = %d, want 0 with no entries", o.cursor)
		}
	})

	t.Run("T-P9 cursor clamps when entries shrink", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.cursor = 2
		o.SetEntries([]Entry{{ID: 1, Name: "Only"}})
		o.keyJustPressed = stubKeys()

		o.Update()
		if o.cursor != 0 {
			t.Fatalf("cursor = %d, want 0 (clamped to last valid index)", o.cursor)
		}
	})

	t.Run("T-P10 Draw with no entries no-ops smoke", func(t *testing.T) {
		o := New()
		o.Open()

		screen := ebiten.NewImage(64, 64)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Draw() panicked with no entries: %v", r)
			}
		}()
		o.Draw(screen)
	})

	t.Run("T-P11 digit hotkey selects by index and closes", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()

		var got int
		o.SetOnSelect(func(id int) { got = id })
		o.keyJustPressed = stubKeys(ebiten.Key1) // index 1 → entries[1].ID = 7

		if alive := o.Update(); alive != false {
			t.Fatalf("Update() = %v, want false after hotkey selection", alive)
		}
		if got != 7 {
			t.Fatalf("selected ID = %d, want 7", got)
		}
		if o.IsOpen() {
			t.Fatalf("IsOpen() = true after hotkey selection, want false")
		}
	})

	t.Run("T-P12 digit hotkey out of range is ignored", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries()) // 3 entries: valid hotkeys 0,1,2
		o.Open()

		var called bool
		o.SetOnSelect(func(_ int) { called = true })
		o.keyJustPressed = stubKeys(ebiten.Key5) // index 5 out of range

		if alive := o.Update(); alive != true {
			t.Fatalf("Update() = %v, want true (overlay stays open)", alive)
		}
		if called {
			t.Fatalf("onSelect called unexpectedly for out-of-range hotkey")
		}
	})

	t.Run("T-P13 digit hotkey 0 selects first entry", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()

		var got int
		o.SetOnSelect(func(id int) { got = id })
		o.keyJustPressed = stubKeys(ebiten.Key0) // index 0 → entries[0].ID = 1

		o.Update()
		if got != 1 {
			t.Fatalf("selected ID = %d, want 1", got)
		}
	})

	t.Run("T-P14 digit hotkey with nil onSelect is safe", func(t *testing.T) {
		o := New()
		o.SetEntries(threeEntries())
		o.Open()
		o.keyJustPressed = stubKeys(ebiten.Key0)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Update() panicked with nil onSelect on digit key: %v", r)
			}
		}()
		o.Update()
	})
}
