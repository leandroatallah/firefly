package actorinspector

import (
	"image/color"
	"sort"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/overlayutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// statsProvider is implemented by actors that expose debug stats lines.
type statsProvider interface {
	DebugStats() []string
}

// ActorSource is the minimal interface the overlay needs from the actor manager.
type ActorSource interface {
	ForEach(func(actors.ActorEntity))
	Find(id string) (actors.ActorEntity, bool)
}

// Overlay is an in-game two-panel actor inspector toggled by F5.
// Left panel: sorted actor ID list. Right panel: stats for the selected actor.
type Overlay struct {
	overlayutil.Base
	cursor         int
	ids            []string // sorted snapshot rebuilt each Open
	face           *text.GoTextFace
	keyJustPressed func(ebiten.Key) bool
	source         func() ActorSource
}

// New creates an Overlay. source is called each time the overlay opens
// to snapshot the current actor list.
func New(source func() ActorSource) *Overlay {
	return &Overlay{
		keyJustPressed: func(k ebiten.Key) bool { return inpututil.IsKeyJustPressed(k) },
		source:         source,
	}
}

// SetFont sets the font face used for drawing.
func (o *Overlay) SetFont(f *text.GoTextFace) { o.face = f }

// Open makes the overlay visible and snapshots the current actor list.
func (o *Overlay) Open() {
	o.Base.Open()
	o.ids = o.snapshotIDs()
	if o.cursor >= len(o.ids) {
		o.cursor = 0
	}
}

// Close hides the overlay.
func (o *Overlay) Close() { o.Base.Close() }

// IsOpen reports whether the overlay is currently visible.
func (o *Overlay) IsOpen() bool { return o.Base.IsOpen() }

func (o *Overlay) snapshotIDs() []string {
	m := o.source()
	if m == nil {
		return nil
	}
	var ids []string
	m.ForEach(func(a actors.ActorEntity) { ids = append(ids, a.ID()) })
	sort.Strings(ids)
	return ids
}

// Update handles input. Returns true while the overlay is open.
func (o *Overlay) Update() bool {
	if !o.IsOpen() {
		return false
	}

	if o.keyJustPressed(ebiten.KeyF5) || o.keyJustPressed(ebiten.KeyEscape) {
		o.Close()
		return false
	}

	n := len(o.ids)
	if n > 0 {
		if o.keyJustPressed(ebiten.KeyArrowUp) {
			o.cursor = (o.cursor - 1 + n) % n
		}
		if o.keyJustPressed(ebiten.KeyArrowDown) {
			o.cursor = (o.cursor + 1) % n
		}
	}

	return true
}

// Draw renders the overlay. No-op when closed.
func (o *Overlay) Draw(screen *ebiten.Image) {
	if !o.IsOpen() {
		return
	}

	overlayutil.DrawDimPanel(screen)

	if o.face == nil {
		return
	}

	const (
		lineH  = 14
		yStart = 14
		xPad   = 10
	)
	splitX := float64(screen.Bounds().Dx()) * 0.30

	// --- Left panel: actor list ---
	y := float64(yStart)
	overlayutil.DrawText(screen, o.face, "--- Actors ---", xPad, y, color.RGBA{180, 180, 180, 255})
	y += lineH

	for i, id := range o.ids {
		c := color.Color(color.White)
		if i == o.cursor {
			c = color.RGBA{255, 255, 0, 255}
		}
		overlayutil.DrawText(screen, o.face, id, xPad, y, c)
		y += lineH
	}

	// --- Right panel: stats ---
	if len(o.ids) == 0 {
		return
	}
	selectedID := o.ids[o.cursor]

	m := o.source()
	if m == nil {
		return
	}
	actor, found := m.Find(selectedID)
	if !found {
		return
	}

	ry := float64(yStart)
	overlayutil.DrawText(screen, o.face, "--- "+selectedID+" ---", splitX+xPad, ry, color.RGBA{180, 180, 180, 255})
	ry += lineH

	if sp, ok := actor.(statsProvider); ok {
		for _, line := range sp.DebugStats() {
			overlayutil.DrawText(screen, o.face, line, splitX+xPad, ry, color.White)
			ry += lineH
		}
	} else {
		overlayutil.DrawText(screen, o.face, "id    = "+actor.ID(), splitX+xPad, ry, color.White)
	}
}
