package beatemupphasescene

import (
	"errors"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/hajimehoshi/ebiten/v2"
)

// Options configures a BeatemupPhaseScene created via NewWithOptions.
// All fields except Ctx and PlayerFactory are optional.
type Options[T any] struct {
	// Ctx is the shared application context (required).
	Ctx *app.AppContext
	// PlayerFactory constructs the player on scene start (required).
	PlayerFactory func(*app.AppContext) (T, error)
	// ItemMap maps item types to their constructors; pass as any compatible ItemMap value (optional).
	ItemMap any
	// EnemyMap maps enemy types to their constructors; pass as any compatible EnemyMap value (optional).
	EnemyMap any
	// NpcMap maps NPC types to their constructors; pass as any compatible NpcMap value (optional).
	NpcMap any
	// DebugDrawHook is called at the end of Draw when non-nil (optional).
	DebugDrawHook func(*ebiten.Image)
	// RebootSceneType is the scene navigated to after the player dies (required for death routing).
	RebootSceneType navigation.SceneType
	// MenuSceneType is the scene navigated to when the player exits to the main menu (required for pause->menu).
	MenuSceneType navigation.SceneType
	// InitActors is an optional hook called after tilemap load to spawn entities.
	InitActors func(*scene.TilemapScene)
}

// NewWithOptions creates a BeatemupPhaseScene from the provided options.
// Returns an error if Ctx is nil, PlayerFactory is nil, or if the PlayerFactory
// returns an error.
func NewWithOptions(opts Options[Player]) (*BeatemupPhaseScene, error) {
	if opts.Ctx == nil {
		return nil, errors.New("beatemupphasescene: Options.Ctx must not be nil")
	}
	if opts.PlayerFactory == nil {
		return nil, errors.New("beatemupphasescene: Options.PlayerFactory must not be nil")
	}

	// Validate the factory by calling it now so the error propagates to the caller.
	_, err := opts.PlayerFactory(opts.Ctx)
	if err != nil {
		return nil, err
	}

	cfg := config.Get()
	sw := float64(cfg.ScreenWidth)
	sh := float64(cfg.ScreenHeight)

	cam := enginecamera.NewController(sw/2, sh/2)
	s := newScene(cam, nil, sw, sh)
	s.appCtx = opts.Ctx
	s.playerFactory = opts.PlayerFactory
	s.initActors = opts.InitActors
	s.rebootScene = opts.RebootSceneType
	s.menuScene = opts.MenuSceneType
	s.debugDrawHook = opts.DebugDrawHook

	return s, nil
}
