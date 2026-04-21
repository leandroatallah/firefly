package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	engineskill "github.com/boilerplate/ebiten-template/internal/engine/physics/skill"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

// movementChecker is the subset of body methods the shootingContributor needs.
type movementChecker interface {
	IsWalking() bool
	IsGoingUp() bool
	IsFalling() bool
}

// characterWithSkills is the subset of character methods needed to wire state contributors.
type characterWithSkills interface {
	Skills() []engineskill.Skill
	AddStateContributor(sc actors.StateContributor)
}

// activeChecker reports whether a skill is currently active.
type activeChecker interface {
	IsActive() bool
}

// staticActive wraps a bool for use in tests where no real skill exists.
type staticActive struct{ v bool }

func (s *staticActive) IsActive() bool { return s.v }

type dashContributor struct{ s activeChecker }

func (d *dashContributor) ContributeState(_ actors.ActorStateEnum) (actors.ActorStateEnum, bool) {
	if d.s.IsActive() {
		return gamestates.StateDashing, true
	}
	return 0, false
}

type shootingContributor struct {
	s   activeChecker
	chr movementChecker
}

func (sc *shootingContributor) ContributeState(_ actors.ActorStateEnum) (actors.ActorStateEnum, bool) {
	if !sc.s.IsActive() {
		return 0, false
	}
	switch {
	case sc.chr.IsGoingUp():
		return actors.JumpingShooting, true
	case sc.chr.IsFalling():
		return actors.FallingShooting, true
	case sc.chr.IsWalking():
		return actors.WalkingShooting, true
	default:
		return actors.IdleShooting, true
	}
}

// NewShootingContributorForTest is exported for white-box tests only.
func NewShootingContributorForTest(active bool, chr movementChecker) actors.StateContributor {
	return &shootingContributor{s: &staticActive{v: active}, chr: chr}
}

// NewDashContributorForTest is exported for white-box tests only.
func NewDashContributorForTest(active bool) actors.StateContributor {
	return &dashContributor{s: &staticActive{v: active}}
}

// WireStateContributors registers state contributors for dash and shooting
// so state transitions reflect active skills.
func WireStateContributors(char characterWithSkills, mov movementChecker) {
	for _, s := range char.Skills() {
		switch sk := s.(type) {
		case *engineskill.DashSkill:
			char.AddStateContributor(&dashContributor{s: sk})
		case *engineskill.ShootingSkill:
			char.AddStateContributor(&shootingContributor{s: sk, chr: mov})
		}
	}
}
