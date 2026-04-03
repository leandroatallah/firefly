package gamestates

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/tween"
)

var StateDashing actors.ActorStateEnum

func init() {
	StateDashing = actors.RegisterState("dash", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b} // placeholder; DashState is constructed directly
	})
}

// DashConfig holds configuration for the dash state.
type DashConfig struct {
	Speed          int
	DurationFrames int
	BlockDistance  int
	Cooldown       int
}

// DashState drives a tween-based dash deceleration.
type DashState struct {
	body        contractsbody.MovableCollidable
	space       contractsbody.BodiesSpace
	cfg         DashConfig
	tween       *tween.InOutSineTween
	dashing     bool
	airDashUsed bool
}

func NewDashState(body contractsbody.MovableCollidable, space contractsbody.BodiesSpace, cfg DashConfig) *DashState {
	return &DashState{body: body, space: space, cfg: cfg}
}

// OnStart begins the dash. If already dashing, it is a no-op.
func (s *DashState) OnStart(_ int) {
	if s.dashing {
		return
	}

	// Block-distance check: query in facing direction.
	pos := s.body.Position()
	var queryRect image.Rectangle
	if s.body.FaceDirection() == animation.FaceDirectionRight {
		queryRect = image.Rect(pos.Max.X, pos.Min.Y, pos.Max.X+s.cfg.BlockDistance, pos.Max.Y)
	} else {
		queryRect = image.Rect(pos.Min.X-s.cfg.BlockDistance, pos.Min.Y, pos.Min.X, pos.Max.Y)
	}
	if len(s.space.Query(queryRect)) > 0 {
		return
	}

	s.body.SetFreeze(true)
	s.body.SetVelocity(0, 0)

	if s.body.IsFalling() {
		s.airDashUsed = true
	}

	s.tween = tween.NewInOutSineTween(float64(s.cfg.Speed), 0, s.cfg.DurationFrames)
	s.dashing = true
}

// Update advances the dash by one frame and returns the next state.
func (s *DashState) Update() actors.ActorStateEnum {
	if !s.dashing {
		return s.nextIdleOrFalling()
	}

	// Check for wall collision via ApplyValidPosition.
	_, _, wallBlocked := s.body.ApplyValidPosition(0, true, s.space)
	if wallBlocked {
		s.finish()
		return s.nextIdleOrFalling()
	}

	speed := s.tween.Tick()
	dir := 1
	if s.body.FaceDirection() == animation.FaceDirectionLeft {
		dir = -1
	}
	s.body.SetVelocity(int(speed)*dir, 0)

	if s.tween.Done() {
		s.finish()
		return s.nextIdleOrFalling()
	}

	return StateDashing
}

// OnFinish ensures cleanup is done.
func (s *DashState) OnFinish() {
	s.finish()
}

func (s *DashState) finish() {
	if !s.dashing {
		return
	}
	s.dashing = false
	s.body.SetFreeze(false)
}

func (s *DashState) nextIdleOrFalling() actors.ActorStateEnum {
	if s.body.IsFalling() {
		return actors.Falling
	}
	return actors.Idle
}
