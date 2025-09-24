package scene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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
		s.Manager.GoToScene(SceneMenu, nil)
	}

	return nil
}

func (s *SandboxScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xcc, 0xcc, 0xdd, 0xff})
	s.player.Draw(screen)

	for _, b := range s.boundaries {
		b.(physics.Obstacle).Draw(screen)
		b.(physics.Obstacle).DrawCollisionBox(screen)
	}
}

func (s *SandboxScene) OnStart() {
	const wallWidth = 20

	s.player = physics.NewPlayer()

	obstacleFactory := physics.NewDefaultObstacleFactory()

	// Boundaries
	wallTop, err := obstacleFactory.Create(physics.ObstacleWallTop)
	if err != nil {
		log.Fatal(err)
	}
	wallLeft, err := obstacleFactory.Create(physics.ObstacleWallLeft)
	if err != nil {
		log.Fatal(err)
	}
	wallRight, err := obstacleFactory.Create(physics.ObstacleWallRight)
	if err != nil {
		log.Fatal(err)
	}
	wallDown, err := obstacleFactory.Create(physics.ObstacleWallDown)
	if err != nil {
		log.Fatal(err)
	}

	// Enemies
	enemyRect := physics.NewObstacleRect(
		physics.NewRect(100, 100, 32, 32),
	).AddCollision()

	s.AddBoundaries(
		wallTop,
		wallLeft,
		wallRight,
		wallDown,
		enemyRect,
	)
}

func (s *SandboxScene) OnFinish() {}
