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

// intPtr is a small helper so tests can build *int values for XFlipped.
func intPtr(v int) *int { return &v }

// T-B1 (story 070): ApplyRenderOffsets forwards the XFlipped pointer from
// schemas.SpriteOffset into the character's render-offset registry. The
// per-facing resolution is observable at draw time via Character.RenderOffset
// which returns the X resolved against the current facing direction.
// Covers AC-5, AC-10.
func TestApplyRenderOffsets_ForwardsXFlippedPerFacing(t *testing.T) {
	tests := []struct {
		name        string
		assetOffset schemas.SpriteOffset
		facing      animation.FacingDirectionEnum
		wantPt      image.Point
	}{
		{
			name:        "XFlipped set, facing right -> uses X",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2, XFlipped: intPtr(6)},
			facing:      animation.FaceDirectionRight,
			wantPt:      image.Pt(-4, 2),
		},
		{
			name:        "XFlipped set, facing left -> uses XFlipped",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2, XFlipped: intPtr(6)},
			facing:      animation.FaceDirectionLeft,
			wantPt:      image.Pt(6, 2),
		},
		{
			name:        "XFlipped nil, facing right -> uses X",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2},
			facing:      animation.FaceDirectionRight,
			wantPt:      image.Pt(-4, 2),
		},
		{
			name:        "XFlipped nil, facing left -> falls back to X (068 regression)",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2},
			facing:      animation.FaceDirectionLeft,
			wantPt:      image.Pt(-4, 2),
		},
		{
			name:        "XFlipped=0 explicit, facing left -> 0",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2, XFlipped: intPtr(0)},
			facing:      animation.FaceDirectionLeft,
			wantPt:      image.Pt(0, 2),
		},
		{
			name:        "XFlipped=0 explicit, facing right -> uses X",
			assetOffset: schemas.SpriteOffset{X: -4, Y: 2, XFlipped: intPtr(0)},
			facing:      animation.FaceDirectionRight,
			wantPt:      image.Pt(-4, 2),
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

			// Drive facing direction on the character itself; RenderOffset
			// resolves X against the current facing at call time.
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
