package items

import "fmt"

type ItemFactory[T Item] struct {
	itemMap ItemMap[T]
}

func NewItemFactory[T Item](itemMap ItemMap[T]) *ItemFactory[T] {
	return &ItemFactory[T]{itemMap: itemMap}
}

func (f *ItemFactory[T]) Create(itemType ItemType, x, y int, id string) (Item, error) {
	itemFunc, ok := f.itemMap[itemType]
	if !ok {
		return nil, fmt.Errorf("unknown item type")
	}

	item := itemFunc(x, y, id)

	return item, nil
}
