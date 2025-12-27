package items

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
)

type StatData struct {
	Id string `json:"id"`
}

type ItemData struct {
	SpriteData schemas.SpriteData `json:"sprites"`
	StatData   StatData           `json:"stats"`
}

func ParseJsonItem(path string) (schemas.SpriteData, StatData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schemas.SpriteData{}, StatData{}, err
	}

	var itemData ItemData
	if err := json.Unmarshal(data, &itemData); err != nil {
		return schemas.SpriteData{}, StatData{}, err
	}

	return itemData.SpriteData, itemData.StatData, nil
}
