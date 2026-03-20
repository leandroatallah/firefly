package jsonutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
)

type TestStats struct {
	Health int `json:"health"`
	Speed  int `json:"speed"`
}

func TestParseSpriteAndStats(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	jsonContent := `{
		"sprites": {
			"body_rect": {"x": 0, "y": 0, "width": 32, "height": 32},
			"assets": {
				"idle": {
					"path": "idle.png",
					"collision_rect": []
				}
			},
			"frame_rate": 8,
			"facing_direction": 0
		},
		"stats": {
			"health": 100,
			"speed": 10
		}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	spriteData, stats, err := ParseSpriteAndStats[TestStats](os.DirFS(tmpDir), "test.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Spot-check key fields to verify parsing worked
	if spriteData.FrameRate != 8 {
		t.Errorf("expected FrameRate 8, got %d", spriteData.FrameRate)
	}
	if spriteData.FacingDirection != animation.FaceDirectionRight {
		t.Errorf("expected FacingDirection right, got %d", spriteData.FacingDirection)
	}
	if stats.Health != 100 {
		t.Errorf("expected Health 100, got %d", stats.Health)
	}
}

func TestParseSpriteAndStats_FileNotFound(t *testing.T) {
	_, _, err := ParseSpriteAndStats[TestStats](os.DirFS("/"), "nonexistent/file.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestParseSpriteAndStats_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(testFile, []byte("not valid json"), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, _, err = ParseSpriteAndStats[TestStats](os.DirFS(tmpDir), "invalid.json")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseSpriteAndStats_MissingFields(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "missing.json")

	jsonContent := `{
		"sprites": {
			"body_rect": {"x": 0, "y": 0, "width": 32, "height": 32}
		}
	}`

	err := os.WriteFile(testFile, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, stats, err := ParseSpriteAndStats[TestStats](os.DirFS(tmpDir), "missing.json")
	if err != nil {
		t.Fatalf("expected no error (missing fields should use defaults), got %v", err)
	}

	// Assets will be nil when not provided in JSON - this is expected
	if stats.Health != 0 {
		t.Errorf("expected Health 0 (default), got %d", stats.Health)
	}
}
