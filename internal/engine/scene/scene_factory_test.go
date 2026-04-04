package scene

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestDefaultSceneFactory_Create(t *testing.T) {
	ctx := &app.AppContext{}
	sceneType := navigation.SceneType(1)

	sceneMap := navigation.SceneMap{
		sceneType: func() navigation.Scene {
			return &mocks.MockScene{}
		},
	}

	factory := NewDefaultSceneFactory(sceneMap)
	factory.SetAppContext(ctx)

	// Test 1: Successful creation
	s1, err := factory.Create(sceneType, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s1 == nil {
		t.Fatal("expected scene to be created")
	}
	mockS1 := s1.(*mocks.MockScene)
	if mockS1.AppContextSet != ctx {
		t.Error("AppContext not set on created scene")
	}

	// Test 2: Cached instance
	s2, _ := factory.Create(sceneType, false)
	if s2 != s1 {
		t.Error("expected cached instance when freshInstance=false")
	}

	// Test 3: Fresh instance
	s3, _ := factory.Create(sceneType, true)
	if s3 == s1 {
		t.Error("expected new instance when freshInstance=true")
	}

	// Test 4: Unknown scene type
	_, err = factory.Create(999, false)
	if err == nil {
		t.Error("expected error for unknown scene type")
	}
}
