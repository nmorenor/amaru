package ui

import (
	"amaru/archetype"
	"amaru/assets"

	"golang.org/x/image/colornames"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	menuTitle    = "Welcome!"
	hostLabel    = "Host"
	joinLabel    = "Join"
	cancelLabel  = "Cancel"
	refreshLabel = "Refresh"
	aboutLabel   = "About"
)

type StartMenuOption int

const (
	NoOption StartMenuOption = iota
	Host
	Join
	About
)

type StartMenu struct {
	SelectedOption StartMenuOption
	container      *widget.Container
	Ui             *ebitenui.UI
	hostButton     *widget.Button
	joinButton     *widget.Button
	aboutButton    *widget.Button
}

func NewStartMenu() *StartMenu {
	startMenu := &StartMenu{
		container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
	}

	parentContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(25, 5),
		)),
	)

	welcomeLabel := widget.NewLabel(widget.LabelOpts.Text(menuTitle, assets.MainMidFont, &widget.LabelColor{
		Disabled: assets.BlueColor,
		Idle:     assets.BlueColor,
	}))

	welcomeLabelContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
	)))

	welcomeLabelContainer.AddChild(welcomeLabel)
	welcomeLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	buttonsContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(2),
		widget.GridLayoutOpts.Spacing(10, 3),
		widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
	)))

	startMenu.hostButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(hostLabel, assets.MainFont, &widget.ButtonTextColor{
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
			startMenu.SelectedOption = Host
		}),
	)
	buttonsContainer.AddChild(startMenu.hostButton)

	startMenu.joinButton = widget.NewButton(
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
			startMenu.SelectedOption = Join
		}),
	)
	buttonsContainer.AddChild(startMenu.joinButton)

	aboutContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(1),
		widget.GridLayoutOpts.Spacing(10, 3),
		widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
	)))

	startMenu.aboutButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(aboutLabel, assets.MainFont, &widget.ButtonTextColor{
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
			startMenu.SelectedOption = About
		}),
	)
	aboutContainer.AddChild(startMenu.aboutButton)

	parentContainer.AddChild(welcomeLabelContainer)
	parentContainer.AddChild(buttonsContainer)
	parentContainer.AddChild(aboutContainer)

	startMenu.container.AddChild(parentContainer)
	parentContainer.GetWidget().LayoutData = widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	startMenu.Ui = &ebitenui.UI{
		Container: startMenu.container,
	}
	return startMenu
}

func (s *StartMenu) Draw(screen *ebiten.Image) {
	s.Ui.Draw(screen)
}

func (s *StartMenu) Container() *widget.Container {
	return s.container
}

func (s *StartMenu) Update() {
	s.Ui.Update()
	hostButtonRect := s.hostButton.GetWidget().Rect
	joinButtonRect := s.joinButton.GetWidget().Rect
	aboutButtonRect := s.aboutButton.GetWidget().Rect
	mx, my := ebiten.CursorPosition()
	if (hostButtonRect.Min.X <= mx && mx <= hostButtonRect.Max.X && hostButtonRect.Min.Y <= my && my <= hostButtonRect.Max.Y) ||
		(joinButtonRect.Min.X <= mx && mx <= joinButtonRect.Max.X && joinButtonRect.Min.Y <= my && my <= joinButtonRect.Max.Y) ||
		(aboutButtonRect.Min.X <= mx && mx <= aboutButtonRect.Max.X && aboutButtonRect.Min.Y <= my && my <= aboutButtonRect.Max.Y) {
		archetype.UpdateCursorImage(true)
	} else {
		archetype.UpdateCursorImage(false)
	}
}
