package animation

// NOTE: SpriteState is a any type to represent an int enum.

// SpriteState represents the current animation state of a sprite.
type SpriteState interface{}

// FacingDirectionEnum indicates the horizontal direction an entity is facing.
type FacingDirectionEnum int

const (
	// FaceDirectionRight indicates the entity is facing right.
	FaceDirectionRight FacingDirectionEnum = iota
	// FaceDirectionLeft indicates the entity is facing left.
	FaceDirectionLeft
)
