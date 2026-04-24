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
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	_ "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
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
