package scene

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"amaru/net"
	"amaru/system"
	"amaru/ui"
	"context"
	"time"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

type WinnerMenu struct {
	world     *donburi.World
	game      *component.GameData
	systems   []System
	drawables []Drawable

	screenWidth  int
	screenHeight int

	offscreen *ebiten.Image
	uiHandler *ui.WinnerUI
}

func NewWinnerMenu(gameData *component.GameData) *WinnerMenu {
	gameData.GameOver = false
	var winnerParticipant *net.SessionParticipant
	maxPoints := 0
	for _, participant := range gameData.Session.RemoteClient.GameData.SessionParticipants {
		if participant.Score > maxPoints {
			winnerParticipant = participant
			maxPoints = participant.Score
		}
	}
	winner := "No winner"
	if winnerParticipant != nil {
		winner = *winnerParticipant.Name
	}

	menu := &WinnerMenu{
		screenWidth:  gameData.Settings.ScreenWidth,
		screenHeight: gameData.Settings.ScreenHeight,
		offscreen:    ebiten.NewImage(gameData.Settings.ScreenWidth, gameData.Settings.ScreenHeight),
		uiHandler:    ui.NewWinnerUI(winner, gameData),
	}

	archetype.StopAudioMenu()

	menu.loadMenu(gameData)

	return menu
}

func (menu *WinnerMenu) loadMenu(gameData *component.GameData) {
	lastIndex := gameData.Session.RemoteClient.GameData.LevelIndex
	selectedLevelIndex := engine.RandomIntRange(0, assets.GameLevelLoader.LevelsSize)
	for selectedLevelIndex == lastIndex {
		selectedLevelIndex = engine.RandomIntRange(0, assets.GameLevelLoader.LevelsSize)
	}
	assets.GameLevelLoader.LoadLevel(selectedLevelIndex)
	gameData.Session.RemoteClient.GameData.LevelIndex = selectedLevelIndex
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

	gameData.Session.RemoteClient.RemoteGameData.AddListener(func(ctx context.Context, rwm net.RemoteGameDataMessage) {
		gameData.Session.RemoteClient.GameData = rwm.Msg
	})

	menu.world = engine.Ptr(menu.createWorld(gameData))
	menu.game = gameData

	uiRender.Initialize(*menu.world)
}

func (menu *WinnerMenu) UpdateLayout(width, height int) {
	// do nothing
}

func (menu *WinnerMenu) createWorld(gameData *component.GameData) donburi.World {
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
	space, shapes := archetype.SetupSpaceForLevel(selectedLevel)

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
	// host creates waste then send location to remote players
	// animals are set on level state
	wasteSize := engine.RandomIntRange(component.MinWasteSize, component.MaxWasteSize)

	debugEntity := world.Create(component.Debug)
	debugComponent := component.Debug.Get(world.Entry(debugEntity))
	debugComponent.Shapes = shapes

	physics := world.Entry(world.Create(component.Physics))
	component.Physics.Get(physics).Space = space

	boxes := lo.Map(shapes, func(shape *cp.Shape, idx int) cp.BB {
		return shape.BB()
	})

	animalBoxes := lo.Map(selectedLevel.Animals, func(animal assets.Path, idx int) cp.BB {
		return archetype.CreateBoxFromPath(space, animal, component.AnimalCollisionType).BB()
	})
	boxes = append(boxes, animalBoxes...)
	wasteList := archetype.PlaceWasteComponents(world, space, wasteSize, debugComponent, boxes, float64(selectedLevel.Background.Bounds().Dx()), float64(selectedLevel.Background.Bounds().Dy()))
	wastePaths := lo.Map(wasteList, func(waste *component.WasteData, idx int) *net.WasteLocation {
		return &net.WasteLocation{
			Id:        waste.Id,
			Location:  waste.Path,
			Collected: false,
		}
	})
	locations := map[string]*net.WasteLocation{}
	for _, location := range wastePaths {
		locations[location.Id] = location
	}

	if menu.world == nil {
		game := world.Entry(world.Create(component.Game))
		component.Game.SetValue(game, *gameData)
	}

	gameData.Session.RemoteClient.RemoteChat.AddListener(func(ctx context.Context, rcm net.RemoteChatMessage) {
		gameData.ChatMessages.Add(&rcm.Msg)
	})

	gameData.Session.RemoteClient.RemoteGameData.AddListener(func(ctx context.Context, rcm net.RemoteGameDataMessage) {
		gameData.Session.RemoteClient.GameData = rcm.Msg
		if gameData.Session.RemoteClient.GameData.OnGameState {
			gameData.GameOver = true
		}
	})

	gameData.Session.RemoteClient.SessionEnd.AddListener(func(ctx context.Context, val int) {
		gameData.Session.End = true
	})

	if gameData.Session.Type == component.SessionTypeHost {
		gameData.Session.RemoteClient.GameData.Counter = 15
		gameData.Session.RemoteClient.GameData.Frames = 0
		gameData.Session.RemoteClient.GameData.OnGameState = false
		gameData.Session.RemoteClient.GameData.WasteLocations = locations
		go gameData.Session.RemoteClient.SendGameDataMessage(*gameData.Session.RemoteClient.GameData)
	} else {
		go gameData.Session.RemoteClient.RequestGameData()
	}

	archetype.StopAudioMenu()
	archetype.StopAudioGame()
	archetype.PlayWinnerAudio(gameData)

	return world
}

func (menu *WinnerMenu) renderUI(image *ebiten.Image) *ebiten.Image {
	menu.uiHandler.Draw(image)
	return image
}

func (menu *WinnerMenu) NextScene() archetype.Scene {

	if menu.game.GameOver {
		menu.game.Session.RemoteClient.ResetListeners()
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		menu.uiHandler.Ui.Container.RemoveChildren()
		menu.uiHandler.Ui = nil

		if menu.game.Session.Type == component.SessionTypeHost && menu.game.Session.RemoteClient.GameData != nil {
			menu.game.Session.RemoteClient.GameData.Counter = 30
			menu.game.Session.RemoteClient.GameData.Frames = 0
			menu.game.Session.RemoteClient.GameData.OnGameState = true
			go menu.game.Session.RemoteClient.SendGameDataMessage(*menu.game.Session.RemoteClient.GameData)
		}
		menu.game.Session.JustJoined = false
		return NewGame(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight, menu.game)
	}
	if menu.game.Session.End {
		CleanWorld(menu.world)
		menu.game.Session.RemoteClient.ResetListeners()
		return NewStartMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight)
	}
	return menu
}

func (menu *WinnerMenu) Update() {
	if menu.game.Muted {
		archetype.StopAudioGame()
		archetype.StopAudioMenu()
		archetype.StopWinnerAudio()
	} else {
		archetype.StopAudioGame()
		archetype.StopAudioMenu()
		archetype.PlayWinnerAudio(menu.game)
	}
	if menu.game.Session.JustJoined {
		menu.game.Session.JustJoined = false
	}

	for _, s := range menu.systems {
		s.Update(*menu.world)
	}
	menu.uiHandler.Update()
}

func (menu *WinnerMenu) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range menu.drawables {
		s.Draw(*menu.world, screen)
	}
}
