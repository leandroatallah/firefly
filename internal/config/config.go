package config

// TODO: Use a env file
const (
	ScreenWidth   = 320
	ScreenHeight  = 180
	Unit          = 16
	DefaultVolume = 0.5
	MainFontFace  = "assets/pressstart2p.ttf"
)

type AppConfig struct {
	ScreenWidth  int
	ScreenHeight int
	Unit         int

	DefaultVolume float64

	MainFontFace string
}

var cfg AppConfig

func init() {
	cfg = AppConfig{
		ScreenWidth:  ScreenWidth,
		ScreenHeight: ScreenHeight,
		Unit:         Unit,

		DefaultVolume: DefaultVolume,

		MainFontFace: MainFontFace,
	}
}

func Get() AppConfig {
	return cfg
}
