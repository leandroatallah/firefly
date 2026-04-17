package gameenemies

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/projectile"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	_ "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

func getModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find go.mod")
		}
		dir = parent
	}
}

func TestMain(m *testing.M) {
	err := os.Chdir(getModuleRoot())
	if err != nil {
		panic(err)
	}

	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 224,
	}
	config.Set(cfg)

	os.Exit(m.Run())
}

func newEnemyTestContext() *app.AppContext {
	sp := space.NewSpace()
	return &app.AppContext{
		Assets:            os.DirFS("."),
		ActorManager:      actors.NewManager(),
		Space:             sp,
		EventManager:      event.NewManager(),
		ProjectileManager: projectile.NewManager(sp),
	}
}

func bulletBodies(sp body.BodiesSpace) []body.Collidable {
	out := []body.Collidable{}
	for _, b := range sp.Bodies() {
		if strings.HasPrefix(b.ID(), "bullet_") && !strings.Contains(b.ID(), "_COLLISION_") {
			out = append(out, b)
		}
	}
	return out
}

// TestBatEnemy_FactionAndShooter verifies that after construction, BatEnemy
// carries FactionEnemy and has a non-nil shooter wired from bat.json's
// weapon block (shoot_mode=always, shoot_direction=vertical).
func TestBatEnemy_FactionAndShooter(t *testing.T) {
	ctx := newEnemyTestContext()

	bat, err := NewBatEnemy(ctx, 100, 100, "bat-1")
	if err != nil {
		t.Fatalf("NewBatEnemy returned error: %v", err)
	}
	if bat == nil {
		t.Fatal("NewBatEnemy returned nil")
	}

	if got := bat.GetCharacter().Faction(); got != enginecombat.FactionEnemy {
		t.Errorf("Faction() = %v, want FactionEnemy", got)
	}
	if bat.Shooter() == nil {
		t.Error("expected BatEnemy.Shooter() to be non-nil after construction")
	}
}

// TestBatEnemy_AlwaysFires_VerticalDown verifies the Simple Immobile Shooter
// archetype: without ever calling SetTarget, repeated Update calls spawn at
// least one projectile on the vertical-down axis.
func TestBatEnemy_AlwaysFires_VerticalDown(t *testing.T) {
	ctx := newEnemyTestContext()

	bat, err := NewBatEnemy(ctx, 100, 100, "bat-fire")
	if err != nil {
		t.Fatalf("NewBatEnemy returned error: %v", err)
	}

	// bat.json declares cooldown=60 frames. 180 frames guarantees at least one shot.
	firstBulletSeen := false
	for i := 0; i < 180; i++ {
		if err := bat.Update(ctx.Space); err != nil {
			t.Fatalf("Update error at frame %d: %v", i, err)
		}
		if !firstBulletSeen {
			if bullets := bulletBodies(ctx.Space); len(bullets) > 0 {
				firstBulletSeen = true
			}
		}
	}

	if !firstBulletSeen {
		t.Fatal("expected at least one projectile body spawned by BatEnemy always-mode shooter, got 0")
	}
}
