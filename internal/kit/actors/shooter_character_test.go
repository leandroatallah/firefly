package kitactors

import "testing"

func TestShooterCharacter_UpdateShooter_CallsShooterUpdate(t *testing.T) {
	shooter := &mockEnemyShooter{}
	sc := NewShooterCharacter(shooter)
	sc.UpdateShooter()
	if shooter.updateCalled != 1 {
		t.Errorf("expected shooter.Update() called once, got %d", shooter.updateCalled)
	}
}

func TestShooterCharacter_UpdateShooter_NilShooterDoesNotPanic(t *testing.T) {
	sc := NewShooterCharacter(nil)
	sc.UpdateShooter() // must not panic
}

func TestShooterCharacter_Shooter_ReturnsConstructorValue(t *testing.T) {
	shooter := &mockEnemyShooter{}
	sc := NewShooterCharacter(shooter)
	if sc.Shooter() != shooter {
		t.Error("Shooter() did not return the value passed to NewShooterCharacter")
	}
}

func TestShooterCharacter_SetShooter_UpdatesReturnedValue(t *testing.T) {
	sc := NewShooterCharacter(nil)
	shooter := &mockEnemyShooter{}
	sc.SetShooter(shooter)
	if sc.Shooter() != shooter {
		t.Error("Shooter() did not return the value set via SetShooter")
	}
}
