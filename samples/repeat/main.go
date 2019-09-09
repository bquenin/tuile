package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/bquenin/tmxmap"
	"github.com/bquenin/tuile"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 200
)

var (
	frameBuffer = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	engine      *tuile.Engine
	layer       *tuile.Layer
	x, y        = 0, 0
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
	layer.SetOrigin(x, y)

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
	engine.SetBackgroundColor(color.RGBA{R: 0x66, G: 0xCC, B: 0xFF})
	engine.SetPlot(func(x, y int, r, g, b, a byte) {
		i := frameBuffer.PixOffset(x, y)
		frameBuffer.Pix[i] = r
		frameBuffer.Pix[i+1] = g
		frameBuffer.Pix[i+2] = b
		frameBuffer.Pix[i+3] = a
	})

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

	if err := ebiten.Run(update, screenWidth, screenHeight, 4, "repeat"); err != nil {
		log.Fatal(err)
	}
}
