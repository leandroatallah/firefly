package skill

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
)

func ptrBool(b bool) *bool {
	return &b
}

type mockShooter struct{}

func (m *mockShooter) SpawnBullet(x16, y16, vx16, vy16 int, owner interface{}) {}

func TestFromConfig_AllSkillsEnabled(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Movement: &schemas.MovementConfig{Enabled: ptrBool(true)},
		Jump:     &schemas.JumpConfig{Enabled: ptrBool(true), JumpCutMultiplier: 0.4},
		Dash:     &schemas.DashConfig{Enabled: ptrBool(true)},
		Shooting: &schemas.ShootingConfig{Enabled: ptrBool(true), CooldownFrames: 15},
	}
	deps := SkillDeps{Shooter: &mockShooter{}}

	skills := FromConfig(cfg, deps)

	if len(skills) != 4 {
		t.Errorf("expected 4 skills, got %d", len(skills))
	}

	// Verify Jump skill has correct cut multiplier
	var jumpSkill *JumpSkill
	for _, s := range skills {
		if js, ok := s.(*JumpSkill); ok {
			jumpSkill = js
			break
		}
	}
	if jumpSkill == nil {
		t.Fatal("jump skill not found")
	}
	if jumpSkill.jumpCutMultiplier != 0.4 {
		t.Errorf("expected jump cut multiplier 0.4, got %f", jumpSkill.jumpCutMultiplier)
	}

	// Verify Shooting skill has correct cooldown
	var shootingSkill *ShootingSkill
	for _, s := range skills {
		if ss, ok := s.(*ShootingSkill); ok {
			shootingSkill = ss
			break
		}
	}
	if shootingSkill == nil {
		t.Fatal("shooting skill not found")
	}
	if shootingSkill.cooldown != 15 {
		t.Errorf("expected shooting cooldown 15, got %d", shootingSkill.cooldown)
	}
}

func TestFromConfig_NilConfig(t *testing.T) {
	deps := SkillDeps{}
	skills := FromConfig(nil, deps)

	if len(skills) != 0 {
		t.Errorf("expected empty slice for nil config, got %d skills", len(skills))
	}
}

func TestFromConfig_DisabledSkills(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Movement: &schemas.MovementConfig{Enabled: ptrBool(false)},
		Jump:     &schemas.JumpConfig{Enabled: ptrBool(false)},
		Dash:     &schemas.DashConfig{Enabled: ptrBool(false)},
		Shooting: &schemas.ShootingConfig{Enabled: ptrBool(false)},
	}
	deps := SkillDeps{Shooter: &mockShooter{}}

	skills := FromConfig(cfg, deps)

	if len(skills) != 0 {
		t.Errorf("expected 0 skills when all disabled, got %d", len(skills))
	}
}

func TestFromConfig_MissingShooter(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Movement: &schemas.MovementConfig{Enabled: ptrBool(true)},
		Jump:     &schemas.JumpConfig{Enabled: ptrBool(true)},
		Dash:     &schemas.DashConfig{Enabled: ptrBool(true)},
		Shooting: &schemas.ShootingConfig{Enabled: ptrBool(true)},
	}
	deps := SkillDeps{Shooter: nil}

	skills := FromConfig(cfg, deps)

	if len(skills) != 3 {
		t.Errorf("expected 3 skills (shooting skipped), got %d", len(skills))
	}

	// Verify shooting skill is not present
	for _, s := range skills {
		if _, ok := s.(*ShootingSkill); ok {
			t.Fatal("shooting skill should not be present when Shooter is nil")
		}
	}
}

func TestFromConfig_PartialConfig(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Jump: &schemas.JumpConfig{Enabled: ptrBool(true)},
		Dash: &schemas.DashConfig{Enabled: ptrBool(true)},
	}
	deps := SkillDeps{}

	skills := FromConfig(cfg, deps)

	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
}

func TestFromConfig_JumpCallbackAttached(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Jump: &schemas.JumpConfig{Enabled: ptrBool(true)},
	}
	deps := SkillDeps{
		OnJump: func(b interface{}) {},
	}

	skills := FromConfig(cfg, deps)

	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}

	jumpSkill, ok := skills[0].(*JumpSkill)
	if !ok {
		t.Fatal("expected JumpSkill")
	}

	if jumpSkill.OnJump == nil {
		t.Fatal("OnJump callback not attached")
	}
}

func TestFromConfig_NilSubConfigs(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Movement: nil,
		Jump:     nil,
		Dash:     nil,
		Shooting: nil,
	}
	deps := SkillDeps{}

	skills := FromConfig(cfg, deps)

	if len(skills) != 0 {
		t.Errorf("expected 0 skills for all nil sub-configs, got %d", len(skills))
	}
}
