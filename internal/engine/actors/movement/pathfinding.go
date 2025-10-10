package movement

import (
	"image"
	"math"

	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

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

// Functional Options Pattern
// WithObstacles is an option to provide obstacles for pathfinding states.
func WithObstacles(obstacles []physics.Body) MovementStateOption {
	return func(ms MovementState) {
		if chaseState, ok := ms.(*ChaseMovementState); ok {
			chaseState.obstacles = obstacles
		}
	}
}
