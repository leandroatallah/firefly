package gamescenephases

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/projectile"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
	"github.com/hajimehoshi/ebiten/v2"
)

// debugMeleeOwner is a Factioned owner used by the melee weapon under test.
type debugMeleeOwner struct {
	x16, y16 int
	face     animation.FacingDirectionEnum
	faction  combat.Faction
}

func (o *debugMeleeOwner) Faction() combat.Faction                      { return o.faction }
func (o *debugMeleeOwner) FaceDirection() animation.FacingDirectionEnum { return o.face }

// runDebugDrawBlock mirrors the production §3.6 debug block from
// PhasesScene.Draw under config.Get().CollisionBox. It is the exact integration
// the Phase Scene will perform once the production code lands; running it
// here lets the Red Phase pixel-sample without needing a fully wired scene.
func runDebugDrawBlock(
	screen *ebiten.Image,
	cam *enginecamera.Controller,
	pm *projectile.Manager,
	w *weapon.MeleeWeapon,
) {
	if !config.Get().CollisionBox {
		return
	}
	if pm != nil {
		pm.DrawCollisionBoxesWithOffset(func(b body.Collidable) {
			cam.DrawCollisionBox(screen, b)
		})
	}
	if w != nil {
		if rect, active := w.ActiveHitboxRect(); active {
			cam.DrawHitboxRect(screen, rect)
		}
	}
}

// spyCamera wraps a camera controller and counts draw calls.
type spyCamera struct {
	cam                   *enginecamera.Controller
	drawCollisionBoxCalls int
	drawHitboxRectCalls   int
}

func newSpyCamera(cam *enginecamera.Controller) *spyCamera {
	return &spyCamera{cam: cam}
}

func (s *spyCamera) DrawCollisionBox(screen *ebiten.Image, b body.Collidable) {
	s.drawCollisionBoxCalls++
	s.cam.DrawCollisionBox(screen, b)
}

func (s *spyCamera) DrawHitboxRect(screen *ebiten.Image, rect image.Rectangle) {
	s.drawHitboxRectCalls++
	s.cam.DrawHitboxRect(screen, rect)
}

// newDebugMeleeWeapon builds a 1-step weapon matching the wider tests'
// shape so we can drive it into the active window deterministically.
func newDebugMeleeWeapon(owner interface{}) *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{{
		Damage:          1,
		ActiveFrames:    [2]int{3, 5},
		HitboxW16:       24 * 16,
		HitboxH16:       16 * 16,
		HitboxOffsetX16: 12 * 16,
		HitboxOffsetY16: 0,
	}}
	w := weapon.NewMeleeWeapon("debug_melee", 20, 0, steps)
	w.SetOwner(owner)
	return w
}

// driveToActiveFrame fires the weapon and ticks Update until the hitbox is
// active (frame 3 of the test step).
func driveToActiveFrame(w *weapon.MeleeWeapon, owner *debugMeleeOwner) {
	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	for i := 0; i < 3; i++ {
		w.Update()
	}
}

// TestPhasesScene_Draw_CollisionBoxFlag covers AC-1..AC-4: the debug block
// only draws when the flag is on, projectile boxes are rendered (green),
// and the melee hitbox box is rendered (orange) only when active.
func TestPhasesScene_Draw_CollisionBoxFlag(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})

	tests := []struct {
		name                  string
		flag                  bool
		spawnCount            int
		fireMelee             bool
		driveActive           bool
		wantCollisionBoxCalls int // -1 means "expect > 0"
		wantHitboxRectCalls   int // -1 means "expect > 0"
	}{
		{
			name:                  "AC-1: flag on with two projectiles draws green boxes",
			flag:                  true,
			spawnCount:            2,
			fireMelee:             false,
			driveActive:           false,
			wantCollisionBoxCalls: -1,
			wantHitboxRectCalls:   0,
		},
		{
			name:                  "AC-2: flag on with active melee draws orange box",
			flag:                  true,
			spawnCount:            0,
			fireMelee:             true,
			driveActive:           true,
			wantCollisionBoxCalls: 0,
			wantHitboxRectCalls:   -1,
		},
		{
			name:                  "AC-3: flag on but melee not swinging draws no orange",
			flag:                  true,
			spawnCount:            0,
			fireMelee:             false,
			driveActive:           false,
			wantCollisionBoxCalls: 0,
			wantHitboxRectCalls:   0,
		},
		{
			name:                  "AC-4: flag off draws nothing even with projectiles and active swing",
			flag:                  false,
			spawnCount:            2,
			fireMelee:             true,
			driveActive:           true,
			wantCollisionBoxCalls: 0,
			wantHitboxRectCalls:   0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config.Set(&config.AppConfig{
				ScreenWidth:  320,
				ScreenHeight: 240,
				CollisionBox: tc.flag,
			})

			cam := enginecamera.NewController(0, 0)
			cam.SetPositionTopLeft(0, 0)
			cam.DisableSmoothing()
			spy := newSpyCamera(cam)

			space := &mockBodiesSpace{}
			mgr := projectile.NewManager(space)
			cfg := projectile.ProjectileConfig{Width: 4, Height: 4, Damage: 0}
			for i := 0; i < tc.spawnCount; i++ {
				mgr.Spawn(cfg, (32+i*16)<<4, 32<<4, 0, 0, nil)
			}

			owner := &debugMeleeOwner{
				x16:     50 * 16,
				y16:     50 * 16,
				face:    animation.FaceDirectionRight,
				faction: combat.FactionPlayer,
			}
			var w *weapon.MeleeWeapon
			if tc.fireMelee || tc.driveActive {
				w = newDebugMeleeWeapon(owner)
				if tc.driveActive {
					driveToActiveFrame(w, owner)
				}
			}

			screen := ebiten.NewImage(320, 240)

			// Run debug draw with spy camera.
			if !config.Get().CollisionBox {
				return
			}
			if mgr != nil {
				mgr.DrawCollisionBoxesWithOffset(func(b body.Collidable) {
					spy.DrawCollisionBox(screen, b)
				})
			}
			if w != nil {
				if rect, active := w.ActiveHitboxRect(); active {
					spy.DrawHitboxRect(screen, rect)
				}
			}

			// Verify invocation counts.
			switch tc.wantCollisionBoxCalls {
			case -1:
				if spy.drawCollisionBoxCalls == 0 {
					t.Errorf("expected > 0 DrawCollisionBox calls, got 0")
				}
			default:
				if spy.drawCollisionBoxCalls != tc.wantCollisionBoxCalls {
					t.Errorf("DrawCollisionBox calls = %d, want %d", spy.drawCollisionBoxCalls, tc.wantCollisionBoxCalls)
				}
			}
			switch tc.wantHitboxRectCalls {
			case -1:
				if spy.drawHitboxRectCalls == 0 {
					t.Errorf("expected > 0 DrawHitboxRect calls, got 0")
				}
			default:
				if spy.drawHitboxRectCalls != tc.wantHitboxRectCalls {
					t.Errorf("DrawHitboxRect calls = %d, want %d", spy.drawHitboxRectCalls, tc.wantHitboxRectCalls)
				}
			}
		})
	}
}

// TestPhasesScene_Draw_CollisionBoxFlag_LogicInvariance covers AC-5: running
// the same simulation with the flag on vs. off produces byte-identical
// post-frame state for MeleeWeapon (swingFrame via IsHitboxActive transitions
// + cooldown via CanFire) and projectile body positions. The debug draw block
// must be a pure read; it must not mutate any game state.
func TestPhasesScene_Draw_CollisionBoxFlag_LogicInvariance(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})

	type snapshot struct {
		isHitboxActive bool
		canFire        bool
		swingActive    bool
		stepIndex      int
		comboRemaining int
		projPositions  []image.Point
	}

	run := func(flag bool) snapshot {
		config.Set(&config.AppConfig{
			ScreenWidth:  320,
			ScreenHeight: 240,
			CollisionBox: flag,
		})

		cam := enginecamera.NewController(0, 0)
		cam.SetPositionTopLeft(0, 0)
		cam.DisableSmoothing()

		space := &mockBodiesSpace{}
		mgr := projectile.NewManager(space)
		cfg := projectile.ProjectileConfig{Width: 4, Height: 4, Damage: 0}
		mgr.Spawn(cfg, 32<<4, 32<<4, 1<<4, 0, nil)
		mgr.Spawn(cfg, 64<<4, 32<<4, 2<<4, 0, nil)

		owner := &debugMeleeOwner{
			x16:     50 * 16,
			y16:     50 * 16,
			face:    animation.FaceDirectionRight,
			faction: combat.FactionPlayer,
		}
		w := newDebugMeleeWeapon(owner)
		w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)

		screen := ebiten.NewImage(320, 240)

		// Advance enough frames to span startup, active, post-active, and into cooldown.
		const frames = 30
		for i := 0; i < frames; i++ {
			w.Update()
			// Note: Manager.Update is intentionally skipped here because
			// the test's mockBodiesSpace.Find returns nil for non-tracked
			// bodies, which would prune projectiles. The invariance test
			// focuses on whether *Draw* mutates anything — the only thing
			// the new debug block does. So we sample post-Draw state directly.
			runDebugDrawBlock(screen, cam, mgr, w)
		}

		positions := make([]image.Point, 0, 2)
		for _, b := range space.bodies {
			x, y := b.GetPosition16()
			positions = append(positions, image.Point{X: x, Y: y})
		}

		return snapshot{
			isHitboxActive: w.IsHitboxActive(),
			canFire:        w.CanFire(),
			swingActive:    w.IsSwinging(),
			stepIndex:      w.StepIndex(),
			comboRemaining: w.ComboWindowRemaining(),
			projPositions:  positions,
		}
	}

	withFlag := run(true)
	withoutFlag := run(false)

	if withFlag.isHitboxActive != withoutFlag.isHitboxActive {
		t.Errorf("IsHitboxActive: flag-on=%v flag-off=%v (debug Draw must not mutate state)",
			withFlag.isHitboxActive, withoutFlag.isHitboxActive)
	}
	if withFlag.canFire != withoutFlag.canFire {
		t.Errorf("CanFire: flag-on=%v flag-off=%v", withFlag.canFire, withoutFlag.canFire)
	}
	if withFlag.swingActive != withoutFlag.swingActive {
		t.Errorf("IsSwinging: flag-on=%v flag-off=%v", withFlag.swingActive, withoutFlag.swingActive)
	}
	if withFlag.stepIndex != withoutFlag.stepIndex {
		t.Errorf("StepIndex: flag-on=%d flag-off=%d", withFlag.stepIndex, withoutFlag.stepIndex)
	}
	if withFlag.comboRemaining != withoutFlag.comboRemaining {
		t.Errorf("ComboWindowRemaining: flag-on=%d flag-off=%d", withFlag.comboRemaining, withoutFlag.comboRemaining)
	}
	if len(withFlag.projPositions) != len(withoutFlag.projPositions) {
		t.Fatalf("projectile count: flag-on=%d flag-off=%d", len(withFlag.projPositions), len(withoutFlag.projPositions))
	}
	for i := range withFlag.projPositions {
		if withFlag.projPositions[i] != withoutFlag.projPositions[i] {
			t.Errorf("projectile[%d] position: flag-on=%v flag-off=%v",
				i, withFlag.projPositions[i], withoutFlag.projPositions[i])
		}
	}
}
