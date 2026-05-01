package kitactors

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"

// ShooterCharacter is a reusable trait that holds an EnemyShooter and provides update logic.
type ShooterCharacter struct {
	shooter combat.EnemyShooter
}

// NewShooterCharacter creates a new ShooterCharacter with the given shooter.
func NewShooterCharacter(shooter combat.EnemyShooter) *ShooterCharacter {
	return &ShooterCharacter{shooter: shooter}
}

// Shooter returns the EnemyShooter (may be nil).
func (s *ShooterCharacter) Shooter() combat.EnemyShooter {
	return s.shooter
}

// SetShooter assigns the shooter field.
func (s *ShooterCharacter) SetShooter(shooter combat.EnemyShooter) {
	s.shooter = shooter
}

// UpdateShooter calls shooter.Update() if shooter is not nil; otherwise no-op.
func (s *ShooterCharacter) UpdateShooter() {
	if s.shooter != nil {
		s.shooter.Update()
	}
}
