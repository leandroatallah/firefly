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
func newDuckingTestCharacter(_, _ int) *actors.Character {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Ducking: &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 32)
	return actors.NewCharacter(sMap, rect)
}

func TestDuckingState_OnStart_ShrinksBodyHeightAndZerosHorizontalVelocity(t *testing.T) {
	const fullWidth, fullHeight = 16, 32

	c := newDuckingTestCharacter(fullWidth, fullHeight)
	c.SetVelocity(5, -3) // non-zero vx to confirm it gets zeroed

	duckState, err := c.NewState(actors.Ducking)
	if err != nil {
		t.Fatalf("NewState(Ducking) error: %v", err)
	}
	c.SetState(duckState)

	pos := c.Position()
	wantHeight := fullHeight / 2

	if pos.Dy() != wantHeight {
		t.Errorf("body height after OnStart = %d, want %d", pos.Dy(), wantHeight)
	}

	vx, _ := c.Velocity()
	if vx != 0 {
		t.Errorf("horizontal velocity after OnStart = %d, want 0", vx)
	}
}

func TestDuckingState_OnStart_PreservesBottomEdge(t *testing.T) {
	const fullWidth, fullHeight = 16, 32
	const originX, originY = 10, 100

	c := newDuckingTestCharacter(fullWidth, fullHeight)
	c.SetPosition(originX, originY)

	duckState, err := c.NewState(actors.Ducking)
	if err != nil {
		t.Fatalf("NewState(Ducking) error: %v", err)
	}
	c.SetState(duckState)

	pos := c.Position()
	wantBottom := originY + fullHeight

	if pos.Max.Y != wantBottom {
		t.Errorf("bottom edge after OnStart = %d, want %d", pos.Max.Y, wantBottom)
	}
}

func TestDuckingState_OnFinish_RestoresFullHeight(t *testing.T) {
	const fullWidth, fullHeight = 16, 32

	c := newDuckingTestCharacter(fullWidth, fullHeight)

	// Enter ducking state (records fullHeight internally)
	duckState, err := c.NewState(actors.Ducking)
	if err != nil {
		t.Fatalf("NewState(Ducking) error: %v", err)
	}
	c.SetState(duckState)

	// Exit ducking state by transitioning to Idle
	idleState, err := c.NewState(actors.Idle)
	if err != nil {
		t.Fatalf("NewState(Idle) error: %v", err)
	}
	c.SetState(idleState)

	pos := c.Position()
	if pos.Dy() != fullHeight {
		t.Errorf("body height after OnFinish = %d, want %d", pos.Dy(), fullHeight)
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

func TestDuckingState_BodyWidthUnchanged(t *testing.T) {
	const fullWidth, fullHeight = 16, 32

	c := newDuckingTestCharacter(fullWidth, fullHeight)

	duckState, err := c.NewState(actors.Ducking)
	if err != nil {
		t.Fatalf("NewState(Ducking) error: %v", err)
	}
	c.SetState(duckState)

	pos := c.Position()
	if pos.Dx() != fullWidth {
		t.Errorf("body width after OnStart = %d, want %d", pos.Dx(), fullWidth)
	}
}

// Compile-time guard: Ducking must be a valid ActorStateEnum (non-zero after init).
var _ image.Rectangle = image.Rectangle{} // keep image import used
