package phases

import "github.com/leandroatallah/firefly/internal/engine/contracts/navigation"

type GoalType string

type Phase struct {
	ID                  int
	Name                string
	Title               string
	TilemapPath         string
	NextPhaseID         int
	SequencePath        string
	GoalType            GoalType
	SceneType           navigation.SceneType
	BlockPlayerMovement bool
}
