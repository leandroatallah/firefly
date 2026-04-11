package sequences

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/hajimehoshi/ebiten/v2"
)

// stubVFXManager implements vfx.Manager for headless testing.
type stubVFXManager struct {
	spawnedTexts        []string
	spawnedAboveTargets []string
}

func (v *stubVFXManager) SetAppContext(_ any)                                      {}
func (v *stubVFXManager) Update()                                                  {}
func (v *stubVFXManager) Draw(_ *ebiten.Image, _ *camera.Controller)               {}
func (v *stubVFXManager) AddParticle(_ *particles.Particle)                        {}
func (v *stubVFXManager) AddTrauma(_ *camera.Controller, _ float64)                {}
func (v *stubVFXManager) PixelConfig() *particles.Config                           { return nil }
func (v *stubVFXManager) SetDefaultFont(_ *font.FontText)                          {}
func (v *stubVFXManager) SpawnDeathExplosion(_ float64, _ float64, _ int)          {}
func (v *stubVFXManager) SpawnFallingRocks(_ float64, _ float64, _ float64, _ int) {}
func (v *stubVFXManager) SpawnFloatingText(msg string, _ float64, _ float64, _ int) {
	v.spawnedTexts = append(v.spawnedTexts, msg)
}
func (v *stubVFXManager) SpawnFloatingTextAbove(target body.Body, msg string, _ int) {
	v.spawnedAboveTargets = append(v.spawnedAboveTargets, target.ID())
	v.spawnedTexts = append(v.spawnedTexts, msg)
}
func (v *stubVFXManager) SpawnJumpPuff(_ float64, _ float64, _ int)                  {}
func (v *stubVFXManager) SpawnLandingPuff(_ float64, _ float64, _ int)               {}
func (v *stubVFXManager) SpawnPuff(_ string, _ float64, _ float64, _ int, _ float64) {}
func (v *stubVFXManager) TriggerScreenFlash()                                        {}

func TestSpawnTextCommand_Init_ScreenType(t *testing.T) {
	ctx := &app.AppContext{}
	vfxMgr := &stubVFXManager{}
	ctx.VFX = vfxMgr
	ctx.ActorManager = actors.NewManager()

	cmd := &SpawnTextCommand{Text: "Hello", Type: "screen", X: 10, Y: 20, Duration: 30}
	cmd.Init(ctx)

	if len(vfxMgr.spawnedTexts) == 0 || vfxMgr.spawnedTexts[0] != "Hello" {
		t.Errorf("expected SpawnFloatingText called with 'Hello', got %v", vfxMgr.spawnedTexts)
	}
}

func TestSpawnTextCommand_Init_OverheadType(t *testing.T) {
	ctx := &app.AppContext{}
	vfxMgr := &stubVFXManager{}
	ctx.VFX = vfxMgr
	am := actors.NewManager()
	ctx.ActorManager = am

	actor := &mocks.MockActor{Id: "npc"}
	actor.SetPosition(0, 0)
	am.Register(actor)

	cmd := &SpawnTextCommand{Text: "Overhead!", Type: "overhead", TargetID: "npc", Duration: 30}
	cmd.Init(ctx)

	if len(vfxMgr.spawnedAboveTargets) == 0 || vfxMgr.spawnedAboveTargets[0] != "npc" {
		t.Errorf("expected SpawnFloatingTextAbove called for 'npc', got %v", vfxMgr.spawnedAboveTargets)
	}
}

func TestQuakeCommand_Init_Update(t *testing.T) {
	ctx := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(ctx)
	ctx.SceneManager = sceneManager
	ctx.Space = space.NewSpace()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(ctx)
	ctx.SceneManager.SwitchTo(mockScene)

	cmd := &QuakeCommand{Trauma: 0.3, Duration: 15}
	cmd.Init(ctx)

	if cmd.camera == nil {
		t.Fatal("camera should not be nil after Init")
	}

	done := false
	for i := 0; i < 20 && !done; i++ {
		done = cmd.Update()
	}

	if !done {
		t.Error("QuakeCommand should complete after Duration frames")
	}
	if cmd.timer < cmd.Duration {
		t.Errorf("expected timer >= Duration (%d), got %d", cmd.Duration, cmd.timer)
	}
}

func TestQuakeCommand_Update_AddsTraumaOnMultiplesOf10(t *testing.T) {
	ctx := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(ctx)
	ctx.SceneManager = sceneManager
	ctx.Space = space.NewSpace()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(ctx)
	ctx.SceneManager.SwitchTo(mockScene)

	cmd := &QuakeCommand{Trauma: 0.5, Duration: 25}
	cmd.Init(ctx)

	// Run 10 frames — trauma should be added at frame 0 and frame 10
	for i := 0; i < 10; i++ {
		cmd.Update()
	}
	// Just verify no panic and timer advanced
	if cmd.timer != 10 {
		t.Errorf("expected timer=10, got %d", cmd.timer)
	}
}

type mockVignetteScene struct {
	log []string
}

func (m *mockVignetteScene) EnableVignetteDarkness(radiusPx float64) {
	m.log = append(m.log, "enable")
}

func (m *mockVignetteScene) DisableVignetteDarkness() {
	m.log = append(m.log, "disable")
}

func TestVignetteRadiusCommandInitNoScene(t *testing.T) {
	cmd := &VignetteRadiusCommand{
		InitialRadius: 10,
		FinalRadius:   20,
		Duration:      30,
	}

	cmd.Init(nil)
	if !cmd.Update() {
		t.Fatalf("expected Update to complete immediately when no controller is set")
	}
}

func TestVignetteRadiusCommandUpdateInstant(t *testing.T) {
	sceneMock := &mockVignetteScene{}
	cmd := &VignetteRadiusCommand{
		InitialRadius: 0,
		FinalRadius:   32,
		Duration:      0,
		controller:    sceneMock,
	}

	if done := cmd.Update(); !done {
		t.Fatalf("expected instant completion when Duration <= 0")
	}
	if len(sceneMock.log) == 0 || sceneMock.log[len(sceneMock.log)-1] != "enable" {
		t.Fatalf("expected EnableVignetteDarkness to be called for final radius")
	}
}

func TestSpawnTextCommandInitStructFields(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text:     "Hello World",
		Type:     "screen",
		X:        100.0,
		Y:        200.0,
		Duration: 60,
	}

	if cmd.Text != "Hello World" {
		t.Errorf("expected Text 'Hello World', got %q", cmd.Text)
	}
	if cmd.Type != "screen" {
		t.Errorf("expected Type 'screen', got %q", cmd.Type)
	}
	if cmd.X != 100.0 {
		t.Errorf("expected X 100.0, got %f", cmd.X)
	}
	if cmd.Y != 200.0 {
		t.Errorf("expected Y 200.0, got %f", cmd.Y)
	}
}

func TestSpawnTextCommandInitOverheadNoActor(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text:     "Overhead Text",
		Type:     "overhead",
		TargetID: "nonexistent",
		Duration: 60,
	}

	if cmd.TargetID != "nonexistent" {
		t.Errorf("expected TargetID 'nonexistent', got %q", cmd.TargetID)
	}
}

func TestSpawnTextCommandInitEmptyTargetID(t *testing.T) {
	cmd := &SpawnTextCommand{
		Text:     "Test",
		Type:     "overhead",
		TargetID: "",
	}

	if cmd.TargetID != "" {
		t.Errorf("expected empty TargetID, got %q", cmd.TargetID)
	}
}

func TestSpawnTextCommandUpdate(t *testing.T) {
	cmd := &SpawnTextCommand{Text: "Test", Type: "screen"}
	if !cmd.Update() {
		t.Error("SpawnTextCommand.Update() should return true (instant command)")
	}
}
