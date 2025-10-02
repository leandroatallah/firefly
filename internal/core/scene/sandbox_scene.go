package scene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/actors/enemies"
	"github.com/leandroatallah/firefly/internal/actors/movement"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/hud"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

const (
	jab8         = "assets/jab8.ogg"
	sketchbookBG = "assets/Sketchbook 2024-06-19.ogg"
)

type SandboxScene struct {
	BaseScene
	player            *actors.Player
	isPlayingJab      bool
	showMenu          bool
	menuDeadzoneCount int
}

func (s *SandboxScene) Update() error {
	space := s.PhysicsSpace()

	for _, i := range space.Bodies() {
		actor, ok := i.(actors.ActorEntity)
		if !ok {
			continue
		}
		if err := actor.Update(space); err != nil {
			return err
		}
	}

	// Key events
	s.menuDeadzoneCount++
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if s.menuDeadzoneCount > 10 {
			s.showMenu = !s.showMenu
			s.menuDeadzoneCount = 0
		}
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

	space := s.PhysicsSpace()
	for _, b := range space.Bodies() {
		// TODO: Fix it
		// TODO: Is player calling Draw twice?
		switch b.(type) {
		case *enemies.BlueEnemy:
			b.(*enemies.BlueEnemy).Draw(screen)
		default:
			b.(physics.Obstacle).Draw(screen)
		}
	}

	// HUD
	statusBar, err := hud.NewStatusBar(s.player)
	if err != nil {
		log.Fatal(err)
	}
	statusBar.Draw(screen)

	if s.showMenu {
		shadow := ebiten.NewImage(config.ScreenWidth, config.ScreenWidth)
		shadow.Fill(color.RGBA{0, 0, 0, 0xCC})
		screen.DrawImage(shadow, nil)

		containerWidth, containerHeight := config.ScreenWidth/3, config.ScreenHeight/2
		container := ebiten.NewImage(containerWidth, containerHeight)
		container.Fill(color.RGBA{0xAA, 0xAA, 0xAA, 0xff})
		containerOp := &ebiten.DrawImageOptions{}
		containerOp.GeoM.Translate(config.ScreenWidth/2, config.ScreenHeight/2)
		containerOp.GeoM.Translate(-float64(containerWidth/2), -float64(containerHeight/2))
		screen.DrawImage(container, containerOp)
	}
}

func (s *SandboxScene) OnStart() {
	// Init space
	space := s.PhysicsSpace()

	// Init audio manager
	s.audiomanager = s.Manager.audioManager
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.SetVolume(0)
		s.audiomanager.PlaySound(sketchbookBG)
	}()

	const wallWidth = 20

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

	s.AddBoundaries(box)

	// Create Player
	s.player = actors.NewPlayer()
	s.PhysicsSpace().AddBody(s.player)

	// Create enemies
	// TODO: It should be a builder
	enemyFactory := enemies.NewEnemyFactory()
	blueEnemy, err := enemyFactory.Create(enemies.BlueEnemyType, 60, 30)
	if err != nil {
		log.Fatal(err)
	}
	// // Set up patrol movement with predefined waypoints
	// waypoints := []*physics.Rect{
	// 	physics.NewRect(100, 100, 32, 32),
	// 	physics.NewRect(200, 100, 32, 32),
	// 	physics.NewRect(200, 200, 32, 32),
	// 	physics.NewRect(100, 200, 32, 32),
	// }
	// predefinedConfig := movement.NewPredefinedWaypointConfig(waypoints, 120) // 2 seconds at 60 FPS
	// blueEnemy.SetMovementState(
	// 	movement.Patrol, s.player, movement.WithWaypointConfig(predefinedConfig),
	// )
	blueEnemy.SetMovementState(movement.Chase, s.player, movement.WithObstacles(space.Bodies()))

	s.AddBoundaries(blueEnemy.(physics.Body))
}

func (s *SandboxScene) OnFinish() {
	s.audiomanager.PauseMusic(sketchbookBG)
}
