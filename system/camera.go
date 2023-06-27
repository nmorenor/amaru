package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"

	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
)

type Camera struct {
	levelIndex int
	game       *component.GameData
}

func NewCamera(selectedLevelIndex int) *Camera {
	return &Camera{
		levelIndex: selectedLevelIndex,
	}
}

func (c *Camera) Update(w donburi.World) {
	if c.game == nil {
		c.game = component.MustFindGame(w)
		if c.game == nil {
			return
		}
	}
	camera := archetype.MustFindCamera(w)
	cam := transform.Transform.Get(camera)

	if component.Camera.Get(camera).Disabled {
		return
	}
	level := assets.AvailableLevels[c.levelIndex]
	width := level.Background.Bounds().Dx()
	height := level.Background.Bounds().Dy()
	_, player := archetype.MustFindLocalPlayer(w)

	if player != nil {
		playerTransform := transform.Transform.Get(player)

		if (cam.LocalPosition.X+c.game.LeftOffset) > (playerTransform.LocalPosition.X-float64(width/2)) || cam.LocalPosition.Y > (playerTransform.LocalPosition.Y-float64(height/2)) {
			if playerTransform.LocalPosition.X < float64(c.game.Settings.ScreenWidth) {
				cam.LocalPosition.X = (playerTransform.LocalPosition.X - float64(width/4))
			} else if cam.LocalPosition.X > (float64(width) * 0.375) {
				cam.LocalPosition.X = float64(width) * 0.375
			}

			cam.LocalPosition.Y = playerTransform.LocalPosition.Y - float64(height/4)
			return
		}

	}

	speed := c.game.Speed
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		cam.LocalPosition.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		cam.LocalPosition.X += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		cam.LocalPosition.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		cam.LocalPosition.Y += speed
	}
}
