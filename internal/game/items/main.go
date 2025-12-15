package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/items"
)

const (
	CollectibleCoinType items.ItemType = iota
	SignpostType
)

func InitItemMap(ctx *core.AppContext) items.ItemMap {
	itemMap := map[items.ItemType]func(x, y int) items.Item{
		CollectibleCoinType: func(x, y int) items.Item {
			return NewCollectibleCoinItem(ctx, x, y)
		},
		SignpostType: func(x, y int) items.Item {
			return NewSignpostItem(ctx, x, y)
		},
	}
	return itemMap
}
