package main

import (
	"fmt"
	"github.com/bquenin/tmxmap"
	"image/color"
	"log"

	"github.com/bquenin/tuile"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x--
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x++
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y++
	}
	overworld.SetOrigin(x<<2, y<<2)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Draw the frame
	frame := engine.DrawFrame()

	// Display it on screen
	_ = screen.ReplacePixels(frame.Pix)

	// Draw the message
	msg := fmt.Sprintf("TPS: %f\n", ebiten.CurrentTPS())
	_ = ebitenutil.DebugPrint(screen, msg)
	return nil
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

	if err := ebiten.Run(update, screenWidth, screenHeight, 4, "scrolling"); err != nil {
		log.Fatal(err)
	}
}
