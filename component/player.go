package component

import (
	"amaru/net"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/donburi"
)

type PlayerInputs struct {
	Up    ebiten.Key
	Right ebiten.Key
	Down  ebiten.Key
	Left  ebiten.Key
	Shoot ebiten.Key
}

type PlayerSettings struct {
	Inputs PlayerInputs
}

type PlayerLabelData struct {
	Name  string
	Color color.Color
}

type RemotePlayerMessage struct {
	Source       string
	Position     *net.Point
	LastPosition *net.Point
	Vector       net.Point
	Animation    string
}

type PlayerData struct {
	ID                  string
	Name                string
	Local               bool
	Body                *cp.Body
	Shape               *cp.Shape
	Space               *cp.Space
	PlayerSettings      *PlayerSettings
	Collision           bool
	PlayerCollision     bool
	LastPlayerCollision *time.Time
	LastDirection       *cp.Vector
	Label               *donburi.Entry
	OutOfBounds         bool
}

var Player = donburi.NewComponentType[PlayerData]()
var PlayerLabel = donburi.NewComponentType[PlayerLabelData]()
