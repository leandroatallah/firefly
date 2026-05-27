package platformer_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
	"testing/fstest"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/kit/actors/platformer"
)

// tinyPlatformerPNG returns the bytes of a 1x1 transparent PNG so the
// sprite loader can decode a real image without an on-disk asset.
func tinyPlatformerPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0, 0, 0, 0})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

// newPlatformerRenderOffsetFixtures builds the minimal inputs to
// NewPlatformerCharacter with a single "idle" asset wired to actors.Idle.
// Callers can attach a RenderOffset to the asset before constructing.
func newPlatformerRenderOffsetFixtures(t *testing.T) (fstest.MapFS, map[string]animation.SpriteState, schemas.SpriteData, *bodyphysics.Rect) {
	t.Helper()
	pngBytes := tinyPlatformerPNG(t)
	fsys := fstest.MapFS{
		"idle.png": &fstest.MapFile{Data: pngBytes},
	}
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
	}
	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {Path: "idle.png"},
		},
		FrameRate:       1,
		FacingDirection: animation.FaceDirectionRight,
	}
	bodyRect := bodyphysics.NewRect(0, 0, 8, 8)
	return fsys, stateMap, spriteData, bodyRect
}

// T-P1 (story 070): NewPlatformerCharacter wires builder.ApplyRenderOffsets so
// that a SpriteData carrying render_offset registers an offset on the returned
// PlatformerCharacter. Today the platformer kit silently no-ops these JSON
// declarations — AC-6 closes that gap so the schema is honoured uniformly
// across kit genres.
func TestNewPlatformerCharacter_AppliesRenderOffsetsFromSpriteData(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newPlatformerRenderOffsetFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.RenderOffset = &schemas.SpriteOffset{X: -2, Y: 0}
	spriteData.Assets["idle"] = asset

	pf, err := platformer.NewPlatformerCharacter(fsys, stateMap, spriteData, bodyRect)
	if err != nil {
		t.Fatalf("NewPlatformerCharacter: %v", err)
	}
	if pf == nil {
		t.Fatal("NewPlatformerCharacter returned nil PlatformerCharacter")
	}

	got, ok := pf.RenderOffset(actors.Idle)
	if !ok {
		t.Fatalf("RenderOffset(Idle) ok = false; want true after construction with render_offset {x:-2, y:0} (AC-6 wiring missing)")
	}
	if got != image.Pt(-2, 0) {
		t.Errorf("RenderOffset(Idle) = %v, want (-2, 0)", got)
	}
}

// T-P1b (story 070): Render offset survives the platformer constructor
// pipeline and auto-mirrors X when the embedded Character faces left.
// Guards against a future refactor that drops facing-aware resolution.
func TestNewPlatformerCharacter_AppliesRenderOffsetsAutoMirrorsLeft(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newPlatformerRenderOffsetFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.RenderOffset = &schemas.SpriteOffset{X: -2, Y: 0}
	spriteData.Assets["idle"] = asset

	pf, err := platformer.NewPlatformerCharacter(fsys, stateMap, spriteData, bodyRect)
	if err != nil {
		t.Fatalf("NewPlatformerCharacter: %v", err)
	}

	pf.SetAcceleration(0, 0)
	pf.SetFaceDirection(animation.FaceDirectionLeft)

	got, ok := pf.RenderOffset(actors.Idle)
	if !ok {
		t.Fatalf("RenderOffset(Idle) ok = false; want true")
	}
	if got != image.Pt(2, 0) {
		t.Errorf("RenderOffset(Idle) facing-left = %v, want (2, 0) (auto-mirror)", got)
	}
}

// T-P2 (story 070): NewPlatformerCharacter without render_offset is a no-op
// — RenderOffset returns ok=false for every registered state. Pins AC-7:
// existing platformer JSONs that omit render_offset must continue to render
// identically.
func TestNewPlatformerCharacter_NoRenderOffsetIsNoOp(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newPlatformerRenderOffsetFixtures(t)
	// Intentionally leave RenderOffset nil on every asset.

	pf, err := platformer.NewPlatformerCharacter(fsys, stateMap, spriteData, bodyRect)
	if err != nil {
		t.Fatalf("NewPlatformerCharacter: %v", err)
	}

	for key, st := range stateMap {
		enum, ok := st.(actors.ActorStateEnum)
		if !ok {
			t.Fatalf("stateMap[%q] is not ActorStateEnum (got %T)", key, st)
		}
		if got, ok := pf.RenderOffset(enum); ok {
			t.Errorf("RenderOffset(%v) = (%v, ok=true) for asset %q without render_offset; want ok=false",
				enum, got, key)
		}
	}
}
