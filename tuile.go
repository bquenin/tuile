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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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
			if layer.disabled {
				continue
			}
			if layer.transformed {
				t.drawLayerLineAffine(line, layer)
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

func (t *Engine) drawLayerLine(line int, layer *Layer) {
	yTile := layer.origin.Y + line
	if yTile < 0 || yTile >= layer.pixelHeight {
		if !layer.repeat {
			return
		}
		// https://maurobringolf.ch/2017/12/a-neat-trick-to-compute-modulo-of-negative-numbers/
		yTile = (yTile%layer.pixelHeight + layer.pixelHeight) % layer.pixelHeight
	}

	for x := 0; x < t.width; {
		xTile := layer.origin.X + x
		if xTile < 0 || xTile >= layer.pixelWidth {
			if !layer.repeat {
				x++
				continue
			}
			// https://maurobringolf.ch/2017/12/a-neat-trick-to-compute-modulo-of-negative-numbers/
			xTile = (xTile%layer.pixelWidth + layer.pixelWidth) % layer.pixelWidth
		}
		tile := layer.tiles[yTile/layer.tileHeight*layer.width+xTile/layer.tileWidth]
		if tile.Nil {
			x++
			continue
		}

		yImage := int(tile.ID) / layer.tileSet.Columns
		yImage *= layer.tileHeight
		yImage += yTile % layer.tileHeight

		xImage := int(tile.ID) % layer.tileSet.Columns
		xImage *= layer.tileWidth
		for xx := xTile % layer.tileWidth; xx < layer.tileWidth && x < t.width; xx, x = xx+1, x+1 {
			var src int
			if tile.HorizontalFlip {
				src = layer.image.PixOffset(xImage+layer.tileWidth-1-xx, yImage)
			} else {
				src = layer.image.PixOffset(xImage+xx, yImage)
			}

			r, g, b, a := layer.image.Palette[layer.image.Pix[src]].RGBA()
			if a == 0 {
				continue
			}

			dst := t.pixels.PixOffset(x, line)
			t.pixels.Pix[dst] = uint8(r)
			t.pixels.Pix[dst+1] = uint8(g)
			t.pixels.Pix[dst+2] = uint8(b)
		}
	}
}

func (t *Engine) drawLayerLineAffine(line int, layer *Layer) {
	left, right := layer.transform(
		VInt(layer.origin.X, layer.origin.Y+line),
		VInt(layer.origin.X+t.width, layer.origin.Y+line),
	)

	x1, y1 := left.X, left.Y
	x2, y2 := right.X, right.Y

	dx := (x2 - x1) / float64(t.width)
	dy := (y2 - y1) / float64(t.width)

	for x := 0; x < t.width; x, x1, y1 = x+1, x1+dx, y1+dy {
		if !layer.repeat && (x1 < 0 || int(x1) >= layer.pixelWidth || y1 < 0 || int(y1) >= layer.pixelHeight) {
			continue
		}
		xTile := abs(int(x1)+layer.pixelWidth) % layer.pixelWidth
		yTile := abs(int(y1)+layer.pixelHeight) % layer.pixelHeight

		tile := layer.tiles[yTile/layer.tileHeight*layer.width+xTile/layer.tileWidth]
		if tile.Nil {
			continue
		}

		yImage := int(tile.ID) / layer.tileSet.Columns
		yImage *= layer.tileHeight
		yImage += yTile % layer.tileHeight

		xImage := int(tile.ID) % layer.tileSet.Columns
		xImage *= layer.tileWidth

		var src int
		if tile.HorizontalFlip {
			src = layer.image.PixOffset(xImage+layer.tileWidth-1-(xTile%layer.tileWidth), yImage)
		} else {
			src = layer.image.PixOffset(xImage+xTile%layer.tileWidth, yImage)
		}

		r, g, b, a := layer.image.Palette[layer.image.Pix[src]].RGBA()
		if a == 0 {
			continue
		}

		dst := t.pixels.PixOffset(x, line)
		t.pixels.Pix[dst] = uint8(r)
		t.pixels.Pix[dst+1] = uint8(g)
		t.pixels.Pix[dst+2] = uint8(b)
	}
}
