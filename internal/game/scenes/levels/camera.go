package gamescenelevels

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

func (s *LevelsScene) CamDebug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.cam.Kamera().Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		s.cam.Kamera().Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		s.cam.Kamera().Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		s.cam.Kamera().ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		s.cam.Kamera().ZoomFactor *= 1.02
	}
}

func (s *LevelsScene) SetCamTargetPointToSpace() {
	tPos := s.cam.Target().Position()
	targetRect := physics.NewObstacleRect(physics.NewRect(tPos.Min.X, tPos.Min.Y, tPos.Dx(), tPos.Dy()))
	targetBody := physics.NewPhysicsBody(targetRect)
	targetBody.SetID("TARGET")
	s.PhysicsSpace().AddBody(targetBody)
}

func (s *LevelsScene) DrawCamTargetPoint(screen *ebiten.Image) {
	tPos := s.cam.Target().Position()
	targetImage := ebiten.NewImage(tPos.Dx(), tPos.Dy())
	targetImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(tPos.Min.X), float64(tPos.Min.Y))
	s.cam.Draw(targetImage, opts, screen)
}
