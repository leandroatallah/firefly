package body

// Passthrough marks a body or its owner that projectiles should ignore entirely.
// Entities implementing this interface will not be hit, damaged, or trigger
// impact effects when a projectile overlaps them.
type Passthrough interface {
	IsPassthrough() bool
}
