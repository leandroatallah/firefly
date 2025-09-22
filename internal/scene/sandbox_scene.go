package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/physics"
)

type SandboxScene struct {
	BaseScene
	player *physics.Player
}

func (s *SandboxScene) Update() error {
	if s.player != nil {
		s.player.Update(s.boundaries)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.nextScene = &MenuScene{}
	}

	return nil
}

func (s *SandboxScene) Draw(screen *ebiten.Image) {
	s.player.Draw(screen)

	for _, b := range s.boundaries {
		b.(physics.Obstacle).Draw(screen)
	}
}

func (s *SandboxScene) OnStart() {
	const wallWidth = 20

	s.player = physics.NewPlayer()

	// Boundaries
	wallTop := physics.NewObstacleRect(
		physics.NewRect(0, 0, config.ScreenWidth, wallWidth),
		[]*physics.CollisionArea{physics.NewCollisionArea(physics.NewRect(0, 0, config.ScreenWidth, wallWidth))},
	)
	wallLeft := physics.NewObstacleRect(
		physics.NewRect(0, 0, wallWidth, config.ScreenHeight), nil,
	)
	wallRight := physics.NewObstacleRect(
		physics.NewRect(config.ScreenWidth-wallWidth, 0, wallWidth, config.ScreenHeight), nil,
	)
	wallDown := physics.NewObstacleRect(
		physics.NewRect(0, config.ScreenHeight-wallWidth, config.ScreenWidth, wallWidth), nil,
	)

	// Enemies
	enemyRect := physics.NewObstacleRect(physics.NewRect(100, 100, 32, 32), nil)

	s.AddBoundaries(
		wallTop,
		wallLeft,
		wallRight,
		wallDown,
		enemyRect,
	)
}

func (s *SandboxScene) OnFinish() {}
