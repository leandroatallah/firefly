package builder

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
)

// newConfigurableActor returns an ActorEntity with a Character attached so that
// ConfigureEnemyWeapon can wire an EnemyShooter whose owner is the character.
func newConfigurableActor() *mockActorWithCollision {
	actor := newMockActorWithCollision()
	character := actors.NewCharacter(sprites.SpriteMap{}, bodyphysics.NewRect(0, 0, 16, 16))
	actor.SetCharacter(character)
	actor.SetID("enemy-under-test")
	return actor
}

func baseEnemyWeaponCfg() *schemas.EnemyWeaponConfig {
	return &schemas.EnemyWeaponConfig{
		ProjectileType: "bullet_small",
		Speed:          6,
		Cooldown:       90,
		Damage:         1,
		Range:          160,
		ShootMode:      "on_sight",
		ShootDirection: "horizontal",
	}
}

func TestConfigureEnemyWeapon_NilConfig(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}

	shooter, err := ConfigureEnemyWeapon(actor, nil, mgr)
	if err != nil {
		t.Fatalf("expected nil error for nil config, got %v", err)
	}
	if shooter != nil {
		t.Errorf("expected nil shooter for nil config, got %v", shooter)
	}
}

func TestConfigureEnemyWeapon_MissingManager(t *testing.T) {
	actor := newConfigurableActor()
	cfg := baseEnemyWeaponCfg()

	shooter, err := ConfigureEnemyWeapon(actor, cfg, nil)
	if err == nil {
		t.Fatal("expected error when projectile manager is nil, got nil")
	}
	if shooter != nil {
		t.Errorf("expected nil shooter when manager is nil, got %v", shooter)
	}
}

func TestConfigureEnemyWeapon_Builds_OnSightHorizontal(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := baseEnemyWeaponCfg()

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shooter == nil {
		t.Fatal("expected non-nil EnemyShooter")
	}
	if shooter.Mode() != combat.ShootModeOnSight {
		t.Errorf("Mode(): got %v, want ShootModeOnSight", shooter.Mode())
	}
	if shooter.Direction() != body.ShootDirectionStraight {
		t.Errorf("Direction(): got %v, want ShootDirectionStraight", shooter.Direction())
	}
	if shooter.Range() != cfg.Range {
		t.Errorf("Range(): got %d, want %d", shooter.Range(), cfg.Range)
	}
	if _, ok := shooter.ShootState(); ok {
		t.Error("ShootState(): expected inactive gate (empty ShootState), but active flag is true")
	}
}

func TestConfigureEnemyWeapon_Builds_AlwaysVertical(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := &schemas.EnemyWeaponConfig{
		ProjectileType: "bullet_small",
		Speed:          5,
		Cooldown:       60,
		Damage:         1,
		Range:          0,
		ShootMode:      "always",
		ShootDirection: "vertical",
	}

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shooter == nil {
		t.Fatal("expected non-nil EnemyShooter")
	}
	if shooter.Mode() != combat.ShootModeAlways {
		t.Errorf("Mode(): got %v, want ShootModeAlways", shooter.Mode())
	}
	if shooter.Direction() != body.ShootDirectionDown {
		t.Errorf("Direction(): got %v, want ShootDirectionDown", shooter.Direction())
	}
}

func TestConfigureEnemyWeapon_DefaultsOnEmptyStrings(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := &schemas.EnemyWeaponConfig{
		ProjectileType: "bullet_small",
		Speed:          4,
		Cooldown:       30,
		Damage:         1,
		Range:          80,
		// ShootMode and ShootDirection omitted: defaults are on_sight/horizontal.
	}

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shooter.Mode() != combat.ShootModeOnSight {
		t.Errorf("default Mode(): got %v, want ShootModeOnSight", shooter.Mode())
	}
	if shooter.Direction() != body.ShootDirectionStraight {
		t.Errorf("default Direction(): got %v, want ShootDirectionStraight", shooter.Direction())
	}
}

func TestConfigureEnemyWeapon_InvalidShootMode(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := baseEnemyWeaponCfg()
	cfg.ShootMode = "sometimes"

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err == nil {
		t.Fatal("expected error for invalid shoot_mode, got nil")
	}
	if shooter != nil {
		t.Errorf("expected nil shooter on invalid shoot_mode, got %v", shooter)
	}
}

func TestConfigureEnemyWeapon_InvalidShootDirection(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := baseEnemyWeaponCfg()
	cfg.ShootDirection = "diagonal"

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err == nil {
		t.Fatal("expected error for invalid shoot_direction, got nil")
	}
	if shooter != nil {
		t.Errorf("expected nil shooter on invalid shoot_direction, got %v", shooter)
	}
}

func TestConfigureEnemyWeapon_UnknownShootState(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := baseEnemyWeaponCfg()
	cfg.ShootState = "definitely_not_a_registered_state_name"

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err == nil {
		t.Fatal("expected error for unknown shoot_state, got nil")
	}
	if shooter != nil {
		t.Errorf("expected nil shooter on unknown shoot_state, got %v", shooter)
	}
}

func TestConfigureEnemyWeapon_ValidShootState(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := baseEnemyWeaponCfg()
	cfg.ShootState = "walk" // corresponds to actors.Walking

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shooter == nil {
		t.Fatal("expected non-nil EnemyShooter")
	}
	enum, ok := shooter.ShootState()
	if !ok {
		t.Fatal("ShootState(): active flag should be true for non-empty shoot_state")
	}
	wantEnum, wantOk := actors.GetStateEnum("walk")
	if !wantOk {
		t.Fatal("precondition failed: actors.GetStateEnum(\"walk\") returned ok=false")
	}
	if enum != wantEnum {
		t.Errorf("ShootState() enum: got %v, want %v", enum, wantEnum)
	}
}

// Compile-time guard: ConfigureEnemyWeapon returns an EnemyShooter whose methods
// include Target/SetTarget so the game layer can prime the on_sight range gate.
func TestConfigureEnemyWeapon_ShooterHasTargetAPI(t *testing.T) {
	actor := newConfigurableActor()
	mgr := &mockProjectileManager{}
	cfg := baseEnemyWeaponCfg()

	shooter, err := ConfigureEnemyWeapon(actor, cfg, mgr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shooter.Target() != nil {
		t.Errorf("Target() before SetTarget: got %v, want nil", shooter.Target())
	}
	// Call SetTarget with a nil body.MovableCollidable — should not panic.
	var nilTarget body.MovableCollidable
	shooter.SetTarget(nilTarget)
}
