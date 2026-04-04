package body

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewObstacleRect(t *testing.T) {
	rect := NewRect(0, 0, 20, 30)
	obs := NewObstacleRect(rect)

	if obs == nil {
		t.Fatal("NewObstacleRect returned nil")
	}
	if obs.MovableBody == nil {
		t.Error("MovableBody was not initialized")
	}
	if obs.CollidableBody == nil {
		t.Error("CollidableBody was not initialized")
	}
}

func TestObstacleRect_Ownership(t *testing.T) {
	rect := NewRect(0, 0, 20, 30)
	obs := NewObstacleRect(rect)

	// Body's owner should be MovableBody (access via MovableBody)
	if obs.MovableBody.Body.Owner() != obs.MovableBody {
		t.Error("expected Body's owner to be MovableBody")
	}

	// MovableBody's owner should be ObstacleRect
	if obs.MovableBody.Owner() != obs {
		t.Error("expected MovableBody's owner to be ObstacleRect")
	}

	// CollidableBody's owner should be ObstacleRect
	if obs.CollidableBody.Owner() != obs {
		t.Error("expected CollidableBody's owner to be ObstacleRect")
	}
}

func TestObstacleRect_ID(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	obs.SetID("obstacle-1")
	if obs.ID() != "obstacle-1" {
		t.Errorf("expected ID 'obstacle-1'; got '%s'", obs.ID())
	}
}

func TestObstacleRect_Position(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	obs.SetPosition(50, 60)

	pos := obs.Position()
	if pos.Min.X != 50 || pos.Min.Y != 60 {
		t.Errorf("expected position (50,60); got (%d,%d)", pos.Min.X, pos.Min.Y)
	}
}

func TestObstacleRect_GetPositionMin(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	obs.SetPosition(100, 200)

	x, y := obs.GetPositionMin()
	if x != 100 || y != 200 {
		t.Errorf("expected (100,200); got (%d,%d)", x, y)
	}
}

func TestObstacleRect_SetPosition16(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	x16 := 100 << 16
	y16 := 200 << 16
	obs.SetPosition16(x16, y16)

	gotX16, gotY16 := obs.GetPosition16()
	if gotX16 != x16 || gotY16 != y16 {
		t.Errorf("expected (%d,%d); got (%d,%d)", x16, y16, gotX16, gotY16)
	}
}

func TestObstacleRect_GetShape(t *testing.T) {
	rect := NewRect(0, 0, 25, 35)
	obs := NewObstacleRect(rect)

	shape := obs.GetShape()
	if shape != rect {
		t.Errorf("expected shape %v; got %v", rect, shape)
	}

	if shape.Width() != 25 || shape.Height() != 35 {
		t.Errorf("expected 25x35; got %dx%d", shape.Width(), shape.Height())
	}
}

func TestObstacleRect_AddCollisionBodies_Empty(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))
	obs.SetID("obstacle")

	// Call with no arguments - should create default collision body
	obs.AddCollisionBodies()

	shapes := obs.CollisionShapes()
	if len(shapes) != 1 {
		t.Errorf("expected 1 collision shape; got %d", len(shapes))
	}

	// Check the collision body has correct ID
	expectedID := "obstacle_COLLISION_0"
	if shapes[0].ID() != expectedID {
		t.Errorf("expected ID '%s'; got '%s'", expectedID, shapes[0].ID())
	}
}

func TestObstacleRect_AddCollisionBodies_WithRects(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))
	obs.SetID("obstacle")

	other := NewCollidableBodyFromRect(NewRect(0, 0, 5, 5))
	other.SetID("collision-1")
	obs.AddCollisionBodies(other)

	shapes := obs.CollisionShapes()
	if len(shapes) != 1 {
		t.Errorf("expected 1 collision shape; got %d", len(shapes))
	}
	if shapes[0] != other {
		t.Error("expected collision shape to be the one provided")
	}
}

func TestObstacleRect_Draw(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))
	obs.SetPosition(5, 5)

	screen := ebiten.NewImage(100, 100)

	// Should not panic
	obs.Draw(screen)
}

func TestObstacleRect_Image(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 20, 30))

	img := obs.Image()
	if img == nil {
		t.Fatal("Image() returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 20 || bounds.Dy() != 30 {
		t.Errorf("expected image 20x30; got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestObstacleRect_ImageOptions(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	opts := obs.ImageOptions()
	if opts == nil {
		t.Fatal("ImageOptions() returned nil")
	}
}

func TestObstacleRect_UpdateImageOptions(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))
	obs.SetPosition(50, 60)

	obs.UpdateImageOptions()

	opts := obs.ImageOptions()
	// Check that GeoM has been set (translation)
	// The exact values depend on ebiten's GeoM implementation
	// We just verify it doesn't panic and options are updated
	if opts == nil {
		t.Error("ImageOptions should not be nil after UpdateImageOptions")
	}
}

func TestObstacleRect_ImplementsObstacle(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	// Verify it implements body.Obstacle interface
	var _ body.Obstacle = obs
}

func TestObstacleRect_ForwardedMethods(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	// Test that all forwarded methods work correctly
	obs.SetID("test")
	if obs.ID() != "test" {
		t.Errorf("ID forwarding failed")
	}

	obs.SetPosition(10, 20)
	x, y := obs.GetPositionMin()
	if x != 10 || y != 20 {
		t.Errorf("Position forwarding failed: got (%d,%d)", x, y)
	}

	shape := obs.GetShape()
	if shape == nil {
		t.Error("GetShape forwarding failed")
	}
}

func TestObstacleRect_CollidableInterface(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))
	obs.SetID("test")

	// Test CollidableBody methods are accessible
	obs.SetIsObstructive(true)
	if !obs.IsObstructive() {
		t.Error("SetIsObstructive/IsObstructive not working")
	}

	// Test touchable
	obs.SetTouchable(obs)
	if obs.GetTouchable() != obs {
		t.Error("SetTouchable/GetTouchable not working")
	}
}

func TestObstacleRect_MovableInterface(t *testing.T) {
	obs := NewObstacleRect(NewRect(0, 0, 10, 10))

	// Test MovableBody methods are accessible
	err := obs.SetSpeed(100)
	if err != nil {
		t.Errorf("SetSpeed failed: %v", err)
	}

	if obs.Speed() != 100 {
		t.Errorf("expected speed 100; got %d", obs.Speed())
	}

	// Test movement
	obs.OnMoveRight(5)
	accX, accY := obs.Acceleration()
	if accX <= 0 {
		t.Errorf("expected positive accelerationX; got %d", accX)
	}
	if accY != 0 {
		t.Errorf("expected zero accelerationY; got %d", accY)
	}
}
