package gameentitytypes

type FreezeSkillUser interface {
	ActivateFreezeSkill()
}

type GrowSkillUser interface {
	ActivateGrowSkill()
	IsGrowActive() bool
}

type StarSkillUser interface {
	ActivateStarSkill()
	IsStarActive() bool
}

type InvincibleSkillUser interface {
	IsStarActive() bool
	IsGrowActive() bool
}
