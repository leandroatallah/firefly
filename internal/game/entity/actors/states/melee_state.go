package gamestates

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// State enum for melee attack.
//
//nolint:gochecknoglobals
var StateMeleeAttack actors.ActorStateEnum

func init() {
	StateMeleeAttack = actors.RegisterState("melee_attack", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b}
	})
	// Pre-register per-step states so the sprite builder can resolve them at startup.
	MeleeAttackStepStates(3)
}

// meleeWeaponIface captures the MeleeWeapon surface needed by MeleeAttackState.
type meleeWeaponIface interface {
	combat.Weapon
	IsHitboxActive() bool
	ApplyHitbox(space contractsbody.BodiesSpace)
	StepIndex() int
	ComboWindowRemaining() int
	ResetCombo()
}

// meleeOwnerIface is the minimum owner interface needed by MeleeAttackState.
type meleeOwnerIface interface {
	contractsbody.Collidable
	FaceDirection() animation.FacingDirectionEnum
}

// MeleeAttackState is the actor state active during a melee swing.
type MeleeAttackState struct {
	owner      meleeOwnerIface
	space      contractsbody.BodiesSpace
	weapon     meleeWeaponIface
	returnTo   actors.ActorStateEnum
	animFrames int
	frame      int
	stepUsed   int
}

// NewMeleeAttackState constructs a MeleeAttackState.
// returnTo is the state to resume after the animation finishes (StateGrounded or Falling).
func NewMeleeAttackState(owner meleeOwnerIface, space contractsbody.BodiesSpace, w meleeWeaponIface, returnTo actors.ActorStateEnum) *MeleeAttackState {
	return &MeleeAttackState{
		owner:    owner,
		space:    space,
		weapon:   w,
		returnTo: returnTo,
	}
}

// SetAnimationFrames sets the total number of animation frames for the swing.
func (s *MeleeAttackState) SetAnimationFrames(n int) { s.animFrames = n }

// StepUsed returns the combo step index that was active when OnStart was called.
func (s *MeleeAttackState) StepUsed() int { return s.stepUsed }

// OnStart captures the current combo step and resets the frame counter.
// Fire is owned by the caller (ClimberPlayer.Update) and must be called before OnStart.
func (s *MeleeAttackState) OnStart(_ int) {
	s.frame = 0
	s.stepUsed = s.weapon.StepIndex()
}

// OnFinish is a no-op (weapon cooldown is self-managed).
func (s *MeleeAttackState) OnFinish() {}

// Update advances the weapon and state by one frame.
func (s *MeleeAttackState) Update() actors.ActorStateEnum {
	s.weapon.Update()
	if s.weapon.IsHitboxActive() {
		s.weapon.ApplyHitbox(s.space)
	}
	s.frame++
	if s.frame >= s.animFrames {
		return s.returnTo
	}
	return StateMeleeAttack
}

// State satisfies actors.ActorState.
func (s *MeleeAttackState) State() actors.ActorStateEnum { return StateMeleeAttack }

// GetAnimationCount satisfies actors.ActorState.
func (s *MeleeAttackState) GetAnimationCount(currentCount int) int { return currentCount - s.frame }

// IsAnimationFinished satisfies actors.ActorState.
func (s *MeleeAttackState) IsAnimationFinished() bool { return s.frame >= s.animFrames }

// TryMeleeFromFalling is a helper for wiring melee triggers from the Falling
// state. It returns the new state and true if a melee attack should begin.
func TryMeleeFromFalling(w meleeWeaponIface, meleePressed bool) (actors.ActorStateEnum, bool) {
	if !meleePressed || !w.CanFire() {
		return 0, false
	}
	return StateMeleeAttack, true
}

// ResetComboOnInterrupt resets the combo chain when dash or jump is pressed during the combo window.
func ResetComboOnInterrupt(w interface {
	ComboWindowRemaining() int
	ResetCombo()
}, dashPressed, jumpPressed bool) {
	if (dashPressed || jumpPressed) && w.ComboWindowRemaining() > 0 {
		w.ResetCombo()
	}
}

// MeleeAttackStepStates returns a slice of n state enums for per-step melee attack states,
// registered under the name pattern "melee_attack_step_<i>". Repeated calls with the same
// n return the same enum values (idempotent via GetStateEnum).
func MeleeAttackStepStates(n int) []actors.ActorStateEnum {
	out := make([]actors.ActorStateEnum, n)
	for i := range out {
		name := fmt.Sprintf("melee_attack_step_%d", i)
		if s, ok := actors.GetStateEnum(name); ok {
			out[i] = s
		} else {
			out[i] = actors.RegisterState(name, func(b actors.BaseState) actors.ActorState {
				return &actors.IdleState{BaseState: b}
			})
		}
	}
	return out
}
