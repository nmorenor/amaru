package component

import (
	"amaru/assets"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
)

type WasteData struct {
	Id        string
	Path      assets.Path
	Collected bool
	Shape     *cp.Shape
	Entry     *donburi.Entry
}

var Waste = donburi.NewComponentType[WasteData]()
