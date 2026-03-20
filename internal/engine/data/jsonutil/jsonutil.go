package jsonutil

import (
	"encoding/json"
	"io/fs"

	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
)

type SpriteAndStats[T any] struct {
	SpriteData schemas.SpriteData `json:"sprites"`
	StatData   T                  `json:"stats"`
}

func ParseSpriteAndStats[T any](fsys fs.FS, path string) (schemas.SpriteData, T, error) {
	var zero T
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return schemas.SpriteData{}, zero, err
	}
	var payload SpriteAndStats[T]
	if err := json.Unmarshal(data, &payload); err != nil {
		return schemas.SpriteData{}, zero, err
	}
	return payload.SpriteData, payload.StatData, nil
}
