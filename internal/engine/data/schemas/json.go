package schemas

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
)

// ShapeRect defines a rectangular shape with position and dimensions.
type ShapeRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Rect returns the coordinates and dimensions of the rectangle.
func (s ShapeRect) Rect() (x, y, width, height int) {
	return s.X, s.Y, s.Width, s.Height
}

// AssetData holds information about a single asset, including its path and collision areas.
type AssetData struct {
	Path           string      `json:"path"`
	CollisionRects []ShapeRect `json:"collision_rect"`
	Loop           *bool       `json:"loop,omitempty"`
}

// MovementConfig defines horizontal movement parameters.
type MovementConfig struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	HorizontalSpeed float64 `json:"horizontal_speed,omitempty"`
}

// JumpConfig defines jump ability parameters.
type JumpConfig struct {
	Enabled           *bool   `json:"enabled,omitempty"`
	JumpCutMultiplier float64 `json:"jump_cut_multiplier,omitempty"`
	CoyoteTimeFrames  int     `json:"coyote_time_frames,omitempty"`
	JumpBufferFrames  int     `json:"jump_buffer_frames,omitempty"`
}

// DashConfig defines dash ability parameters.
type DashConfig struct {
	Enabled    *bool `json:"enabled,omitempty"`
	DurationMs int   `json:"duration_ms,omitempty"`
	CooldownMs int   `json:"cooldown_ms,omitempty"`
	Speed      int   `json:"speed,omitempty"`
	CanAirDash *bool `json:"can_air_dash,omitempty"`
}

// ShootingConfig defines shooting ability parameters.
type ShootingConfig struct {
	Enabled         *bool `json:"enabled,omitempty"`
	CooldownFrames  int   `json:"cooldown_frames,omitempty"`
	ProjectileSpeed int   `json:"projectile_speed,omitempty"`
	ProjectileRange int   `json:"projectile_range,omitempty"`
	Directions      int   `json:"directions,omitempty"`
}

// SkillsConfig defines all skill configurations for an entity.
type SkillsConfig struct {
	Movement *MovementConfig `json:"movement,omitempty"`
	Jump     *JumpConfig     `json:"jump,omitempty"`
	Dash     *DashConfig     `json:"dash,omitempty"`
	Shooting *ShootingConfig `json:"shooting,omitempty"`
}

// EnemyWeaponConfig defines the weapon configuration for an enemy entity.
type EnemyWeaponConfig struct {
	ProjectileType string `json:"projectile_type"`
	Speed          int    `json:"speed"`
	Cooldown       int    `json:"cooldown"`
	Damage         int    `json:"damage"`
	Range          int    `json:"range"`
	ShootMode      string `json:"shoot_mode,omitempty"`
	ShootDirection string `json:"shoot_direction,omitempty"`
	ShootState     string `json:"shoot_state,omitempty"`
}

// SpriteData contains all data related to a sprite's appearance and behavior,
// including its body rectangle, assets for different states, animation frame rate, and initial facing direction.
type SpriteData struct {
	BodyRect        ShapeRect                     `json:"body_rect"`
	Assets          map[string]AssetData          `json:"assets"`
	FrameRate       int                           `json:"frame_rate"`
	FacingDirection animation.FacingDirectionEnum `json:"facing_direction"` // 0 - right, 1 - left
	Skills          *SkillsConfig                 `json:"skills,omitempty"`
	Weapon          *EnemyWeaponConfig            `json:"weapon,omitempty"`
}

// ParticleData defines the configuration for a particle effect.
// A particle entry is either image-based (Image set) or pixel-based (Pixel set).
type ParticleData struct {
	Image       string             `json:"image,omitempty"`
	FrameWidth  int                `json:"frame_width,omitempty"`
	FrameHeight int                `json:"frame_height,omitempty"`
	FrameRate   int                `json:"frame_rate,omitempty"` // Ticks per frame
	Scale       float64            `json:"scale,omitempty"`
	Pixel       *PixelParticleData `json:"pixel,omitempty"`
}

// PixelParticleData defines a pixel-based (no image asset) particle.
// Color is a hex string limited to "#000000" or "#FFFFFF" (1-bit palette).
type PixelParticleData struct {
	Size           int    `json:"size"`
	Color          string `json:"color"`
	LifetimeFrames int    `json:"lifetime_frames"`
}
