package config

import (
	"flag"

	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

const (
	ScreenWidth   = 320
	ScreenHeight  = 180
	DefaultVolume = 0.5
	MainFontFace  = "assets/fonts/pressstart2p.ttf"
)

type PhysicsConfig struct {
	// HorizontalInertia controls the smoothness of horizontal movement. Higher values lead to more sliding.
	// 0.0 means instant movement.
	HorizontalInertia float64
	// AirFrictionMultiplier controls how much friction is applied in the air, as a factor of ground friction.
	// 0.0 means no air friction; 1.0 means same as ground.
	AirFrictionMultiplier float64
	// AirControlMultiplier controls how much acceleration is applied in the air.
	// < 1.0 for less air control, > 1.0 for more.
	AirControlMultiplier float64
	// CoyoteTimeFrames is the number of frames the player can still jump after leaving a ledge.
	CoyoteTimeFrames int
	// JumpBufferFrames is the number of frames a jump input is remembered before landing.
	JumpBufferFrames int
	// JumpForce is the initial vertical velocity applied when jumping.
	JumpForce int
	// JumpCutMultiplier is the factor applied to vertical velocity when the jump button is released mid-air.
	// Should be between 0.0 and 1.0.
	JumpCutMultiplier float64
	// UpwardGravity is the gravity force applied when the actor is moving up.
	UpwardGravity int
	// DownwardGravity is the gravity force applied when the actor is falling. A higher value than UpwardGravity creates a snappier jump.
	DownwardGravity int
	// MaxFallSpeed is the terminal velocity for falling.
	MaxFallSpeed int
}

type AppConfig struct {
	ScreenWidth  int
	ScreenHeight int
	Physics      PhysicsConfig

	DefaultVolume float64

	MainFontFace string
	CamDebug     bool
	CollisionBox bool
	NoSound      bool
}

var cfg AppConfig

func init() {
	defaultPhysics := PhysicsConfig{
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

	cfg = AppConfig{
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		Physics:      defaultPhysics,

		DefaultVolume: DefaultVolume,

		MainFontFace: MainFontFace,
	}

	Parse()
}

func Parse() {
	flag.BoolVar(&cfg.CamDebug, "cam-debug", false, "Enable camera debug")
	flag.BoolVar(&cfg.CollisionBox, "collision-box", false, "Enable collision box debug")
	flag.BoolVar(&cfg.NoSound, "no-sound", false, "Disable game sound")
}

func Get() *AppConfig {
	return &cfg
}
