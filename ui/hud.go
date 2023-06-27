package ui

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

const (
	sendPlaceHolder = "Send"
)

type listResources struct {
	image        *widget.ScrollContainerImage
	track        *widget.SliderTrackImage
	trackPadding widget.Insets
	handle       *widget.ButtonImage
	handleSize   int
	entryPadding widget.Insets
}

type HudUi struct {
	container          *widget.Container
	ui                 *ebitenui.UI
	closeButton        *widget.Button
	audioButton        *widget.Button
	Close              bool
	Audio              bool
	Game               *component.GameData
	World              *donburi.World
	remainingTimeLabel *widget.Label
	playerPointsLabel  *widget.Label
}

func NewHudUI() *HudUi {

	container := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)
	hudUi := &HudUi{}

	playerPointsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(0),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left: 10,
				Top:  5,
			}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
		})),
	)
	hudUi.playerPointsLabel = widget.NewLabel(
		widget.LabelOpts.Text("0", assets.MainBigFont, &widget.LabelColor{
			Disabled: colornames.White,
			Idle:     colornames.White,
		}),
	)
	playerPointsContainer.AddChild(hudUi.playerPointsLabel)
	container.AddChild(playerPointsContainer)

	remainingContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{
			Top: 5,
		}))),
	)
	hudUi.remainingTimeLabel = widget.NewLabel(
		widget.LabelOpts.Text("00", assets.MainBigFont, &widget.LabelColor{
			Disabled: colornames.White,
			Idle:     colornames.White,
		}),
	)
	hudUi.remainingTimeLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionStart,
	}
	remainingContainer.AddChild(hudUi.remainingTimeLabel)
	container.AddChild(remainingContainer)

	compositeContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//Set how much padding before displaying content
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
		})),
	)

	container.AddChild(compositeContainer)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(-10),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   2,
				Bottom: 1,
			}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
		})),
	)

	closeButtonContainer, closeButton := archetype.CreateSVGImageButton(assets.CloseKey, "Close", nil, func() {
		hudUi.Close = true
	})
	audioButtonContainer, audioButton := archetype.CreateSVGImageButton(assets.AudioKey, "Audio", nil, func() {
		hudUi.Audio = true
	})

	// Add the button containers to the row layout
	buttonContainer.AddChild(closeButtonContainer)
	buttonContainer.AddChild(audioButtonContainer)

	// Add the buttons to the AnchorLayout
	compositeContainer.AddChild(buttonContainer)

	hudUi.container = container
	hudUi.closeButton = closeButton
	hudUi.audioButton = audioButton

	hudUi.ui = &ebitenui.UI{
		Container: hudUi.container,
	}

	return hudUi
}

func newTextArea(text string, widgetOpts ...widget.WidgetOpt) *widget.TextArea {
	width := 64.0
	height := 64.0
	areaImage := ebiten.NewImage(int(width), int(height))
	borderWidth := 2.0
	aimg := archetype.DrawRoundedRect(areaImage, borderWidth, borderWidth, width-2*borderWidth, height-2*borderWidth, 1, colornames.White, assets.BlueColor, borderWidth)
	track := archetype.NewColoredEbitenImage(4, 40, archetype.ColorToRGBA(assets.BlueColor))
	handle := archetype.NewColoredEbitenImage(5, 5, archetype.ColorToRGBA(assets.BlueColor))
	lres := &listResources{
		image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(aimg, [3]int{12, 12, 12}, [3]int{16, 16, 32}),
			Disabled: image.NewNineSlice(aimg, [3]int{12, 12, 12}, [3]int{12, 12, 12}),
			Mask:     image.NewNineSlice(aimg, [3]int{12, 12, 12}, [3]int{12, 12, 12}),
		},

		track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(track, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(track, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(track, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		trackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handle, 0, 5),
			Hover:    image.NewNineSliceSimple(handle, 0, 5),
			Pressed:  image.NewNineSliceSimple(handle, 0, 5),
			Disabled: image.NewNineSliceSimple(handle, 0, 5),
		},

		handleSize: 5,

		entryPadding: widget.Insets{
			Left:   8,
			Right:  8,
			Top:    2,
			Bottom: 8,
		},
	}
	return widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.TextAreaOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(lres.image)),
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(lres.track, lres.handle),
			widget.SliderOpts.MinHandleSize(lres.handleSize),
			widget.SliderOpts.TrackPadding(lres.trackPadding),
		),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.VerticalScrollMode(widget.PositionAtEnd),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontFace(assets.MainFont),
		widget.TextAreaOpts.FontColor(colornames.Black),
		widget.TextAreaOpts.TextPadding(lres.entryPadding),
		widget.TextAreaOpts.Text(text),
	)
}

func (s *HudUi) Reset() {
}

func (s *HudUi) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
}

func (s *HudUi) Container() *widget.Container {
	return s.container
}

func (s *HudUi) Update() {
	if s.World != nil {
		if s.Game == nil {
			s.Game = component.MustFindGame(*s.World)
		}
		s.Game.Session.RemoteClient.GameData.Frames++

		if int(ebiten.ActualTPS()) > 0 && s.Game.Session.RemoteClient.GameData.Frames%(int(ebiten.ActualTPS())) == 0 {
			s.Game.Session.RemoteClient.GameData.Counter--
		}

		s.remainingTimeLabel.Label = fmt.Sprintf("%02d", s.Game.Session.RemoteClient.GameData.Counter)

		player, _ := archetype.MustFindLocalPlayer(*s.World)
		if player != nil {
			sessionParticipant := s.Game.Session.RemoteClient.GameData.SessionParticipants[player.ID]
			if sessionParticipant != nil {
				s.playerPointsLabel.Label = fmt.Sprintf("%d", sessionParticipant.Score)
			}
		}

		game := component.MustFindGame(*s.World)
		if game != nil {
			if s.Game.Session.Type == component.SessionTypeHost && s.Game.Session.RemoteClient.GameData.Counter <= 0 || game.WasteSize == game.CollectedWaste {
				game.GameOver = true
			}
			s.ui.Container.GetWidget().LayoutData = widget.RowLayoutData{
				MaxWidth:  int(game.LeftOffset),
				MaxHeight: game.Settings.ScreenHeight,
			}

			closeButtonRect := s.closeButton.GetWidget().Rect
			audioButtonRect := s.audioButton.GetWidget().Rect

			mx, my := ebiten.CursorPosition()
			if game.CursorOverButton {
				return
			}
			cursorOverButton := (closeButtonRect.Min.X <= mx && mx <= closeButtonRect.Max.X && closeButtonRect.Min.Y <= my && my <= closeButtonRect.Max.Y) ||
				(audioButtonRect.Min.X <= mx && mx <= audioButtonRect.Max.X && audioButtonRect.Min.Y <= my && my <= audioButtonRect.Max.Y)

			if cursorOverButton != game.CursorOverButton {
				game.CursorOverButton = cursorOverButton
			}
		}
	}
	s.ui.Update()
}
