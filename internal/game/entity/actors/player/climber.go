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

	for _, s := range gamestates.MeleeAttackStepStates(3) {
		if state == s {
			return !c.IsAnimationFinished()
		}
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
	meleeBuffered bool
	meleeAnimWait int

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
	isGrounded := !p.IsFalling() && !p.IsGoingUp()

	if p.melee != nil {
		p.melee.Update()

		if p.meleeAnimWait > 0 {
			p.meleeAnimWait--
		}

		if (cmds.Dash || cmds.Jump) && p.melee.ComboWindowRemaining() > 0 {
			p.melee.ResetCombo()
			p.meleeBuffered = false
			p.meleeAnimWait = 0
		}

		meleePressed := cmds.Melee && !p.meleeHeldPrev
		if meleePressed && isGrounded && (p.melee.IsSwinging() || p.meleeAnimWait > 0) {
			p.meleeBuffered = true
		}

		canAct := p.melee.CanFire() && !p.melee.IsSwinging() && p.meleeAnimWait == 0 && !p.IsDucking()
		wantFire := canAct && (meleePressed || (isGrounded && p.meleeBuffered && p.melee.ComboWindowRemaining() > 0))
		if wantFire {
			if isGrounded && p.melee.ComboWindowRemaining() > 0 {
				p.melee.AdvanceCombo()
			}
			x16, y16 := p.GetPosition16()
			p.melee.Fire(x16, y16, p.FaceDirection(), body.ShootDirectionStraight, 0)
			p.spawnMeleeVFX(x16, y16)
			p.meleeBuffered = false
			p.meleeAnimWait = p.meleeStepAnimDuration(p.melee.StepIndex())
		}

		if !p.melee.IsSwinging() && p.melee.ComboWindowRemaining() == 0 && p.meleeAnimWait == 0 {
			p.meleeBuffered = false
		}

		if p.melee.IsHitboxActive() {
			p.melee.ApplyHitbox(space)
		}
		p.meleeHeldPrev = cmds.Melee
	}

	// Check for ducking input
	duckHeld := cmds.Down
	p.SetDucking(duckHeld && !p.IsFalling() && !p.IsGoingUp())

	isMeleeActive := p.melee != nil && (p.melee.IsSwinging() || p.meleeAnimWait > 0) && !p.IsFalling() && !p.IsGoingUp()
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

func (p *ClimberPlayer) Hurt(_ int) {
	if p.melee != nil {
		p.melee.ResetCombo()
		p.meleeBuffered = false
		p.meleeAnimWait = 0
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
// Character.Faction().
func (p *ClimberPlayer) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
	p.melee = w
	p.meleeVFX = vfxMgr
	if w != nil {
		w.SetOwner(p)
		stepStates := gamestates.MeleeAttackStepStates(len(w.Steps()))
		p.AddStateContributor(&meleeContributor{w: w, stepStates: stepStates})
	}
}

// meleeStepAnimDuration returns the total game-frames for the sprite animation
// of the given combo step, derived from the sprite sheet width and frame rate.
func (p *ClimberPlayer) meleeStepAnimDuration(stepIdx int) int {
	if p.melee == nil {
		return 0
	}
	stepStates := gamestates.MeleeAttackStepStates(len(p.melee.Steps()))
	if stepIdx < 0 || stepIdx >= len(stepStates) {
		return 0
	}
	char := p.GetCharacter()
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
	p.meleeVFX.SpawnDirectionalPuff("melee_slash", px, py, p.FaceDirection() == animation.FaceDirectionRight, 1, 0.0)
}
