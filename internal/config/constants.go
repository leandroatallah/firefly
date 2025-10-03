package config

const (
	ScreenWidth   = 480
	ScreenHeight  = 360
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