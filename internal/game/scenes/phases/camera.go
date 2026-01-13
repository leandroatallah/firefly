package gamescenephases

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func (s *PhasesScene) CamDebug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.Camera().Kamera().Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		s.Camera().Kamera().Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		s.Camera().Kamera().Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		s.Camera().Kamera().ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		s.Camera().Kamera().ZoomFactor *= 1.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		s.Camera().Kamera().CenterOffsetX *= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		s.Camera().Kamera().CenterOffsetX /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		s.Camera().Kamera().CenterOffsetY /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		s.Camera().Kamera().CenterOffsetY *= 1.02
	}

}

func (s *PhasesScene) DrawCamTargetPoint(screen *ebiten.Image) {
	tPos := s.Camera().Target().Position()
	targetImage := ebiten.NewImage(tPos.Dx(), tPos.Dy())
	targetImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(tPos.Min.X), float64(tPos.Min.Y))
	s.Camera().Draw(targetImage, opts, screen)
}
