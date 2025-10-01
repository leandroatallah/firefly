package movement

import (
	"container/heap"
	"image"
	"math"

	"github.com/leandroatallah/firefly/internal/systems/physics"
)

// ChaseMovementState implements the A* pathfinding algorithm to chase a target.
type ChaseMovementState struct {
	BaseMovementState
	count     int
	path      []image.Point
	obstacles []physics.Body
}

func NewChaseMovementState(base BaseMovementState) *ChaseMovementState {
	return &ChaseMovementState{BaseMovementState: base}
}

// --- A* Node and Priority Queue ---

// Node represents a point in the search grid for A*.
type Node struct {
	point    image.Point
	parent   *Node
	gCost    int // Distance from starting node
	hCost    int // Heuristic distance to end node
	fCost    int // gCost + hCost
	mapIndex int // Index of the item in the priority queue
}

// PriorityQueue implements heap.Interface and holds Nodes.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest fCost, so we use less than here.
	return pq[i].fCost < pq[j].fCost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].mapIndex = i
	pq[j].mapIndex = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	node := x.(*Node)
	node.mapIndex = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil     // avoid memory leak
	node.mapIndex = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

// --- Helper Functions ---
func euclideanDistance(a, b image.Point) int {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return int(math.Sqrt(dx*dx + dy*dy))
}

// reconstructPath builds the path from the end node back to the start.
func reconstructPath(endNode *Node) []image.Point {
	path := []image.Point{}
	for current := endNode; current != nil; current = current.parent {
		path = append(path, current.point)
	}
	// Reverse the path to get it from start to end
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func (s *ChaseMovementState) Move() {
	s.count++

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

	// visited nodes
	closedSet := make(map[image.Point]*Node)

	actorSize := s.actor.Position().Size()

	for openSet.Len() > 0 {
		// Get the node with the lowest F cost
		currentNode := heap.Pop(openSet).(*Node)

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

			// Temporary fCost
			tentativeGCost := currentNode.gCost + euclideanDistance(currentNode.point, neighborPoint)

			// TODO: Check it to optimize.
			// For now, we don't have a good way to check if the neighbor is in the open set
			// without iterating, which is inefficient. A* can be optimized by using a map
			// for the open set as well, but for simplicity, we'll just add the node.
			// This can lead to duplicates but is functionally closer to correct.

			neighborNode := &Node{
				point:  neighborPoint,
				parent: currentNode,
				gCost:  tentativeGCost,
				hCost:  euclideanDistance(neighborPoint, targetPos),
			}
			neighborNode.fCost = neighborNode.gCost + neighborNode.hCost
			heap.Push(openSet, neighborNode)
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

// Functional Options Pattern
// WithObstacles is an option to provide obstacles for pathfinding states.
func WithObstacles(obstacles []physics.Body) MovementStateOption {
	return func(ms MovementState) {
		if chaseState, ok := ms.(*ChaseMovementState); ok {
			chaseState.obstacles = obstacles
		}
	}
}
