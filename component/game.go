package component

import (
	"amaru/engine"
	"amaru/net"

	vpad "github.com/kemokemo/ebiten-virtualpad"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type SessionType int

const (
	MaxWasteSize          = 260
	MinWasteSize          = 220
	WastePoints           = 1
	AnimalPoints          = 2
	PlayerCollisionPoints = 2
	MaxPlayers            = 4

	SessionTypeHost SessionType = iota
	SessionTypeJoin
)

type GameData struct {
	GameOver         bool
	Settings         Settings
	Speed            float64
	LeftOffset       float64
	Session          *SessionData
	CursorOverButton bool
	ChatMessages     *engine.Queue[net.ChatMessage]
	WasteSize        int
	CollectedWaste   int
	Dpad             *vpad.DirectionalPad
	Muted            bool
}

type SessionData struct {
	Type          SessionType
	End           bool
	JustJoined    bool
	SessionID     *string
	UserName      *string
	RemoteClient  *net.RemoteClient
	PlayerMessage map[string]*RemotePlayerMessage
}

type Settings struct {
	ScreenWidth  int
	ScreenHeight int
}

var Game = donburi.NewComponentType[GameData]()

func MustFindGame(w donburi.World) *GameData {
	game, ok := query.NewQuery(filter.Contains(Game)).First(w)
	if !ok {
		panic("game not found")
	}
	return Game.Get(game)
}
