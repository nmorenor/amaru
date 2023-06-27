package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteLayer int

const (
	SpriteLayerBackground SpriteLayer = iota
	SpriteLayerDefault
	SpriteLayerForeground
	SpriteLayerUI
)

type SpritePivot int

const (
	SpritePivotCenter SpritePivot = iota
	SpritePivotTopLeft
	SpritePivotScreenCenter
)

type SpriteData struct {
	Image *ebiten.Image
	Layer SpriteLayer
	Pivot SpritePivot

	Hidden bool

	ColorOverride *ColorOverride
}

type UIImageTransformer func(*ebiten.Image) *ebiten.Image

type UISpriteData struct {
	Image     *ebiten.Image
	Layer     SpriteLayer
	Pivot     SpritePivot
	UIHandler UIImageTransformer

	Hidden bool

	ColorOverride *ColorOverride
}

type ColorOverride struct {
	R, G, B, A float64
}

func (s *SpriteData) Show() {
	s.Hidden = false
}

func (s *SpriteData) Hide() {
	s.Hidden = true
}

var Sprite = donburi.NewComponentType[SpriteData]()
var UISprite = donburi.NewComponentType[UISpriteData]()
