package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/boilerplate/ebiten-template/internal/game/entity/actors/methods"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

// climberStateTransitionLogic provides custom state handling for the ClimberPlayer.
func climberStateTransitionLogic(c *actors.Character) bool {
	state := c.State()

	if state == gamestates.Rising && c.IsAnimationFinished() {
		c.SetNewStateFatal(actors.Idle)
		return true
	}

	if state == gamestates.Growing && c.IsAnimationFinished() {
		c.SetNewStateFatal(actors.Idle)
		if p, ok := c.Owner().(interface{ SetScale(float64) }); ok {
			p.SetScale(2.0)
		}
		return true
	}

	if state == gamestates.Shrinking && c.IsAnimationFinished() {
		c.SetNewStateFatal(actors.Idle)
		if p, ok := c.Owner().(interface {
			SetScale(float64)
			SetSize(int, int)
			GetPosition16() (int, int)
			SetPosition16(int, int)
			RefreshCollisions()
		}); ok {
			p.SetScale(1.0)

			// Capture current bottom-center position
			x16, y16 := p.GetPosition16()
			centerX16 := x16 + (32*16)/2 // Transition frame is 32x32
			bottomY16 := y16 + (32 * 16)

			// Restore natural size
			p.SetSize(16, 16)

			// Re-position to maintain bottom-center
			newX16 := centerX16 - (16*16)/2
			newY16 := bottomY16 - (16 * 16)
			p.SetPosition16(newX16, newY16)

			p.RefreshCollisions()
		}
		return true
	}

	if state == gamestates.Exiting || state == gamestates.Lying || state == gamestates.Rising || state == gamestates.Growing || state == gamestates.Shrinking {
		return true
	}

	return false
}

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed   int
	freezeSkill *gameskill.FreezeSkill
	growSkill   *gameskill.GrowSkill
	starSkill   *gameskill.StarSkill

	*gameplayermethods.PlayerDeathBehavior
}

// NewClimberPlayer creates a new climber player.
func NewClimberPlayer(ctx *app.AppContext) (platformer.PlatformerActorEntity, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/player/climber.json")
	if err != nil {
		return nil, err
	}

	character.SetStateTransitionHandler(climberStateTransitionLogic)

	player := &ClimberPlayer{
		PlatformerCharacter: character,
		freezeSkill:         gameskill.NewFreezeSkill(),
		growSkill:           gameskill.NewGrowSkill(),
		starSkill:           gameskill.NewStarSkill(),
	}
	// Set the owner on the embedded character so LastOwner() works correctly
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	// Configure Star Skill VFX
	player.starSkill.OnActive = func() {
		if ctx.VFX != nil && ctx.FrameCount%4 == 0 {
			rect := player.Position()
			centerX := float64(rect.Min.X + rect.Dx()/2)
			centerY := float64(rect.Min.Y + rect.Dy()/2)
			gamevfx.SpawnStarParticles(ctx.VFX, centerX, centerY, 3)
		}
	}

	// Configure Freeze Skill VFX
	player.freezeSkill.OnActive = func() {
		if ctx.VFX != nil && ctx.FrameCount%5 == 0 {
			rect := player.Position()
			centerX := float64(rect.Min.X + rect.Dx()/2)
			centerY := float64(rect.Min.Y + rect.Dy()/2)
			gamevfx.SpawnFreezeAuraParticles(ctx.VFX, centerX, centerY, 3)
		}
	}

	// Configure Grow Skill VFX
	player.growSkill.OnActive = func() {
		if ctx.VFX != nil && ctx.FrameCount%5 == 0 {
			rect := player.Position()
			centerX := float64(rect.Min.X + rect.Dx()/2)
			centerY := float64(rect.Min.Y + rect.Dy()/2)
			gamevfx.SpawnGrowAuraParticles(ctx.VFX, centerX, centerY, 3)
		}
	}

	character.AddSkill(player.freezeSkill)
	character.AddSkill(player.growSkill)
	character.AddSkill(player.starSkill)

	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "player"); err != nil {
		return nil, err
	}
	player.baseSpeed = player.Speed()

	if err = builder.ApplyPlatformerPhysics(player, player); err != nil {
		return nil, err
	}

	character.StateCollisionManager.RefreshCollisions()
	player.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(player)

	return player, nil
}

func (p *ClimberPlayer) ActivateFreezeSkill() {
	if p.freezeSkill != nil {
		p.freezeSkill.RequestActivation()
	}
}

func (p *ClimberPlayer) ActivateFreezeSkillWithItem(item interface{}) {
	if p.freezeSkill != nil {
		ctx := p.AppContext()
		if ctx != nil && ctx.Space != nil {
			p.freezeSkill.RequestActivationWithItem(item.(body.Collidable), ctx.Space)
		} else {
			p.freezeSkill.RequestActivation()
		}
	}
}

func (p *ClimberPlayer) ActivateGrowSkill() {
	if p.growSkill != nil {
		p.growSkill.RequestActivation()
	}
}

func (p *ClimberPlayer) ActivateGrowSkillWithItem(item interface{}) {
	if p.growSkill != nil {
		// Get the physics space from the app context
		ctx := p.AppContext()
		if ctx != nil && ctx.Space != nil {
			if collidable, ok := item.(interface{}); ok {
				p.growSkill.RequestActivationWithItem(collidable.(body.Collidable), ctx.Space)
			} else {
				p.growSkill.RequestActivation()
			}
		} else {
			p.growSkill.RequestActivation()
		}
	}
}

func (p *ClimberPlayer) ActivateStarSkill() {
	if p.starSkill != nil {
		p.starSkill.RequestActivation()
	}
}

func (p *ClimberPlayer) ActivateStarSkillWithItem(item interface{}) {
	if p.starSkill != nil {
		ctx := p.AppContext()
		if ctx != nil && ctx.Space != nil {
			p.starSkill.RequestActivationWithItem(item.(body.Collidable), ctx.Space)
		} else {
			p.starSkill.RequestActivation()
		}
	}
}

func (p *ClimberPlayer) ResetSkills() {
	if p.freezeSkill != nil {
		p.freezeSkill.Reset()
	}
	if p.growSkill != nil {
		p.growSkill.Reset(p)
	}
	if p.starSkill != nil {
		p.starSkill.Reset()
	}
}

func (p *ClimberPlayer) IsGrowActive() bool {
	return p.growSkill != nil && p.growSkill.IsActive()
}

func (p *ClimberPlayer) IsStarActive() bool {
	return p.starSkill != nil && p.starSkill.IsActive()
}

func (p *ClimberPlayer) Update(space body.BodiesSpace) error {
	p.SetHorizontalInertia(-1.0)
	p.SetSpeed(p.baseSpeed)
	p.SetJumpForceMultiplier(1.0)
	return p.Character.Update(space)
}

func (p *ClimberPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *ClimberPlayer) Hurt(damage int) {
	if p.IsStarActive() || p.IsGrowActive() {
		return
	}
	if p.State() == gamestates.Dying {
		return
	}

	p.SetNewStateFatal(gamestates.Dying)
}

func (p *ClimberPlayer) OnTouch(other body.Collidable) {
	// Handle Star Skill and Grow Skill collision with enemies
	if p.IsStarActive() || p.IsGrowActive() {
		ctx := p.AppContext()
		owner := other.LastOwner()
		if enemy, ok := owner.(gameentitytypes.EnemyActor); ok && enemy.IsEnemy() {
			// Kill enemy
			pos := enemy.Position()
			centerX := float64(pos.Min.X + pos.Dx()/2)
			centerY := float64(pos.Min.Y + pos.Dy()/2)

			if ctx.VFX != nil {
				ctx.AudioManager.PlaySoundAtVolume("assets/audio/Fire_Hit_01.ogg", 0.5)
				ctx.VFX.SpawnDeathExplosion(centerX, centerY, 15)
			}

			if ctx.ActorManager != nil {
				ctx.ActorManager.Unregister(enemy)
			}

			// Remove from physics space
			if space := ctx.Space; space != nil {
				space.QueueForRemoval(enemy)
			}
		}
	}
}

func (p *ClimberPlayer) OnBlock(other body.Collidable) {
	// Required to implement body.Touchable to avoid recursion if we rely on embedded CollidableBody.OnBlock
}
