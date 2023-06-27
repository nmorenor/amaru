package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	vpad "github.com/kemokemo/ebiten-virtualpad"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/jakecoffman/cp"

	"amaru/assets"
	"amaru/component"
	"amaru/engine"
)

var (
	StopUpAction = &component.Animation{
		Name:   "upstop",
		Frames: []int{0},
	}

	StopDownAction = &component.Animation{
		Name:   "downstop",
		Frames: []int{6},
	}

	StopLeftAction = &component.Animation{
		Name:   "leftstop",
		Frames: []int{9},
	}

	StopRightAction = &component.Animation{
		Name:   "rightstop",
		Frames: []int{4},
	}

	WalkUpAction = &component.Animation{
		Name:   "up",
		Frames: []int{0, 1, 2},
		Loop:   true,
	}

	WalkDownAction = &component.Animation{
		Name:   "down",
		Frames: []int{6, 7, 8},
		Loop:   true,
	}

	WalkLeftAction = &component.Animation{
		Name:   "left",
		Frames: []int{9, 10, 11},
		Loop:   true,
	}

	WalkRightAction = &component.Animation{
		Name:   "right",
		Frames: []int{3, 4, 5},
		Loop:   true,
	}

	Actions = []*component.Animation{
		StopUpAction,
		StopDownAction,
		StopLeftAction,
		StopRightAction,
		WalkUpAction,
		WalkDownAction,
		WalkLeftAction,
		WalkRightAction,
	}

	ActionsByKey = make(map[string]*component.Animation)
)

func MustLoadPlayerActions() {
	for _, next := range Actions {
		if next == nil {
			continue
		}
		ActionsByKey[next.Name] = next
	}
}

func NewPlayer(w donburi.World, space *cp.Space, startPosition math.Vec2, intAnim string, name string, id string, localPlayer bool) *donburi.Entry {
	mass := 1.0
	moment := cp.MomentForBox(mass, 32, 32)
	body := space.AddBody(cp.NewBody(mass, moment))

	body.SetPosition(cp.Vector{X: startPosition.X, Y: (startPosition.Y)})

	shape := space.AddShape(cp.NewBox(body, 32, 32, 0))
	shape.SetFriction(1)
	shape.SetCollisionType(component.PlayerCollisionType)

	inputs := &component.PlayerSettings{
		Inputs: component.PlayerInputs{
			Up:    ebiten.KeyUp,
			Right: ebiten.KeyRight,
			Down:  ebiten.KeyDown,
			Left:  ebiten.KeyLeft,
			Shoot: ebiten.KeyEnter,
		},
	}
	if !localPlayer {
		inputs = nil
	}

	player := component.PlayerData{
		ID:             id,
		Name:           name,
		Local:          localPlayer,
		Body:           body,
		Shape:          shape,
		Space:          space,
		PlayerSettings: inputs,
	}
	pEntity := newPlayerFromPlayerData(w, math.Vec2{X: startPosition.X, Y: startPosition.Y}, intAnim, &player)
	shape.UserData = pEntity

	debug := component.MustFindDebug(w)
	debug.Shapes = append(debug.Shapes, shape)
	return pEntity
}

func newPlayerFromPlayerData(w donburi.World, startPosition math.Vec2, animation string, playerData *component.PlayerData) *donburi.Entry {
	player := w.Entry(
		w.Create(
			component.Player,
			component.Sprite,
			component.AnimationComponent,
			transform.Transform,
			component.Bounds,
		),
	)
	spriteData := component.SpriteData{}
	spriteData.Layer = component.SpriteLayerDefault
	spriteData.Pivot = component.SpritePivotTopLeft
	component.Sprite.SetValue(player, spriteData)

	animationEntry := NewAnimationComponent(w, player, player, assets.BoatSpriteSheet.Drawables(), 0.1)
	animationData := component.AnimationComponent.Get(animationEntry)
	animationData.AddAnimations(Actions)
	animationData.SelectAnimationByName(animation)

	spriteData.Image = animationData.Cell()
	component.Sprite.SetValue(player, spriteData)

	transform.Transform.Get(player).LocalPosition = startPosition
	component.Bounds.Get(player).Disabled = false

	playerLabel := w.Entry(
		w.Create(
			component.PlayerLabel,
			transform.Transform,
		),
	)
	labelData := component.PlayerLabel.Get(playerLabel)
	labelData.Name = playerData.Name
	labelData.Color = colornames.Orange
	if !playerData.Local {
		labelData.Color = colornames.Fuchsia
	}

	playerData.Label = playerLabel
	component.Player.SetValue(player, *playerData)
	transform.Transform.Get(playerLabel).LocalPosition = math.Vec2{X: startPosition.X, Y: startPosition.Y}

	return player
}

func FindPlayerByName(w donburi.World, name string) (*component.PlayerData, *donburi.Entry) {
	var foundPlayer *component.PlayerData
	var foundPlayerEntry *donburi.Entry
	query.NewQuery(filter.Contains(component.Player)).Each(w, func(e *donburi.Entry) {
		player := component.Player.Get(e)
		if player.Name == name {
			foundPlayer = player
			foundPlayerEntry = e
		}
	})

	return foundPlayer, foundPlayerEntry
}

func MustFindLocalPlayer(w donburi.World) (*component.PlayerData, *donburi.Entry) {
	var foundPlayer *component.PlayerData
	var foundPlayerEntry *donburi.Entry
	query.NewQuery(filter.Contains(component.Player)).Each(w, func(e *donburi.Entry) {
		player := component.Player.Get(e)
		if player.Local {
			foundPlayer = player
			foundPlayerEntry = e
		}
	})

	if foundPlayer == nil {
		panic("local player not found")
	}

	return foundPlayer, foundPlayerEntry
}

func FindRmotePlayers(w donburi.World) []*component.PlayerData {
	foundPlayers := make([]*component.PlayerData, 0)
	query.NewQuery(filter.Contains(component.Player)).Each(w, func(e *donburi.Entry) {
		player := component.Player.Get(e)
		if player.Local {
			foundPlayers = append(foundPlayers, player)
		}
	})

	return foundPlayers
}

func SetAnimation(w donburi.World, gameData *component.GameData, entry *donburi.Entry, player *component.PlayerData) *component.Animation {
	animationComponent := component.AnimationComponent.Get(entry)
	input := component.Input.Get(MustFindInput(w))
	var result *component.Animation
	if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) || GamePadIsUp(gameData) {
		animationComponent.SelectAnimationByAction(WalkUpAction)
		result = WalkUpAction
	} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) || GamePadIsDown(gameData) {
		animationComponent.SelectAnimationByAction(WalkDownAction)
		result = WalkDownAction
	} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) || GamePadIsLeft(gameData) {
		animationComponent.SelectAnimationByAction(WalkLeftAction)
		result = WalkLeftAction
	} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) || GamePadIsRight(gameData) {
		animationComponent.SelectAnimationByAction(WalkRightAction)
		result = WalkRightAction
	} else if input.IsKeyReleased(player.PlayerSettings.Inputs.Up) {
		animationComponent.SelectAnimationByAction(StopUpAction)
		result = StopUpAction
		if input.IsKeyReleased(player.PlayerSettings.Inputs.Left) {
			animationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) {
			animationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) {
			animationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) {
			animationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	} else if input.IsKeyReleased(player.PlayerSettings.Inputs.Down) {
		animationComponent.SelectAnimationByAction(StopDownAction)
		result = StopDownAction
		if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) {
			animationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) {
			animationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) {
			animationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) {
			animationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	} else if input.IsKeyReleased(player.PlayerSettings.Inputs.Left) {
		animationComponent.SelectAnimationByAction(StopLeftAction)
		result = StopLeftAction
		if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) {
			animationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) {
			animationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) {
			animationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) {
			animationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	} else if input.IsKeyReleased(player.PlayerSettings.Inputs.Right) {
		animationComponent.SelectAnimationByAction(StopRightAction)
		result = StopRightAction
		if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) {
			animationComponent.SelectAnimationByAction(WalkLeftAction)
			result = WalkLeftAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) {
			animationComponent.SelectAnimationByAction(WalkRightAction)
			result = WalkRightAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) {
			animationComponent.SelectAnimationByAction(WalkUpAction)
			result = WalkUpAction
		} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) {
			animationComponent.SelectAnimationByAction(WalkDownAction)
			result = WalkDownAction
		}
	}
	return result
}

func GamePadIsUp(gameData *component.GameData) bool {
	if gameData.Dpad == nil {
		return false
	}
	return gameData.Dpad.GetDirection() == vpad.Upper
}

func GamePadIsDown(gameData *component.GameData) bool {
	if gameData.Dpad == nil {
		return false
	}
	return gameData.Dpad.GetDirection() == vpad.Lower
}

func GamePadIsLeft(gameData *component.GameData) bool {
	if gameData.Dpad == nil {
		return false
	}
	return gameData.Dpad.GetDirection() == vpad.Left
}

func GamePadIsRight(gameData *component.GameData) bool {
	if gameData.Dpad == nil {
		return false
	}
	return gameData.Dpad.GetDirection() == vpad.Right
}

func GetSpeed(w donburi.World, gameData *component.GameData, player *component.PlayerData) (p cp.Vector, changed bool) {
	input := component.Input.Get(MustFindInput(w))
	p.X = input.Axis.X
	p.Y = input.Axis.Y
	origX, origY := p.X, p.Y

	if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) || GamePadIsUp(gameData) {
		p.Y = -1
	} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) || GamePadIsDown(gameData) {
		p.Y = 1
	}
	if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) || GamePadIsLeft(gameData) {
		p.X = -1
	} else if ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) || GamePadIsRight(gameData) {
		p.X = 1
	}
	if !engine.IsMobileBrowser() && !ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Left) && !ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Right) {
		p.X = 0
	}
	if engine.IsMobileBrowser() && !GamePadIsLeft(gameData) && !GamePadIsRight(gameData) {
		p.X = 0
	}
	if !engine.IsMobileBrowser() && !ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Up) && !ebiten.IsKeyPressed(player.PlayerSettings.Inputs.Down) {
		p.Y = 0
	}
	if engine.IsMobileBrowser() && !GamePadIsUp(gameData) && !GamePadIsDown(gameData) {
		p.Y = 0
	}

	changed = changed || p.X != origX || p.Y != origY
	return
}
