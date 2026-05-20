package beatemup

import (
	"io/fs"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	kitactors "github.com/boilerplate/ebiten-template/internal/kit/actors"
)

type BeatEmUpCharacter struct {
	*actors.Character
	*kitactors.MeleeCharacter
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
	c.SetMovementModel(physicsmovement.NewBeatEmUpMovementModel(blocker))
	c.SetFaceDirection(spriteData.FacingDirection)
	c.SetFrameRate(spriteData.FrameRate)
	c.SetOwner(be)
	c.SetMovementTransitionHandler(beatemupMovementTransitions)
	return be, nil
}

func beatemupMovementTransitions(c *actors.Character) {
	vx, vy := c.Velocity()
	threshold := config.Get().Physics.DownwardGravity
	isMoving := vx > threshold || vx < -threshold || vy > threshold || vy < -threshold
	state := c.State()
	set := func(s actors.ActorStateEnum) { c.SetNewStateFatal(s) }
	switch {
	case state != actors.Walking && isMoving:
		set(actors.Walking)
	case state != actors.Idle && !isMoving:
		set(actors.Idle)
	}
}
