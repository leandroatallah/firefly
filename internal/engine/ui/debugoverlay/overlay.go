package debugoverlay

import (
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/debug"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// DebugOverlay is an in-game overlay that lists all registered debug flags
// and allows toggling them interactively.
type DebugOverlay struct {
	open           bool
	cursor         int
	face           *text.GoTextFace
	keyJustPressed func(ebiten.Key) bool
}

// New creates a DebugOverlay wired to the real Ebitengine key-input backend.
func New() *DebugOverlay {
	return &DebugOverlay{
		keyJustPressed: func(k ebiten.Key) bool {
			return inpututil.IsKeyJustPressed(k)
		},
	}
}

// Open makes the overlay visible.
func (o *DebugOverlay) Open() { o.open = true }

// Close hides the overlay.
func (o *DebugOverlay) Close() { o.open = false }

// IsOpen reports whether the overlay is currently visible.
func (o *DebugOverlay) IsOpen() bool { return o.open }

// SetFont sets the font face used when drawing entry labels.
func (o *DebugOverlay) SetFont(f *text.GoTextFace) { o.face = f }

// Update handles input for the overlay. It returns true when the overlay
// consumed the frame (i.e. the overlay is still open after processing),
// and false when the overlay is closed or was already closed.
func (o *DebugOverlay) Update() bool {
	if !o.open {
		return false
	}

	entries := debug.List()
	n := len(entries)

	if n > 0 {
		if o.keyJustPressed(ebiten.KeyArrowUp) {
			o.cursor = (o.cursor - 1 + n) % n
		}
		if o.keyJustPressed(ebiten.KeyArrowDown) {
			o.cursor = (o.cursor + 1) % n
		}
		if o.keyJustPressed(ebiten.KeySpace) || o.keyJustPressed(ebiten.KeyEnter) {
			if entries[o.cursor].Ptr != nil {
				*entries[o.cursor].Ptr = !*entries[o.cursor].Ptr
			}
		}
	} else {
		o.cursor = 0
	}

	if o.keyJustPressed(ebiten.KeyF1) || o.keyJustPressed(ebiten.KeyEscape) {
		o.open = false
		return false
	}

	if n == 0 {
		o.cursor = 0
	} else if o.cursor >= n {
		o.cursor = n - 1
	}

	return true
}

// Draw renders the overlay onto screen. No-op when the overlay is closed.
func (o *DebugOverlay) Draw(screen *ebiten.Image) {
	if !o.open {
		return
	}

	// semi-transparent panel
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	panel := ebiten.NewImage(w, h)
	panel.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(panel, nil)

	if o.face == nil {
		return
	}

	entries := debug.List()
	const (
		xPad     = 10
		yStart   = 14
		lineH    = 14
		groupGap = 8
	)
	y := float64(yStart)
	var prevGroup debug.Group = -1
	for i, e := range entries {
		if e.Group != prevGroup {
			if prevGroup != -1 {
				y += groupGap
			}
			header := "--- " + e.Group.String() + " ---"
			op := &text.DrawOptions{}
			op.GeoM.Translate(xPad, y)
			op.ColorScale.ScaleWithColor(color.RGBA{180, 180, 180, 255})
			text.Draw(screen, header, o.face, op)
			y += lineH
			prevGroup = e.Group
		}

		mark := "[ ]"
		if e.Ptr != nil && *e.Ptr {
			mark = "[x]"
		}
		line := mark + " " + e.Name

		op := &text.DrawOptions{}
		op.GeoM.Translate(xPad, y)
		if i == o.cursor {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 255})
		} else {
			op.ColorScale.ScaleWithColor(color.White)
		}
		text.Draw(screen, line, o.face, op)
		y += lineH
	}
}
