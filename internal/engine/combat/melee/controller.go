package melee

import (
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// animDurationFn computes how many game-frames the animation for the given
// combo step should last. Provided at construction so sprite-sheet logic stays
// in the game layer.
type animDurationFn func(stepIdx int) int

// Controller bundles all per-actor melee state: weapon, attack state, input
// bookkeeping. Install it on a Character to make any actor melee-capable.
type Controller struct {
	weapon          *weapon.MeleeWeapon
	state           *State
	meleeAttackEnum actors.ActorStateEnum
	stepStates      []actors.ActorStateEnum
	animDuration    animDurationFn

	heldPrev bool
	buffered bool
	animWait int
}

// New constructs a Controller.
// meleeAttackEnum is the game-registered state for the main melee swing.
// stepStates are the per-step animation states (game-registered).
// animDuration computes frame count for a given combo step index.
func New(
	w *weapon.MeleeWeapon,
	st *State,
	meleeAttackEnum actors.ActorStateEnum,
	stepStates []actors.ActorStateEnum,
	animDuration animDurationFn,
) *Controller {
	return &Controller{
		weapon:          w,
		state:           st,
		meleeAttackEnum: meleeAttackEnum,
		stepStates:      stepStates,
		animDuration:    animDuration,
	}
}

// Install wires the controller into a Character: adds it as a StateContributor
// and installs it as the StateTransitionHandler.
func (c *Controller) Install(char *actors.Character) {
	char.AddStateContributor(c)
	char.SetStateTransitionHandler(c.handleTransition)
}

// ContributeState implements actors.StateContributor. Returns the per-step
// animation state while the weapon is in startup or mid-swing.
func (c *Controller) ContributeState(_ actors.ActorStateEnum) (actors.ActorStateEnum, bool) {
	if !c.weapon.IsSwinging() && !c.weapon.IsInStartup() {
		return 0, false
	}
	step := c.weapon.StepIndex()
	if step >= 0 && step < len(c.stepStates) {
		return c.stepStates[step], true
	}
	return 0, false
}

// handleTransition is installed as the Character's StateTransitionHandler. It
// drives State.Update while the character is in any melee-owned state (the
// meleeAttackEnum or any per-step enum) and exits to the return state when the
// swing animation finishes.
func (c *Controller) handleTransition(char *actors.Character) bool {
	if c.state == nil {
		return false
	}

	state := char.State()
	if !c.isMeleeState(state) {
		return false
	}

	next := c.state.Update()
	if !c.isMeleeState(next) {
		char.SetNewStateFatal(next)
		return false
	}
	return true
}

// isMeleeState reports whether s is the parent meleeAttackEnum or any step state.
func (c *Controller) isMeleeState(s actors.ActorStateEnum) bool {
	if s == c.meleeAttackEnum {
		return true
	}
	for _, x := range c.stepStates {
		if s == x {
			return true
		}
	}
	return false
}

// SetSpace propagates the physics space to the attack state so ApplyHitbox
// resolves correctly. Call once per frame before the character Update.
func (c *Controller) SetSpace(space contractsbody.BodiesSpace) {
	if c.state != nil {
		c.state.SetSpace(space)
	}
}

// Tick advances the weapon one frame when the character is NOT in a melee state.
// During a melee state, State.Update (called by handleTransition) owns weapon.Update.
// This ensures cooldown and combo-window frames drain even after the animation ends.
// Call once per frame before HandleInput.
func (c *Controller) Tick(char *actors.Character) {
	if c.weapon == nil || c.isMeleeState(char.State()) {
		return
	}
	c.weapon.Update()
}

// HandleInput processes one frame of melee input. meleeHeld is the raw button
// state; dashPressed / jumpPressed are edge-detected interrupts; isGrounded and
// isDucking describe the owner's current posture.
//
// Returns true if a melee attack was entered this frame.
func (c *Controller) HandleInput(meleeHeld, dashPressed, jumpPressed, isGrounded, isDucking bool) bool {
	if c.animWait > 0 {
		c.animWait--
	}

	if (dashPressed || jumpPressed) && c.weapon.ComboWindowRemaining() > 0 {
		c.weapon.ResetCombo()
		c.buffered = false
		c.animWait = 0
	}

	meleePressed := meleeHeld && !c.heldPrev
	if meleePressed && isGrounded && (c.weapon.IsSwinging() || c.animWait > 0) {
		c.buffered = true
	}

	entered := false
	wantFire := c.weapon.CanFire() && !c.weapon.IsSwinging() && c.animWait == 0 && !isDucking &&
		(meleePressed || (isGrounded && c.buffered && c.weapon.ComboWindowRemaining() > 0))
	if wantFire {
		c.fire(isGrounded)
		entered = true
	}

	if !c.weapon.IsSwinging() && c.weapon.ComboWindowRemaining() == 0 && c.animWait == 0 {
		c.buffered = false
	}

	c.heldPrev = meleeHeld
	return entered
}

// fire advances the combo (if applicable) and transitions the character into
// the melee attack state. The state's OnStart owns Fire + VFX.
func (c *Controller) fire(isGrounded bool) {
	if isGrounded && c.weapon.ComboWindowRemaining() > 0 {
		c.weapon.AdvanceCombo()
	}
	animFrames := 0
	if c.animDuration != nil {
		animFrames = c.animDuration(c.weapon.StepIndex())
	}
	c.buffered = false
	c.animWait = animFrames
	if c.state != nil {
		c.state.SetAnimationFrames(animFrames)
	}
}

// EnterAttackState transitions the character into the per-step melee state
// matching the current weapon step. Callers must invoke this immediately after
// HandleInput returns true. Using the step state directly (rather than the
// parent meleeAttackEnum) lets sprite resolution find the correct combo-step
// image, since the State instance is registered for every step enum.
func (c *Controller) EnterAttackState(char *actors.Character) {
	step := c.weapon.StepIndex()
	target := c.meleeAttackEnum
	if step >= 0 && step < len(c.stepStates) {
		target = c.stepStates[step]
	}
	char.SetNewStateFatal(target)
}

// IsBlockingMovement reports whether horizontal movement should be suppressed:
// the weapon is mid-swing or in the post-swing animWait window.
func (c *Controller) IsBlockingMovement() bool {
	return c.weapon.IsSwinging() || c.animWait > 0
}

// StepCount returns the number of combo steps in the weapon.
func (c *Controller) StepCount() int {
	return len(c.weapon.Steps())
}

// OnInterrupt resets the combo chain. Call when the player takes damage or an
// action (dash, jump, etc.) interrupts the melee.
func (c *Controller) OnInterrupt() {
	c.weapon.ResetCombo()
	c.buffered = false
	c.animWait = 0
}
