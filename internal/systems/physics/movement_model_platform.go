package physics

type PlatformMovementModel struct{}

func NewPlatformMovementModel() *PlatformMovementModel {
	return &PlatformMovementModel{}
}

func (m *PlatformMovementModel) Update(body *PhysicsBody, space *Space) error {
	return nil
}

func (m *PlatformMovementModel) InputHandler(body *PhysicsBody) {
}
