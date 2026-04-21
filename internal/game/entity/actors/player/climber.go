package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	gameplayermethods "github.com/boilerplate/ebiten-template/internal/game/entity/actors/methods"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

// climberStateTransitionLogic provides custom state handling for the ClimberPlayer.
func climberStateTransitionLogic(c *actors.Character) bool {
	state := c.State()

	if state == gamestates.Exiting || state == gamestates.Dying || state == gamestates.Dead {
		return true
	}

	return false
}

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed     int
	spriteData    *schemas.SpriteData
	inventory     interface{}
	melee         *weapon.MeleeWeapon
	meleeVFX      vfx.Manager
	meleeHeldPrev bool

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
	}
	// Set the owner on the embedded character so LastOwner() works correctly
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "player"); err != nil {
		return nil, err
	}
	player.baseSpeed = player.Speed()

	if err = builder.ApplyPlatformerPhysics(player, player); err != nil {
		return nil, err
	}

	character.RefreshCollisions()
	player.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(player)

	// Store spriteData for later skill configuration
	player.spriteData = &spriteData

	return player, nil
}

func (p *ClimberPlayer) Update(space body.BodiesSpace) error {
	cmds := input.CommandsReader()

	if p.melee != nil {
		p.melee.Update()
		meleePressed := cmds.Melee && !p.meleeHeldPrev
		if meleePressed && p.melee.CanFire() && !p.IsDucking() {
			x16, y16 := p.GetPosition16()
			p.melee.Fire(x16, y16, p.FaceDirection(), body.ShootDirectionStraight, 0)
			p.spawnMeleeVFX(x16, y16)
		}
		if p.melee.IsHitboxActive() {
			p.melee.ApplyHitbox(space)
		}
		p.meleeHeldPrev = cmds.Melee
	}

	// Check for ducking input
	duckHeld := cmds.Down
	p.SetDucking(duckHeld && !p.IsFalling() && !p.IsGoingUp())

	// When ducking, prevent horizontal movement but allow facing direction change
	if p.IsDucking() {
		p.SetSpeed(0)
		p.SetHorizontalInertia(0)

		// Allow facing direction change while ducking
		if cmds.Left {
			p.SetFaceDirection(animation.FaceDirectionLeft)
		} else if cmds.Right {
			p.SetFaceDirection(animation.FaceDirectionRight)
		}
	} else {
		p.SetHorizontalInertia(-1.0)
		p.SetSpeed(p.baseSpeed)
	}

	p.SetJumpForceMultiplier(1.0)

	return p.Character.Update(space)
}

func (p *ClimberPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *ClimberPlayer) GetSpriteData() *schemas.SpriteData {
	return p.spriteData
}

func (p *ClimberPlayer) Hurt(damage int) {
	if p.State() == gamestates.Dying || p.State() == gamestates.Dead {
		return
	}

	p.SetNewStateFatal(gamestates.Dying)
}

func (p *ClimberPlayer) OnTouch(other body.Collidable) {
	// Standard player touch behavior
}

func (p *ClimberPlayer) OnBlock(other body.Collidable) {
	// Required to implement body.Touchable to avoid recursion if we rely on embedded CollidableBody.OnBlock
}

func (p *ClimberPlayer) SetInventory(inv interface{}) {
	p.inventory = inv
}

func (p *ClimberPlayer) Inventory() interface{} {
	return p.inventory
}

// SetMelee wires the player's melee weapon and VFX manager. The owner is set
// on the weapon so the faction gate in ApplyHitbox resolves against
// Character.Faction().
func (p *ClimberPlayer) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
	p.melee = w
	p.meleeVFX = vfxMgr
	if w != nil {
		w.SetOwner(p)
	}
}

func (p *ClimberPlayer) spawnMeleeVFX(x16, y16 int) {
	if p.meleeVFX == nil {
		return
	}
	offsetX16 := fp16.To16(12)
	if p.FaceDirection() == animation.FaceDirectionLeft {
		offsetX16 = -offsetX16
	}
	px := float64(fp16.From16(x16 + offsetX16))
	py := float64(fp16.From16(y16))
	p.meleeVFX.SpawnPuff("melee_slash", px, py, 6, 4)
}
