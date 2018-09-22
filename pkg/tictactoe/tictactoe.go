package tictactoe

import (
	"fmt"
	"os"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/juanvallejo/tictactoe/pkg/tictactoe/grid"
	"github.com/juanvallejo/tictactoe/pkg/tictactoe/shape"
)

const (
	winWidth  = 800
	winHeight = 600

	cellMargin  = 45
	shapeMargin = 25

	winTextSize = 4
)

var winBgcolor = colornames.Darkslategrey
var winTextAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

func NewGame() {
	config := pixelgl.WindowConfig{
		Title:  "Tic Tac Toe",
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
	textContext := text.New(pixel.V(bounds.Max.X/2, bounds.Max.Y/2), winTextAtlas)

	g := grid.NewGrid(pixel.V(0, 0), bounds.Max.X, bounds.Max.Y, grid.MaxCells, cellMargin)

	for !window.Closed() {
		context.Clear()
		textContext.Clear()

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			gameWin = handleMouseClick(window, context, textContext, g, shapeDecider, bounds, gameWin)
		}

		g.Render(context)
		context.Draw(window)
		textContext.Draw(window, pixel.IM.Scaled(textContext.Orig, winTextSize))
		window.Update()
	}
}

func handleMouseClick(window *pixelgl.Window, context *imdraw.IMDraw, textContext *text.Text, g grid.Grid, shapeDecider *shape.ShapeDecider, bounds pixel.Rect, gameWin bool) bool {
	if gameWin {
		window.Clear(winBgcolor)
		g.Reset()
		return false
	}

	if cell := g.AtVector(window.MousePosition()); cell != nil {
		cell.Set(shape.NewShape(pixel.V(cell.Start().X, cell.Start().Y), shapeDecider.Next(), (bounds.Max.X-cellMargin*2)/grid.MaxCells, (bounds.Max.Y-cellMargin*2)/grid.MaxCells, shapeMargin))
		return g.CheckWin(context, textContext)
	}

	return false
}
