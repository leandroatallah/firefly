// Package shadow renders ground-plane shadows for airborne beat-em-up actors.
package shadow

import (
	"image/color"
	"math"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Tunables (SPEC §2). Implementer must keep these exact values.
const (
	ShadowAlpha           = 0.50
	ShadowBaseWidthRatio  = 0.75
	ShadowBaseHeight      = 4
	ShadowMinScale        = 0.30
	ShadowAltitudeFalloff = 64.0
)

// ShadowColor is semi-transparent black at ShadowAlpha opacity.
//
//nolint:gochecknoglobals
var ShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: shadowAlphaByte()}

func shadowAlphaByte() uint8 {
	v := 255.0 * ShadowAlpha
	return uint8(v)
}

// AltitudeBody is the minimum body interface the shadow needs.
type AltitudeBody interface {
	GetPositionMin() (x, y int)
	GetShape() body.Shape
	Altitude() int
}

// Bounds describes the computed oval rectangle (in world-pixel coords).
type Bounds struct {
	CenterX, CenterY float64
	Width, Height    float64
}

// drawOval is the production drawer wired into ovalDrawerFn.
func drawOval(screen *ebiten.Image, cam *camera.Controller, b Bounds, c color.Color) {
	w := math.Round(b.Width)
	h := math.Round(b.Height)
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	iw := int(w)
	ih := int(h)

	img := ebiten.NewImage(iw, ih)
	vector.DrawFilledCircle(img, float32(w/2), float32(h/2), float32(w/2), c, false)

	r, g, bl, a := c.RGBA()
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetR(float32(r) / 0xffff)
	op.ColorScale.SetG(float32(g) / 0xffff)
	op.ColorScale.SetB(float32(bl) / 0xffff)
	op.ColorScale.SetA(float32(a) / 0xffff)

	// Scale to get the oval shape (circle was drawn using w; scale y to h/w ratio).
	scaleX := 1.0
	scaleY := h / w
	op.GeoM.Scale(scaleX, scaleY)

	// Translate so the oval is centered at world position (CenterX, CenterY).
	op.GeoM.Translate(b.CenterX-w/2, b.CenterY-h/2)

	cam.Draw(img, op, screen)
}

// ovalDrawerFn is the swappable drawer used by Draw/DrawAll.
//
//nolint:gochecknoglobals
var ovalDrawerFn func(screen *ebiten.Image, cam *camera.Controller, b Bounds, c color.Color) = drawOval

// SetOvalDrawerForTest replaces ovalDrawerFn and returns a restore func.
func SetOvalDrawerForTest(f func(*ebiten.Image, *camera.Controller, Bounds, color.Color)) (restore func()) {
	prev := ovalDrawerFn
	ovalDrawerFn = f
	return func() { ovalDrawerFn = prev }
}

// ScaleFor returns the linear scale factor for the given altitude.
func ScaleFor(alt int) float64 {
	if alt <= 0 {
		return 1.0
	}
	t := float64(alt) / ShadowAltitudeFalloff
	if t > 1.0 {
		t = 1.0
	}
	return 1.0 - t*(1.0-ShadowMinScale)
}

// ComputeBounds returns the oval bounds for b at its current altitude.
func ComputeBounds(b AltitudeBody) Bounds {
	x, y := b.GetPositionMin()
	w, h := b.GetShape().Width(), b.GetShape().Height()
	s := ScaleFor(b.Altitude())
	cx := float64(x) + float64(w)/2
	cy := float64(y) + float64(h) // foot midpoint on ground plane
	bw := float64(w) * ShadowBaseWidthRatio * s
	bh := float64(ShadowBaseHeight) * s
	return Bounds{cx, cy, bw, bh}
}

// Draw renders a single shadow for b.
func Draw(screen *ebiten.Image, cam *camera.Controller, b AltitudeBody) bool {
	if b == nil {
		return false
	}
	if b.Altitude() <= 0 {
		return false
	}
	bn := ComputeBounds(b)
	ovalDrawerFn(screen, cam, bn, ShadowColor)
	return true
}

// DrawAll iterates bodies and draws shadows for airborne Altitudables.
func DrawAll(screen *ebiten.Image, cam *camera.Controller, bodies []body.Collidable) {
	for _, c := range bodies {
		a, ok := c.(AltitudeBody)
		if !ok {
			continue
		}
		Draw(screen, cam, a)
	}
}
