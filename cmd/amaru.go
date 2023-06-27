package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"amaru/archetype"
	"amaru/assets"
	"amaru/engine"
	"amaru/scene"
)

var (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	scene        archetype.Scene
	updateTicker *time.Ticker
}

func NewGame() *Game {
	assets.MustLoadAssets()
	archetype.MustLoadPlayerActions()

	g := &Game{
		updateTicker: time.NewTicker(time.Second / 60),
	}
	g.scene = scene.NewStartMenu(screenWidth, screenHeight)
	return g
}

func (g *Game) Update() error {
	g.scene = g.scene.NextScene()
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetFullscreen(false)
	ebiten.SetTPS(60)
	rand.Seed(time.Now().UTC().UnixNano())
	err := engine.InitClipboard()
	if err != nil {
		panic(err)
	}

	err = ebiten.RunGame(NewGame())
	if err != nil {
		log.Fatal(err)
	}
}
