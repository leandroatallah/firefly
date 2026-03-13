package text

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
)

type FloatingText interface {
	Update() error
	Draw(screen *ebiten.Image, cam *camera.Controller)
	IsComplete() bool
	SetFont(f *font.FontText)
	SetColor(c color.Color)
}

type FloatingTextBase struct {
	Text        string
	Duration    int
	MaxDuration int
	Font        *font.FontText
	Size        float64
	Color       color.Color
	removed     bool
}

func (ft *FloatingTextBase) IsComplete() bool {
	return ft.removed
}

func (ft *FloatingTextBase) Update() error {
	ft.Duration--
	if ft.Duration <= 0 {
		ft.removed = true
	}
	return nil
}

func (ft *FloatingTextBase) SetFont(f *font.FontText) {
	ft.Font = f
}

func (ft *FloatingTextBase) SetColor(c color.Color) {
	ft.Color = c
}

// SimpleFloatingText is a simple floating text implementation using FloatingTextBase.
type SimpleFloatingText struct {
	FloatingTextBase
	X, Y      float64
	VelocityY float64
}

// NewFloatingText creates a new floating text effect (stationary by default).
func NewFloatingText(msg string, x, y float64, duration int) *SimpleFloatingText {
	return &SimpleFloatingText{
		FloatingTextBase: FloatingTextBase{
			Text:        msg,
			Duration:    duration,
			MaxDuration: duration,
			Size:        8,
			Color:       color.White,
		},
		X:         x,
		Y:         y,
		VelocityY: 0, // No movement by default
	}
}

// NewFloatingTextWithVelocity creates a floating text effect with vertical movement.
func NewFloatingTextWithVelocity(msg string, x, y float64, duration int, velocityY float64) *SimpleFloatingText {
	ft := NewFloatingText(msg, x, y, duration)
	ft.VelocityY = velocityY
	return ft
}

// Update updates the floating text position and duration.
func (ft *SimpleFloatingText) Update() error {
	ft.Y += ft.VelocityY
	return ft.FloatingTextBase.Update()
}

// Draw draws the floating text at its current position.
func (ft *SimpleFloatingText) Draw(screen *ebiten.Image, cam *camera.Controller) {
	ft.FloatingTextBase.DrawText(screen, ft.X, ft.Y, cam)
}

func (ft *FloatingTextBase) DrawText(screen *ebiten.Image, x, y float64, cam *camera.Controller) {
	if ft.Font == nil {
		return
	}

	face := ft.Font.NewFace(ft.Size)
	textWidth, _ := text.Measure(ft.Text, face, 0)

	op := &text.DrawOptions{}

	if cam != nil {
		centerX, centerY := cam.Kamera().Center()
		screenW, screenH := float64(480), float64(640)
		if cfg := config.Get(); cfg != nil {
			screenW = float64(cfg.ScreenWidth)
			screenH = float64(cfg.ScreenHeight)
		}
		topLeftX := centerX - screenW/2
		topLeftY := centerY - screenH/2

		screenX := x - topLeftX
		screenY := y - topLeftY
		screenX -= textWidth / 2

		op.GeoM.Translate(screenX, screenY)
	} else {
		op.GeoM.Translate(x-textWidth/2, y)
	}

	op.ColorScale.ScaleWithColor(ft.Color)

	ft.Font.Draw(screen, ft.Text, ft.Size, op)
}
