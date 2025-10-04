package scene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

const (
	bgSound = "assets/Sketchbook 2024-06-19.ogg"
)

type PlatformScene struct {
	BaseScene
	count  int
	player actors.PlayerEntity
	space  *physics.Space
}

func (s *PlatformScene) OnStart() {
	// Init audio manager
	s.audiomanager = s.Manager.AudioManager()
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.SetVolume(0)
		s.audiomanager.PlaySound(bgSound)
	}()

	// Init boundaries
	s.space = s.PhysicsSpace()

	// Ground
	groundHeight := 40
	createPlatform(
		physics.NewRect(0, config.ScreenHeight-groundHeight, config.ScreenWidth, groundHeight),
		s.space,
	)

	// Flying platform
	createPlatform(
		physics.NewRect(120, 140, 100, 30),
		s.space,
	)
	createPlatform(
		physics.NewRect(280, 200, 100, 30),
		s.space,
	)

	player, err := createPlayer(s.space)
	if err != nil {
		log.Fatal(err)
	}
	s.player = player

	// Parse space bodies to scene boundaries
	for _, o := range s.space.Bodies() {
		s.AddBoundaries(o)
	}

	// Create player
	s.space.AddBody(s.player)
}

func (s *PlatformScene) Update() error {
	space := s.PhysicsSpace()

	s.count++

	for _, i := range space.Bodies() {
		actor, ok := i.(actors.ActorEntity)
		if !ok {
			continue
		}
		if err := actor.Update(space); err != nil {
			return err
		}
	}

	return nil
}

func (s *PlatformScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x33, 0x33, 0xff})

	space := s.PhysicsSpace()

	for _, b := range space.Bodies() {
		switch b.(type) {
		case actors.PlayerEntity:
			b.(actors.PlayerEntity).Draw(screen)
		case physics.Obstacle:
			b.(physics.Obstacle).DrawCollisionBox(screen)
		}
	}
}

func (s *PlatformScene) OnFinish() {
	s.audiomanager.PauseMusic(bgSound)
}

func createPlatform(rect *physics.Rect, space *physics.Space) {
	o := physics.NewObstacleRect(rect).AddCollision()
	o.SetIsObstructive(true)

	space.AddBody(o)
}

func createPlayer(space *physics.Space) (actors.PlayerEntity, error) {
	player, err := actors.NewPlayer(actors.Platform)
	if err != nil {
		return nil, err
	}
	space.AddBody(player)
	return player, nil
}
