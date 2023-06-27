package component

import (
	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
)

type AnimalData struct {
	X         float64
	Y         float64
	Collected bool
	Shape     *cp.Shape
	Entry     *donburi.Entry
}

var Animal = donburi.NewComponentType[AnimalData]()
