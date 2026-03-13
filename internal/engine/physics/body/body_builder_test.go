package body

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
)

// mockEntity implements body.Collidable for testing
type mockEntity struct {
	*CollidableBody
	collisionRects []body.Collidable
}

func newMockEntity() *mockEntity {
	cb := NewCollidableBodyFromRect(NewRect(0, 0, 10, 10))
	return &mockEntity{
		CollidableBody: cb,
		collisionRects: make([]body.Collidable, 0),
	}
}

func (m *mockEntity) addCollisionRect(state animation.SpriteState, rect body.Collidable) {
	m.collisionRects = append(m.collisionRects, rect)
}

func TestSetCollisionBodies(t *testing.T) {
	entity := newMockEntity()

	// Create sprite data with collision rects
	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path: "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{
					{X: 0, Y: 0, Width: 5, Height: 10},
					{X: 5, Y: 0, Width: 5, Height: 10},
				},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
	}

	idCounter := 0
	idProvider := func(assetKey string, index int) string {
		idCounter++
		return assetKey + "_" + string(rune('0'+index))
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	if len(entity.collisionRects) != 2 {
		t.Errorf("expected 2 collision rects; got %d", len(entity.collisionRects))
	}
}

func TestSetCollisionBodies_WithTouchable(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"walk": {
				Path:           "assets/walk.png",
				CollisionRects: []schemas.ShapeRect{},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"walk": "walk",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_0"
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	// Entity should have touchable set
	if entity.GetTouchable() == nil {
		t.Error("expected touchable to be set")
	}
}

func TestSetCollisionBodies_MissingState(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"attack": {
				Path: "assets/attack.png",
				CollisionRects: []schemas.ShapeRect{
					{X: 0, Y: 0, Width: 10, Height: 10},
				},
			},
		},
	}

	// State map doesn't include "attack"
	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_0"
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	// No collision rects should be added since state is missing
	if len(entity.collisionRects) != 0 {
		t.Errorf("expected 0 collision rects (missing state); got %d", len(entity.collisionRects))
	}
}

func TestSetCollisionBodies_EmptyCollisionRects(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path:           "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_0"
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	if len(entity.collisionRects) != 0 {
		t.Errorf("expected 0 collision rects (empty); got %d", len(entity.collisionRects))
	}
}

func TestSetCollisionBodies_MultipleAssets(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path: "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{
					{X: 0, Y: 0, Width: 5, Height: 10},
				},
			},
			"walk": {
				Path: "assets/walk.png",
				CollisionRects: []schemas.ShapeRect{
					{X: 0, Y: 0, Width: 6, Height: 10},
					{X: 6, Y: 0, Width: 4, Height: 10},
				},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
		"walk": "walk",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_" + string(rune('0'+index))
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	// Should have 1 (idle) + 2 (walk) = 3 collision rects
	if len(entity.collisionRects) != 3 {
		t.Errorf("expected 3 collision rects; got %d", len(entity.collisionRects))
	}
}

func TestSetCollisionBodies_CollisionRectPosition(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path: "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{
					{X: 10, Y: 20, Width: 5, Height: 10},
				},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_0"
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	if len(entity.collisionRects) != 1 {
		t.Fatalf("expected 1 collision rect; got %d", len(entity.collisionRects))
	}

	pos := entity.collisionRects[0].Position()
	if pos.Min.X != 10 || pos.Min.Y != 20 {
		t.Errorf("expected collision rect at (10,20); got (%d,%d)", pos.Min.X, pos.Min.Y)
	}

	if pos.Dx() != 5 || pos.Dy() != 10 {
		t.Errorf("expected collision rect 5x10; got %dx%d", pos.Dx(), pos.Dy())
	}
}

func TestSetCollisionBodies_IDProvider(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path: "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{
					{X: 0, Y: 0, Width: 5, Height: 10},
					{X: 5, Y: 0, Width: 5, Height: 10},
				},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_idx" + string(rune('0'+index))
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	SetCollisionBodies(entity, spriteData, stateMap, idProvider, addCollisionRect)

	if len(entity.collisionRects) != 2 {
		t.Fatalf("expected 2 collision rects; got %d", len(entity.collisionRects))
	}

	if entity.collisionRects[0].ID() != "idle_idx0" {
		t.Errorf("expected ID 'idle_idx0'; got '%s'", entity.collisionRects[0].ID())
	}
	if entity.collisionRects[1].ID() != "idle_idx1" {
		t.Errorf("expected ID 'idle_idx1'; got '%s'", entity.collisionRects[1].ID())
	}
}

func TestSetCollisionBodies_NilEntity(t *testing.T) {
	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path:           "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{},
			},
		},
	}

	stateMap := map[string]animation.SpriteState{
		"idle": "idle",
	}

	idProvider := func(assetKey string, index int) string {
		return assetKey + "_0"
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {}

	// Nil entity will cause a panic when calling SetTouchable
	// This is expected behavior - the function requires a valid entity
	defer func() {
		if r := recover(); r == nil {
			t.Error("SetCollisionBodies should panic with nil entity")
		}
	}()

	SetCollisionBodies(nil, spriteData, stateMap, idProvider, addCollisionRect)
}

func TestSetCollisionBodies_NilParameters(t *testing.T) {
	entity := newMockEntity()

	spriteData := schemas.SpriteData{
		Assets: map[string]schemas.AssetData{
			"idle": {
				Path:           "assets/idle.png",
				CollisionRects: []schemas.ShapeRect{},
			},
		},
	}

	// Test with nil stateMap
	idProvider := func(assetKey string, index int) string {
		return assetKey + "_0"
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		entity.addCollisionRect(state, rect)
	}

	// Should not panic with nil stateMap
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetCollisionBodies panicked: %v", r)
		}
	}()

	SetCollisionBodies(entity, spriteData, nil, idProvider, addCollisionRect)
}
