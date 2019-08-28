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
	screenWidth  = 320
	screenHeight = 200
)

var (
	engine                 *tuile.Engine
	background, foreground *tuile.Layer
	x, y                   int
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
	background.SetOrigin(x, screenHeight/2+y)
	foreground.SetOrigin(x*2, screenHeight/2+y*2)

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

	tileMap, err := tmxmap.Load("../assets/track1_bg.tmx")
	if err != nil {
		log.Fatal(err)
	}

	background, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}
	engine.AddLayer(background)

	foreground, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}
	engine.AddLayer(foreground)

	if err := ebiten.Run(update, screenWidth, screenHeight, 4, "scrolling"); err != nil {
		log.Fatal(err)
	}
}
