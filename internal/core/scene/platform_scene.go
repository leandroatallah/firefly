package scene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/config"
)

const (
	bgSound = "assets/Sketchbook 2024-06-19.ogg"
)

type PlatformScene struct {
	BaseScene
	count  int
	player actors.PlayerEntity
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

	s.DrawGround(screen)

	if s.player != nil {
		s.player.Draw(screen)
	}
}

func (s *PlatformScene) OnStart() {
	space := s.PhysicsSpace()

	// Init audio manager
	s.audiomanager = s.Manager.AudioManager()
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.SetVolume(0)
		s.audiomanager.PlaySound(bgSound)
	}()

	var err error
	s.player, err = actors.NewPlayer(actors.Platform)
	if err != nil {
		log.Fatal(err)
	}
	space.AddBody(s.player)
}

func (s *PlatformScene) OnFinish() {
	s.audiomanager.PauseMusic(bgSound)
}

func (s *PlatformScene) DrawGround(screen *ebiten.Image) {
	groundHeight := 120
	ground := ebiten.NewImage(config.ScreenWidth, groundHeight)
	ground.Fill(color.RGBA{0x99, 0x99, 0x99, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(config.ScreenHeight-groundHeight))
	screen.DrawImage(ground, op)
}
