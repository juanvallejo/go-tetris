package tetris

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/juanvallejo/tetris/pkg/tetris/grid"
	"github.com/juanvallejo/tetris/pkg/tetris/shape"
)

const (
	winWidth  = 800
	winHeight = 600

	cellMargin  = 45
	shapeMargin = 25
)

var winBgcolor = colornames.Darkslategrey

func NewGame() {
	config := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
		VSync:  true,
	}

	window, err := pixelgl.NewWindow(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	window.Clear(winBgcolor)

	gameWin := false
	shapeDecider := shape.NewShapeDecider(shape.CrossShape)
	bounds := window.Bounds()
	context := imdraw.New(nil)

	g := grid.NewGrid(pixel.V(0, 0), bounds.Max.X, bounds.Max.Y, grid.MaxCells, cellMargin)

	for !window.Closed() {
		context.Clear()

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			gameWin = handleMouseClick(window, context, g, shapeDecider, bounds, gameWin)
		}

		g.Render(context)
		context.Draw(window)
		window.Update()
	}
}

func handleMouseClick(window *pixelgl.Window, context *imdraw.IMDraw, g grid.Grid, shapeDecider *shape.ShapeDecider, bounds pixel.Rect, gameWin bool) bool {
	if gameWin {
		window.Clear(winBgcolor)
		g.Reset()
		return false
	}

	if cell := g.AtVector(window.MousePosition()); cell != nil {
		cell.Set(shape.NewShape(pixel.V(cell.Start().X, cell.Start().Y), shapeDecider.Next(), (bounds.Max.X-cellMargin*2)/grid.MaxCells, (bounds.Max.Y-cellMargin*2)/grid.MaxCells, shapeMargin))
		return g.CheckWin(context)
	}

	return false
}
