package space

type TilemapDimensionsProvider interface {
	GetTilemapWidth() int
	GetTilemapHeight() int
}
