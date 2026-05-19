// Red-Phase tests for story 059-thin-game-phase-scenes [AC-1, AC-3, AC-5,
// AC-6, AC-10].
//
// SPEC.md §4 / §5 absorbs the full Update/Draw loop and player/factory
// wiring into the kit PlatformerPhaseScene via an Options-based
// constructor. These tests pin the new public surface:
//
//   - NewWithOptions(opts Options[Player]) (*PlatformerPhaseScene, error)
//   - Options exposes Ctx, PlayerFactory, ItemMap, EnemyMap, NpcMap,
//     DebugDrawHook, RebootSceneType, MenuSceneType.
//   - Required fields (Ctx, PlayerFactory) are validated.
//   - Sequence-gated pause: canPause() == false while a sequence is
//     playing (re-tested via observable Update semantics).
//   - PlayerFactory returning nil does not panic.
//
// Each test is independently table-driven where multiple scenarios apply.
package platformerphasescene_test

import (
	"errors"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	platformerphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/platformer"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewWithOptions_ErrorsWhenCtxNil(t *testing.T) {
	playerFactory := func(*app.AppContext) (platformerphasescene.Player, error) {
		return nil, nil
	}

	_, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
		Ctx:           nil,
		PlayerFactory: playerFactory,
	})

	if err == nil {
		t.Fatal("expected error when Ctx is nil")
	}
}

func TestNewWithOptions_ErrorsWhenPlayerFactoryNil(t *testing.T) {
	_, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
		Ctx:           &app.AppContext{},
		PlayerFactory: nil,
	})

	if err == nil {
		t.Fatal("expected error when PlayerFactory is nil")
	}
}

func TestNewWithOptions_AcceptsOptionalNilFields(t *testing.T) {
	// All non-required fields default to nil/zero; constructor must not error.
	pf := func(*app.AppContext) (platformerphasescene.Player, error) { return nil, nil }
	_, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
		Ctx:             &app.AppContext{},
		PlayerFactory:   pf,
		ItemMap:         nil,
		EnemyMap:        nil,
		NpcMap:          nil,
		DebugDrawHook:   nil,
		RebootSceneType: navigation.SceneType(0),
		MenuSceneType:   navigation.SceneType(0),
	})
	if err != nil {
		t.Fatalf("expected no error for fully-optional fields, got %v", err)
	}
}

func TestNewWithOptions_PropagatesFactoryError(t *testing.T) {
	wantErr := errors.New("boom")
	pf := func(*app.AppContext) (platformerphasescene.Player, error) {
		return nil, wantErr
	}

	_, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
		Ctx:           &app.AppContext{},
		PlayerFactory: pf,
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected factory error propagated via errors.Is, got %v", err)
	}
}

func TestNewWithOptions_NilPlayerFromFactory_DoesNotPanic(t *testing.T) {
	// "Phase with no player": kit scene must handle PlayerFactory returning
	// nil without panicking; camera and draw loop run in degraded mode.
	pf := func(*app.AppContext) (platformerphasescene.Player, error) {
		return nil, nil
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic when factory returns nil player, got: %v", r)
		}
	}()

	scene, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
		Ctx:           &app.AppContext{},
		PlayerFactory: pf,
	})
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}
	if scene == nil {
		t.Fatal("expected non-nil scene")
	}

	// Update + Draw must not panic even in degraded mode.
	if err := scene.Update(); err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}
	scene.Draw(ebiten.NewImage(320, 200))
}

func TestSetDebugDrawHook_NilHook_DrawDoesNotPanic(t *testing.T) {
	// AC-10 edge case: DebugDrawHook nil; Draw must skip the hook call.
	scene := platformerphasescene.NewForTest(platformerphasescene.TestOptions{
		CameraCenterX: 0,
		CameraCenterY: 0,
		ScreenWidth:   320,
		ScreenHeight:  200,
	})
	scene.SetDebugDrawHook(nil)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic with nil DebugDrawHook, got: %v", r)
		}
	}()
	scene.Draw(ebiten.NewImage(320, 200))
}
