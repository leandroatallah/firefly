package kitactors

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/hajimehoshi/ebiten/v2"
)

// mockEnemyShooter implements combat.EnemyShooter for testing ShooterCharacter.
type mockEnemyShooter struct {
	updateCalled int
}

func (m *mockEnemyShooter) Update()                         { m.updateCalled++ }
func (m *mockEnemyShooter) SetTarget(_ combat.TargetBody)   {}
func (m *mockEnemyShooter) Target() combat.TargetBody       { return nil }
func (m *mockEnemyShooter) Range() int                      { return 0 }
func (m *mockEnemyShooter) Mode() combat.ShootMode          { return 0 }
func (m *mockEnemyShooter) Direction() body.ShootDirection  { return 0 }
func (m *mockEnemyShooter) ShootState() (interface{}, bool) { return nil, false }
func (m *mockEnemyShooter) TryFire() bool                   { return false }

// mockPlatformerActor implements platformer.PlatformerActorEntity for testing PlayerDeathBehavior.
type mockPlatformerActor struct {
	health int
}

// body.Alive
func (m *mockPlatformerActor) Health() int               { return m.health }
func (m *mockPlatformerActor) MaxHealth() int            { return 0 }
func (m *mockPlatformerActor) SetHealth(h int)           { m.health = h }
func (m *mockPlatformerActor) SetMaxHealth(_ int)        {}
func (m *mockPlatformerActor) LoseHealth(_ int)          {}
func (m *mockPlatformerActor) RestoreHealth(_ int)       {}
func (m *mockPlatformerActor) Invulnerable() bool        { return false }
func (m *mockPlatformerActor) SetInvulnerability(_ bool) {}

// body.Body / body.Ownable
func (m *mockPlatformerActor) ID() string                 { return "" }
func (m *mockPlatformerActor) SetID(_ string)             {}
func (m *mockPlatformerActor) Position() image.Rectangle  { return image.Rectangle{} }
func (m *mockPlatformerActor) SetPosition(_, _ int)       {}
func (m *mockPlatformerActor) SetPosition16(_, _ int)     {}
func (m *mockPlatformerActor) SetSize(_, _ int)           {}
func (m *mockPlatformerActor) Scale() float64             { return 1 }
func (m *mockPlatformerActor) SetScale(_ float64)         {}
func (m *mockPlatformerActor) GetPosition16() (int, int)  { return 0, 0 }
func (m *mockPlatformerActor) GetPositionMin() (int, int) { return 0, 0 }
func (m *mockPlatformerActor) GetShape() body.Shape       { return m }
func (m *mockPlatformerActor) Width() int                 { return 0 }
func (m *mockPlatformerActor) Height() int                { return 0 }
func (m *mockPlatformerActor) Owner() interface{}         { return nil }
func (m *mockPlatformerActor) SetOwner(_ interface{})     {}
func (m *mockPlatformerActor) LastOwner() interface{}     { return nil }

// body.Drawable
func (m *mockPlatformerActor) Image() *ebiten.Image                   { return nil }
func (m *mockPlatformerActor) ImageOptions() *ebiten.DrawImageOptions { return nil }
func (m *mockPlatformerActor) UpdateImageOptions()                    {}

// body.Touchable
func (m *mockPlatformerActor) OnTouch(_ body.Collidable) {}
func (m *mockPlatformerActor) OnBlock(_ body.Collidable) {}

// body.Collidable
func (m *mockPlatformerActor) GetTouchable() body.Touchable                        { return m }
func (m *mockPlatformerActor) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (m *mockPlatformerActor) CollisionPosition() []image.Rectangle                { return nil }
func (m *mockPlatformerActor) CollisionShapes() []body.Collidable                  { return nil }
func (m *mockPlatformerActor) IsObstructive() bool                                 { return false }
func (m *mockPlatformerActor) SetIsObstructive(_ bool)                             {}
func (m *mockPlatformerActor) AddCollision(_ ...body.Collidable)                   {}
func (m *mockPlatformerActor) ClearCollisions()                                    {}
func (m *mockPlatformerActor) SetTouchable(_ body.Touchable)                       {}
func (m *mockPlatformerActor) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}

// body.Movable
func (m *mockPlatformerActor) MoveX(_ int)                                      {}
func (m *mockPlatformerActor) MoveY(_ int)                                      {}
func (m *mockPlatformerActor) OnMoveLeft(_ int)                                 {}
func (m *mockPlatformerActor) OnMoveRight(_ int)                                {}
func (m *mockPlatformerActor) OnMoveUpLeft(_ int)                               {}
func (m *mockPlatformerActor) OnMoveDownLeft(_ int)                             {}
func (m *mockPlatformerActor) OnMoveUpRight(_ int)                              {}
func (m *mockPlatformerActor) OnMoveDownRight(_ int)                            {}
func (m *mockPlatformerActor) OnMoveUp(_ int)                                   {}
func (m *mockPlatformerActor) OnMoveDown(_ int)                                 {}
func (m *mockPlatformerActor) Velocity() (int, int)                             { return 0, 0 }
func (m *mockPlatformerActor) SetVelocity(_, _ int)                             {}
func (m *mockPlatformerActor) Acceleration() (int, int)                         { return 0, 0 }
func (m *mockPlatformerActor) SetAcceleration(_, _ int)                         {}
func (m *mockPlatformerActor) SetSpeed(_ int) error                             { return nil }
func (m *mockPlatformerActor) SetMaxSpeed(_ int) error                          { return nil }
func (m *mockPlatformerActor) Speed() int                                       { return 0 }
func (m *mockPlatformerActor) MaxSpeed() int                                    { return 0 }
func (m *mockPlatformerActor) Immobile() bool                                   { return false }
func (m *mockPlatformerActor) SetImmobile(_ bool)                               {}
func (m *mockPlatformerActor) SetFreeze(_ bool)                                 {}
func (m *mockPlatformerActor) Freeze() bool                                     { return false }
func (m *mockPlatformerActor) FaceDirection() animation.FacingDirectionEnum     { return 0 }
func (m *mockPlatformerActor) SetFaceDirection(_ animation.FacingDirectionEnum) {}
func (m *mockPlatformerActor) IsIdle() bool                                     { return false }
func (m *mockPlatformerActor) IsWalking() bool                                  { return false }
func (m *mockPlatformerActor) IsFalling() bool                                  { return false }
func (m *mockPlatformerActor) IsGoingUp() bool                                  { return false }
func (m *mockPlatformerActor) CheckMovementDirectionX()                         {}
func (m *mockPlatformerActor) TryJump(_ int)                                    {}
func (m *mockPlatformerActor) SetJumpForceMultiplier(_ float64)                 {}
func (m *mockPlatformerActor) JumpForceMultiplier() float64                     { return 1 }
func (m *mockPlatformerActor) SetHorizontalInertia(_ float64)                   {}
func (m *mockPlatformerActor) HorizontalInertia() float64                       { return 1 }

// actors.Controllable
func (m *mockPlatformerActor) BlockMovement()          {}
func (m *mockPlatformerActor) UnblockMovement()        {}
func (m *mockPlatformerActor) IsMovementBlocked() bool { return false }

// actors.Stateful
func (m *mockPlatformerActor) State() actors.ActorStateEnum { return 0 }
func (m *mockPlatformerActor) SetState(_ actors.ActorState) {}
func (m *mockPlatformerActor) SetMovementState(_ movement.MovementStateEnum, _ body.MovableCollidable, _ ...movement.MovementStateOption) {
}
func (m *mockPlatformerActor) SwitchMovementState(_ movement.MovementStateEnum) {}
func (m *mockPlatformerActor) MovementState() movement.MovementState            { return nil }
func (m *mockPlatformerActor) NewState(_ actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}

// actors.Damageable
func (m *mockPlatformerActor) Hurt(_ int) {}

// actors.ActorEntity extras
func (m *mockPlatformerActor) Update(_ body.BodiesSpace) error                  { return nil }
func (m *mockPlatformerActor) MovementModel() physicsmovement.MovementModel     { return nil }
func (m *mockPlatformerActor) SetMovementModel(_ physicsmovement.MovementModel) {}
func (m *mockPlatformerActor) GetCharacter() *actors.Character                  { return nil }

// context.ContextProvider
func (m *mockPlatformerActor) SetAppContext(_ any)         {}
func (m *mockPlatformerActor) AppContext() *app.AppContext { return nil }

// platformer.PlatformerActorEntity event methods
func (m *mockPlatformerActor) OnDie()                        {}
func (m *mockPlatformerActor) OnJump()                       {}
func (m *mockPlatformerActor) OnLand()                       {}
func (m *mockPlatformerActor) OnFall()                       {}
func (m *mockPlatformerActor) SetOnJump(_ func(image.Point)) {}
func (m *mockPlatformerActor) SetOnFall(_ func(image.Point)) {}
func (m *mockPlatformerActor) SetOnLand(_ func(image.Point)) {}
