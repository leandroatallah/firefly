package space

// DefaultLaneHalfWidth is the recommended default vertical lane tolerance, in
// pixels, used for depth-aware collision gating in the 2.5D beat-em-up plane.
//
// NOTE: This constant is exposed for downstream documentation and default
// configuration only. The implementation of the depth-lane gate inside
// HasCollision is intentionally NOT yet wired in this file — it is added in
// the Feature Implementer step of story 062-depth-aware-collision.
const DefaultLaneHalfWidth = 8

// DepthLaneBody is the opt-in interface bodies implement to participate in
// depth-aware collision gating.
//
// When BOTH bodies in a collision check implement DepthLaneBody, HasCollision
// must additionally require that the vertical distance between their ground
// positions does not exceed the larger of the two LaneHalfWidth values.
//
// GroundY returns the body's logical ground Y coordinate in world space (in
// pixels). For airborne bodies (story 061), this is the shadow position on
// the floor — NOT the visually-offset screen Y.
//
// LaneHalfWidth returns the half-width of this body's depth lane, in pixels.
// A value of 0 means an exact-equal GroundY is required for a match.
type DepthLaneBody interface {
	GroundY() int
	LaneHalfWidth() int
}
