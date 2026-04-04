package sequences

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
)

type mockSceneWithCamera struct {
	mocks.MockScene
	cam *camera.Controller
}

func (m *mockSceneWithCamera) Camera() *camera.Controller {
	return m.cam
}

func TestCameraSetTargetCommand_SmoothTransition(t *testing.T) {
	// 1. Setup
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	target := &mocks.MockActor{Id: "test_target"}
	target.SetPosition(100, 100)
	appContext.Space.AddBody(target)

	cmd := &CameraSetTargetCommand{
		TargetID: "test_target",
		Duration: 60,
	}

	// 2. Init
	cmd.Init(appContext)

	if cmd.camera == nil {
		t.Fatal("camera should not be nil after Init")
	}
	if cmd.target == nil {
		t.Fatal("target should not be nil after Init")
	}
	if cmd.camera.IsFollowing() {
		t.Fatal("camera should not be following during smooth transition")
	}

	startX, startY := cmd.camera.Kamera().Center()
	if startX != 0 || startY != 0 {
		t.Fatalf("expected camera to start at (0,0), got (%f, %f)", startX, startY)
	}

	// 3. Update (simulate frames)
	finished := false
	for i := 0; i < 60; i++ {
		finished = cmd.Update()
		if finished {
			break
		}
	}

	if !finished {
		t.Fatal("command should have finished after 60 frames")
	}

	// 4. Verify final state
	if !cmd.camera.IsFollowing() {
		t.Fatal("camera should be following after transition is complete")
	}

	finalX, finalY := cmd.camera.Kamera().Center()
	targetX, targetY := target.GetPositionMin()
	w, h := target.Width(), target.Height()
	expectedX := float64(targetX) + float64(w)/2
	expectedY := float64(targetY) + float64(h)/2

	// Use a small tolerance for floating point comparison
	if finalX < expectedX-0.1 || finalX > expectedX+0.1 || finalY < expectedY-0.1 || finalY > expectedY+0.1 {
		t.Fatalf("expected camera to be centered on target (%f, %f), got (%f, %f)", expectedX, expectedY, finalX, finalY)
	}
}

func TestCameraSetTargetCommand_InstantTransition(t *testing.T) {
	// 1. Setup
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	target := &mocks.MockActor{Id: "test_target"}
	target.SetPosition(100, 100)
	appContext.Space.AddBody(target)

	cmd := &CameraSetTargetCommand{
		TargetID: "test_target",
		Duration: 0, // Instant
	}

	// 2. Init
	cmd.Init(appContext)

	if !cmd.camera.IsFollowing() {
		t.Fatal("camera should be following immediately on instant transition")
	}

	// 3. Update
	finished := cmd.Update()
	if !finished {
		t.Fatal("instant command should finish in one update")
	}

	// 4. Verify final state
	finalX, finalY := cmd.camera.Kamera().Center()
	targetX, targetY := target.GetPositionMin()
	w, h := target.Width(), target.Height()
	expectedX := float64(targetX) + float64(w)/2
	expectedY := float64(targetY) + float64(h)/2

	if finalX != expectedX || finalY != expectedY {
		t.Fatalf("expected camera to be centered on target (%f, %f), got (%f, %f)", expectedX, expectedY, finalX, finalY)
	}
}

func TestCameraMoveCommand_Init(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraMoveCommand{
		X:        100.0,
		Y:        200.0,
		Duration: 30,
	}

	cmd.Init(appContext)

	// Camera should be found via the interface assertion
	if cmd.camera == nil {
		t.Error("camera should not be nil")
	}
}

func TestCameraMoveCommand_Update_NilCamera(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	// No camera in this scene

	cmd := &CameraMoveCommand{
		X:        100.0,
		Y:        200.0,
		Duration: 10,
	}

	cmd.Init(appContext)

	finished := cmd.Update()
	if !finished {
		t.Error("command with nil camera should finish immediately")
	}
}

func TestCameraResetCommand_Init(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cam := camera.NewController(0, 0)
	cam.Kamera().ZoomFactor = 2.0 // Start zoomed in

	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraResetCommand{
		DefaultZoom: 1.0,
		Duration:    30,
	}

	cmd.Init(appContext)

	// Camera should be found via the interface assertion
	if cmd.camera == nil {
		t.Error("camera should not be nil")
	}
}

func TestCameraResetCommand_Update_NilCamera(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	// No camera

	cmd := &CameraResetCommand{
		DefaultZoom: 1.0,
		Duration:    10,
	}

	cmd.Init(appContext)

	finished := cmd.Update()
	if !finished {
		t.Error("command with nil camera should finish immediately")
	}
}

func TestCameraZoomCommand_Init_NilCamera(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cmd := &CameraZoomCommand{
		Zoom:     2.0,
		Duration: 30,
		Delay:    10,
	}

	// Should not panic with nil camera
	cmd.Init(appContext)
}

func TestCameraZoomCommand_Update_NilCamera(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cmd := &CameraZoomCommand{
		Zoom:     2.0,
		Duration: 30,
	}

	cmd.Init(appContext)

	finished := cmd.Update()
	if !finished {
		t.Error("command with nil camera should finish immediately")
	}
}

func TestCameraZoomCommand_WithTarget(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()
	appContext.ActorManager = actors.NewManager()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	target := &mocks.MockActor{Id: "zoom_target"}
	target.SetPosition(100, 100)
	appContext.Space.AddBody(target)
	appContext.ActorManager.Register(target)

	cmd := &CameraZoomCommand{
		Zoom:     2.0,
		Duration: 0, // Instant
		TargetID: "zoom_target",
	}

	// Camera should be found via the interface assertion
	cmd.Init(appContext)

	if cmd.camera == nil {
		t.Error("camera should not be nil")
	}
}

func TestCameraZoomCommand_AllPhases_InstantZoom(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()
	appContext.ActorManager = actors.NewManager()

	cam := camera.NewController(0, 0)
	cam.Kamera().ZoomFactor = 1.0

	// Give the camera a follow target so origFollowTarget is non-nil
	followTarget := &mocks.MockActor{Id: "follow_target"}
	followTarget.SetPosition(0, 0)
	cam.SetFollowTarget(followTarget)
	cam.SetFollowing(true)

	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraZoomCommand{
		Zoom:     2.0,
		Duration: 0, // instant zoom in
		Delay:    0, // instant wait
		// OutDuration defaults to Duration=0 → instant zoom out
	}
	cmd.Init(appContext)

	if cmd.camera == nil {
		t.Fatal("camera should not be nil after Init")
	}

	// Phase 0 (instant) → phase 1 → phase 2 (instant) → done
	done := cmd.Update()
	if !done {
		done = cmd.Update()
	}
	if !done {
		t.Error("CameraZoomCommand with Duration=0 and Delay=0 should complete quickly")
	}
}

func TestCameraZoomCommand_AllPhases_TimedZoom(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()
	appContext.ActorManager = actors.NewManager()

	cam := camera.NewController(0, 0)
	cam.Kamera().ZoomFactor = 1.0

	// Give the camera a follow target so origFollowTarget is non-nil
	followTarget := &mocks.MockActor{Id: "follow_target"}
	followTarget.SetPosition(0, 0)
	cam.SetFollowTarget(followTarget)
	cam.SetFollowing(true)

	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraZoomCommand{
		Zoom:        2.0,
		Duration:    3,
		Delay:       2,
		OutDuration: 3,
	}
	cmd.Init(appContext)

	// Drive through all three phases
	done := false
	for i := 0; i < 20 && !done; i++ {
		done = cmd.Update()
	}

	if !done {
		t.Error("CameraZoomCommand should complete after driving through all phases")
	}
	// Zoom should be restored to original
	if cam.Kamera().ZoomFactor != 1.0 {
		t.Errorf("expected zoom restored to 1.0, got %f", cam.Kamera().ZoomFactor)
	}
}

func TestCameraShakeCommand_InitAddsTrauma(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraShakeCommand{Trauma: 0.5}
	cmd.Init(appContext)

	if cmd.camera == nil {
		t.Fatal("camera should not be nil after Init")
	}
	// AddTrauma was called — Update should return true immediately
	if !cmd.Update() {
		t.Error("CameraShakeCommand.Update() should return true (instant command)")
	}
}

func TestCameraMoveCommand_Update_DurationPath(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraMoveCommand{X: 100, Y: 200, Duration: 5}
	cmd.Init(appContext)

	done := false
	for i := 0; i < 10 && !done; i++ {
		done = cmd.Update()
	}

	if !done {
		t.Error("CameraMoveCommand should complete after Duration frames")
	}
	x, y := cam.Kamera().Center()
	if x != 100 || y != 200 {
		t.Errorf("expected camera at (100,200), got (%f,%f)", x, y)
	}
}

func TestCameraResetCommand_Update_DurationPath(t *testing.T) {
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager

	cam := camera.NewController(0, 0)
	cam.Kamera().ZoomFactor = 2.0
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	cmd := &CameraResetCommand{DefaultZoom: 1.0, Duration: 5}
	cmd.Init(appContext)

	done := false
	for i := 0; i < 10 && !done; i++ {
		done = cmd.Update()
	}

	if !done {
		t.Error("CameraResetCommand should complete after Duration frames")
	}
	if cam.Kamera().ZoomFactor != 1.0 {
		t.Errorf("expected zoom 1.0, got %f", cam.Kamera().ZoomFactor)
	}
}
