package component

import (
	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
)

const (
	PlayerCollisionType cp.CollisionType = iota
	BoxCollisionType
	WasteCollisionType
	AnimalCollisionType
)

type PhysicsData struct {
	Space *cp.Space
}

var Physics = donburi.NewComponentType[PhysicsData]()
