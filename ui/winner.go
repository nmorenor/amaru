package ui

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"fmt"
	"math"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"golang.org/x/image/colornames"
)

type WinnerUI struct {
	container          *widget.Container
	Ui                 *ebitenui.UI
	sendButton         *widget.Button
	inputText          *widget.TextInput
	Game               *component.GameData
	World              *donburi.World
	MessageValue       *string
	MessageDone        bool
	textAreaLayoutData widget.RowLayoutData
	textArea           *widget.TextArea
	remainingTimeLabel *widget.Label
}

func NewWinnerUI(winner string, gameData *component.GameData) *WinnerUI {
	winnerUI := &WinnerUI{
		container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		),
		Game: gameData,
	}

	remainingContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{
			Top: 5,
		}))),
	)
	winnerUI.remainingTimeLabel = widget.NewLabel(
		widget.LabelOpts.Text("00", assets.MainBigFont, &widget.LabelColor{
			Disabled: colornames.White,
			Idle:     colornames.White,
		}),
	)
	winnerUI.remainingTimeLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionStart,
	}
	remainingContainer.AddChild(winnerUI.remainingTimeLabel)
	winnerUI.container.AddChild(remainingContainer)

	centerContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	parentContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.GridLayoutOpts.Spacing(25, 5),
		)),
	)
	centerContainer.AddChild(parentContainer)

	winnerLabel := widget.NewLabel(widget.LabelOpts.Text(fmt.Sprintf("Winner: %s", winner), assets.MainMidFont, &widget.LabelColor{
		Disabled: assets.BlueColor,
		Idle:     assets.BlueColor,
	}))

	winnerLabelContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(0)),
	)))

	winnerLabelContainer.AddChild(winnerLabel)
	winnerLabel.GetWidget().LayoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	chatContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(1),
		widget.GridLayoutOpts.Spacing(0, 2),
	)))

	textAreaContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   false,
			MaxWidth:  (gameData.Settings.ScreenHeight / 2) - 24,
			MaxHeight: (gameData.Settings.ScreenHeight / 2) - 16,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(0),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   8,
				Top:    8,
				Bottom: 8,
			}),
		)),
	)
	chatContainer.AddChild(textAreaContainer)
	winnerUI.textAreaLayoutData = widget.RowLayoutData{
		Stretch:   false,
		MaxWidth:  (gameData.Settings.ScreenHeight / 2) - 24,
		MaxHeight: (gameData.Settings.ScreenHeight / 2) - 16,
	}
	winnerUI.textArea = newTextArea("")
	winnerUI.textArea.GetWidget().LayoutData = winnerUI.textAreaLayoutData
	textAreaContainer.AddChild(winnerUI.textArea)

	inputContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(0),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   10,
				Bottom: 0,
			}),
		)),
	)
	chatContainer.AddChild(inputContainer)

	winnerUI.inputText = widget.NewTextInput(
		widget.TextInputOpts.Image(archetype.CreateRoundedTextInputImages(200, 50, 5, colornames.White, colornames.Fuchsia, colornames.Grey, 5)),
		widget.TextInputOpts.Placeholder(sendPlaceHolder),
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
			winnerUI.MessageValue = engine.Ptr(args.TextInput.GetText())
		}),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			if len(args.TextInput.GetText()) > 0 {
				winnerUI.MessageValue = engine.Ptr(args.TextInput.GetText())
				winnerUI.MessageDone = true
				winnerUI.UpdateTextArea(fmt.Sprintf("\n[color=27BDF5]%s:[/color]", *winnerUI.Game.Session.UserName))
				winnerUI.UpdateTextArea(args.TextInput.GetText())
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
	winnerUI.container.GetWidget().LayoutData = widget.WidgetOpts.LayoutData(widget.RowLayoutData{
		MaxHeight: 64,
		Position:  widget.RowLayoutPositionStart,
		Stretch:   true,
	})
	inputContainer.AddChild(winnerUI.inputText)

	sendButtonInsets := &widget.Insets{
		Left:   15,
		Right:  15,
		Top:    0,
		Bottom: 5,
	}
	sendButtonContainer, sendButton := archetype.CreateSVGImageButton(assets.SendKey, "Send", sendButtonInsets, func() {
		if winnerUI.MessageValue != nil && len(*winnerUI.MessageValue) > 0 {
			winnerUI.MessageDone = true
			archetype.PlayButtonClickAudio()
		}
	})
	sendButtonContainer.GetWidget().LayoutData = widget.RowLayoutData{
		Position: widget.RowLayoutPositionEnd,
		Stretch:  false,
	}

	winnerUI.sendButton = sendButton
	inputContainer.AddChild(sendButtonContainer)

	parentContainer.AddChild(winnerLabelContainer)
	parentContainer.AddChild(chatContainer)

	winnerUI.container.AddChild(centerContainer)
	parentContainer.GetWidget().LayoutData = widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
	}

	winnerUI.Ui = &ebitenui.UI{
		Container: winnerUI.container,
	}

	winnerUI.inputText.Focus(true)

	return winnerUI
}

func (s *WinnerUI) UpdateTextArea(text string) {
	nextText := s.textArea.GetText() + "\n" + text
	s.textArea.SetText(nextText)
}

func (s *WinnerUI) Reset() {
	s.inputText.SetText("")
	s.MessageDone = false
	s.MessageValue = nil
}

func (s *WinnerUI) Draw(screen *ebiten.Image) {
	s.Ui.Draw(screen)
}

func (s *WinnerUI) Container() *widget.Container {
	return s.container
}

func (s *WinnerUI) Update() {
	if s.MessageDone {
		s.Game.Session.RemoteClient.SendChatMessage(*s.MessageValue)
		s.Reset()
	}

	s.Game.Session.RemoteClient.GameData.Frames++
	if int(ebiten.ActualTPS()) > 0 && s.Game.Session.RemoteClient.GameData.Frames%(int(ebiten.ActualTPS())) == 0 {
		s.Game.Session.RemoteClient.GameData.Counter--
	}

	s.remainingTimeLabel.Label = fmt.Sprintf("%02d", int(math.Max(float64(s.Game.Session.RemoteClient.GameData.Counter), 0)))

	for s.Game.ChatMessages.Length() > 0 {
		message := s.Game.ChatMessages.Remove()
		from := s.Game.Session.RemoteClient.Participants[message.Source]
		if from != nil {
			s.UpdateTextArea(fmt.Sprintf("\n[color=43FF64]%s:[/color]", *from))
			s.UpdateTextArea(message.Message)
		}
	}

	textAreaWidget := s.textArea.GetWidget()
	textAreaWidget.MinWidth = (s.Game.Settings.ScreenHeight / 2) + 64
	textAreaWidget.MinHeight = (s.Game.Settings.ScreenHeight / 2) - 128
	s.textAreaLayoutData.MaxWidth = (s.Game.Settings.ScreenHeight / 2) + 64
	s.textAreaLayoutData.MaxHeight = (s.Game.Settings.ScreenHeight / 2) - 128
	s.inputText.GetWidget().MinWidth = (s.Game.Settings.ScreenHeight / 2) + 32

	textAreaWidget.LayoutData = s.textAreaLayoutData

	if s.Game.Session.Type == component.SessionTypeHost && s.Game.Session.RemoteClient.GameData.Counter <= 0 {
		s.Game.GameOver = true
	}
	sendButtonRect := s.sendButton.GetWidget().Rect
	mx, my := ebiten.CursorPosition()
	if sendButtonRect.Min.X <= mx && mx <= sendButtonRect.Max.X && sendButtonRect.Min.Y <= my && my <= sendButtonRect.Max.Y {
		archetype.UpdateCursorImage(true)
	} else {
		archetype.UpdateCursorImage(false)
	}
	s.Ui.Update()
}
