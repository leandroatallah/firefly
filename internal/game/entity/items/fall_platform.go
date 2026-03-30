package gameitems

import (
	"fmt"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/jsonutil"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
)

var (
	Shaking items.ItemStateEnum
	Break   items.ItemStateEnum
)

func init() {
	Shaking = items.RegisterState("shaking", func(b items.BaseState) items.ItemState {
		return &ShakingState{BaseState: b}
	})
	Break = items.RegisterState("break", func(b items.BaseState) items.ItemState {
		return &BreakState{BaseState: b}
	})
}

type FallingPlatformItem struct {
	items.BaseItem

	count int
}

func NewFallingPlatformItem(ctx *app.AppContext, x, y int, id string) (*FallingPlatformItem, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[items.StatData](ctx.Assets, "assets/entities/items/fall_platform.json")
	if err != nil {
		return nil, err
	}

	// Custom initialization to include Shaking state
	stateMap := map[string]animation.SpriteState{
		"shaking": Shaking,
		"break":   Break,
	}

	base, err := CreateAnimatedItem(ctx.Assets, id, spriteData, stateMap)
	if err != nil {
		return nil, err
	}

	fp := &FallingPlatformItem{
		BaseItem: *base,
	}

	fp.SetPosition(x, y)
	fp.SetAppContext(ctx)
	fp.SetOwner(fp)
	fp.SetIsObstructive(true)

	if err = SetItemBodies(fp, spriteData, stateMap); err != nil {
		return nil, fmt.Errorf("SetItemBodies: %w", err)
	}
	if err = SetItemStats(fp, statData); err != nil {
		return nil, fmt.Errorf("SetItemStats: %w", err)
	}

	fp.StateCollisionManager.RefreshCollisions()

	return fp, nil
}

func (c *FallingPlatformItem) Update(space body.BodiesSpace) error {
	switch c.State() {
	case Shaking:
		c.count++
		if c.count >= timing.FromDuration(time.Second) {
			state, err := items.NewState(c, Break)
			if err != nil {
				return err
			}
			c.SetState(state)
		}
	case Break:
		if c.IsAnimationFinished() {
			c.SetRemoved(true)
		}
	}
	return c.BaseItem.Update(space)
}

func (c *FallingPlatformItem) ResetCount() {
	c.count = 0
}

func (c *FallingPlatformItem) OnTouch(other body.Collidable) {
	if c.IsRemoved() {
		return
	}

	if c.State() != items.Idle {
		return
	}

	player, found := c.AppContext().ActorManager.GetPlayer()
	if !found {
		return
	}

	if other.ID() != player.ID() {
		return
	}

	// Only trigger if the player is above the platform (stepping on it)
	playerRect := other.Position()
	platformRect := c.Position()
	if playerRect.Max.Y > platformRect.Min.Y+4 {
		return
	}

	state, err := items.NewState(c, Shaking)
	if err != nil {
		return
	}
	c.SetState(state)
}

type ShakingState struct {
	items.BaseState
}

func (s *ShakingState) OnStart() {
	s.Item().(*FallingPlatformItem).ResetCount()
}

type BreakState struct {
	items.BaseState
}

func (s *BreakState) OnStart() {}
