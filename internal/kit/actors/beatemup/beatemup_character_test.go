package beatemup_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
)

// newTestFixtures returns the minimal arguments accepted by the new
// NewBeatEmUpCharacter constructor (see SPEC §3). Assets map is empty so
// no PNG decoding is needed in the unit-test environment.
func newTestFixtures() (fstest.MapFS, map[string]animation.SpriteState, schemas.SpriteData, *bodyphysics.Rect) {
	fsys := fstest.MapFS{}
	stateMap := map[string]animation.SpriteState{}
	spriteData := schemas.SpriteData{
		Assets:          map[string]schemas.AssetData{},
		FrameRate:       1,
		FacingDirection: 0, // right
	}
	bodyRect := bodyphysics.NewRect(0, 0, 8, 8)
	return fsys, stateMap, spriteData, bodyRect
}

// T-B1: NewBeatEmUpCharacter returns a non-nil character whose embedded
// Character and MeleeCharacter are both initialised.
func TestNewBeatEmUpCharacter_EmbedsCharacterAndMeleeCharacter(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewBeatEmUpCharacter returned nil character")
	}
	if c.Character == nil {
		t.Error("expected embedded *actors.Character to be initialised")
	}
	if c.MeleeCharacter == nil {
		t.Error("expected embedded *MeleeCharacter to be initialised")
	}
}

// T-B2: BeatEmUpCharacter owns a BeatEmUpMovementModel — not a PlatformMovementModel.
func TestNewBeatEmUpCharacter_OwnsBeatEmUpMovementModel(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	model := c.MovementModel()
	if model == nil {
		t.Fatal("expected MovementModel() to be non-nil after construction")
	}
	if _, ok := model.(*physicsmovement.BeatEmUpMovementModel); !ok {
		t.Fatalf("expected *BeatEmUpMovementModel, got %T", model)
	}
	if _, ok := model.(*physicsmovement.PlatformMovementModel); ok {
		t.Fatal("BeatEmUpCharacter must not own a PlatformMovementModel")
	}
}

// T-B3: Update on a zero-input frame does not panic and leaves the body at rest.
// Verifies the model is properly wired into Character.Update (no panic on
// type assertion) and that there is no gravity accumulation when idle.
func TestBeatEmUpCharacter_Update_ZeroInput_NoPanic_NoDrift(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	// Empty input — no movement commands.
	origReader := input.CommandsReader
	defer func() { input.CommandsReader = origReader }()
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{} }

	sp := space.NewSpace()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Update panicked on zero-input frame: %v", r)
		}
	}()

	if err := c.Update(sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	vx, vy := c.Velocity()
	if vx != 0 || vy != 0 {
		t.Errorf("expected zero velocity on idle frame; got (%d, %d)", vx, vy)
	}
}

// T-B4: Registering an EightDirectionalMovementSkill causes HandleInput to be
// invoked each frame: with Right pressed, the body accelerates rightward.
func TestBeatEmUpCharacter_EightDirSkill_RightInputProducesRightwardVelocity(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	c.SetSpeed(200)
	c.AddSkill(kitskills.NewEightDirectionalMovementSkill())

	origReader := input.CommandsReader
	defer func() { input.CommandsReader = origReader }()
	input.CommandsReader = func() input.PlayerCommands { return input.PlayerCommands{Right: true} }

	sp := space.NewSpace()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Update panicked with EightDirectionalMovementSkill registered: %v", r)
		}
	}()

	if err := c.Update(sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	vx, _ := c.Velocity()
	if vx <= 0 {
		t.Errorf("expected positive X velocity after Right input; got vx=%d", vx)
	}
}

// --- Footprint test helpers (story 064) -----------------------------------

// tinyPNGBytes returns the bytes of a 1x1 transparent PNG. It exists so that
// LoadSprites can decode a real image while keeping unit tests independent of
// any on-disk asset.
func tinyPNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0, 0, 0, 0})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

// newFootprintFixtures builds fixtures wired so that key "idle" → actors.Idle
// in stateMap and Assets["idle"] points to a real (tiny) PNG file in the
// in-memory fsys. The asset's FootprintRect is left nil here — callers set it
// before invoking the constructor.
func newFootprintFixtures(t *testing.T) (fstest.MapFS, map[string]animation.SpriteState, schemas.SpriteData, *bodyphysics.Rect) {
	t.Helper()
	png := tinyPNGBytes(t)
	fsys := fstest.MapFS{
		"idle.png": &fstest.MapFile{Data: png},
	}
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
	}
	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {Path: "idle.png"},
		},
		FrameRate:       1,
		FacingDirection: animation.FaceDirectionRight,
	}
	bodyRect := bodyphysics.NewRect(0, 0, 8, 8)
	return fsys, stateMap, spriteData, bodyRect
}

// T-F1: Footprint() returns the world-offset rect when the current state has
// a footprint_rect declared in JSON.
func TestBeatEmUpCharacter_Footprint_ReturnsWorldOffsetRect(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.FootprintRect = &schemas.ShapeRect{X: 2, Y: 30, Width: 12, Height: 6}
	spriteData.Assets["idle"] = asset

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	c.SetID("beatemup-actor")
	c.SetPosition(100, 200)

	got := c.Footprint()
	want := image.Rect(102, 230, 114, 236)
	if got != want {
		t.Errorf("Footprint() = %v, want %v (world-offset of {X:2,Y:30,W:12,H:6} at body min (100,200))", got, want)
	}
}

// T-F2: Footprint() falls back to the union of CollidableBody.CollisionPosition()
// when no footprint_rect is declared for the current state.
func TestBeatEmUpCharacter_Footprint_FallsBackToCollisionUnion(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	// Intentionally do NOT set FootprintRect — fallback path under test.

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	c.SetID("beatemup-actor")

	// Attach a single collision body equal to the actor body Position().
	collision := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(0, 0, 8, 8))
	collision.SetID("seed")
	c.AddCollision(collision)

	got := c.Footprint()
	want := c.Position()
	if got != want {
		t.Errorf("Footprint() = %v, want %v (Position() because the single collision rect matches the body)", got, want)
	}
}

// T-F3: Footprint() falls back to Position() when no footprint AND no
// collision rects are configured for the current state.
func TestBeatEmUpCharacter_Footprint_FallsBackToPositionWhenNoCollisions(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}

	got := c.Footprint()
	want := image.Rect(0, 0, 8, 8)
	if got != want {
		t.Errorf("Footprint() = %v, want %v (no footprint, no collision rects → body Position())", got, want)
	}
}

// T-F4: CollisionPosition() returns ONLY the footprint world-rect when a
// footprint exists. Movement vs. world collision checks must consume the
// feet area, not the full body.
func TestBeatEmUpCharacter_CollisionPosition_FootprintOnly(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.FootprintRect = &schemas.ShapeRect{X: 2, Y: 30, Width: 12, Height: 6}
	spriteData.Assets["idle"] = asset

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	c.SetID("beatemup-actor")
	c.SetPosition(100, 200)

	rs := c.CollisionPosition()
	if len(rs) != 1 {
		t.Fatalf("CollisionPosition() len = %d, want 1 (only footprint)", len(rs))
	}
	want := image.Rect(102, 230, 114, 236)
	if rs[0] != want {
		t.Errorf("CollisionPosition()[0] = %v, want %v", rs[0], want)
	}
}

// T-F5: CollisionPosition() falls back to the embedded *CollidableBody
// behavior when no footprint exists for the current state.
func TestBeatEmUpCharacter_CollisionPosition_FallsBackToEmbedded(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	// No FootprintRect — fallback path under test.

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	c.SetID("beatemup-actor")

	collision := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(0, 0, 8, 8))
	collision.SetID("seed")
	c.AddCollision(collision)

	got := c.CollisionPosition()
	want := c.CollidableBody.CollisionPosition()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CollisionPosition() = %v, want %v (must equal embedded behavior when no footprint)", got, want)
	}
}

// T-F6: A zero-size footprint_rect is treated as absent — Footprint() must
// take the fallback path. This locks the spec's edge-case decision.
func TestBeatEmUpCharacter_Footprint_ZeroSizeTreatedAsAbsent(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.FootprintRect = &schemas.ShapeRect{X: 0, Y: 0, Width: 0, Height: 0}
	spriteData.Assets["idle"] = asset

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}

	got := c.Footprint()
	want := c.Position()
	if got != want {
		t.Errorf("Footprint() = %v, want %v (zero-size footprint must be ignored)", got, want)
	}
}

// T-F7: A state change updates Footprint()'s target. Footprint is read fresh
// from the current state — when the actor leaves "idle" (the only state with
// a footprint here), Footprint() falls back.
func TestBeatEmUpCharacter_Footprint_StateChangeUpdatesTarget(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.FootprintRect = &schemas.ShapeRect{X: 2, Y: 30, Width: 12, Height: 6}
	spriteData.Assets["idle"] = asset

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	c.SetID("beatemup-actor")
	c.SetPosition(100, 200)

	r1 := c.Footprint()
	wantIdle := image.Rect(102, 230, 114, 236)
	if r1 != wantIdle {
		t.Errorf("Footprint() before state change = %v, want %v", r1, wantIdle)
	}

	c.SetNewStateFatal(actors.Walking)

	r2 := c.Footprint()
	// In Walking there is no footprint registered → fallback to Position().
	// There are also no collision rects registered → fallback to body Position().
	wantWalking := c.Position()
	if r2 != wantWalking {
		t.Errorf("Footprint() after switching to Walking = %v, want %v (fallback)", r2, wantWalking)
	}
	if r1 == r2 {
		t.Error("expected Footprint() to differ across states; got identical rects")
	}
}

// T-F8: Facing-left does NOT mirror the footprint (parity with collision_rect,
// which is also stored as-authored and only mirrored at draw time).
func TestBeatEmUpCharacter_Footprint_FaceLeftDoesNotMirror(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newFootprintFixtures(t)
	asset := spriteData.Assets["idle"]
	asset.FootprintRect = &schemas.ShapeRect{X: 2, Y: 30, Width: 12, Height: 6}
	spriteData.Assets["idle"] = asset

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	c.SetID("beatemup-actor")
	c.SetPosition(100, 200)
	c.SetFaceDirection(animation.FaceDirectionLeft)

	got := c.Footprint()
	want := image.Rect(102, 230, 114, 236)
	if got != want {
		t.Errorf("Footprint() facing left = %v, want %v (mirroring must NOT be applied)", got, want)
	}
}

// newObstacleBody returns a real *CollidableBody positioned to occupy `rect`.
// Using the production CollidableBody type (instead of a hand-rolled stub) keeps
// the test wired through the same Collidable surface area the engine consumes.
func newObstacleBody(t *testing.T, id string, rect image.Rectangle) *bodyphysics.CollidableBody {
	t.Helper()
	b := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(0, 0, rect.Dx(), rect.Dy()))
	b.SetID(id)
	b.SetPosition(rect.Min.X, rect.Min.Y)
	b.SetIsObstructive(true)
	return b
}

// T-I1: space.HasCollision consults the footprint for a beatemup actor.
//
//   - Case A: body bbox overlaps the obstacle but the feet do not → no collision.
//   - Case B: the feet overlap the obstacle → collision.
//
// This pins the observable behavior: the engine pipeline must see ONLY the
// footprint when one is declared, so movement-vs-world checks no longer use
// the full body bbox for beat-em-up actors.
func TestBeatEmUpCharacter_HasCollision_UsesFootprint(t *testing.T) {
	cases := []struct {
		name          string
		footprint     schemas.ShapeRect // local (unoffset)
		bodyPosition  image.Point
		bodySize      image.Point
		obstacle      image.Rectangle
		wantCollision bool
	}{
		{
			name:          "case A: body overlaps but feet are clear",
			footprint:     schemas.ShapeRect{X: 4, Y: 80, Width: 24, Height: 8}, // world: (4,180)-(28,188)
			bodyPosition:  image.Pt(0, 100),
			bodySize:      image.Pt(200, 100), // body: (0,100)-(200,200)
			obstacle:      image.Rect(0, 100, 200, 140),
			wantCollision: false,
		},
		{
			name:          "case B: feet overlap obstacle",
			footprint:     schemas.ShapeRect{X: 4, Y: 10, Width: 24, Height: 25}, // world: (4,110)-(28,135)
			bodyPosition:  image.Pt(0, 100),
			bodySize:      image.Pt(200, 100),
			obstacle:      image.Rect(0, 100, 200, 140),
			wantCollision: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fsys, stateMap, spriteData, _ := newFootprintFixtures(t)
			asset := spriteData.Assets["idle"]
			fp := tc.footprint
			asset.FootprintRect = &fp
			spriteData.Assets["idle"] = asset

			bodyRect := bodyphysics.NewRect(0, 0, tc.bodySize.X, tc.bodySize.Y)
			c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
			if err != nil {
				t.Fatalf("NewBeatEmUpCharacter: %v", err)
			}
			c.SetID("beatemup-actor")
			c.SetPosition(tc.bodyPosition.X, tc.bodyPosition.Y)

			obstacle := newObstacleBody(t, "obstacle", tc.obstacle)

			got := space.HasCollision(c, obstacle)
			if got != tc.wantCollision {
				t.Errorf("HasCollision = %v, want %v\n  body=%v footprint(world)=%v obstacle=%v",
					got, tc.wantCollision, c.Position(), c.CollisionPosition(), tc.obstacle)
			}
		})
	}
}
