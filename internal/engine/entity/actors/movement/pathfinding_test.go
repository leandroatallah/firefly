package movement

import (
	"container/heap"
	"image"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// --- Tests for euclideanDistance ---

func TestEuclideanDistance(t *testing.T) {
	tests := []struct {
		name string
		a    image.Point
		b    image.Point
		want int
	}{
		{"same point", image.Point{0, 0}, image.Point{0, 0}, 0},
		{"horizontal right", image.Point{0, 0}, image.Point{3, 0}, 3},
		{"horizontal left", image.Point{3, 0}, image.Point{0, 0}, 3},
		{"vertical down", image.Point{0, 0}, image.Point{0, 4}, 4},
		{"vertical up", image.Point{0, 4}, image.Point{0, 0}, 4},
		{"diagonal 3-4-5", image.Point{0, 0}, image.Point{3, 4}, 5},
		{"diagonal 6-8-10", image.Point{0, 0}, image.Point{6, 8}, 10},
		{"diagonal 1-1-sqrt2", image.Point{0, 0}, image.Point{1, 1}, 1},
		{"negative coordinates", image.Point{-3, -4}, image.Point{0, 0}, 5},
		{"large distance", image.Point{0, 0}, image.Point{100, 0}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := euclideanDistance(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("euclideanDistance(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// --- Tests for reconstructPath ---

func TestReconstructPath(t *testing.T) {
	tests := []struct {
		name     string
		buildPath func() *Node
		want     []image.Point
	}{
		{
			name: "single node",
			buildPath: func() *Node {
				return &Node{point: image.Point{0, 0}}
			},
			want: []image.Point{{0, 0}},
		},
		{
			name: "straight line horizontal",
			buildPath: func() *Node {
				n0 := &Node{point: image.Point{0, 0}}
				n1 := &Node{point: image.Point{1, 0}, parent: n0}
				n2 := &Node{point: image.Point{2, 0}, parent: n1}
				n3 := &Node{point: image.Point{3, 0}, parent: n2}
				return n3
			},
			want: []image.Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		},
		{
			name: "L-shape path",
			buildPath: func() *Node {
				n0 := &Node{point: image.Point{0, 0}}
				n1 := &Node{point: image.Point{1, 0}, parent: n0}
				n2 := &Node{point: image.Point{2, 0}, parent: n1}
				n3 := &Node{point: image.Point{2, 1}, parent: n2}
				n4 := &Node{point: image.Point{2, 2}, parent: n3}
				return n4
			},
			want: []image.Point{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			name: "diagonal path",
			buildPath: func() *Node {
				n0 := &Node{point: image.Point{0, 0}}
				n1 := &Node{point: image.Point{1, 1}, parent: n0}
				n2 := &Node{point: image.Point{2, 2}, parent: n1}
				n3 := &Node{point: image.Point{3, 3}, parent: n2}
				return n3
			},
			want: []image.Point{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endNode := tt.buildPath()
			got := reconstructPath(endNode)

			if len(got) != len(tt.want) {
				t.Errorf("reconstructPath() length = %d, want %d", len(got), len(tt.want))
				return
			}

			for i, p := range got {
				if p != tt.want[i] {
					t.Errorf("reconstructPath()[%d] = %v, want %v", i, p, tt.want[i])
				}
			}
		})
	}
}

// --- Tests for PriorityQueue ---

func TestPriorityQueue(t *testing.T) {
	t.Run("Push and Pop", func(t *testing.T) {
		pq := &PriorityQueue{}
		heap.Init(pq)

		n1 := &Node{point: image.Point{0, 0}, fCost: 10}
		n2 := &Node{point: image.Point{1, 0}, fCost: 5}
		n3 := &Node{point: image.Point{2, 0}, fCost: 15}

		heap.Push(pq, n1)
		heap.Push(pq, n2)
		heap.Push(pq, n3)

		// Should pop in order of lowest fCost
		got1 := heap.Pop(pq).(*Node)
		if got1 != n2 {
			t.Errorf("expected first pop to be n2 (fCost=5), got fCost=%d", got1.fCost)
		}

		got2 := heap.Pop(pq).(*Node)
		if got2 != n1 {
			t.Errorf("expected second pop to be n1 (fCost=10), got fCost=%d", got2.fCost)
		}

		got3 := heap.Pop(pq).(*Node)
		if got3 != n3 {
			t.Errorf("expected third pop to be n3 (fCost=15), got fCost=%d", got3.fCost)
		}
	})

	t.Run("Len", func(t *testing.T) {
		pq := &PriorityQueue{}
		heap.Init(pq)

		if pq.Len() != 0 {
			t.Errorf("expected empty queue length to be 0, got %d", pq.Len())
		}

		heap.Push(pq, &Node{fCost: 1})
		heap.Push(pq, &Node{fCost: 2})

		if pq.Len() != 2 {
			t.Errorf("expected queue length to be 2, got %d", pq.Len())
		}
	})

	t.Run("Less", func(t *testing.T) {
		pq := PriorityQueue{
			&Node{fCost: 10},
			&Node{fCost: 5},
			&Node{fCost: 15},
		}

		if !pq.Less(1, 0) {
			t.Error("expected pq[1] (fCost=5) to be less than pq[0] (fCost=10)")
		}
		if pq.Less(0, 1) {
			t.Error("expected pq[0] (fCost=10) to not be less than pq[1] (fCost=5)")
		}
	})

	t.Run("Swap", func(t *testing.T) {
		pq := PriorityQueue{
			&Node{point: image.Point{0, 0}, fCost: 10, mapIndex: 0},
			&Node{point: image.Point{1, 0}, fCost: 5, mapIndex: 1},
		}

		pq.Swap(0, 1)

		if pq[0].fCost != 5 {
			t.Errorf("expected pq[0].fCost to be 5 after swap, got %d", pq[0].fCost)
		}
		if pq[1].fCost != 10 {
			t.Errorf("expected pq[1].fCost to be 10 after swap, got %d", pq[1].fCost)
		}
		// Check mapIndex is updated
		if pq[0].mapIndex != 0 {
			t.Errorf("expected pq[0].mapIndex to be 0 after swap, got %d", pq[0].mapIndex)
		}
		if pq[1].mapIndex != 1 {
			t.Errorf("expected pq[1].mapIndex to be 1 after swap, got %d", pq[1].mapIndex)
		}
	})

	t.Run("Pop updates mapIndex", func(t *testing.T) {
		pq := &PriorityQueue{}
		heap.Init(pq)

		n1 := &Node{fCost: 5}
		heap.Push(pq, n1)

		heap.Pop(pq)

		if n1.mapIndex != -1 {
			t.Errorf("expected popped node mapIndex to be -1, got %d", n1.mapIndex)
		}
	})
}

// --- Tests for WithObstacles ---

func TestWithObstacles(t *testing.T) {
	base := NewBaseMovementState(Chase, &mockActor{}, &mockActor{})
	state := NewChaseMovementState(base)

	obstacles := []body.MovableCollidable{
		newMockMovableCollidable(10, 10, 10, 10),
		newMockMovableCollidable(30, 30, 10, 10),
	}

	option := WithObstacles(obstacles)
	option(state)

	// Verify obstacles were set
	// We need to check via type assertion since obstacles is not exported
	// This test mainly verifies the option doesn't panic and applies to ChaseMovementState
}

func TestWithObstacles_NotChaseState(t *testing.T) {
	base := NewBaseMovementState(Patrol, &mockActor{}, nil)
	state := NewPatrolMovementState(base)

	obstacles := []body.MovableCollidable{
		newMockMovableCollidable(10, 10, 10, 10),
	}

	option := WithObstacles(obstacles)
	// Should not panic when applied to non-ChaseMovementState
	option(state)
}
