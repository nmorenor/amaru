package archetype

import (
	"amaru/assets"
	"image"
	"image/color"

	uiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func DrawMainMenuRoundedRect(screen *ebiten.Image, x, y, width, height, borderRadius float64, fillColor, borderColor color.Color, borderWidth float64, label string) *ebiten.Image {
	// Create a gg.Context to draw the rounded rectangle and border
	img := image.NewRGBA(image.Rect(0, 0, int(width+2*borderWidth), int(height+2*borderWidth)))
	gc := gg.NewContextForRGBA(img)

	// Draw the rounded border
	gc.SetColor(borderColor)
	gc.DrawRoundedRectangle(borderWidth/2, borderWidth/2, width+borderWidth, height+borderWidth, borderRadius+borderWidth/2)
	gc.SetLineWidth(borderWidth)
	gc.Stroke()

	// Draw the rounded rectangle
	gc.SetColor(fillColor)
	gc.DrawRoundedRectangle(borderWidth, borderWidth, width, height, borderRadius)
	gc.Fill()
	menuRectangle := ebiten.NewImageFromImage(img)

	// Draw the img onto the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-borderWidth/2, y-borderWidth/2)
	if label != "" {
		gc.SetFontFace(assets.MainMidFont)
		gc.SetColor(borderColor)
		textWidth, textHeight := gc.MeasureString(label)
		textX := x + (width-textWidth)/2
		textY := y - textHeight // Adjust the textY position
		text.Draw(
			screen,
			label,
			assets.MainMidFont,
			int(textX),
			int(textY),
			borderColor,
		)
	}
	return menuRectangle
}

func CreateRoundedButtonImages(width, height, borderRadius float64, fillColor, borderColor, pressedBorderColor color.Color, hoverBorderColor color.Color, borderWidth float64) *widget.ButtonImage {
	idleImg := ebiten.NewImage(int(width), int(height))
	idleImg = DrawRoundedRect(idleImg, borderWidth, borderWidth, width-2*borderWidth, height-2*borderWidth, borderRadius, fillColor, borderColor, borderWidth)

	hoverImg := ebiten.NewImage(int(width), int(height))
	hoverImg = DrawRoundedRect(hoverImg, borderWidth, borderWidth, width-2*borderWidth, height-2*borderWidth, borderRadius, fillColor, hoverBorderColor, borderWidth)

	pressedImg := ebiten.NewImage(int(width), int(height))
	pressedImg = DrawRoundedRect(pressedImg, borderWidth, borderWidth, width-2*borderWidth, height-2*borderWidth, borderRadius, fillColor, pressedBorderColor, borderWidth)

	// Create the NineSlice images for idle and pressed states
	widths := [3]int{int(borderRadius), int(width) - 2*int(borderRadius), int(borderWidth)}
	heights := [3]int{int(borderRadius), int(height) - 2*int(borderRadius), int(borderWidth)}
	idleNineSlice := uiimage.NewNineSlice(idleImg, widths, heights)
	hoverNineSlice := uiimage.NewNineSlice(hoverImg, widths, heights)
	pressedNineSlice := uiimage.NewNineSlice(pressedImg, widths, heights)

	return &widget.ButtonImage{
		Idle:    idleNineSlice,
		Hover:   hoverNineSlice,
		Pressed: pressedNineSlice,
	}
}

func CreateRoundedTextInputImages(width, height, borderRadius float64, fillColor, borderColor, disabledBorderColor color.Color, borderWidth float64) *widget.TextInputImage {
	idleImg := ebiten.NewImage(int(width), int(height))
	idleImg = DrawRoundedRect(idleImg, borderWidth, borderWidth, width-2*borderWidth, height-2*borderWidth, borderRadius, fillColor, borderColor, borderWidth)

	disabledImg := ebiten.NewImage(int(width), int(height))
	disabledImg = DrawRoundedRect(disabledImg, borderWidth, borderWidth, width-2*borderWidth, height-2*borderWidth, borderRadius, fillColor, disabledBorderColor, borderWidth)

	// Create the NineSlice images for idle and pressed states
	widths := [3]int{int(borderRadius), int(width) - 2*int(borderRadius), int(borderWidth)}
	heights := [3]int{int(borderRadius), int(height) - 2*int(borderRadius), int(borderWidth)}
	idleNineSlice := uiimage.NewNineSlice(idleImg, widths, heights)
	hoverNineSlice := uiimage.NewNineSlice(disabledImg, widths, heights)

	return &widget.TextInputImage{
		Disabled: hoverNineSlice,
		Idle:     idleNineSlice,
	}
}

func DrawRoundedRect(screen *ebiten.Image, x, y, width, height, borderRadius float64, fillColor, borderColor color.Color, borderWidth float64) *ebiten.Image {
	// Create a gg.Context to draw the rounded rectangle and border
	img := image.NewRGBA(image.Rect(0, 0, int(width+2*borderWidth), int(height+2*borderWidth)))
	gc := gg.NewContextForRGBA(img)

	// Draw the rounded border
	gc.SetColor(borderColor)
	gc.DrawRoundedRectangle(borderWidth/2, borderWidth/2, width+borderWidth, height+borderWidth, borderRadius+borderWidth/2)
	gc.SetLineWidth(borderWidth)
	gc.Stroke()

	// Draw the rounded rectangle
	gc.SetColor(fillColor)
	gc.DrawRoundedRectangle(borderWidth, borderWidth, width, height, borderRadius)
	gc.Fill()

	// Draw a transparent inner rectangle
	gc.SetColor(color.Transparent)
	gc.DrawRectangle(borderWidth+borderRadius, borderWidth, width-2*borderRadius, height)
	gc.Fill()
	return ebiten.NewImageFromImage(img)
}

func UpdateCursorImage(isPointer bool) {
	if isPointer {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}

func NewTransparentEbitenImage(width, height int) *ebiten.Image {
	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill it with a transparent color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// Convert it to an ebiten.Image
	ebitenImage := ebiten.NewImageFromImage(img)

	return ebitenImage
}

func NewColoredEbitenImage(width, height int, color color.RGBA) *ebiten.Image {
	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill it with a transparent color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color)
		}
	}

	// Convert it to an ebiten.Image
	ebitenImage := ebiten.NewImageFromImage(img)

	return ebitenImage
}

func ColorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r / 257), // or r >> 8
		G: uint8(g / 257), // or g >> 8
		B: uint8(b / 257), // or b >> 8
		A: uint8(a / 257), // or a >> 8
	}
}

func CreateSVGImageButton(svgKey string, tooltipText string, targetInsets *widget.Insets, onClick func()) (*widget.Container, *widget.Button) {
	containerInsets := targetInsets
	if containerInsets == nil {
		containerInsets = &widget.Insets{
			Left:   15,
			Right:  15,
			Top:    5,
			Bottom: 5,
		}
	}
	tt := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(uiimage.NewNineSlice(NewToolTipImage(), [3]int{19, 6, 13}, [3]int{19, 5, 13})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(*containerInsets),
			widget.RowLayoutOpts.Spacing(2),
		)),
	)
	text := widget.NewText(
		widget.TextOpts.Text(tooltipText, assets.MainFont, assets.BlueColor),
	)
	tt.AddChild(text)
	image := assets.MustLoadImage(svgKey)

	// Create the button without text
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.ToolTip(widget.NewToolTip(
			widget.ToolTipOpts.Content(tt),
		))),
		widget.ButtonOpts.Image(loadButtonImage()),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if onClick != nil {
				onClick()
			}
		}),
		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),
	)
	imageWidget := widget.NewGraphic(widget.GraphicOpts.Image(image))

	// Bundle the button and the graphics widget together using a stacked layout container
	buttonStackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		// instruct the container's anchor layout to center the button both horizontally and vertically;
		// since our button is a 2-widget object, we add the anchor info to the wrapping container
		// instead of the button
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	)

	// Add button and image widget to the stacked layout
	buttonStackedLayout.AddChild(button)
	buttonStackedLayout.AddChild(imageWidget)

	return buttonStackedLayout, button
}

func loadButtonImage() *widget.ButtonImage {
	idle := uiimage.NewNineSliceColor(color.NRGBA{R: 255, G: 255, B: 255, A: 0})
	hover := uiimage.NewNineSliceColor(color.NRGBA{R: 255, G: 255, B: 255, A: 0})
	pressed := uiimage.NewNineSliceColor(color.NRGBA{R: 255, G: 255, B: 255, A: 0})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func NewToolTipImage() *ebiten.Image {
	width, height := 38, 37
	radius := float64(5)
	dc := gg.NewContext(width, height)

	dc.SetColor(color.White)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), radius)

	// Draw square corners on the bottom left, top right, and bottom right
	dc.DrawRectangle(0, float64(height)-radius, radius, radius)                     // Bottom Left
	dc.DrawRectangle(float64(width)-radius, 0, radius, radius)                      // Top Right
	dc.DrawRectangle(float64(width)-radius, float64(height)-radius, radius, radius) // Bottom Right

	dc.Fill()

	img := ebiten.NewImageFromImage(dc.Image())
	return img
}
