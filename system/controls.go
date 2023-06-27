package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"amaru/archetype"
	"amaru/component"
)

type Controls struct {
	game *component.GameData
}

func NewControls() *Controls {
	return &Controls{}
}

func (i *Controls) Update(w donburi.World) {
	if i.game == nil {
		i.game = component.MustFindGame(w)
		if i.game == nil {
			return
		}
	}
	input := component.Input.Get(archetype.MustFindInput(w))
	player, _ := archetype.MustFindLocalPlayer(w)

	input.Axis.X = 0
	input.Axis.Y = 0
	if input.IsKeyReleased(player.PlayerSettings.Inputs.Right) || ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) || archetype.GamePadIsRight(i.game) {
		input.Axis.X = 1
	} else if input.IsKeyReleased(player.PlayerSettings.Inputs.Left) || ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) || archetype.GamePadIsLeft(i.game) {
		input.Axis.X = -1
	}

	if input.IsKeyReleased(player.PlayerSettings.Inputs.Up) || ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) || archetype.GamePadIsUp(i.game) {
		input.Axis.Y = -1
	} else if input.IsKeyReleased(player.PlayerSettings.Inputs.Down) || ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) || archetype.GamePadIsDown(i.game) {
		input.Axis.Y = 1
	}

	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		input.PrevKeyState[key] = ebiten.IsKeyPressed(key)
	}

}
