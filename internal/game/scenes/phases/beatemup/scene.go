// Package gamebeatemupphase wires the kit BeatemupPhaseScene with
// game-specific factories (player, enemies, items, NPCs).
package gamebeatemupphase

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
	beatemupphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/beatemup"
)

// NewBeatemupPhaseScene constructs the kit scene wired with game-specific factories.
func NewBeatemupPhaseScene(ctx *app.AppContext) *beatemupphasescene.BeatemupPhaseScene {
	s, err := beatemupphasescene.NewWithOptions(beatemupphasescene.Options[beatemupphasescene.Player]{
		Ctx:           ctx,
		PlayerFactory: newCodyPlayer,
		InitActors: func(ts *scene.TilemapScene) {
			scene.InitItems(ts, items.NewItemFactory(gameitems.InitItemMap(ctx)))
			scene.InitEnemies(ts, enemies.NewEnemyFactory(gameenemies.InitEnemyMap(ctx)))
			scene.InitNPCs(ts, npcs.NewNpcFactory(gamenpcs.InitNpcMap(ctx)))
		},
		RebootSceneType: scenestypes.ScenePhaseReboot,
		MenuSceneType:   scenestypes.SceneMenu,
	})
	if err != nil {
		log.Fatal(err)
	}
	return s
}
