package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/yohamta/donburi/features/math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"amaru/engine"

	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

var (
	//go:embed fonts/ark-pixel-16px-monospaced-latin.ttf
	mainFontData []byte

	//go:embed images/close-circle.svg
	closeCircleData []byte

	//go:embed images/send.svg
	sendData []byte

	//go:embed images/audio.svg
	audioData []byte

	//go:embed meta/about.json
	AboutData []byte

	//go:embed *
	assetsFS embed.FS

	audioContext *audio.Context

	// AvailableLevels []Level
	GameLevelLoader *LevelLoader

	MessagesFont      font.Face
	MessagesSmallFont font.Face
	MessagesBigFont   font.Face

	MainFont      font.Face
	MainMidFont   font.Face
	MainSmallFont font.Face
	MainBigFont   font.Face

	BoatSpriteSheet   *engine.Spritesheet
	WasteSpriteSheet  *engine.Spritesheet
	TurtleSpriteSheet *engine.Spritesheet
	SealSpriteSheet   *engine.Spritesheet
	DirectionalPad    *ebiten.Image
	DirectionalBtn    *ebiten.Image

	BlueColor  color.Color
	GreenColor color.Color

	BlueColorHex  = "#27bdf5cc"
	GreenColorHex = "#43ff64d9"

	images   = map[string]*ebiten.Image{}
	CloseKey = "close-circle"
	AudioKey = "audio"
	SendKey  = "send"

	MenuAdioPlayer       *audio.Player
	GameAudioPlayer      *audio.Player
	WavesAudioPlayer     *audio.Player
	HeronAudioPlayer     *audio.Player
	ShipAudioPlayer      *audio.Player
	CollectedAudioPlayer *audio.Player
	ButtonClickPlayer    *audio.Player
)

type AudioStream interface {
	io.ReadSeeker
	Length() int64
}

type Level struct {
	OverPlayer   *ebiten.Image
	Background   *ebiten.Image
	Paths        map[uint32]Path
	Animals      []Path
	PlayersStart []Path
}

type Path struct {
	Points []math.Vec2
	Loops  bool
}

func (p *Path) TetraCenter() math.Vec2 {
	if len(p.Points) != 4 {
		// This method assumes there are exactly 4 points
		return math.Vec2{}
	}

	var centerX, centerY float64
	for _, point := range p.Points {
		centerX += point.X
		centerY += point.Y
	}

	centerX /= 4
	centerY /= 4

	return math.Vec2{X: centerX, Y: centerY}
}

type LevelLoader struct {
	LevelsSize        int
	CurrentLevel      *Level
	CurrentLevelIndex int
}

func MustLoadAssets() {
	GameLevelLoader = newLevelLoader()
	BlueColor, _ = colorful.Hex(BlueColorHex)
	GreenColor, _ = colorful.Hex(GreenColorHex)

	MainFont = mustLoadFont(mainFontData, 16)
	MainMidFont = mustLoadFont(mainFontData, 32)
	MainSmallFont = mustLoadFont(mainFontData, 8)
	MainBigFont = mustLoadFont(mainFontData, 48)

	BoatSpriteSheet = mustLoadNorthSpriteSheet("images/boat.png", 32, 32)
	WasteSpriteSheet = mustLoadNorthSpriteSheet("images/waste-sheet.png", 32, 32)
	TurtleSpriteSheet = mustLoadNorthSpriteSheet("images/turtle-sheet.png", 32, 32)
	SealSpriteSheet = mustLoadNorthSpriteSheet("images/foca-sheet.png", 32, 32)
	DirectionalPad = MustLoadImageFromFS("images/directional_pad.png")
	DirectionalBtn = MustLoadImageFromFS("images/directional_button.png")

	audioContext = audio.NewContext(44100)
	MenuAdioPlayer = mustLoadAudioPlayer("audio/menu.ogg", audioContext)
	GameAudioPlayer = mustLoadAudioPlayer("audio/main-game.ogg", audioContext)
	WavesAudioPlayer = mustLoadAudioPlayer("audio/waves.ogg", audioContext)
	HeronAudioPlayer = mustLoadAudioPlayer("audio/heron.ogg", audioContext)
	ShipAudioPlayer = mustLoadAudioPlayer("audio/ship.ogg", audioContext)
	CollectedAudioPlayer = mustLoadAudioPlayer("audio/collected.ogg", audioContext)
	ButtonClickPlayer = mustLoadAudioPlayer("audio/button-click.ogg", audioContext)
}

func newLevelLoader() *LevelLoader {
	levelPaths, err := fs.Glob(assetsFS, "levels/level*.tmx")
	if err != nil {
		panic(err)
	}
	return &LevelLoader{
		LevelsSize: len(levelPaths),
	}
}

func mustLoadNorthSpriteSheet(image string, width int, height int) *engine.Spritesheet {
	img := MustLoadImageFromFS(image)

	return engine.NewSpritesheetFromTexture(img, width, height)
}

func mustLoadFont(data []byte, size int) font.Face {
	f, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	return face
}

func mustNewEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func (l *LevelLoader) LoadLevel(index int) *Level {
	if index == l.CurrentLevelIndex && l.CurrentLevel != nil {
		return l.CurrentLevel
	}
	levelPaths, err := fs.Glob(assetsFS, "levels/level*.tmx")
	if err != nil {
		panic(err)
	}
	targetPath := levelPaths[index]

	l.CurrentLevelIndex = index
	l.CurrentLevel = engine.Ptr(l.MustLoadLevel(targetPath))

	return l.CurrentLevel
}

func (l *LevelLoader) MustLoadLevel(levelPath string) Level {
	levelMap, err := tiled.LoadFile(levelPath, tiled.WithFileSystem(assetsFS))
	if err != nil {
		panic(err)
	}

	nextLevel := Level{}

	paths := map[uint32]Path{}
	animals := []Path{}
	playerStarts := []Path{}
	for _, og := range levelMap.ObjectGroups {
		for _, o := range og.Objects {
			if o.Width != 0 && o.Height != 0 && len(o.PolyLines) == 0 && len(o.Polygons) == 0 {
				box := l.MustLoadBox(o)
				if o.Class == "playerStart" {
					playerStarts = append(playerStarts, box)
				} else if o.Class == "animal" {
					animals = append(animals, box)
				} else {
					paths[o.ID] = box
				}
			}
			if len(o.PolyLines) > 0 {
				var points []math.Vec2
				for _, p := range o.PolyLines {
					for _, pp := range *p.Points {
						points = append(points, math.Vec2{
							X: o.X + pp.X,
							Y: o.Y + pp.Y,
						})
					}
				}
				paths[o.ID] = Path{
					Loops:  false,
					Points: points,
				}
			}
			if len(o.Polygons) > 0 {
				var points []math.Vec2
				for _, p := range o.Polygons {
					for _, pp := range *p.Points {
						points = append(points, math.Vec2{
							X: o.X + pp.X,
							Y: o.Y + pp.Y,
						})
					}
				}
				paths[o.ID] = Path{
					Loops:  true,
					Points: points,
				}
			}
		}
	}

	renderer, err := render.NewRendererWithFileSystem(levelMap, assetsFS)
	if err != nil {
		panic(err)
	}
	overPlayerRenderer, err := render.NewRendererWithFileSystem(levelMap, assetsFS)
	if err != nil {
		panic(err)
	}

	// err = renderer.RenderVisibleLayers()
	for i := range levelMap.Layers {
		if levelMap.Layers[i].Class == "overPlayer" {
			err = overPlayerRenderer.RenderLayer(i)
			if err != nil {
				panic(err)
			}
			continue
		}
		err = renderer.RenderLayer(i)
		if err != nil {
			panic(err)
		}
	}

	nextLevel.Background = ebiten.NewImageFromImage(renderer.Result)
	nextLevel.OverPlayer = ebiten.NewImageFromImage(overPlayerRenderer.Result)
	nextLevel.Paths = paths
	nextLevel.Animals = animals
	nextLevel.PlayersStart = playerStarts

	return nextLevel
}

func (l *LevelLoader) MustLoadBox(o *tiled.Object) Path {
	y := o.Y - 32
	points := []math.Vec2{
		{X: o.X, Y: y},
		{X: o.X + o.Width, Y: y},
		{X: o.X + o.Width, Y: y + o.Height},
		{X: o.X, Y: y + o.Height},
	}

	return Path{
		Points: points,
		Loops:  true,
	}

}

func MustLoadImageFromFS(filePath string) *ebiten.Image {
	file, err := assetsFS.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file from embed.FS %s", filePath))
	}

	return mustNewEbitenImage(file)
}

func MustLoadSvgs() {
	closedCircle, err := SvgToEbitenImage(closeCircleData, 32, 32)
	if err != nil {
		panic("Failed to read file from embed.FS close-circle.svg")
	}
	audioImage, err := SvgToEbitenImage(audioData, 32, 32)
	if err != nil {
		panic("Failed to read file from embed.FS audio.svg")
	}
	sendImage, err := SvgToEbitenImage(sendData, 32, 32)
	if err != nil {
		panic("Failed to read file from embed.FS send.svg")
	}
	images[CloseKey] = closedCircle
	images[AudioKey] = audioImage
	images[SendKey] = sendImage
}

func MustLoadImage(assetKey string) *ebiten.Image {
	return images[assetKey]
}

func SvgToEbitenImage(svgData []byte, targetWidth, targetHeight int) (*ebiten.Image, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData), oksvg.IgnoreErrorMode)
	if err != nil {
		return nil, err
	}

	// Get the original dimensions of the SVG
	originalWidth := icon.ViewBox.W
	originalHeight := icon.ViewBox.H

	// Calculate the scale factors
	scaleX := float64(targetWidth) / originalWidth
	scaleY := float64(targetHeight) / originalHeight

	// Set the target dimensions for the SVG icon
	icon.SetTarget(0, 0, originalWidth*scaleX, originalHeight*scaleY)

	img := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	scannerGV := rasterx.NewScannerGV(targetWidth, targetHeight, img, img.Bounds())
	raster := rasterx.NewDasher(targetWidth, targetHeight, scannerGV)

	icon.Draw(raster, 1.0)

	ebitenImage := ebiten.NewImageFromImage(img)

	return ebitenImage, nil
}

func mustLoadAudioPlayer(filePath string, audioContext *audio.Context) *audio.Player {
	file, err := assetsFS.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file from embed.FS %s", filePath))
	}
	var stream AudioStream
	stream, err = vorbis.DecodeWithoutResampling(bytes.NewReader(file))
	if err != nil {
		panic(fmt.Sprintf("Failed to read audio file from embed.FS %s", filePath))
	}
	audioPlayer, err := audioContext.NewPlayer(stream)
	if err != nil {
		panic(fmt.Sprintf("Failed to create audio player for %s", filePath))
	}
	return audioPlayer
}
