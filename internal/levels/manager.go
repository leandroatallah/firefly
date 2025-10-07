package levels

import "fmt"

type Manager struct {
	levels       map[int]Level
	CurrentLevel int
}

func NewManager() *Manager {
	return &Manager{
		levels: make(map[int]Level),
	}
}

func (m *Manager) AddLevel(level Level) {
	m.levels[level.ID] = level
}

func (m *Manager) GetLevel(id int) (Level, error) {
	level, ok := m.levels[id]
	if !ok {
		return Level{}, fmt.Errorf("level with id %d not found", id)
	}
	return level, nil
}

func (m *Manager) GetCurrentLevel() (Level, error) {
	return m.GetLevel(m.CurrentLevel)
}

func (m *Manager) SetCurrentLevel(id int) error {
	_, err := m.GetLevel(id)
	if err != nil {
		return err
	}
	m.CurrentLevel = id
	return nil
}

func (m *Manager) AdvanceToNextLevel() error {
	level, err := m.GetCurrentLevel()
	if err != nil {
		return err
	}

	if level.NextLevelID == 0 {
		return fmt.Errorf("no next level defined for level %d", level.ID)
	}

	return m.SetCurrentLevel(level.NextLevelID)
}
