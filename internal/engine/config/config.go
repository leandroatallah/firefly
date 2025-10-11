package config

type BaseConfig struct {
	ScreenWidth  int
	ScreenHeight int
	Unit         int
}

// Config holds the main configuration for the application.
// It's intended to be a read-only struct passed around via the AppContext.
type Config struct {
	// Use embedded BaseConfig to allow extended fields only for Config.
	BaseConfig
}

// NewConfig creates a new Config struct with default values
func NewConfig(base BaseConfig) *Config {
	return &Config{BaseConfig: base}
}
