package ui

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/engine"
	"time"

	"golang.org/x/image/colornames"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"
)

type TextInputMenu struct {
	container         *widget.Container
	Ui                *ebitenui.UI
	okButton          *widget.Button
	okButtonLabel     string
	cancelButton      *widget.Button
	pasteButtonWidget *widget.Button
	cancelButtonLabel string
	Change            float32
	canSubmit         bool

	Value     *string
	Done      bool
	Cancel    bool
	inputText *widget.TextInput
}

func NewTextInputMenu(title string, label string, cancelLabel string, placeHolder string) *TextInputMenu {
	textInputMenu := &TextInputMenu{
		okButtonLabel:     label,
		cancelButtonLabel: cancelLabel,
		container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
				widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
			)),
		),
	}

	parentContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			widget.GridLayoutOpts.Spacing(25, 5),
		)),
	)

	welcomeLabel := widget.NewLabel(widget.LabelOpts.Text(title, assets.MainFont, &widget.LabelColor{
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

	textInputMenu.inputText = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.TextInputOpts.Image(archetype.CreateRoundedTextInputImages(200, 50, 5, colornames.White, colornames.Fuchsia, colornames.Grey, 5)),
		widget.TextInputOpts.Placeholder(placeHolder),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          assets.BlueColor,
			Disabled:      colornames.Grey,
			Caret:         colornames.Gray,
			DisabledCaret: colornames.Grey,
		}),
		widget.TextInputOpts.Face(assets.MainFont),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Color(colornames.Gray),
			widget.CaretOpts.Size(assets.MainFont, 16),
		),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			textInputMenu.Value = engine.Ptr(args.TextInput.GetText())
		}),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			if textInputMenu.canSubmit && len(args.TextInput.GetText()) > 0 {
				textInputMenu.Value = engine.Ptr(args.TextInput.GetText())
				textInputMenu.Done = true
			}
		}),
		widget.TextInputOpts.Padding(widget.Insets{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}),
		widget.TextInputOpts.RepeatInterval(150*time.Millisecond),
	)

	// Create a new Container with AnchorLayout
	userNameInputContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Spacing(0),
		widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(0)))))

	// Add the TextInput instance to the userNameInputContainer
	userNameInputContainer.AddChild(textInputMenu.inputText)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(10, 3),
			),
		),
	)
	buttonsContainer.GetWidget().LayoutData = widget.WidgetOpts.LayoutData(widget.RowLayoutData{
		Position: widget.RowLayoutPositionCenter,
		Stretch:  true,
	})

	textInputMenu.okButton = widget.NewButton(
		widget.ButtonOpts.Image(archetype.CreateRoundedButtonImages(200, 50, 5, colornames.White, assets.BlueColor, assets.BlueColor, assets.GreenColor, 5)),
		widget.ButtonOpts.Text(textInputMenu.okButtonLabel, assets.MainFont, &widget.ButtonTextColor{
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
			if textInputMenu.canSubmit && len(textInputMenu.inputText.GetText()) > 0 {
				textInputMenu.Done = true
			}
		}),
	)

	textInputMenu.cancelButton = widget.NewButton(
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
			textInputMenu.Cancel = true
		}),
	)

	buttonsContainer.AddChild(textInputMenu.okButton)
	buttonsContainer.AddChild(textInputMenu.cancelButton)

	parentContainer.AddChild(welcomeLabelContainer)
	parentContainer.AddChild(userNameInputContainer)
	parentContainer.AddChild(buttonsContainer)

	textInputMenu.container.AddChild(parentContainer)
	parentContainer.GetWidget().LayoutData = widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	textInputMenu.Ui = &ebitenui.UI{
		Container: textInputMenu.container,
	}
	textInputMenu.inputText.Focus(true)
	return textInputMenu
}

func (s *TextInputMenu) Draw(screen *ebiten.Image) {
	s.Ui.Draw(screen)
}

func (s *TextInputMenu) Container() *widget.Container {
	return s.container
}

func (s *TextInputMenu) Update() {
	s.Ui.Update()
	s.canSubmit = true
	okButtonRect := s.okButton.GetWidget().Rect
	cancelButtonRect := s.cancelButton.GetWidget().Rect
	mx, my := ebiten.CursorPosition()
	if (okButtonRect.Min.X <= mx && mx <= okButtonRect.Max.X && okButtonRect.Min.Y <= my && my <= okButtonRect.Max.Y) ||
		(cancelButtonRect.Min.X <= mx && mx <= cancelButtonRect.Max.X && cancelButtonRect.Min.Y <= my && my <= cancelButtonRect.Max.Y) ||
		(s.pasteButtonWidget != nil && s.pasteButtonWidget.GetWidget().Rect.Min.X <= mx && mx <= s.pasteButtonWidget.GetWidget().Rect.Max.X && s.pasteButtonWidget.GetWidget().Rect.Min.Y <= my && my <= s.pasteButtonWidget.GetWidget().Rect.Max.Y) {
		archetype.UpdateCursorImage(true)
	} else {
		archetype.UpdateCursorImage(false)
	}

	ctrlPressed := ebiten.IsKeyPressed(ebiten.KeyControl)
	leftMetaPressed := ebiten.IsKeyPressed(ebiten.KeyMetaLeft)
	rightMetaPressed := ebiten.IsKeyPressed(ebiten.KeyMetaRight)
	vPressed := ebiten.IsKeyPressed(ebiten.KeyV)

	if (ctrlPressed || leftMetaPressed || rightMetaPressed) && vPressed {
		cvalue, err := engine.ReadClipboard()
		if err == nil {
			value := string(cvalue[:])
			s.inputText.SetText(value)
			s.inputText.CursorMoveEnd()
		}
	}
}
