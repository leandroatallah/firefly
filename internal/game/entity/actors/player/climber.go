package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	meleeengine "github.com/boilerplate/ebiten-template/internal/engine/combat/melee"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	gameplayermethods "github.com/boilerplate/ebiten-template/internal/game/entity/actors/methods"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed  int
	spriteData *schemas.SpriteData
	inventory  interface{}

	// melee is kept as a separate field (not in inventory) because Controller
	// bundles weapon lifecycle, hitbox activation, input buffering, and combo
	// state — none of which are present on the combat.Weapon interface.
	melee *meleeengine.Controller

	*gameplayermethods.PlayerDeathBehavior
}

// NewClimberPlayer creates a new climber player.
func NewClimberPlayer(ctx *app.AppContext) (platformer.PlatformerActorEntity, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/player/climber.json")
	if err != nil {
		return nil, err
	}

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
	isGrounded := !p.IsFalling() && !p.IsGoingUp()

	if p.melee != nil {
		p.melee.SetSpace(space)
		p.melee.Tick(p.GetCharacter())

		if p.melee.HandleInput(cmds.Melee, cmds.Dash, cmds.Jump, isGrounded, p.IsDucking()) {
			p.melee.EnterAttackState(p.GetCharacter())
		}
	}

	// Check for ducking input
	duckHeld := cmds.Down
	p.SetDucking(duckHeld && !p.IsFalling() && !p.IsGoingUp())

	isMeleeActive := p.melee != nil && p.melee.IsBlockingMovement() && !p.IsFalling() && !p.IsGoingUp()
	isGroundedShooting := isGrounded && (p.State() == actors.IdleShooting || p.State() == actors.WalkingShooting)

	// When ducking, mid-melee, or shooting on ground: lock horizontal movement but allow facing change
	if p.IsDucking() || isMeleeActive || isGroundedShooting {
		p.SetSpeed(0)
		p.SetHorizontalInertia(0)

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

func (p *ClimberPlayer) MeleeController() *meleeengine.Controller { return p.melee }

func (p *ClimberPlayer) Hurt(_ int) {
	if p.melee != nil {
		p.melee.OnInterrupt()
	}

	if p.State() == gamestates.Dying || p.State() == gamestates.Dead {
		return
	}

	p.SetNewStateFatal(gamestates.Dying)
}

func (p *ClimberPlayer) OnTouch(_ body.Collidable) {
	// Standard player touch behavior
}

func (p *ClimberPlayer) OnBlock(_ body.Collidable) {
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
// Character.Faction(). A MeleeController is installed on the character,
// making it the sole owner of the swing lifecycle and state machine hooks.
func (p *ClimberPlayer) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
	if w == nil {
		return
	}
	w.SetOwner(p)

	char := p.GetCharacter()
	stepStates := gamestates.MeleeAttackStepStates(len(w.Steps()))

	st := meleeengine.InstallState(char, p, w, vfxMgr, gamestates.StateMeleeAttack, gamestates.StateGrounded, actors.Falling, stepStates)

	p.melee = meleeengine.New(w, st, gamestates.StateMeleeAttack, stepStates, p.meleeStepAnimDuration)
	p.melee.Install(char)
}

// meleeStepAnimDuration returns the total game-frames for the sprite animation
// of the given combo step, derived from the sprite sheet width and frame rate.
func (p *ClimberPlayer) meleeStepAnimDuration(stepIdx int) int {
	if p.melee == nil {
		return 0
	}
	char := p.GetCharacter()
	stepStates := gamestates.MeleeAttackStepStates(p.melee.StepCount())
	if stepIdx < 0 || stepIdx >= len(stepStates) {
		return 0
	}
	sprite := char.GetSpriteByState(stepStates[stepIdx])
	if sprite == nil || sprite.Image == nil {
		return 0
	}
	rect := char.Position()
	if rect.Dx() == 0 {
		return 0
	}
	frameCount := sprite.Image.Bounds().Dx() / rect.Dx()
	frameRate := char.FrameRate()
	if frameRate == 0 {
		frameRate = 1
	}
	return frameCount * frameRate
}
