package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/bquenin/tmxmap"
	"github.com/bquenin/tuile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 200
)

var (
	engine *tuile.Engine
	layer  *tuile.Layer
	x, y   = 0, 0
)

type Game struct {
	offscreen *ebiten.Image
}

func NewGame() *Game {
	return &Game{
		offscreen: ebiten.NewImage(screenWidth, screenHeight),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y--
	}
	layer.SetOrigin(x, y)

	// Render off-screen
	g.offscreen.ReplacePixels(engine.Render())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)

	// Draw the message
	msg := fmt.Sprintf("TPS: %f\n", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func main() {
	engine = tuile.NewEngine(screenWidth, screenHeight)
	engine.SetBackgroundColor(color.RGBA{R: 0x66, G: 0xCC, B: 0xFF})

	tileMap, err := tmxmap.Load("../assets/clouds.tmx")
	if err != nil {
		log.Fatal(err)
	}

	layer, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}
	layer.SetRepeat(true)
	engine.AddLayer(layer)

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
