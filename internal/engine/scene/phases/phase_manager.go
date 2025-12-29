package phases

import "fmt"

type Manager struct {
	phases       map[int]Phase
	CurrentPhase int
}

func NewManager() *Manager {
	return &Manager{
		phases: make(map[int]Phase),
	}
}

func (m *Manager) AddPhase(p Phase) {
	m.phases[p.ID] = p
}

func (m *Manager) GetPhase(id int) (Phase, error) {
	p, ok := m.phases[id]
	if !ok {
		return Phase{}, fmt.Errorf("phase with id %d not found", id)
	}
	return p, nil
}

func (m *Manager) GetCurrentPhase() (Phase, error) {
	return m.GetPhase(m.CurrentPhase)
}

func (m *Manager) SetCurrentPhase(id int) error {
	_, err := m.GetPhase(id)
	if err != nil {
		return err
	}
	m.CurrentPhase = id
	return nil
}

func (m *Manager) AdvanceToNextPhase() error {
	p, err := m.GetCurrentPhase()
	if err != nil {
		return err
	}

	if p.NextPhaseID == 0 {
		return fmt.Errorf("no next phase defined for phase %d", p.ID)
	}

	return m.SetCurrentPhase(p.NextPhaseID)
}
