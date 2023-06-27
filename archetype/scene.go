package archetype

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
	UpdateLayout(width, height int)
	NextScene() Scene
}
