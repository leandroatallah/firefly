package sequences

// Command is an action to be executed in a sequence.
// It can be initialized, and it is updated every frame until it is done.
type Command interface {
	// Init is called once when the command begins.
	// It can be used to set up initial state and get resources from the app context.
	Init(appContext any)

	// Update is called every frame.
	// It should return true when the command is finished.
	Update() bool
}

type Sequence interface {
	Commands() []Command
	Interruptible() bool
	OneTime() bool
	GetPath() string
}

type Player interface {
	IsPlaying() bool
	IsOver() bool
	Play(sequence Sequence)
	PlaySequence(filePath string)
	Stop()
	Update()
}
