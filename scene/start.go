package scene

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"amaru/system"
	"amaru/ui"
	"time"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const (
	borderWidth = 7
	menuTitle   = "Amaru!, Reverse Pollution!"
)

type StartMenu struct {
	world     *donburi.World
	game      *component.GameData
	systems   []System
	drawables []Drawable

	screenWidth  int
	screenHeight int

	offscreen *ebiten.Image
	uiHandler *ui.StartMenu
}

func NewStartMenu(screenWidth int, screenHeight int) *StartMenu {
	menu := &StartMenu{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		uiHandler:    ui.NewStartMenu(),
	}

	menu.loadMenu()

	return menu
}

func (menu *StartMenu) loadMenu() {
	selectedLevelIndex := engine.RandomIntRange(0, assets.GameLevelLoader.LevelsSize)
	assets.GameLevelLoader.LoadLevel(selectedLevelIndex)
	render := system.NewRenderer()
	uiRender := system.NewUIRenderer()

	menu.systems = []System{
		system.NewCamera(),
		render,
		uiRender,
	}

	menu.drawables = []Drawable{
		render,
		uiRender,
	}

	menu.world = engine.Ptr(menu.createWorld(selectedLevelIndex))
	menu.game = component.MustFindGame(*menu.world)
	menu.game.Session = nil // reset session
	uiRender.Initialize(*menu.world)
}

func (menu *StartMenu) UpdateLayout(width, height int) {
	// do nothing
}

func (menu *StartMenu) createWorld(selectedLevelIndex int) donburi.World {
	rectX := float64(menu.screenWidth/2) - (float64(menu.screenWidth/2) / 2)
	rectY := float64(menu.screenHeight/2) - (float64(menu.screenHeight/2) / 2)

	menuContainerImage := archetype.DrawMainMenuRoundedRect(menu.offscreen, rectX, rectY, float64(menu.screenWidth/2), float64(menu.screenHeight/2), 5, colornames.White, assets.BlueColor, borderWidth, menuTitle)
	world := donburi.NewWorld()

	archetype.NewInput(world)

	level := world.Entry(world.Create(component.Level))
	component.Level.Get(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	cameraEntry := archetype.NewCamera(world, menu.screenWidth, menu.screenHeight, math.Vec2{
		X: 0,
		Y: 0,
	})

	selectedLevel := assets.GameLevelLoader.CurrentLevel

	component.Camera.Get(cameraEntry).Disabled = true

	levelEntry := world.Entry(
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(levelEntry, component.SpriteData{
		Image: selectedLevel.Background,
		Layer: component.SpriteLayerBackground,
		Pivot: component.SpritePivotScreenCenter,
	})
	overPlayerEntry := world.Entry(
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(overPlayerEntry, component.SpriteData{
		Image: selectedLevel.OverPlayer,
		Layer: component.SpriteLayerForeground,
		Pivot: component.SpritePivotScreenCenter,
	})
	menuEntry := world.Entry(
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(menuEntry, component.SpriteData{
		Image: menu.offscreen,
		Layer: component.SpriteLayerUI,
		Pivot: component.SpritePivotTopLeft,
	})

	menuUIEntry := world.Entry(
		world.Create(transform.Transform, component.UISprite),
	)
	component.UISprite.SetValue(menuUIEntry, component.UISpriteData{
		Image:     menuContainerImage,
		Layer:     component.SpriteLayerUI,
		Pivot:     component.SpritePivotScreenCenter,
		UIHandler: menu.renderUI,
	})

	if menu.world == nil {
		game := world.Entry(world.Create(component.Game))
		component.Game.SetValue(game, component.GameData{
			Settings: component.Settings{
				ScreenWidth:  menu.screenWidth,
				ScreenHeight: menu.screenHeight,
			},
			Speed:      3.0,
			LeftOffset: 0,
		})
	}

	archetype.StopAudioGame()
	archetype.PlayAudioMenu()

	return world
}

func (menu *StartMenu) renderUI(image *ebiten.Image) *ebiten.Image {
	menu.uiHandler.Draw(image)
	return image
}

func (menu *StartMenu) NextScene() archetype.Scene {

	if menu.uiHandler.SelectedOption == ui.About {
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		menu.uiHandler.Ui.Container.RemoveChildren()
		menu.uiHandler.Ui = nil
		return NewAboutMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight)
	}

	if menu.uiHandler.SelectedOption == ui.Host {
		menu.game.Session = &component.SessionData{
			Type:          component.SessionTypeHost,
			PlayerMessage: make(map[string]*component.RemotePlayerMessage),
		}
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		menu.uiHandler.Ui.Container.RemoveChildren()
		menu.uiHandler.Ui = nil
		return NewHostUserNameMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight, menu.game.Session)
	}
	if menu.uiHandler.SelectedOption == ui.Join {
		menu.game.Session = &component.SessionData{
			Type:          component.SessionTypeJoin,
			PlayerMessage: make(map[string]*component.RemotePlayerMessage),
		}
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		return NewJoinUserNameMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight, menu.game.Session)
	}
	return menu
}

func (menu *StartMenu) Update() {
	archetype.PlayAudioMenu()
	for _, s := range menu.systems {
		s.Update(*menu.world)
	}
	menu.uiHandler.Update()
}

func (menu *StartMenu) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range menu.drawables {
		s.Draw(*menu.world, screen)
	}
}
