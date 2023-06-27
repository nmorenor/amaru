package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/engine"
)

type LevelData struct {
	ProgressionTimer *engine.Timer
	ReachedEnd       bool
	Progressed       bool
}

var Level = donburi.NewComponentType[LevelData]()

func MustFindLevel(w donburi.World) *donburi.Entry {
	level, ok := query.NewQuery(filter.Contains(Level)).First(w)
	if !ok {
		panic("no level found")
	}

	return level
}
