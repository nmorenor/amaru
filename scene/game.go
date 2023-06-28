package scene

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	vpad "github.com/kemokemo/ebiten-virtualpad"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"amaru/archetype"
	"amaru/assets"
	"amaru/component"
	"amaru/engine"
	"amaru/net"
	"amaru/system"
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Game struct {
	gameData  *component.GameData
	world     donburi.World
	systems   []System
	drawables []Drawable
	space     *cp.Space

	shapes []*cp.Shape

	screenWidth  int
	screenHeight int
}

func NewGame(screenWidth int, screenHeight int, gameData *component.GameData) *Game {
	g := &Game{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		shapes:       make([]*cp.Shape, 0),
		gameData:     gameData,
	}
	assets.MustLoadSvgs()

	space := cp.NewSpace()
	space.SetGravity(cp.Vector{X: 0, Y: 0})

	level := assets.GameLevelLoader.LoadLevel(gameData.Session.RemoteClient.GameData.LevelIndex)
	space, shapes := archetype.SetupSpaceForLevel(level)
	g.space = space
	g.shapes = shapes

	g.loadLevel()

	// reset all players score
	for id, sparticipant := range g.gameData.Session.RemoteClient.GameData.SessionParticipants {
		if sparticipant.Score > 0 {
			g.gameData.Session.RemoteClient.GameData.SessionParticipants[id].Score = 0
		}
	}

	return g
}

func (g *Game) NextScene() archetype.Scene {
	if g.gameData.GameOver {
		CleanWorld(&g.world)
		g.shapes = nil
		g.systems = nil
		g.gameData.Session.RemoteClient.ResetListeners()
		return NewWinnerMenu(g.gameData)
	}

	if g.gameData.Session.End {
		CleanWorld(&g.world)
		g.shapes = nil
		if g.gameData.Session != nil && g.gameData.Session.RemoteClient != nil {
			g.gameData.Session.RemoteClient.GameData.WasteLocations = nil
			g.gameData.Session.RemoteClient.GameData.OnGameState = false
			g.gameData.Session.RemoteClient.ResetListeners()
		}
		g.systems = nil
		return NewStartMenu(g.gameData.Settings.ScreenWidth, g.gameData.Settings.ScreenHeight)
	}
	return g
}

func (g *Game) loadLevel() {
	render := system.NewRenderer()
	debug := system.NewDebug()
	remote := system.NewRemoteSystem()
	hud := system.NewHUD()

	g.systems = []System{
		system.NewCamera(),
		system.NewCameraBounds(),
		system.NewAnimation(),
		remote,
		system.NewBounds(),
		system.NewPlayer(g.space),
		system.NewControls(),
		hud,
		render,
		debug,
	}

	g.drawables = []Drawable{
		render,
		debug,
		hud,
	}

	g.world = g.createWorld()
	remote.Initialize(g.gameData, g.world)
	if g.gameData.Session.Type == component.SessionTypeJoin {
		go g.gameData.Session.RemoteClient.RequestGameData()
	}
}

func (g *Game) UpdateLayout(width, height int) {
	// do nothing
}

func (g *Game) createWorld() donburi.World {
	levelAsset := assets.GameLevelLoader.CurrentLevel

	world := donburi.NewWorld()

	archetype.NewInput(world)

	level := world.Entry(world.Create(component.Level))
	component.Level.Get(level).ProgressionTimer = engine.NewTimer(time.Second * 3)

	archetype.NewCamera(world, levelAsset.Background.Bounds().Dx(), levelAsset.Background.Bounds().Dy(), math.Vec2{
		X: float64(levelAsset.Background.Bounds().Dx() / 2),
		Y: float64(levelAsset.Background.Bounds().Dy() / 2),
	})

	levelEntry := world.Entry(
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(levelEntry, component.SpriteData{
		Image: levelAsset.Background,
		Layer: component.SpriteLayerBackground,
		Pivot: component.SpritePivotTopLeft,
	})

	overPlayerEntry := world.Entry(
		world.Create(transform.Transform, component.Sprite),
	)
	component.Sprite.SetValue(overPlayerEntry, component.SpriteData{
		Image: levelAsset.OverPlayer,
		Layer: component.SpriteLayerForeground,
		Pivot: component.SpritePivotTopLeft,
	})

	g.gameData.GameOver = false
	g.gameData.WasteSize = len(g.gameData.Session.RemoteClient.GameData.WasteLocations)
	game := world.Entry(world.Create(component.Game))
	component.Game.SetValue(game, *g.gameData)

	g.gameData.GameOver = false
	g.gameData.WasteSize = len(g.gameData.Session.RemoteClient.GameData.WasteLocations)
	gameEntry := world.Entry(world.Create(component.Game))
	component.Game.SetValue(gameEntry, *g.gameData)

	g.gameData = component.MustFindGame(world)
	if engine.IsMobileBrowser() {
		g.gameData.Dpad = engine.Ptr(vpad.NewDirectionalPad(assets.DirectionalPad, assets.DirectionalBtn, engine.ConvertToRGBA(assets.BlueColor)))
		g.gameData.Dpad.SetLocation(g.gameData.Settings.ScreenWidth-230, g.gameData.Settings.ScreenHeight-230)
	}

	debugEntity := world.Create(component.Debug)
	debugComponent := component.Debug.Get(world.Entry(debugEntity))
	debugComponent.Shapes = g.shapes

	pPos := engine.RandomIntRange(0, len(levelAsset.PlayersStart))
	startPos := levelAsset.PlayersStart[pPos].TetraCenter()
	archetype.NewPlayer(world, g.space, startPos, component.DefaultPlayerAnimation, *g.gameData.Session.UserName, *g.gameData.Session.RemoteClient.Client.Id, true)
	if g.gameData.Session.RemoteClient.GameData.SessionParticipants[*g.gameData.Session.RemoteClient.Client.Id] == nil {
		g.gameData.Session.RemoteClient.GameData.SessionParticipants[*g.gameData.Session.RemoteClient.Client.Id] = &net.SessionParticipant{
			Id:        *g.gameData.Session.RemoteClient.Client.Id,
			Name:      g.gameData.Session.UserName,
			Position:  &net.Point{X: startPos.X, Y: startPos.Y},
			Anim:      engine.Ptr(component.DefaultPlayerAnimation),
			HasPlayer: false,
		}
	}
	go func() {
		g.gameData.Session.RemoteClient.SendInitialPositionDataMessage(net.Point{X: startPos.X, Y: startPos.Y})
	}()
	physics := world.Entry(world.Create(component.Physics))
	component.Physics.Get(physics).Space = g.space

	archetype.PlaceAnimalComponents(world, g.space, debugComponent, levelAsset.Animals, float64(levelAsset.Background.Bounds().Dx()), float64(levelAsset.Background.Bounds().Dy()))
	if g.gameData.Session.Type == component.SessionTypeHost {
		for _, loc := range g.gameData.Session.RemoteClient.GameData.WasteLocations {
			archetype.PlaceRemoteWasteFromPath(world, g.space, debugComponent, loc.Id, loc.Location, loc.Collected)
		}
	}

	if g.gameData.Session.Type == component.SessionTypeHost {
		g.gameData.Session.RemoteClient.GameData.OnGameState = true
	}

	archetype.SetupColliders(world)
	if !g.gameData.Muted {
		archetype.StopAudioMenu()
		archetype.StopWinnerAudio()
		archetype.PlayAudioGame()
	}

	return world
}

func (g *Game) Update() {
	if g.gameData.Muted {
		archetype.StopAudioGame()
		archetype.StopWinnerAudio()
		archetype.StopAudioGame()
	} else {
		archetype.StopAudioMenu()
		archetype.StopWinnerAudio()
		archetype.PlayAudioGame()
	}

	if g.gameData.Dpad != nil {
		g.gameData.Dpad.Update()
	}
	for _, s := range g.systems {
		s.Update(g.world)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range g.drawables {
		s.Draw(g.world, screen)
	}
	if g.gameData.Dpad != nil {
		g.gameData.Dpad.Draw(screen)
	}
}

func CleanWorld(w *donburi.World) {
	q := query.NewQuery(
		filter.Contains(component.Level),
	)
	world := *w
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	cam := archetype.MustFindCamera(world)
	world.Remove(cam.Entity())
	q = query.NewQuery(
		filter.Contains(transform.Transform, component.UISprite),
	)
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(
		filter.Contains(transform.Transform, component.Sprite),
	)
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(filter.Contains(
		component.PlayerLabel,
	))
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(filter.Contains(
		component.Input,
	))
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(filter.Contains(
		component.Player,
	))
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(filter.Contains(
		component.Waste,
	))
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(filter.Contains(
		component.Animal,
	))
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
	q = query.NewQuery(filter.Contains(
		component.Physics,
	))
	q.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})
}
