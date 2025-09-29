package scene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/actors/enemies"
	"github.com/leandroatallah/firefly/internal/actors/movement"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

const (
	jab8         = "assets/jab8.ogg"
	sketchbookBG = "assets/Sketchbook 2024-06-19.ogg"
)

type SandboxScene struct {
	BaseScene
	player       *actors.Player
	isPlayingJab bool
}

func (s *SandboxScene) Update() error {
	if s.player != nil {
		s.player.Update(s.boundaries)
	}
	for _, i := range s.boundaries {
		actor, ok := i.(actors.ActorEntity)
		if ok {
			actor.Update(s.boundaries)
		}
	}

	// Key events
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.Manager.NavigateTo(SceneMenu, nil)
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
	if s.player != nil {
		s.player.Draw(screen)
	}

	for _, b := range s.boundaries {
		// TODO: Fix it
		switch b.(type) {
		case *enemies.BlueEnemy:
			b.(*enemies.BlueEnemy).Draw(screen)
		default:
			b.(physics.Obstacle).Draw(screen)
		}
	}
}

func (s *SandboxScene) OnStart() {
	s.audiomanager = s.Manager.audioManager
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.SetVolume(0)
		s.audiomanager.PlaySound(sketchbookBG)
	}()

	const wallWidth = 20

	s.player = actors.NewPlayer()

	// TODO: It should be a builder
	enemyFactory := enemies.NewEnemyFactory()
	blueEnemy, err := enemyFactory.Create(enemies.BlueEnemyType, 60, 30)
	if err != nil {
		log.Fatal(err)
	}
	blueEnemy.SetMovementState(movement.DumbChase, s.player)

	obstacleFactory := physics.NewDefaultObstacleFactory()

	// Boundaries
	boundaries := []physics.ObstacleType{
		physics.ObstacleWallTop,
		physics.ObstacleWallLeft,
		physics.ObstacleWallRight,
		physics.ObstacleWallDown,
	}
	for _, o := range boundaries {
		b, err := obstacleFactory.Create(o)
		if err != nil {
			log.Fatal(err)
		}
		b.SetIsObstructive(true)
		s.AddBoundaries(b)
	}

	box := physics.NewObstacleRect(
		physics.NewRect(100, 100, 32, 32),
	).AddCollision()
	box.SetIsObstructive(true)

	s.AddBoundaries(
		// TODO: Should it be added here?
		blueEnemy.(physics.Body),
		box,
	)
}

func (s *SandboxScene) OnFinish() {
	s.audiomanager.PauseMusic(sketchbookBG)
}
