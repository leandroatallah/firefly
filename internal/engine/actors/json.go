package actors

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/schemas"
)

type StatData struct {
	Health   int `json:"health"`
	Speed    int `json:"speed"`
	MaxSpeed int `json:"max_speed"`
}

type PlayerData struct {
	SpriteData schemas.SpriteData `json:"sprites"`
	StatData   StatData           `json:"stats"`
}

func ParseJsonPlayer(path string) (schemas.SpriteData, StatData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schemas.SpriteData{}, StatData{}, err
	}

	var playerData PlayerData
	if err := json.Unmarshal(data, &playerData); err != nil {
		return schemas.SpriteData{}, StatData{}, err
	}

	return playerData.SpriteData, playerData.StatData, nil
}
