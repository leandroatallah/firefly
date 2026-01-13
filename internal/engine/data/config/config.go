package config

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

func Set(newCfg AppConfig) {
	cfg = newCfg
}

func Get() *AppConfig {
	return &cfg
}