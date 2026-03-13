package gameitems

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
)

const (
	FallingPlatformType items.ItemType = "FALL_PLATFORM"
)

func InitItemMap(ctx *app.AppContext) items.ItemMap[items.Item] {
	itemMap := map[items.ItemType]func(x, y int, id string) items.Item{
		FallingPlatformType: func(x, y int, id string) items.Item {
			item, err := NewFallingPlatformItem(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return item
		},
	}
	return itemMap
}
