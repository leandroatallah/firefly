package body

// TODO: Should these be concrete?
type MovableCollidable interface {
	Movable
	Collidable
}

type MovableCollidableTouchable interface {
	MovableCollidable
	Touchable
}

type MovableCollidableAlive interface {
	MovableCollidable
	Alive
}

type MovableCollidableTouchableAlive interface {
	MovableCollidableTouchable
	Alive
}
