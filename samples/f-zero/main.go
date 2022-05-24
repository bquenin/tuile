package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/bquenin/tmxmap"
	"github.com/bquenin/tuile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 256
	screenHeight = 224
)

var (
	frameBuffer = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	engine      *tuile.Engine
	track       *tuile.Layer
	x, y, θ     = .0, .0, math.Pi
	top, bottom = 0.2, 5.0
	ratio       = .5
)

func lerp(x2, x1, x3, y1, y3 float64) float64 {
	return (x2-x1)*(y3-y1)/x3 - x1 + y1
}

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
		θ += 0.04
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		θ -= 0.04
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		x -= math.Sin(θ) * 8
		y += math.Cos(θ) * 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		x += math.Sin(θ) * 8
		y -= math.Cos(θ) * 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		ratio += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		ratio -= 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		top += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		top -= 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		bottom += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		bottom -= 0.01
	}
	track.SetOrigin(int(x), int(y))
	track.SetRotation(θ)

	// Draw the frame
	engine.DrawFrame()

	// Render it off-screen
	g.offscreen.ReplacePixels(frameBuffer.Pix)

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
	engine.SetBackgroundColor(color.Black)
	engine.SetHBlank(hBlank)
	engine.SetPlot(func(x, y int, r, g, b, a byte) {
		i := frameBuffer.PixOffset(x, y)
		frameBuffer.Pix[i] = r
		frameBuffer.Pix[i+1] = g
		frameBuffer.Pix[i+2] = b
		frameBuffer.Pix[i+3] = a
	})

	tileMap, err := tmxmap.Load("../assets/f-zero/mc1.tmx")
	if err != nil {
		log.Fatal(err)
	}

	track, err = tuile.NewLayer(tileMap)
	if err != nil {
		log.Fatal(err)
	}
	track.SetTranslation(screenWidth/2, screenHeight)
	engine.AddLayer(track)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func hBlank(line int) {
	scale := lerp(float64(line), 0, float64(screenHeight), top, bottom)
	track.SetScale(scale*ratio, scale*ratio)
}
