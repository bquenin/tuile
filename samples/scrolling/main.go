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
	screenWidth  = 256
	screenHeight = 240
)

var (
	engine    *tuile.Engine
	overworld *tuile.Layer
	x, y      int
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
	overworld.SetOrigin(x<<2, y<<2)

	// Render off-screen
	g.offscreen.ReplacePixels(engine.Render())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)

	// Draw the message
	msg := fmt.Sprintf("TPS: %.f\nUP/DOWN/LEFT/RIGHT to move", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func main() {
	engine = tuile.NewEngine(screenWidth, screenHeight)
	engine.SetBackgroundColor(color.Black)

	tileMap, err := tmxmap.Load("../assets/zelda3/overworld.tmx")
	if err != nil {
		log.Fatal(err)
	}

	overworld, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}

	engine.AddLayer(overworld)

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
