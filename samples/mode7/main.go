package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/bquenin/tmxmap"
	"github.com/bquenin/tuile"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 256
	screenHeight = 240
)

var (
	engine  *tuile.Engine
	track   *tuile.Layer
	x, y, θ = .0, .0, math.Pi
	ratio   = 4.0
)

func lerp(x2, x1, x3, y1, y3 float64) float64 {
	return (x2-x1)*(y3-y1)/x3 - x1 + y1
}

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		θ += 0.04
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		θ -= 0.04
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		x -= math.Sin(θ) * 4
		y += math.Cos(θ) * 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		x += math.Sin(θ) * 4
		y -= math.Cos(θ) * 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		ratio += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		ratio -= 0.02
	}
	track.SetOrigin(int(x), int(y))
	track.SetRotation(θ)

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
	engine.SetHBlank(hBlank)

	tileMap, err := tmxmap.Load("../assets/smk/track.tmx")
	if err != nil {
		log.Fatal(err)
	}

	track, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}
	track.SetTranslation(screenWidth/2, screenHeight)
	engine.AddLayer(track)

	if err := ebiten.Run(update, screenWidth, screenHeight, 4, "mode7"); err != nil {
		log.Fatal(err)
	}
}

func hBlank(line int) {
	scale := lerp(float64(line), 0, float64(screenHeight), 0.2, 5)
	track.SetScale(scale*ratio, scale*ratio)
}
