package movement

import (
	"container/heap"
	"image"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// ChaseMovementState implements the A* pathfinding algorithm to chase a target.
type ChaseMovementState struct {
	BaseMovementState
	count     int
	path      []image.Point
	obstacles []body.MovableCollidable
}

func NewChaseMovementState(base BaseMovementState) *ChaseMovementState {
	return &ChaseMovementState{BaseMovementState: base}
}

func (s *ChaseMovementState) Move() {
	s.count++

	if s.actor.Immobile() {
		return
	}

	calculatePathRate := 30 // 0.5 seconds in 60 fps

	if s.count == 0 || s.count%calculatePathRate == 0 {
		s.calculatePath()
	}

	if len(s.path) == 0 {
		return // No path found or path is empty
	}

	targetPoint := s.path[0]
	actorPos := s.actor.Position().Min

	// If we are close enough to the next point, remove it from the path
	threshold := 10 // Use a small threshold
	if euclideanDistance(actorPos, targetPoint) < threshold {
		s.path = s.path[1:]
		// Check if it has reached the end of path
		if len(s.path) == 0 {
			return
		}
		targetPoint = s.path[0]
	}

	// Move towards the next point in the path
	dx := targetPoint.X - actorPos.X
	dy := targetPoint.Y - actorPos.Y

	// Add a deadzone to prevent zig-zag movement when the actor is close to the path node.
	// This stops the actor from flipping direction if it slightly overshoots the axis.
	deadzone := 3 // A small pixel tolerance.

	directions := MovementDirections{
		Up:    dy < -deadzone,
		Down:  dy > deadzone,
		Left:  dx < -deadzone,
		Right: dx > deadzone,
	}
	executeMovement(s.actor, directions)
}

func (s *ChaseMovementState) calculatePath() {
	startPos := s.actor.Position().Min
	targetPos := s.target.Position().Min

	startNode := &Node{
		point: startPos,
		gCost: 0,
		hCost: euclideanDistance(startPos, targetPos),
	}
	startNode.fCost = startNode.gCost + startNode.hCost

	openSet := &PriorityQueue{}
	heap.Init(openSet)
	heap.Push(openSet, startNode)
	openSetMap := map[image.Point]*Node{startPos: startNode}

	// visited nodes
	closedSet := make(map[image.Point]*Node)

	actorSize := s.actor.Position().Size()

	for openSet.Len() > 0 {
		// Get the node with the lowest F cost
		currentNode := heap.Pop(openSet).(*Node)
		delete(openSetMap, currentNode.point)

		// Check if we've reached the destination
		if euclideanDistance(currentNode.point, targetPos) < actorSize.X { // Close enough
			s.path = reconstructPath(currentNode)
			return
		}

		closedSet[currentNode.point] = currentNode

		// Explore neighbors
		for _, neighborPoint := range s.getNeighbors(currentNode.point, actorSize) {
			// If neighbor is in the closed set, skip it
			if _, exists := closedSet[neighborPoint]; exists {
				continue
			}

			tentativeGCost := currentNode.gCost + euclideanDistance(currentNode.point, neighborPoint)

			neighborNode, inOpenSet := openSetMap[neighborPoint]

			if inOpenSet && tentativeGCost >= neighborNode.gCost {
				continue // This path is not better
			}

			// This path is the best until now. Record it
			if neighborNode == nil {
				neighborNode = &Node{point: neighborPoint}
				openSetMap[neighborPoint] = neighborNode
			}

			neighborNode.parent = currentNode
			neighborNode.gCost = tentativeGCost
			neighborNode.hCost = euclideanDistance(neighborPoint, targetPos)
			neighborNode.fCost = neighborNode.gCost + neighborNode.hCost

			if !inOpenSet {
				heap.Push(openSet, neighborNode)
			} else {
				heap.Fix(openSet, neighborNode.mapIndex)
			}
		}
	}
}

// isTraversable checks if a given point is a valid and unoccupied position.
func (s *ChaseMovementState) isTraversable(point image.Point, size image.Point) bool {
	// Basic bounds checking
	if point.X < 0 || point.Y < 0 {
		// TODO: Check against map boundaries if they exist
		return false
	}

	// Obstacle detection
	neighborRect := image.Rect(point.X, point.Y, point.X+size.X, point.Y+size.Y)
	for _, obstacle := range s.obstacles {
		if obstacle == s.target {
			continue
		}
		if obstacle != nil && obstacle.Position().Overlaps(neighborRect) {
			return false
		}
	}

	return true
}

// getNeighbors returns valid neighbors for pathfinding.
// It now correctly handles diagonal movements, preventing corner-cutting.
func (s *ChaseMovementState) getNeighbors(point image.Point, size image.Point) []image.Point {
	neighbors := []image.Point{}
	stepX, stepY := size.X, size.Y

	// Check straight directions first
	up := point.Add(image.Point{X: 0, Y: -stepY})
	down := point.Add(image.Point{X: 0, Y: stepY})
	left := point.Add(image.Point{X: -stepX, Y: 0})
	right := point.Add(image.Point{X: stepX, Y: 0})

	isUpTraversable := s.isTraversable(up, size)
	isDownTraversable := s.isTraversable(down, size)
	isLeftTraversable := s.isTraversable(left, size)
	isRightTraversable := s.isTraversable(right, size)

	if isUpTraversable {
		neighbors = append(neighbors, up)
	}
	if isDownTraversable {
		neighbors = append(neighbors, down)
	}
	if isLeftTraversable {
		neighbors = append(neighbors, left)
	}
	if isRightTraversable {
		neighbors = append(neighbors, right)
	}

	// Only add diagonal moves if the two adjacent straight moves are also traversable.
	// This prevents cutting corners of obstacles.
	if isUpTraversable && isLeftTraversable {
		upLeft := point.Add(image.Point{-stepX, -stepY})
		if s.isTraversable(upLeft, size) {
			neighbors = append(neighbors, upLeft)
		}
	}
	if isUpTraversable && isRightTraversable {
		upRight := point.Add(image.Point{stepX, -stepY})
		if s.isTraversable(upRight, size) {
			neighbors = append(neighbors, upRight)
		}
	}
	if isDownTraversable && isLeftTraversable {
		downLeft := point.Add(image.Point{-stepX, stepY})
		if s.isTraversable(downLeft, size) {
			neighbors = append(neighbors, downLeft)
		}
	}
	if isDownTraversable && isRightTraversable {
		downRight := point.Add(image.Point{stepX, stepY})
		if s.isTraversable(downRight, size) {
			neighbors = append(neighbors, downRight)
		}
	}

	return neighbors
}
