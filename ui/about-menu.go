package ui

import (
	"amaru/archetype"
	"amaru/assets"
	"encoding/json"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/browser"
)

const (
	aboutTitle = "About"
	backLabel  = "Back"
)

type AboutData struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Collaborators struct {
	Collabs []*AboutData `json:"collabs"`
}

type AboutMenuUI struct {
	container     *widget.Container
	Ui            *ebitenui.UI
	backButton    *widget.Button
	collabButtons []*widget.Button
	collabs       *Collaborators
	collabHovered bool
	Back          bool
}

func NewAboutMenuUI() *AboutMenuUI {
	var collaboratorsData Collaborators
	if err := json.Unmarshal(assets.AboutData, &collaboratorsData); err != nil {
		panic(err)
	}
	aboutMenu := &AboutMenuUI{
		container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		),
		collabs: &collaboratorsData,
	}

	parentContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(25, 5),
		)),
	)

	aboutLabel := widget.NewLabel(widget.LabelOpts.Text(aboutTitle, assets.MainMidFont, &widget.LabelColor{
		Disabled: assets.BlueColor,
		Idle:     assets.BlueColor,
	}))

	aboutLabelContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
	)))

	aboutLabelContainer.AddChild(aboutLabel)
	aboutLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	buttonsContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(2),
		widget.GridLayoutOpts.Spacing(10, 3),
	)))

	normal := &widget.ButtonTextColor{
		Idle:     assets.BlueColor,
		Disabled: assets.BlueColor,
	}
	hovered := &widget.ButtonTextColor{
		Idle:     assets.GreenColor,
		Disabled: assets.GreenColor,
	}
	for next := range collaboratorsData.Collabs {
		collab := collaboratorsData.Collabs[next]

		collabLabel := widget.NewLabel(
			widget.LabelOpts.Text(collab.Description, assets.MainFont, &widget.LabelColor{
				Disabled: assets.BlueColor,
				Idle:     assets.BlueColor,
			},
			))
		collabLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		}
		buttonsContainer.AddChild(collabLabel)

		collabButton := widget.NewButton(
			widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, color.White, color.White, color.White, 5)),
			widget.ButtonOpts.Text(collab.Name, assets.MainFont, normal),
			widget.ButtonOpts.TextPadding(widget.Insets{
				Top:    0,
				Bottom: 0,
				Left:   2,
				Right:  2,
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
				archetype.PlayButtonClickAudio()
				browser.OpenURL(collab.Url)
			}),
			widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
				args.Button.TextColor = hovered
				aboutMenu.collabHovered = true
			}),
			widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
				args.Button.TextColor = normal
				aboutMenu.collabHovered = false
			}),
		)
		buttonsContainer.AddChild(collabButton)
		aboutMenu.collabButtons = append(aboutMenu.collabButtons, collabButton)
	}

	backContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(1),
		widget.GridLayoutOpts.Spacing(10, 3),
		widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
		widget.GridLayoutOpts.Padding(widget.Insets{
			Top: 10,
		}),
	)))

	aboutMenu.backButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(backLabel, assets.MainFont, &widget.ButtonTextColor{
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
			aboutMenu.Back = true
			archetype.PlayButtonClickAudio()
		}),
	)
	backContainer.AddChild(aboutMenu.backButton)

	parentContainer.AddChild(aboutLabelContainer)
	parentContainer.AddChild(buttonsContainer)
	parentContainer.AddChild(backContainer)

	aboutMenu.container.AddChild(parentContainer)
	parentContainer.GetWidget().LayoutData = widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	aboutMenu.Ui = &ebitenui.UI{
		Container: aboutMenu.container,
	}
	return aboutMenu
}

func (s *AboutMenuUI) Draw(screen *ebiten.Image) {
	s.Ui.Draw(screen)
}

func (s *AboutMenuUI) Container() *widget.Container {
	return s.container
}

func (s *AboutMenuUI) Update() {
	s.Ui.Update()
	backButton := s.backButton.GetWidget().Rect
	mx, my := ebiten.CursorPosition()
	if (backButton.Min.X <= mx && mx <= backButton.Max.X && backButton.Min.Y <= my && my <= backButton.Max.Y) ||
		s.collabHovered {
		archetype.UpdateCursorImage(true)
	} else {
		archetype.UpdateCursorImage(false)
	}
}
