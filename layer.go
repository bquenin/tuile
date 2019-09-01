package tuile

import (
	"errors"
	"image"

	"github.com/bquenin/tmxmap"
)

// Layer structure
type Layer struct {
	origin                  image.Point
	width, height           int
	pixelWidth, pixelHeight int
	tileWidth, tileHeight   int
	tiles                   []*tmxmap.TileInfo
	tileSet                 *tmxmap.TileSet
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
		//tileMap:     tileMap,
		tiles:       tileMap.Layers[0].Tiles,
		tileSet:     &tileMap.TileSets[0],
		image:       image,
		width:       tileMap.Layers[0].Width,
		height:      tileMap.Layers[0].Height,
		pixelWidth:  tileMap.Layers[0].Width * tileMap.TileSets[0].TileWidth,
		pixelHeight: tileMap.Layers[0].Height * tileMap.TileSets[0].TileHeight,
		tileWidth:   tileMap.TileSets[0].TileWidth,
		tileHeight:  tileMap.TileSets[0].TileHeight,
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
