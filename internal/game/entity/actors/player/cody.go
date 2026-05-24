package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/jsonutil"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	kitactors "github.com/boilerplate/ebiten-template/internal/kit/actors"
	beatemupkit "github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	meleeengine "github.com/boilerplate/ebiten-template/internal/kit/combat/melee"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
	kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"
)

type CodyPlayer struct {
	*beatemupkit.BeatEmUpCharacter
	baseSpeed  int
	spriteData *schemas.SpriteData
	inventory  interface{}

	*kitactors.PlayerDeathBehavior
}

func NewCodyPlayer(ctx *app.AppContext) (beatemupkit.BeatEmUpActorEntity, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData](ctx.Assets, "assets/entities/player/cody.json")
	if err != nil {
		return nil, err
	}

	stateMap, err := builder.BuildStateMap(spriteData)
	if err != nil {
		return nil, err
	}

	rect := builder.BodyRectFromSpriteData(spriteData)
	character, err := beatemupkit.NewBeatEmUpCharacter(ctx.Assets, stateMap, spriteData, rect, nil)
	if err != nil {
		return nil, err
	}

	player := &CodyPlayer{
		BeatEmUpCharacter: character,
	}
	player.SetOwner(player)
	character.SetOwner(player)

	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "player"); err != nil {
		return nil, err
	}
	player.baseSpeed = player.Speed()

	character.RefreshCollisions()
	player.PlayerDeathBehavior = kitactors.NewPlayerDeathBehavior(player)

	player.spriteData = &spriteData

	return player, nil
}

func (p *CodyPlayer) Update(space body.BodiesSpace) error {
	cmds := input.CommandsReader()

	melee := p.MeleeController()
	if melee != nil {
		melee.SetSpace(space)
		melee.Tick(p.GetCharacter())

		// TODO: HandleInput has too many unused parameters for different game-genre actors
		if melee.HandleInput(cmds.Melee, cmds.Dash, cmds.Jump, true, false) {
			melee.EnterAttackState(p.GetCharacter())
		}
	}

	isMeleeActive := melee != nil && melee.IsBlockingMovement()
	isShooting := p.State() == actors.IdleShooting || p.State() == actors.WalkingShooting

	if isMeleeActive || isShooting {
		p.SetSpeed(0)
	} else {
		p.SetSpeed(p.baseSpeed)
	}

	return p.Character.Update(space)
}

func (p *CodyPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *CodyPlayer) GetSpriteData() *schemas.SpriteData {
	return p.spriteData
}

func (p *CodyPlayer) Hurt(_ int) {
	melee := p.MeleeController()
	if melee != nil {
		melee.OnInterrupt()
	}

	if p.State() == actors.Dying || p.State() == actors.Dead {
		return
	}

	p.SetNewStateFatal(actors.Dying)
}

func (p *CodyPlayer) OnTouch(_ body.Collidable) {}

func (p *CodyPlayer) OnBlock(_ body.Collidable) {}

func (p *CodyPlayer) SetInventory(inv interface{}) {
	p.inventory = inv
}

func (p *CodyPlayer) Inventory() interface{} {
	return p.inventory
}

func (p *CodyPlayer) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
	if w == nil {
		return
	}
	w.SetOwner(p)

	char := p.GetCharacter()
	stepStates := kitstates.MeleeAttackStepStates(len(w.Steps()))

	st := meleeengine.InstallState(char, p, w, nil, kitstates.StateMeleeAttack, kitstates.StateGrounded, actors.Falling, stepStates)

	controller := meleeengine.New(w, st, kitstates.StateMeleeAttack, stepStates, p.meleeStepAnimDuration)
	p.SetMeleeController(controller)
	controller.Install(char)
}

func (p *CodyPlayer) meleeStepAnimDuration(stepIdx int) int {
	melee := p.MeleeController()
	if melee == nil {
		return 0
	}
	char := p.GetCharacter()
	stepStates := kitstates.MeleeAttackStepStates(melee.StepCount())
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
