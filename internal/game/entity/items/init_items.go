package gameitems

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
)

const (
	FallingPlatformType items.ItemType = "FALL_PLATFORM"
	WeaponCannonType    items.ItemType = "ITEM_WEAPON_CANNON"
	GrowPowerUpType     items.ItemType = "GROW_POWER_UP"
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
		WeaponCannonType: func(x, y int, id string) items.Item {
			return itemFactoryOrFatal(NewWeaponCannonItem(ctx, x, y, id))
		},
		GrowPowerUpType: func(x, y int, id string) items.Item {
			return itemFactoryOrFatal(NewPowerUpItem(ctx, x, y, id, "assets/entities/items/item-power-grow.json", func() {
			}))
		},
	}
	return itemMap
}
