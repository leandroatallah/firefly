package phases

type Phase struct {
	ID          int
	Name        string
	TilemapPath string
	NextPhaseID int
}
