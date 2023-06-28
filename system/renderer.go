package system

import (
	"sort"

	"github.com/fogleman/gg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
)

type Render struct {
	query      *query.Query
	labelQuery *query.Query
	offscreen  *ebiten.Image
	game       *component.GameData
	debug      *component.DebugData
}

func NewRenderer() *Render {
	level := assets.GameLevelLoader.CurrentLevel
	return &Render{
		query: query.NewQuery(
			filter.Contains(transform.Transform, component.Sprite),
		),
		labelQuery: query.NewQuery(
			filter.Contains(transform.Transform, component.PlayerLabel),
		),
		offscreen: ebiten.NewImage(level.Background.Bounds().Dx(), level.Background.Bounds().Dy()),
	}
}

func (r *Render) Update(w donburi.World) {
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

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	camera := archetype.MustFindCamera(w)
	cameraPos := transform.Transform.Get(camera).LocalPosition

	r.offscreen.Clear()

	var entries []*donburi.Entry
	r.query.Each(w, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	byLayer := lo.GroupBy(entries, func(entry *donburi.Entry) int {
		return int(component.Sprite.Get(entry).Layer)
	})
	layers := lo.Keys(byLayer)
	sort.Ints(layers)

	for _, layer := range layers {
		for _, entry := range byLayer[layer] {
			sprite := component.Sprite.Get(entry)

			if sprite.Hidden {
				continue
			}

			w := sprite.Image.Bounds().Dx()
			h := sprite.Image.Bounds().Dy()
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

			colorm.DrawImage(r.offscreen, sprite.Image, colormm, op)
		}
	}

	r.labelQuery.Each(w, func(entry *donburi.Entry) {
		position := transform.WorldPosition(entry)
		name := component.PlayerLabel.Get(entry).Name
		color := component.PlayerLabel.Get(entry).Color
		width, height := r.offscreen.Bounds().Dx(), r.offscreen.Bounds().Dy()
		dc := gg.NewContext(width, height)
		dc.SetFontFace(assets.MainFont)
		// Calculate the text width and height
		textWidth, textHeight := dc.MeasureString(name)
		// Calculate the starting point for centered text
		startX := (position.X + 16) - textWidth/2
		startY := (position.Y - 16) + textHeight/2
		text.Draw(
			r.offscreen,
			name,
			assets.MainFont,
			int(startX),
			int(startY),
			color,
		)
	})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(r.offscreen, op)
}
