// Package gameplatformerphase wires the kit PlatformerPhaseScene with
// game-specific factories (player, enemies, items, NPCs).
package gameplatformerphase

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/enemies"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/npcs"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	gameenemies "github.com/boilerplate/ebiten-template/internal/game/entity/actors/enemies"
	gamenpcs "github.com/boilerplate/ebiten-template/internal/game/entity/actors/npcs"
	gameitems "github.com/boilerplate/ebiten-template/internal/game/entity/items"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
	platformerphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/platformer"
)

// NewPlatformerPhaseScene constructs the kit scene wired with game-specific factories.
func NewPlatformerPhaseScene(ctx *app.AppContext) *platformerphasescene.PlatformerPhaseScene {
	s, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
		Ctx:           ctx,
		PlayerFactory: newClimberPlayer,
		InitActors: func(ts *scene.TilemapScene) {
			scene.InitItems(ts, items.NewItemFactory(gameitems.InitItemMap(ctx)))
			scene.InitEnemies(ts, enemies.NewEnemyFactory(gameenemies.InitEnemyMap(ctx)))
			scene.InitNPCs(ts, npcs.NewNpcFactory(gamenpcs.InitNpcMap(ctx)))
		},
		DebugDrawHook:   makeClimberDebugHook(ctx),
		RebootSceneType: scenestypes.ScenePhaseReboot,
		MenuSceneType:   scenestypes.SceneMenu,
	})
	if err != nil {
		log.Fatal(err)
	}
	return s
}
