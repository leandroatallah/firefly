package gamescenephases

type deathSequencePhase int

const (
	deathSequencePhaseWaiting deathSequencePhase = iota
	deathSequencePhaseMoving
)

// deathSequence holds all state for the death and respawn sequence.
type deathSequence struct {
	active        bool
	phase         deathSequencePhase
	waitTimer     int
	waitDuration  int
	playerStartX  float64
	playerStartY  float64
	cameraStartX  float64
	cameraStartY  float64
	cameraTargetX float64
	cameraTargetY float64
	timer         int
	duration      int
}
