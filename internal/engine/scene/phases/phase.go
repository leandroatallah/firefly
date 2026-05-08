package phases

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"

// Genre is a dumb int type identifying the gameplay genre of a phase.
// The engine stores the value but assigns no meaning to it; the kit layer
// owns the named constants (GenrePlatformer, GenreBeatemup, …).
type Genre int

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
	Genre               Genre
	BlockPlayerMovement bool
}
