package phaseoverlay

import (
	"image/color"
	"strconv"

	"github.com/boilerplate/ebiten-template/internal/engine/ui/overlayutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Entry is a single selectable phase in the overlay.
type Entry struct {
	ID   int
	Name string
}

// PhaseOverlay is an in-game overlay that lists registered phases and jumps to
// the selected one when the player confirms with Enter. It is purely a UI
// component: the actual phase switch is delegated to the onSelect callback so
// the engine stays decoupled from game-specific phase wiring.
type PhaseOverlay struct {
	overlayutil.Base
	cursor         int
	entries        []Entry
	face           *text.GoTextFace
	keyJustPressed func(ebiten.Key) bool
	onSelect       func(id int)
}

// New creates a PhaseOverlay wired to the real Ebitengine key-input backend.
func New() *PhaseOverlay {
	return &PhaseOverlay{
		keyJustPressed: func(k ebiten.Key) bool {
			return inpututil.IsKeyJustPressed(k)
		},
	}
}

// Open makes the overlay visible.
func (o *PhaseOverlay) Open() { o.Base.Open() }

// Close hides the overlay.
func (o *PhaseOverlay) Close() { o.Base.Close() }

// IsOpen reports whether the overlay is currently visible.
func (o *PhaseOverlay) IsOpen() bool { return o.Base.IsOpen() }

// SetFont sets the font face used when drawing entry labels.
func (o *PhaseOverlay) SetFont(f *text.GoTextFace) { o.face = f }

// SetEntries sets the list of phases shown in the overlay.
func (o *PhaseOverlay) SetEntries(entries []Entry) { o.entries = entries }

// SetOnSelect registers the callback invoked with the selected phase ID when
// the player confirms a choice with Enter.
func (o *PhaseOverlay) SetOnSelect(fn func(id int)) { o.onSelect = fn }

// Update handles input for the overlay. It returns true when the overlay
// consumed the frame (i.e. the overlay is still open after processing), and
// false when the overlay is closed or was closed this frame.
func (o *PhaseOverlay) Update() bool {
	if !o.IsOpen() {
		return false
	}

	if o.keyJustPressed(ebiten.KeyF2) || o.keyJustPressed(ebiten.KeyEscape) {
		o.Close()
		return false
	}

	n := len(o.entries)
	if n > 0 {
		if o.keyJustPressed(ebiten.KeyArrowUp) {
			o.cursor = (o.cursor - 1 + n) % n
		}
		if o.keyJustPressed(ebiten.KeyArrowDown) {
			o.cursor = (o.cursor + 1) % n
		}
		if o.keyJustPressed(ebiten.KeyEnter) {
			if o.cursor >= 0 && o.cursor < n && o.onSelect != nil {
				o.onSelect(o.entries[o.cursor].ID)
			}
			o.Close()
			return false
		}

		digitKeys := [10]ebiten.Key{
			ebiten.Key0, ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
			ebiten.Key5, ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
		}
		for idx, k := range digitKeys {
			if idx < n && o.keyJustPressed(k) {
				if o.onSelect != nil {
					o.onSelect(o.entries[idx].ID)
				}
				o.Close()
				return false
			}
		}
	}

	if n == 0 {
		o.cursor = 0
	} else if o.cursor >= n {
		o.cursor = n - 1
	}

	return true
}

// Draw renders the overlay onto screen. No-op when the overlay is closed.
func (o *PhaseOverlay) Draw(screen *ebiten.Image) {
	if !o.IsOpen() {
		return
	}

	overlayutil.DrawDimPanel(screen)

	if o.face == nil {
		return
	}

	const (
		xPad   = 10
		yStart = 14
		lineH  = 14
	)
	y := float64(yStart)

	overlayutil.DrawText(screen, o.face, "--- Jump to Phase ---", xPad, y, color.RGBA{180, 180, 180, 255})
	y += lineH

	for i, e := range o.entries {
		c := color.Color(color.White)
		if i == o.cursor {
			c = color.RGBA{255, 255, 0, 255}
		}
		overlayutil.DrawText(screen, o.face, strconv.Itoa(i)+" "+e.Name, xPad, y, c)
		y += lineH
	}
}
