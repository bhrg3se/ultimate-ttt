package game

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"uttt/resources/fonts"
)

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// TileData represents a tile information like a value and a position.
type TileData struct {
	value int
	x     int
	y     int
}

// Tile represents a tile information including TileData and animation states.
type Tile struct {
	current TileData
}

// Pos returns the tile's current position.
// Pos is used only at testing so far.
func (t *Tile) Pos() (int, int) {
	return t.current.x, t.current.y
}

// Value returns the tile's current value.
// Value is used only at testing so far.
func (t *Tile) Value() int {
	return t.current.value
}

// NewTile creates a new Tile object.
func NewTile(value int, x, y int) *Tile {
	return &Tile{
		current: TileData{
			value: value,
			x:     x,
			y:     y,
		},
	}
}

func tileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
	var result *Tile

	const x_offset = 0
	const y_offset = 0

	xx := x_offset + (x-tileMargin)/(tileSize+tileMargin)
	yy := y_offset + (y-tileMargin)/(tileSize+tileMargin)

	for t := range tiles {
		if t.current.x != xx || t.current.y != yy {
			continue
		}
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

const (
	maxMovingCount  = 5
	maxPoppingCount = 6
)

// Update updates the tile's animation states.
func (t *Tile) Update(user int) error {
	t.current.value = user
	return nil
}

const (
	tileSize   = 80
	tileMargin = 4
)

var (
	tileImage = ebiten.NewImage(tileSize, tileSize)
)

func init() {
	tileImage.Fill(color.White)
}

// Draw draws the current tile to the given boardImage.
func (t *Tile) Draw(boardImage *ebiten.Image) {
	i, j := t.current.x, t.current.y

	v := t.current.value

	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin

	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(tileBackgroundColor(v))
	boardImage.DrawImage(tileImage, op)

	str := ""
	if v == 0 {
		str = "O"
	} else if v == 1 {
		str = "X"
	}

	f := mplusBigFont
	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}

	w := font.MeasureString(f, str).Floor()
	h := (f.Metrics().Ascent + f.Metrics().Descent).Floor()
	x += (tileSize - w) / 2
	y += (tileSize-h)/2 + f.Metrics().Ascent.Floor()
	text.Draw(boardImage, str, f, x, y, tileColor(v))
}
