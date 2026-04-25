package gamestates

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
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

// meleeVFXSpawner is the minimum surface needed by MeleeAttackState to render
// the slash VFX. It is satisfied by vfx.Manager.
type meleeVFXSpawner interface {
	SpawnDirectionalPuff(typeKey string, x, y float64, faceRight bool, count int, randRange float64)
}

// meleeOwnerIface is the minimum owner interface needed by MeleeAttackState.
type meleeOwnerIface interface {
	contractsbody.Collidable
	FaceDirection() animation.FacingDirectionEnum
	IsFalling() bool
	IsGoingUp() bool
	IsDucking() bool
}

// MeleeAttackState is the actor state active during a melee swing.
type MeleeAttackState struct {
	owner      meleeOwnerIface
	space      contractsbody.BodiesSpace
	weapon     meleeWeaponIface
	vfx        meleeVFXSpawner // may be nil
	returnTo   actors.ActorStateEnum
	animFrames int
	frame      int
	stepUsed   int
}

// NewMeleeAttackState constructs a MeleeAttackState. vfx may be nil.
// returnTo is computed dynamically from the owner's grounded state at OnStart time.
func NewMeleeAttackState(owner meleeOwnerIface, space contractsbody.BodiesSpace, w meleeWeaponIface, vfx meleeVFXSpawner) *MeleeAttackState {
	return &MeleeAttackState{
		owner:  owner,
		space:  space,
		weapon: w,
		vfx:    vfx,
	}
}

// SetAnimationFrames sets the total number of animation frames for the swing.
func (s *MeleeAttackState) SetAnimationFrames(n int) { s.animFrames = n }

// SetSpace updates the BodiesSpace used by ApplyHitbox. Call once per frame
// before the character's state machine is ticked.
func (s *MeleeAttackState) SetSpace(sp contractsbody.BodiesSpace) { s.space = sp }

// StepUsed returns the combo step index that was active when OnStart was called.
func (s *MeleeAttackState) StepUsed() int { return s.stepUsed }

// OnStart fires the weapon, spawns VFX, and resets the frame counter.
// If the owner is ducking the swing is aborted: no Fire, no VFX, frame is set
// to animFrames so the very next Update resolves to returnTo immediately.
func (s *MeleeAttackState) OnStart(_ int) {
	s.frame = 0

	// Resolve dynamic returnTo from owner grounded state.
	if !s.owner.IsFalling() && !s.owner.IsGoingUp() {
		s.returnTo = StateGrounded
	} else {
		s.returnTo = actors.Falling
	}

	// Ducking abort: skip Fire and VFX; let Update resolve on next tick.
	if s.owner.IsDucking() {
		s.frame = s.animFrames
		return
	}

	// Capture step index BEFORE Fire (Fire may mutate internal bookkeeping).
	s.stepUsed = s.weapon.StepIndex()

	x16, y16 := s.owner.GetPosition16()
	faceDir := s.owner.FaceDirection()
	s.weapon.Fire(x16, y16, faceDir, contractsbody.ShootDirectionStraight, 0)

	if s.vfx != nil {
		offsetX16 := fp16.To16(12)
		if faceDir == animation.FaceDirectionLeft {
			offsetX16 = -offsetX16
		}
		px := float64(fp16.From16(x16 + offsetX16))
		py := float64(fp16.From16(y16))
		s.vfx.SpawnDirectionalPuff("melee_slash", px, py, faceDir == animation.FaceDirectionRight, 1, 0.0)
	}
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

// InstallMeleeAttackState constructs a MeleeAttackState, registers it as the
// per-actor instance for StateMeleeAttack on the given character, and returns
// the instance so the caller can store it for Update-time space injection.
func InstallMeleeAttackState(char *actors.Character, owner meleeOwnerIface, w meleeWeaponIface, vfx meleeVFXSpawner) *MeleeAttackState {
	st := NewMeleeAttackState(owner, nil /*space injected at Update time via SetSpace*/, w, vfx)
	char.SetStateInstance(StateMeleeAttack, st)
	return st
}

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
