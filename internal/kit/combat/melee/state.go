package melee

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

// weaponIface captures the MeleeWeapon surface needed by State.
type weaponIface interface {
	combat.Weapon
	IsHitboxActive() bool
	ApplyHitbox(space contractsbody.BodiesSpace)
	StepIndex() int
	ComboWindowRemaining() int
	ResetCombo()
}

// vfxSpawner is the minimum surface needed to render the slash VFX.
// Satisfied by vfx.Manager.
type vfxSpawner interface {
	SpawnDirectionalPuff(typeKey string, x, y float64, faceRight bool, count int, randRange float64)
}

// ownerIface is the minimum owner interface needed by State.
type ownerIface interface {
	contractsbody.Collidable
	FaceDirection() animation.FacingDirectionEnum
	IsFalling() bool
	IsGoingUp() bool
	IsDucking() bool
}

// State is the actor state active during a melee swing.
type State struct {
	owner               ownerIface
	space               contractsbody.BodiesSpace
	weapon              weaponIface
	vfx                 vfxSpawner
	meleeAttackEnum     actors.ActorStateEnum
	groundedReturnState actors.ActorStateEnum
	fallingReturnState  actors.ActorStateEnum
	returnTo            actors.ActorStateEnum
	stepStates          []actors.ActorStateEnum
	animFrames          int
	frame               int
	startCount          int
	stepUsed            int
}

// SetStepStates assigns the per-step state enums used by State() and Update()
// to report the active sprite key. Without this, sprite resolution falls back
// to the meleeAttackEnum (which typically has no registered sprite).
func (s *State) SetStepStates(stepStates []actors.ActorStateEnum) {
	s.stepStates = stepStates
}

// activeStepEnum returns the state enum that corresponds to the current swing
// step, or meleeAttackEnum if no step states are configured / index is invalid.
func (s *State) activeStepEnum() actors.ActorStateEnum {
	if s.stepUsed >= 0 && s.stepUsed < len(s.stepStates) {
		return s.stepStates[s.stepUsed]
	}
	return s.meleeAttackEnum
}

// NewState constructs a State.
// meleeAttackEnum is the state this node represents (game-registered).
// groundedReturnState / fallingReturnState are the states to transition to when the
// swing finishes, chosen dynamically from the owner's airborne status at OnStart time.
// vfx may be nil.
func NewState(
	owner ownerIface,
	space contractsbody.BodiesSpace,
	w weaponIface,
	vfx vfxSpawner,
	meleeAttackEnum, groundedReturnState, fallingReturnState actors.ActorStateEnum,
) *State {
	return &State{
		owner:               owner,
		space:               space,
		weapon:              w,
		vfx:                 vfx,
		meleeAttackEnum:     meleeAttackEnum,
		groundedReturnState: groundedReturnState,
		fallingReturnState:  fallingReturnState,
	}
}

// SetAnimationFrames sets the total number of animation frames for the swing.
func (s *State) SetAnimationFrames(n int) { s.animFrames = n }

// SetSpace updates the BodiesSpace used by ApplyHitbox. Call once per frame
// before the character's state machine is ticked.
func (s *State) SetSpace(sp contractsbody.BodiesSpace) { s.space = sp }

// StepUsed returns the combo step index that was active when OnStart was called.
func (s *State) StepUsed() int { return s.stepUsed }

// OnStart fires the weapon, spawns VFX, and resets the frame counter.
// If the owner is ducking the swing is aborted: no Fire, no VFX, frame is set
// to animFrames so the very next Update resolves to the return state immediately.
func (s *State) OnStart(currentCount int) {
	s.startCount = currentCount
	s.frame = 0

	if !s.owner.IsFalling() && !s.owner.IsGoingUp() {
		s.returnTo = s.groundedReturnState
	} else {
		s.returnTo = s.fallingReturnState
	}

	if s.owner.IsDucking() {
		s.frame = s.animFrames
		return
	}

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
func (s *State) OnFinish() {}

// Update advances the weapon and state by one frame.
func (s *State) Update() actors.ActorStateEnum {
	s.weapon.Update()
	if s.weapon.IsHitboxActive() {
		s.weapon.ApplyHitbox(s.space)
	}
	s.frame++
	if s.frame >= s.animFrames {
		return s.returnTo
	}
	return s.activeStepEnum()
}

// State satisfies actors.ActorState. It reports the per-step enum so sprite
// resolution selects the correct combo-step image.
func (s *State) State() actors.ActorStateEnum { return s.activeStepEnum() }

// GetAnimationCount satisfies actors.ActorState.
func (s *State) GetAnimationCount(currentCount int) int { return currentCount - s.startCount }

// IsAnimationFinished satisfies actors.ActorState.
func (s *State) IsAnimationFinished() bool { return s.frame >= s.animFrames }

// InstallState constructs a State, registers it as the per-actor instance for
// meleeAttackEnum AND each step state on the given character, and returns it.
// Registering for each step state ensures Character.NewState returns the same
// State instance whichever per-step enum is requested, so Update() advances a
// single shared frame counter across the swing.
func InstallState(
	char *actors.Character,
	owner ownerIface,
	w weaponIface,
	vfx vfxSpawner,
	meleeAttackEnum, groundedReturnState, fallingReturnState actors.ActorStateEnum,
	stepStates []actors.ActorStateEnum,
) *State {
	st := NewState(owner, nil, w, vfx, meleeAttackEnum, groundedReturnState, fallingReturnState)
	st.SetStepStates(stepStates)
	char.SetStateInstance(meleeAttackEnum, st)
	for _, s := range stepStates {
		char.SetStateInstance(s, st)
	}
	return st
}
