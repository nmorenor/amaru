package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/jakecoffman/cp"

	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
)

type Debug struct {
	query          *query.Query
	debug          *component.DebugData
	space          *cp.Space
	offscreen      *ebiten.Image
	offscreenBoxes *ebiten.Image
	game           *component.GameData
}

func NewDebug(levelIndex int) *Debug {
	return &Debug{
		query: query.NewQuery(
			filter.Contains(transform.Transform, component.Sprite),
		),
		offscreen:      ebiten.NewImage(assets.AvailableLevels[levelIndex].Background.Bounds().Dx(), assets.AvailableLevels[levelIndex].Background.Bounds().Dy()),
		offscreenBoxes: ebiten.NewImage(assets.AvailableLevels[levelIndex].Background.Bounds().Dx(), assets.AvailableLevels[levelIndex].Background.Bounds().Dy()),
	}
}

func (d *Debug) Update(w donburi.World) {
	if d.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).First(w)
		if !ok {
			return
		}

		d.debug = component.Debug.Get(debug)
	}
	if d.game == nil {
		d.game = component.MustFindGame(w)
		if d.game == nil {
			return
		}
	}
	if d.space == nil {
		pdata, _ := archetype.MustFindPhysics(w)
		d.space = pdata.Space
		if d.space == nil {
			return
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		d.debug.Enabled = !d.debug.Enabled
	}
}

func (d *Debug) drawBox(screen *ebiten.Image, shape *cp.Shape, clr color.Color) {
	box := shape.Class.(*cp.PolyShape)
	lineWidth := float32(5)

	body := shape.Body()

	for i := 0; i < 4; i++ {
		vert1 := box.Vert(i)
		vert2 := box.Vert((i + 1) % 4)

		// Convert the vertices to world coordinates
		worldVert1 := body.LocalToWorld(vert1)
		worldVert2 := body.LocalToWorld(vert2)

		vector.StrokeLine(screen, float32(worldVert1.X), float32(worldVert1.Y+32), float32(worldVert2.X), float32(worldVert2.Y+32), lineWidth, clr, false)
	}
}

func (d *Debug) Draw(w donburi.World, screen *ebiten.Image) {
	if d.debug == nil || !d.debug.Enabled {
		return
	}
	d.offscreen.Clear()
	d.offscreenBoxes.Clear()
	op := &ebiten.DrawImageOptions{}
	if d.debug.Enabled {
		boxColor := colornames.Lime
		for _, shape := range d.debug.Shapes {
			if shape.UserData != nil {
				entry := shape.UserData.(*donburi.Entry)
				if !w.Valid(entry.Entity()) {
					continue
				}
				if entry.HasComponent(component.Animal) {
					boxColor = colornames.Fuchsia
				}
				if entry.HasComponent(component.Waste) {
					boxColor = colornames.Orange
				}
				if entry.HasComponent(component.Player) && component.Player.Get(entry).Local {
					boxColor = colornames.Red
				}
				if entry.HasComponent(component.Player) && !component.Player.Get(entry).Local {
					boxColor = colornames.Cyan
				}
			}
			d.drawBox(d.offscreenBoxes, shape, boxColor)
		}
	}
	d.offscreen.DrawImage(d.offscreenBoxes, op)

	camera := archetype.MustFindCamera(w)
	cameraPos := transform.Transform.Get(camera).LocalPosition
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(d.offscreen, op)
}
