package main

import (
	"fmt"
	"github.com/bquenin/tmxmap"
	"image"
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
	frameBuffer = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	engine      *tuile.Engine
	overworld   *tuile.Layer
	x, y        int
)

func update(screen *ebiten.Image) error {
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

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Draw the frame
	engine.DrawFrame()

	// Display it on screen
	_ = screen.ReplacePixels(frameBuffer.Pix)

	// Draw the message
	msg := fmt.Sprintf("TPS: %f, x: %d, y: %d\n", ebiten.CurrentTPS(), x, y)
	_ = ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	engine = tuile.NewEngine(screenWidth, screenHeight)
	engine.SetBackgroundColor(color.Black)
	engine.SetPlot(func(x, y int, r, g, b, a byte) {
		i := frameBuffer.PixOffset(x, y)
		frameBuffer.Pix[i] = r
		frameBuffer.Pix[i+1] = g
		frameBuffer.Pix[i+2] = b
		frameBuffer.Pix[i+3] = a
	})

	tileMap, err := tmxmap.Load("../assets/zelda3/overworld.tmx")
	if err != nil {
		log.Fatal(err)
	}

	overworld, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}
	//overworld.SetRepeat(true)
	engine.AddLayer(overworld)

	if err := ebiten.Run(update, screenWidth, screenHeight, 4, "scrolling"); err != nil {
		log.Fatal(err)
	}
}
