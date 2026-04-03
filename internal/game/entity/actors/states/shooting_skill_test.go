package gamestates_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

func TestShootingSkill_CooldownGating(t *testing.T) {
	tests := []struct {
		name           string
		cooldownFrames int
		updateCalls    int
		wantSpawns     int
	}{
		{"no double-spawn within cooldown", 3, 3, 1},
		{"cooldown resets after window expires", 2, 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spawnCount := 0
			shooter := &mocks.MockShooter{
				SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {
					spawnCount++
				},
			}

			cfg := gamestates.ShootingConfig{
				CooldownFrames: tt.cooldownFrames,
				SpawnOffsetX16: 16 << 4,
				BulletSpeedX16: 32 << 4,
				YOffset:        4,
			}
			skill := gamestates.NewShootingSkill(cfg, shooter)

			body := &MockBody{
				GetPosition16Func: func() (int, int) { return 100 << 4, 50 << 4 },
				FaceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
				OwnerFunc:         func() interface{} { return nil },
			}

			for i := 0; i < tt.updateCalls; i++ {
				skill.Update(body, nil)
			}

			if spawnCount != tt.wantSpawns {
				t.Errorf("got %d spawns, want %d", spawnCount, tt.wantSpawns)
			}
		})
	}
}

func TestShootingSkill_AlternatingYOffset(t *testing.T) {
	var yOffsets []int
	shooter := &mocks.MockShooter{
		SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {
			yOffsets = append(yOffsets, y16)
		},
	}

	cfg := gamestates.ShootingConfig{
		CooldownFrames: 0,
		SpawnOffsetX16: 0,
		BulletSpeedX16: 32 << 4,
		YOffset:        4,
	}
	skill := gamestates.NewShootingSkill(cfg, shooter)

	body := &MockBody{
		GetPosition16Func: func() (int, int) { return 0, 50 << 4 },
		FaceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
		OwnerFunc:         func() interface{} { return nil },
	}

	for i := 0; i < 4; i++ {
		skill.Update(body, nil)
	}

	want := []int{(50 << 4) + 4, (50 << 4) - 4, (50 << 4) + 4, (50 << 4) - 4}
	if len(yOffsets) != len(want) {
		t.Fatalf("got %d offsets, want %d", len(yOffsets), len(want))
	}
	for i, got := range yOffsets {
		if got != want[i] {
			t.Errorf("offset[%d]: got %d, want %d", i, got, want[i])
		}
	}
}

func TestShootingSkill_ReleaseRepressWithinCooldown(t *testing.T) {
	spawnCount := 0
	shooter := &mocks.MockShooter{
		SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {
			spawnCount++
		},
	}

	cfg := gamestates.ShootingConfig{
		CooldownFrames: 3,
		SpawnOffsetX16: 0,
		BulletSpeedX16: 32 << 4,
		YOffset:        4,
	}
	skill := gamestates.NewShootingSkill(cfg, shooter)

	body := &MockBody{
		GetPosition16Func: func() (int, int) { return 0, 0 },
		FaceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
		OwnerFunc:         func() interface{} { return nil },
	}

	skill.Update(body, nil)
	skill.Update(body, nil)

	if spawnCount != 1 {
		t.Errorf("got %d spawns, want 1", spawnCount)
	}
}

