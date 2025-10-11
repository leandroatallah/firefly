package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/items"
)

const (
	CollectibleCoinType items.ItemType = iota
	SignpostType
)

func InitItemMap() items.ItemMap {
	enemyMap := map[items.ItemType]func(x, y int) items.Item{
		CollectibleCoinType: func(x, y int) items.Item {
			return NewCollectibleCoinItem(x, y)
		},
		SignpostType: func(x, y int) items.Item {
			return NewSignpostItem(x, y)
		},
	}
	return enemyMap
}
