package enemies

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
)

type EnemyData struct {
	SpriteData schemas.SpriteData `json:"sprites"`
	StatData   actors.StatData    `json:"stats"`
}

func ParseJsonEnemy(path string) (schemas.SpriteData, actors.StatData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schemas.SpriteData{}, actors.StatData{}, err
	}

	var enemyData EnemyData
	if err := json.Unmarshal(data, &enemyData); err != nil {
		return schemas.SpriteData{}, actors.StatData{}, err
	}

	return enemyData.SpriteData, enemyData.StatData, nil
}
