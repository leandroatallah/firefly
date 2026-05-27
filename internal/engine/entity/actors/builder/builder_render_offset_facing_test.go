package builder

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
)

// T-B1 (story 070): ApplyRenderOffsets forwards SpriteOffset {X,Y} into the
// character's render-offset registry. Per-facing X resolution is observable at
// draw time via Character.RenderOffset which auto-mirrors X when facing left.
// Covers AC-5, AC-10.
func TestApplyRenderOffsets_FacingAware(t *testing.T) {
	tests := []struct {
		name        string
		assetOffset schemas.SpriteOffset
		facing      animation.FacingDirectionEnum
		wantPt      image.Point
	}{
		{
			name:        "facing right uses X as-is",
			assetOffset: schemas.SpriteOffset{X: 10, Y: 2},
			facing:      animation.FaceDirectionRight,
			wantPt:      image.Pt(10, 2),
		},
		{
			name:        "facing left auto-mirrors X",
			assetOffset: schemas.SpriteOffset{X: 10, Y: 2},
			facing:      animation.FaceDirectionLeft,
			wantPt:      image.Pt(-10, 2),
		},
		{
			name:        "facing left auto-mirrors negative X",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2},
			facing:      animation.FaceDirectionLeft,
			wantPt:      image.Pt(4, 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := newMockActorWithCollision()
			character := actors.NewCharacter(sprites.SpriteMap{}, bodyphysics.NewRect(0, 0, 16, 16))
			actor.SetCharacter(character)

			offset := tt.assetOffset
			spriteData := schemas.SpriteData{
				Assets: map[string]schemas.AssetData{
					"idle": {Path: "i.png", RenderOffset: &offset},
				},
			}
			stateMap, err := BuildStateMap(spriteData)
			if err != nil {
				t.Fatalf("BuildStateMap: %v", err)
			}

			ApplyRenderOffsets(actor, spriteData, stateMap)

			character.SetAcceleration(0, 0)
			character.SetFaceDirection(tt.facing)

			got, ok := character.RenderOffset(actors.Idle)
			if !ok {
				t.Fatalf("RenderOffset(Idle) ok = false after ApplyRenderOffsets; want true")
			}
			if got != tt.wantPt {
				t.Errorf("RenderOffset(Idle) facing=%v = %v, want %v",
					tt.facing, got, tt.wantPt)
			}
		})
	}
}
