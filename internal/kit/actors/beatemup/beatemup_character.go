package beatemup

import (
	"image"
	"io/fs"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	kitactors "github.com/boilerplate/ebiten-template/internal/kit/actors"
)

type BeatEmUpCharacter struct {
	*actors.Character
	*kitactors.MeleeCharacter
	footprints map[actors.ActorStateEnum]image.Rectangle // local rect, NOT world-offset
}

// buildFootprints constructs a per-state footprint map from asset data.
// Only assets with a non-nil FootprintRect of positive size and a matching
// stateMap entry are included. The stored rectangles are in local (body-relative)
// coordinates; world offset is applied at call time in Footprint/CollisionPosition.
func buildFootprints(
	assets map[string]schemas.AssetData,
	stateMap map[string]animation.SpriteState,
) map[actors.ActorStateEnum]image.Rectangle {
	out := make(map[actors.ActorStateEnum]image.Rectangle)
	for key, asset := range assets {
		if asset.FootprintRect == nil {
			continue
		}
		st, ok := stateMap[key]
		if !ok {
			continue
		}
		enumSt, ok := st.(actors.ActorStateEnum)
		if !ok {
			continue
		}
		r := asset.FootprintRect
		if r.Width <= 0 || r.Height <= 0 {
			continue
		}
		out[enumSt] = image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
	}
	return out
}

func NewBeatEmUpCharacter(
	fsys fs.FS,
	stateMap map[string]animation.SpriteState,
	spriteData schemas.SpriteData,
	bodyRect *bodyphysics.Rect,
	blocker physicsmovement.PlayerMovementBlocker,
) (*BeatEmUpCharacter, error) {
	s, err := sprites.GetSpritesFromAssets(fsys, spriteData.Assets, stateMap)
	if err != nil {
		return nil, err
	}
	c := actors.NewCharacter(s, bodyRect)
	be := &BeatEmUpCharacter{
		Character:      c,
		MeleeCharacter: kitactors.NewMeleeCharacter(),
	}
	be.footprints = buildFootprints(spriteData.Assets, stateMap)
	c.SetMovementModel(physicsmovement.NewBeatEmUpMovementModel(blocker))
	c.SetFaceDirection(spriteData.FacingDirection)
	c.SetFrameRate(spriteData.FrameRate)
	c.SetOwner(be)
	c.SetMovementTransitionHandler(beatemupMovementTransitions)
	return be, nil
}

// Footprint returns the current state's footprint rectangle in world coordinates.
// Falls back to the union of the actor's collision rects when no footprint is
// declared for the current state. If no collision rects exist either, returns
// the body Position().
func (c *BeatEmUpCharacter) Footprint() image.Rectangle {
	st := c.State()
	if local, ok := c.footprints[st]; ok {
		minX, minY := c.GetPositionMin()
		return local.Add(image.Pt(minX, minY))
	}
	// Fallback: union of collision rects via explicit embedded selector.
	rects := c.CollidableBody.CollisionPosition()
	if len(rects) == 0 {
		return c.Position()
	}
	u := rects[0]
	for i := 1; i < len(rects); i++ {
		u = u.Union(rects[i])
	}
	return u
}

// CollisionPosition shadows the embedded *CollidableBody method.
// When a footprint exists for the current state, it returns only the footprint
// so movement-vs-world and actor-vs-actor checks use the feet area.
// When absent, returns the embedded full collision rects.
func (c *BeatEmUpCharacter) CollisionPosition() []image.Rectangle {
	st := c.State()
	if local, ok := c.footprints[st]; ok {
		minX, minY := c.GetPositionMin()
		return []image.Rectangle{local.Add(image.Pt(minX, minY))}
	}
	return c.CollidableBody.CollisionPosition()
}

// GroundY implements space.DepthLaneBody.
// Returns the body's pre-altitude ground Y (the depth position on the floor
// plane). Altitude is kept separate and must not be subtracted here — depth
// gating compares floor positions, not screen positions.
func (c *BeatEmUpCharacter) GroundY() int {
	_, y16 := c.GetPosition16()
	return y16 / 16 // == fp16.From16(y16); pre-altitude depth Y
}

// LaneHalfWidth implements space.DepthLaneBody.
// Returns space.DefaultLaneHalfWidth (8 px) — a fixed tolerance that covers
// the character's footprint depth extent for standard beat-em-up lane matching.
func (c *BeatEmUpCharacter) LaneHalfWidth() int {
	return space.DefaultLaneHalfWidth
}

func beatemupMovementTransitions(c *actors.Character) {
	vx, vy := c.Velocity()
	threshold := config.Get().Physics.DownwardGravity
	isMovingGround := vx > threshold || vx < -threshold || vy > threshold || vy < -threshold
	state := c.State()
	vAlt16 := c.VAltitude16()
	altitude := c.Altitude()
	airborne := altitude > 0
	set := func(s actors.ActorStateEnum) { c.SetNewStateFatal(s) }
	switch {
	case state == actors.Landing:
		if !c.IsAnimationFinished() {
			return
		}
		if isMovingGround {
			set(actors.Walking)
		} else {
			set(actors.Idle)
		}
		return
	case state == actors.Jumping && vAlt16 >= 0 && airborne:
		set(actors.Falling)
		return
	case vAlt16 < 0:
		if state != actors.Jumping {
			set(actors.Jumping)
		}
		return
	case vAlt16 > 0 && airborne:
		if state != actors.Falling {
			set(actors.Falling)
		}
		return
	case state == actors.Falling && !airborne:
		set(actors.Landing)
		return
	case airborne:
		return
	case state != actors.Walking && isMovingGround:
		set(actors.Walking)
	case state != actors.Idle && !isMovingGround:
		set(actors.Idle)
	}
}
