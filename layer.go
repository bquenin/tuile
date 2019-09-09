package tuile

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

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
	Image                   *image.Paletted
	repeat                  bool
	disabled                bool
	transformed             bool
	angle                   float64 //radians
	translation             Vector
	scale                   Vector
}

func NewLayerWithReader(tileMap *tmxmap.Map, reader io.Reader) (*Layer, error) {
	var err error
	tileMap.TileSets[0].Image.Image, _, err = image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return NewLayer(tileMap)
}

// NewLayer instantiates a new layer
func NewLayer(tileMap *tmxmap.Map) (*Layer, error) {
	image, ok := tileMap.TileSets[0].Image.Image.(*image.Paletted)
	if !ok {
		return nil, errors.New("tileset Image is not paletted")
	}
	return &Layer{
		tiles:       tileMap.Layers[0].Tiles,
		tileSet:     &tileMap.TileSets[0],
		Image:       image,
		width:       tileMap.Layers[0].Width,
		height:      tileMap.Layers[0].Height,
		pixelWidth:  tileMap.Layers[0].Width * tileMap.TileSets[0].TileWidth,
		pixelHeight: tileMap.Layers[0].Height * tileMap.TileSets[0].TileHeight,
		tileWidth:   tileMap.TileSets[0].TileWidth,
		tileHeight:  tileMap.TileSets[0].TileHeight,
		repeat:      false,
		scale:       VInt(1, 1),
	}, nil
}

func (l *Layer) SetOrigin(x int, y int) {
	l.origin = image.Pt(x, y)
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

func (l *Layer) SetRotation(angle float64) {
	l.transformed = true
	l.angle = angle
}

func (l *Layer) SetScale(sx, sy float64) {
	l.transformed = true
	l.scale = V(sx, sy)
}

func (l *Layer) SetTranslation(dx, dy float64) {
	l.transformed = true
	l.translation = V(dx, dy)
}

func (l *Layer) transform(left, right Vector) (Vector, Vector) {
	dx, dy := float64(l.origin.X)+l.translation.X, float64(l.origin.Y)+l.translation.Y
	translate := IM.Translate(V(-dx, -dy))
	rotation := translate.Rotate(l.angle)
	scale := rotation.Scale(l.scale)
	result := scale.Translate(V(dx, dy))
	return left.Mul(result), right.Mul(result)
}
