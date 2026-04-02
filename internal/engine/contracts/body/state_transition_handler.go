package body

type ShootDirection int

const (
	ShootDirectionStraight ShootDirection = iota
	ShootDirectionUp
	ShootDirectionDown
	ShootDirectionDiagonalUpForward
	ShootDirectionDiagonalDownForward
	ShootDirectionDiagonalUpBack
	ShootDirectionDiagonalDownBack
)

type StateTransitionHandler interface {
	TransitionToShooting(direction ShootDirection)
	TransitionFromShooting()
}
