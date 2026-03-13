package scene

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/tilemap"
)

func TestScreenFlipper_GridGeneration(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	cam := camera.NewController(0, 0)
	player := &mocks.MockActor{Id: "player"}
	tm := &tilemap.Tilemap{
		Width:      20,
		Height:     20,
		Tilewidth:  16,
		Tileheight: 16,
	} // 320x320 map
	ctx := &app.AppContext{}

	sf := NewScreenFlipper(cam, player, tm, ctx)
	sf.ensureRooms()

	// mw=320, mh=320. cols=ceil(320/320)=1, rows=ceil(320/240)=2.
	expectedRooms := 2
	if len(sf.rooms) != expectedRooms {
		t.Errorf("expected %d rooms; got %d", expectedRooms, len(sf.rooms))
	}
}

func TestScreenFlipper_SnapAndTrigger(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240, ScreenFlipSpeed: 0.1})
	cam := camera.NewController(0, 0)
	player := &mocks.MockActor{Id: "player"}
	player.SetPosition(160, 120)
	
	tm := &tilemap.Tilemap{
		Width:      40, // 640 wide
		Height:     20, // 320 high
		Tilewidth:  16,
		Tileheight: 16,
	}
	ctx := &app.AppContext{}

	sf := NewScreenFlipper(cam, player, tm, ctx)
	
	// Test Snap
	sf.SnapToCurrentRoom()
	if sf.currentRoom == nil || sf.currentRoom.Min.X != 0 {
		t.Fatal("failed to snap to initial room")
	}

	// Move player to the right edge
	player.SetPosition(330, 120)
	sf.Update() // checkTrigger -> triggerFlip
	
	if !sf.IsFlipping() {
		t.Error("expected IsFlipping to be true after trigger")
	}
	
	// Update flip to completion
	for sf.IsFlipping() {
		sf.Update()
	}
	
	if sf.currentRoom.Min.X != 320 {
		t.Errorf("expected current room to be second room; got minX=%d", sf.currentRoom.Min.X)
	}
}

func TestScreenFlipper_InstantFlip(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	cam := camera.NewController(0, 0)
	player := &mocks.MockActor{Id: "player"}
	player.SetPosition(160, 120)
	
	tm := &tilemap.Tilemap{
		Width: 40, Height: 20, Tilewidth: 16, Tileheight: 16,
	}
	ctx := &app.AppContext{}

	sf := NewScreenFlipper(cam, player, tm, ctx)
	sf.FlipStrategy = func(dx, dy int) FlipType { return FlipTypeInstant }
	
	// Start in first room
	sf.SnapToCurrentRoom()
	
	// Trigger flip
	player.SetPosition(330, 120)
	sf.Update()
	
	if sf.IsFlipping() {
		t.Error("expected instant flip to not set isFlipping=true")
	}
	if sf.currentRoom.Min.X != 320 {
		t.Error("failed to flip to next room instantly")
	}
}

func TestScreenFlipper_Hooks(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240, ScreenFlipSpeed: 1.0})
	cam := camera.NewController(0, 0)
	player := &mocks.MockActor{Id: "player"}
	tm := &tilemap.Tilemap{
		Width: 40, Height: 20, Tilewidth: 16, Tileheight: 16,
	}
	
	sf := NewScreenFlipper(cam, player, tm, nil)
	
	startCalled := false
	finishCalled := false
	sf.OnFlipStart = func() { startCalled = true }
	sf.OnFlipFinish = func() { finishCalled = true }
	
	sf.SnapToCurrentRoom()
	player.SetPosition(330, 120)
	sf.Update() // trigger
	
	if !startCalled {
		t.Error("OnFlipStart not called")
	}
	
	sf.Update() // complete (speed 1.0)
	if !finishCalled {
		t.Error("OnFlipFinish not called")
	}
}

func TestScreenFlipper_NilChecks(t *testing.T) {
	sf := NewScreenFlipper(nil, nil, nil, nil)
	// Should not panic
	sf.Update()
	sf.SnapToCurrentRoom()
}

func TestScreenFlipper_VerticalFlip(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240, ScreenFlipSpeed: 1.0})
	cam := camera.NewController(0, 0)
	player := &mocks.MockActor{Id: "player"}
	tm := &tilemap.Tilemap{
		Width: 20, Height: 40, Tilewidth: 16, Tileheight: 16,
	} // 320x640 map
	
	sf := NewScreenFlipper(cam, player, tm, nil)
	sf.SnapToCurrentRoom()
	
	// Flip Down
	player.SetPosition(160, 250)
	sf.Update() // trigger
	sf.Update() // complete
	
	if sf.currentRoom.Min.Y != 240 {
		t.Errorf("vertical flip failed; got minY=%d", sf.currentRoom.Min.Y)
	}
	
	// Flip Up
	player.SetPosition(160, 230)
	sf.Update() // trigger
	sf.Update() // complete
	if sf.currentRoom.Min.Y != 0 {
		t.Errorf("vertical flip up failed; got minY=%d", sf.currentRoom.Min.Y)
	}
}
