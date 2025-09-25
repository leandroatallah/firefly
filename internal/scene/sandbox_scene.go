package scene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/physics"
)

const (
	jab8         = "assets/jab8.ogg"
	sketchbookBG = "assets/Sketchbook 2024-06-19.ogg"
)

type SandboxScene struct {
	BaseScene
	player       *physics.Player
	isPlayingJab bool
}

func (s *SandboxScene) Update() error {
	if s.player != nil {
		s.player.Update(s.boundaries)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.Manager.GoToScene(SceneMenu, nil)
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !s.isPlayingJab {
		s.audiomanager.PlaySound(jab8)
		s.isPlayingJab = true
		go func() {
			time.Sleep(200 * time.Millisecond)
			s.isPlayingJab = false
		}()
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
	s.audiomanager = s.Manager.audioManager
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.PlaySound(sketchbookBG)
	}()

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

func (s *SandboxScene) OnFinish() {
	s.audiomanager.PauseMusic(sketchbookBG)
}
