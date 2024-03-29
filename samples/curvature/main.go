package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/bquenin/tmxmap"
	"github.com/bquenin/tuile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 512
	screenHeight = 448
)

var (
	engine                      *tuile.Engine
	clouds, overworld           *tuile.Layer
	x, y                        = 0, 64
	offsets                     = [screenHeight]float64{}
	cloudsRatio, overworldRatio = 64.0, 160.0
)

func lerp(x2, x1, x3, y1, y3 int) float64 {
	return float64((x2-x1)*(y3-y1))/float64(x3-x1) + float64(y1)
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

	// Render off-screen
	g.offscreen.ReplacePixels(engine.Render())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)

	// Draw the message
	msg := fmt.Sprintf("TPS: %.f\nUP/DOWN/LEFT/RIGHT to move\nQ/A to curve clouds\nW/S to curve world", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func main() {
	for n := 0; n < screenHeight; n++ {
		offsets[n] = math.Tan(lerp(n, 0, screenHeight, 105.0, 180.0) * math.Pi / 180)
	}
	engine = tuile.NewEngine(screenWidth, screenHeight)
	engine.SetBackgroundColor(color.Black)
	engine.SetHBlank(hBlank)

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

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func hBlank(line int) {
	clouds.SetOrigin(x<<2, y<<2+int(offsets[line]*cloudsRatio)-line)
	overworld.SetOrigin(x<<2, y<<1+int(offsets[line]*overworldRatio)-line)
}
