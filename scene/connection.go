package scene

import (
	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"amaru/net"
	"amaru/system"
	"context"
	"time"

	"github.com/jakecoffman/cp"
	"github.com/nmorenor/chezmoi-net/client"
	cnet "github.com/nmorenor/chezmoi-net/net"
	"github.com/samber/lo"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"golang.org/x/image/colornames"
)

const (
	connectingLabel = "Connecting..."
)

type ConnectingMenu struct {
	world     *donburi.World
	game      *component.GameData
	systems   []System
	drawables []Drawable

	screenWidth  int
	screenHeight int

	offscreen *ebiten.Image
	connected bool

	remoteClient *net.RemoteClient

	wateLocations map[string]*net.WasteLocation
}

func NewConnectingMenuMenu(screenWidth int, screenHeight int, session *component.SessionData) *ConnectingMenu {
	menu := &ConnectingMenu{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
	}

	menu.loadMenu(session)

	return menu
}

func (menu *ConnectingMenu) loadMenu(session *component.SessionData) {
	selectedLevelIndex := engine.RandomIntRange(0, assets.GameLevelLoader.LevelsSize)
	assets.GameLevelLoader.LoadLevel(selectedLevelIndex)
	render := system.NewRenderer()

	menu.systems = []System{
		system.NewCamera(),
		render,
	}

	menu.drawables = []Drawable{
		render,
	}

	menu.world = engine.Ptr(menu.createWorld(selectedLevelIndex, session))
	menu.game = component.MustFindGame(*menu.world)
}

func (menu *ConnectingMenu) UpdateLayout(width, height int) {
	// do nothing
}

func (menu *ConnectingMenu) createWorld(levelIndex int, session *component.SessionData) donburi.World {
	rectX := float64(menu.screenWidth/2) - (float64(menu.screenWidth/2) / 2)
	rectY := float64(menu.screenHeight/2) - (float64(menu.screenHeight/2) / 2)

	menuContainerWidth := float64(menu.screenWidth / 2)
	menuContainerHeight := float64(menu.screenHeight / 2)
	menuContainerImage := archetype.DrawMainMenuRoundedRect(menu.offscreen, rectX, rectY, menuContainerWidth, menuContainerHeight, 5, colornames.White, assets.BlueColor, borderWidth, menuTitle)

	textSize := text.BoundString(assets.MainBigFont, connectingLabel)
	textWidth, _ := textSize.Dx(), textSize.Dy()
	textX := (float64(menuContainerImage.Bounds().Dx()) - float64(textWidth)) / 2
	textY := float64(menuContainerImage.Bounds().Dy() / 2)
	text.Draw(
		menuContainerImage,
		connectingLabel,
		assets.MainBigFont,
		int(textX),
		int(textY),
		assets.BlueColor,
	)

	world := donburi.NewWorld()

	archetype.NewInput(world)

	level := world.Entry(world.Create(component.Level))
	component.Level.Get(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	cameraEntry := archetype.NewCamera(world, menu.screenWidth, menu.screenHeight, math.Vec2{
		X: 0,
		Y: 0,
	})

	component.Camera.Get(cameraEntry).Disabled = true

	selectedLevel := assets.GameLevelLoader.CurrentLevel

	space, shapes := archetype.SetupSpaceForLevel(selectedLevel)

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
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(menuUIEntry, component.SpriteData{
		Image: menuContainerImage,
		Layer: component.SpriteLayerUI,
		Pivot: component.SpritePivotScreenCenter,
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
	menu.wateLocations = map[string]*net.WasteLocation{}
	for _, location := range wastePaths {
		menu.wateLocations[location.Id] = location
	}

	if menu.world == nil {
		game := world.Entry(world.Create(component.Game))
		component.Game.SetValue(game, component.GameData{
			Settings: component.Settings{
				ScreenWidth:  menu.screenWidth,
				ScreenHeight: menu.screenHeight,
			},
			Speed:        3.0,
			LeftOffset:   100,
			Session:      session,
			ChatMessages: engine.NewQueue[net.ChatMessage](),
			WasteSize:    wasteSize,
		})
		menu.game = component.Game.Get(game)
	}

	archetype.StopAudioMenu()
	archetype.PlayWavesAudio()

	return world
}

func (menu *ConnectingMenu) StartSession() {
	menu.remoteClient = net.NewRemoteClient(client.NewClient(cnet.New(net.ConnectionURL)), *menu.game.Session.UserName, menu.game.Session.Type == component.SessionTypeHost)
	if menu.game.Session.SessionID != nil {
		menu.remoteClient.Session = menu.game.Session.SessionID
		menu.remoteClient.Client.Session = menu.game.Session.SessionID
	}
	menu.game.Session.RemoteClient = menu.remoteClient
	menu.game.Session.RemoteClient.GameData.Counter = 30
	menu.game.Session.RemoteClient.GameData.Frames = 0
	if menu.game.Session.Type == component.SessionTypeHost {
		menu.game.Session.RemoteClient.GameData.LevelIndex = assets.GameLevelLoader.CurrentLevelIndex
		menu.game.Session.RemoteClient.GameData.WasteLocations = menu.wateLocations
	}
	menu.remoteClient.SessionJoin.AddListener(func(ctx context.Context, sjm net.SessionJoinMessage) {
		if menu.game.Session.RemoteClient.GameData == nil {
			return
		}
		if *menu.remoteClient.Client.Id == sjm.Target {
			return
		}
		participant := net.SessionParticipant{
			Id:       sjm.Target,
			Name:     sjm.Client.Participants[sjm.Target],
			Position: sjm.Position,
			Anim:     sjm.Anim,
		}
		if menu.game.Session.RemoteClient.GameData.SessionParticipants == nil {
			menu.game.Session.RemoteClient.GameData.SessionParticipants = map[string]*net.SessionParticipant{}
		}
		menu.game.Session.RemoteClient.GameData.SessionParticipants[participant.Id] = &participant
	})
	go func(c *client.Client) {
		c.Connect()
	}(menu.remoteClient.Client)
}

func (menu *ConnectingMenu) NextScene() archetype.Scene {
	if menu.connected {

		if menu.game.Session.Type == component.SessionTypeJoin && !menu.game.Session.RemoteClient.GameData.OnGameState {
			CleanWorld(menu.world)
			menu.world = nil
			menu.systems = nil
			menu.drawables = nil
			return NewWinnerMenu(menu.game)
		}
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		return NewGame(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight, menu.game)
	}
	if menu.remoteClient.InvalidSession {
		CleanWorld(menu.world)
		menu.world = nil
		menu.systems = nil
		menu.drawables = nil
		return NewStartMenu(menu.game.Settings.ScreenWidth, menu.game.Settings.ScreenHeight)
	}
	return menu
}

func (menu *ConnectingMenu) Update() {
	archetype.StopAudioMenu()
	archetype.PlayWavesAudio()
	if menu.remoteClient == nil {
		menu.StartSession()
	} else {
		if menu.remoteClient.Ready && !menu.remoteClient.InitAndReady {
			if menu.game.Session.Type == component.SessionTypeJoin {
				menu.game.Session.JustJoined = true
			}
			go menu.remoteClient.Initialize()
		}
		if menu.remoteClient.Ready && menu.remoteClient.InitAndReady {
			menu.connected = true
		}
		if menu.remoteClient.InvalidSession {
			menu.remoteClient.Client.Close()
		}
	}

	for _, s := range menu.systems {
		s.Update(*menu.world)
	}
}

func (menu *ConnectingMenu) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range menu.drawables {
		s.Draw(*menu.world, screen)
	}
}
