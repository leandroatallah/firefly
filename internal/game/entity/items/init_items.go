package gameitems

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
)

const (
	FallingPlatformType items.ItemType = "FALL_PLATFORM"
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
	}
	return itemMap
}
