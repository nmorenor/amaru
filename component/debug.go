package component

import (
	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type DebugData struct {
	Enabled bool
	Shapes  []*cp.Shape
}

var Debug = donburi.NewComponentType[DebugData]()

func MustFindDebug(w donburi.World) *DebugData {
	debug, ok := query.NewQuery(filter.Contains(Debug)).First(w)
	if !ok {
		panic("debug not found")
	}
	return Debug.Get(debug)
}
