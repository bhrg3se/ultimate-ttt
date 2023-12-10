package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	ScreenWidth  = (80 * 9) + (9 * 4)
	ScreenHeight = (80 * 9) + (9 * 4)
	boardSize    = 9
)

// Game represents a game state.
type Game struct {
	currentUser     int
	nextMiniBoardID int
	input           *Input
	board           *Board
	boardImage      *ebiten.Image
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{
		nextMiniBoardID: -1,
		input:           NewInput(),
	}
	var err error
	g.board, err = NewBoard(boardSize)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tile := g.board.tileAt(x, y)

		if tile.Value() != -1 {
			fmt.Println("already clicked")
			return nil
		}

		if g.nextMiniBoardID >= 0 && g.nextMiniBoardID != g.getMiniBoard(tile.Pos()) {
			fmt.Println("not allowed ", g.nextMiniBoardID, g.getMiniBoard(tile.Pos()))
			return nil
		}

		err := tile.Update(g.currentUser)
		if err != nil {
			return err
		}

		err = g.board.Update(tile)
		if err != nil {
			return err
		}
		g.currentUser = (g.currentUser + 1) % 2
		//miniBoardID := g.getMiniBoard(tile.Pos())
		g.nextMiniBoardID = g.getNextMiniBoard(tile.Pos())
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(g.board.Size())
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))

	vector.StrokeLine(g.boardImage, 3*tileSize+4*tileMargin, 0, 3*tileSize+4*tileMargin, 9*tileSize+10*tileMargin, tileMargin, frame2Color, true)
	vector.StrokeLine(g.boardImage, 6*tileSize+7*tileMargin, 0, 6*tileSize+7*tileMargin, 9*tileSize+10*tileMargin, tileMargin, frame2Color, true)
	vector.StrokeLine(g.boardImage, 0, 3*tileSize+4*tileMargin, 9*tileSize+10*tileMargin, 3*tileSize+4*tileMargin, tileMargin, frame2Color, true)
	vector.StrokeLine(g.boardImage, 0, 6*tileSize+7*tileMargin, 9*tileSize+10*tileMargin, 6*tileSize+7*tileMargin, tileMargin, frame2Color, true)

	g.drawNextMiniBoard(g.boardImage)
	screen.DrawImage(g.boardImage, op)
}

func (g *Game) getMiniBoard(x, y int) int {
	x = x / 3
	y = y / 3
	return (3 * y) + x
}
func (g *Game) getNextMiniBoard(x, y int) int {
	x = x % 3
	y = y % 3
	return (3 * y) + x
}

func (g *Game) drawNextMiniBoard(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	i := g.nextMiniBoardID % 3
	j := g.nextMiniBoardID / 3

	x := (i*tileSize + (i+1)*tileMargin) * 3
	y := (j*tileSize + (j+1)*tileMargin) * 3

	op.GeoM.Scale(3, 3)

	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(nextMiniBoardColor)
	screen.DrawImage(tileImage, op)
}
