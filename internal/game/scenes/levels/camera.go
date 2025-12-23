package gamescenelevels

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		s.cam.Kamera().CenterOffsetX *= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		s.cam.Kamera().CenterOffsetX /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		s.cam.Kamera().CenterOffsetY /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		s.cam.Kamera().CenterOffsetY *= 1.02
	}

}

func (s *LevelsScene) SetCamTargetPointToSpace() {
	// tPos := s.cam.Target().Position()
	// rect := physics.NewRect(tPos.Min.X, tPos.Min.Y, tPos.Dx(), tPos.Dy())
	// targetRect := physics.NewObstacleRect(rect)
	// targetRect.SetBody(rect)
	// targetBody := physics.NewMovableBody(targetRect.Body)
	// // FIX: Uncomment and fix this
	// // targetBody := &actors.Character{
	// // 	MovableCollidableAliveBody: *physics.NewMovableCollidableAliveBody(targetRect),
	// // }
	// targetBody.SetID("TARGET")
	// // NOTE: Camera target should be added to space for debugging purpose.
	// // s.PhysicsSpace().AddBody(targetBody)
	// s.PhysicsSpace().AddBody(s.player)
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
