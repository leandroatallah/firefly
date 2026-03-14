package gamesetup

import (
	"flag"
	"time"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

const (
	// Celeste is 320 x 180
	ScreenWidth   = 256
	ScreenHeight  = 240
	DefaultVolume = 0.5
	MainFontFace  = "assets/fonts/pressstart2p.ttf"
	SmallFontFace = "assets/fonts/tiny5.ttf"
)

func NewConfig() *config.AppConfig {
	defaultPhysics := config.PhysicsConfig{
		SpeedMultiplier:       0.25,
		HorizontalInertia:     2.0,
		AirFrictionMultiplier: 0.5,
		AirControlMultiplier:  0.25,
		CoyoteTimeFrames:      timing.FromDuration(100 * time.Millisecond), // 6 frames
		JumpBufferFrames:      timing.FromDuration(100 * time.Millisecond), // 6 frames
		JumpForce:             6,
		JumpCutMultiplier:     0.5,
		UpwardGravity:         6,
		DownwardGravity:       8,
		MaxFallSpeed:          fp16.To16(3),
	}

	cfg := &config.AppConfig{
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		Physics:      defaultPhysics,

		DefaultVolume: DefaultVolume,

		MainFontFace:        MainFontFace,
		SmallFontFace:       SmallFontFace,
		ScreenFlipSpeed:     1.0 / 60.0,
		FadeHoldDuration:    0,
		FadeVisibleDuration: 0,

		EnableSpeechSkip:          false,
		EnableTypingSounds:        true,
		TypingSoundVolume:         0.6,
		TypingSoundCooldownFrames: 15,
	}

	flag.BoolVar(&cfg.CamDebug, "cam-debug", false, "Enable camera debug")
	flag.BoolVar(&cfg.CollisionBox, "collision-box", false, "Enable collision box debug")
	flag.BoolVar(&cfg.NoSound, "no-sound", false, "Disable game sound")
	flag.BoolVar(&cfg.SkipIntro, "skip-intro", false, "Skip all intros and start in phases")
	// Temporary debug flags; only speech-skip is expected to remain as a global debug override.
	flag.BoolVar(&cfg.EnableSpeechSkip, "speech-skip", true, "Enable skipping speech typing with Enter")
	flag.BoolVar(&cfg.EnableTypingSounds, "typing-sounds", true, "Enable typing sound effects")
	flag.Float64Var(&cfg.TypingSoundVolume, "typing-sound-volume", 0.6, "Typing sound effect volume multiplier")
	flag.IntVar(&cfg.TypingSoundCooldownFrames, "typing-sound-cooldown", 5, "Frames between typing sounds")

	return cfg
}
