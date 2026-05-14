package beatemup_test

import (
	"testing"
	"testing/fstest"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
)

// newTestFixtures returns the minimal arguments accepted by the new
// NewBeatEmUpCharacter constructor (see SPEC §3). Assets map is empty so
// no PNG decoding is needed in the unit-test environment.
func newTestFixtures() (fstest.MapFS, map[string]animation.SpriteState, schemas.SpriteData, *bodyphysics.Rect) {
	fsys := fstest.MapFS{}
	stateMap := map[string]animation.SpriteState{}
	spriteData := schemas.SpriteData{
		Assets:          map[string]schemas.AssetData{},
		FrameRate:       1,
		FacingDirection: 0, // right
	}
	bodyRect := bodyphysics.NewRect(0, 0, 8, 8)
	return fsys, stateMap, spriteData, bodyRect
}

// T-B1: NewBeatEmUpCharacter returns a non-nil character whose embedded
// Character and MeleeCharacter are both initialised.
func TestNewBeatEmUpCharacter_EmbedsCharacterAndMeleeCharacter(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewBeatEmUpCharacter returned nil character")
	}
	if c.Character == nil {
		t.Error("expected embedded *actors.Character to be initialised")
	}
	if c.MeleeCharacter == nil {
		t.Error("expected embedded *MeleeCharacter to be initialised")
	}
}

// T-B2: BeatEmUpCharacter owns a BeatEmUpMovementModel — not a PlatformMovementModel.
func TestNewBeatEmUpCharacter_OwnsBeatEmUpMovementModel(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	model := c.MovementModel()
	if model == nil {
		t.Fatal("expected MovementModel() to be non-nil after construction")
	}
	if _, ok := model.(*physicsmovement.BeatEmUpMovementModel); !ok {
		t.Fatalf("expected *BeatEmUpMovementModel, got %T", model)
	}
	if _, ok := model.(*physicsmovement.PlatformMovementModel); ok {
		t.Fatal("BeatEmUpCharacter must not own a PlatformMovementModel")
	}
}

// T-B3: Update on a zero-input frame does not panic and leaves the body at rest.
// Verifies the model is properly wired into Character.Update (no panic on
// type assertion) and that there is no gravity accumulation when idle.
func TestBeatEmUpCharacter_Update_ZeroInput_NoPanic_NoDrift(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	// Empty input — no movement commands.
	origReader := input.CommandsReader
	defer func() { input.CommandsReader = origReader }()
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{} }

	sp := space.NewSpace()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Update panicked on zero-input frame: %v", r)
		}
	}()

	if err := c.Update(sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	vx, vy := c.Velocity()
	if vx != 0 || vy != 0 {
		t.Errorf("expected zero velocity on idle frame; got (%d, %d)", vx, vy)
	}
}

// T-B4: Registering an EightDirectionalMovementSkill causes HandleInput to be
// invoked each frame: with Right pressed, the body accelerates rightward.
func TestBeatEmUpCharacter_EightDirSkill_RightInputProducesRightwardVelocity(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	c.SetSpeed(200)
	c.AddSkill(kitskills.NewEightDirectionalMovementSkill())

	origReader := input.CommandsReader
	defer func() { input.CommandsReader = origReader }()
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Right: true} }

	sp := space.NewSpace()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Update panicked with EightDirectionalMovementSkill registered: %v", r)
		}
	}()

	if err := c.Update(sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	vx, _ := c.Velocity()
	if vx <= 0 {
		t.Errorf("expected positive X velocity after Right input; got vx=%d", vx)
	}
}
