package actors

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

type Character struct {
	sprites.SpriteEntity

	*bodyphysics.MovableBody
	*bodyphysics.CollidableBody
	*bodyphysics.AliveBody
	*space.StateCollisionManager[ActorStateEnum]

	Touchable body.Touchable

	count                int
	state                ActorState
	movementState        movement.MovementState
	movementModel        physicsmovement.MovementModel
	movementBlockers     int
	invulnerabilityTimer int
	imageOptions         *ebiten.DrawImageOptions

	skills []skill.Skill

	// Original sprite frame size for correct rendering when body size changes
	// TODO: Turns it on a behavioral struct. It is too specific.
	originalFrameWidth  int
	originalFrameHeight int

	StateTransitionHandler func(*Character) bool
	OnStateChange          func(oldState, newState ActorStateEnum)
	bodyphysics.Ownership
}

// SetStateTransitionHandler sets a function that can override the default state transition logic.
func (c *Character) SetStateTransitionHandler(handler func(*Character) bool) {
	c.StateTransitionHandler = handler
}

func NewCharacter(s sprites.SpriteMap, bodyRect *bodyphysics.Rect) *Character { // Modified signature
	spriteEntity := sprites.NewSpriteEntity(s)
	b := bodyphysics.NewBody(bodyRect)
	movable := bodyphysics.NewMovableBody(b)
	collidable := bodyphysics.NewCollidableBody(b)
	alive := bodyphysics.NewAliveBody(b)
	c := &Character{
		MovableBody:    movable,
		CollidableBody: collidable,
		AliveBody:      alive,

		SpriteEntity:        spriteEntity,
		imageOptions:        &ebiten.DrawImageOptions{},
		originalFrameWidth:  bodyRect.Width(),
		originalFrameHeight: bodyRect.Height(),
	}
	// Set the owner for all body components to this Character
	// Body.Owner -> MovableBody (chosen as the primary physical representation)
	b.SetOwner(movable)
	movable.SetOwner(c)
	collidable.SetOwner(c)
	alive.SetOwner(c)

	c.StateCollisionManager = space.NewStateCollisionManager[ActorStateEnum](c)

	state, err := c.NewState(Idle)
	if err != nil {
		log.Fatal(err)
	}
	c.SetState(state)
	return c
}

// Forwarding methods for Body to avoid ambiguous selector
// Always route via the MovableBody component
func (c *Character) ID() string {
	return c.MovableBody.ID()
}
func (c *Character) SetID(id string) {
	c.MovableBody.SetID(id)
}
func (c *Character) Position() image.Rectangle {
	return c.MovableBody.Position()
}
func (c *Character) SetPosition(x, y int) {
	c.CollidableBody.SetPosition(x, y)
}
func (c *Character) SetPosition16(x16, y16 int) {
	c.CollidableBody.SetPosition16(x16, y16)
}
func (c *Character) SetSize(width, height int) {
	c.MovableBody.SetSize(width, height)
}
func (c *Character) Scale() float64 {
	return c.MovableBody.Scale()
}
func (c *Character) SetScale(scale float64) {
	c.MovableBody.SetScale(scale)
}
func (c *Character) GetPosition16() (int, int) {
	return c.MovableBody.GetPosition16()
}
func (c *Character) GetPositionMin() (int, int) {
	return c.MovableBody.GetPositionMin()
}
func (c *Character) GetShape() body.Shape {
	return c.MovableBody.GetShape()
}

// Builder methods
func (c *Character) State() ActorStateEnum {
	return c.state.State()
}

func (c *Character) IsAnimationFinished() bool {
	if c.state == nil {
		return true
	}
	return c.state.IsAnimationFinished()
}

func (c *Character) AddCollisionRect(state ActorStateEnum, rect body.Collidable) {
	c.StateCollisionManager.AddCollisionRect(state, rect)
}

func (c *Character) GetCharacter() *Character {
	return c
}

func (c *Character) NewState(state ActorStateEnum) (ActorState, error) {
	return NewState(c, state)
}

func (c *Character) SetNewStateFatal(state ActorStateEnum) {
	s, err := NewState(c, state)
	if err != nil {
		log.Fatalf("Failed to create new state %v: %v", s, err)
	}
	c.SetState(s)
}

// SetState set a new Character state and update current collision shapes.
func (c *Character) SetState(state ActorState) {
	if c.state == nil || c.state.State() != state.State() {
		var oldState ActorStateEnum
		if c.state != nil {
			oldState = c.state.State()
		} else {
			oldState = -1 // Unknown/First state
		}

		c.state = state
		c.state.OnStart(c.count)
		c.StateCollisionManager.RefreshCollisions()

		if c.OnStateChange != nil {
			c.OnStateChange(oldState, c.state.State())
		}
	}
}

func (c *Character) SetMovementState(
	state movement.MovementStateEnum,
	target body.MovableCollidable,
	options ...movement.MovementStateOption,
) {
	movementState, err := movement.NewMovementState(c, state, target, options...)
	if err != nil {
		log.Fatal(err)
	}

	c.movementState = movementState
	c.movementState.OnStart()
}
func (c *Character) SwitchMovementState(state movement.MovementStateEnum) {
	target := c.MovementState().Target()
	movementState, err := movement.NewMovementState(c, state, target)
	if err != nil {
		log.Fatal(err)
	}
	c.movementState = movementState
}

func (c *Character) MovementState() movement.MovementState {
	return c.movementState
}

func (c *Character) Update(space body.BodiesSpace) error {
	if c.Freeze() {
		return nil
	}
	c.count++

	for _, s := range c.skills {
		if activeSkill, ok := s.(skill.ActiveSkill); ok {
			activeSkill.HandleInput(c, c.movementModel.(*physicsmovement.PlatformMovementModel), space)
		}
		s.Update(c, c.movementModel.(*physicsmovement.PlatformMovementModel))
	}

	// Handle movement by Movement State - must happen BEFORE UpdateMovement
	if c.movementState != nil {
		c.movementState.Move(space)
	}

	// Update physics and apply movement
	c.UpdateMovement(space)

	c.handleState()

	return nil
}

func (c *Character) UpdateMovement(space body.BodiesSpace) {
	if c.movementModel != nil {
		c.movementModel.Update(c, space)
	}
}

func (c *Character) UpdateImageOptions() {
	if c.imageOptions == nil {
		return
	}
	c.imageOptions.GeoM.Reset()

	if s := c.Scale(); s != 0 && s != 1.0 {
		c.imageOptions.GeoM.Scale(s, s)
	}

	accX, _ := c.Acceleration()
	fDirection := c.FaceDirection()

	if accX > 0 {
		fDirection = animation.FaceDirectionRight
	} else if accX < 0 {
		fDirection = animation.FaceDirectionLeft
	}

	c.SetFaceDirection(fDirection)

	if fDirection == animation.FaceDirectionLeft {
		width := c.Position().Dx()
		c.imageOptions.GeoM.Scale(-1, 1)
		c.imageOptions.GeoM.Translate(float64(width), 0)
	}

	// Apply character position
	x, y := c.GetPositionMin()
	c.imageOptions.GeoM.Translate(
		float64(x),
		float64(y),
	)
}

func (c *Character) handleState() {
	if c.state == nil {
		return
	}

	// Handle invulnerability timer
	if c.invulnerabilityTimer > 0 {
		c.invulnerabilityTimer--
		if c.invulnerabilityTimer == 0 {
			c.SetInvulnerability(false)
		}
	}

	// Allow game-specific logic to override the default behavior
	if c.StateTransitionHandler != nil && c.StateTransitionHandler(c) {
		return
	}

	setNewState := func(s ActorStateEnum) {
		state, err := c.NewState(s)
		if err != nil {
			log.Fatal(err)
		}
		c.SetState(state)
	}

	state := c.state.State()

	switch {
	case state == Hurted:
		isAnimationOver := c.state.IsAnimationFinished()
		if isAnimationOver {
			setNewState(Idle)
		}
	case state == Landing:
		isAnimationOver := c.state.IsAnimationFinished()
		if c.IsWalking() {
			setNewState(Walking)
		} else if isAnimationOver {
			setNewState(Idle)
		}
	case state == Jumping:
		isAnimationOver := c.state.IsAnimationFinished()
		if isAnimationOver {
			setNewState(Idle)
		}
	case state == Falling && !c.IsFalling():
		setNewState(Landing)
	case c.IsGoingUp():
		setNewState(Jumping)
	case state != Falling && c.IsFalling():
		setNewState(Falling)
	case state != Walking && c.IsWalking():
		setNewState(Walking)
	case state != Idle && c.IsIdle():
		setNewState(Idle)
	}
}

func (c *Character) Hurt(damage int) {
	if c.Invulnerable() {
		return
	}

	c.LoseHealth(damage)

	// Switch to Hurt state
	state, err := c.NewState(Hurted)
	if err != nil {
		log.Fatal(err)
	}
	c.SetState(state)
	c.SetInvulnerability(true)
	c.invulnerabilityTimer = 120 // 2 seconds at 60fps
}

func (c *Character) SetTouchable(t body.Touchable) {
	c.Touchable = t
}

func (c *Character) Image() *ebiten.Image {
	sprite := c.GetSpriteByState(c.state.State())
	if sprite == nil || sprite.Image == nil {
		// Try to fallback to idle sprite
		sprite = c.GetSpriteByState(Idle)
	}
	if sprite == nil || sprite.Image == nil {
		sprite = c.GetFirstSprite()
	}

	// Use original frame size for sprite extraction, not current body size
	// This ensures correct sprite rendering when body size changes (e.g., grow skill)
	pos := c.Position()
	frameRect := image.Rect(pos.Min.X, pos.Min.Y, pos.Min.X+c.originalFrameWidth, pos.Min.Y+c.originalFrameHeight)
	stateDurationCount := c.state.GetAnimationCount(c.count)
	return c.AnimatedSpriteImage(sprite, frameRect, stateDurationCount, c.SpriteEntity.FrameRate())
}

func (c *Character) ImageOptions() *ebiten.DrawImageOptions {
	return c.imageOptions
}

// BlockMovement increases the count of systems blocking movement.
func (p *Character) BlockMovement() {
	p.movementBlockers++
}

// UnblockMovement decreases the count.
func (p *Character) UnblockMovement() {
	p.movementBlockers--
	if p.movementBlockers < 0 {
		p.movementBlockers = 0
	}
}

// IsPlayerMovementBlocked checks if any system is currently blocking movement.
func (p *Character) IsMovementBlocked() bool {
	return p.movementBlockers > 0
}

// Movement Model methods
func (c *Character) SetMovementModel(model physicsmovement.MovementModel) {
	c.movementModel = model
}

func (c *Character) MovementModel() physicsmovement.MovementModel {
	return c.movementModel
}

func (c *Character) AddSkill(s skill.Skill) {
	c.skills = append(c.skills, s)
}

func (c *Character) RemoveSkill(s skill.Skill) {
	for i, skill := range c.skills {
		if skill == s {
			c.skills = append(c.skills[:i], c.skills[i+1:]...)
			return
		}
	}
}

func (c *Character) ClearSkills() {
	c.skills = nil
}
