package body

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

func TestNewDefaultObstacleFactory(t *testing.T) {
	factory := NewDefaultObstacleFactory()

	if factory == nil {
		t.Fatal("NewDefaultObstacleFactory returned nil")
	}
	// obstacleMap is nil by default and will be initialized when needed
	// or can be set directly by the user
	if factory.obstacleMap == nil {
		// This is expected - the map is nil until obstacles are registered
		t.Logf("obstacleMap is nil by default (expected)")
	}
}

func TestDefaultObstacleFactory_Create_UnknownType(t *testing.T) {
	factory := NewDefaultObstacleFactory()

	_, err := factory.Create(999)
	if err == nil {
		t.Error("expected error for unknown obstacle type")
	}
}

func TestDefaultObstacleFactory_Create_RegisteredType(t *testing.T) {
	factory := NewDefaultObstacleFactory()

	// Register a custom obstacle type
	const testType ObstacleType = 1
	factory.obstacleMap = ObstacleMap{
		testType: func() body.Obstacle {
			return NewObstacleRect(NewRect(0, 0, 10, 10))
		},
	}

	obstacle, err := factory.Create(testType)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if obstacle == nil {
		t.Fatal("Create returned nil obstacle")
	}

	pos := obstacle.Position()
	if pos.Dx() != 10 || pos.Dy() != 10 {
		t.Errorf("expected obstacle 10x10; got %dx%d", pos.Dx(), pos.Dy())
	}
}

func TestDefaultObstacleFactory_Create_MultipleTypes(t *testing.T) {
	factory := NewDefaultObstacleFactory()

	const (
		typeA ObstacleType = 1
		typeB ObstacleType = 2
	)

	factory.obstacleMap = ObstacleMap{
		typeA: func() body.Obstacle {
			return NewObstacleRect(NewRect(10, 10, 20, 20))
		},
		typeB: func() body.Obstacle {
			return NewObstacleRect(NewRect(30, 30, 40, 40))
		},
	}

	obstacleA, err := factory.Create(typeA)
	if err != nil {
		t.Fatalf("Create typeA failed: %v", err)
	}

	obstacleB, err := factory.Create(typeB)
	if err != nil {
		t.Fatalf("Create typeB failed: %v", err)
	}

	posA := obstacleA.Position()
	posB := obstacleB.Position()

	if posA.Dx() != 20 || posA.Dy() != 20 {
		t.Errorf("expected obstacle A 20x20; got %dx%d", posA.Dx(), posA.Dy())
	}

	if posB.Dx() != 40 || posB.Dy() != 40 {
		t.Errorf("expected obstacle B 40x40; got %dx%d", posB.Dx(), posB.Dy())
	}
}

func TestDefaultObstacleFactory_Create_EachCallReturnsNewInstance(t *testing.T) {
	factory := NewDefaultObstacleFactory()

	const testType ObstacleType = 1
	factory.obstacleMap = ObstacleMap{
		testType: func() body.Obstacle {
			return NewObstacleRect(NewRect(0, 0, 10, 10))
		},
	}

	obstacle1, _ := factory.Create(testType)
	obstacle2, _ := factory.Create(testType)

	if obstacle1 == obstacle2 {
		t.Error("expected each Create call to return a new instance")
	}
}

func TestDefaultObstacleFactory_ImplementsInterface(t *testing.T) {
	factory := NewDefaultObstacleFactory()

	// Verify it implements ObstacleFactory interface
	var _ ObstacleFactory = factory
}

func TestObstacleType_Constants(t *testing.T) {
	// Just verify ObstacleType can be used
	var _ ObstacleType = 0
	var _ ObstacleType = 1
}

func TestObstacleMap_Type(t *testing.T) {
	// Verify ObstacleMap type
	var _ ObstacleMap = make(map[ObstacleType]func() body.Obstacle)
}
