package platformer

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
)

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
