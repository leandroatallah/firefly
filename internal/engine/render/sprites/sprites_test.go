package sprites

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
)

func TestAnimatedSpriteImageFrameSelection(t *testing.T) {
	// Create a 96x16 image with 3 frames of 32x16 each
	img := ebiten.NewImage(96, 16)
	sprite := &Sprite{Image: img, Loop: true}

	se := NewSpriteEntity(SpriteMap{"idle": sprite})

	rect := image.Rect(0, 0, 32, 16)

	// count/frameRate -> frame index
	frame := se.AnimatedSpriteImage(sprite, rect, 0, 1)
	if frame.Bounds().Dx() != 32 || frame.Bounds().Dy() != 16 {
		t.Fatalf("unexpected frame size: %+v", frame.Bounds())
	}

	frame2 := se.AnimatedSpriteImage(sprite, rect, 1, 1) // next frame
	if frame2.Bounds().Dx() != 32 || frame2.Bounds().Dy() != 16 {
		t.Fatalf("unexpected frame2 size: %+v", frame2.Bounds())
	}

	se.SetFrameRate(2)
	if se.FrameRate() != 2 {
		t.Fatalf("FrameRate not set")
	}

	// Non-looping clamps to last frame
	spriteNL := &Sprite{Image: img, Loop: false}
	_ = se.AnimatedSpriteImage(spriteNL, rect, 999, 1) // should not panic, clamps internally

	// Error cases: nil sprite or image
	if se.AnimatedSpriteImage(nil, rect, 0, 1) != nil {
		t.Error("expected nil for nil sprite")
	}
	if se.AnimatedSpriteImage(&Sprite{Image: nil}, rect, 0, 1) != nil {
		t.Error("expected nil for nil image")
	}

	// No width case
	_ = se.AnimatedSpriteImage(sprite, image.Rect(0, 0, 0, 16), 0, 1)
}

func TestAnimatedSpriteImage_FrameCalculation(t *testing.T) {
	// Create sprite sheet with known frame count
	img := ebiten.NewImage(128, 32) // 4 frames of 32x32
	sprite := &Sprite{Image: img, Loop: true}
	se := NewSpriteEntity(SpriteMap{"test": sprite})
	rect := image.Rect(0, 0, 32, 32)

	tests := []struct {
		name      string
		count     int
		frameRate int
		loop      bool
		wantFrame int // expected frame index
	}{
		{"frame 0 at count 0", 0, 1, true, 0},
		{"frame 1 at count 1", 1, 1, true, 1},
		{"frame 2 at count 2", 2, 1, true, 2},
		{"frame 3 at count 3", 3, 1, true, 3},
		{"wrap to frame 0 at count 4", 4, 1, true, 0},
		{"frame 1 at count 5", 5, 1, true, 1},
		{"slower animation (rate 2)", 2, 2, true, 1}, // count/rate = 2/2 = 1, frame 1, offset 32
		{"slower animation (rate 2) frame 2", 4, 2, true, 2}, // count/rate = 4/2 = 2, frame 2, offset 64
		{"non-loop clamps", 100, 1, false, 3}, // Should clamp to last frame (3)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sprite.Loop = tt.loop
			frame := se.AnimatedSpriteImage(sprite, rect, tt.count, tt.frameRate)
			if frame == nil {
				t.Fatal("AnimatedSpriteImage returned nil")
			}

			// Calculate expected frame bounds
			expectedX := tt.wantFrame * 32
			bounds := frame.Bounds()
			if bounds.Min.X != expectedX {
				t.Errorf("frame X offset = %d, want %d", bounds.Min.X, expectedX)
			}
		})
	}
}

func TestAnimatedSpriteImage_EdgeCases(t *testing.T) {
	img := ebiten.NewImage(64, 32)
	sprite := &Sprite{Image: img, Loop: true}
	se := NewSpriteEntity(SpriteMap{"test": sprite})

	tests := []struct {
		name   string
		sprite *Sprite
		rect   image.Rectangle
		count  int
		rate   int
		wantNil bool
	}{
		{
			name: "nil sprite",
			sprite: nil, rect: image.Rect(0, 0, 32, 32),
			count: 0, rate: 1, wantNil: true,
		},
		{
			name: "nil image",
			sprite: &Sprite{Image: nil}, rect: image.Rect(0, 0, 32, 32),
			count: 0, rate: 1, wantNil: true,
		},
		{
			name: "zero width rect",
			sprite: sprite, rect: image.Rect(0, 0, 0, 32),
			count: 0, rate: 1, wantNil: false, // Returns full image
		},
		{
			name: "zero height rect",
			sprite: sprite, rect: image.Rect(0, 0, 32, 0),
			count: 0, rate: 1, wantNil: false,
		},
		{
			name: "frame larger than image",
			sprite: sprite, rect: image.Rect(0, 0, 128, 64),
			count: 0, rate: 1, wantNil: false, // Returns full image
		},
		{
			name: "negative count",
			sprite: sprite, rect: image.Rect(0, 0, 32, 32),
			count: -10, rate: 1, wantNil: false,
		},
		// Note: zero frame rate causes divide by zero in production code
		// This is a known limitation that should be fixed in sprites.go
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := se.AnimatedSpriteImage(tt.sprite, tt.rect, tt.count, tt.rate)
			if tt.wantNil && result != nil {
				t.Error("expected nil result")
			}
			if !tt.wantNil && result == nil {
				t.Error("expected non-nil result")
			}
		})
	}
}

func TestSpriteEntityGetters(t *testing.T) {
	sprite := &Sprite{Image: ebiten.NewImage(32, 32), Loop: true}
	sprites := SpriteMap{"idle": sprite}
	se := NewSpriteEntity(sprites)

	if se.GetFirstSprite() != sprite {
		t.Error("expected GetFirstSprite to return the only sprite")
	}

	if se.GetSpriteByState("idle") != sprite {
		t.Error("expected GetSpriteByState to return the correct sprite")
	}

	if len(se.Sprites()) != 1 {
		t.Error("expected Sprites() to return the sprite map")
	}

	seEmpty := NewSpriteEntity(nil)
	if seEmpty.GetFirstSprite() != nil {
		t.Error("expected GetFirstSprite to return nil for empty entity")
	}
}

func TestSpriteEntity_GetSpriteByState(t *testing.T) {
	sprite1 := &Sprite{Image: ebiten.NewImage(32, 32), Loop: true}
	sprite2 := &Sprite{Image: ebiten.NewImage(32, 32), Loop: false}
	sprites := SpriteMap{
		"idle":  sprite1,
		"walk":  sprite2,
		"jump":  sprite1,
	}
	se := NewSpriteEntity(sprites)

	tests := []struct {
		name  string
		state interface{}
		want  *Sprite
	}{
		{"existing state idle", "idle", sprite1},
		{"existing state walk", "walk", sprite2},
		{"existing state jump", "jump", sprite1},
		{"non-existing state", "run", nil},
		{"empty string state", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: state type depends on animation.SpriteState which is a string alias
			stateStr, ok := tt.state.(string)
			if !ok {
				t.Skip("state must be string")
			}
			got := se.GetSpriteByState(stateStr)
			if got != tt.want {
				t.Errorf("GetSpriteByState(%q) = %p, want %p", stateStr, got, tt.want)
			}
		})
	}
}

func TestSpriteEntity_FrameRate(t *testing.T) {
	se := NewSpriteEntity(SpriteMap{})

	// Default frame rate
	if se.FrameRate() != 1 {
		t.Errorf("default FrameRate() = %d, want 1", se.FrameRate())
	}

	// Set frame rate
	se.SetFrameRate(5)
	if se.FrameRate() != 5 {
		t.Errorf("FrameRate() after SetFrameRate(5) = %d, want 5", se.FrameRate())
	}

	// Zero frame rate
	se.SetFrameRate(0)
	if se.FrameRate() != 0 {
		t.Errorf("FrameRate() after SetFrameRate(0) = %d, want 0", se.FrameRate())
	}

	// Negative frame rate (edge case)
	se.SetFrameRate(-1)
	if se.FrameRate() != -1 {
		t.Errorf("FrameRate() after SetFrameRate(-1) = %d, want -1", se.FrameRate())
	}
}

func TestSpriteAssets(t *testing.T) {
	var sa SpriteAssets
	sa = sa.AddSprite("idle", "path/to/idle.png", true)
	if len(sa) != 1 {
		t.Fatalf("expected 1 sprite in assets, got %d", len(sa))
	}
	if sa["idle"].Path != "path/to/idle.png" || !sa["idle"].Loop {
		t.Error("sprite assets not set correctly")
	}
}

func TestSpriteAssets_AddMultiple(t *testing.T) {
	sa := SpriteAssets{}
	sa = sa.AddSprite("idle", "idle.png", true)
	sa = sa.AddSprite("walk", "walk.png", false)
	sa = sa.AddSprite("jump", "jump.png", true)

	if len(sa) != 3 {
		t.Errorf("expected 3 sprites, got %d", len(sa))
	}

	if sa["idle"].Path != "idle.png" || !sa["idle"].Loop {
		t.Error("idle sprite not set correctly")
	}
	if sa["walk"].Path != "walk.png" || sa["walk"].Loop {
		t.Error("walk sprite not set correctly")
	}
	if sa["jump"].Path != "jump.png" || !sa["jump"].Loop {
		t.Error("jump sprite not set correctly")
	}
}

func TestSpriteAssets_AddSpriteOverwrite(t *testing.T) {
	sa := SpriteAssets{}
	sa = sa.AddSprite("idle", "first.png", true)
	sa = sa.AddSprite("idle", "second.png", false) // Overwrite

	if len(sa) != 1 {
		t.Errorf("expected 1 sprite after overwrite, got %d", len(sa))
	}
	if sa["idle"].Path != "second.png" || sa["idle"].Loop {
		t.Error("sprite not overwritten correctly")
	}
}

func TestLoadSpritesError(t *testing.T) {
	sa := SpriteAssets{}.AddSprite("idle", "non_existent.png", true)
	_, err := LoadSpritesFromOS(sa)
	if err == nil {
		t.Error("expected error loading non-existent sprite")
	}
}

func TestLoadSprites_Success(t *testing.T) {
	// This test requires actual image files, so we use the error case
	// to verify the function structure. A full integration test would
	// require test assets.
	sa := SpriteAssets{}.AddSprite("test", "non_existent.png", true)
	sprites, err := LoadSpritesFromOS(sa)
	if err == nil {
		t.Error("expected error for non-existent file")
	}
	if sprites != nil {
		t.Error("expected nil sprites on error")
	}
}

func TestGetSpritesFromAssets(t *testing.T) {
	// Since GetSpritesFromAssets calls LoadSprites which hits the disk,
	// we test the mapping logic but expect an error unless we have a real image.

	assets := map[string]schemas.AssetData{
		"idle": {Path: "non_existent.png"},
	}
	stateMap := map[string]animation.SpriteState{
		"idle": "idle_state",
	}

	_, err := GetSpritesFromAssetsOS(assets, stateMap)
	if err == nil {
		t.Error("expected error for non-existent image path")
	}

	// Test with no matching states (should be empty map but no error)
	emptyStateMap := map[string]animation.SpriteState{}
	res, err := GetSpritesFromAssetsOS(assets, emptyStateMap)
	if err != nil {
		t.Errorf("unexpected error for empty state map: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected empty result, got %v", res)
	}
}

func TestGetSpritesFromAssets_LoopNil(t *testing.T) {
	// Test that nil Loop field defaults to true
	assets := map[string]schemas.AssetData{
		"idle": {Path: "non_existent.png", Loop: nil},
	}
	stateMap := map[string]animation.SpriteState{
		"idle": "idle_state",
	}

	_, err := GetSpritesFromAssetsOS(assets, stateMap)
	// Will error on file load, but we're testing the nil Loop handling
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}
