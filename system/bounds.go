package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
)

type Bounds struct {
	levelIndex int
	query      *query.Query
	game       *component.GameData
}

func NewBounds(selectedLevelIndex int) *Bounds {
	return &Bounds{
		levelIndex: selectedLevelIndex,
		query: query.NewQuery(filter.Contains(
			transform.Transform,
			component.Sprite,
			component.Bounds,
		)),
	}
}

func (b *Bounds) Update(w donburi.World) {
	if b.game == nil {
		b.game = component.MustFindGame(w)
		if b.game == nil {
			return
		}
	}

	camera := archetype.MustFindCamera(w)
	cameraPos := transform.Transform.Get(camera).LocalPosition

	b.query.Each(w, func(entry *donburi.Entry) {
		bounds := component.Bounds.Get(entry)
		if bounds.Disabled {
			return
		}

		t := transform.Transform.Get(entry)
		sprite := component.Sprite.Get(entry)

		w := sprite.Image.Bounds().Dx()
		h := sprite.Image.Bounds().Dy()

		width, height := float64(w), float64(h)

		var minX, maxX, minY, maxY float64
		level := assets.AvailableLevels[b.levelIndex]
		levelWidth := level.Background.Bounds().Dx()
		levelHeight := level.Background.Bounds().Dy()

		switch sprite.Pivot {
		case component.SpritePivotTopLeft:
			minX = 0
			maxX = float64(levelWidth) - width

			minY = 0
			maxY = float64(levelHeight) - height
		case component.SpritePivotCenter:
			minX = cameraPos.X + width/2
			maxX = float64(levelWidth) - width/2

			minY = 0
			maxY = float64(levelHeight) - height/2
		}

		targetX, outx := engine.Clamp(t.LocalPosition.X, minX, maxX)
		targetY, outy := engine.Clamp(t.LocalPosition.Y, minY, maxY)
		if (outx || outy) && entry.HasComponent(component.Player) && component.Player.Get(entry).Local {
			lastDirection := component.Player.Get(entry).LastDirection
			if lastDirection != nil && lastDirection.X == -1 && t.LocalPosition.X <= minX {
				component.Player.Get(entry).OutOfBounds = true
			} else if lastDirection != nil && lastDirection.X == 1 && t.LocalPosition.X >= maxX {
				component.Player.Get(entry).OutOfBounds = true
			} else if lastDirection != nil && lastDirection.Y == -1 && t.LocalPosition.Y <= minY {
				component.Player.Get(entry).OutOfBounds = true
			} else if lastDirection != nil && lastDirection.Y == 1 && t.LocalPosition.Y >= maxY {
				component.Player.Get(entry).OutOfBounds = true
			} else {
				component.Player.Get(entry).OutOfBounds = false
			}
		} else if !(outx || outy) && (entry.HasComponent(component.Player) && component.Player.Get(entry).Local) {
			component.Player.Get(entry).OutOfBounds = false
		}
		t.LocalPosition.X = targetX
		t.LocalPosition.Y = targetY
	})
}
