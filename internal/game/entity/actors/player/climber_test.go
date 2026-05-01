package gameplayer

import (
	"os"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"
)

// AC3 bullet 2 — Taking damage mid-combo resets the chain to step 1.
// Verifies ClimberPlayer.Hurt calls weapon.ResetCombo before state transition.
func TestClimberPlayer_Hurt_ResetsCombo(t *testing.T) {
	ctx := &app.AppContext{
		Assets:       os.DirFS("."),
		ActorManager: actors.NewManager(),
		Space:        space.NewSpace(),
		EventManager: event.NewManager(),
	}

	p, err := NewClimberPlayer(ctx)
	if err != nil {
		t.Fatalf("failed to create climber player: %v", err)
	}
	climber, ok := p.(*ClimberPlayer)
	if !ok {
		t.Fatalf("NewClimberPlayer returned %T, want *ClimberPlayer", p)
	}

	// Build a 3-step combo weapon and attach it to the player.
	steps := []weapon.ComboStep{
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 24 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 12 * 16, HitboxOffsetY16: 0},
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 28 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 14 * 16, HitboxOffsetY16: 0},
		{Damage: 2, ActiveFrames: [2]int{3, 5}, HitboxW16: 32 * 16, HitboxH16: 20 * 16, HitboxOffsetX16: 16 * 16, HitboxOffsetY16: 0},
	}
	w := weapon.NewMeleeWeapon("player_melee", 0 /*cooldown*/, 15 /*combo window*/, steps)
	climber.SetMelee(w, nil)

	// Open the combo window at step 1 (index 1) so a reset is observable.
	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	for i := 0; i <= 5+1; i++ {
		w.Update()
	}
	if !w.AdvanceCombo() {
		t.Fatalf("precondition: AdvanceCombo returned false; combo window must be open")
	}
	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	for i := 0; i <= 5+1; i++ {
		w.Update()
	}
	if w.StepIndex() != 1 {
		t.Fatalf("precondition: StepIndex() = %d, want 1", w.StepIndex())
	}
	if w.ComboWindowRemaining() == 0 {
		t.Fatalf("precondition: combo window must be open before Hurt")
	}

	climber.Hurt(1)

	if w.StepIndex() != 0 {
		t.Errorf("after Hurt, StepIndex() = %d, want 0 (combo must reset)", w.StepIndex())
	}
	if w.ComboWindowRemaining() != 0 {
		t.Errorf("after Hurt, ComboWindowRemaining() = %d, want 0 (combo must reset)", w.ComboWindowRemaining())
	}
}

// US-042 RED-5 — ClimberPlayer.Update must not drive the melee weapon directly.
// Instead, on a melee press it transitions to StateMeleeAttack via the state
// machine (Fire/Update/IsHitboxActive/ApplyHitbox are owned by MeleeAttackState).
//
// Observable assertion: after a single melee press, exactly one swing has begun
// (weapon.IsSwinging() == true once and never re-fires within the same Update);
// and the climber's state is StateMeleeAttack (or one of its per-step states).
func TestClimberPlayer_Update_MeleePressTransitionsToMeleeAttackState(t *testing.T) {
	origReader := input.CommandsReader
	t.Cleanup(func() { input.CommandsReader = origReader })

	ctx := &app.AppContext{
		Assets:       os.DirFS("."),
		ActorManager: actors.NewManager(),
		Space:        space.NewSpace(),
		EventManager: event.NewManager(),
	}

	p, err := NewClimberPlayer(ctx)
	if err != nil {
		t.Fatalf("failed to create climber player: %v", err)
	}
	climber := p.(*ClimberPlayer)

	steps := []weapon.ComboStep{
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 24 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 12 * 16, HitboxOffsetY16: 0},
	}
	w := weapon.NewMeleeWeapon("player_melee", 0 /*cooldown*/, 0 /*window*/, steps)
	climber.SetMelee(w, nil)

	// One frame of "melee not pressed" to seed meleeHeldPrev=false.
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{} }
	if err := climber.Update(ctx.Space); err != nil {
		t.Fatalf("seed Update failed: %v", err)
	}
	if w.IsSwinging() {
		t.Fatalf("precondition: weapon must not be swinging before melee press")
	}

	// Now press melee. The state machine must drive Fire — climber.Update must
	// not invoke Fire directly (RED-5). After the refactor, Fire still happens
	// once because MeleeAttackState.OnStart owns it.
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Melee: true} }
	if err := climber.Update(ctx.Space); err != nil {
		t.Fatalf("press Update failed: %v", err)
	}

	if !w.IsSwinging() {
		t.Errorf("after melee press: weapon.IsSwinging() = false, want true (state OnStart must Fire)")
	}

	// The state machine must reflect the melee attack — either StateMeleeAttack
	// or one of the per-step states registered for it.
	st := climber.State()
	stepStates := kitstates.MeleeAttackStepStates(len(w.Steps()))
	isMeleeState := st == kitstates.StateMeleeAttack
	for _, s := range stepStates {
		if st == s {
			isMeleeState = true
		}
	}
	if !isMeleeState {
		t.Errorf("after melee press: climber.State() = %v, want StateMeleeAttack (or per-step variant)", st)
	}
}
