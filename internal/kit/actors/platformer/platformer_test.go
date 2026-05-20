package platformer

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			UpwardGravity:   2,
		},
	})
}

type stubSpace struct{ body.BodiesSpace }

func (s *stubSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return nil
}
func (s *stubSpace) Query(_ image.Rectangle) []body.Collidable        { return nil }
func (s *stubSpace) ResolveCollisions(_ body.Collidable) (bool, bool) { return false, false }

func newPlatformerTestCharacter(states ...actors.ActorStateEnum) *actors.Character {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Walking: &sprites.Sprite{Image: img},
		actors.Jumping: &sprites.Sprite{Image: img},
		actors.Falling: &sprites.Sprite{Image: img},
		actors.Landing: &sprites.Sprite{Image: img},
		actors.Hurted:  &sprites.Sprite{Image: img},
		actors.Dying:   &sprites.Sprite{Image: img},
		actors.Dead:    &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxHealth(100)
	c.SetHealth(100)
	c.SetMovementTransitionHandler(platformerMovementTransitions)
	if len(states) > 0 && states[0] != actors.Idle {
		s, _ := c.NewState(states[0])
		c.SetState(s)
	}
	return c
}

// JumpingToIdle/LandingToIdle/FallingToLanding use no movement model so
// UpdateMovement is a no-op and onGround defaults to true.

func TestPlatformerMovement_JumpingToIdle(t *testing.T) {
	c := newPlatformerTestCharacter(actors.Jumping)
	c.Update(nil)
	if c.State() != actors.Idle {
		t.Errorf("expected Idle, got %v", c.State())
	}
}

func TestPlatformerMovement_LandingToIdle(t *testing.T) {
	c := newPlatformerTestCharacter(actors.Landing)
	c.Update(nil)
	if c.State() != actors.Idle {
		t.Errorf("expected Idle, got %v", c.State())
	}
}

func TestPlatformerMovement_FallingToLanding(t *testing.T) {
	c := newPlatformerTestCharacter(actors.Falling)
	c.Update(nil)
	if c.State() != actors.Landing {
		t.Errorf("expected Landing, got %v", c.State())
	}
}

func TestPlatformerMovement_JumpFlicker(t *testing.T) {
	c := newPlatformerTestCharacter(actors.Jumping)
	model := physicsmovement.NewPlatformMovementModel(nil)
	model.SetOnGround(false)
	c.SetMovementModel(model)
	c.SetVelocity(0, -100)
	c.Update(&stubSpace{})
	if c.State() != actors.Jumping {
		t.Errorf("expected Jumping while vy < 0 and airborne, got %v", c.State())
	}
}

func TestPlatformerMovement_AirPeakStaysJumping(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  256,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 1.0,
			DownwardGravity: 5,
			UpwardGravity:   4,
			MaxFallSpeed:    100,
		},
	})

	c := newPlatformerTestCharacter(actors.Jumping)
	c.SetMaxSpeed(100)
	model := physicsmovement.NewPlatformMovementModel(nil)
	model.SetOnGround(false)
	c.SetMovementModel(model)
	c.SetVelocity(0, -4)
	c.Update(&stubSpace{})
	if c.State() != actors.Jumping {
		_, curVy := c.Velocity()
		t.Errorf("expected Jumping at air peak (vy=%d), got %v", curVy, c.State())
	}

	c.SetVelocity(0, 5)
	c.Update(&stubSpace{})
	if c.State() != actors.Falling {
		t.Errorf("expected Falling when vy >= threshold, got %v", c.State())
	}
}

func TestPlatformerCharacter_SetOnJump(t *testing.T) {
	character := &PlatformerCharacter{}
	called := false

	character.SetOnJump(func(pos image.Point) {
		called = true
	})

	if character.jumpHandler == nil {
		t.Fatal("Expected jumpHandler to be set")
	}

	// Call the handler
	character.jumpHandler(image.Point{X: 10, Y: 20})

	if !called {
		t.Error("Expected jumpHandler to be called")
	}
}

func TestPlatformerCharacter_SetOnLand(t *testing.T) {
	character := &PlatformerCharacter{}
	called := false

	character.SetOnLand(func(pos image.Point) {
		called = true
	})

	if character.landHandler == nil {
		t.Fatal("Expected landHandler to be set")
	}

	character.landHandler(image.Point{X: 10, Y: 20})

	if !called {
		t.Error("Expected landHandler to be called")
	}
}

func TestPlatformerCharacter_SetOnFall(t *testing.T) {
	character := &PlatformerCharacter{}
	called := false

	character.SetOnFall(func(pos image.Point) {
		called = true
	})

	if character.fallHandler == nil {
		t.Fatal("Expected fallHandler to be set")
	}

	character.fallHandler(image.Point{X: 10, Y: 20})

	if !called {
		t.Error("Expected fallHandler to be called")
	}
}

func TestPlatformerCharacter_OnJump(t *testing.T) {
	t.Run("with handler", func(t *testing.T) {
		// Create a minimal character for testing
		spriteMap := sprites.SpriteMap{}
		char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
		character := &PlatformerCharacter{
			Character: char,
		}

		var capturedPos image.Point

		character.SetOnJump(func(pos image.Point) {
			capturedPos = pos
		})

		// Set position
		character.SetPosition(100, 200)

		character.OnJump()

		// Expected: bottom center of character (100+8, 200+16) = (108, 216)
		expectedPos := image.Point{X: 108, Y: 216}
		if capturedPos != expectedPos {
			t.Errorf("Expected position %v, got %v", expectedPos, capturedPos)
		}
	})

	t.Run("without handler", func(t *testing.T) {
		spriteMap := sprites.SpriteMap{}
		char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
		character := &PlatformerCharacter{
			Character: char,
		}
		character.SetPosition(100, 200)

		// Should not panic
		character.OnJump()
	})
}

func TestPlatformerCharacter_OnLand(t *testing.T) {
	t.Run("with handler", func(t *testing.T) {
		spriteMap := sprites.SpriteMap{}
		char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
		character := &PlatformerCharacter{
			Character: char,
		}

		var capturedPos image.Point

		character.SetOnLand(func(pos image.Point) {
			capturedPos = pos
		})

		character.SetPosition(50, 100)
		character.OnLand()

		// Expected: bottom center (50+8, 100+16) = (58, 116)
		expectedPos := image.Point{X: 58, Y: 116}
		if capturedPos != expectedPos {
			t.Errorf("Expected position %v, got %v", expectedPos, capturedPos)
		}
	})

	t.Run("without handler", func(t *testing.T) {
		spriteMap := sprites.SpriteMap{}
		char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
		character := &PlatformerCharacter{
			Character: char,
		}
		character.SetPosition(50, 100)

		// Should not panic
		character.OnLand()
	})
}

func TestPlatformerCharacter_OnFall(t *testing.T) {
	t.Run("with handler", func(t *testing.T) {
		spriteMap := sprites.SpriteMap{}
		char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
		character := &PlatformerCharacter{
			Character: char,
		}

		var capturedPos image.Point

		character.SetOnFall(func(pos image.Point) {
			capturedPos = pos
		})

		character.SetPosition(75, 150)
		character.OnFall()

		// Expected: bottom center (75+8, 150+16) = (83, 166)
		expectedPos := image.Point{X: 83, Y: 166}
		if capturedPos != expectedPos {
			t.Errorf("Expected position %v, got %v", expectedPos, capturedPos)
		}
	})

	t.Run("without handler", func(t *testing.T) {
		spriteMap := sprites.SpriteMap{}
		char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
		character := &PlatformerCharacter{
			Character: char,
		}
		character.SetPosition(75, 150)

		// Should not panic
		character.OnFall()
	})
}

func TestPlatformerCharacter_CoinCount(t *testing.T) {
	character := &PlatformerCharacter{}

	if character.coinCount != 0 {
		t.Errorf("Expected initial coinCount 0, got %d", character.coinCount)
	}

	character.coinCount = 5
	if character.coinCount != 5 {
		t.Errorf("Expected coinCount 5, got %d", character.coinCount)
	}
}

func TestPlatformerCharacter_MovementBlockers(t *testing.T) {
	character := &PlatformerCharacter{}

	if character.movementBlockers != 0 {
		t.Errorf("Expected initial movementBlockers 0, got %d", character.movementBlockers)
	}

	character.movementBlockers = 3
	if character.movementBlockers != 3 {
		t.Errorf("Expected movementBlockers 3, got %d", character.movementBlockers)
	}
}

func TestPlatformerCharacter_Fields(t *testing.T) {
	spriteMap := sprites.SpriteMap{}
	char := actors.NewCharacter(spriteMap, bodyphysics.NewRect(0, 0, 16, 16))
	character := &PlatformerCharacter{
		Character: char,
	}

	// Test initial values
	if character.Character != char {
		t.Error("Expected Character to be set")
	}
	if character.coinCount != 0 {
		t.Errorf("Expected coinCount 0, got %d", character.coinCount)
	}
	if character.movementBlockers != 0 {
		t.Errorf("Expected movementBlockers 0, got %d", character.movementBlockers)
	}
}
