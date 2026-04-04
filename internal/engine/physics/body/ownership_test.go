package body

import (
	"testing"
)

func TestOwnership_SetOwner(t *testing.T) {
	var o Ownership

	owner := "test-owner"
	o.SetOwner(owner)

	if o.owner != owner {
		t.Errorf("expected owner '%v'; got '%v'", owner, o.owner)
	}
}

func TestOwnership_Owner(t *testing.T) {
	var o Ownership

	if o.Owner() != nil {
		t.Errorf("expected nil owner by default; got %v", o.Owner())
	}

	owner := "my-owner"
	o.SetOwner(owner)

	if o.Owner() != owner {
		t.Errorf("expected owner '%v'; got '%v'", owner, o.Owner())
	}
}

func TestOwnership_LastOwner_NoOwner(t *testing.T) {
	var o Ownership

	last := o.LastOwner()
	if last != nil {
		t.Errorf("expected nil LastOwner; got %v", last)
	}
}

func TestOwnership_LastOwner_SingleOwner(t *testing.T) {
	var o Ownership
	owner := "direct-owner"
	o.SetOwner(owner)

	last := o.LastOwner()
	if last != owner {
		t.Errorf("expected LastOwner '%v'; got '%v'", owner, last)
	}
}

func TestOwnership_LastOwner_TwoLevelHierarchy(t *testing.T) {
	var child Ownership
	var parent Ownership

	parent.SetOwner("grandparent")
	child.SetOwner(&parent)

	// LastOwner traverses to the top-most owner
	last := child.LastOwner()
	if last != "grandparent" {
		t.Errorf("expected LastOwner grandparent; got %v", last)
	}
}

func TestOwnership_LastOwner_ThreeLevelHierarchy(t *testing.T) {
	var level1 Ownership
	var level2 Ownership
	var level3 Ownership

	level1.SetOwner("root")
	level2.SetOwner(&level1)
	level3.SetOwner(&level2)

	// LastOwner traverses to the top-most owner
	last := level3.LastOwner()
	if last != "root" {
		t.Errorf("expected LastOwner root; got %v", last)
	}
}

func TestOwnership_LastOwner_WithNonOwnershipOwner(t *testing.T) {
	var o Ownership

	// Owner that doesn't implement Owner() method
	owner := "string-owner"
	o.SetOwner(owner)

	last := o.LastOwner()
	if last != owner {
		t.Errorf("expected LastOwner '%v'; got '%v'", owner, last)
	}
}

func TestOwnership_LastOwner_CircularReference(t *testing.T) {
	var a Ownership
	var b Ownership

	// Create circular reference: a -> b -> a
	a.SetOwner(&b)
	b.SetOwner(&a)

	// Should not hang - should detect cycle and return
	last := a.LastOwner()
	if last == nil {
		t.Error("expected LastOwner to return something, not nil")
	}
}

func TestOwnership_LastOwner_SelfReference(t *testing.T) {
	var o Ownership

	// Self-reference
	o.SetOwner(&o)

	// Should not hang - should detect cycle and return
	last := o.LastOwner()
	if last == nil {
		t.Error("expected LastOwner to return something, not nil")
	}
}

func TestOwnership_LastOwner_ComplexHierarchy(t *testing.T) {
	// Create: leaf -> mid1 -> mid2 -> root -> "root-owner"
	var leaf Ownership
	var mid1 Ownership
	var mid2 Ownership
	var root Ownership

	root.SetOwner("root-owner")
	mid2.SetOwner(&root)
	mid1.SetOwner(&mid2)
	leaf.SetOwner(&mid1)

	// LastOwner traverses to the top-most owner
	last := leaf.LastOwner()
	if last != "root-owner" {
		t.Errorf("expected LastOwner root-owner; got %v", last)
	}
}

func TestOwnership_LastOwner_NilInChain(t *testing.T) {
	var level1 Ownership
	var level2 Ownership

	level1.SetOwner(nil)
	level2.SetOwner(&level1)

	last := level2.LastOwner()
	if last != &level1 {
		t.Errorf("expected LastOwner to be level1; got %v", last)
	}
}

func TestOwnership_ChangeOwner(t *testing.T) {
	var o Ownership

	o.SetOwner("owner1")
	if o.Owner() != "owner1" {
		t.Errorf("expected owner1; got %v", o.Owner())
	}

	o.SetOwner("owner2")
	if o.Owner() != "owner2" {
		t.Errorf("expected owner2; got %v", o.Owner())
	}
}

func TestOwnership_LastOwner_AfterOwnerChange(t *testing.T) {
	var child Ownership
	var parent1 Ownership
	var parent2 Ownership

	child.SetOwner(&parent1)
	parent1.SetOwner("grandparent1")

	// LastOwner returns the top-most owner (grandparent1)
	last := child.LastOwner()
	if last != "grandparent1" {
		t.Errorf("expected LastOwner grandparent1; got %v", last)
	}

	// Change parent's owner
	parent1.SetOwner(&parent2)
	parent2.SetOwner("grandparent2")

	// Now LastOwner should traverse to grandparent2
	last = child.LastOwner()
	if last != "grandparent2" {
		t.Errorf("expected LastOwner grandparent2 after change; got %v", last)
	}
}

func TestOwnership_MultipleInstances(t *testing.T) {
	var o1, o2, o3 Ownership

	o1.SetOwner("owner1")
	o2.SetOwner("owner2")
	o3.SetOwner("owner3")

	if o1.Owner() != "owner1" {
		t.Errorf("o1: expected owner1; got %v", o1.Owner())
	}
	if o2.Owner() != "owner2" {
		t.Errorf("o2: expected owner2; got %v", o2.Owner())
	}
	if o3.Owner() != "owner3" {
		t.Errorf("o3: expected owner3; got %v", o3.Owner())
	}
}

// Test with actual body types to ensure integration works
func TestOwnership_WithBodyTypes(t *testing.T) {
	body := NewBody(NewRect(0, 0, 10, 10))
	movable := NewMovableBody(body)
	collidable := NewCollidableBody(body)

	// Set up ownership chain
	body.SetOwner(movable)
	movable.SetOwner(collidable)

	last := body.LastOwner()
	if last != collidable {
		t.Errorf("expected LastOwner to be collidable; got %v", last)
	}
}

func TestOwnership_WithAliveBody(t *testing.T) {
	body := NewBody(NewRect(0, 0, 10, 10))
	alive := NewAliveBody(body)

	body.SetOwner(alive)

	last := body.LastOwner()
	if last != alive {
		t.Errorf("expected LastOwner to be alive; got %v", last)
	}
}
