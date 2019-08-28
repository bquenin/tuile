package tuile

import (
	"image"
	"image/color"
)

type HBlank func(line int)

// Engine structure
type Engine struct {
	hBlank          HBlank
	backgroundColor color.Color
	width           int
	height          int
	pixels          *image.RGBA
	layers          []*Layer
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewEngine instantiates a new tuile engine
func NewEngine(width, height int) *Engine {
	return &Engine{
		width:  width,
		height: height,
		pixels: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (t *Engine) SetHBlank(hBlank HBlank) {
	t.hBlank = hBlank
}

func (t *Engine) SetBackgroundColor(color color.Color) {
	t.backgroundColor = color
}

func (t *Engine) DrawFrame() *image.RGBA {
	for line := 0; line < t.height; line++ {
		if t.hBlank != nil {
			t.hBlank(line)
		}
		if t.backgroundColor != nil {
			t.fillBackgroundLine(line, t.backgroundColor, t.width)
		}
		for _, layer := range t.layers {
			if layer.repeat {
				t.drawLayerLineRepeat(line, layer)
			} else {
				t.drawLayerLine(line, layer)
			}
		}
	}
	return t.pixels
}

func (t *Engine) fillBackgroundLine(line int, color color.Color, width int) {
	r, g, b, _ := color.RGBA()
	for x := 0; x < width; x++ {
		i := t.pixels.PixOffset(x, line)
		t.pixels.Pix[i] = uint8(r)
		t.pixels.Pix[i+1] = uint8(g)
		t.pixels.Pix[i+2] = uint8(b)
	}
}

func (t *Engine) AddLayer(layer ...*Layer) {
	t.layers = append(t.layers, layer...)
}

func (t *Engine) drawTilePixel(x, y int, layer *Layer, tileID int, xTile, yTile int) {
	xImage := tileID % layer.tileMap.TileSets[0].Columns
	yImage := tileID / layer.tileMap.TileSets[0].Columns
	xImage *= layer.tileWidth
	yImage *= layer.tileHeight
	xImage += xTile % layer.tileWidth
	yImage += yTile % layer.tileHeight
	src := layer.image.PixOffset(xImage, yImage)
	r, g, b, a := layer.image.Palette[layer.image.Pix[src]].RGBA()
	if a == 0 {
		return
	}

	dst := t.pixels.PixOffset(x, y)
	t.pixels.Pix[dst] = uint8(r)
	t.pixels.Pix[dst+1] = uint8(g)
	t.pixels.Pix[dst+2] = uint8(b)
}

func (t *Engine) drawLayerLine(line int, layer *Layer) {
	if line < layer.origin.Y || line >= layer.origin.Y+layer.pixelHeight {
		return // Out of vertical bounds
	}

	for x := max(0, layer.origin.X); x < min(t.width, layer.origin.X+layer.pixelWidth); x++ {
		xTile := x - layer.origin.X
		yTile := line - layer.origin.Y
		tile := layer.tileMap.Layers[0].Tiles[yTile/layer.tileHeight*layer.width+xTile/layer.tileWidth]
		if tile.Nil {
			continue
		}
		t.drawTilePixel(x, line, layer, int(tile.ID), xTile, yTile)
	}
}

func (t *Engine) drawLayerLineRepeat(line int, layer *Layer) {
	for x := 0; x < t.width; x++ {
		xTile := (layer.origin.X + x) % layer.pixelWidth
		yTile := (layer.origin.Y + line) % layer.pixelHeight
		if xTile < 0 {
			xTile += layer.pixelWidth
		}
		if yTile < 0 {
			yTile += layer.pixelHeight
		}
		tile := layer.tileMap.Layers[0].Tiles[yTile/layer.tileHeight*layer.width+xTile/layer.tileWidth]
		if tile.Nil {
			continue
		}
		t.drawTilePixel(x, line, layer, int(tile.ID), xTile, yTile)
	}
}
