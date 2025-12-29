// Package movement provides the patrol movement behavior for actors.
package movement

import bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"

type patrolStateEnum int

const (
	patrolIdle patrolStateEnum = iota
	patrolChase
)

// WaypointConfig holds configuration for waypoint generation
type WaypointConfig struct {
	Predefined []*bodyphysics.Rect // Custom waypoints for Predefined type
	IdleDelay  int                 // Delay in frames between waypoint transitions
}

// NewPredefinedWaypointConfig creates a predefined waypoint configuration
func NewPredefinedWaypointConfig(waypoints []*bodyphysics.Rect, idleDelay int) *WaypointConfig {
	return &WaypointConfig{
		Predefined: waypoints,
		IdleDelay:  idleDelay,
	}
}

type PatrolMovementState struct {
	BaseMovementState
	count              int
	currentTargetIndex int
	waypoints          []*bodyphysics.Rect
	patrolState        patrolStateEnum
	movementDirections MovementDirections
	idleDelay          int
	waypointConfig     *WaypointConfig
}

func NewPatrolMovementState(base BaseMovementState) *PatrolMovementState {
	return &PatrolMovementState{BaseMovementState: base}
}

// generateWaypoints creates waypoints based on the configuration
func (s *PatrolMovementState) generateWaypoints() {
	if s.waypointConfig == nil {
		// Fallback to no waypoints if none are provided
		s.waypoints = []*bodyphysics.Rect{}
		return
	}
	s.waypoints = append(s.waypoints, s.waypointConfig.Predefined...)
}

// SetWaypointConfig allows setting a custom waypoint configuration
func (s *PatrolMovementState) SetWaypointConfig(config *WaypointConfig) {
	s.waypointConfig = config
	// Regenerate waypoints with new configuration
	s.waypoints = []*bodyphysics.Rect{}
	s.generateWaypoints()
}

func (s *PatrolMovementState) OnStart() {
	// Initialize waypoint configuration if not set
	if s.waypointConfig == nil {
		s.waypointConfig = &WaypointConfig{
			IdleDelay: 60,
		}
	}

	// Initialize idle delay from configuration
	s.idleDelay = s.waypointConfig.IdleDelay

	// Generate waypoints based on configuration
	s.generateWaypoints()

	if len(s.waypoints) > 0 {
		target := s.CurrentWaypoint()
		o := bodyphysics.NewObstacleRect(target)
		s.movementDirections = calculateMovementDirections(s.actor, o, false)
		s.patrolState = patrolChase
	}

	s.BaseMovementState.OnStart()
}

func (s *PatrolMovementState) Move() {
	s.count++

	if s.actor.Immobile() {
		return
	}

	if len(s.waypoints) == 0 {
		return
	}

	switch s.patrolState {
	case patrolIdle:
		if s.count > s.idleDelay {
			s.currentTargetIndex = (s.currentTargetIndex + 1) % len(s.waypoints)
			target := s.CurrentWaypoint()
			o := bodyphysics.NewObstacleRect(target)
			s.movementDirections = calculateMovementDirections(s.actor, o, false)
			s.patrolState = patrolChase
			s.count = 0
		}
	case patrolChase:
		executeMovement(s.actor, s.movementDirections)
		if s.count > s.idleDelay {
			s.patrolState = patrolIdle
			s.count = 0
		}
	}
}

func (s *PatrolMovementState) CurrentWaypoint() *bodyphysics.Rect {
	current := s.waypoints[s.currentTargetIndex]
	return current
}

// GetWaypointConfig returns the current waypoint configuration
func (s *PatrolMovementState) GetWaypointConfig() *WaypointConfig {
	return s.waypointConfig
}

// GetWaypointCount returns the number of waypoints in the current patrol
func (s *PatrolMovementState) GetWaypointCount() int {
	return len(s.waypoints)
}

// GetCurrentWaypointIndex returns the index of the current waypoint
func (s *PatrolMovementState) GetCurrentWaypointIndex() int {
	return s.currentTargetIndex
}

// Functional Options Pattern
// sets the waypoint configuration for patrol movement
func WithWaypointConfig(config *WaypointConfig) MovementStateOption {
	return func(ms MovementState) {
		if patrolState, ok := ms.(*PatrolMovementState); ok {
			patrolState.SetWaypointConfig(config)
		}
	}
}
