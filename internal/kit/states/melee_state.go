package kitstates

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	meleeengine "github.com/boilerplate/ebiten-template/internal/kit/combat/melee"
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

// MeleeAttackState is a type alias for the engine-side melee State so that
// existing game-layer callers and tests require no changes.
type MeleeAttackState = meleeengine.State

// NewMeleeAttackState constructs a MeleeAttackState. vfx may be nil.
// returnTo is computed dynamically from the owner's grounded state at OnStart time.
func NewMeleeAttackState(owner meleeOwnerIface, space contractsbody.BodiesSpace, w meleeWeaponIface, vfx meleeVFXSpawner) *MeleeAttackState {
	return meleeengine.NewState(owner, space, w, vfx, StateMeleeAttack, StateGrounded, actors.Falling)
}

// InstallMeleeAttackState constructs a MeleeAttackState, registers it as the
// per-actor instance for StateMeleeAttack and every step state on the given
// character, and returns the instance so the caller can store it for
// Update-time space injection.
func InstallMeleeAttackState(char *actors.Character, owner meleeOwnerIface, w meleeWeaponIface, vfx meleeVFXSpawner, stepStates []actors.ActorStateEnum) *MeleeAttackState {
	return meleeengine.InstallState(char, owner, w, vfx, StateMeleeAttack, StateGrounded, actors.Falling, stepStates)
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
