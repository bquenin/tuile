package tuile

import (
	"errors"
	"image"

	"github.com/bquenin/tmxmap"
)

// Layer structure
type Layer struct {
	tileMap                 *tmxmap.Map
	origin                  image.Point
	width, height           int
	pixelWidth, pixelHeight int
	tileWidth, tileHeight   int
	image                   *image.Paletted
	repeat                  bool
	disabled                bool
}

// NewLayer instantiates a new layer
func NewLayer(tileMap *tmxmap.Map) (*Layer, error) {
	image, ok := tileMap.TileSets[0].Image.Image.(*image.Paletted)
	if !ok {
		return nil, errors.New("tileset image is not paletted")
	}
	return &Layer{
		tileMap:     tileMap,
		image:       image,
		width:       tileMap.Width,
		height:      tileMap.Height,
		pixelWidth:  tileMap.Width * tileMap.TileWidth,
		pixelHeight: tileMap.Height * tileMap.TileHeight,
		tileWidth:   tileMap.TileWidth,
		tileHeight:  tileMap.TileHeight,
		repeat:      false,
	}, nil
}

func (l *Layer) SetOrigin(x int, y int) {
	l.origin = image.Point{X: x, Y: y}
}

func (l *Layer) SetRepeat(repeat bool) {
	l.repeat = repeat
}

func (l *Layer) Disable() {
	l.disabled = true
}

func (l *Layer) Enable() {
	l.disabled = false
}
