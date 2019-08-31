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
	engine                      *tuile.Engine
	clouds, overworld           *tuile.Layer
	x, y                        = 0, 64
	offsets                     = [screenHeight]float64{}
	cloudsRatio, overworldRatio = 64.0, 320.0
)

func lerp(x2, x1, x3, y1, y3 int) float64 {
	return float64((x2-x1)*(y3-y1))/float64(x3-x1) + float64(y1)
}

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
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		cloudsRatio++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cloudsRatio--
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		overworldRatio++
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		overworldRatio--
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Draw the frame
	frame := engine.DrawFrame()

	// Display it on screen
	_ = screen.ReplacePixels(frame.Pix)

	// Draw the message
	msg := fmt.Sprintf("TPS: %.f, clouds(q/a): %d, overworld(w/s): %d\n", ebiten.CurrentTPS(), x, y)
	_ = ebitenutil.DebugPrint(screen, msg)
	return nil
}

func main() {
	engine = tuile.NewEngine(screenWidth, screenHeight)
	engine.SetBackgroundColor(color.Black)
	engine.SetHBlank(hBlank)
	for n := 0; n < screenHeight; n++ {
		offsets[n] = math.Tan(lerp(n, 0, screenHeight, 105.0, 180.0) * math.Pi / 180)
	}

	overworldMap, err := tmxmap.Load("../assets/zelda3/overworld.tmx")
	if err != nil {
		log.Fatal(err)
	}
	cloudsMap, err := tmxmap.Load("../assets/clouds.tmx")
	if err != nil {
		log.Fatal(err)
	}

	overworld, err = tuile.NewLayer(overworldMap)
	if err != nil {
		log.Fatal(err)
	}
	engine.AddLayer(overworld)

	clouds, err = tuile.NewLayer(cloudsMap)
	if err != nil {
		log.Fatal(err)
	}
	clouds.SetRepeat(true)
	engine.AddLayer(clouds)

	if err := ebiten.Run(update, screenWidth, screenHeight, 4, "curvature"); err != nil {
		log.Fatal(err)
	}
}

func hBlank(line int) {
	clouds.SetOrigin(x<<2, y<<2+int(offsets[line]*cloudsRatio)-line)
	overworld.SetOrigin(x<<2, y<<1+int(offsets[line]*overworldRatio)-line)
}
