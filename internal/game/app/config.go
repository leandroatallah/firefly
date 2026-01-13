package gamesetup

import (
	"flag"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

const (
	ScreenWidth   = 320
	ScreenHeight  = 192
	DefaultVolume = 0.5
	MainFontFace  = "assets/fonts/pressstart2p.ttf"
)

func NewConfig() config.AppConfig {
	defaultPhysics := config.PhysicsConfig{
		HorizontalInertia:     2.0,
		AirFrictionMultiplier: 0.5,
		AirControlMultiplier:  0.25,
		CoyoteTimeFrames:      6,
		JumpBufferFrames:      6,
		JumpForce:             6,
		JumpCutMultiplier:     0.5,
		UpwardGravity:         6,
		DownwardGravity:       6,
		MaxFallSpeed:          fp16.To16(4),
	}

	cfg := config.AppConfig{
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		Physics:      defaultPhysics,

		DefaultVolume: DefaultVolume,

		MainFontFace: MainFontFace,
	}

	flag.BoolVar(&cfg.CamDebug, "cam-debug", false, "Enable camera debug")
	flag.BoolVar(&cfg.CollisionBox, "collision-box", false, "Enable collision box debug")
	flag.BoolVar(&cfg.NoSound, "no-sound", false, "Disable game sound")

	return cfg
}
