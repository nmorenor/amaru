package archetype

import (
	"amaru/component"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

func NewInput(w donburi.World) *donburi.Entry {
	input := w.Entry(
		w.Create(
			component.Input,
		),
	)

	inputInput := component.InputData{
		PrevKeyState: make(map[ebiten.Key]bool),
		Axis:         &component.Axis{X: 0, Y: 0},
	}

	component.Input.SetValue(input, inputInput)

	return input
}

func MustFindInput(w donburi.World) *donburi.Entry {
	input, ok := query.NewQuery(filter.Contains(component.Input)).First(w)
	if !ok {
		panic("no input found")
	}
	return input
}
