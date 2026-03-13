package config

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	originalCfg := Get()

	newCfg := &AppConfig{
		ScreenWidth:  1280,
		ScreenHeight: 720,
		Physics: PhysicsConfig{
			SpeedMultiplier: 2.0,
		},
		DefaultVolume: 0.5,
	}

	Set(newCfg)

	retrievedCfg := Get()
	if retrievedCfg.ScreenWidth != 1280 {
		t.Errorf("expected ScreenWidth 1280, got %d", retrievedCfg.ScreenWidth)
	}
	if retrievedCfg.ScreenHeight != 720 {
		t.Errorf("expected ScreenHeight 720, got %d", retrievedCfg.ScreenHeight)
	}
	if retrievedCfg.Physics.SpeedMultiplier != 2.0 {
		t.Errorf("expected Physics.SpeedMultiplier 2.0, got %f", retrievedCfg.Physics.SpeedMultiplier)
	}
	if retrievedCfg.DefaultVolume != 0.5 {
		t.Errorf("expected DefaultVolume 0.5, got %f", retrievedCfg.DefaultVolume)
	}

	Set(originalCfg)
}

func TestSetNil(t *testing.T) {
	originalCfg := Get()

	Set(nil)

	retrievedCfg := Get()
	if retrievedCfg != originalCfg {
		t.Error("expected config to remain unchanged when Set(nil) is called")
	}
}
