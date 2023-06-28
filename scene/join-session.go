package scene

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"amaru/system"
	"amaru/ui"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"golang.org/x/image/colornames"
)

type JoinSessionMenu struct {
	world     *donburi.World
	game      *component.GameData
	systems   []System
	drawables []Drawable

	screenWidth  int
	screenHeight int

	offscreen *ebiten.Image
	uiHandler *ui.AvailableSessionsMenu
}

func NewJoinSessionMenu(screenWidth int, screenHeight int, session *component.SessionData, gameData *component.GameData) *JoinSessionMenu {
	assets.MustLoadSvgs()
	menu := &JoinSessionMenu{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		uiHandler:    ui.NewAvailableSessionsMenu(gameData),
	}

	menu.loadMenu(session)

	return menu
}

func (menu *JoinSessionMenu) loadMenu(session *component.SessionData) {
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

	menu.world = engine.Ptr(menu.createWorld())
	menu.game = component.MustFindGame(*menu.world)
	menu.game.Session = session
	uiRender.Initialize(*menu.world)
}

func (menu *JoinSessionMenu) UpdateLayout(width, height int) {
	// do nothing
}

func (menu *JoinSessionMenu) createWorld() donburi.World {
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

	archetype.PlayAudioMenu()

	return world
}

func (menu *JoinSessionMenu) renderUI(image *ebiten.Image) *ebiten.Image {
	menu.uiHandler.Draw(image)
	return image
}

func (menu *JoinSessionMenu) NextScene() archetype.Scene {
	if menu.uiHandler.Done && menu.uiHandler.Session != nil {
		menu.game.Session.SessionID = &menu.uiHandler.Session.ID
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		menu.uiHandler.Ui.Container.RemoveChildren()
		menu.uiHandler.Ui = nil
		return NewConnectingMenuMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight, menu.game.Session)
	}
	if menu.uiHandler.Cancelled {
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		menu.uiHandler.Ui.Container.RemoveChildren()
		menu.uiHandler.Ui = nil
		return NewStartMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight)
	}
	return menu
}

func (menu *JoinSessionMenu) Update() {
	archetype.PlayAudioMenu()
	menu.uiHandler.Update()
	for _, s := range menu.systems {
		s.Update(*menu.world)
	}

}

func (menu *JoinSessionMenu) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range menu.drawables {
		s.Draw(*menu.world, screen)
	}
}
