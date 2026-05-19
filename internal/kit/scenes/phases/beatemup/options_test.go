// Red-Phase tests for story 059-thin-game-phase-scenes [AC-2, AC-4, AC-5,
// AC-6].
//
// SPEC.md §4 / §6 absorbs the full Update/Draw loop and player/factory
// wiring into the kit BeatemupPhaseScene via an Options-based
// constructor. These tests pin the new public surface and confirm:
//
//   - NewWithOptions(opts Options[Player]) (*BeatemupPhaseScene, error)
//   - Options exposes Ctx, PlayerFactory, ItemMap, EnemyMap, NpcMap,
//     DebugDrawHook, RebootSceneType, MenuSceneType.
//   - Beat-em-up scenes do NOT perform fall-death checks (AC-2).
//   - PlayerFactory returning nil does not panic.
package beatemupphasescene_test

import (
	"errors"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	beatemupphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/beatemup"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestBeatemupNewWithOptions_ErrorsWhenCtxNil(t *testing.T) {
	pf := func(*app.AppContext) (beatemupphasescene.Player, error) { return nil, nil }

	_, err := beatemupphasescene.NewWithOptions(beatemupphasescene.Options[beatemupphasescene.Player]{
		Ctx:           nil,
		PlayerFactory: pf,
	})
	if err == nil {
		t.Fatal("expected error when Ctx is nil")
	}
}

func TestBeatemupNewWithOptions_ErrorsWhenPlayerFactoryNil(t *testing.T) {
	_, err := beatemupphasescene.NewWithOptions(beatemupphasescene.Options[beatemupphasescene.Player]{
		Ctx:           &app.AppContext{},
		PlayerFactory: nil,
	})
	if err == nil {
		t.Fatal("expected error when PlayerFactory is nil")
	}
}

func TestBeatemupNewWithOptions_PropagatesFactoryError(t *testing.T) {
	wantErr := errors.New("boom")
	pf := func(*app.AppContext) (beatemupphasescene.Player, error) {
		return nil, wantErr
	}

	_, err := beatemupphasescene.NewWithOptions(beatemupphasescene.Options[beatemupphasescene.Player]{
		Ctx:           &app.AppContext{},
		PlayerFactory: pf,
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected factory error propagated via errors.Is, got %v", err)
	}
}

func TestBeatemupNewWithOptions_NilPlayer_NoPanicOnUpdateDraw(t *testing.T) {
	pf := func(*app.AppContext) (beatemupphasescene.Player, error) {
		return nil, nil
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic when factory returns nil, got: %v", r)
		}
	}()

	scene, err := beatemupphasescene.NewWithOptions(beatemupphasescene.Options[beatemupphasescene.Player]{
		Ctx:           &app.AppContext{},
		PlayerFactory: pf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scene == nil {
		t.Fatal("expected non-nil scene")
	}

	if err := scene.Update(); err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}
	scene.Draw(ebiten.NewImage(320, 200))
}

// AC-2: Beat-em-up Update must NOT perform fall-death checks. Even if the
// (hypothetical) player is positioned far below the camera, the scene
// must not flip a death flag from a fall-check pathway.
func TestBeatemupUpdate_DoesNotPerformFallDeathCheck(t *testing.T) {
	scene := beatemupphasescene.NewForTest(beatemupphasescene.TestOptions{
		CameraCenterX:          0,
		CameraCenterY:          0,
		ScreenWidth:            320,
		ScreenHeight:           200,
		HasPlayerStartPosition: false,
	})

	if err := scene.Update(); err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}
	if scene.DeathActiveForTest() {
		t.Fatal("expected deathActive==false; beatemup Update must not run fall-death")
	}
}
