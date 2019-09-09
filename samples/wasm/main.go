// +build wasm

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"net/http"
	"syscall/js"

	"github.com/bquenin/tmxmap"
	"github.com/bquenin/tuile"
)

const (
	screenWidth  = 256
	screenHeight = 224
)

var (
	engine      *tuile.Engine
	layer       *tuile.Layer
	x, y, θ     = .0, .0, math.Pi
	top, bottom = .2, 5.0
	ratio       = .5

	ctx, frameBuffer js.Value
)

func lerp(x2, x1, x3, y1, y3 float64) float64 {
	return (x2-x1)*(y3-y1)/x3 - x1 + y1
}

func main() {
	// Load resources
	tmxRequest, err := http.NewRequest("GET", "overworld.tmx", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	tmxResponse, err := http.DefaultClient.Do(tmxRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tmxResponse.Body.Close()

	pngRequest, err := http.NewRequest("GET", "overworld.png", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	pngResponse, err := http.DefaultClient.Do(pngRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pngResponse.Body.Close()

	// Initialize tuile
	engine = tuile.NewEngine(screenWidth, screenHeight)
	engine.SetBackgroundColor(color.Black)
	engine.SetHBlank(hBlank)

	tileMap, err := tmxmap.Decode(tmxResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	layer, err = tuile.NewLayerWithReader(tileMap, pngResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	layer.SetTranslation(screenWidth/2, screenHeight)
	engine.AddLayer(layer)

	// Initialize canvas
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "tuile")
	canvas.Set("width", screenWidth)
	canvas.Set("height", screenHeight)
	ctx = canvas.Call("getContext", "2d")
	frameBuffer = js.Global().Get("Uint8Array").New(4 * screenWidth * screenHeight)

	done := make(chan struct{}, 0)

	// Render
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		x++
		y++
		θ += 0.01
		layer.SetOrigin(int(x), int(y))
		layer.SetRotation(θ)

		// Draw the frame
		frame := engine.DrawFrame()

		js.CopyBytesToJS(frameBuffer, frame.Pix)
		clamped := js.Global().Get("Uint8ClampedArray").New(frameBuffer)
		imgData := js.Global().Get("ImageData").New(clamped, screenWidth, screenHeight)
		ctx.Call("putImageData", imgData, 0, 0)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done
}

func hBlank(line int) {
	scale := lerp(float64(line), 0, float64(screenHeight), top, bottom)
	layer.SetScale(scale*ratio, scale*ratio)
}
