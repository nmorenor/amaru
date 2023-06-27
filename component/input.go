package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Axis struct {
	X float64
	Y float64
}

type InputData struct {
	PrevKeyState map[ebiten.Key]bool
	Axis         *Axis
}

func (id *InputData) IsKeyReleased(key ebiten.Key) bool {
	return id.PrevKeyState[key] && !ebiten.IsKeyPressed(key)
}

var Input = donburi.NewComponentType[InputData]()
