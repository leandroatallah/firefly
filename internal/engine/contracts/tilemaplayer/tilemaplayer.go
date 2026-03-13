package tilemaplayer

import "image"

type TilemapDimensionsProvider interface {
	GetTilemapWidth() int
	GetTilemapHeight() int
	GetCameraBounds() (image.Rectangle, bool)
}
