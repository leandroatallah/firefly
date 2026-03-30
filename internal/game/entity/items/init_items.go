package gameitems

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
)

const (
	FallingPlatformType items.ItemType = "FALL_PLATFORM"
	FreezePowerUpType   items.ItemType = "FREEZE_POWER_UP"
	GrowPowerUpType     items.ItemType = "GROW_POWER_UP"
	StarPowerUpType     items.ItemType = "STAR_POWER_UP"
)

func itemFactoryOrFatal(item items.Item, err error) items.Item {
	if err != nil {
		log.Fatal(err)
	}
	return item
}

func InitItemMap(ctx *app.AppContext) items.ItemMap[items.Item] {
	itemMap := map[items.ItemType]func(x, y int, id string) items.Item{
		FallingPlatformType: func(x, y int, id string) items.Item {
			return itemFactoryOrFatal(NewFallingPlatformItem(ctx, x, y, id))
		},
		FreezePowerUpType: func(x, y int, id string) items.Item {
			return itemFactoryOrFatal(NewFreezePowerItem(ctx, x, y, id))
		},
		GrowPowerUpType: func(x, y int, id string) items.Item {
			return itemFactoryOrFatal(NewGrowPowerItem(ctx, x, y, id))
		},
		StarPowerUpType: func(x, y int, id string) items.Item {
			return itemFactoryOrFatal(NewStarPowerItem(ctx, x, y, id))
		},
	}
	return itemMap
}
