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
			t.drawLayerLine(line, layer)
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
		//t.pixels.Pix[i+3] = a
	}
}

func (t *Engine) AddLayer(layer ...*Layer) {
	t.layers = append(t.layers, layer...)
}

func (t *Engine) drawLayerLine(line int, layer *Layer) {
	if line < layer.origin.Y || line >= layer.origin.Y+layer.pixelHeight {
		return // Out of vertical bounds
	}

	for x := max(0, layer.origin.X); x < min(t.width, layer.origin.X+layer.pixelWidth); x++ {
		xTile := (x - layer.origin.X) >> 3
		yTile := (line - layer.origin.Y) >> 3
		tile := layer.tileMap.Layers[0].Tiles[yTile*layer.width+xTile]
		if tile.Nil {
			continue
		}

		xImage := int(tile.ID) % layer.tileMap.TileSets[0].Columns
		yImage := int(tile.ID) / layer.tileMap.TileSets[0].Columns
		xImage *= layer.tileWidth
		yImage *= layer.tileHeight
		xImage += (x - layer.origin.X) % layer.tileWidth
		yImage += (line - layer.origin.Y) % layer.tileHeight
		src := layer.image.PixOffset(xImage, yImage)
		r, g, b, a := layer.image.Palette[layer.image.Pix[src]].RGBA()
		if a == 0 {
			continue
		}

		dst := t.pixels.PixOffset(x, line)
		t.pixels.Pix[dst] = uint8(r)
		t.pixels.Pix[dst+1] = uint8(g)
		t.pixels.Pix[dst+2] = uint8(b)
		//t.pixels.Pix[dst+3] = uint8(a)
	}
}
