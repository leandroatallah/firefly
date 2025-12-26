package gameitems

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/items"
)

const (
	CollectibleCoinType items.ItemType = iota
)

func InitItemMap(ctx *core.AppContext) items.ItemMap {
	itemMap := map[items.ItemType]func(x, y int, id string) items.Item{
		CollectibleCoinType: func(x, y int, id string) items.Item {
			item, err := NewCollectibleCoinItem(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return item
		},
	}
	return itemMap
}
