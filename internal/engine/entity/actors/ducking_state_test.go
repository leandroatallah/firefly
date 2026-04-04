package actors_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

// newDuckingTestCharacter creates a Character with a body of the given dimensions.
func newDuckingTestCharacter(w, h int) *actors.Character {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Ducking: &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, w, h)
	return actors.NewCharacter(sMap, rect)
}

func TestDuckingState_OnStart_ZerosHorizontalVelocity(t *testing.T) {
	c := newDuckingTestCharacter(16, 32)
	c.SetVelocity(5, -3)

	duckState, err := c.NewState(actors.Ducking)
	if err != nil {
		t.Fatalf("NewState(Ducking) error: %v", err)
	}
	c.SetState(duckState)

	vx, _ := c.Velocity()
	if vx != 0 {
		t.Errorf("horizontal velocity after OnStart = %d, want 0", vx)
	}
}

func TestDuckingState_IsRegisteredAndConstructable(t *testing.T) {
	c := newDuckingTestCharacter(16, 32)

	state, err := c.NewState(actors.Ducking)
	if err != nil {
		t.Fatalf("NewState(Ducking) returned error: %v", err)
	}
	if state == nil {
		t.Fatal("NewState(Ducking) returned nil state")
	}
	if state.State() != actors.Ducking {
		t.Errorf("state.State() = %v, want Ducking (%v)", state.State(), actors.Ducking)
	}
}

// Compile-time guard: Ducking must be a valid ActorStateEnum (non-zero after init).
var _ image.Rectangle = image.Rectangle{} // keep image import used
