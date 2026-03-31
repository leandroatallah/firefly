package actors

import (
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
)

const (
	duckHeightRatio = 0.5
)

var Ducking ActorStateEnum

func init() {
	Ducking = RegisterState("duck", func(b BaseState) ActorState { return &DuckingState{BaseState: b} })
}

type DuckingState struct {
	BaseState
	fullHeight int
}

func (s *DuckingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	actor := s.GetActor()
	pos := actor.Position()
	s.fullHeight = pos.Dy()
	duckHeight := s.fullHeight / 2

	newRect := bodyphysics.ResizeFixedBottom(pos, duckHeight)
	actor.SetPosition(newRect.Min.X, newRect.Min.Y)
	actor.SetSize(newRect.Dx(), newRect.Dy())

	_, vy := actor.Velocity()
	actor.SetVelocity(0, vy)
}

func (s *DuckingState) OnFinish() {
	actor := s.GetActor()
	pos := actor.Position()
	newY := pos.Max.Y - s.fullHeight
	actor.SetPosition(pos.Min.X, newY)
	actor.SetSize(pos.Dx(), s.fullHeight)
}
