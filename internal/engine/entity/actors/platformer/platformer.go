package platformer

import (
	"image"
	"io/fs"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/context"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/events"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

type AlivePlayer interface {
	Hurt(damage int)
}

type PlatformerActorEntity interface {
	actors.ActorEntity
	context.ContextProvider

	OnDie()
	OnJump()
	OnLand()
	OnFall()
	SetOnJump(func(image.Point))
	SetOnFall(func(image.Point))
	SetOnLand(func(image.Point))
}

type PlatformerCharacter struct {
	*actors.Character
	app.AppContextHolder

	coinCount        int
	movementBlockers int

	jumpHandler func(image.Point)
	landHandler func(image.Point)
	fallHandler func(image.Point)
}

func (p *PlatformerCharacter) SetOnJump(f func(image.Point)) {
	p.jumpHandler = f
}

func (p *PlatformerCharacter) SetOnFall(f func(image.Point)) {
	p.fallHandler = f
}

func (p *PlatformerCharacter) SetOnLand(f func(image.Point)) {
	p.landHandler = f
}

func (p *PlatformerCharacter) OnJump() {
	if p.jumpHandler != nil {
		rect := p.Position()
		// Bottom center
		pos := image.Point{X: rect.Min.X + rect.Dx()/2, Y: rect.Max.Y}
		p.jumpHandler(pos)
	}
}

func (p *PlatformerCharacter) OnFall() {
	if p.fallHandler != nil {
		rect := p.Position()
		// Bottom center
		pos := image.Point{X: rect.Min.X + rect.Dx()/2, Y: rect.Max.Y}
		p.fallHandler(pos)
	}
}

func (p *PlatformerCharacter) OnLand() {
	if p.landHandler != nil {
		rect := p.Position()
		// Bottom center
		pos := image.Point{X: rect.Min.X + rect.Dx()/2, Y: rect.Max.Y}
		p.landHandler(pos)
	}
}

func (p *PlatformerCharacter) SetGravityEnabled(enabled bool) {
	if model, ok := p.Character.MovementModel().(*physicsmovement.PlatformMovementModel); ok {
		model.SetGravityEnabled(enabled)
	}
}

func (p *PlatformerCharacter) AddSkill(s skill.Skill) {
	p.Character.AddSkill(s)
}

func NewPlatformerCharacter(fsys fs.FS, stateMap map[string]animation.SpriteState, spriteData schemas.SpriteData, bodyRect *bodyphysics.Rect) *PlatformerCharacter {
	s, err := sprites.GetSpritesFromAssets(fsys, spriteData.Assets, stateMap)
	if err != nil {
		return nil
	}
	c := actors.NewCharacter(s, bodyRect)
	pf := &PlatformerCharacter{
		Character: c,
	}

	// Note: Jump event is now published directly from JumpSkill.OnJump callback
	// to ensure it fires for all jump types including coyote jumps
	// Land event is still handled here via state machine
	pf.SetOnLand(func(pos image.Point) {
		if pf.AppContext() != nil {
			pf.AppContext().EventManager.Publish(&events.ActorLandedEvent{
				X: float64(pos.X),
				Y: float64(pos.Y),
			})
		}
	})

	c.SetFaceDirection(spriteData.FacingDirection)
	c.SetFrameRate(spriteData.FrameRate)
	c.SetOwner(pf)

	return pf
}
