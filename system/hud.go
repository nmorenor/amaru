package system

import (
	"image/color"

	"github.com/fogleman/gg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/ui"
)

type HUD struct {
	query         *query.Query
	game          *component.GameData
	shadowOverlay *ebiten.Image
	controls      *ebiten.Image
	borderColor   color.Color
	hudUi         *ui.HudUi
}

func NewHUD() *HUD {
	return &HUD{
		query:       query.NewQuery(filter.Contains(component.Player)),
		borderColor: assets.BlueColor,
		hudUi:       ui.NewHudUI(),
	}
}

func (h *HUD) createRoundedRect(width, height int, borderRadius float64, borderColor color.Color, fillColor color.Color) *ebiten.Image {
	// Create a new gg.Context with the desired width and height
	dc := gg.NewContext(width, height)

	// Draw the custom path for the rectangle with two rounded corners on the right side (for fill)
	dc.MoveTo(0, 0)
	dc.LineTo(float64(width)-borderRadius, 0)
	dc.QuadraticTo(float64(width), 0, float64(width), borderRadius)
	dc.LineTo(float64(width), float64(height)-borderRadius)
	dc.QuadraticTo(float64(width), float64(height), float64(width)-borderRadius, float64(height))
	dc.LineTo(0, float64(height))
	dc.ClosePath()

	// Set the fill color and fill the path
	dc.SetColor(fillColor)
	dc.Fill()

	// Draw the custom path for the rectangle with two larger rounded corners on the right side (for border)
	borderWidth := 10.0
	dc.NewSubPath()
	dc.MoveTo(0, 0)
	dc.LineTo(float64(width)-(borderRadius+borderWidth), 0)
	dc.QuadraticTo(float64(width), 0, float64(width), borderRadius+borderWidth)
	dc.LineTo(float64(width), float64(height)-(borderRadius+borderWidth))
	dc.QuadraticTo(float64(width), float64(height), float64(width)-(borderRadius+borderWidth), float64(height))
	dc.LineTo(0, float64(height))
	dc.ClosePath()

	// Set the border color and stroke the path to create the border
	dc.SetColor(borderColor)
	dc.SetLineWidth(borderWidth) // Set the line width for the border
	dc.Stroke()

	// Convert the gg.Context to an ebiten.Image
	img := ebiten.NewImageFromImage(dc.Image())

	return img
}

func (h *HUD) Update(w donburi.World) {
	if h.game == nil {
		h.game = component.MustFindGame(w)
		if h.game == nil {
			return
		}
	}
	h.game.CursorOverButton = false
	h.hudUi.World = &w
	h.hudUi.Game = h.game

	h.hudUi.Update()
	if h.hudUi.Close {
		h.hudUi.Close = false
		h.game.Session.End = true
		go h.game.Session.RemoteClient.Client.Close()
	}
	if h.hudUi.Audio {
		h.game.Muted = !h.game.Muted
		h.hudUi.Audio = false
	}
	archetype.UpdateCursorImage(h.game.CursorOverButton)
}

func (h *HUD) Draw(w donburi.World, screen *ebiten.Image) {
	if h.shadowOverlay == nil {
		h.shadowOverlay = ebiten.NewImage(h.game.Settings.ScreenWidth, h.game.Settings.ScreenHeight)
		h.shadowOverlay.Fill(colornames.Black)
	}

	if h.controls == nil || h.controls.Bounds().Dx() != int(h.game.LeftOffset) {
		h.controls = h.createRoundedRect(int(h.game.LeftOffset), 48, 10, h.borderColor, colornames.White)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(h.game.Settings.ScreenHeight-48))
	screen.DrawImage(h.controls, op)

	h.hudUi.Draw(screen)
	archetype.UpdateCursorImage(h.game.CursorOverButton)
}
