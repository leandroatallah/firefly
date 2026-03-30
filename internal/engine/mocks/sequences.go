package mocks

import contractseq "github.com/boilerplate/ebiten-template/internal/engine/contracts/sequences"

// MockCommand implements sequences.Command
type MockCommand struct {
	InitCalled    bool
	UpdateCount   int
	CompleteAfter int
}

func (c *MockCommand) Init(appContext any) {
	c.InitCalled = true
}

func (c *MockCommand) Update() bool {
	c.UpdateCount++
	return c.UpdateCount >= c.CompleteAfter
}

// MockSequence implements sequences.Sequence
type MockSequence struct {
	CommandsList  []contractseq.Command
	IsInterruptible bool
	IsOneTime       bool
	Path            string
}

func (s *MockSequence) Commands() []contractseq.Command {
	return s.CommandsList
}

func (s *MockSequence) Interruptible() bool {
	return s.IsInterruptible
}

func (s *MockSequence) OneTime() bool {
	return s.IsOneTime
}

func (s *MockSequence) GetPath() string {
	if s.Path != "" {
		return s.Path
	}
	return "mock_sequence"
}
