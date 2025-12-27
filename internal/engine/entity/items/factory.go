package items

import "fmt"

// To be initialized on game package.
type ItemType int
type ItemMap map[ItemType]func(x, y int, id string) Item

type ItemFactory struct {
	itemMap ItemMap
}

func NewItemFactory(itemMap ItemMap) *ItemFactory {
	return &ItemFactory{itemMap: itemMap}
}

func (f *ItemFactory) Create(itemType ItemType, x, y int, id string) (Item, error) {
	itemFunc, ok := f.itemMap[itemType]
	if !ok {
		return nil, fmt.Errorf("unknown item type")
	}

	item := itemFunc(x, y, id)

	return item, nil
}
