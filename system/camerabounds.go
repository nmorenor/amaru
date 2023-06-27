package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/archetype"
	"amaru/component"
)

type CameraBounds struct {
	query *query.Query
	game  *component.GameData
}

func NewCameraBounds() *CameraBounds {
	return &CameraBounds{
		query: query.NewQuery(filter.Contains(
			component.Camera,
			transform.Transform,
		)),
	}
}

func (b *CameraBounds) Update(w donburi.World) {
	if b.game == nil {
		b.game = component.MustFindGame(w)
		if b.game == nil {
			return
		}
	}
	camera := archetype.MustFindCamera(w)
	cameraCamera := component.Camera.Get(camera)

	width, height := float64(cameraCamera.Width), float64(cameraCamera.Height)

	b.query.Each(w, func(entry *donburi.Entry) {
		t := transform.Transform.Get(entry)
		maxX := (width - float64(b.game.Settings.ScreenWidth)) + b.game.LeftOffset
		maxY := (height - float64(b.game.Settings.ScreenHeight))

		if t.LocalPosition.X < 0 {
			t.LocalPosition.X = 0
		}
		if t.LocalPosition.X > maxX {
			t.LocalPosition.X = maxX
		}

		if t.LocalPosition.Y < 0 {
			t.LocalPosition.Y = 0
		}

		if t.LocalPosition.Y > maxY {
			t.LocalPosition.Y = maxY
		}
	})
}
