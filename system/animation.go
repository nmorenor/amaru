package system

import (
	"amaru/component"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Animation struct {
	query *query.Query
	game  *component.GameData
}

func NewAnimation() *Animation {
	return &Animation{
		query: query.NewQuery(filter.Contains(
			component.AnimationComponent,
		)),
	}
}

func (a *Animation) Update(w donburi.World) {
	if a.game == nil {
		a.game = component.MustFindGame(w)
		if a.game == nil {
			return
		}
	}
	a.query.Each(w, func(entry *donburi.Entry) {
		e := component.AnimationComponent.Get(entry)
		if e.CurrentAnimation == nil {
			if e.Def == nil {
				return
			}
			e.SelectAnimationByAction(e.Def)
		}

		e.Change += (1.0 / 60.0)
		if e.Change >= e.Rate {
			component.Sprite.Get(entry).Image = e.Cell()
			e.NextFrame()
		}
	})
}
