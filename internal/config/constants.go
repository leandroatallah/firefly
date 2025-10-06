package config

const (
	// TODO: To simplify, we can use the Celeste resolution (320x180)
	ScreenWidth   = 640 // 20 32px tiles
	ScreenHeight  = 360 // 11.25 32px tiles (use 12 tiles)
	Unit          = 16
	DefaultVolume = 0.5

	MainFontFace = "assets/pressstart2p.ttf"
)

// Config holds the main configuration for the application.
// It's intended to be a read-only struct passed around via the AppContext.

type Config struct {
	// You can add fields here that might be loaded from a file in the future.
	ScreenWidth  int
	ScreenHeight int
	Unit         int
}

// NewConfig creates a new Config struct with default values.
func NewConfig() *Config {
	return &Config{
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		Unit:         Unit,
	}
}
