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
	gameplayermethods "github.com/boilerplate/ebiten-template/internal/game/entity/actors/methods"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed  int
	spriteData *schemas.SpriteData
	inventory  interface{}

	// melee is kept as a separate field (not in inventory) because MeleeWeapon
	// exposes hitbox lifecycle methods (IsHitboxActive, ApplyHitbox) not present
	// on the combat.Weapon interface.
	melee         *weapon.MeleeWeapon
	meleeState    *gamestates.MeleeAttackState // installed per-actor instance
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

	player := &ClimberPlayer{
		PlatformerCharacter: character,
	}
	// Set the owner on the embedded character so LastOwner() works correctly
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	// Install state transition handler as a method closure so it has access to the player.
	character.SetStateTransitionHandler(player.handleStateTransition)

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

// handleStateTransition is the per-player StateTransitionHandler. It drives
// custom states (MeleeAttack step animations, melee state Update) and blocks
// default transitions when appropriate.
func (p *ClimberPlayer) handleStateTransition(c *actors.Character) bool {
	state := c.State()

	if state == gamestates.Exiting || state == gamestates.Dying || state == gamestates.Dead {
		return true
	}

	// Drive melee attack state update: weapon tick, hitbox, frame counter.
	if state == gamestates.StateMeleeAttack {
		if p.meleeState != nil {
			next := p.meleeState.Update()
			if next != gamestates.StateMeleeAttack {
				c.SetNewStateFatal(next)
				// Do not suppress default transitions after exiting MeleeAttackState:
				// return false so stateContributors can re-route to a per-step state if needed.
				return false
			}
			return true
		}
	}

	for _, s := range gamestates.MeleeAttackStepStates(3) {
		if state == s {
			return !c.IsAnimationFinished()
		}
	}

	return false
}

func (p *ClimberPlayer) Update(space body.BodiesSpace) error {
	cmds := input.CommandsReader()
	isGrounded := !p.IsFalling() && !p.IsGoingUp()

	if p.melee != nil {
		// Propagate the current space so MeleeAttackState.Update can call ApplyHitbox.
		if p.meleeState != nil {
			p.meleeState.SetSpace(space)
		}

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

		wantFire := p.melee.CanFire() && !p.melee.IsSwinging() && p.meleeAnimWait == 0 && !p.IsDucking() &&
			(meleePressed || (isGrounded && p.meleeBuffered && p.melee.ComboWindowRemaining() > 0))
		if wantFire {
			p.tryEnterMeleeState(isGrounded)
		}

		if !p.melee.IsSwinging() && p.melee.ComboWindowRemaining() == 0 && p.meleeAnimWait == 0 {
			p.meleeBuffered = false
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

// tryEnterMeleeState advances the combo (if applicable) and transitions the
// character into StateMeleeAttack. The state's OnStart owns Fire + VFX.
func (p *ClimberPlayer) tryEnterMeleeState(isGrounded bool) {
	if isGrounded && p.melee.ComboWindowRemaining() > 0 {
		p.melee.AdvanceCombo()
	}
	animFrames := p.meleeStepAnimDuration(p.melee.StepIndex())
	p.meleeBuffered = false
	p.meleeAnimWait = animFrames
	if p.meleeState != nil {
		p.meleeState.SetAnimationFrames(animFrames)
	}
	if err := p.SetNewState(gamestates.StateMeleeAttack); err != nil {
		return
	}
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
// Character.Faction(). The MeleeAttackState is installed as a per-actor state
// instance on the character, making the state the sole owner of the swing lifecycle.
func (p *ClimberPlayer) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
	p.melee = w
	if w != nil {
		w.SetOwner(p)
		stepStates := gamestates.MeleeAttackStepStates(len(w.Steps()))
		p.AddStateContributor(&meleeContributor{w: w, stepStates: stepStates})

		// Install the real MeleeAttackState for this actor.
		p.meleeState = gamestates.InstallMeleeAttackState(p.GetCharacter(), p, w, vfxMgr)
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
