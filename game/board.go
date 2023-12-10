package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
)

var taskTerminated = errors.New("uttt: task terminated")

type task func() error

// Board represents the game board.
type Board struct {
	size  int
	tiles map[*Tile]struct{}
	tasks []task
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size:  size,
		tiles: map[*Tile]struct{}{},
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			b.tiles[NewTile(-1, i, j)] = struct{}{}
		}
	}
	return b, nil
}

func (b *Board) tileAt(x, y int) *Tile {
	return tileAt(b.tiles, x, y)
}

// Update updates the board state.
func (b *Board) Update(tile *Tile) error {

	if 0 < len(b.tasks) {
		t := b.tasks[0]
		if err := t(); err == taskTerminated {
			b.tasks = b.tasks[1:]
		} else if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}

// Draw draws the board to the given boardImage.
func (b *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			v := 0
			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(tileBackgroundColor(v))
			boardImage.DrawImage(tileImage, op)
		}
	}
	for t := range b.tiles {
		t.Draw(boardImage)
	}
}
