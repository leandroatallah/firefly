package vfx

import (
	"encoding/json"
	"image/color"
	"os"
	"testing"
	"testing/fstest"
)

// realVFXFS returns a filesystem rooted at the repo's assets/particles dir,
// so tests can load the production vfx.json without copying it.
func realVFXFS(t *testing.T) (string, []byte) {
	t.Helper()
	path := "../../../../../assets/particles/vfx.json"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vfx.json: %v", err)
	}
	return path, data
}

func TestVFXManager_LoadsCombatPixelTypes(t *testing.T) {
	m := NewManagerFromPath("../../../../../assets/particles/vfx.json")
	if m == nil {
		t.Fatal("NewManagerFromPath returned nil")
	}

	for _, typeKey := range []string{"muzzle_flash", "bullet_impact", "bullet_despawn"} {
		if _, ok := m.configs[typeKey]; !ok {
			t.Errorf("config %q not loaded from vfx.json", typeKey)
		}
	}
}

func TestVFXManager_PixelConfigLifetimeAndColor(t *testing.T) {
	m := NewManagerFromPath("../../../../../assets/particles/vfx.json")
	if m == nil {
		t.Fatal("NewManagerFromPath returned nil")
	}

	cases := map[string]int{
		"muzzle_flash":   3,
		"bullet_impact":  6,
		"bullet_despawn": 8,
	}

	for typeKey, wantLifetime := range cases {
		cfg, ok := m.configs[typeKey]
		if !ok {
			t.Errorf("config %q missing", typeKey)
			continue
		}
		if cfg.Lifetime != wantLifetime {
			t.Errorf("%s: Lifetime = %d, want %d", typeKey, cfg.Lifetime, wantLifetime)
		}
		if cfg.Color == nil {
			t.Errorf("%s: Color is nil, want color.White", typeKey)
			continue
		}
		gotR, gotG, gotB, gotA := cfg.Color.RGBA()
		wantR, wantG, wantB, wantA := color.White.RGBA()
		if gotR != wantR || gotG != wantG || gotB != wantB || gotA != wantA {
			t.Errorf("%s: Color RGBA = (%d,%d,%d,%d), want white (%d,%d,%d,%d)",
				typeKey, gotR, gotG, gotB, gotA, wantR, wantG, wantB, wantA)
		}
	}
}

func TestVFXManager_SpawnPuffUsesPixelLifetime(t *testing.T) {
	jsonBytes := []byte(`[
		{"type": "test_pixel", "pixel": {"size": 1, "color": "#FFFFFF", "lifetime_frames": 10}}
	]`)
	fsys := fstest.MapFS{
		"vfx.json": &fstest.MapFile{Data: jsonBytes},
	}

	m := NewManager(fsys, "vfx.json")
	if _, ok := m.configs["test_pixel"]; !ok {
		t.Fatal("test_pixel config not loaded")
	}

	m.SpawnPuff("test_pixel", 0, 0, 1, 0)

	parts := m.system.Particles()
	if len(parts) != 1 {
		t.Fatalf("expected 1 particle, got %d", len(parts))
	}
	p := parts[0]
	if p.Duration != 10 {
		t.Errorf("Duration = %d, want 10", p.Duration)
	}
	if p.MaxDuration != 10 {
		t.Errorf("MaxDuration = %d, want 10", p.MaxDuration)
	}
}

func TestVFXManager_RejectsInvalidPixelColor(t *testing.T) {
	jsonBytes := []byte(`[
		{"type": "bad_color", "pixel": {"size": 1, "color": "#123456", "lifetime_frames": 5}}
	]`)
	fsys := fstest.MapFS{
		"vfx.json": &fstest.MapFile{Data: jsonBytes},
	}

	// Must not panic on invalid color.
	m := NewManager(fsys, "vfx.json")

	cfg, ok := m.configs["bad_color"]
	if !ok {
		t.Fatal("invalid-color entry should still produce a fallback config")
	}
	if cfg.Color == nil {
		t.Fatal("fallback Color is nil; expected white fallback")
	}
	gotR, gotG, gotB, gotA := cfg.Color.RGBA()
	wantR, wantG, wantB, wantA := color.White.RGBA()
	if gotR != wantR || gotG != wantG || gotB != wantB || gotA != wantA {
		t.Errorf("fallback Color RGBA = (%d,%d,%d,%d), want white", gotR, gotG, gotB, gotA)
	}
}

func TestVFXManager_PixelSizeClampedToOne(t *testing.T) {
	jsonBytes := []byte(`[
		{"type": "zero_size", "pixel": {"size": 0, "color": "#FFFFFF", "lifetime_frames": 4}}
	]`)
	fsys := fstest.MapFS{
		"vfx.json": &fstest.MapFile{Data: jsonBytes},
	}

	m := NewManager(fsys, "vfx.json")
	cfg, ok := m.configs["zero_size"]
	if !ok {
		t.Fatal("zero_size config not loaded")
	}
	if cfg.FrameWidth != 1 || cfg.FrameHeight != 1 {
		t.Errorf("size<=0 must clamp to 1; got FrameWidth=%d FrameHeight=%d",
			cfg.FrameWidth, cfg.FrameHeight)
	}
}

func TestVFXJSON_AllPixelEntriesUse1BitPalette(t *testing.T) {
	_, data := realVFXFS(t)

	var entries []VFXConfig
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal vfx.json: %v", err)
	}

	pixelCount := 0
	for _, e := range entries {
		if e.Pixel == nil {
			continue
		}
		pixelCount++
		switch e.Pixel.Color {
		case "#000000", "#FFFFFF":
			// ok
		default:
			t.Errorf("entry %q uses non-1bit color %q", e.Type, e.Pixel.Color)
		}
	}
	if pixelCount < 3 {
		t.Errorf("expected at least 3 pixel-mode entries (muzzle_flash, bullet_impact, bullet_despawn); found %d", pixelCount)
	}
}
