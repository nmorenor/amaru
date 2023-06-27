package ui

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/net"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"golang.org/x/image/colornames"
)

const (
	menuSessions = "Available Sessions:"
)

type AvailableSessionsMenu struct {
	container      *widget.Container
	Ui             *ebitenui.UI
	cancelButton   *widget.Button
	refreshButton  *widget.Button
	joinButton     *widget.Button
	sessions       *[]net.AvailableSession
	listLayoutData widget.RowLayoutData
	list           *widget.List
	gameData       *component.GameData
	shouldUpdate   bool
	Session        *net.AvailableSession
	Done           bool
	Cancelled      bool
}

func NewAvailableSessionsMenu(gameData *component.GameData) *AvailableSessionsMenu {
	availableSessionsMenu := &AvailableSessionsMenu{
		container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
		sessions: &[]net.AvailableSession{},
		gameData: gameData,
	}

	parentContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(25, 5),
		)),
	)
	parentContainer.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	availableSessionsLabel := widget.NewLabel(widget.LabelOpts.Text(menuSessions, assets.MainMidFont, &widget.LabelColor{
		Disabled: assets.BlueColor,
		Idle:     assets.BlueColor,
	}))

	availableSessionsLabelContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
	)))

	availableSessionsLabelContainer.AddChild(availableSessionsLabel)
	availableSessionsLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}
	entries := make([]any, 0, len(*availableSessionsMenu.sessions))
	track := archetype.NewColoredEbitenImage(5, 5, archetype.ColorToRGBA(colornames.Fuchsia))
	trackButtonImage := &widget.ButtonImage{
		Idle:     image.NewNineSliceSimple(track, 0, 5),
		Hover:    image.NewNineSliceSimple(track, 0, 5),
		Pressed:  image.NewNineSliceSimple(track, 0, 5),
		Disabled: image.NewNineSliceSimple(track, 0, 5),
	}
	availableSessionsMenu.list = widget.NewList(
		//Set how wide the list should be
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchVertical:    true,
			}),
		)),
		//Set the entries in the list
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			//Set the background images/color for the list
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.White),
				Disabled: image.NewNineSliceColor(color.White),
				Mask:     image.NewNineSliceColor(color.White),
			}),
		),
		widget.ListOpts.SliderOpts(
			//Set the background images/color for the background of the slider track
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(colornames.Fuchsia),
				Hover: image.NewNineSliceColor(colornames.Fuchsia),
			}, trackButtonImage),
			widget.SliderOpts.MinHandleSize(5),
			//Set how wide the track should be
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		//Hide the horizontal slider
		widget.ListOpts.HideHorizontalSlider(),
		//Set the font for the list options
		widget.ListOpts.EntryFontFace(assets.MainMidFont),
		//Set the colors for the list
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   colornames.White, //Foreground color for the unfocused selected entry
			Unselected:                 color.Black,      //Foreground color for the unfocused unselected entry
			SelectedBackground:         assets.BlueColor, //Background color for the unfocused selected entry
			SelectedFocusedBackground:  assets.BlueColor, //Background color for the focused selected entry
			FocusedBackground:          assets.BlueColor, //Background color for the focused unselected entry
			DisabledUnselected:         assets.BlueColor, //Foreground color for the disabled unselected entry
			DisabledSelected:           assets.BlueColor, //Foreground color for the disabled selected entry
			DisabledSelectedBackground: assets.BlueColor, //Background color for the disabled selected entry
		}),
		//This required function returns the string displayed in the list
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(net.AvailableSession).SessionHostName
		}),
		//Padding for each entry
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		//This handler defines what function to run when a list item is selected.
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			selectedSession := args.Entry.(net.AvailableSession)
			availableSessionsMenu.Session = &selectedSession
		}),
	)

	availableSessionsMenu.listLayoutData = widget.RowLayoutData{
		Stretch:   false,
		MaxWidth:  (gameData.Settings.ScreenHeight / 2) - 24,
		MaxHeight: (gameData.Settings.ScreenHeight / 2) - 16,
	}
	availableSessionsMenu.list.GetWidget().LayoutData = availableSessionsMenu.listLayoutData

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Spacing(10, 3),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true}, []bool{true}),
		),
		))
	buttonsContainer.GetWidget().LayoutData = widget.GridLayoutData{
		HorizontalPosition: widget.GridLayoutPositionCenter,
	}

	availableSessionsMenu.cancelButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(cancelLabel, assets.MainFont, &widget.ButtonTextColor{
			Idle:     assets.BlueColor,
			Disabled: assets.BlueColor,
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}),
		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			availableSessionsMenu.Cancelled = true
			availableSessionsMenu.Done = false
		}),
	)
	buttonsContainer.AddChild(availableSessionsMenu.cancelButton)

	availableSessionsMenu.refreshButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(refreshLabel, assets.MainFont, &widget.ButtonTextColor{
			Idle:     assets.BlueColor,
			Disabled: assets.BlueColor,
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}),
		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			availableSessionsMenu.refreshSessions()
		}),
	)
	buttonsContainer.AddChild(availableSessionsMenu.refreshButton)

	availableSessionsMenu.joinButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(joinLabel, assets.MainFont, &widget.ButtonTextColor{
			Idle:     assets.BlueColor,
			Disabled: assets.BlueColor,
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}),
		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			availableSessionsMenu.Done = true
			availableSessionsMenu.Cancelled = true
		}),
	)
	buttonsContainer.AddChild(availableSessionsMenu.joinButton)

	parentContainer.AddChild(availableSessionsLabelContainer)
	parentContainer.AddChild(availableSessionsMenu.list)
	parentContainer.AddChild(buttonsContainer)

	availableSessionsMenu.container.AddChild(parentContainer)
	parentContainer.GetWidget().LayoutData = widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	availableSessionsMenu.Ui = &ebitenui.UI{
		Container: availableSessionsMenu.container,
	}

	availableSessionsMenu.refreshSessions()

	return availableSessionsMenu
}

func (s *AvailableSessionsMenu) refreshSessions() {
	go func() {
		s.sessions = net.GetAvailableSessions()
		if s.sessions != nil {
			filtered := lo.Filter(*s.sessions, func(session net.AvailableSession, index int) bool {
				return session.Size < component.MaxPlayers
			})
			s.sessions = &filtered
		}
		s.shouldUpdate = s.sessions != nil
	}()
}

func (s *AvailableSessionsMenu) Draw(screen *ebiten.Image) {
	s.Ui.Draw(screen)
}

func (s *AvailableSessionsMenu) Container() *widget.Container {
	return s.container
}

func (s *AvailableSessionsMenu) Update() {
	if s.shouldUpdate {
		entries := make([]any, 0, len(*s.sessions))
		for _, session := range *s.sessions {
			entries = append(entries, session)
		}
		s.list.SetEntries(entries)
		s.shouldUpdate = false
	}
	listWidget := s.list.GetWidget()
	listWidget.MinWidth = (s.gameData.Settings.ScreenHeight / 2) + 64
	listWidget.MinHeight = (s.gameData.Settings.ScreenHeight / 2) - 128
	s.listLayoutData.MaxWidth = (s.gameData.Settings.ScreenHeight / 2) + 64
	s.listLayoutData.MaxHeight = (s.gameData.Settings.ScreenHeight / 2) - 128
	listWidget.LayoutData = s.listLayoutData

	s.Ui.Update()
	cancelButtonRect := s.cancelButton.GetWidget().Rect
	joinButtonRect := s.joinButton.GetWidget().Rect
	refreshButtonRect := s.refreshButton.GetWidget().Rect
	mx, my := ebiten.CursorPosition()
	if (cancelButtonRect.Min.X <= mx && mx <= cancelButtonRect.Max.X && cancelButtonRect.Min.Y <= my && my <= cancelButtonRect.Max.Y) ||
		(joinButtonRect.Min.X <= mx && mx <= joinButtonRect.Max.X && joinButtonRect.Min.Y <= my && my <= joinButtonRect.Max.Y) ||
		(refreshButtonRect.Min.X <= mx && mx <= refreshButtonRect.Max.X && refreshButtonRect.Min.Y <= my && my <= refreshButtonRect.Max.Y) {
		archetype.UpdateCursorImage(true)
	} else {
		archetype.UpdateCursorImage(false)
	}
}
