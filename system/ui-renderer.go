package system

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/archetype"
	"amaru/component"
)

type UIRender struct {
	query     *query.Query
	offscreen *ebiten.Image
	game      *component.GameData
	debug     *component.DebugData
}

func NewUIRenderer() *UIRender {
	return &UIRender{
		query: query.NewQuery(
			filter.Contains(transform.Transform, component.UISprite),
		),
	}
}

func (r *UIRender) Initialize(w donburi.World) {
	if r.game == nil {
		r.game = component.MustFindGame(w)
		if r.game == nil {
			return
		}
	}
	r.offscreen = ebiten.NewImage(r.game.Settings.ScreenWidth, r.game.Settings.ScreenHeight)
}

func (r *UIRender) Update(w donburi.World) {
	if r.game == nil {
		r.game = component.MustFindGame(w)
		if r.game == nil {
			return
		}
	}
	if r.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).First(w)
		if !ok {
			return
		}

		r.debug = component.Debug.Get(debug)
	}
}

func (r *UIRender) Draw(w donburi.World, screen *ebiten.Image) {
	camera := archetype.MustFindCamera(w)
	cameraPos := transform.Transform.Get(camera).LocalPosition

	r.offscreen.Clear()

	var entries []*donburi.Entry
	r.query.Each(w, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	byLayer := lo.GroupBy(entries, func(entry *donburi.Entry) int {
		return int(component.UISprite.Get(entry).Layer)
	})
	layers := lo.Keys(byLayer)
	sort.Ints(layers)

	for _, layer := range layers {
		for _, entry := range byLayer[layer] {
			sprite := component.UISprite.Get(entry)

			if sprite.Hidden {
				continue
			}

			nextImage := sprite.Image

			w := nextImage.Bounds().Dx()
			h := nextImage.Bounds().Dy()
			halfW, halfH := float64(w)/2, float64(h)/2

			op := &colorm.DrawImageOptions{}
			position := transform.WorldPosition(entry)

			x := position.X
			y := position.Y

			switch sprite.Pivot {
			case component.SpritePivotCenter:
				x -= halfW
				y -= halfH
			case component.SpritePivotScreenCenter:
				x = float64(r.game.Settings.ScreenWidth-w) / 2
				y = float64(r.game.Settings.ScreenHeight-h) / 2
			}

			scale := transform.WorldScale(entry)
			op.GeoM.Translate(-halfW, -halfH)
			op.GeoM.Scale(scale.X, scale.Y)
			op.GeoM.Translate(halfW, halfH)

			colormm := colorm.ColorM{}
			if sprite.ColorOverride != nil {
				colormm.Scale(0, 0, 0, sprite.ColorOverride.A)
				colormm.Translate(sprite.ColorOverride.R, sprite.ColorOverride.G, sprite.ColorOverride.B, 0)
			}

			op.GeoM.Translate(x, y)

			colorm.DrawImage(r.offscreen, nextImage, colormm, op)
			r.offscreen = sprite.UIHandler(r.offscreen)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(r.offscreen, op)
}
