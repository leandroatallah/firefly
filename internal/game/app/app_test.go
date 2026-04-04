package gamesetup

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Reset the global flag set between tests if needed
	// But since TestMain runs once, we just run the tests
	os.Exit(m.Run())
}

func TestNewConfig(t *testing.T) {
	// We need to reset flag.CommandLine before calling NewConfig if other tests call it
	// or if we call it multiple times in this test.
	// Since this is the only test calling it so far, it's fine.

	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig returned nil")
	}

	if cfg.ScreenWidth != ScreenWidth {
		t.Errorf("expected ScreenWidth %d, got %d", ScreenWidth, cfg.ScreenWidth)
	}

	if cfg.DefaultVolume != DefaultVolume {
		t.Errorf("expected DefaultVolume %f, got %f", DefaultVolume, cfg.DefaultVolume)
	}
}

func TestGetPhases(t *testing.T) {
	phases := GetPhases()
	if len(phases) == 0 {
		t.Fatal("GetPhases returned empty slice")
	}

	// Basic check for phase 1
	if phases[0].ID != 1 {
		t.Errorf("expected phase 1 ID to be 1, got %d", phases[0].ID)
	}
}
